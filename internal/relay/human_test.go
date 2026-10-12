// Synthetic members exercise friendship epochs, signed and owner text, retention, revocation and legacy migration.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func humanFixture(t *testing.T) (*State, map[string]string, map[string]ed25519.PrivateKey, time.Time) {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	st := newState()
	st.HumanModel = 1
	ids := map[string]string{}
	keys := map[string]ed25519.PrivateKey{}
	for _, person := range []string{"a", "b", "c"} {
		v := st.signIn(emailIssuer, person+"@example.invalid", person+"@example.invalid", "", "", now)
		if v.Problem != "" {
			t.Fatal(v.Problem)
		}
		id, _, e := st.session(v.Session, now)
		if e != nil {
			t.Fatal(e)
		}
		ids[person] = id
		m := st.Members[id]
		m.Name = "합성 회원 " + person
		for n := 1; n <= 2; n++ {
			aid := fmt.Sprintf("agent_%s%d", person, n)
			pub, key, e := ed25519.GenerateKey(rand.Reader)
			if e != nil {
				t.Fatal(e)
			}
			keys[aid] = key
			st.Agents[aid] = &Agent{Name: "도우미", Owner: m.Owner, Keys: map[string]*Key{"key1": {Public: pub, Credential: hashToken(aid)}}}
		}
	}
	return st, ids, keys, now
}

func acceptHumans(t *testing.T, st *State, ids map[string]string, now time.Time) *Friend {
	t.Helper()
	st.requestFriend(ids["a"], st.Members[ids["b"]].FriendCode, now)
	f := st.Friends[pairID(ids["a"], ids["b"])]
	if f == nil {
		t.Fatal("request missing")
	}
	if e := st.decideFriend(ids["b"], f.ID, "accept", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	return f
}

func TestHumanNamesAndIdentity(t *testing.T) {
	for _, raw := range []string{"", " \t ", strings.Repeat("가", 41), "a\x00b", "a\u202eb", "a\u2066b", string([]byte{255})} {
		if _, ok := validName(raw); ok {
			t.Fatalf("accepted unsafe name %q", raw)
		}
	}
	if got, ok := validName("  " + strings.Repeat("가", 40) + "  "); !ok || len([]rune(got)) != 40 {
		t.Fatal("unicode trim/length")
	}
	st, ids, _, now := humanFixture(t)
	before := st.Members[ids["a"]]
	owner := before.Owner
	before.Name = "같은 이름"
	again := st.signIn(emailIssuer, "a@example.invalid", "a@example.invalid", "", "", now)
	id, _, e := st.session(again.Session, now)
	if e != nil || id != ids["a"] || st.Members[id].Owner != owner || st.Members[id].Name != "같은 이름" {
		t.Fatal("identity changed")
	}
	other := st.signIn(googleIssuer, "different-sub", "a@example.invalid", "", "", now)
	oid, _, _ := st.session(other.Session, now)
	if oid == id {
		t.Fatal("merged by email")
	}
}

func TestHumanFriendTransitionsAndDiscovery(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	a, b := ids["a"], ids["b"]
	st.requestFriend(a, st.Members[b].FriendCode, now)
	f := st.Friends[pairID(a, b)]
	exp := f.Exp
	generation := f.Generation
	st.requestFriend(a, st.Members[b].FriendCode, now.Add(time.Hour))
	st.requestFriend(b, st.Members[a].FriendCode, now.Add(time.Hour))
	if f.State != "pending" || !f.Exp.Equal(exp) || f.Generation != generation || len(st.Friends) != 1 {
		t.Fatal("duplicate/reverse request changed consent or expiry")
	}
	if _, _, ok := st.authority("agent_a1", "agent_b1"); ok {
		t.Fatal("pending authorizes")
	}
	if e := st.decideFriend(a, f.ID, "accept", generation, now); e == nil {
		t.Fatal("sender accepted")
	}
	if e := st.decideFriend(ids["c"], f.ID, "deny", generation, now); e == nil {
		t.Fatal("third party decided")
	}
	if e := st.decideFriend(b, f.ID, "accept", generation, now); e != nil {
		t.Fatal(e)
	}
	if e := st.decideFriend(b, f.ID, "accept", generation, now); e != nil || len(st.Friends) != 1 {
		t.Fatal("duplicate accept")
	}
	if _, _, ok := st.authority("agent_a2", "agent_b2"); !ok {
		t.Fatal("future/all agent consent")
	}
	if _, _, ok := st.authority("agent_a1", "agent_c1"); ok {
		t.Fatal("transitive grant")
	}
	if e := st.decideFriend(a, f.ID, "block", generation, now); e != nil {
		t.Fatal(e)
	}
	st.Pairs[pairID("agent_a1", "agent_b1")] = &Pair{A: "agent_a1", B: "agent_b1", State: "active", Generation: 1}
	if _, _, ok := st.authority("agent_a1", "agent_b1"); ok {
		t.Fatal("legacy bypasses block")
	}
	st.requestFriend(b, st.Members[a].FriendCode, now)
	if f.State == "pending" || !f.BlockA {
		t.Fatal("blocked request restored friend")
	}
	if e := st.decideFriend(b, f.ID, "unblock", generation, now); e == nil {
		t.Fatal("other unblocked")
	}
	if e := st.decideFriend(a, f.ID, "unblock", generation, now); e != nil {
		t.Fatal(e)
	}
	if _, _, ok := st.authority("agent_a1", "agent_b1"); ok {
		t.Fatal("unblock restored legacy")
	}
	st.requestFriend(a, st.Members[b].FriendCode, now)
	if f.Generation != generation+1 {
		t.Fatal("generation not advanced")
	}
	if e := st.decideFriend(b, f.ID, "accept", generation, now); e == nil {
		t.Fatal("old accept")
	}
	if e := st.decideFriend(b, f.ID, "deny", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	st.requestFriend(a, st.Members[b].FriendCode, now)
	if e := st.decideFriend(a, f.ID, "cancel", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	st.requestFriend(a, st.Members[b].FriendCode, now)
	st.sweep(now.Add(24 * time.Hour))
	if f.State != "expired" {
		t.Fatal("pending lifetime")
	}

	// Missing and blocked codes leave the same safe outgoing-only state, never target identity.
	for _, code := range []string{"missing-code", st.Members[a].FriendCode} {
		fresh, fi, _, at := humanFixture(t)
		fresh.requestFriend(fi["a"], code, at)
		view, e := fresh.humanView(fi["a"], &Session{Verified: at}, "friends", httptest.NewRequest("GET", "/home/friends", nil), at)
		if e != nil {
			t.Fatal(e)
		}
		for _, r := range view["Requests"].([]humanFriend) {
			if r.Name != "요청 대기" || r.Incoming {
				t.Fatal("code discovery reveals identity")
			}
		}
	}
}

func TestHumanConversationLifetimeReplayAndRevocation(t *testing.T) {
	st, ids, keys, now := humanFixture(t)
	f := acceptHumans(t, st, ids, now)
	if _, e := st.sendOwnerText(ids["a"], "agent_b1", "agent_a1", "forged", "owner-attempt-0001", now); e == nil {
		t.Fatal("foreign sender")
	}
	var first *Message
	for i := 0; i < 6; i++ {
		person, from, to := "a", "agent_a1", "agent_b1"
		if i%2 == 1 {
			person, from, to = "b", to, from
		}
		at := now.Add(time.Duration(i) * 181 * time.Second)
		st.sweep(at)
		m, e := st.sendOwnerText(ids[person], from, to, fmt.Sprintf("합성 대화 %d", i), fmt.Sprintf("owner-attempt-%04d", i), at)
		if e != nil {
			t.Fatal("continuous send", i, e)
		}
		if i == 0 {
			first = m
		}
		leased, e := st.leaseMatching(to, at, false, true)
		if e != nil || leased == nil {
			t.Fatal("lease", e)
		}
		lease := leased.(map[string]any)
		c := command{ID: m.Receipt.ID, Token: lease["lease_token"].(string)}
		if _, e = st.operateAs("ack", to, "agent", c, at, false); e == nil {
			t.Fatal("ACK before durable copy")
		}
		if _, e = st.operateAs("persist", to, "agent", c, at, false); e != nil {
			t.Fatal(e)
		}
		if _, e = st.operateAs("ack", to, "agent", c, at, false); e != nil {
			t.Fatal(e)
		}
		if m.Receipt.State != "delivered" || m.Conversation == nil || m.Conversation.Authorship != "owner" || len(m.Envelope) != 0 {
			t.Fatal("receipt/copy/wire contract")
		}
	}
	if len(st.historical(st.Members[ids["a"]].Owner, "agent_a1", "agent_b1", now.Add(time.Hour))) != 6 {
		t.Fatal("history lost with TTL")
	}
	m, e := st.sendOwnerText(ids["a"], "agent_a1", "agent_b1", "4096", "owner-attempt-replay", now.Add(time.Hour))
	if e != nil {
		t.Fatal(e)
	}
	old, e := st.sendOwnerText(ids["a"], "agent_a1", "agent_b1", "4096", "owner-attempt-replay", now.Add(time.Hour+time.Second))
	if e != nil || old.Receipt.ID != m.Receipt.ID || !old.Receipt.Exp.Equal(m.Receipt.Exp) {
		t.Fatal("replay extended TTL or duplicated")
	}
	for _, change := range [][2]string{{"changed", "agent_b1"}, {"4096", "agent_b2"}} {
		if _, e := st.sendOwnerText(ids["a"], "agent_a1", change[1], change[0], "owner-attempt-replay", now.Add(time.Hour)); e == nil || e.Error() != "idempotency_conflict" {
			t.Fatal("semantic replay conflict", e)
		}
	}
	for _, body := range []string{" \n ", strings.Repeat("한", 1366), string([]byte{255})} {
		if _, e := st.sendOwnerText(ids["a"], "agent_a1", "agent_b1", body, "owner-invalid-test", now); e == nil {
			t.Fatal("invalid body")
		}
	}
	if _, e := st.sendOwnerText(ids["a"], "agent_a1", "agent_b1", strings.Repeat("한", 1365)+"x", "owner-max-body-0001", now); e != nil {
		t.Fatal("4096 boundary", e)
	}
	signed := signedText(remoteUUID(now), "agent_b1", "agent_a1", "signed-human-0001", "", "signed synthetic", now.Add(180*time.Second), keys["agent_b1"])
	if _, e := st.ingestPublicText("agent_b1", signed, now); e != nil {
		t.Fatal(e)
	}
	if st.Messages[signed.ID].Conversation.Authorship != "agent" {
		t.Fatal("signed author mislabeled")
	}
	if e := st.decideFriend(ids["a"], f.ID, "remove", f.Generation, now.Add(time.Hour)); e != nil {
		t.Fatal(e)
	}
	if st.current(m) || st.eventPending("agent_b1", now.Add(time.Hour)) != "" {
		t.Fatal("revoked delivery/events")
	}
	if _, e := st.sendOwnerText(ids["a"], "agent_a1", "agent_b1", "4096", "owner-attempt-replay", now.Add(time.Hour)); e == nil {
		t.Fatal("revoked replay")
	}
	if len(st.historical(st.Members[ids["c"]].Owner, "agent_a1", "agent_b1", now)) != 0 {
		t.Fatal("third party history")
	}
	before := first.Conversation.PersonB
	st.Members[ids["b"]].Name = "changed current name"
	if first.Conversation.PersonB != before {
		t.Fatal("historical identity changed")
	}
	st.requestFriend(ids["a"], st.Members[ids["b"]].FriendCode, now.Add(time.Hour))
	if e := st.decideFriend(ids["b"], f.ID, "accept", f.Generation, now.Add(time.Hour)); e != nil {
		t.Fatal(e)
	}
	if st.current(m) {
		t.Fatal("old epoch revived")
	}
	st.sweep(now.Add(24 * time.Hour))
	if st.Messages[first.Receipt.ID] != nil {
		t.Fatal("24h body retained")
	}
	if st.Messages[m.Receipt.ID] == nil {
		t.Fatal("independent retention shortened")
	}
}

func TestHumanLegacyMigrationAndExecutionBoundary(t *testing.T) {
	st, ids, keys, now := humanFixture(t)
	st.HumanModel = 0
	st.Pairs[pairID("agent_a1", "agent_b1")] = &Pair{A: "agent_a1", B: "agent_b1", State: "active", Generation: 7}
	st.Pairs[pairID("agent_a2", "agent_b2")] = &Pair{A: "agent_a2", B: "agent_b2", State: "pending", Generation: 4}
	credential := st.Agents["agent_a1"].Keys["key1"].Credential
	raw, e := json.Marshal(st)
	if e != nil {
		t.Fatal(e)
	}
	restored := newState()
	if e = json.Unmarshal(raw, restored); e != nil {
		t.Fatal(e)
	}
	if e = restored.migrateHuman(now); e != nil {
		t.Fatal(e)
	}
	st = restored
	if st.HumanModel != 1 || st.Members[ids["a"]].SelfConnect || st.Agents["agent_a1"].Keys["key1"].Credential != credential {
		t.Fatal("migration changed identity/key/consent")
	}
	if _, _, ok := st.authority("agent_a1", "agent_b1"); !ok {
		t.Fatal("exact legacy lost")
	}
	if _, _, ok := st.authority("agent_a2", "agent_b2"); ok {
		t.Fatal("pending expanded")
	}
	if _, _, ok := st.authority("agent_a1", "agent_a2"); ok {
		t.Fatal("self scope silently expanded")
	}
	if _, e = st.operateAs("invites", st.Members[ids["a"]].Owner, "owner", command{Agent: "agent_a1", Target: "agent_b2"}, now, false); e == nil {
		t.Fatal("legacy new grant")
	}
	businessRaw := wire(t, keys["agent_a1"], remoteUUID(now), "agent_a1", "agent_b1", "schedule.query", "human-business-0001", "", queryBody(), now.Add(180*time.Second))
	env, e := Parse(businessRaw)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = st.ingest("agent_a1", env, "", now); e != nil {
		t.Fatal(e)
	}
	msg := st.Messages[env.ID]
	if msg.Conversation != nil {
		t.Fatal("business copied into chat")
	}
	lease, _ := st.leaseMatching("agent_b1", now, false, false)
	c := command{ID: env.ID, Token: lease.(map[string]any)["lease_token"].(string)}
	_, _ = st.operateAs("persist", "agent_b1", "agent", c, now, false)
	_, _ = st.operateAs("ack", "agent_b1", "agent", c, now, false)
	claim, e := st.operateAs("claim", "agent_b1", "agent", command{ID: env.ID}, now, false)
	if e != nil {
		t.Fatal(e)
	}
	f := acceptHumans(t, st, ids, now)
	if st.current(msg) {
		t.Fatal("friend acceptance inherited old legacy claim")
	}
	if _, e = st.operateAs("authorize", "agent_b1", "agent", command{ID: env.ID, Claim: claim.(map[string]string)["claim"]}, now, false); e == nil {
		t.Fatal("old claim authorizes")
	}
	if e = st.decideFriend(ids["a"], f.ID, "block", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	st.Members[ids["a"]].Active = false
	if _, _, ok := st.authority("agent_a1", "agent_b1"); ok {
		t.Fatal("inactive member revives legacy")
	}
	st.HumanModel = 2
	if e = st.migrateHuman(now); e == nil {
		t.Fatal("unknown schema allowed")
	}
}

func TestHumanGateRevocationAndReaccept(t *testing.T) {
	for _, boundary := range []string{"remove", "block", "key", "member"} {
		t.Run(boundary, func(t *testing.T) {
			st, ids, keys, now := humanFixture(t)
			f := acceptHumans(t, st, ids, now)
			send := func(id, from, to, intent, key, parent, claim string, body map[string]any) *Envelope {
				t.Helper()
				e, err := Parse(wire(t, keys[from], id, from, to, intent, key, parent, body, now.Add(time.Minute)))
				if err != nil {
					t.Fatal(err)
				}
				if _, err = st.ingest(from, e, claim, now); err != nil {
					t.Fatal(err)
				}
				return e
			}
			m := send(firstID, "agent_a1", "agent_b1", "schedule.query", "friend-gate-query-01", "", "", queryBody())
			leased, err := st.leaseMatching("agent_b1", now, false, false)
			if err != nil || leased == nil {
				t.Fatal("lease", err)
			}
			c := command{ID: m.ID, Token: leased.(map[string]any)["lease_token"].(string)}
			for _, action := range []string{"persist", "ack"} {
				if _, err = st.operateAs(action, "agent_b1", "agent", c, now, false); err != nil {
					t.Fatal(err)
				}
			}
			claimed, err := st.operateAs("claim", "agent_b1", "agent", command{ID: m.ID}, now, false)
			if err != nil {
				t.Fatal(err)
			}
			claim := claimed.(map[string]string)["claim"]
			send(gateID, "agent_b1", "agent_b1", "relay.approval.request", "friend-gate-request-01", m.ID, claim, map[string]any{"reason": "judgment_required", "request_digest": st.Messages[m.ID].Receipt.Digest})
			g := st.Gates[gateID]
			if _, valid := verifiedGate(st, g, now); !valid {
				t.Fatal("valid friend gate not verified")
			}
			g.State = "approved"
			switch boundary {
			case "remove", "block":
				if err := st.decideFriend(ids["a"], f.ID, boundary, f.Generation, now); err != nil {
					t.Fatal(err)
				}
			case "key":
				st.Agents["agent_a1"].Keys["key1"].Revoked = true
			case "member":
				st.Members[ids["a"]].Active = false
			}
			if _, valid := verifiedGate(st, g, now); valid {
				t.Fatal("revoked gate remains valid")
			}
			for _, action := range []string{"authorize", "gate-consume"} {
				id := m.ID
				if action == "gate-consume" {
					id = gateID
				}
				if _, err := st.operateAs(action, "agent_b1", "agent", command{ID: id, Claim: claim}, now, false); err == nil {
					t.Fatal("revoked execution authority", action)
				}
			}
			if boundary == "remove" {
				acceptHumans(t, st, ids, now)
				if st.current(st.Messages[m.ID]) {
					t.Fatal("reaccept revived old gate parent")
				}
				result, err := Parse(wire(t, keys["agent_b1"], secondID, "agent_b1", "agent_a1", "relay.result", "friend-result-old-01", m.ID, map[string]any{"status": "denied"}, now.Add(time.Minute)))
				if err != nil {
					t.Fatal(err)
				}
				if _, err = st.ingest("agent_b1", result, claim, now); err == nil {
					t.Fatal("old claim emitted result after reaccept")
				}
			}
		})
	}
}

func TestHumanFriendPrivacyAndCapacity(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	for _, code := range []string{"unknown-code", st.Members[ids["a"]].FriendCode} {
		st.requestFriend(ids["a"], code, now)
		f := st.Friends[pairID(ids["a"], "undisclosed:"+hashToken(code))]
		if f == nil || f.State != "pending" || f.NameB != "" {
			t.Fatal("code outcome disclosed target existence")
		}
	}
	st.requestFriend(ids["a"], st.Members[ids["b"]].FriendCode, now)
	f := st.Friends[pairID(ids["a"], ids["b"])]
	for i := 0; i < 20; i++ {
		st.Friends[fmt.Sprint("capacity", i)] = &Friend{A: ids["a"], B: fmt.Sprint("other", i), State: "active"}
	}
	if err := st.decideFriend(ids["b"], f.ID, "accept", f.Generation, now); err == nil || err.Error() != "capacity" || f.State != "pending" {
		t.Fatal("member friend limit changed pending", err)
	}
	if err := st.decideFriend(ids["a"], f.ID, "cancel", f.Generation, now); err != nil {
		t.Fatal("full capacity blocked cleanup", err)
	}
	st, ids, _, now = humanFixture(t)
	for i := 0; i < 11; i++ {
		st.requestFriend(ids["a"], fmt.Sprint("unknown-code-", i), now)
	}
	if len(st.Friends) != 10 {
		t.Fatal("pending sender limit", len(st.Friends))
	}
	st.sweep(now.Add(24 * time.Hour))
	st.requestFriend(ids["a"], "fresh-code", now.Add(24*time.Hour))
	if len(st.Friends) != 11 {
		t.Fatal("expired pending limit not recovered")
	}
}

func TestHumanRemoteOwnerAndAgentAuthorship(t *testing.T) {
	st, ids, keys, now := humanFixture(t)
	f := acceptHumans(t, st, ids, now)
	i := &remoteInstall{Agent: "agent_b1"}
	owner, err := st.sendOwnerText(ids["a"], "agent_a1", i.Agent, "synthetic web text", "remote-owner-test-01", now)
	if err != nil {
		t.Fatal(err)
	}
	receive := func() map[string]any {
		t.Helper()
		v, err := st.remoteTool(i, "knowslink_receive", nil, now)
		if err != nil {
			t.Fatal(err)
		}
		return v.(map[string]any)
	}
	v := receive()
	if v["untrusted"] != true || v["authorship"].(map[string]string)[owner.Receipt.ID] != "owner" || owner.Receipt.State != "delivered" || len(i.Inbox) != 1 {
		t.Fatal("remote owner author/durable ACK contract")
	}
	signed := signedText(firstID, "agent_a1", i.Agent, "remote-signed-test-01", "", "synthetic agent text", now.Add(time.Minute), keys["agent_a1"])
	if _, err := st.ingestPublicText("agent_a1", signed, now); err != nil {
		t.Fatal(err)
	}
	v = receive()
	if v["authorship"].(map[string]string)[firstID] != "agent" || len(v["messages"].([]json.RawMessage)) != 2 {
		t.Fatal("signed author or retained response-loss Inbox")
	}
	if err := st.decideFriend(ids["a"], f.ID, "block", f.Generation, now); err != nil {
		t.Fatal(err)
	}
	v = receive()
	if len(v["messages"].([]json.RawMessage)) != 0 || len(v["authorship"].(map[string]string)) != 0 {
		t.Fatal("remote revoked Inbox disclosed")
	}
	if len(st.historical(st.Members[ids["a"]].Owner, "agent_a1", i.Agent, now)) != 2 {
		t.Fatal("remote ACK/block erased separate historical copy")
	}
}
