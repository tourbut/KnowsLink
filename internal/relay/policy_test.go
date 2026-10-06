// Policy checks pin the agreed re-invite rules and the record saturation guidance (SAR-PUBLIC-AGENTS-001-POLICY).
package relay

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Saturation refusals name the real limit and a next step that can work, without quota, payment or exact-time promises.
func TestSaturationGuidance(t *testing.T) {
	st, now := connectionFixture()
	for i := 0; i < agentKeyRecords; i++ {
		st.Agents["agent_a"].Keys[fmt.Sprint("k", i)] = &Key{Revoked: i > 0, Changed: now}
	}
	for _, mode := range []string{"register", "rotate"} {
		if _, e := st.beginConnection("owner", "agent_a", supportedClient, mode, now); e == nil || e.Error() != "capacity" {
			t.Fatal("key records", e)
		}
		problem, create := connectLimit(st, "owner", "agent_a", mode)
		if problem != keyRecordsFull || !create {
			t.Fatal("key-full guidance", mode, problem, create)
		}
	}
	// No replacement fits: the refusal adds the actual agent limit and drops the create control.
	for i := 1; i < ownerAgentRecords; i++ {
		st.Agents[fmt.Sprint("agent_r", i)] = &Agent{Owner: "owner", Keys: map[string]*Key{}, Revoked: true, Changed: now}
	}
	problem, create := connectLimit(st, "owner", "agent_a", "register")
	if create || !strings.HasPrefix(problem, keyRecordsFull) || !strings.Contains(problem, agentLimits["records"]) {
		t.Fatal("key-full without replacement", problem, create)
	}
	if st.agentLimit("owner") != "records" {
		t.Fatal("owner records", st.agentLimit("owner"))
	}
	// Only the actual cleanup after the retention frees the record space; the old key stays refused throughout.
	st.Agents["agent_r1"].Keys["old"] = &Key{Revoked: true, Changed: now, Credential: hashToken("old-credential")}
	st.sweep(now.Add(revokedRetention - time.Second))
	if _, e := st.principal("old-credential", "agent"); e == nil {
		t.Fatal("revoked key during retention")
	}
	if st.agentLimit("owner") != "records" {
		t.Fatal("record space before retention")
	}
	st.sweep(now.Add(revokedRetention))
	if st.agentLimit("owner") != "" || st.Agents["agent_r1"] != nil || st.Agents["agent_a"] == nil {
		t.Fatal("cleanup", st.agentLimit("owner"))
	}
	if _, e := st.principal("old-credential", "agent"); e == nil {
		t.Fatal("revoked key after cleanup")
	}
	// Three active keys: register names the active key limit; the live agent's key records are untouched.
	st, now = connectionFixture()
	for i := 0; i < 3; i++ {
		st.Agents["agent_a"].Keys[fmt.Sprint("k", i)] = &Key{Changed: now}
	}
	if problem, create := connectLimit(st, "owner", "agent_a", "register"); create || !strings.Contains(problem, "활성 키는 agent당 3개") {
		t.Fatal("active keys", problem)
	}
	for i := 1; i < 5; i++ {
		st.Agents[fmt.Sprint("agent_", i)] = &Agent{Owner: "owner", Keys: map[string]*Key{}}
	}
	if st.agentLimit("owner") != "active" || !strings.Contains(agentProblem("active"), "회원당 5개") {
		t.Fatal("active agents", st.agentLimit("owner"))
	}
	if agentProblem("total") != safeProblem(fault("capacity")) {
		t.Fatal("service-wide limit")
	}
	for _, m := range []string{keyRecordsFull, agentLimits["active"], agentLimits["records"]} {
		for _, banned := range []string{"결제", "업그레이드", "요금", "정확히", "quota"} {
			if strings.Contains(m, banned) {
				t.Fatal(banned, m)
			}
		}
	}
	if !strings.Contains(keyRecordsFull, "기다려도 이 agent에 새 키 공간은 생기지 않습니다") || !strings.Contains(agentLimits["records"], "최소 24시간") {
		t.Fatal("record guidance")
	}
}

// Relationship policy (POLICY D02 PS-07): repeats keep the current invite, terminated relationships take a manual
// re-invite with a new generation and a new explicit acceptance, and no old decision or message reactivates it.
func TestRelationshipPolicy(t *testing.T) {
	st, now := connectionFixture()
	st.Owners["other"] = &Owner{Active: true}
	st.Agents["agent_b"] = &Agent{Owner: "other", Keys: map[string]*Key{}}
	id := pairID("agent_a", "agent_b")
	invite := func(owner, from, to string, at time.Time) (*Pair, error) {
		v, e := st.operateAs("invites", owner, "owner", command{Agent: from, Target: to}, at, false)
		p, _ := v.(*Pair)
		return p, e
	}
	decide := func(owner, decision string, at time.Time) error {
		_, e := st.operateAs("invite-decision", owner, "owner", command{Agent: "agent_a", Target: "agent_b", Decision: decision}, at, false)
		return e
	}
	pendingCount := func() int {
		n := 0
		for _, p := range st.Pairs {
			if p.State == "pending" {
				n++
			}
		}
		return n
	}
	message := func() *Message {
		return &Message{Receipt: Receipt{From: "agent_a", To: "agent_b"}, Kid: "k", Generation: st.Pairs[id].Generation}
	}
	st.Agents["agent_a"].Keys["k"] = &Key{Changed: now}
	st.Agents["agent_b"].Keys["k"] = &Key{Changed: now}

	p, e := invite("owner", "agent_a", "agent_b", now)
	if e != nil || p.State != "pending" || p.Generation != 1 {
		t.Fatal("first invite", e)
	}
	exp := p.Exp
	// Same-direction and opposite-direction repeats keep the pending: no new invite, generation, deadline or acceptance.
	for _, c := range [][3]string{{"owner", "agent_a", "agent_b"}, {"other", "agent_b", "agent_a"}} {
		if q, e := invite(c[0], c[1], c[2], now.Add(time.Hour)); e != nil || q != st.Pairs[id] || q.State != "pending" || q.Generation != 1 || !q.Exp.Equal(exp) || q.Recipient != "agent_b" || pendingCount() != 1 {
			t.Fatal("repeat changed pending", c, e)
		}
	}
	pendingMessage := message()
	if st.current(pendingMessage) {
		t.Fatal("pending allowed a message")
	}
	// Only the current recipient owner decides; the inviter's owner cannot accept for the recipient.
	if decide("owner", "accept", now) == nil || st.Pairs[id].State != "pending" {
		t.Fatal("inviter decided")
	}
	// The deadline boundary: accepted just before it, refused at it (the transaction sweeps first).
	saved, _ := json.Marshal(st)
	st.sweep(exp.Add(-time.Nanosecond))
	if decide("other", "accept", exp.Add(-time.Nanosecond)) != nil || st.Pairs[id].State != "active" {
		t.Fatal("accept before deadline")
	}
	st = newState()
	_ = json.Unmarshal(saved, st)
	st.sweep(exp)
	if decide("other", "accept", exp) == nil || st.Pairs[id].State != "expired" || pendingCount() != 0 {
		t.Fatal("accept at deadline")
	}
	// Expired: a manual re-invite is a new pending with a new deadline and generation.
	if p, e = invite("owner", "agent_a", "agent_b", exp); e != nil || p.Generation != 2 || !p.Exp.Equal(exp.Add(24*time.Hour)) {
		t.Fatal("re-invite after expiry", e)
	}
	if decide("other", "accept", exp) != nil {
		t.Fatal("explicit accept")
	}
	active := message()
	if !st.current(active) {
		t.Fatal("active refused current message")
	}
	// Active repeats keep the relationship, generation and active count from either side.
	for _, c := range [][3]string{{"owner", "agent_a", "agent_b"}, {"other", "agent_b", "agent_a"}} {
		if q, e := invite(c[0], c[1], c[2], exp); e != nil || q.State != "active" || q.Generation != 2 || len(st.Pairs) != 1 || pendingCount() != 0 {
			t.Fatal("active repeat", c, e)
		}
	}
	// Either owner's unpair ends it; the old acceptance and messages never carry into the next generation.
	for round, owner := range []string{"owner", "other"} {
		if _, e := st.operateAs("unpair", owner, "owner", command{Agent: "agent_a", Target: "agent_b"}, exp, false); e != nil || st.Pairs[id].State != "revoked" {
			t.Fatal("unpair", owner, e)
		}
		if st.current(active) || decide("other", "accept", exp) == nil || decide("other", "deny", exp) == nil {
			t.Fatal("ended relationship still usable", owner)
		}
		// The member who unpaired, or the other side, may re-invite manually; it waits for a new acceptance.
		from, to, inviter := "agent_a", "agent_b", "owner"
		if round == 1 {
			from, to, inviter = "agent_b", "agent_a", "other"
		}
		p, e = invite(inviter, from, to, exp)
		if e != nil || p.State != "pending" || p.Generation != int64(3+round*2) || st.current(active) || st.current(message()) {
			t.Fatal("re-invite after unpair", owner, e)
		}
		recipient := st.Agents[p.Recipient].Owner
		if _, e := st.operateAs("invite-decision", inviter, "owner", command{Agent: from, Target: to, Decision: "accept"}, exp, false); e == nil && inviter != recipient {
			t.Fatal("inviter accepted the re-invite")
		}
		if _, e := st.operateAs("invite-decision", recipient, "owner", command{Agent: from, Target: to, Decision: "accept"}, exp, false); e != nil {
			t.Fatal("re-accept", e)
		}
		if st.current(active) || !st.current(message()) || st.Pairs[id].Generation != int64(3+round*2) {
			t.Fatal("new generation", owner)
		}
		active = message()
		// Leave the next round a fresh generation to end.
		if round == 0 {
			if _, e := st.operateAs("unpair", "other", "owner", command{Agent: "agent_a", Target: "agent_b"}, exp, false); e != nil {
				t.Fatal(e)
			}
			if p, e = invite("owner", "agent_a", "agent_b", exp); e != nil || decide("other", "accept", exp) != nil || p.Generation != 4 {
				t.Fatal("round reset", e)
			}
			active = message()
		}
	}
	// Deny ends the invite and frees its pending slot; a late accept of the denied invite changes nothing.
	if _, e := st.operateAs("unpair", "owner", "owner", command{Agent: "agent_a", Target: "agent_b"}, exp, false); e != nil {
		t.Fatal(e)
	}
	p, _ = invite("owner", "agent_a", "agent_b", exp)
	denied := p.Generation
	if decide("other", "deny", exp) != nil || st.Pairs[id].State != "denied" || pendingCount() != 0 || decide("other", "accept", exp) == nil {
		t.Fatal("deny")
	}
	// A third owner cannot decide or unpair someone else's relationship.
	st.Owners["third"] = &Owner{Active: true}
	if p, _ = invite("owner", "agent_a", "agent_b", exp); decide("third", "accept", exp) == nil || p.Generation != denied+1 {
		t.Fatal("foreign decision")
	}
	if _, e := st.operateAs("unpair", "third", "owner", command{Agent: "agent_a", Target: "agent_b"}, exp, false); e == nil || st.Pairs[id].State != "pending" {
		t.Fatal("foreign unpair")
	}
	// Same owner: two agents still need an explicit acceptance.
	st.Agents["agent_c"] = &Agent{Owner: "owner", Keys: map[string]*Key{}}
	if q, e := invite("owner", "agent_a", "agent_c", exp); e != nil || q.State != "pending" {
		t.Fatal("same-owner auto accept", e)
	}
	// At the sender pending cap the next re-invite creates no pending or generation.
	if decide("other", "deny", exp) != nil {
		t.Fatal(e)
	}
	for i := 0; i < 9; i++ {
		st.Pairs[fmt.Sprint("fill", i)] = &Pair{A: "agent_a", B: "agent_c", State: "pending", Exp: exp.Add(time.Hour)}
	}
	gen := st.Pairs[id].Generation
	if _, e := invite("owner", "agent_a", "agent_b", exp); e == nil || e.Error() != "capacity" || st.Pairs[id].State != "denied" || st.Pairs[id].Generation != gen {
		t.Fatal("capacity re-invite", e)
	}
	// The state survives a restart; a revoked agent or an inactive owner cannot be revived by a re-invite.
	for k := range st.Pairs {
		if strings.HasPrefix(k, "fill") {
			delete(st.Pairs, k)
		}
	}
	raw, _ := json.Marshal(st)
	st = newState()
	_ = json.Unmarshal(raw, st)
	if st.Pairs[id].State != "denied" || st.Pairs[id].Generation != gen {
		t.Fatal("restart")
	}
	st.Owners["other"].Active = false
	if _, e := invite("owner", "agent_a", "agent_b", exp); e == nil {
		t.Fatal("inactive owner re-invited")
	}
	st.Owners["other"].Active = true
	st.Agents["agent_b"].Revoked = true
	if _, e := invite("owner", "agent_a", "agent_b", exp); e == nil {
		t.Fatal("revoked agent re-invited")
	}
	if _, e := invite("other", "agent_b", "agent_a", exp); e == nil {
		t.Fatal("revoked agent invited")
	}
}
