//go:build integration

// Real Postgres HTTP checks prove email login, member continuity, session end, member isolation, and closed synthetic signup.
package relay

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type inbox struct {
	sync.Mutex
	codes map[string]string
	fail  bool
}

func (m *inbox) send(_ context.Context, to, _, body string) error {
	m.Lock()
	defer m.Unlock()
	if m.fail {
		return errors.New("provider down")
	}
	m.codes[to] = regexp.MustCompile(`\d{6}`).FindString(body)
	return nil
}
func (m *inbox) code(to string) string {
	m.Lock()
	defer m.Unlock()
	return m.codes[to]
}

type browser struct {
	t       *testing.T
	h       http.Handler
	ip      string
	cookies map[string]string
}

func (b *browser) do(method, path string, form url.Values, header ...string) *httptest.ResponseRecorder {
	b.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
	req.RemoteAddr = b.ip + ":40000"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	for name, value := range b.cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	w := httptest.NewRecorder()
	b.h.ServeHTTP(w, req)
	for _, c := range w.Result().Cookies() {
		if c.MaxAge < 0 {
			delete(b.cookies, c.Name)
			continue
		}
		if !c.Secure || !c.HttpOnly || c.SameSite != http.SameSiteStrictMode || c.Path != "/" || !strings.HasPrefix(c.Name, "__Host-") {
			b.t.Fatalf("weak cookie %v", c)
		}
		b.cookies[c.Name] = c.Value
	}
	return w
}
func expect(t *testing.T, w *httptest.ResponseRecorder, status int, contains ...string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("got %d want %d: %s", w.Code, status, w.Body)
	}
	for _, s := range contains {
		if !strings.Contains(w.Body.String()+w.Header().Get("Location"), s) {
			t.Fatalf("missing %q in %s %s", s, w.Header().Get("Location"), w.Body)
		}
	}
}

func identityService(t *testing.T, pool *pgxpool.Pool) (*Service, *inbox) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `UPDATE relay_state SET data='{}',epoch=epoch+1,clock=clock_timestamp()`); err != nil {
		t.Fatal(err)
	}
	mail := &inbox{codes: map[string]string{}}
	return &Service{Pool: pool, Mail: mail.send}, mail
}
func (s *Service) mutateState(t *testing.T, fn func(*State)) {
	t.Helper()
	if _, err := s.transaction(context.Background(), func(st *State, _ time.Time) (any, error) { fn(st); return nil, nil }); err != nil {
		t.Fatal(err)
	}
}

// login runs the real start/verify flow and returns the member ID shown on the home page.
func login(t *testing.T, b *browser, mail *inbox, email string) string {
	t.Helper()
	expect(t, b.do("POST", "/auth/start", url.Values{"email": {email}}), 303, "/auth/verify")
	expect(t, b.do("GET", "/auth/verify", nil), 200, maskEmail(strings.ToLower(email)), "아직 로그인되지 않았습니다")
	expect(t, b.do("POST", "/auth/verify", url.Values{"code": {mail.code(strings.ToLower(email))}}), 303, "/home")
	w := b.do("GET", "/home", nil)
	expect(t, w, 200, maskEmail(strings.ToLower(email)), "회원 식별자")
	return regexp.MustCompile(`mem_[A-Za-z0-9_-]+`).FindString(w.Body.String())
}

func TestEmailIdentity(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	newBrowser := func(s *Service, ip string) *browser {
		return &browser{t, s.Handler(), ip, map[string]string{}}
	}

	t.Run("verified_signup_home_and_continuity", func(t *testing.T) {
		s, mail := identityService(t, pool)
		b := newBrowser(s, "192.0.2.1")
		expect(t, b.do("GET", "/", nil), 200, "이메일로 가입", "관리자 아이디, 서버 접속, 별도 초대는 필요하지 않습니다")
		expect(t, b.do("GET", "/home", nil), 303, "/?n=expired")
		expect(t, b.do("POST", "/auth/start", url.Values{"email": {"Name <x@y.co>"}}), 422, "형식")
		expect(t, b.do("POST", "/auth/start", url.Values{"email": {"Member.One@Example.com"}}), 303)
		s.mutateState(t, func(st *State) {
			if len(st.Owners) != 0 || len(st.Members) != 0 {
				t.Fatal("owner or member issued before verification")
			}
		})
		code := mail.code("member.one@example.com")
		wrong := "000000"
		if code == wrong {
			wrong = "111111"
		}
		expect(t, b.do("POST", "/auth/verify", url.Values{"code": {wrong}}), 401, "남은 시도: 4회")
		expect(t, b.do("GET", "/home", nil), 303, "/?n=expired")
		expect(t, b.do("POST", "/auth/verify", url.Values{"code": {code}}), 303, "/home")
		w := b.do("GET", "/home", nil)
		expect(t, w, 200, "m***@example.com", "아직 연결한 agent가 없습니다", "상태: 준비 중", "agent의 키와 자격은 철회하지 않습니다")
		if strings.Contains(w.Body.String(), "member.one") || strings.Contains(w.Body.String(), "/v1/owners") || strings.Contains(w.Body.String(), "Basic") {
			t.Fatal("home leaked address or owner-admin path")
		}
		first := regexp.MustCompile(`mem_[A-Za-z0-9_-]+`).FindString(w.Body.String())
		expect(t, b.do("POST", "/auth/verify", url.Values{"code": {code}}), 401, "만료")
		// A new browser re-login reaches the same member; the spacing limit is reset to model a later login.
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		other := newBrowser(s, "192.0.2.2")
		if login(t, other, mail, "member.one@EXAMPLE.com") != first {
			t.Fatal("re-login changed member")
		}
		// Restart: a new process on the same database keeps the session.
		restarted := newBrowser(&Service{Pool: pool, Mail: mail.send}, "192.0.2.2")
		restarted.cookies = other.cookies
		expect(t, restarted.do("GET", "/home", nil), 200, first)
		// Synthetic signup is closed unless explicitly enabled; a session cookie is not an owner or agent credential.
		expect(t, b.do("POST", "/v1/owners", nil), 403)
		expect(t, b.do("POST", "/v1/agents", nil), 401)
		expect(t, b.do("GET", "/owner", nil), 401)
	})

	t.Run("concurrent_first_signup_and_one_use_code", func(t *testing.T) {
		s, _ := identityService(t, pool)
		now := time.Now().UTC()
		s.mutateState(t, func(st *State) {
			for _, p := range []string{"p1", "p2"} {
				st.Challenges[hashToken(p)] = &Challenge{"race@example.com", hashToken(p + ":123456"), now.Add(codeTTL), 0}
			}
			st.Challenges[hashToken("p3")] = &Challenge{"once@example.com", hashToken("p3:654321"), now.Add(codeTTL), 0}
		})
		var wg sync.WaitGroup
		codes := make(chan int, 10)
		for i, p := range []string{"p1", "p2", "p3", "p3", "p3", "p3", "p3", "p3"} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				b := &browser{t, s.Handler(), "198.51.100." + string(rune('1'+i)), map[string]string{pendingCookie: p}}
				code := "123456"
				if p == "p3" {
					code = "654321"
				}
				w := b.do("POST", "/auth/verify", url.Values{"code": {code}})
				if p != "p3" && w.Code != 303 {
					t.Error("race signup", w.Code)
				}
				if p == "p3" {
					codes <- w.Code
				}
			}()
		}
		wg.Wait()
		close(codes)
		wins := 0
		for c := range codes {
			if c == 303 {
				wins++
			}
		}
		s.mutateState(t, func(st *State) {
			if wins != 1 || len(st.Members) != 2 || len(st.Owners) != 2 || len(st.Identities) != 2 {
				t.Fatalf("wins %d members %d owners %d", wins, len(st.Members), len(st.Owners))
			}
		})
	})

	t.Run("logout_all_reauth_expiry_and_agent_credential_survives", func(t *testing.T) {
		s, mail := identityService(t, pool)
		one, two := newBrowser(s, "203.0.113.1"), newBrowser(s, "203.0.113.2")
		member := login(t, one, mail, "owner@example.com")
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		login(t, two, mail, "owner@example.com")
		s.mutateState(t, func(st *State) {
			st.Agents["agent_kept"] = &Agent{st.Members[member].Owner, hashToken("agent-secret"), map[string]*Key{"key1": {Public: make([]byte, 32)}}}
		})
		old := map[string]string{}
		for k, v := range one.cookies {
			old[k] = v
		}
		expect(t, one.do("POST", "/auth/logout", nil), 303, "/?n=logout")
		expect(t, one.do("GET", "/?n=logout", nil), 200, "연결한 agent의 키와 자격은 그대로입니다")
		replay := newBrowser(s, "203.0.113.1")
		replay.cookies = old
		expect(t, replay.do("GET", "/home", nil), 303, "/?n=expired")
		agent := httptest.NewRequest("GET", "/v1/contacts", nil)
		agent.Header.Set("Authorization", "Bearer agent-secret")
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, agent)
		expect(t, w, 200)
		// Logout-all needs an email check within five minutes.
		s.mutateState(t, func(st *State) {
			for _, session := range st.Sessions {
				session.Verified = session.Verified.Add(-reauthWindow)
			}
			st.Rates = map[string][]time.Time{}
		})
		expect(t, two.do("GET", "/home", nil), 200, "이메일 다시 확인", "<code>agent_kept</code>")
		expect(t, two.do("POST", "/auth/logout-all", nil), 403, "5분 안에")
		third := newBrowser(s, "203.0.113.3")
		login(t, third, mail, "owner@example.com")
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		expect(t, two.do("POST", "/auth/reauth", nil), 303, "/auth/verify")
		expect(t, two.do("POST", "/auth/verify", url.Values{"code": {mail.code("owner@example.com")}}), 303, "/home")
		expect(t, two.do("GET", "/home", nil), 200, "모든 브라우저에서 로그아웃")
		expect(t, two.do("POST", "/auth/logout-all", nil), 303, "/?n=all")
		expect(t, third.do("GET", "/home", nil), 303, "/?n=expired")
		// Idle expiry.
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		login(t, one, mail, "owner@example.com")
		s.mutateState(t, func(st *State) {
			for _, session := range st.Sessions {
				session.Seen = session.Seen.Add(-sessionIdle)
			}
		})
		expect(t, one.do("GET", "/home", nil), 303, "/?n=expired")
		w = httptest.NewRecorder()
		s.Handler().ServeHTTP(w, agent)
		expect(t, w, 200)
	})

	t.Run("send_failure_limits_and_csrf", func(t *testing.T) {
		s, mail := identityService(t, pool)
		b := newBrowser(s, "192.0.2.50")
		mail.fail = true
		expect(t, b.do("POST", "/auth/start", url.Values{"email": {"fail@example.com"}}), 503, "보내지 못했습니다", "로그인되지 않았습니다")
		s.mutateState(t, func(st *State) {
			if len(st.Challenges) != 0 {
				t.Fatal("unsent code retained")
			}
		})
		mail.fail = false
		expect(t, b.do("POST", "/auth/start", url.Values{"email": {"fail@example.com"}}), 429, "이후 다시 시도")
		unconfigured := newBrowser(&Service{Pool: pool}, "192.0.2.51")
		expect(t, unconfigured.do("POST", "/auth/start", url.Values{"email": {"none@example.com"}}), 503)
		// Anonymous source budget counts refused requests: 30 per rolling minute.
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		flood := newBrowser(s, "192.0.2.60")
		for i := 0; i < 30; i++ {
			flood.do("POST", "/auth/verify", url.Values{"code": {"000000"}})
		}
		expect(t, flood.do("POST", "/auth/verify", url.Values{"code": {"000000"}}), 429)
		expect(t, newBrowser(s, "192.0.2.61").do("POST", "/auth/start", url.Values{"email": {"ok@example.com"}}), 303)
		// The edge client-IP header is honored only when configured.
		s.ClientIPHeader = "CF-Connecting-IP"
		expect(t, flood.do("POST", "/auth/start", url.Values{"email": {"edge@example.com"}}, "CF-Connecting-IP", "192.0.2.99"), 303)
		s.ClientIPHeader = ""
		expect(t, flood.do("POST", "/auth/start", url.Values{"email": {"edge2@example.com"}}, "CF-Connecting-IP", "192.0.2.98"), 429)
		// Cross-site browser writes are refused before any state change.
		cross := newBrowser(s, "192.0.2.70")
		expect(t, cross.do("POST", "/auth/start", url.Values{"email": {"cross@example.com"}}, "Sec-Fetch-Site", "cross-site"), 403)
		if mail.code("cross@example.com") != "" {
			t.Fatal("cross-site request sent mail")
		}
	})

	t.Run("refused_principal_does_not_spend_shared_rate", func(t *testing.T) {
		s, mail := identityService(t, pool)
		alice, bob := newBrowser(s, "192.0.2.90"), newBrowser(s, "192.0.2.91")
		login(t, alice, mail, "alice@example.com")
		login(t, bob, mail, "bob@example.com")
		s.mutateState(t, func(st *State) {
			st.Rates = map[string][]time.Time{}
			for _, session := range st.Sessions {
				session.Verified = session.Verified.Add(-reauthWindow)
			}
		})
		// Each flood exceeds the shared budget on its own: 250 > 200 new and 150 > 100 cleanup.
		flood := newBrowser(s, "198.51.100.7")
		for i := 0; i < 250; i++ {
			flood.do("POST", "/auth/verify", url.Values{"code": {"000000"}})
			alice.do("GET", "/home", nil)
		}
		for i := 0; i < 150; i++ {
			alice.do("POST", "/auth/logout-all", nil)
		}
		expect(t, flood.do("POST", "/auth/verify", url.Values{"code": {"000000"}}), 429)
		expect(t, alice.do("GET", "/home", nil), 429)
		expect(t, alice.do("POST", "/auth/logout-all", nil), 429)
		// A new process on the same database keeps the counts; other principals still log in, load home, and log out.
		restarted := &Service{Pool: pool, Mail: mail.send}
		flood.h, alice.h, bob.h = restarted.Handler(), restarted.Handler(), restarted.Handler()
		expect(t, flood.do("POST", "/auth/verify", url.Values{"code": {"000000"}}), 429)
		expect(t, alice.do("GET", "/home", nil), 429)
		expect(t, newBrowser(restarted, "203.0.113.50").do("POST", "/auth/start", url.Values{"email": {"carol@example.com"}}), 303, "/auth/verify")
		expect(t, bob.do("GET", "/home", nil), 200, "회원 식별자")
		expect(t, bob.do("POST", "/auth/logout", nil), 303, "/?n=logout")
	})

	t.Run("member_gate_isolation_and_decision", func(t *testing.T) {
		f := setup(t, pool)
		mail := &inbox{codes: map[string]string{}}
		f.s.Mail = mail.send
		f.handler = f.s.Handler()
		alice, bob := &browser{t, f.handler, "192.0.2.80", map[string]string{}}, &browser{t, f.handler, "192.0.2.81", map[string]string{}}
		bobMember := login(t, bob, mail, "bob@example.com")
		login(t, alice, mail, "alice@example.com")
		// Bind agent_b to Bob's member owner; the frozen gate flow then targets Bob's home.
		f.mutate(func(st *State) { st.Agents["agent_b"].Owner = st.Members[bobMember].Owner })
		receipt := f.send(f.message(firstID, "idempotency-key-01"), "agent_a", "", 200)
		claim := f.deliver(firstID)
		h := wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "approval-key-0001", firstID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute))
		f.send(h, "agent_b", claim, 200)
		expect(t, bob.do("GET", "/home", nil), 200, "/home/gates/"+gateID, "<code>agent_b</code>")
		expect(t, alice.do("GET", "/home/gates/"+gateID, nil), 403)
		expect(t, alice.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrf(alice.cookies[sessionCookie], gateID)}, "decision": {"approve"}}), 409)
		expect(t, bob.do("GET", "/home/gates/"+gateID, nil), 200, "granularity_min", "Approve 승인", `href="/home"`)
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {"bad"}, "decision": {"approve"}}), 403)
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrf(bob.cookies[sessionCookie], gateID)}, "decision": {"approve"}}, "Sec-Fetch-Site", "cross-site"), 403)
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrf(bob.cookies[sessionCookie], gateID)}, "decision": {"approve"}}), 303)
		result := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 200)
		if result["executable"] != false || result["disclosure"] != false {
			t.Fatal("gate enabled effect")
		}
		expect(t, bob.do("POST", "/auth/logout", nil), 303)
		expect(t, bob.do("GET", "/home/gates/"+gateID, nil), 303, "/?n=expired")
	})
}
