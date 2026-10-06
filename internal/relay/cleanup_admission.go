// Cleanup admission rests on an honored credential and the caller's own record, so rejected or anonymous cleanup is new work.
package relay

import (
	"crypto/hmac"
	"net/http"
	"strings"
	"time"
)

var crossOrigin = http.NewCrossOriginProtection()

// cleanupCredential returns the credential a cleanup request rests on, or "" when its path, origin, gate CSRF or body
// already makes it new work. It reads no state, so it also picks the local channel before any DB admission.
func cleanupCredential(r *http.Request) string {
	if !cleanupRequest(r) || crossOrigin.Check(r) != nil {
		return ""
	}
	token := ownerToken(r)
	if strings.HasPrefix(r.URL.Path, "/home") || strings.HasPrefix(r.URL.Path, "/auth") {
		token = readCookie(r, sessionCookie)
	}
	if strings.HasPrefix(r.URL.Path, "/v1/") && Strict(commandBody(r), &command{}) != nil {
		return ""
	}
	if gate := gatePath(r.URL.Path); gate != "" && !hmac.Equal([]byte(r.PostForm.Get("csrf")), []byte(csrf(token, gate))) {
		return ""
	}
	return token
}

// cleanupTarget mirrors each cleanup handler's ownership check. The handler stays the authority; this only picks the budget.
func (st *State) cleanupTarget(r *http.Request, principal string, now time.Time) bool {
	var c command
	_ = Strict(commandBody(r), &c)
	owner := principal
	if strings.HasPrefix(r.URL.Path, "/home") || strings.HasPrefix(r.URL.Path, "/auth") {
		m := st.Members[principal]
		if m == nil {
			return false
		}
		owner = m.Owner
		c = command{Agent: r.PostForm.Get("agent"), Target: r.PostForm.Get("target"), Kid: r.PostForm.Get("kid")}
	}
	own := func(agent string) bool { a := st.Agents[agent]; return a != nil && a.Owner == owner }
	p := st.Pairs[pairID(c.Agent, c.Target)]
	switch r.URL.Path {
	case "/v1/ack", "/v1/test/ack", "/v1/text/ack":
		m, err := st.leased(principal, c.ID, c.Token, now)
		return err == nil && (m.Receipt.Intent == publicTextIntent) == (r.URL.Path == "/v1/text/ack")
	case "/v1/owner-revoke":
		return st.Owners[principal] != nil
	case "/auth/logout", "/auth/logout-all":
		return true
	case "/v1/key-revoke", "/home/key-revoke":
		return own(c.Agent) && st.Agents[c.Agent].Keys[c.Kid] != nil
	case "/home/agent-revoke":
		return own(c.Agent)
	case "/v1/unpair", "/home/unpair":
		return p != nil && (own(p.A) || own(p.B))
	case "/home/invite-deny", "/home/invite-decision":
		return p != nil && own(p.Recipient)
	case "/home/cancel":
		conn := st.Connections[r.PostForm.Get("connection")]
		return conn != nil && conn.Owner == owner
	}
	g := st.Gates[gatePath(r.URL.Path)]
	return g != nil && g.Owner == owner
}

func gatePath(path string) string {
	for _, prefix := range []string{"/home/gates/", "/owner/gates/"} {
		if id, ok := strings.CutPrefix(path, prefix); ok {
			return strings.TrimSuffix(id, "/deny")
		}
	}
	return ""
}

// liveCredentials indexes the credentials the committed state still honors, so the pre-DB channel choice cannot be won by
// an invented token. Admission re-verifies inside the transaction.
func (st *State) liveCredentials() map[string]time.Time {
	live := map[string]time.Time{}
	active := func(owner string) bool { o := st.Owners[owner]; return o != nil && o.Active }
	for _, o := range st.Owners {
		if o.Active && o.Credential != "" {
			live[o.Credential] = time.Time{}
		}
	}
	for _, a := range st.Agents {
		if a.Revoked || !active(a.Owner) {
			continue
		}
		for _, k := range a.Keys {
			if !k.Revoked {
				live[k.Credential] = time.Time{}
				if a.Credential != "" {
					live[a.Credential] = time.Time{}
				}
			}
		}
	}
	for id, session := range st.Sessions {
		if m := st.Members[session.Member]; m != nil && m.Active && active(m.Owner) {
			exp := session.Created.Add(sessionAbsolute)
			if idle := session.Seen.Add(sessionIdle); idle.Before(exp) {
				exp = idle
			}
			live[id] = exp
		}
	}
	return live
}
func (s *Service) liveCredential(token string) bool {
	live := s.live.Load()
	if live == nil || token == "" {
		return false
	}
	exp, ok := (*live)[hashToken(token)]
	return ok && (exp.IsZero() || time.Now().Before(exp))
}
