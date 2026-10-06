// Cleanup admission classification: honored credential, same-origin form, gate CSRF, valid body and own record.
package relay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Review M1: cleanup admission needs an honored credential, same-origin form, gate CSRF, valid body and the caller's own record.
func TestCleanupAdmissionNeedsVerifiedOwnRecord(t *testing.T) {
	st, _, now := textFixture()
	st.Sessions[hashToken("session-a")] = &Session{Member: "agent_a", Created: now, Seen: now}
	st.Sessions[hashToken("session-b")] = &Session{Member: "agent_b", Created: now, Seen: now}
	st.Gates["gate"] = &Gate{Owner: "agent_a", State: "pending"}
	st.Agents["agent_c"].Revoked = true
	s := &Service{}
	live := st.liveCredentials()
	s.live.Store(&live)
	for token, want := range map[string]bool{"agent_a": true, "session-a": true, "agent_c": false, "invented": false, "": false} {
		if s.liveCredential(token) != want {
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
	}
}
