// Member connections bind a short approval to one owner, agent, client and new key; only confirmed PoP creates credentials.
package relay

import (
	"crypto/ed25519"
	"encoding/base64"
	"log"
	"sort"
	"time"
)

const connectionTTL = 10 * time.Minute
const supportedClient = "node-local"

// A revoked agent stays at least the frozen 24h revoke-metadata period, then the sweep deletes it with its keys and pairs.
// Member agent IDs are random and never issued again, so deleting the whole agent cannot reassign a (from,kid) (C1).
// Revoked keys of a live agent are never deleted, for the same reason.
// The key record cap is technical state protection, not a product quota: it holds the active key limit (3) plus
// revoked records, so churn cannot grow the single state row without bound. A key-saturated agent is replaced by a
// new agent. Agent count has no per-member cap (SAR-CONNECT-SIMPLE-001 L-01/L-02): the 24h new-agent rate and the
// service capacity boundary below bound the row instead.
const (
	revokedRetention = 24 * time.Hour
	agentKeyRecords  = 20
	// agentRate new agents per member in a rolling 24h, revoked ones included (L-04).
	agentRate = 20
	// The live-agent boundary of the single state row (L-03, D12): a warning log from agentWarn, new agents refused
	// at agentCeiling. Existing agents, keys and friends never depend on it.
	agentWarn    = 1000
	agentCeiling = 2000
)

type Connection struct {
	Owner, Agent, Client, Mode, Kid, Public, State string
	Exp                                            time.Time
	Device                                         bool // Client-started requests never put the private polling token in a browser.
	// A client-started attempt keeps the client's host hint and the approved name: its agent exists only once the
	// client saves its key (A-01). Refused and Retry record an L-03/L-04 refusal so the client can report it.
	Host    string    `json:",omitempty"`
	Name    string    `json:",omitempty"`
	Refused string    `json:",omitempty"`
	Retry   time.Time `json:",omitzero"`
}

// agentRefused is an L-03/L-04 refusal of a new agent; Retry is when agent_rate frees a slot.
type agentRefused struct {
	Code  string
	Retry time.Time
}

func (e agentRefused) Error() string { return e.Code }

func connectionBytes(token string, c *Connection) []byte {
	return []byte("KNOWSLINK-CONNECT\x00" + token + "\x00" + c.Owner + "\x00" + c.Agent + "\x00" + c.Client + "\x00" + c.Mode + "\x00" + c.Kid + "\x00" + c.Public)
}

func (st *State) liveAgents() int {
	n := 0
	for _, a := range st.Agents {
		if !a.Revoked && st.Owners[a.Owner] != nil && st.Owners[a.Owner].Active {
			n++
		}
	}
	return n
}

// agentLimit refuses a new agent for owner at the service capacity boundary or the 24h rate, else returns nil.
func (st *State) agentLimit(owner string, now time.Time) error {
	if st.liveAgents() >= agentCeiling {
		return agentRefused{Code: "service_capacity"}
	}
	var made []time.Time
	for _, a := range st.Agents {
		if a.Owner == owner && !a.Created.IsZero() && now.Sub(a.Created) < 24*time.Hour {
			made = append(made, a.Created)
		}
	}
	if len(made) < agentRate {
		return nil
	}
	sort.Slice(made, func(i, j int) bool { return made[i].Before(made[j]) })
	return agentRefused{"agent_rate", made[len(made)-agentRate].Add(24 * time.Hour)}
}

// addAgent stores an agent agentLimit admitted; every creation path goes through it so the rate sees each one.
func (st *State) addAgent(id string, a *Agent, now time.Time) {
	a.Created = now
	st.Agents[id] = a
	if live := st.liveAgents(); live >= agentWarn {
		log.Printf("agent capacity warning: live=%d warn=%d refuse=%d", live, agentWarn, agentCeiling)
	}
}

// open reports whether an attempt can still progress. An approved client attempt stays completable until the sweep
// removes it a day after its deadline, so a lost answer can still be saved (A-02, F-07).
func (c *Connection) open(now time.Time) bool {
	if c.Device && c.State == "approved" {
		return now.Before(c.Exp.Add(24 * time.Hour))
	}
	return now.Before(c.Exp) && c.State != "cancelled" && c.State != "consumed" && c.State != "refused"
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
	if token == "" || c == nil || !c.open(now) {
		return nil, fault("invalid_auth")
	}
	a := st.Agents[c.Agent]
	// The agent is made when its key is saved (A-01). Only an attempt approved by this flow carries Name; an older
	// attempt's agent already existed, so a missing one was revoked and swept and its ID is never made again (C1).
	reserved := c.Device && c.State == "approved" && a == nil && c.Name != ""
	if (!reserved && (a == nil || a.Owner != c.Owner || a.Revoked)) || st.Owners[c.Owner] == nil || !st.Owners[c.Owner].Active {
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
	if a == nil {
		if err := st.agentLimit(c.Owner, now); err != nil {
			return nil, err
		}
		a = &Agent{Name: c.Name, Owner: c.Owner, Keys: map[string]*Key{}}
		st.addAgent(c.Agent, a, now)
	}
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
	a.Client = c.Client
	c.State = "consumed"
	st.sweep(now)
	return map[string]string{"agent": c.Agent, "kid": c.Kid, "credential": credential, "state": "connected"}, nil
}
func (st *State) sweepConnections(now time.Time) {
	for id, c := range st.Connections {
		if !now.Before(c.Exp) {
			// An approved client attempt keeps its state so the key can still be saved; it goes with the 24h delete (A-02).
			if c.State != "consumed" && c.State != "cancelled" && c.State != "refused" && !(c.Device && c.State == "approved") {
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
	// Agent-create replay entries map to agent IDs (message entries map to message IDs) and go with their agent.
	if len(deleted) > 0 {
		for key, value := range st.Idempotency {
			if deleted[value] {
				delete(st.Idempotency, key)
			}
		}
	}
}
