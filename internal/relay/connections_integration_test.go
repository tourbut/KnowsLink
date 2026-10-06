//go:build integration

// Public-member HTTP checks exercise confirmed onboarding, cross-account denial, races and revocation on actual Postgres.
package relay

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func publicCall(t *testing.T, h http.Handler, path, token string, body any, status int) map[string]any {
	t.Helper()
	raw, _ := json.Marshal(body)
	r := httptest.NewRequest("POST", path, strings.NewReader(string(raw)))
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	expect(t, w, status)
	out := map[string]any{}
	if w.Body.String() != "null\n" {
		_ = json.Unmarshal(w.Body.Bytes(), &out)
	}
	return out
}
func browserAgent(t *testing.T, b *browser) string {
	t.Helper()
	expect(t, b.do("POST", "/home/agents", nil), 303)
	w := b.do("GET", "/home", nil)
	expect(t, w, 200, "미연결")
	ids := regexp.MustCompile(`<code>(agent_[a-f0-9]+)</code>`).FindAllStringSubmatch(w.Body.String(), -1)
	return ids[len(ids)-1][1]
}
func browserGrant(t *testing.T, b *browser, agent, mode string) (string, string) {
	t.Helper()
	w := b.do("POST", "/home/connect", url.Values{"agent": {agent}, "client": {supportedClient}, "mode": {mode}})
	expect(t, w, 200, "연결 대기", "개인키")
	token := regexp.MustCompile(`readonly value="([A-Za-z0-9_-]{43})"`).FindStringSubmatch(w.Body.String())[1]
	return token, hashToken(token)
}
func clientPrepare(t *testing.T, h http.Handler, token string, kid string) (ed25519.PrivateKey, string) {
	t.Helper()
	info := publicCall(t, h, "/v1/connect/info", "", map[string]string{"token": token, "client": supportedClient}, 200)
	public, private, _ := ed25519.GenerateKey(rand.Reader)
	c := &Connection{Owner: info["owner"].(string), Agent: info["agent"].(string), Client: supportedClient, Mode: info["mode"].(string), Kid: kid, Public: base64.RawURLEncoding.EncodeToString(public)}
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c)))
	publicCall(t, h, "/v1/connect/prepare", "", map[string]string{"token": token, "client": supportedClient, "kid": kid, "public": c.Public, "proof": proof}, 200)
	return private, proof
}
func TestPublicAgentHTTP(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, e := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	newBrowser := func(s *Service, ip string) *browser { return &browser{t, s.Handler(), ip, map[string]string{}} }
	t.Run("onboarding_owner_isolation_rotation_and_all_revoked_paths", func(t *testing.T) {
		s, mail := identityService(t, pool)
		a, b := newBrowser(s, "192.0.2.1"), newBrowser(s, "192.0.2.2")
		member := login(t, a, mail, "alice@example.com")
		login(t, b, mail, "bob@example.com")
		agent := browserAgent(t, a)
		expect(t, b.do("POST", "/home/connect", url.Values{"agent": {agent}, "client": {supportedClient}, "mode": {"register"}}), 403)
		expect(t, a.do("POST", "/home/connect", url.Values{"agent": {agent}, "client": {"unknown"}, "mode": {"register"}}), 422, "미지원")
		token, id := browserGrant(t, a, agent, "register")
		private, proof := clientPrepare(t, s.Handler(), token, "key1")
		expect(t, a.do("GET", "/home/connections/"+id, nil), 200, "공개키 확인", "SHA256:")
		expect(t, b.do("GET", "/home/connections/"+id, nil), 403)
		expect(t, b.do("POST", "/home/confirm", url.Values{"connection": {id}}), 401)
		publicCall(t, s.Handler(), "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 403)
		expect(t, a.do("POST", "/home/confirm", url.Values{"connection": {id}}, "Sec-Fetch-Site", "cross-site"), 403)
		expect(t, a.do("POST", "/home/confirm", url.Values{"connection": {id}}), 303)
		connected := publicCall(t, s.Handler(), "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 200)
		old := connected["credential"].(string)
		publicCall(t, s.Handler(), "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 401)
		publicCall(t, s.Handler(), "/v1/pull", old, map[string]string{}, 200)
		publicCall(t, s.Handler(), "/v1/owner-revoke", old, map[string]string{}, 401)
		expect(t, a.do("GET", "/home", nil), 200, "연결 완료", "활성")
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		token, id = browserGrant(t, a, agent, "rotate")
		_, proof = clientPrepare(t, s.Handler(), token, "key2")
		expect(t, a.do("POST", "/home/confirm", url.Values{"connection": {id}}), 303)
		connected = publicCall(t, s.Handler(), "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 200)
		fresh := connected["credential"].(string)
		// New service instance on the same database retains the key-specific denial.
		restarted := (&Service{Pool: pool}).Handler()
		for _, path := range []string{"pull", "persist", "ack", "claim", "authorize", "gate-consume"} {
			publicCall(t, restarted, "/v1/"+path, old, map[string]string{}, 401)
		}
		for _, intent := range []string{"schedule.query", "relay.result"} {
			body := queryBody()
			reply := ""
			if intent == "relay.result" {
				body = map[string]any{"status": "denied"}
				reply = firstID
			}
			message := wire(t, private, firstID, agent, "agent_other", intent, "public-revoke-test-01", reply, body, time.Now().Add(time.Minute))
			publicCall(t, restarted, "/v1/send", old, json.RawMessage(message), 401)
		}
		publicCall(t, restarted, "/v1/pull", fresh, map[string]string{}, 200)
		expect(t, a.do("POST", "/home/key-revoke", url.Values{"agent": {agent}, "kid": {"key2"}}), 303)
		publicCall(t, restarted, "/v1/pull", fresh, map[string]string{}, 401)
		expect(t, a.do("GET", "/home", nil), 200, "철회", "미연결")
		s.mutateState(t, func(st *State) {
			if st.Agents[agent].Owner != st.Members[member].Owner {
				t.Fatal("owner changed")
			}
			st.Rates = map[string][]time.Time{}
			for _, ss := range st.Sessions {
				ss.Verified = ss.Verified.Add(-reauthWindow)
			}
		})
		expect(t, a.do("POST", "/home/connect", url.Values{"agent": {agent}, "client": {supportedClient}, "mode": {"register"}}), 422, "5분")
	})
	t.Run("same_owner_first_accept_generation_pending_expiry_and_cleanup_budget", func(t *testing.T) {
		s, mail := identityService(t, pool)
		a := newBrowser(s, "198.51.100.1")
		member := login(t, a, mail, "same@example.com")
		one := browserAgent(t, a)
		expect(t, a.do("POST", "/home/agents", nil), 303)
		var two string
		s.mutateState(t, func(st *State) {
			for id := range st.Agents {
				if id != one {
					two = id
				}
			}
			st.Rates = map[string][]time.Time{}
		})
		form := url.Values{"agent": {one}, "target": {two}, "generation": {"1"}}
		expect(t, a.do("POST", "/home/invites", form), 303, "/home?n=invited")
		var exp time.Time
		s.mutateState(t, func(st *State) { exp = st.Pairs[pairID(one, two)].Exp })
		// Repeats from either side keep the pending invite and say so.
		expect(t, a.do("POST", "/home/invites", form), 303, "/home?n=invite-pending")
		expect(t, a.do("POST", "/home/invites", url.Values{"agent": {two}, "target": {one}}), 303, "/home?n=invite-pending")
		expect(t, a.do("GET", "/home?n=invite-pending", nil), 200, "이미 대기 중인 초대", "기한도 그대로", "받은 초대입니다")
		s.mutateState(t, func(st *State) {
			p := st.Pairs[pairID(one, two)]
			if p.State != "pending" || p.Generation != 1 || !p.Exp.Equal(exp) || len(st.Pairs) != 1 {
				t.Fatal("implicit accept or repeat changed the invite")
			}
		})
		decision := url.Values{"agent": {one}, "target": {two}, "decision": {"accept"}, "generation": {"1"}}
		var wg sync.WaitGroup
		results := make(chan int, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				b := &browser{t, s.Handler(), "198.51.100.2", map[string]string{sessionCookie: a.cookies[sessionCookie]}}
				results <- b.do("POST", "/home/invite-decision", decision).Code
			}()
		}
		wg.Wait()
		close(results)
		for c := range results {
			if c != 303 {
				t.Fatal(c)
			}
		}
		s.mutateState(t, func(st *State) {
			if len(st.Pairs) != 1 || st.Pairs[pairID(one, two)].Generation != 1 {
				t.Fatal("duplicate pair")
			}
			st.Rates["http:new"] = make([]time.Time, 201)
			for i := range st.Rates["http:new"] {
				st.Rates["http:new"][i] = time.Now()
			}
		})
		expect(t, a.do("POST", "/home/invites", form), 429)
		expect(t, a.do("POST", "/home/unpair", form), 303)
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		expect(t, a.do("GET", "/home", nil), 200, "연결이 끝났습니다", "새 초대 보내기")
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		expect(t, a.do("POST", "/home/invites", form), 303)
		expect(t, a.do("POST", "/home/invite-decision", decision), 403)
		decision.Set("generation", "2")
		expect(t, a.do("POST", "/home/invite-decision", decision), 303)
		expect(t, a.do("POST", "/home/invites", form), 303, "/home?n=invite-active")
		s.mutateState(t, func(st *State) {
			if st.Pairs[pairID(one, two)].Generation != 2 || st.Pairs[pairID(one, two)].State != "active" {
				t.Fatal("old generation reused")
			}
			st.Rates = map[string][]time.Time{}
			st.Pairs[pairID(one, two)].State = "pending"
			st.Pairs[pairID(one, two)].Exp = time.Now().Add(-time.Second)
		})
		expect(t, a.do("POST", "/home/invite-decision", decision), 403)
		expect(t, a.do("POST", "/home/invites", form), 303)
		s.mutateState(t, func(st *State) {
			if st.Pairs[pairID(one, two)].Generation != 3 {
				t.Fatal("expiry reused generation")
			}
			for i := 0; i < 3; i++ {
				st.Agents[fmt.Sprint("agent_extra", i)] = &Agent{Owner: st.Members[member].Owner, Keys: map[string]*Key{}}
			}
		})
		expect(t, a.do("POST", "/home/agents", nil), 409)
	})
	t.Run("malformed_connection_requests_spend_anonymous_budget", func(t *testing.T) {
		s, _ := identityService(t, pool)
		for i := 0; i < 31; i++ {
			r := httptest.NewRequest("POST", "/v1/connect/prepare", strings.NewReader("{"))
			w := httptest.NewRecorder()
			s.Handler().ServeHTTP(w, r)
			expected := 422
			if i == 30 {
				expected = 429
			}
			expect(t, w, expected)
			if i == 30 {
				// The connect API reports the actual retry instant like the /v1 rate limiter, not a fixed 60s.
				var body map[string]string
				_ = json.Unmarshal(w.Body.Bytes(), &body)
				retryAt, e1 := time.Parse(time.RFC3339, body["retry_at"])
				header, e2 := http.ParseTime(w.Header().Get("Retry-After"))
				if body["error"] != "rate_limited" || e1 != nil || e2 != nil || !header.Equal(retryAt) || !retryAt.After(time.Now()) || retryAt.After(time.Now().Add(time.Minute+time.Second)) {
					t.Fatal("connect 429 retry", body, w.Header().Get("Retry-After"))
				}
			}
		}
	})
	t.Run("invalid_session_pages_spend_anonymous_budget", func(t *testing.T) {
		s, _ := identityService(t, pool)
		stale := &browser{t, s.Handler(), "192.0.2.77", map[string]string{sessionCookie: "stale-session"}}
		for i := 0; i < 30; i++ {
			path := "/home/connections/unknown"
			if i%2 == 1 {
				path = "/home"
			}
			expect(t, stale.do("GET", path, nil), 303, "/?n=expired")
			stale.cookies[sessionCookie] = "stale-session"
		}
		// A restarted relay reads the committed budget: invalid-session GETs were not rolled back.
		stale.h = (&Service{Pool: pool}).Handler()
		for _, path := range []string{"/home/connections/unknown", "/home"} {
			expect(t, stale.do("GET", path, nil), 429, "이후 다시 시도")
		}
		expect(t, stale.do("POST", "/auth/logout", nil), 429)
		expect(t, stale.do("POST", "/auth/reauth", nil), 429)
	})
	t.Run("revoked_record_saturation_concurrency_restart_and_retention", func(t *testing.T) {
		s, mail := identityService(t, pool)
		a := newBrowser(s, "192.0.2.66")
		member := login(t, a, mail, "churn@example.com")
		ids := func() (live, revoked []string) {
			s.mutateState(t, func(st *State) {
				for id, ag := range st.Agents {
					if ag.Revoked {
						revoked = append(revoked, id)
					} else {
						live = append(live, id)
					}
				}
				st.Rates = map[string][]time.Time{}
			})
			return
		}
		for round := 0; round < 2; round++ {
			n := 5 - 4*round
			for i := 0; i < n; i++ {
				expect(t, a.do("POST", "/home/agents", nil), 303)
			}
			live, _ := ids()
			for _, id := range live {
				expect(t, a.do("POST", "/home/agent-revoke", url.Values{"agent": {id}}), 303)
			}
		}
		// Six revoked records and no active agent: the record cap, not the active cap of 5, admits exactly 4.
		var wg sync.WaitGroup
		codes := make(chan int, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				b := &browser{t, s.Handler(), "192.0.2.66", map[string]string{sessionCookie: a.cookies[sessionCookie]}}
				codes <- b.do("POST", "/home/agents", nil).Code
			}()
		}
		wg.Wait()
		close(codes)
		wins := 0
		for c := range codes {
			if c == 303 {
				wins++
			} else if c != 409 {
				t.Fatal(c)
			}
		}
		live, revoked := ids()
		if wins != ownerAgentRecords-6 || len(live) != 4 || len(revoked) != 6 {
			t.Fatal("record cap race", wins, len(live), len(revoked))
		}
		a.h = (&Service{Pool: pool}).Handler()
		expect(t, a.do("POST", "/home/agents", nil), 409, "철회한 agent 기록을 보존하는 중", "최소 24시간", "자기 홈으로 돌아가기")
		// Saturation never blocks cleanup, and revoking adds no record.
		expect(t, a.do("POST", "/home/agent-revoke", url.Values{"agent": {live[0]}}), 303)
		expect(t, a.do("POST", "/home/agent-revoke", url.Values{"agent": {revoked[0]}}), 303)
		s.mutateState(t, func(st *State) {
			owner := st.Members[member].Owner
			if len(st.Agents) != ownerAgentRecords || st.agentCapacity(owner) {
				t.Fatal("cleanup changed records", len(st.Agents))
			}
			for _, ag := range st.Agents {
				if ag.Revoked && ag.Changed.IsZero() {
					t.Fatal("revocation time missing")
				}
			}
			// The first six revocations pass the 24h retention; the latest one stays.
			for _, id := range revoked {
				st.Agents[id].Changed = st.Agents[id].Changed.Add(-revokedRetention)
			}
			st.Rates = map[string][]time.Time{}
		})
		expect(t, a.do("POST", "/home/agents", nil), 303)
		expect(t, a.do("GET", "/home", nil), 200, "철회", "최소 24시간 보존", "목록에서 사라질 수 있습니다")
		s.mutateState(t, func(st *State) {
			for _, id := range revoked {
				if st.Agents[id] != nil {
					t.Fatal("expired revoked agent retained")
				}
			}
			if st.Agents[live[0]] == nil || !st.Agents[live[0]].Revoked {
				t.Fatal("revoked agent deleted before 24h")
			}
			// Saturate a live agent's key records: connect refuses with the replacement path.
			for i := 0; i < agentKeyRecords; i++ {
				st.Agents[live[1]].Keys[fmt.Sprint("k", i)] = &Key{Revoked: true, Changed: time.Now()}
			}
			st.Rates = map[string][]time.Time{}
		})
		expect(t, a.do("GET", "/home", nil), 200, "새 키를 더 연결할 수 없습니다")
		expect(t, a.do("POST", "/home/connect", url.Values{"agent": {live[1]}, "client": {supportedClient}, "mode": {"rotate"}}), 409, "기다려도 이 agent에 새 키 공간은 생기지 않습니다", `action="/home/agents"`, `href="/home"`)
		s.mutateState(t, func(st *State) {
			if len(st.Agents[live[1]].Keys) != agentKeyRecords || len(st.Connections) != 0 {
				t.Fatal("key-full refusal changed records", len(st.Connections))
			}
		})
	})

	t.Run("concurrent_agent_capacity_and_restart", func(t *testing.T) {
		s, mail := identityService(t, pool)
		a := newBrowser(s, "192.0.2.88")
		member := login(t, a, mail, "capacity@example.com")
		var wg sync.WaitGroup
		codes := make(chan int, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				b := &browser{t, s.Handler(), "192.0.2.88", map[string]string{sessionCookie: a.cookies[sessionCookie]}}
				codes <- b.do("POST", "/home/agents", nil).Code
			}()
		}
		wg.Wait()
		close(codes)
		wins := 0
		for c := range codes {
			if c == 303 {
				wins++
			} else if c != 409 {
				t.Fatal(c)
			}
		}
		if wins != 5 {
			t.Fatal("capacity race", wins)
		}
		a.h = (&Service{Pool: pool}).Handler()
		expect(t, a.do("POST", "/home/agents", nil), 409)
		var first string
		s.mutateState(t, func(st *State) {
			if len(st.Agents) != 5 {
				t.Fatal("agent overflow")
			}
			for id := range st.Agents {
				first = id
				break
			}
			st.Rates["http:new"] = make([]time.Time, 201)
			for i := range st.Rates["http:new"] {
				st.Rates["http:new"][i] = time.Now()
			}
		})
		expect(t, a.do("POST", "/home/agent-revoke", url.Values{"agent": {first}}), 303)
		s.mutateState(t, func(st *State) {
			if !st.agentCapacity(st.Members[member].Owner) {
				t.Fatal("revoke did not free slot")
			}
			st.Rates = map[string][]time.Time{}
		})
		expect(t, a.do("POST", "/home/agents", nil), 303)
	})

	t.Run("concurrent_completion_cancel_expire_and_reauth", func(t *testing.T) {
		s, mail := identityService(t, pool)
		a := newBrowser(s, "203.0.113.1")
		login(t, a, mail, "race@example.com")
		agent := browserAgent(t, a)
		token, id := browserGrant(t, a, agent, "register")
		_, proof := clientPrepare(t, s.Handler(), token, "key1")
		expect(t, a.do("POST", "/home/confirm", url.Values{"connection": {id}}), 303)
		var wg sync.WaitGroup
		wins := make(chan bool, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				raw, _ := json.Marshal(map[string]string{"token": token, "client": supportedClient, "proof": proof})
				r := httptest.NewRequest("POST", "/v1/connect/complete", strings.NewReader(string(raw)))
				w := httptest.NewRecorder()
				s.Handler().ServeHTTP(w, r)
				wins <- w.Code == 200
			}()
		}
		wg.Wait()
		close(wins)
		n := 0
		for win := range wins {
			if win {
				n++
			}
		}
		if n != 1 {
			t.Fatal("completion winners", n)
		}
		s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
		token, id = browserGrant(t, a, agent, "register")
		expect(t, a.do("POST", "/home/cancel", url.Values{"connection": {id}}), 303)
		publicCall(t, s.Handler(), "/v1/connect/info", "", map[string]string{"token": token, "client": supportedClient}, 401)
		token, id = browserGrant(t, a, agent, "register")
		s.mutateState(t, func(st *State) { st.Connections[id].Exp = time.Now().Add(-time.Second) })
		expect(t, a.do("GET", "/home/connections/"+id, nil), 200, "만료")
		publicCall(t, s.Handler(), "/v1/connect/info", "", map[string]string{"token": token, "client": supportedClient}, 401)
		token, id = browserGrant(t, a, agent, "register")
		clientPrepare(t, s.Handler(), token, "key3")
		s.mutateState(t, func(st *State) {
			st.Rates = map[string][]time.Time{}
			for _, ss := range st.Sessions {
				ss.Verified = ss.Verified.Add(-reauthWindow)
			}
		})
		expect(t, a.do("POST", "/home/confirm", url.Values{"connection": {id}}), 422, "5분")
	})
}
