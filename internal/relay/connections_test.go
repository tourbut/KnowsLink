// Connection and quota checks cover key-bound credentials, expiry, owner isolation and exact capacity boundaries.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func connectionFixture() (*State, time.Time) {
	st := newState()
	now := time.Now().UTC()
	st.Owners["owner"] = &Owner{Active: true}
	st.Agents["agent_a"] = &Agent{Owner: "owner", Keys: map[string]*Key{}}
	return st, now
}
func prepared(t *testing.T, st *State, now time.Time, mode, kid string) (string, ed25519.PrivateKey, string) {
	t.Helper()
	token, e := st.beginConnection("owner", "agent_a", supportedClient, mode, now)
	if e != nil {
		t.Fatal(e)
	}
	public, private, _ := ed25519.GenerateKey(rand.Reader)
	c := *st.Connections[hashToken(token)]
	c.Kid = kid
	c.Public = base64.RawURLEncoding.EncodeToString(public)
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, &c)))
	if _, e = st.prepareConnection(token, supportedClient, kid, c.Public, proof, now); e != nil {
		t.Fatal(e)
	}
	return token, private, proof
}
func TestConnectionApprovalAndKeyCredentials(t *testing.T) {
	st, now := connectionFixture()
	if _, e := st.beginConnection("other", "agent_a", supportedClient, "register", now); e == nil {
		t.Fatal("foreign owner")
	}
	if _, e := st.beginConnection("owner", "agent_a", "unknown", "register", now); e == nil {
		t.Fatal("unsupported client")
	}
	token, _, proof := prepared(t, st, now, "register", "key1")
	if _, e := st.completeConnection(token, supportedClient, proof, now); e == nil {
		t.Fatal("activation before confirmation")
	}
	c := st.Connections[hashToken(token)]
	c.State = "approved"
	if _, e := st.completeConnection(token, "other", proof, now); e == nil {
		t.Fatal("client substitution")
	}
	v, e := st.completeConnection(token, supportedClient, proof, now)
	if e != nil {
		t.Fatal(e)
	}
	old := v.(map[string]string)["credential"]
	if _, e := st.principal(old, "owner"); e == nil {
		t.Fatal("agent became owner")
	}
	if id, e := st.principal(old, "agent"); e != nil || id != "agent_a" {
		t.Fatal(id, e)
	}
	if _, e := st.completeConnection(token, supportedClient, proof, now); e == nil {
		t.Fatal("grant replay")
	}
	for i := 2; i <= 3; i++ {
		token, _, proof = prepared(t, st, now, "register", fmt.Sprintf("key%d", i))
		st.Connections[hashToken(token)].State = "approved"
		if _, e := st.completeConnection(token, supportedClient, proof, now); e != nil {
			t.Fatal(e)
		}
	}
	if _, e := st.beginConnection("owner", "agent_a", supportedClient, "register", now); e == nil {
		t.Fatal("fourth active key")
	}
	token, _, proof = prepared(t, st, now, "rotate", "key4")
	st.Connections[hashToken(token)].State = "approved"
	v, e = st.completeConnection(token, supportedClient, proof, now)
	if e != nil {
		t.Fatal(e)
	}
	if activeKeys(st.Agents["agent_a"]) != 1 {
		t.Fatal("rotation left old keys")
	}
	if _, e := st.principal(old, "agent"); e == nil {
		t.Fatal("revoked credential accepted")
	}
	fresh := v.(map[string]string)["credential"]
	if _, e := st.principal(fresh, "agent"); e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(st)
	loaded := newState()
	if e = json.Unmarshal(raw, loaded); e != nil {
		t.Fatal(e)
	}
	if _, e := loaded.principal(old, "agent"); e == nil {
		t.Fatal("restart restored old credential")
	}
	loaded.Agents["agent_a"].Keys["key4"].Revoked = true
	if _, e := loaded.principal(fresh, "agent"); e == nil {
		t.Fatal("selected key revocation")
	}
}
func TestConnectionFailureExpiryAndCancellation(t *testing.T) {
	for _, state := range []string{"cancelled", "expired", "consumed"} {
		t.Run(state, func(t *testing.T) {
			st, now := connectionFixture()
			token, _, proof := prepared(t, st, now, "register", "key1")
			c := st.Connections[hashToken(token)]
			c.State = state
			if state == "expired" {
				c.Exp = now
			}
			if _, e := st.completeConnection(token, supportedClient, proof, now); e == nil {
				t.Fatal("invalid connection accepted")
			}
			if len(st.Agents["agent_a"].Keys) != 0 {
				t.Fatal("failure connected key")
			}
		})
	}
	st, now := connectionFixture()
	token, _, proof := prepared(t, st, now, "register", "key1")
	c := st.Connections[hashToken(token)]
	c.State = "approved"
	if _, e := st.completeConnection(token, supportedClient, proof+"x", now); e == nil {
		t.Fatal("bad PoP")
	}
	st.Agents["agent_a"].Owner = "other"
	if _, e := st.completeConnection(token, supportedClient, proof, now); e == nil {
		t.Fatal("owner substitution")
	}
}
func TestAgentAndPairCapacity(t *testing.T) {
	st, now := connectionFixture()
	for i := 1; i < 5; i++ {
		st.Agents[fmt.Sprintf("agent_%d", i)] = &Agent{Owner: "owner", Keys: map[string]*Key{}}
	}
	if st.agentCapacity("owner") {
		t.Fatal("owner sixth agent")
	}
	st.Agents["agent_1"].Revoked = true
	if !st.agentCapacity("owner") {
		t.Fatal("revoke did not release slot")
	}
	st.Owners["other"] = &Owner{Active: true}
	for i := 5; i < 201; i++ {
		st.Agents[fmt.Sprintf("agent_%d", i)] = &Agent{Owner: "other", Keys: map[string]*Key{}}
	}
	if st.agentCapacity("owner") {
		t.Fatal("global 201st agent")
	}
	st, now = connectionFixture()
	st.Owners["other"] = &Owner{Active: true}
	st.Agents["agent_b"] = &Agent{Owner: "other", Keys: map[string]*Key{}}
	for i := 0; i < 200; i++ {
		p := &Pair{A: "agent_a", B: "agent_b", State: "pending", Exp: now.Add(time.Hour)}
		st.Pairs[fmt.Sprint(i)] = p
	}
	if st.pairCapacity("new", "other", false) {
		t.Fatal("global pending cap")
	}
	st.Pairs = map[string]*Pair{}
	for i := 0; i < 10; i++ {
		st.Pairs[fmt.Sprint(i)] = &Pair{A: "agent_a", B: "agent_b", State: "pending", Exp: now.Add(time.Hour)}
	}
	if st.pairCapacity("owner", "other", false) {
		t.Fatal("sender pending cap")
	}
	for _, p := range st.Pairs {
		p.Exp = now
	}
	st.sweep(now)
	if !st.pairCapacity("owner", "other", false) {
		t.Fatal("pending expiry did not free capacity")
	}
	for i := 0; i < 20; i++ {
		st.Pairs[fmt.Sprint(i)] = &Pair{A: "agent_a", B: "agent_b", State: "active"}
	}
	if st.pairCapacity("other", "other", true) {
		t.Fatal("owner active cap")
	}
	if _, e := st.operateAs("unpair", "owner", "owner", command{Agent: "agent_a", Target: "agent_b"}, now, false); e == nil {
		t.Fatal("unrelated map entry accepted")
	}
	st.Pairs[pairID("agent_a", "agent_b")] = &Pair{A: "agent_a", B: "agent_b", Inviter: "agent_a", Recipient: "agent_b", State: "pending", Generation: 1, Exp: now.Add(time.Hour)}
	if _, e := st.operateAs("invite-decision", "other", "owner", command{Agent: "agent_a", Target: "agent_b", Decision: "accept"}, now, false); e == nil {
		t.Fatal("accept exceeded cap")
	}
	if st.Pairs[pairID("agent_a", "agent_b")].State != "pending" {
		t.Fatal("rejected acceptance mutated pair")
	}
	if _, e := st.operateAs("invite-decision", "other", "owner", command{Agent: "agent_a", Target: "agent_b", Decision: "deny"}, now, false); e != nil {
		t.Fatal(e)
	}
	st.Pairs = map[string]*Pair{}
	for i := 0; i < 400; i++ {
		st.Pairs[fmt.Sprint(i)] = &Pair{A: "agent_a", B: "agent_b", State: "active"}
	}
	if st.pairCapacity("new", "new", true) {
		t.Fatal("global active cap")
	}
}
