//go:build integration && langfixture

// Opt-in loopback U3 fixture: synthetic members with a saved English, a saved Korean and no saved language, plus the
// stored states that reach every translated screen. No mail, external account, host install or real counterpart.
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

func TestLanguageBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("LANG_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	const life = 12 * time.Hour
	f := setup(t, messagePool(t))
	members, owners, ids := map[string]string{}, map[string]string{}, map[string]string{}
	email := func(name string) string { return "lang-" + name + "@example.test" }
	now := time.Now()
	f.mutate(func(st *State) {
		for name, locale := range map[string]string{"en": "en", "ko": "ko", "none": ""} {
			v := st.signIn(emailIssuer, email(name), email(name), "", "", now)
			members[name] = st.Sessions[hashToken(v.Session)].Member
			m := st.Members[members[name]]
			m.Name, m.Locale, m.FriendCode, owners[name] = "Synthetic member "+strings.ToUpper(name), locale, "fixture-code-"+name+"-"+randomToken()[:8], m.Owner
		}
		// EN owns agent_b (registered key, legacy pair with agent_a), KO owns agent_a; the pair stays as legacy history.
		st.Agents["agent_b"].Owner, st.Agents["agent_b"].Name = owners["en"], "Synthetic registered helper with a long English name"
		st.Agents["agent_a"].Owner, st.Agents["agent_a"].Name = owners["ko"], strings.Repeat("긴한글이름", 8)
		st.Members[members["en"]].SelfConnect = true
	})
	// EN and KO are friends; one web message from EN to KO opens a conversation, a receipt and the friend screens.
	f.mutate(func(st *State) {
		st.requestFriend(members["en"], st.Members[members["ko"]].FriendCode, now)
		fr := st.Friends[pairID(members["en"], members["ko"])]
		if e := st.decideFriend(members["ko"], fr.ID, "accept", fr.Generation, now); e != nil {
			t.Fatal(e)
		}
		m, e := st.sendOwnerText(members["en"], "agent_b", "agent_a", "Synthetic fixture message", "fixture-attempt-0000000001", now)
		if e != nil {
			t.Fatal(e)
		}
		ids["message"] = m.Receipt.ID
		// Friend rows for the request states: a sent request still open, a denied one, a cancelled one, and an incoming one from NONE to EN.
		st.requestFriend(members["en"], "fixture-unknown-zz99", now)
		st.requestFriend(members["none"], st.Members[members["en"]].FriendCode, now)
		st.requestFriend(members["ko"], "fixture-cancel-ca11", now)
		cancelled := st.Friends[pairID(members["ko"], "undisclosed:"+hashToken("fixture-cancel-ca11"))]
		if e := st.decideFriend(members["ko"], cancelled.ID, "cancel", cancelled.Generation, now); e != nil {
			t.Fatal(e)
		}
	})
	// One pending approval for EN's agent_b, requested by KO's agent_a (the approval fixture's own path).
	parent, gate := "0199a3f2-4c10-7a11-8b22-00000000f001", "0199a3f2-4c10-7a11-8b22-00000000f002"
	receipt := f.send(wire(t, f.private["agent_a"], parent, "agent_a", "agent_b", "schedule.query", "lang-fixture-p01", "", queryBody(), now.Add(180*time.Second)), "agent_a", "", 200)
	claim := f.deliver(parent)
	f.send(wire(t, f.private["agent_b"], gate, "agent_b", "agent_b", "relay.approval.request", "lang-fixture-g01", parent, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, now.Add(time.Minute)), "agent_b", claim, 200)
	ids["gate"] = gate
	public, _, _ := ed25519.GenerateKey(rand.Reader)
	key := base64.RawURLEncoding.EncodeToString(public)
	f.mutate(func(st *State) {
		g := st.Gates[gate]
		g.Exp, st.Messages[parent].Receipt.Exp, st.Messages[gate].Receipt.Exp = now.Add(life), now.Add(life), now.Add(life)
		st.Messages[parent].ClaimToken = claim
		// EN: local connection stages on three more agents; NONE keeps no agent at all.
		for _, state := range []string{"waiting", "prepared", "approved"} {
			id := "agent_fx_" + state
			st.Agents[id] = &Agent{Name: "Synthetic " + state + " helper", Owner: owners["en"], Keys: map[string]*Key{}}
			c := &Connection{Owner: owners["en"], Agent: id, Client: supportedClient, Mode: "register", State: state, Exp: now.Add(life)}
			if state != "waiting" {
				c.Kid, c.Public = "fixture-key", key
			}
			cid := hashToken(randomToken())
			st.Connections[cid] = c
			ids["connection_"+state] = cid
		}
		// Device consent screen (anonymous login variant until EN signs in) and a remote consent attempt for EN.
		device := hashToken(randomToken())
		st.Connections[device] = &Connection{Device: true, Client: supportedClient, Mode: "register", Kid: "fixture-device", Public: key, State: "prepared", Exp: now.Add(life)}
		ids["device"] = device
		rs := st.remote()
		rs.Clients["fixture-client"] = &remoteClient{Name: "Synthetic remote app", Redirects: []string{"https://app.invalid/cb"}, Exp: now.Add(life)}
		rs.Attempts["fixture-attempt"] = &remoteAttempt{Client: "fixture-client", Redirect: "https://app.invalid/cb", Challenge: strings.Repeat("c", 43), State: "s", Resource: "fixture", Owner: owners["en"], Exp: now.Add(life)}
		ids["remote"] = "fixture-attempt"
	})
	allowedGET := map[string]bool{"/": true, "/auth/verify": true, "/home": true, "/home/name": true, "/home/agents": true, "/home/friends": true, "/home/account": true, "/home/legacy": true, "/home/new": true, "/home/receipts": true, "/home/conversation": true}
	allowedPOST := map[string]bool{"/home/language": true, "/home/profile": true, "/home/friends/request": true, "/home/friends/decision": true, "/home/agents": true, "/auth/logout": true}
	pagePrefixes := []string{"/home/connections/", "/connect/", "/home/gates/", "/home/friends/", "/remote/consent/", "/lang/"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if name, ok := strings.CutPrefix(r.URL.Path, "/fixture/as/"); ok {
			if r.Method != "POST" || r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+r.Host) || members[name] == "" {
				http.Error(w, "fixture origin, method or account refused", 403)
				return
			}
			var token string
			f.mutate(func(st *State) { token = st.signIn(emailIssuer, email(name), email(name), "", "", time.Now()).Session })
			setCookie(w, sessionCookie, token, int(sessionAbsolute/time.Second))
			w.WriteHeader(204)
			return
		}
		page := false
		for _, p := range pagePrefixes {
			page = page || strings.HasPrefix(r.URL.Path, p)
		}
		if !(r.Method == "GET" && (allowedGET[r.URL.Path] || page)) && !(r.Method == "POST" && allowedPOST[r.URL.Path]) {
			http.Error(w, "outside language fixture scope", 403)
			return
		}
		f.handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	data, _ := json.Marshal(map[string]any{"origin": server.URL, "accounts": map[string]string{"en": "saved English", "ko": "saved Korean", "none": "no saved language"}, "ids": ids,
		"conversation": "/home/conversation?from=agent_b&to=agent_a", "expires_at": time.Now().Add(life).UTC().Format(time.RFC3339)})
	if err := os.WriteFile(os.Getenv("LANG_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("language fixture ready; synthetic accounts only")
	<-time.After(life)
}
