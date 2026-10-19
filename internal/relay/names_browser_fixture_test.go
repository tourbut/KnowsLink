//go:build integration && namesfixture

// Opt-in loopback U2/U4 fixture: synthetic accounts only, member pages plus explicit friend/name POSTs, no message sending.
package relay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestNamesBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("NAMES_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, _ := identityService(t, pool)
	members := map[string]string{}
	var friendAB, requestCA string
	email := func(name string) string { return "names-" + name + "@example.test" }
	s.mutateState(t, func(st *State) {
		now := time.Now()
		for _, name := range []string{"a", "b", "c"} {
			v := st.signIn(emailIssuer, email(name), email(name), "", "", now)
			id := st.Sessions[hashToken(v.Session)].Member
			members[name], st.Members[id].Name = id, "합성 회원 "+name
		}
		st.Members[members["a"]].Name = strings.Repeat("가", 40)
		owner := func(name string) string { return st.Members[members[name]].Owner }
		// A: two same-name 39-character agents; B: one emoji name; C: no agent at all.
		long := strings.Repeat("긴한글", 13)
		st.Agents["agent_a1"] = &Agent{Owner: owner("a"), Name: long, Client: "Node 22 로컬 클라이언트"}
		st.Agents["agent_a2"] = &Agent{Owner: owner("a"), Name: long, Client: "Node 22 로컬 클라이언트"}
		st.Agents["agent_b1"] = &Agent{Owner: owner("b"), Name: strings.Repeat("😀", 40), Client: "Node 22 로컬 클라이언트"}
		// Legacy: an exact pair with an agent whose owner is no human member, so its name must stay hidden.
		st.Owners["owner-legacy"] = &Owner{Active: true}
		st.Agents["agent_legacy_x"] = &Agent{Owner: "owner-legacy", Name: "보이면 안 되는 이름"}
		st.Pairs[pairID("agent_a1", "agent_legacy_x")] = &Pair{A: "agent_a1", B: "agent_legacy_x", Inviter: "agent_a1", Recipient: "agent_legacy_x", State: "active", Generation: 1}
		st.requestFriend(members["a"], st.Members[members["b"]].FriendCode, now)
		f := st.Friends[pairID(members["a"], members["b"])]
		if e := st.decideFriend(members["b"], f.ID, "accept", f.Generation, now); e != nil {
			t.Fatal(e)
		}
		st.requestFriend(members["c"], st.Members[members["a"]].FriendCode, now)
		friendAB, requestCA = f.ID, st.Friends[pairID(members["c"], members["a"])].ID
	})
	allowedGET := map[string]bool{"/": true, "/home": true, "/home/friends": true, "/home/account": true, "/home/agents": true, "/home/legacy": true, "/home/receipts": true, "/home/new": true}
	allowedPOST := map[string]bool{"/home/friends/request": true, "/home/friends/decision": true, "/home/friend-code": true, "/home/agent-name": true}
	handler := s.Handler()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if name, ok := strings.CutPrefix(r.URL.Path, "/fixture/as/"); ok {
			if r.Method != "POST" || r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+r.Host) || members[name] == "" {
				http.Error(w, "fixture origin, method or account refused", 403)
				return
			}
			var token string
			s.mutateState(t, func(st *State) { token = st.signIn(emailIssuer, email(name), email(name), "", "", time.Now()).Session })
			setCookie(w, sessionCookie, token, int(sessionAbsolute/time.Second))
			w.WriteHeader(204)
			return
		}
		if !(r.Method == "GET" && (allowedGET[r.URL.Path] || strings.HasPrefix(r.URL.Path, "/home/friends/"))) && !(r.Method == "POST" && allowedPOST[r.URL.Path]) {
			http.Error(w, "outside names fixture scope", 403)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	data, _ := json.Marshal(map[string]string{"origin": server.URL, "friend_ab": friendAB, "request_ca": requestCA})
	if err := os.WriteFile(os.Getenv("NAMES_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("names fixture ready; synthetic accounts only")
	<-time.After(24 * time.Hour)
}
