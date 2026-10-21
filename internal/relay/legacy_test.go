// Legacy state regression: claims issued before Message.Deliver existed must never reopen agent processing (C1).
package relay

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
	"time"
)

const legacyFixture = "testdata/legacy_claims.json"
const legacyHuman, legacyGated, legacyAgent = firstID, secondID, "0199a3f2-4c10-7a11-8b22-3344556677cc"

var legacyNow = time.Date(2026, 10, 3, 10, 0, 0, 0, time.UTC)

type legacyState struct {
	State              *State
	PrivateA, PrivateB []byte
}

// seedLegacy runs only on the a6a10c7 source, whose State had no Deliver route:
// RELAY_LEGACY_SEED=1 go test ./internal/relay -run TestLegacyClaims
func seedLegacy(t *testing.T) {
	publicA, privateA, _ := ed25519.GenerateKey(rand.Reader)
	publicB, privateB, _ := ed25519.GenerateKey(rand.Reader)
	st := newState()
	st.Owners["owner_a"] = &Owner{hashToken("legacy_owner_a"), true}
	st.Owners["owner_b"] = &Owner{hashToken("legacy_owner_b"), true}
	st.Agents["agent_a"] = &Agent{Owner: "owner_a", Credential: hashToken("legacy_agent_a"), Keys: map[string]*Key{"key1": {Public: publicA}}}
	st.Agents["agent_b"] = &Agent{Owner: "owner_b", Credential: hashToken("legacy_agent_b"), Keys: map[string]*Key{"key1": {Public: publicB}}}
	st.Pairs[pairID("agent_a", "agent_b")] = &Pair{A: "agent_a", B: "agent_b", State: "active", Generation: 1}
	exp := legacyNow.Add(240 * time.Second)
	send := func(private ed25519.PrivateKey, raw []byte, deliver, claim string) {
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
		m["deliver"] = deliver
		unsigned, _ := json.Marshal(m)
		e, err := Parse(unsigned)
		if err != nil {
			t.Fatal(err)
		}
		m["sig"].(map[string]any)["value"] = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, e.SigningBytes()))
		signed, _ := json.Marshal(m)
		if e, err = Parse(signed); err != nil {
			t.Fatal(err)
		}
		if _, err = st.ingest(e.From, e, claim, legacyNow); err != nil {
			t.Fatal(err)
		}
	}
	for i, id := range []string{legacyHuman, legacyGated, legacyAgent} {
		deliver := "human"
		if id == legacyAgent {
			deliver = "agent"
		}
		send(privateA, wire(t, privateA, id, "agent_a", "agent_b", "schedule.query", "legacy-request-0"+string(rune('1'+i)), "", queryBody(), exp), deliver, "")
		lease, err := st.leaseMessage("agent_b", legacyNow)
		if err != nil || lease == nil {
			t.Fatal("legacy lease failed", err)
		}
		token := lease.(map[string]any)["lease_token"].(string)
		for _, path := range []string{"persist", "ack", "claim"} {
			if _, err = st.operate(path, "legacy_agent_b", command{ID: id, Token: token}, legacyNow); err != nil {
				t.Fatal(path, err)
			}
		}
	}
	gated := st.Messages[legacyGated]
	send(privateB, wire(t, privateB, gateID, "agent_b", "agent_b", "relay.approval.request", "legacy-approval-01", legacyGated, map[string]any{"reason": "judgment_required", "request_digest": gated.Receipt.Digest}, exp), "human", gated.ClaimToken)
	// Same field changes as the owner POST approve handler at a6a10c7.
	h := st.Messages[gateID]
	h.Receipt.State, h.Persisted, h.Envelope, h.Inbox = "delivered", true, nil, nil
	st.Gates[gateID].State = "approved"
	raw, _ := json.MarshalIndent(legacyState{st, privateA, privateB}, "", " ")
	if err := os.WriteFile(legacyFixture, append(raw, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLegacyClaimsCannotReachParentBoundaries(t *testing.T) {
	if os.Getenv("RELAY_LEGACY_SEED") == "1" {
		seedLegacy(t)
		return
	}
	raw, err := os.ReadFile(legacyFixture)
	if err != nil || bytes.Contains(raw, []byte(`"Deliver"`)) {
		t.Fatal("legacy fixture must predate Message.Deliver", err)
	}
	load := func() (*State, ed25519.PrivateKey, ed25519.PrivateKey) {
		var fixture legacyState
		if err := json.Unmarshal(raw, &fixture); err != nil || fixture.State == nil {
			t.Fatal("invalid legacy fixture", err)
		}
		return fixture.State, fixture.PrivateA, fixture.PrivateB
	}
	denied := func(name string, err error) {
		t.Helper()
		if err == nil || err.Error() != "sender_not_allowed" {
			t.Fatalf("%s: legacy claim reused: %v", name, err)
		}
	}
	ingest := func(st *State, private ed25519.PrivateKey, id, to, intent, key, parent string, body map[string]any) error {
		e, err := Parse(wire(t, private, id, "agent_b", to, intent, key, parent, body, legacyNow.Add(time.Minute)))
		if err != nil {
			t.Fatal(err)
		}
		_, err = st.ingest("agent_b", e, st.Messages[parent].ClaimToken, legacyNow)
		return err
	}
	// The human and the unrecorded agent route are both missing Deliver, so both fail closed.
	for _, id := range []string{legacyHuman, legacyAgent} {
		st, _, privateB := load()
		m := st.Messages[id]
		if !m.Claimed || m.ClaimToken == "" || m.Receipt.State != "delivered" {
			t.Fatal("fixture is not a legacy claim", id)
		}
		_, err := st.operate("authorize", "legacy_agent_b", command{ID: id, Claim: m.ClaimToken}, legacyNow)
		denied("authorize "+id, err)
		denied("result "+id, ingest(st, privateB, "0199a3f2-4c10-7a11-8b22-3344556677dd", "agent_a", "relay.result", "legacy-result-0001", id, map[string]any{"status": "denied"}))
		denied("approval "+id, ingest(st, privateB, "0199a3f2-4c10-7a11-8b22-3344556677ee", "agent_b", "relay.approval.request", "legacy-approval-02", id, map[string]any{"reason": "judgment_required", "request_digest": m.Receipt.Digest}))
		if m.Completion != "" || len(st.Gates) != 1 || m.ClaimToken == "" {
			t.Fatal("legacy parent changed")
		}
	}
	st, _, _ := load()
	_, err = st.operate("gate-consume", "legacy_agent_b", command{ID: gateID, Claim: st.Messages[legacyGated].ClaimToken}, legacyNow)
	denied("gate-consume", err)
	if st.Gates[gateID].Consumed {
		t.Fatal("legacy approved gate consumed")
	}
	// A new agent delivery on the same loaded state keeps the full claim/authorize/result path.
	st, privateA, privateB := load()
	const newID = "0199a3f2-4c10-7a11-8b22-3344556677ff"
	e, err := Parse(wire(t, privateA, newID, "agent_a", "agent_b", "schedule.query", "legacy-request-04", "", queryBody(), legacyNow.Add(time.Minute)))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = st.ingest("agent_a", e, "", legacyNow); err != nil {
		t.Fatal(err)
	}
	lease, err := st.leaseMessage("agent_b", legacyNow)
	if err != nil || lease == nil {
		t.Fatal("new agent lease failed", err)
	}
	token := lease.(map[string]any)["lease_token"].(string)
	for _, path := range []string{"persist", "ack", "claim"} {
		if _, err = st.operate(path, "legacy_agent_b", command{ID: newID, Token: token}, legacyNow); err != nil {
			t.Fatal(path, err)
		}
	}
	if _, err = st.operate("authorize", "legacy_agent_b", command{ID: newID, Claim: st.Messages[newID].ClaimToken}, legacyNow); err != nil {
		t.Fatal("new agent authorize", err)
	}
	if err = ingest(st, privateB, "0199a3f2-4c10-7a11-8b22-3344556677dd", "agent_a", "relay.result", "legacy-result-0001", newID, map[string]any{"status": "denied"}); err != nil || st.Messages[newID].Completion != "denied" {
		t.Fatal("new agent result", err)
	}
}
