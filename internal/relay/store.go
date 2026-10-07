// Postgres shared state serializes authorization, inbox, idempotency, and claims across relay processes.
package relay

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tourbut/KnowsLink/internal/database"
)

type Owner struct {
	Credential string
	Active     bool
}
type Key struct {
	Public     []byte
	Revoked    bool
	Changed    time.Time
	Credential string
}
type Agent struct {
	Owner      string
	Credential string
	Keys       map[string]*Key
	Revoked    bool
	Changed    time.Time
}
type Pair struct {
	A, B, Inviter, Recipient, State string
	Generation                      int64
	Exp                             time.Time
}
type Receipt struct {
	ID       string    `json:"id"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Intent   string    `json:"intent"`
	Digest   string    `json:"digest"`
	Exp      time.Time `json:"exp"`
	Accepted time.Time `json:"accepted_at"`
	State    string    `json:"transport"`
}
type Message struct {
	Receipt         Receipt
	Key, Kid        string
	Parent          string
	Deliver         string
	Envelope        json.RawMessage `json:"Envelope,omitempty"`
	Inbox           json.RawMessage `json:"Inbox,omitempty"`
	Generation      int64
	Attempts        int
	LeaseGeneration int64
	LeaseToken      string
	LeaseUntil      time.Time
	Persisted       bool
	Claimed         bool
	ClaimToken      string
	Completion      string
	ReplyID         string
}
type Gate struct {
	ID, Owner, Parent, Digest, From, To, Policy, State string
	Generation                                         int64
	Exp                                                time.Time
	Consumed                                           bool
}
type State struct {
	GoogleAttempts map[string]*googleAttempt
	Owners         map[string]*Owner
	Agents         map[string]*Agent
	Pairs          map[string]*Pair
	Messages       map[string]*Message
	Idempotency    map[string]string
	Gates          map[string]*Gate
	Members        map[string]*Member
	Identities     map[string]string
	Sessions       map[string]*Session
	Challenges     map[string]*Challenge
	Rates          map[string][]time.Time
	HTTP           map[string]admission
	Connections    map[string]*Connection
	TestAgents     map[string]bool `json:"-"`
}

func newState() *State {
	return &State{GoogleAttempts: map[string]*googleAttempt{}, Owners: map[string]*Owner{}, Agents: map[string]*Agent{}, Pairs: map[string]*Pair{}, Messages: map[string]*Message{}, Idempotency: map[string]string{}, Gates: map[string]*Gate{},
		Members: map[string]*Member{}, Identities: map[string]string{}, Sessions: map[string]*Session{}, Challenges: map[string]*Challenge{}, Rates: map[string][]time.Time{}, Connections: map[string]*Connection{}, HTTP: map[string]admission{}}
}

type Service struct {
	httpOnce           sync.Once
	httpNew, httpClean chan struct{}
	live               atomic.Pointer[committed]
	reading            chan struct{} // one snapshot refresh read at a time
	readAt             time.Time     // start of the last refresh read, guarded by reading
	cleaning           sync.Map      // fairness units holding the local cleanup channel
	orphans            sync.Map      // shared admission tokens whose finish transaction failed
	Pool               *pgxpool.Pool
	TestAgents         map[string]bool
	// SyntheticSignup keeps the local /v1/owners fixture; public members require verified identity.
	SyntheticSignup bool
	Google          *GoogleLogin
	Mail            Mailer
	ClientIPHeader  string
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
func pairID(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "/" + b
}
func (s *Service) transaction(ctx context.Context, operation func(*State, time.Time) (any, error)) (any, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fault("unavailable")
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := database.New(tx)
	// ponytail: one global row lock bounds local MVP throughput; normalize entities before public scale.
	row, err := q.LockRelay(ctx)
	if err != nil {
		return nil, fault("unavailable")
	}
	now := row.Now.Time.UTC()
	if !row.Now.Valid || !row.Clock.Valid || now.Before(row.Clock.Time) {
		return nil, fault("abnormal_clock")
	}
	state := newState()
	if err = json.Unmarshal(row.Data, state); err != nil {
		return nil, fault("unavailable")
	}
	state.TestAgents = s.TestAgents
	state.sweep(now)
	orphans := s.reclaim(state)
	value, opErr := operation(state, now)
	// Even a rejected request commits only cleanup and its spent request budget, never partial operation state.
	if opErr != nil {
		rates := state.Rates
		state = newState()
		if err = json.Unmarshal(row.Data, state); err != nil {
			return nil, fault("unavailable")
		}
		state.TestAgents = s.TestAgents
		state.sweep(now)
		s.reclaim(state)
		state.Rates = rates
	}
	raw, err := json.Marshal(state)
	if err != nil {
		return nil, fault("unavailable")
	}
	changed, err := q.SaveRelay(ctx, database.SaveRelayParams{Data: raw, Clock: row.Now, Epoch: row.Epoch})
	if err != nil || changed != 1 {
		return nil, fault("unavailable")
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fault("unavailable")
	}
	for _, token := range orphans {
		s.orphans.Delete(token)
	}
	s.remember(state, now, row.Epoch+1)
	return value, opErr
}

// reclaim drops shared admissions this process failed to finish, at its next commit (the 1s sweep at the latest).
func (s *Service) reclaim(st *State) (tokens []string) {
	s.orphans.Range(func(token, _ any) bool {
		delete(st.HTTP, token.(string))
		tokens = append(tokens, token.(string))
		return true
	})
	return tokens
}
func (st *State) current(m *Message) bool {
	if m.Receipt.Intent == "relay.test.message" && (!st.TestAgents[m.Receipt.From] || !st.TestAgents[m.Receipt.To]) {
		return false
	}
	a, b := st.Agents[m.Receipt.From], st.Agents[m.Receipt.To]
	if a == nil || b == nil || a.Revoked || b.Revoked {
		return false
	}
	ownerA, ownerB := st.Owners[a.Owner], st.Owners[b.Owner]
	if ownerA == nil || ownerB == nil || !ownerA.Active || !ownerB.Active {
		return false
	}
	recipientKey := false
	for _, key := range b.Keys {
		if !key.Revoked {
			recipientKey = true
			break
		}
	}
	if !recipientKey {
		return false
	}
	if m.Receipt.Intent == publicTextIntent && (!st.publicOwner(a.Owner) || !st.publicOwner(b.Owner)) {
		return false
	}
	key := a.Keys[m.Kid]
	if key == nil || key.Revoked {
		return false
	}
	if m.Receipt.Intent == publicTextIntent && m.Parent != "" {
		parent := st.Messages[m.Parent]
		return parent != nil && parent.Parent == "" && parent.Receipt.Intent == publicTextIntent && parent.Generation == m.Generation && st.current(parent)
	}
	if m.Receipt.Intent == "relay.result" {
		parent := st.Messages[m.Parent]
		return parent != nil && parent.Receipt.Intent != "relay.result" && parent.Receipt.Intent != "relay.approval.request" && st.current(parent)
	}
	if m.Receipt.Intent == "relay.approval.request" {
		g := st.Gates[m.Receipt.ID]
		if g == nil {
			return false
		}
		parent := st.Messages[g.Parent]
		return parent != nil && st.current(parent)
	}
	pair := st.Pairs[pairID(m.Receipt.From, m.Receipt.To)]
	return pair != nil && pair.State == "active" && pair.Generation == m.Generation
}
func (st *State) sweep(now time.Time) {
	for token, a := range st.HTTP {
		if !now.Before(a.Exp) {
			delete(st.HTTP, token)
		}
	}
	st.sweepIdentity(now)
	st.sweepConnections(now)
	for id, m := range st.Messages {
		if now.Sub(m.Receipt.Accepted) >= 24*time.Hour {
			delete(st.Idempotency, m.Receipt.From+"/"+m.Key)
			delete(st.Messages, id)
			delete(st.Gates, id)
			continue
		}
		if !st.current(m) {
			if m.Receipt.State != "delivered" {
				m.Receipt.State = "failed:revoked"
			}
			m.Envelope = nil
			m.Inbox = nil
			m.LeaseToken = ""
			m.ClaimToken = ""
		}
		if !now.Before(m.Receipt.Exp) {
			if m.Completion == "" && m.Claimed {
				m.Completion = "failed:expired"
			}
			if m.Receipt.State != "delivered" {
				m.Receipt.State = "failed:expired"
			}
			m.Envelope = nil
			m.Inbox = nil
			m.LeaseToken = ""
			m.ClaimToken = ""
		}
		if m.Receipt.State == "leased" && !now.Before(m.LeaseUntil) {
			m.LeaseToken = ""
			m.Persisted = false
			if m.Attempts >= 3 {
				m.Receipt.State = "failed:max_attempts"
				m.Envelope = nil
				m.Inbox = nil
			} else {
				m.Receipt.State = "queued"
			}
		}
	}
	for _, g := range st.Gates {
		parent := st.Messages[g.Parent]
		h := st.Messages[g.ID]
		if !now.Before(g.Exp) {
			g.State = "expired"
			continue
		}
		if parent != nil && h != nil && (!st.current(parent) || !st.current(h)) {
			g.State = "revoked"
			continue
		}
		if parent == nil || h == nil || len(parent.Envelope) == 0 {
			if g.State == "pending" || (g.State == "approved" && !g.Consumed) {
				g.State = "unavailable"
			}
			continue
		}
	}
}
func (st *State) principal(token, kind string) (string, error) {
	hash := hashToken(token)
	if token == "" {
		return "", fault("invalid_auth")
	}
	if kind == "owner" {
		for id, o := range st.Owners {
			if o.Active && o.Credential == hash {
				return id, nil
			}
		}
	} else {
		for id, a := range st.Agents {
			o := st.Owners[a.Owner]
			if o == nil || !o.Active || a.Revoked {
				continue
			}
			for _, k := range a.Keys {
				if !k.Revoked && (k.Credential == hash || (a.Credential == hash && a.Credential != "")) {
					return id, nil
				}
			}
		}
	}
	return "", fault("invalid_auth")
}
func require(condition bool, code string) error {
	if !condition {
		return fault(code)
	}
	return nil
}
func policy(intent string) string {
	if intent == "schedule.query" {
		return "schedule.query: disclosure policy absent; deny even after gate approve"
	}
	return "schedule.commit: human gate required; non-executable stub even after approve"
}
func (st *State) parentRouting(agent, id string) (*Message, error) {
	m := st.Messages[id]
	// The one parent boundary for authorize, gate-consume, H, and R: claims stored before Deliver existed fail closed (C1).
	if m == nil || m.Receipt.To != agent || m.Deliver != "agent" || m.Receipt.Intent == "relay.result" || m.Receipt.Intent == "relay.approval.request" || m.Receipt.Intent == "relay.test.message" || m.Receipt.Intent == publicTextIntent || m.Receipt.State != "delivered" || !m.Claimed || !st.current(m) {
		return nil, fault("sender_not_allowed")
	}
	return m, nil
}
func (st *State) parentFor(agent, id string, now time.Time) (*Message, error) {
	m, err := st.parentRouting(agent, id)
	if err != nil {
		return nil, err
	}
	if !now.Before(m.Receipt.Exp) {
		return nil, fault("expired")
	}
	return m, nil
}
func (st *State) leaseMessage(agent string, now time.Time) (any, error) {
	return st.leaseMessageFor(agent, now, false)
}
func (st *State) leaseMessageFor(agent string, now time.Time, testOnly bool) (any, error) {
	return st.leaseMatching(agent, now, testOnly, false)
}
func (st *State) leaseMatching(agent string, now time.Time, testOnly, publicOnly bool) (any, error) {
	var selected *Message
	for _, m := range st.Messages {
		if (m.Receipt.Intent == publicTextIntent) != publicOnly {
			continue
		}
		if testOnly && m.Receipt.Intent != "relay.test.message" {
			continue
		}
		if m.Receipt.To == agent && m.Receipt.State == "queued" && m.Deliver == "agent" && st.current(m) && now.Before(m.Receipt.Exp) {
			if selected == nil || m.Receipt.Accepted.Before(selected.Receipt.Accepted) {
				selected = m
			}
		}
	}
	if selected == nil {
		return nil, nil
	}
	m := selected
	m.Attempts++
	m.LeaseGeneration++
	m.LeaseToken = randomToken()
	m.LeaseUntil = now.Add(30 * time.Second)
	if m.Receipt.Exp.Before(m.LeaseUntil) {
		m.LeaseUntil = m.Receipt.Exp
	}
	m.Receipt.State = "leased"
	return map[string]any{"envelope": m.Envelope, "lease_token": m.LeaseToken, "lease_until": m.LeaseUntil, "generation": m.LeaseGeneration, "attempts": m.Attempts}, nil
}
func (st *State) leased(agent, id, token string, now time.Time) (*Message, error) {
	m := st.Messages[id]
	if m == nil || m.Receipt.To != agent || m.Deliver != "agent" || m.Receipt.State != "leased" || token == "" || m.LeaseToken != token || !now.Before(m.LeaseUntil) || !now.Before(m.Receipt.Exp) || !st.current(m) {
		return nil, fault("invalid_lease")
	}
	return m, nil
}
func (st *State) ingest(agent string, e *Envelope, claim string, now time.Time) (any, error) {
	if e.Intent == "relay.test.message" && (!st.TestAgents[e.From] || !st.TestAgents[e.To]) {
		return nil, fault("sender_not_allowed")
	}
	a := st.Agents[e.From]
	if agent != e.From || a == nil {
		return nil, fault("sender_not_allowed")
	}
	key := a.Keys[e.Sig.Kid]
	if key == nil || key.Revoked {
		return nil, fault("invalid_auth")
	}
	if !e.Verify(key.Public) {
		return nil, fault("invalid_signature")
	}
	b := st.Agents[e.To]
	if b == nil || st.Owners[b.Owner] == nil || !st.Owners[b.Owner].Active {
		return nil, fault("sender_not_allowed")
	}
	var generation int64
	var parent *Message
	if e.Intent == "relay.approval.request" {
		var err error
		parent, err = st.parentRouting(agent, e.ReplyTo)
		if err != nil {
			return nil, err
		}
		if e.From != e.To || e.To != parent.Receipt.To || e.Body["request_digest"] != parent.Receipt.Digest {
			return nil, fault("sender_not_allowed")
		}
		generation = parent.Generation
	} else {
		// No owner inbox serves direct human delivery, so an agent credential must never process it (C1).
		if e.Deliver != "agent" {
			return nil, fault("sender_not_allowed")
		}
		pair := st.Pairs[pairID(e.From, e.To)]
		if pair == nil || pair.State != "active" || e.From == e.To {
			return nil, fault("human_invite_required")
		}
		generation = pair.Generation
		if e.Intent == "relay.result" {
			var err error
			parent, err = st.parentRouting(agent, e.ReplyTo)
			if err != nil {
				return nil, err
			}
			if e.From != parent.Receipt.To || e.To != parent.Receipt.From || generation != parent.Generation {
				return nil, fault("sender_not_allowed")
			}
			if e.Body["status"] == "done" {
				return nil, fault("disclosure_denied")
			}
		}
	}
	digest := e.Digest()
	scope := e.From + "/" + e.Key
	if id, ok := st.Idempotency[scope]; ok {
		m := st.Messages[id]
		if m.Receipt.Digest != digest {
			return nil, fault("idempotency_conflict")
		}
		if m.Generation != generation || !st.current(m) {
			return nil, fault("sender_not_allowed")
		}
		return m.Receipt, nil
	}
	if parent != nil && (claim == "" || parent.ClaimToken != claim) {
		return nil, fault("sender_not_allowed")
	}
	exp, _ := instant(e.Exp)
	if !exp.After(now) {
		return nil, fault("expired")
	}
	if exp.After(now.Add(300 * time.Second)) {
		return nil, fault("ttl_too_long")
	}
	if st.Messages[e.ID] != nil {
		return nil, fault("id_collision")
	}
	if parent != nil && exp.After(parent.Receipt.Exp) {
		return nil, fault("expired")
	}
	if !st.messageCapacity(e.Intent == "relay.approval.request") {
		return nil, fault("capacity")
	}
	if e.Intent == "relay.approval.request" {
		if len(parent.Envelope) == 0 {
			return nil, fault("unavailable")
		}
		for _, g := range st.Gates {
			if g.Parent == e.ReplyTo {
				return nil, fault("duplicate_gate")
			}
		}
		st.Gates[e.ID] = &Gate{ID: e.ID, Owner: b.Owner, Parent: e.ReplyTo, Digest: digest, From: parent.Receipt.From, To: parent.Receipt.To, Policy: policy(parent.Receipt.Intent), State: "pending", Generation: generation, Exp: exp}
		st.Gates[e.ID].Digest = parent.Receipt.Digest
	}
	raw, err := Canonical(e.Raw)
	if err != nil {
		return nil, fault("invalid_schema")
	}
	m := &Message{Receipt: Receipt{e.ID, e.From, e.To, e.Intent, digest, exp, now, "queued"}, Key: e.Key, Kid: e.Sig.Kid, Envelope: raw, Generation: generation, Parent: e.ReplyTo, Deliver: e.Deliver}
	if !st.current(m) {
		return nil, fault("sender_not_allowed")
	}
	st.Messages[e.ID] = m
	st.Idempotency[scope] = e.ID
	if e.Intent == "relay.result" {
		if parent.Completion != "" {
			return nil, fault("duplicate_result")
		}
		parent.Completion = fmt.Sprint(e.Body["status"])
		parent.Envelope = nil
		parent.Inbox = nil
		parent.ClaimToken = ""
	}
	return m.Receipt, nil
}
