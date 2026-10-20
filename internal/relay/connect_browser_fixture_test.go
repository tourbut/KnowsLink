//go:build integration && connectfixture

// Opt-in loopback U8/O-NM-03 fixture: synthetic accounts only, member agent/friend pages plus explicit connection and
// friend-request POSTs. No host install, mail, external account or real counterpart is involved.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestConnectBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("CONNECT_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	const life = 12 * time.Hour
	s, _ := identityService(t, messagePool(t))
	members, owners, agents, connections, codes := map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}
	email := func(name string) string { return "connect-" + name + "@example.test" }
	s.mutateState(t, func(st *State) {
		now := time.Now()
		for _, name := range []string{"a", "b", "c", "d", "e"} {
			v := st.signIn(emailIssuer, email(name), email(name), "", "", now)
			members[name] = st.Sessions[hashToken(v.Session)].Member
			m := st.Members[members[name]]
			m.Name, m.FriendCode, owners[name] = "합성 회원 "+strings.ToUpper(name), "fixture-code-"+name+"-"+randomToken()[:8], m.Owner
		}
		public, _, _ := ed25519.GenerateKey(rand.Reader)
		key := base64.RawURLEncoding.EncodeToString(public)
		agent := func(label, owner, name string) string {
			id := "agent_fx_" + label
			st.Agents[id] = &Agent{Name: name, Owner: owner, Keys: map[string]*Key{}}
			agents[label] = id
			return id
		}
		attempt := func(label, agentID, state string, device bool, exp time.Time) {
			c := &Connection{Owner: st.Agents[agentID].Owner, Agent: agentID, Client: supportedClient, Mode: "register", State: state, Exp: exp, Device: device}
			if state != "waiting" {
				c.Kid, c.Public = "fixture-key", key
			}
			id := hashToken(randomToken())
			st.Connections[id] = c
			connections[label] = id
		}
		// B's remote installs come first: approval respects the five-active-agent limit.
		for _, label := range []string{"remote_returning", "remote_lost", "remote_exchanged"} {
			code, e := st.approveRemote(&remoteAttempt{Client: "fixture-client", Redirect: "https://app.invalid/cb", Resource: "fixture", Owner: owners["b"], Exp: now.Add(time.Minute)}, now)
			if e != nil {
				t.Fatal(e)
			}
			c := st.Remote.Codes[hashToken(code)]
			agents[label] = st.Remote.Installs[c.Install].Agent
			st.Agents[agents[label]].Name = "합성 원격 " + label
			if label == "remote_returning" {
				c.Exp = now.Add(life)
			} else {
				delete(st.Remote.Codes, hashToken(code))
			}
			if label == "remote_exchanged" {
				st.Remote.Tokens[hashToken(randomToken())] = &remoteToken{c.Install, "fixture-client", "fixture", "refresh", now.Add(life)}
			}
		}
		// A: one open local stage per agent plus a registered key; a fifth agent slot stays free to try create twice.
		attempt("waiting", agent("waiting", owners["a"], "합성 대기 도우미"), "waiting", false, now.Add(life))
		attempt("prepared", agent("prepared", owners["a"], "합성 지문 확인 도우미"), "prepared", false, now.Add(life))
		attempt("approved", agent("approved", owners["a"], strings.Repeat("긴한글이름", 8)), "approved", false, now.Add(life))
		registered := agent("registered", owners["a"], "Synthetic registered helper with a long English name")
		st.Agents[registered].Keys["fixture-key"] = &Key{Public: public, Changed: now}
		attempt("consumed", registered, "consumed", false, now.Add(-time.Minute))
		// B: ended, uncertain and unused agents; C keeps no agent at all.
		attempt("expired", agent("expired", owners["b"], "합성 만료 도우미"), "expired", false, now.Add(-time.Minute))
		attempt("cancelled", agent("cancelled", owners["b"], "합성 취소 도우미"), "cancelled", false, now.Add(time.Minute))
		attempt("device_lost", agent("device_lost", owners["b"], "합성 저장 확인 도우미"), "approved", true, now.Add(-time.Minute))
		// O-NM-03 rows on A: existing, unknown, malformed, a cancelled one and a request from before tails were kept.
		st.requestFriend(members["a"], st.Members[members["b"]].FriendCode, now)
		for _, code := range []string{"fixture-unknown-zz99", "ab", "fixture-before-tail", "fixture-cancel-ca11"} {
			st.requestFriend(members["a"], code, now)
		}
		st.Friends[pairID(members["a"], "undisclosed:"+hashToken("fixture-before-tail"))].CodeTail = ""
		cancelled := st.Friends[pairID(members["a"], "undisclosed:"+hashToken("fixture-cancel-ca11"))]
		if e := st.decideFriend(members["a"], cancelled.ID, "cancel", cancelled.Generation, now); e != nil {
			t.Fatal(e)
		}
		st.requestFriend(members["c"], st.Members[members["a"]].FriendCode, now)
		// F1/V-CN-01 re-request paths on A, ended a day ago: D blocked A's pending request, and E now holds ten
		// synthetic incoming requests (receive limit). Re-sending either code must leave one A row, like ca11.
		before := now.Add(-25 * time.Hour)
		st.requestFriend(members["a"], st.Members[members["d"]].FriendCode, before)
		blocked := st.Friends[pairID(members["a"], members["d"])]
		if e := st.decideFriend(members["d"], blocked.ID, "block", blocked.Generation, before); e != nil {
			t.Fatal(e)
		}
		st.requestFriend(members["a"], st.Members[members["e"]].FriendCode, before)
		for i := 0; i < 10; i++ {
			sender := "fixture-sender-" + string(rune('0'+i))
			st.Friends[pairID(sender, members["e"])] = &Friend{ID: randomToken()[:22], A: sender, B: members["e"], Inviter: sender, NameA: "합성 발신자 " + string(rune('0'+i)), State: "pending", Generation: 1, Exp: now.Add(life)}
		}
		codes["blocked"], codes["receive_limit"], codes["unknown_cancelled"] = st.Members[members["d"]].FriendCode, st.Members[members["e"]].FriendCode, "fixture-cancel-ca11"
	})
	allowedGET := map[string]bool{"/": true, "/home": true, "/home/agents": true, "/home/friends": true, "/home/account": true, "/home/legacy": true, "/home/new": true, "/home/receipts": true}
	allowedPOST := map[string]bool{"/home/agents": true, "/home/connect": true, "/home/confirm": true, "/home/cancel": true, "/home/friends/request": true, "/home/friends/decision": true}
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
		page := strings.HasPrefix(r.URL.Path, "/home/connections/") || strings.HasPrefix(r.URL.Path, "/connect/")
		if !(r.Method == "GET" && (allowedGET[r.URL.Path] || page)) && !(r.Method == "POST" && allowedPOST[r.URL.Path]) {
			http.Error(w, "outside connect fixture scope", 403)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	data, _ := json.Marshal(map[string]any{"origin": server.URL, "agents": agents, "connections": connections, "rerequest_codes": codes, "accounts": []string{"a", "b", "c", "d", "e"}, "expires_at": time.Now().Add(life).UTC().Format(time.RFC3339)})
	if err := os.WriteFile(os.Getenv("CONNECT_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("connect fixture ready; synthetic accounts only")
	<-time.After(life)
}
