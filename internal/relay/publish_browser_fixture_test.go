//go:build integration && publishfixture

// Opt-in loopback SAR-PUBLISH-001 fixture: synthetic legacy calendar state, the D018 startup cleanup against a local
// revoke endpoint, policy pages, settings and schedule.query approval in both languages. No Google or mail traffic.
package relay

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPublishBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("PUBLISH_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	pool := messagePool(t)
	f := setup(t, pool)
	const pendingGate = "0199a3f2-4c10-7a11-8b22-3344556677ab"
	const queryID = "0199a3f2-4c10-7a11-8b22-334455667702"
	// Two schedule.query gates for Bob: gateID gets a legacy lookup result (cleanup denies it), pendingGate stays
	// pending so it can be approved after cleanup (U5: gate only, no disclosure).
	var member string
	f.mutate(func(st *State) {
		v := st.signIn(googleIssuer, "sub-fixture-bob", "bob@example.test", "", "", time.Now())
		member = st.Sessions[hashToken(v.Session)].Member
		st.Members[member].Name = "밥"
		st.Agents["agent_b"].Owner = st.Members[member].Owner
		st.Agents["agent_b"].Name = "밥의 도우미"
	})
	var claims sync.Map
	for _, g := range []struct{ parent, gate, key string }{{firstID, gateID, "publish-fixture-a"}, {queryID, pendingGate, "publish-fixture-b"}} {
		receipt := f.send(f.message(g.parent, g.key+"-0001"), "agent_a", "", 200)
		claim := f.deliver(g.parent)
		claims.Store(g.gate, claim)
		f.send(wire(t, f.private["agent_b"], g.gate, "agent_b", "agent_b", "relay.approval.request", g.key+"-gate", g.parent, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
	}
	f.mutate(func(st *State) {
		// Stored deadlines stretch for the fixture lifetime only; signatures, digests and generations stay as sent.
		for _, id := range []string{gateID, pendingGate} {
			g, life := st.Gates[id], time.Now().Add(12*time.Hour)
			g.Exp, st.Messages[g.Parent].Receipt.Exp, st.Messages[id].Receipt.Exp = life, life, life
		}
	})
	key := bytes.Repeat([]byte{7}, 32)
	seal := func(plain, owner string) string {
		block, _ := aes.NewCipher(key)
		gcm, _ := cipher.NewGCM(block)
		nonce := make([]byte, gcm.NonceSize())
		return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(plain), []byte(owner)))
	}
	// seed writes the JSON the SAR-SCHEDULE-001 build stored: two connections and an undelivered result on gateID.
	seed := func(ctx context.Context) error {
		calendars := `{"` + member + `":{"Sealed":"` + seal("synthetic-refresh-1", member) + `","State":"connected","Changed":"2026-10-21T00:00:00Z"},"mem_other":{"Sealed":"` + seal("synthetic-refresh-2", "mem_other") + `","State":"reconnect","Changed":"2026-10-21T00:00:00Z"}}`
		lookup := `{"Started":"2026-10-21T00:00:00Z","Finished":"2026-10-21T00:00:01Z","State":"done","Disclosure":{"busy":[],"granularity_min":30,"window":{"end":"2026-10-03T11:00:00Z","start":"2026-10-03T10:00:00Z"}}}`
		_, err := pool.Exec(ctx, `UPDATE relay_state SET data=jsonb_set(jsonb_set(data,'{Calendars}',$1::jsonb),ARRAY['Gates',$2::text,'Lookup'],$3::jsonb),epoch=epoch+1`, calendars, gateID, lookup)
		if err == nil {
			_, err = pool.Exec(ctx, `UPDATE relay_state SET data=jsonb_set(data,ARRAY['Gates',$1::text,'State'],'"approved"'),epoch=epoch+1`, gateID)
		}
		return err
	}
	if err := seed(context.Background()); err != nil {
		t.Fatal(err)
	}
	var refuse atomic.Bool
	var revokes atomic.Int64
	revoke := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		revokes.Add(1)
		if refuse.Load() {
			w.WriteHeader(500)
		}
	}))
	defer revoke.Close()
	// summary reports what is stored without any token, email or busy value.
	summary := func(ctx context.Context) map[string]any {
		var raw []byte
		_ = pool.QueryRow(ctx, `SELECT data::text FROM relay_state`).Scan(&raw)
		out := map[string]any{"calendar_key_in_state": bytes.Contains(raw, []byte(`"Calendars"`)), "lookup_in_state": bytes.Contains(raw, []byte(`"Lookup"`)), "disclosure_in_state": bytes.Contains(raw, []byte(`"Disclosure"`)), "revoke_requests": revokes.Load()}
		f.mutate(func(st *State) {
			out["connections"] = len(st.Calendars)
			out["gate_with_result"] = st.Gates[gateID].State
			out["gate_pending"] = st.Gates[pendingGate].State
		})
		return out
	}
	handler := f.s.Handler()
	stop := make(chan struct{})
	var once sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/fixture/") {
			if r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+r.Host) {
				http.Error(w, "fixture origin refused", 403)
				return
			}
			reply := func(v any, err error) {
				if err != nil {
					http.Error(w, "fixture step failed", 503)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(v)
			}
			switch r.Method + " " + r.URL.Path {
			case "GET /fixture/state":
				reply(summary(r.Context()), nil)
			case "POST /fixture/seed":
				reply(nil, seed(r.Context()))
			case "POST /fixture/retire":
				// key=absent runs the no-key path; revoke=fail makes the local revoke endpoint refuse every token.
				refuse.Store(r.URL.Query().Get("revoke") == "fail")
				secret := base64.StdEncoding.EncodeToString(key)
				if r.URL.Query().Get("key") == "absent" {
					secret = ""
				}
				var logs bytes.Buffer
				prev := log.Writer()
				log.SetOutput(&logs)
				err := f.s.retireCalendars(r.Context(), secret, revoke.URL, revoke.Client())
				log.SetOutput(prev)
				out := summary(r.Context())
				out["cleanup_record"] = strings.TrimSpace(logs.String())
				reply(out, err)
			case "POST /fixture/as/ko", "POST /fixture/as/en":
				var session string
				_, err := f.s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
					v := st.signIn(googleIssuer, "sub-fixture-bob", "bob@example.test", "", "", now)
					session = v.Session
					st.Members[member].Locale = strings.TrimPrefix(r.URL.Path, "/fixture/as/")
					return nil, nil
				})
				if err != nil {
					http.Error(w, "fixture transaction failed", 503)
					return
				}
				setCookie(w, sessionCookie, session, int(sessionAbsolute/time.Second))
				http.Redirect(w, r, "/home/account", 303)
			case "POST /fixture/consume":
				// The agent side of an approved gate: the response shows there is no disclosure (RC-04).
				id := r.URL.Query().Get("gate")
				claim, _ := claims.Load(id)
				c, _ := claim.(string)
				body, _ := json.Marshal(map[string]any{"id": id, "claim": c})
				rec := f.request("POST", "/v1/gate-consume", f.tokens["agent_b"], body, "")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(rec.Code)
				_, _ = w.Write(rec.Body.Bytes())
			case "POST /fixture/stop":
				w.WriteHeader(204)
				once.Do(func() { close(stop) })
			default:
				http.NotFound(w, r)
			}
			return
		}
		path := r.URL.Path
		if strings.HasPrefix(path, "/home/gates/") || strings.HasPrefix(path, "/lang/") || strings.HasPrefix(path, "/home/calendar/") {
			path = path[:strings.LastIndex(path, "/")] + "/*"
		}
		allowed := map[string]bool{
			"GET /": true, "GET /privacy": true, "GET /terms": true, "GET /lang/*": true, "GET /home": true, "GET /home/account": true,
			"GET /home/legacy": true, "GET /home/gates/*": true, "POST /home/gates/*": true, "POST /home/calendar/*": true, "POST /auth/logout": true,
		}
		if !allowed[r.Method+" "+path] && !(r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/deny")) {
			http.Error(w, "outside publish fixture scope", 403)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	data, _ := json.Marshal(map[string]string{"url": server.URL, "gate_with_result": gateID, "gate_pending": pendingGate})
	if err := os.WriteFile(os.Getenv("PUBLISH_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("publish fixture ready; synthetic legacy calendar state, local revoke endpoint only")
	select {
	case <-stop:
	case <-time.After(24 * time.Hour):
	}
}
