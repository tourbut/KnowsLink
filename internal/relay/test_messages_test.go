// Trial-message regression covers explicit enablement, strict text schema, signature, replay, TTL and one-use receive.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestTrialMessageSafety(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	pa, sa, _ := ed25519.GenerateKey(rand.Reader)
	pb, _, _ := ed25519.GenerateKey(rand.Reader)
	st := newState()
	st.Owners["a"] = &Owner{hashToken("owner-a"), true}
	st.Owners["b"] = &Owner{hashToken("owner-b"), true}
	st.Agents["agent_a"] = &Agent{Owner: "a", Credential: hashToken("token-a"), Keys: map[string]*Key{"key1": {Public: pa}}}
	st.Agents["agent_b"] = &Agent{Owner: "b", Credential: hashToken("token-b"), Keys: map[string]*Key{"key1": {Public: pb}}}
	st.Pairs[pairID("agent_a", "agent_b")] = &Pair{A: "agent_a", B: "agent_b", State: "active", Generation: 1}
	message := func(id, key, text string, exp time.Time) *Envelope {
		raw := wire(t, sa, id, "agent_a", "agent_b", "relay.test.message", key, "", map[string]any{"text": text}, exp)
		e, err := Parse(raw)
		if err != nil {
			t.Fatal(err)
		}
		return e
	}
	e := message(firstID, "trial-idempotent-01", "ping", now.Add(time.Minute))
	if _, err := st.ingest("agent_a", e, "", now); err == nil {
		t.Fatal("default held accepted trial")
	}
	st.TestAgents = map[string]bool{"agent_a": true, "agent_b": true}
	if _, err := st.ingest("agent_b", e, "", now); err == nil {
		t.Fatal("wrong principal")
	}
	if _, err := st.ingest("agent_a", e, "", now); err != nil {
		t.Fatal(err)
	}
	replay := message(secondID, "trial-idempotent-01", "ping", now.Add(2*time.Minute))
	value, err := st.ingest("agent_a", replay, "", now)
	if err != nil || value.(Receipt).ID != firstID {
		t.Fatal("duplicate receipt", err)
	}
	conflict := message(secondID, "trial-idempotent-01", "changed", now.Add(time.Minute))
	if _, err = st.ingest("agent_a", conflict, "", now); err == nil || err.Error() != "idempotency_conflict" {
		t.Fatal("duplicate conflict", err)
	}
	for _, expiry := range []time.Time{now, now.Add(301 * time.Second)} {
		if _, err = st.ingest("agent_a", message(secondID, "trial-expiry-key", "ping", expiry), "", now); err == nil {
			t.Fatal("invalid TTL")
		}
	}
	lease, err := st.operate("test-pull", "token-b", command{}, now)
	if err != nil {
		t.Fatal(err)
	}
	token := lease.(map[string]any)["lease_token"].(string)
	c := command{ID: firstID, Token: token}
	if _, err = st.operate("test-ack", "token-b", c, now); err == nil {
		t.Fatal("ACK before persist")
	}
	for _, operation := range []string{"test-persist", "test-ack", "test-claim"} {
		if _, err = st.operate(operation, "token-b", c, now); err != nil {
			t.Fatal(operation, err)
		}
	}
	if _, err = st.operate("test-claim", "token-a", c, now); err == nil {
		t.Fatal("cross-agent claim")
	}
	if _, err = st.operate("test-claim", "token-b", c, now); err == nil {
		t.Fatal("repeated claim")
	}
	if _, err = st.parentFor("agent_b", firstID, now); err == nil {
		t.Fatal("trial became business parent")
	}
	if len(st.Messages[firstID].Envelope) != 0 || len(st.Messages[firstID].Inbox) != 0 {
		t.Fatal("claimed trial payload retained")
	}
	st.Pairs[pairID("agent_a", "agent_b")].State = "revoked"
	if _, err = st.ingest("agent_a", e, "", now); err == nil {
		t.Fatal("revoked pair replay")
	}
}

func TestTrialSchemaAndConfiguration(t *testing.T) {
	for _, value := range []string{"agent_a", "agent_a,agent_a", "agent_a,AGENT_B", "agent_a,agent_b,agent_c"} {
		if _, err := TestAgentAllowlist(value); err == nil {
			t.Fatal("invalid allowlist", value)
		}
	}
	if result, err := TestAgentAllowlist(""); err != nil || len(result) != 0 {
		t.Fatal("default")
	}
	_, private, _ := ed25519.GenerateKey(rand.Reader)
	raw := wire(t, private, firstID, "agent_a", "agent_b", "relay.test.message", "trial-schema-key", "", map[string]any{"text": "ping"}, time.Now().Add(time.Minute))
	for _, change := range []func(map[string]any){
		func(m map[string]any) { m["body"] = map[string]any{"text": strings.Repeat("x", 4097)} },
		func(m map[string]any) { m["body"] = map[string]any{"text": " "} },
		func(m map[string]any) { m["body"] = map[string]any{"text": "ping", "execute": true} },
		func(m map[string]any) { m["deliver"] = "human" },
		func(m map[string]any) { m["reply_to"] = secondID },
		func(m map[string]any) { m["evidence"] = []any{map[string]any{"ref": "https://example.com"}} },
		func(m map[string]any) { m["ext"] = map[string]any{"test.foo": "effect"} },
	} {
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
		change(m)
		bad, _ := json.Marshal(m)
		if _, err := Parse(bad); err == nil {
			t.Fatal("open trial schema")
		}
	}
}
