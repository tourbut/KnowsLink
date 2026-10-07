// Cleanup admission classification: honored credential, same-origin form, gate CSRF, valid body and own record.
package relay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Review M1: cleanup admission needs an honored credential, same-origin form, gate CSRF, valid body and the caller's own record.
func TestCleanupAdmissionNeedsVerifiedOwnRecord(t *testing.T) {
	st, _, now := textFixture()
	st.Sessions[hashToken("session-a")] = &Session{Member: "agent_a", Created: now, Seen: now}
	st.Sessions[hashToken("session-b")] = &Session{Member: "agent_b", Created: now, Seen: now}
	st.Gates["gate"] = &Gate{Owner: "agent_a", State: "pending"}
	st.Agents["agent_c"].Revoked = true
	for id, persisted := range map[string]bool{"persisted": true, "unpersisted": false} {
		st.Messages[id] = &Message{Receipt: Receipt{ID: id, From: "agent_a", To: "agent_b", Intent: "relay.request", State: "leased", Exp: now.Add(time.Minute)},
			Kid: "key1", Generation: 1, Deliver: "agent", LeaseToken: "lease", LeaseUntil: now.Add(30 * time.Second), Persisted: persisted}
	}
	s := &Service{}
	s.remember(st, now, 1)
	live := st.liveCredentials()
	for token, want := range map[string]bool{"agent_a": true, "session-a": true, "agent_c": false, "invented": false, "": false} {
		if _, ok := live[hashToken(token)]; ok != want {
			t.Fatal("live credential", token)
		}
	}
	request := func(path, body, cookie, auth string, header ...string) *http.Request {
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		if strings.Contains(path, "/home") {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			_ = r.ParseForm()
		}
		r.AddCookie(&http.Cookie{Name: sessionCookie, Value: cookie})
		r.Header.Set("Authorization", "Bearer "+auth)
		for i := 0; i+1 < len(header); i += 2 {
			r.Header.Set(header[i], header[i+1])
		}
		return r.WithContext(context.WithValue(r.Context(), requestBodyKey{}, []byte(body)))
	}
	deny := "decision=deny&csrf=" + csrf("session-a", "gate")
	revoke := `{"agent":"agent_a","kid":"key1"}`
	for _, c := range []struct {
		r          *http.Request
		credential string
		principal  string
		own        bool
	}{
		{request("/v1/key-revoke", revoke, "", "agent_a"), "agent_a", "agent_a", true},
		{request("/v1/key-revoke", revoke, "", "agent_b"), "agent_b", "agent_b", false},
		{request("/v1/key-revoke", `{"agent":`, "", "agent_a"), "", "agent_a", false},
		{request("/v1/key-revoke", revoke, "", "agent_a", "Sec-Fetch-Site", "cross-site"), "", "agent_a", true},
		{request("/v1/ack", `{"id":"none","token":"none"}`, "", "agent_b"), "agent_b", "agent_b", false},
		{request("/v1/ack", `{"id":"persisted","token":"lease"}`, "", "agent_b"), "agent_b", "agent_b", true},
		// FIX-2 (Sol): the ACK handler refuses an unpersisted lease (409 invalid_lease), so it is new work.
		{request("/v1/ack", `{"id":"unpersisted","token":"lease"}`, "", "agent_b"), "agent_b", "agent_b", false},
		{request("/v1/pull", "{}", "", "agent_a"), "", "agent_a", false},
		{request("/home/gates/gate/deny", deny, "session-a", ""), "session-a", "agent_a", true},
		{request("/home/gates/gate/deny", deny, "session-b", ""), "", "agent_b", false},
		{request("/home/gates/gate/deny", "decision=deny&csrf="+csrf("session-b", "gate"), "session-b", ""), "session-b", "agent_b", false},
		{request("/home/unpair", "agent=agent_a&target=agent_b", "session-a", ""), "session-a", "agent_a", true},
		{request("/home/unpair", "agent=agent_a&target=agent_b", "session-b", "", "Sec-Fetch-Site", "same-origin"), "session-b", "agent_b", true},
		{request("/home/agent-revoke", "agent=agent_a", "session-b", ""), "session-b", "agent_b", false},
	} {
		if got := cleanupCredential(c.r); got != c.credential {
			t.Fatalf("%s credential %q want %q", c.r.URL.Path, got, c.credential)
		}
		if st.cleanupTarget(c.r, c.principal, now) != c.own {
			t.Fatal("own cleanup target", c.r.URL.Path, c.principal)
		}
		// Review H2: the pre-DB channel choice applies the same own-record check to the committed snapshot.
		if got := s.cleanupOwner(c.r, c.credential); (got != "") != (c.credential != "" && c.own) {
			t.Fatalf("%s local cleanup owner %q", c.r.URL.Path, got)
		}
	}
	// An owner's agents share one fairness unit apart from the owner's own control.
	if s.cleanupOwner(request("/home/unpair", "agent=agent_a&target=agent_b", "session-a", ""), "session-a") != "agent_a" ||
		s.cleanupOwner(request("/v1/ack", `{"id":"persisted","token":"lease"}`, "", "agent_b"), "agent_b") != agentsUnit("agent_b") {
		t.Fatal("cleanup fairness unit")
	}
	// A session cookie does not authorize /v1, an agent credential does not authorize /owner, and expiry is honored.
	if s.cleanupOwner(request("/v1/key-revoke", revoke, "", "session-a"), "session-a") != "" ||
		s.cleanupOwner(request("/owner/gates/gate/deny", "", "", "agent_a"), "agent_a") != "" {
		t.Fatal("credential kind crossed paths")
	}
	st.Sessions[hashToken("session-a")].Seen = now.Add(-sessionIdle)
	s.live.Store(nil)
	s.remember(st, now, 1)
	if s.cleanupOwner(request("/home/unpair", "agent=agent_a&target=agent_b", "session-a", ""), "session-a") != "" {
		t.Fatal("idle session kept cleanup channel")
	}
}

// Review L2 and FIX-2 (Sol): a goroutine storing an earlier commit late cannot replace the newer committed snapshot,
// even at the same DB time, and a DB restored to a lower epoch still refreshes at its next, later commit.
func TestRememberKeepsNewestCommit(t *testing.T) {
	st, _, now := textFixture()
	s := &Service{}
	s.remember(st, now, 5)
	for _, older := range []struct {
		at    time.Time
		epoch int64
	}{{now.Add(-time.Millisecond), 4}, {now, 4}, {now, 5}} {
		s.remember(newState(), older.at, older.epoch)
		if s.live.Load().st != st {
			t.Fatal("older commit replaced newer snapshot", older)
		}
	}
	for _, newer := range []struct {
		at    time.Time
		epoch int64
	}{{now, 6}, {now.Add(time.Millisecond), 1}} {
		next := newState()
		s.remember(next, newer.at, newer.epoch)
		if s.live.Load().st != next {
			t.Fatal("newer commit not stored", newer)
		}
	}
}
