//go:build integration && receiptfixture

// Opt-in loopback receipt fixture permits synthetic account switches and failures, never message sending or external auth.
package relay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReceiptBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("RECEIPT_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, _ := identityService(t, pool)
	tokens, members := map[string]string{}, map[string]string{}
	s.mutateState(t, func(st *State) {
		now := time.Now()
		for _, name := range []string{"a", "b", "c"} {
			v := st.signIn(emailIssuer, "receipt-"+name+"@example.test", "receipt-"+name+"@example.test", "", "", now)
			tokens[name] = v.Session
			id := st.Sessions[hashToken(v.Session)].Member
			members[name], st.Members[id].Name = id, "합성 회원 "+name
			st.Agents["agent_"+name] = &Agent{Owner: st.Members[id].Owner, Name: strings.Repeat("긴한글", 13), Client: "Node 22 로컬 클라이언트"}
		}
		st.Agents["agent_a2"] = &Agent{Owner: st.Members[members["a"]].Owner, Name: st.Agents["agent_a"].Name, Client: "Node 22 로컬 클라이언트"}
		st.Messages[firstID] = &Message{Receipt: Receipt{ID: firstID, From: "agent_a", To: "agent_b", State: "delivered", Accepted: now, Exp: now.Add(180 * time.Second)}, Completion: "received"}
		st.Messages[gateID] = &Message{Receipt: Receipt{ID: gateID, From: "agent_b", To: "agent_c", State: "delivered", Accepted: now, Exp: now.Add(180 * time.Second)}}
	})
	var outage atomic.Bool
	handler := s.Handler()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if strings.HasPrefix(r.URL.Path, "/fixture/") {
			if r.Method != "POST" || r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+r.Host) {
				http.Error(w, "fixture origin or method refused", 403)
				return
			}
			name := strings.TrimPrefix(r.URL.Path, "/fixture/as/")
			if tokens[name] != "" {
				var token string
				s.mutateState(t, func(st *State) {
					v := st.signIn(emailIssuer, "receipt-"+name+"@example.test", "receipt-"+name+"@example.test", "", "", time.Now())
					token = v.Session
				})
				setCookie(w, sessionCookie, token, int(sessionAbsolute/time.Second))
			} else {
				switch r.URL.Path {
				case "/fixture/outage":
					outage.Store(true)
				case "/fixture/recover":
					outage.Store(false)
					s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
				case "/fixture/expire", "/fixture/limit", "/fixture/revoke":
					s.mutateState(t, func(st *State) {
						if session := st.Sessions[hashToken(readCookie(r, sessionCookie))]; session != nil {
							switch r.URL.Path {
							case "/fixture/expire":
								session.Seen = time.Now().Add(-sessionIdle)
							case "/fixture/limit":
								for range 40 {
									st.Rates["http:member:"+session.Member] = append(st.Rates["http:member:"+session.Member], time.Now())
								}
							case "/fixture/revoke":
								st.Agents["agent_a"].Revoked = true
							}
						}
					})
				default:
					http.NotFound(w, r)
					return
				}
			}
			w.WriteHeader(204)
			return
		}
		if r.Method != "GET" || (r.URL.Path != "/" && r.URL.Path != "/home/receipts" && r.URL.Path != "/home/agents" && r.URL.Path != "/home/legacy") {
			http.Error(w, "outside receipt fixture scope", 403)
			return
		}
		if outage.Load() {
			ctx, cancel := context.WithCancel(r.Context())
			cancel()
			r = r.WithContext(ctx)
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	data, _ := json.Marshal(map[string]string{"origin": server.URL, "owned": firstID, "foreign": gateID, "missing": secondID})
	if err := os.WriteFile(os.Getenv("RECEIPT_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("receipt fixture ready; synthetic accounts and local failures only")
	<-time.After(24 * time.Hour)
}
