//go:build integration

// Google client HTTP flow uses a signed local provider and isolated Postgres; no real accounts or secrets are needed.
package relay

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

func TestGoogleDeviceHTTP(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" {
		t.Fatal("isolated synthetic database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, _ := identityService(t, pool)
	s.Mail = nil
	g, sign := googleFixture(t)
	s.Google = g
	var auth url.Values
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if oauth2.S256ChallengeFromVerifier(r.PostFormValue("code_verifier")) != auth.Get("code_challenge") {
			http.Error(w, "PKCE", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(googleClaims(time.Now(), auth.Get("nonce")))})
	}))
	defer provider.Close()
	g.client = provider.Client()
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "local-fixture", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: provider.URL, AuthStyle: oauth2.AuthStyleInParams}}
	cookies := map[string]*http.Cookie{}
	do := func(method, path, body string, browser, cross bool) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		if browser {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Sec-Fetch-Site", "same-origin")
			if cross {
				r.Header.Set("Sec-Fetch-Site", "cross-site")
			}
			for _, c := range cookies {
				if !cross || c.SameSite != http.SameSiteStrictMode {
					r.AddCookie(c)
				}
			}
		} else {
			r.Header.Set("Content-Type", "application/json")
		}
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		for _, c := range w.Result().Cookies() {
			if c.MaxAge < 0 {
				delete(cookies, c.Name)
			} else {
				cookies[c.Name] = c
			}
		}
		return w
	}
	var firstAgent, firstCredential, firstID string
	for i := range 2 {
		fixture, now := connectionFixture()
		token, c, private := deviceFixture(t, fixture, now)
		proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
		input, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "kid": c.Kid, "public": c.Public, "proof": proof})
		expect(t, do("POST", "/v1/connect/start", string(input), false, false), 200, "requested")
		id := hashToken(token)
		expect(t, do("GET", "/connect/"+id, "", true, false), 200, "Google로 로그인")
		w := do("POST", "/auth/google", "connection="+id, true, false)
		expect(t, w, 303, googleIssuer)
		u, _ := url.Parse(w.Header().Get("Location"))
		auth = u.Query()
		callback := "/auth/google/callback?code=fixture&state=" + auth.Get("state")
		expect(t, do("GET", callback, "", true, true), 200, "/connect/"+id)
		expect(t, do("GET", "/connect/"+id, "", true, false), 200, "/home/device-confirm", fingerprint(decodePublic(c.Public)))
		expect(t, do("POST", "/home/device-confirm", "connection="+id+"&name=Synthetic+device", true, true), 403)
		expect(t, do("POST", "/home/confirm", "connection="+id, true, false), 403)
		if i == 0 {
			s.mutateState(t, func(st *State) {
				st.Sessions[hashToken(cookies[sessionCookie].Value)].Verified = time.Now().Add(-6 * time.Minute)
			})
			expect(t, do("POST", "/home/device-confirm", "connection="+id+"&name=Synthetic+device", true, false), 422)
			s.mutateState(t, func(st *State) { st.Sessions[hashToken(cookies[sessionCookie].Value)].Verified = time.Now() })
		}
		expect(t, do("POST", "/home/device-confirm", "connection="+id+"&name=Synthetic+device", true, false), 303)
		expect(t, do("POST", "/home/device-confirm", "connection="+id+"&name=Synthetic+device", true, false), 401)
		poll, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "proof": proof})
		w = do("POST", "/v1/connect/poll", string(poll), false, false)
		expect(t, w, 200, "approved")
		var approved map[string]string
		if json.Unmarshal(w.Body.Bytes(), &approved) != nil || approved["name"] != "Synthetic device" {
			t.Fatal("poll response", approved)
		}
		c.Owner, c.Agent = approved["owner"], approved["agent"]
		completion := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c)))
		payload, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "proof": completion})
		w = do("POST", "/v1/connect/complete", string(payload), false, false)
		expect(t, w, 200, "connected")
		// W-01: the member who started the link sees the finished connection and its agent name.
		expect(t, do("GET", "/connect/"+id, "", true, false), 200, "이 연결은 완료됐습니다", "Synthetic device")
		firstID = id
		var connected map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &connected)
		// U2-A: the browser name labels the client's agent, and the refused retry made no extra agent.
		s.mutateState(t, func(st *State) {
			named := 0
			for _, a := range st.Agents {
				if a.Name == "Synthetic device" {
					named++
				}
			}
			if st.Agents[connected["agent"]] == nil || st.Agents[connected["agent"]].Name != "Synthetic device" || named != i+1 {
				t.Fatal("device name not bound to exactly one client agent")
			}
		})
		if i == 0 {
			firstAgent, firstCredential = connected["agent"], connected["credential"]
		} else if firstAgent == connected["agent"] || firstCredential == connected["credential"] {
			t.Fatal("clients share authority")
		}
		expect(t, do("POST", "/v1/connect/complete", string(payload), false, false), 401)
		expect(t, do("GET", callback, "", true, true), 401)
	}
	s.mutateState(t, func(st *State) {
		if len(st.Members) != 1 || len(st.Agents) != 2 || len(st.Pairs) != 0 {
			t.Fatal("identity/agent/pair isolation")
		}
	})

	begin := func(host string) (string, string, func() *httptest.ResponseRecorder) {
		t.Helper()
		fixture, now := connectionFixture()
		token, c, private := deviceFixture(t, fixture, now)
		proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
		input, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "kid": c.Kid, "public": c.Public, "proof": proof, "host": host})
		expect(t, do("POST", "/v1/connect/start", string(input), false, false), 200, "requested")
		id := hashToken(token)
		w := do("POST", "/auth/google", "connection="+id, true, false)
		u, _ := url.Parse(w.Header().Get("Location"))
		auth = u.Query()
		expect(t, do("GET", "/auth/google/callback?code=fixture&state="+auth.Get("state"), "", true, true), 200, "/connect/"+id)
		poll, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "proof": proof})
		finish := func() *httptest.ResponseRecorder { // saves the key of an approved attempt
			var approved map[string]string
			_ = json.Unmarshal(do("POST", "/v1/connect/poll", string(poll), false, false).Body.Bytes(), &approved)
			c.Owner, c.Agent = approved["owner"], approved["agent"]
			completion := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c)))
			payload, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "proof": completion})
			return do("POST", "/v1/connect/complete", string(payload), false, false)
		}
		return id, string(poll), finish
	}
	// A-04 default name, A-01 no agent before the key is saved, M-02 the attempt and its removal, W-01 the revisit.
	s.mutateState(t, func(st *State) {
		for _, m := range st.Members {
			m.Name = "합성 회원" // the agents page asks for a member name first
		}
	})
	id, poll, _ := begin("claude")
	expect(t, do("GET", "/connect/"+id, "", true, false), 200, `value="Claude Code"`, "1. 지문 비교", "지문이 같습니다. 연결합니다")
	expect(t, do("POST", "/home/device-confirm", "connection="+id+"&name=Claude+Code", true, false), 303)
	expect(t, do("GET", "/connect/"+id, "", true, false), 200, "연결을 승인했습니다", "Claude Code")
	expect(t, do("GET", "/home/agents", "", true, false), 200, "진행 중인 연결", "동의 완료 · 키 저장 확인 필요", "이 시도 정리")
	s.mutateState(t, func(st *State) {
		if len(st.Agents) != 2 {
			t.Fatal("approval made an agent before its key was saved")
		}
	})
	expect(t, do("POST", "/home/attempt-remove", "connection="+id, true, false), 303)
	expect(t, do("POST", "/v1/connect/poll", poll, false, false), 401)
	if body := do("GET", "/home/agents", "", true, false).Body.String(); strings.Contains(body, "진행 중인 연결") {
		t.Fatal("removed attempt still listed")
	}
	// Signed out, a finished link names nothing (W-01).
	session := cookies[sessionCookie]
	delete(cookies, sessionCookie)
	out := do("GET", "/connect/"+firstID, "", true, false)
	expect(t, out, 401, "연결 수단이 만료")
	if body := out.Body.String(); strings.Contains(body, "Synthetic device") || strings.Contains(body, "완료") {
		t.Fatal("signed-out revisit revealed the connection")
	}
	cookies[sessionCookie] = session
	// M-02: another member's attempt and a consumed one are never removed.
	s.mutateState(t, func(st *State) {
		st.Connections["foreign-attempt"] = &Connection{Owner: "other-owner", Client: supportedClient, Device: true, State: "prepared", Exp: time.Now().Add(connectionTTL)}
	})
	expect(t, do("POST", "/home/attempt-remove", "connection=foreign-attempt", true, false), 403)
	expect(t, do("POST", "/home/attempt-remove", "connection="+firstID, true, false), 403)
	s.mutateState(t, func(st *State) {
		if st.Connections["foreign-attempt"] == nil || st.Connections[firstID] == nil || st.Connections[firstID].State != "consumed" {
			t.Fatal("attempt-remove removed an attempt it does not own")
		}
		delete(st.Connections, "foreign-attempt")
	})
	// L-04 at key save: with 19 agents today, two consented attempts both pass consent; the second saved key is
	// refused with when to retry (E agent_rate, CS-AC-04).
	first, _, finishFirst := begin("")
	second, _, finishSecond := begin("")
	s.mutateState(t, func(st *State) {
		for _, m := range st.Members {
			made := 0
			for _, a := range st.Agents {
				if a.Owner == m.Owner && time.Since(a.Created) < 24*time.Hour {
					made++
				}
			}
			for i := made; i < agentRate-1; i++ {
				st.Agents[fmt.Sprint("agent_today", i)] = &Agent{Owner: m.Owner, Keys: map[string]*Key{}, Created: time.Now()}
			}
		}
		st.Rates = map[string][]time.Time{}
	})
	expect(t, do("POST", "/home/device-confirm", "connection="+first+"&name=First", true, false), 303)
	expect(t, do("POST", "/home/device-confirm", "connection="+second+"&name=Second", true, false), 303)
	expect(t, finishFirst(), 200, "connected")
	expect(t, finishSecond(), 429, `"error":"agent_rate"`, `"retry_at"`)
	// L-04 at consent: past the 24h rate the browser refuses and the polling client learns when to retry (E agent_rate).
	id, poll, _ = begin("")
	s.mutateState(t, func(st *State) {
		for _, m := range st.Members {
			for i := range agentRate {
				st.Agents[fmt.Sprint("agent_seed", i)] = &Agent{Owner: m.Owner, Keys: map[string]*Key{}, Created: time.Now()}
			}
		}
		st.Rates = map[string][]time.Time{}
	})
	expect(t, do("POST", "/home/device-confirm", "connection="+id+"&name=Seed", true, false), 429, "오늘 새 에이전트를 많이 만들었습니다")
	expect(t, do("POST", "/v1/connect/poll", poll, false, false), 429, `"error":"agent_rate"`, `"retry_at"`)
	expect(t, do("GET", "/connect/"+id, "", true, false), 429, "오늘 새 에이전트를 많이 만들었습니다", `href="/home/agents"`, "내 에이전트로 돌아가기")
}

func decodePublic(value string) []byte {
	public, _ := base64.RawURLEncoding.DecodeString(value)
	return public
}
