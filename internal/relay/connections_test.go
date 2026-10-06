// Connection and quota checks cover key-bound credentials, expiry, owner isolation and exact capacity boundaries.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
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
func TestRevokedRecordRetention(t *testing.T) {
	st, now := connectionFixture()
	// Rotation churn: one active key plus revoked keys reach the record cap; only new keys are refused.
	for i := 0; i < agentKeyRecords; i++ {
		token, _, proof := prepared(t, st, now, "rotate", fmt.Sprint("key", i))
		st.Connections[hashToken(token)].State = "approved"
		if _, e := st.completeConnection(token, supportedClient, proof, now.Add(time.Duration(i)*time.Second)); e != nil {
			t.Fatal(i, e)
		}
	}
	a := st.Agents["agent_a"]
	if len(a.Keys) != agentKeyRecords || activeKeys(a) != 1 {
		t.Fatal("rotation records", len(a.Keys))
	}
	for _, mode := range []string{"register", "rotate"} {
		if _, e := st.beginConnection("owner", "agent_a", supportedClient, mode, now); e == nil || e.Error() != "capacity" {
			t.Fatal("saturated key records accepted", mode, e)
		}
	}
	key0 := a.Keys["key0"].Changed
	if _, e := st.operateAs("key-revoke", "owner", "owner", command{Agent: "agent_a", Kid: fmt.Sprint("key", agentKeyRecords-1)}, now.Add(time.Hour), false); e != nil {
		t.Fatal("revocation blocked while saturated", e)
	}
	if _, e := st.operateAs("key-revoke", "owner", "owner", command{Agent: "agent_a", Kid: "key0"}, now.Add(2*time.Hour), false); e != nil || a.Keys["key0"].Changed != key0 {
		t.Fatal("repeated revoke moved revocation time", e)
	}
	// A live agent keeps every revoked kid, so no kid is ever reassigned (C1); the cap stays.
	st.sweep(now.Add(30 * 24 * time.Hour))
	if len(a.Keys) != agentKeyRecords || a.Keys["key0"] == nil {
		t.Fatal("live agent lost revoked key records")
	}
	if _, e := st.beginConnection("owner", "agent_a", supportedClient, "register", now.Add(30*24*time.Hour)); e == nil {
		t.Fatal("key records grew past cap")
	}

	// Agent churn: the owner record cap holds active and revoked agents; the sweep frees it after 24h.
	st, now = connectionFixture()
	st.Owners["other"] = &Owner{Active: true}
	st.Agents["agent_b"] = &Agent{Owner: "other", Keys: map[string]*Key{}}
	st.Agents["agent_c"] = &Agent{Owner: "other", Keys: map[string]*Key{}}
	for _, target := range []string{"agent_b", "agent_c"} {
		st.Pairs[pairID("agent_a", target)] = &Pair{A: "agent_a", B: target, Inviter: "agent_a", Recipient: target, State: "active", Generation: 3}
	}
	st.Pairs[pairID("agent_b", "agent_c")] = &Pair{A: "agent_b", B: "agent_c", Inviter: "agent_b", Recipient: "agent_c", State: "revoked", Generation: 4}
	st.Agents["agent_a"].Revoked = true
	st.Agents["agent_a"].Changed = now
	for i := 1; i < ownerAgentRecords; i++ {
		st.Agents[fmt.Sprint("agent_x", i)] = &Agent{Owner: "owner", Keys: map[string]*Key{}, Revoked: i > 4, Changed: now.Add(time.Hour)}
	}
	if st.agentCapacity("owner") {
		t.Fatal("owner records exceeded")
	}
	// A revocation stored before Changed existed starts retention at the next sweep instead of being deleted.
	legacy := []byte(`{"Agents":{"agent_old":{"Owner":"owner","Keys":{},"Revoked":true}}}`)
	restored := newState()
	if e := json.Unmarshal(legacy, restored); e != nil {
		t.Fatal(e)
	}
	restored.sweep(now)
	if restored.Agents["agent_old"] == nil || !restored.Agents["agent_old"].Changed.Equal(now) {
		t.Fatal("legacy revoked agent lost retention")
	}
	raw, _ := json.Marshal(st)
	st = newState()
	if e := json.Unmarshal(raw, st); e != nil {
		t.Fatal(e)
	}
	st.sweep(now.Add(revokedRetention))
	if st.Agents["agent_a"] != nil || st.Pairs[pairID("agent_a", "agent_b")] != nil || st.Pairs[pairID("agent_a", "agent_c")] != nil {
		t.Fatal("revoked agent or its pairs retained after restart and 24h")
	}
	if p := st.Pairs[pairID("agent_b", "agent_c")]; p == nil || p.Generation != 4 {
		t.Fatal("live pair history lost its generation")
	}
	if !st.agentCapacity("owner") {
		t.Fatal("sweep did not free owner records")
	}
	for _, path := range []string{"invite-decision", "unpair"} {
		if _, e := st.operateAs(path, "other", "owner", command{Agent: "agent_a", Target: "agent_b", Decision: "deny"}, now, false); e == nil {
			t.Fatal("deleted pair accepted", path)
		}
	}
	// Defensive: a pair left pointing at a missing agent is refused, not dereferenced.
	st.Pairs[pairID("agent_gone", "agent_b")] = &Pair{A: "agent_gone", B: "agent_b", Recipient: "agent_b", State: "pending"}
	for _, path := range []string{"invite-decision", "unpair"} {
		if _, e := st.operateAs(path, "other", "owner", command{Agent: "agent_gone", Target: "agent_b", Decision: "accept"}, now, false); e == nil {
			t.Fatal("missing agent pair accepted", path)
		}
	}
	// The synthetic owner path obeys the same key record cap.
	st.Agents["agent_b"].Keys = map[string]*Key{}
	for i := 0; i < agentKeyRecords; i++ {
		st.Agents["agent_b"].Keys[fmt.Sprint("k", i)] = &Key{Revoked: true, Changed: now}
	}
	public, private, _ := ed25519.GenerateKey(rand.Reader)
	encoded := base64.RawURLEncoding.EncodeToString(public)
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, popBytes("other", "agent_b", "fresh", encoded)))
	if _, e := st.operateAs("keys", "other", "owner", command{Agent: "agent_b", Kid: "fresh", Public: encoded, Proof: proof}, now, false); e == nil || e.Error() != "capacity" {
		t.Fatal("synthetic key path exceeded records", e)
	}
}
func TestMemberPagesShowNextSteps(t *testing.T) {
	page := func(status int, name string, v map[string]any) string {
		w := httptest.NewRecorder()
		render(w, status, name, v)
		if w.Code != status {
			t.Fatal(w.Code)
		}
		return w.Body.String()
	}
	has := func(body string, want ...string) {
		t.Helper()
		for _, s := range want {
			if !strings.Contains(body, s) {
				t.Fatalf("missing %q in %s", s, body)
			}
		}
	}
	lacks := func(body, s string) {
		t.Helper()
		if strings.Contains(body, s) {
			t.Fatalf("unexpected %q in %s", s, body)
		}
	}
	// Refusals keep the safe problem text and add the next step it names.
	w := httptest.NewRecorder()
	refused(w, "내 agent", refusal{422, safeProblem(fault("reauth_required")), true})
	has(w.Body.String(), `role="alert"`, `action="/auth/reauth"`, `href="/home"`, "자기 홈으로 돌아가기")
	w = httptest.NewRecorder()
	refused(w, "내 agent", refusal{422, safeProblem(fault("unsupported_client")), false})
	has(w.Body.String(), `href="/home"`)
	lacks(w.Body.String(), "/auth/reauth")
	w = httptest.NewRecorder()
	refused(w, "내 agent", refusal{401, "로그인 세션이 유효하지 않습니다. 다시 로그인하세요.", false})
	has(w.Body.String(), `href="/"`, "로그인 화면으로 이동")
	// Cancel appears only while the connection can still be cancelled.
	for _, state := range []string{"waiting", "prepared", "approved", "consumed", "cancelled", "expired"} {
		body := page(200, "connection", map[string]any{"Title": "연결", "ConnectionID": "c", "Agent": "agent_a", "State": state})
		if open := state == "waiting" || state == "prepared" || state == "approved"; open != strings.Contains(body, `action="/home/cancel"`) {
			t.Fatal("cancel control", state)
		}
		has(body, `href="/home"`)
	}
	// Fingerprints wrap inside the card, and invite deadlines use the KST clock.
	exp := time.Date(2026, 10, 7, 4, 19, 59, 162233000, time.UTC)
	body := page(200, "home", map[string]any{"Title": "홈", "AgentDetails": []memberAgent{{ID: "agent_a", Status: "연결 완료", Keys: []memberKey{{"key_1", "SHA256:7_B9RdVzaKQSQRTfdqCvAAAAAAAAAAAAAAAAAAAAAAAAAAA", "활성"}}}},
		"Pairs": []memberPair{{&Pair{A: "agent_a", B: "agent_b", Inviter: "agent_a", Recipient: "agent_b", State: "pending", Generation: 1, Exp: exp}, true}, {&Pair{A: "agent_a", B: "agent_c", State: "active"}, false}}})
	has(body, "<code>SHA256:7_B9", "초대 기한: 2026-10-07 13:19:59 KST", "초대 기한: 없음", "overflow-wrap:anywhere", "input,select{display:block;width:100%")
	lacks(body, "+0000 UTC")
}
