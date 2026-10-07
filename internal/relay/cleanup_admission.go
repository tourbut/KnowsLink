// Cleanup admission rests on an honored credential and the caller's own record, so rejected, foreign or anonymous cleanup is new work.
package relay

import (
	"context"
	"crypto/hmac"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/tourbut/KnowsLink/internal/database"
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
		return err == nil && m.Persisted && (m.Receipt.Intent == publicTextIntent) == (r.URL.Path == "/v1/text/ack")
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

// committed is the newest relay state this process committed or read, with the credentials it still honors, for the
// pre-DB channel choice. at and epoch are the commit's DB time and row epoch.
type committed struct {
	st    *State
	at    time.Time
	epoch int64
	creds map[string]honored
}

// honored is the principal a credential resolves to; exp is zero unless the credential is a session.
type honored struct {
	kind, id string
	exp      time.Time
}

// remember keeps the newest commit. Transactions commit in row-lock order with non-decreasing DB time and an epoch that
// grows by one, so (time, epoch) orders commits even at the same DB time, and a slower goroutine storing an earlier
// commit cannot replace a later one. Time leads so a DB restored to a lower epoch still refreshes at its next commit.
func (s *Service) remember(st *State, at time.Time, epoch int64) {
	next := &committed{st, at, epoch, st.liveCredentials()}
	for {
		old := s.live.Load()
		if old != nil && (at.Before(old.at) || at.Equal(old.at) && epoch <= old.epoch) || s.live.CompareAndSwap(old, next) {
			return
		}
	}
}

// refresh reads the committed relay state when no read has started since the request arrived, so a cleanup that another
// process's commit or a restart made provable is judged on current records instead of waiting for the 1s sweep. One
// read runs at a time and covers every request that arrived before it started; the read takes no lock or slot, so
// invented credentials add at most one concurrent read per process, and waiters hold only their own connection.
func (s *Service) refresh(ctx context.Context, arrived time.Time) {
	select {
	case s.reading <- struct{}{}:
		defer func() { <-s.reading }()
	case <-ctx.Done():
		return
	}
	if !s.readAt.Before(arrived) {
		return
	}
	started := time.Now()
	row, err := database.New(s.Pool).ReadRelay(ctx)
	st := newState()
	if err != nil || json.Unmarshal(row.Data, st) != nil {
		return
	}
	st.TestAgents = s.TestAgents
	s.remember(st, row.Clock.Time.UTC(), row.Epoch)
	s.readAt = started
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
// the fairness unit whose own record this cleanup ends, or "" when the request is new work: an invented, foreign or
// lease-less cleanup never reaches the reserved local channel. boundedHTTP refreshes a stale or missing snapshot once
// before it treats a cleanup credential as new work.
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
		return agentsUnit(c.st.Agents[h.id].Owner)
	case "session":
		return c.st.Members[h.id].Owner
	}
	return h.id
}

// agentsUnit is the cleanup fairness unit of an owner's agents. Their ACKs share one slot apart from the owner's own
// revoke, unpair, deny and logout, so an agent repeating ACKs cannot keep its owner from revoking it.
func agentsUnit(owner string) string { return owner + "/agents" }
