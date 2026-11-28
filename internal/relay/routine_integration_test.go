//go:build integration

// Isolated Postgres and Redis: an owner registers a routine through the web, then one arrival makes exactly one routine POST.
package relay

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRoutineWakePostgres(t *testing.T) {
	pool := messagePool(t)
	s, mail := identityService(t, pool)
	s.Google, _ = googleFixture(t)
	// RemoteEvents runs only with an HTTPS remote origin, as in production.
	s.Google.config.RedirectURL = "https://example.com/auth/google/callback"
	var err error
	if s.Sealer, err = NewSealer(base64.StdEncoding.EncodeToString(make([]byte, 32))); err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	posts := []*http.Request{}
	bodies := []string{}
	fake := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		posts, bodies = append(posts, r), append(bodies, string(body))
		mu.Unlock()
		w.Write([]byte(`{"runUuid":"fake"}`))
	}))
	defer fake.Close()
	// The stored URL stays the real routine host; only this test dials it to the fake server.
	routineClient = func() *http.Client {
		return &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, network, fake.Listener.Addr().String())
			},
		}}
	}
	defer func() { routineClient = remoteHTTPClient }()

	h := s.Handler()
	b := &browser{t, h, "192.0.2.41", map[string]string{}}
	email := "routine@humanpeer.invalid"
	expect(t, b.do("POST", "/auth/start", url.Values{"email": {email}}), 303)
	expect(t, b.do("POST", "/auth/verify", url.Values{"code": {mail.code(email)}}), 303)
	expect(t, b.do("POST", "/home/profile", url.Values{"name": {"합성 routine 주인"}}), 303)
	var agent string
	var senderKey ed25519.PrivateKey
	s.mutateState(t, func(st *State) {
		now := time.Now()
		id, _, e := st.session(b.cookies[sessionCookie], now)
		if e != nil {
			t.Fatal(e)
		}
		_, install, _, _ := remoteFixture(t, st, st.Members[id].Owner, now)
		agent = install.Agent
		pub, key, _ := ed25519.GenerateKey(rand.Reader)
		senderKey = key
		st.Owners["routine_sender"] = &Owner{Active: true}
		st.Members["routine_sender"] = &Member{Owner: "routine_sender", Active: true}
		st.Agents["agent_routine_sender"] = &Agent{Owner: "routine_sender", Keys: map[string]*Key{"key1": {Public: pub, Credential: hashToken("routine_sender")}}}
		st.Pairs[pairID("agent_routine_sender", agent)] = &Pair{A: "agent_routine_sender", B: agent, State: "active", Generation: 1}
	})
	expect(t, b.do("GET", "/home/agents", nil), 200, `id="routines"`, "등록된 도착 알림이 없습니다", "knowslink_receive를 빈 결과가 나올 때까지")

	form := url.Values{"agent": {agent}, "url": {routineTestURL}, "key": {routineTestKey}, "consent": {"yes"}}
	stale := func(age time.Duration) {
		s.mutateState(t, func(st *State) {
			_, session, e := st.session(b.cookies[sessionCookie], time.Now())
			if e != nil {
				t.Fatal(e)
			}
			session.Verified = time.Now().Add(-age)
		})
	}
	stale(6 * time.Minute)
	refused := b.do("POST", "/home/routine-set", form)
	expect(t, refused, 403, `role="alert"`, "5분 안의 로그인 신원 재확인", `action="/auth/reauth"`, `id="routines"`)
	if strings.Contains(refused.Body.String(), routineTestKey) {
		t.Fatal("refusal echoed the key")
	}
	expect(t, b.do("POST", "/home/routine-remove", url.Values{"agent": {agent}}), 403, `role="alert"`)
	s.mutateState(t, func(st *State) {
		if len(st.Remote.Routines) != 0 {
			t.Fatal("stale session stored a routine")
		}
	})
	stale(0)
	bad := url.Values{"agent": {agent}, "url": {"https://127.0.0.1/automations/webhook/x"}, "key": {routineTestKey}, "consent": {"yes"}}
	expect(t, b.do("POST", "/home/routine-set", bad), 422, `role="alert"`, "routine webhook URL 형식이 아닙니다")
	expect(t, b.do("POST", "/home/routine-set", form), 303, "/home/agents?n=routine#routines")
	page := b.do("GET", "/home/agents?n=routine", nil)
	expect(t, page, 200, "도착 알림 설정을 바꿨습니다", "URL 끝 <code>-123</code>", "정상")
	if strings.Contains(page.Body.String(), routineTestKey) || strings.Contains(page.Body.String(), "abc_DEF") {
		t.Fatal("page showed the key or full URL")
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { s.RemoteEvents(ctx); close(done) }()
	defer func() { cancel(); <-done }()
	s.mutateState(t, func(st *State) {
		now := time.Now()
		text := signedText(firstID, "agent_routine_sender", agent, "routine-integration", "", "secret body from sender", now.Add(170*time.Second), senderKey)
		if _, e := st.ingestPublicText("agent_routine_sender", text, now); e != nil {
			t.Fatal(e)
		}
	})
	s.remoteNotify(context.Background())
	count := func() int { mu.Lock(); defer mu.Unlock(); return len(posts) }
	for deadline := time.Now().Add(10 * time.Second); count() == 0 && time.Now().Before(deadline); {
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(1500 * time.Millisecond)
	if count() != 1 {
		t.Fatalf("routine POSTs: %d", count())
	}
	r := posts[0]
	if r.Method != "POST" || r.URL.Path != "/automations/webhook/abc_DEF-123" || r.Host != "api2.cursor.sh" || r.Header.Get("Authorization") != "Bearer "+routineTestKey || r.Header.Get("X-Automation-Key") != routineTestKey {
		t.Fatal("request", r.Method, r.Host, r.URL, r.Header)
	}
	var body map[string]any
	if json.Unmarshal([]byte(bodies[0]), &body) != nil || body["event"] != "knowslink_text_arrived" || body["pending"] != 1.0 || len(body) != 3 || strings.Contains(bodies[0], "secret body") || strings.Contains(bodies[0], "agent_routine_sender") {
		t.Fatal("body", bodies[0])
	}
	expect(t, b.do("GET", "/home/agents", nil), 200, "마지막 성공")
	expect(t, b.do("POST", "/home/routine-remove", url.Values{"agent": {agent}}), 303)
	expect(t, b.do("GET", "/home/agents", nil), 200, "등록된 도착 알림이 없습니다")
}
