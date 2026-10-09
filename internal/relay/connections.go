// Member connections bind a short approval to one owner, agent, client and new key; only confirmed PoP creates credentials.
package relay

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"
)

const connectionTTL = 10 * time.Minute
const supportedClient = "node-local"

// A revoked agent stays at least the frozen 24h revoke-metadata period, then the sweep deletes it with its keys and pairs.
// Member agent IDs are random and never issued again, so deleting the whole agent cannot reassign a (from,kid) (C1).
// Revoked keys of a live agent are never deleted, for the same reason.
// The record caps are technical state protection, not product quotas: they hold the active limits (agent 5, key 3)
// plus revoked records, so churn cannot grow the single state row without bound. New agents and keys are refused
// at a cap; revocation never adds a record and stays available. A key-saturated agent is replaced by a new agent.
const (
	revokedRetention  = 24 * time.Hour
	ownerAgentRecords = 10
	agentKeyRecords   = 20
)

type Connection struct {
	Owner, Agent, Client, Mode, Kid, Public, State string
	Exp                                            time.Time
	Device                                         bool // Client-started requests never put the private polling token in a browser.
}

func connectionBytes(token string, c *Connection) []byte {
	return []byte("KNOWSLINK-CONNECT\x00" + token + "\x00" + c.Owner + "\x00" + c.Agent + "\x00" + c.Client + "\x00" + c.Mode + "\x00" + c.Kid + "\x00" + c.Public)
}
func (st *State) agentCapacity(owner string) bool { return st.agentLimit(owner) == "" }

// agentLimit names the limit refusing a new agent for owner: "total", "active", "records", or "" when one fits.
func (st *State) agentLimit(owner string) string {
	total, own, records := 0, 0, 0
	for _, a := range st.Agents {
		if a.Owner == owner {
			records++
		}
		if !a.Revoked && st.Owners[a.Owner] != nil && st.Owners[a.Owner].Active {
			total++
			if a.Owner == owner {
				own++
			}
		}
	}
	switch {
	case total >= 200:
		return "total"
	case records >= ownerAgentRecords:
		return "records"
	case own >= 5:
		return "active"
	}
	return ""
}
func (st *State) pairCapacity(a, b string, active bool) bool {
	total, ownA, ownB := 0, 0, 0
	for _, p := range st.Pairs {
		if (active && p.State != "active") || (!active && p.State != "pending") {
			continue
		}
		left, right := st.Agents[p.A], st.Agents[p.B]
		if left == nil || right == nil || left.Revoked || right.Revoked || st.Owners[left.Owner] == nil || st.Owners[right.Owner] == nil || !st.Owners[left.Owner].Active || !st.Owners[right.Owner].Active {
			continue
		}
		total++
		if active {
			if left.Owner == a || right.Owner == a {
				ownA++
			}
			if left.Owner == b || right.Owner == b {
				ownB++
			}
		} else if left.Owner == a {
			ownA++
		}
	}
	if active {
		return total < 400 && ownA < 20 && ownB < 20
	}
	return total < 200 && ownA < 10
}

// revokeKeys keeps the first revocation time as each key's revoke metadata.
func revokeKeys(a *Agent, now time.Time) {
	for _, k := range a.Keys {
		if !k.Revoked {
			k.Revoked = true
			k.Changed = now
		}
	}
}
func activeKeys(a *Agent) int {
	n := 0
	for _, k := range a.Keys {
		if !k.Revoked {
			n++
		}
	}
	return n
}
func (st *State) beginConnection(owner, agent, client, mode string, now time.Time) (string, error) {
	a := st.Agents[agent]
	if a == nil || a.Owner != owner || a.Revoked {
		return "", fault("sender_not_allowed")
	}
	if client != supportedClient || (mode != "register" && mode != "rotate") {
		return "", fault("unsupported_client")
	}
	if (mode == "register" && activeKeys(a) >= 3) || len(a.Keys) >= agentKeyRecords {
		return "", fault("capacity")
	}
	for _, c := range st.Connections {
		if c.Agent == agent && c.State != "consumed" {
			c.State = "cancelled"
		}
	}
	token := randomToken()
	st.Connections[hashToken(token)] = &Connection{Owner: owner, Agent: agent, Client: client, Mode: mode, State: "waiting", Exp: now.Add(connectionTTL)}
	return token, nil
}
func (st *State) connection(token string, now time.Time) (*Connection, error) {
	c := st.Connections[hashToken(token)]
	if token == "" || c == nil || !now.Before(c.Exp) || c.State == "cancelled" || c.State == "consumed" {
		return nil, fault("invalid_auth")
	}
	a := st.Agents[c.Agent]
	if a == nil || a.Owner != c.Owner || a.Revoked || st.Owners[c.Owner] == nil || !st.Owners[c.Owner].Active {
		return nil, fault("invalid_auth")
	}
	return c, nil
}
func (st *State) prepareConnection(token, client, kid, public, proof string, now time.Time) (any, error) {
	c, err := st.connection(token, now)
	if err != nil {
		return nil, err
	}
	if client != c.Client || c.State != "waiting" || !kidPattern.MatchString(kid) || len(kid) > 128 {
		return nil, fault("sender_not_allowed")
	}
	candidate := *c
	candidate.Kid = kid
	candidate.Public = public
	key, e := base64.RawURLEncoding.DecodeString(public)
	sig, e2 := base64.RawURLEncoding.DecodeString(proof)
	if e != nil || e2 != nil || len(key) != 32 || !ed25519.Verify(key, connectionBytes(token, &candidate), sig) {
		return nil, fault("invalid_signature")
	}
	if st.Agents[c.Agent].Keys[kid] != nil {
		return nil, fault("key_exists")
	}
	c.Kid = kid
	c.Public = public
	c.State = "prepared"
	return map[string]string{"state": "prepared"}, nil
}
func (st *State) completeConnection(token, client, proof string, now time.Time) (any, error) {
	c, err := st.connection(token, now)
	if err != nil {
		return nil, err
	}
	if c.Client != client || c.State != "approved" {
		return nil, fault("sender_not_allowed")
	}
	public, _ := base64.RawURLEncoding.DecodeString(c.Public)
	sig, e := base64.RawURLEncoding.DecodeString(proof)
	if e != nil || !ed25519.Verify(public, connectionBytes(token, c), sig) {
		return nil, fault("invalid_signature")
	}
	a := st.Agents[c.Agent]
	if a.Keys[c.Kid] != nil {
		return nil, fault("key_exists")
	}
	if (c.Mode == "register" && activeKeys(a) >= 3) || len(a.Keys) >= agentKeyRecords {
		return nil, fault("capacity")
	}
	if c.Mode == "rotate" {
		revokeKeys(a, now)
	}
	credential := randomToken()
	a.Keys[c.Kid] = &Key{Public: public, Changed: now, Credential: hashToken(credential)}
	// A public connection never retains the synthetic agent-wide credential.
	a.Credential = ""
	c.State = "consumed"
	st.sweep(now)
	return map[string]string{"agent": c.Agent, "kid": c.Kid, "credential": credential, "state": "connected"}, nil
}
func (st *State) sweepConnections(now time.Time) {
	for id, c := range st.Connections {
		if !now.Before(c.Exp) {
			if c.State != "consumed" && c.State != "cancelled" {
				c.State = "expired"
			}
			if now.Sub(c.Exp) >= 24*time.Hour {
				delete(st.Connections, id)
			}
		}
	}
	for _, p := range st.Pairs {
		if p.State == "pending" {
			if p.Exp.IsZero() {
				p.Exp = now.Add(24 * time.Hour)
			}
			if !now.Before(p.Exp) {
				p.State = "expired"
			}
		}
		a, b := st.Agents[p.A], st.Agents[p.B]
		if a != nil && b != nil && (a.Revoked || b.Revoked || st.Owners[a.Owner] == nil || st.Owners[b.Owner] == nil || !st.Owners[a.Owner].Active || !st.Owners[b.Owner].Active) {
			p.State = "revoked"
		}
	}
	deleted := map[string]bool{}
	for id, a := range st.Agents {
		if a.Revoked && a.Changed.IsZero() {
			a.Changed = now // A revocation saved before its time was recorded starts the full retention now.
		}
		if a.Revoked && now.Sub(a.Changed) >= revokedRetention {
			delete(st.Agents, id)
			deleted[id] = true
		}
	}
	// Pairs go with their agent so no relationship path dereferences a deleted agent.
	// Pairs between live agents stay, so a new invite always gets a later generation.
	for id, p := range st.Pairs {
		if deleted[p.A] || deleted[p.B] {
			delete(st.Pairs, id)
		}
	}
}
