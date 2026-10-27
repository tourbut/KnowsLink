//go:build integration && connectsimpleqa

// SAR-CONNECT-SIMPLE-001-QA loopback fixture: the real relay Handler on isolated Postgres, a local stand-in for Google
// (account chooser + signed ID tokens) and a few control endpoints for steps a person's wall clock or Google would do.
// Synthetic accounts only. No product code is changed; the file is copied into internal/relay only to run.
package relay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestConnectSimpleQAFixture(t *testing.T) {
	addr := os.Getenv("CSQA_ADDR")
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || addr == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	base := "http://" + addr
	s, _ := identityService(t, messagePool(t))
	s.Mail = nil
	g, sign := googleFixture(t)
	s.Google = g
	g.client = http.DefaultClient
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "local-fixture", RedirectURL: base + "/auth/google/callback", Scopes: []string{"openid", "email"},
		Endpoint: oauth2.Endpoint{AuthURL: base + "/fixture/google/auth", TokenURL: base + "/fixture/google/token", AuthStyle: oauth2.AuthStyleInParams}}
	accounts := []string{"a", "b", "c", "d", "e"}
	subject := func(x string) string { return "google-qa-" + x }
	mail := func(x string) string { return "csqa-" + x + "@example.test" }
	challenge := sync.Map{} // code -> PKCE challenge
	var dropComplete atomic.Bool
	var hits, starts atomic.Int64 // relay API requests / connection starts since the last reset
	handler := s.Handler()
	member := func(st *State, x string) *Member { return st.Members[st.Identities[googleIssuer+"|"+subject(x)]] }
	mux := http.NewServeMux()
	mux.HandleFunc("/fixture/google/auth", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var b strings.Builder
		b.WriteString("<!doctype html><meta charset=utf-8><title>Fixture Google 계정 선택</title><h1>Fixture Google 계정 선택 (합성)</h1><ul>")
		for _, x := range append(append([]string{}, accounts...), "newbie") {
			code := x + "." + q.Get("nonce")
			challenge.Store(code, q.Get("code_challenge"))
			fmt.Fprintf(&b, `<li><a id="pick-%s" href="/auth/google/callback?code=%s&state=%s">계정 %s (%s)</a></li>`, x, url.QueryEscape(code), url.QueryEscape(q.Get("state")), strings.ToUpper(x), mail(x))
		}
		fmt.Fprintf(&b, `<li><a id="pick-cancel" href="/auth/google/callback?error=access_denied&state=%s">취소</a></li></ul>`, url.QueryEscape(q.Get("state")))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(b.String()))
	})
	mux.HandleFunc("/fixture/google/token", func(w http.ResponseWriter, r *http.Request) {
		code := r.PostFormValue("code")
		want, _ := challenge.Load(code)
		x, nonce, _ := strings.Cut(code, ".")
		if want == nil || oauth2.S256ChallengeFromVerifier(r.PostFormValue("code_verifier")) != want {
			http.Error(w, "PKCE", 400)
			return
		}
		c := googleClaims(time.Now(), nonce)
		c["sub"], c["email"] = subject(x), mail(x)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(c)})
	})
	mux.HandleFunc("/fixture/as/", func(w http.ResponseWriter, r *http.Request) {
		x := strings.TrimPrefix(r.URL.Path, "/fixture/as/")
		if r.Method != "POST" || len(x) != 1 || !strings.Contains("abcde", x) {
			http.Error(w, "fixture method or account refused", 403)
			return
		}
		var token string
		s.mutateState(t, func(st *State) {
			token = st.signIn(googleIssuer, subject(x), mail(x), "", "", time.Now()).Session
			st.Members[st.Identities[googleIssuer+"|"+subject(x)]].Name = "합성 회원 " + strings.ToUpper(x)
		})
		setCookie(w, sessionCookie, token, int(sessionAbsolute/time.Second))
		_, _ = w.Write([]byte(token))
	})
	mux.HandleFunc("/fixture/state", func(w http.ResponseWriter, r *http.Request) {
		x := r.URL.Query().Get("m")
		out := map[string]any{}
		s.mutateState(t, func(st *State) {
			m := member(st, x)
			if m == nil {
				return
			}
			var agents, attempts []map[string]any
			for id, a := range st.Agents {
				if a.Owner == m.Owner {
					agents = append(agents, map[string]any{"id": id, "name": a.Name, "revoked": a.Revoked, "keys": len(a.Keys)})
				}
			}
			for id, c := range st.Connections {
				if c.Owner == m.Owner && c.Device {
					attempts = append(attempts, map[string]any{"id": id, "state": c.State, "agent": c.Agent, "name": c.Name, "host": c.Host, "exp": c.Exp.Format(time.RFC3339)})
				}
			}
			friends := 0
			for _, f := range st.Friends {
				if (f.A == st.Identities[googleIssuer+"|"+subject(x)] || f.B == st.Identities[googleIssuer+"|"+subject(x)]) && f.State == "active" {
					friends++
				}
			}
			out = map[string]any{"member": st.Identities[googleIssuer+"|"+subject(x)], "friendCode": m.FriendCode, "agents": agents, "attempts": attempts, "friends": friends}
		})
		_ = json.NewEncoder(w).Encode(out)
	})
	// Moves one device attempt's expiry (the 10-minute and +24h clocks cannot be waited for); nothing else changes.
	mux.HandleFunc("/fixture/exp", func(w http.ResponseWriter, r *http.Request) {
		d, err := time.ParseDuration(r.URL.Query().Get("delta"))
		if err != nil || r.Method != "POST" {
			http.Error(w, "delta or method refused", 400)
			return
		}
		ok := false
		s.mutateState(t, func(st *State) {
			if c := st.Connections[r.URL.Query().Get("connection")]; c != nil && c.Device {
				c.Exp, ok = time.Now().Add(d), true
			}
		})
		if !ok {
			http.Error(w, "unknown device attempt", 404)
		}
	})
	// Fills the active-member cap with synthetic accounts so a new Google account hits the signup-closed path (W-02).
	mux.HandleFunc("/fixture/fill-members", func(w http.ResponseWriter, r *http.Request) {
		n := 0
		s.mutateState(t, func(st *State) {
			now := time.Now()
			for i := 0; len(st.Members) < activeMembers; i++ {
				st.signIn(emailIssuer, fmt.Sprintf("csqa-filler-%d@example.test", i), fmt.Sprintf("csqa-filler-%d@example.test", i), "", "", now)
				n++
			}
		})
		fmt.Fprintf(w, "%d", n)
	})
	// The next /v1/connect/complete never reaches the relay and the client sees a dropped connection (key not stored).
	mux.HandleFunc("/fixture/drop-complete", func(w http.ResponseWriter, r *http.Request) {
		dropComplete.Store(r.URL.Query().Get("on") == "1")
	})
	// Seeds an agent without a name, the shape an older relay left behind (A-05); no key, nothing else changes.
	mux.HandleFunc("/fixture/nameless", func(w http.ResponseWriter, r *http.Request) {
		id := ""
		s.mutateState(t, func(st *State) {
			if m := member(st, r.URL.Query().Get("m")); m != nil && r.Method == "POST" {
				id = newMemberAgentID()
				st.addAgent(id, &Agent{Owner: m.Owner, Keys: map[string]*Key{}}, time.Now())
			}
		})
		fmt.Fprint(w, id)
	})
	mux.HandleFunc("/fixture/hits", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("reset") == "1" {
			hits.Store(0)
			starts.Store(0)
		}
		fmt.Fprintf(w, `{"api":%d,"starts":%d}`, hits.Load(), starts.Load())
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/") {
			hits.Add(1)
		}
		if r.URL.Path == "/v1/connect/start" {
			starts.Add(1)
		}
		if r.URL.Path == "/v1/connect/complete" && dropComplete.Load() {
			if conn, _, err := http.NewResponseController(w).Hijack(); err == nil {
				_ = conn.Close()
			}
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		handler.ServeHTTP(w, r)
	})
	server := &http.Server{Addr: addr, Handler: mux}
	go func() { _ = server.ListenAndServe() }()
	t.Log("connect-simple QA fixture ready on " + addr + "; synthetic accounts only")
	<-time.After(5 * time.Hour)
}
