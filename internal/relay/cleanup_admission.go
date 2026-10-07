// Cleanup admission rests on an honored credential and the caller's own record, so rejected, foreign or anonymous cleanup is new work.
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
	case "/home/invite-deny", "/home/invite-decision", "/v1/invite-decision":
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

// committed is the last state this process committed, with the credentials it still honors, for the pre-DB channel choice.
type committed struct {
	st    *State
	at    time.Time
	creds map[string]honored
}

// honored is the principal a credential resolves to; exp is zero unless the credential is a session.
type honored struct {
	kind, id string
	exp      time.Time
}

// remember keeps the newest commit. Transactions commit in row-lock order with non-decreasing DB time, so a slower
// goroutine storing an earlier commit cannot replace a later one.
func (s *Service) remember(st *State, at time.Time) {
	next := &committed{st, at, st.liveCredentials()}
	for {
		old := s.live.Load()
		if old != nil && at.Before(old.at) || s.live.CompareAndSwap(old, next) {
			return
		}
	}
}

// liveCredentials indexes the credentials the committed state still honors, so the pre-DB channel choice cannot be won by
// an invented token. Admission re-verifies inside the transaction.
func (st *State) liveCredentials() map[string]honored {
	live := map[string]honored{}
	active := func(owner string) bool { o := st.Owners[owner]; return o != nil && o.Active }
	for id, o := range st.Owners {
		if o.Active && o.Credential != "" {
			live[o.Credential] = honored{kind: "owner", id: id}
		}
	}
	for id, a := range st.Agents {
		if a.Revoked || !active(a.Owner) {
			continue
		}
		for _, k := range a.Keys {
			if !k.Revoked {
				live[k.Credential] = honored{kind: "agent", id: id}
				if a.Credential != "" {
					live[a.Credential] = honored{kind: "agent", id: id}
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
			live[id] = honored{"session", session.Member, exp}
		}
	}
	return live
}

// cleanupOwner mirrors requestBuckets against the last committed state, read only and without a transaction. It returns
// the owner whose own record this cleanup ends, or "" when the request is new work: an invented, foreign or lease-less
// cleanup never reaches the reserved local channel. A stale or missing snapshot only sends cleanup to the new channel;
// any local commit, at least the 1s retention sweep, refreshes it.
func (s *Service) cleanupOwner(r *http.Request, token string) string {
	c := s.live.Load()
	if c == nil || token == "" {
		return ""
	}
	h, ok := c.creds[hashToken(token)]
	now := time.Now()
	sessionPath := strings.HasPrefix(r.URL.Path, "/home") || strings.HasPrefix(r.URL.Path, "/auth")
	if !ok || !h.exp.IsZero() && !now.Before(h.exp) || sessionPath != (h.kind == "session") ||
		h.kind == "agent" && (bearer(r) != token || strings.HasPrefix(r.URL.Path, "/owner")) || !c.st.cleanupTarget(r, h.id, now) {
		return ""
	}
	switch h.kind {
	case "agent":
		return c.st.Agents[h.id].Owner
	case "session":
		return c.st.Members[h.id].Owner
	}
	return h.id
}
