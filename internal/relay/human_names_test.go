// U2/U4 boundaries: names appear only where authority allows, discriminators widen, and GETs never build relationships.
package relay

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVisibleAgentNamesOnlyOwnerAndCurrentFriend(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	owner := st.Members[ids["a"]].Owner
	unknown := func(id string) {
		t.Helper()
		if v := st.visibleAgent(owner, id); v.Name != "확인할 수 없는 에이전트" || v.Short != id || v.Client != "" {
			t.Fatalf("stranger agent named: %+v", v)
		}
	}
	if v := st.visibleAgent(owner, "agent_a1"); v.Name != "도우미" {
		t.Fatal("own agent unnamed")
	}
	unknown("agent_b1")
	unknown("agent_missing")
	f := acceptHumans(t, st, ids, now)
	if v := st.visibleAgent(owner, "agent_b1"); v.Name != "도우미" || v.Short != shortID("agent_b1") {
		t.Fatal("current friend agent unnamed")
	}
	unknown("agent_c1")
	if e := st.decideFriend(ids["b"], f.ID, "block", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	unknown("agent_b1")
}

func TestLegacyAgentLabelsWidenCollidingDiscriminators(t *testing.T) {
	st := newState()
	ids := []string{"agent_collision_5591", "agent_collision_11122"}
	if shortID(ids[0]) != shortID(ids[1]) {
		t.Fatal("fixture must collide")
	}
	for _, id := range ids {
		st.Agents[id] = &Agent{Owner: "owner-a", Name: "같은 이름"}
	}
	for _, a := range ownAgents(st, "owner-a") {
		if a.Short != a.ID || a.Name != "같은 이름" {
			t.Fatalf("colliding label kept a short discriminator: %+v", a)
		}
	}
}

func TestFriendViewsNeverCreateRelationships(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	for id, a := range st.Agents {
		if a.Owner == st.Members[ids["a"]].Owner {
			delete(st.Agents, id)
		}
	}
	code := st.Members[ids["b"]].FriendCode
	for _, path := range []string{"friends", "account", "new", "agents"} {
		r := httptest.NewRequest("GET", "/home/"+path+"?code="+code+"&consent=yes&decision=accept", nil)
		if _, e := st.humanView(ids["a"], &Session{Verified: now}, path, r, now); e != nil {
			t.Fatal(path, e)
		}
	}
	if len(st.Friends) != 0 {
		t.Fatal("GET created a relationship")
	}
	// An agentless member can still request, and acceptance needs no agent on either side.
	f := acceptHumans(t, st, ids, now)
	if f.State != "active" {
		t.Fatal("agentless friendship refused")
	}
	// A request accepted after its 24h deadline fails; a fresh request is a new generation.
	st2, ids2, _, now2 := humanFixture(t)
	st2.requestFriend(ids2["a"], st2.Members[ids2["b"]].FriendCode, now2)
	g := st2.Friends[pairID(ids2["a"], ids2["b"])]
	if e := st2.decideFriend(ids2["b"], g.ID, "accept", g.Generation, g.Exp); e == nil {
		t.Fatal("expired request accepted")
	}
	old := g.Generation
	st2.sweep(g.Exp)
	st2.requestFriend(ids2["a"], st2.Members[ids2["b"]].FriendCode, g.Exp)
	if g.State != "pending" || g.Generation != old+1 {
		t.Fatal("re-request after expiry must start a new generation")
	}
}

func TestFriendCodeRotationKeepsReceivedRequests(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	a, b, c := ids["a"], ids["b"], ids["c"]
	old := st.Members[b].FriendCode
	st.requestFriend(a, old, now)
	st.Members[b].FriendCode = randomToken()
	st.requestFriend(c, old, now)
	if st.Friends[pairID(c, b)] != nil || st.Friends[pairID(a, b)].State != "pending" {
		t.Fatal("rotation must stop discovery but keep the received request")
	}
	st.requestFriend(c, " "+st.Members[b].FriendCode, now)
	if st.Friends[pairID(c, b)] != nil {
		t.Fatal("codes are exact; the handler trims, not the matcher")
	}
}

func TestFriendPagesShareOwnCodeAndKeepRequestInput(t *testing.T) {
	v := humanPageFixture()
	v["Problem"], v["Invalid"], v["RequestCode"] = "친구 코드를 입력하세요.", "friend-code", `"><b>`
	body := renderHumanPage(t, "human-friends", v)
	for _, want := range []string{`<label for="my-code">내 친구 코드</label><textarea id="my-code" class="share" readonly`, `>synthetic-friend-code</textarea>`, "보내는 사람", "합성 회원", `value="&#34;&gt;&lt;b&gt;" aria-describedby="friend-help problem" aria-invalid="true" autofocus`, "연결 범위", "친구 요청 보내기"} {
		if !strings.Contains(body, want) {
			t.Fatal("missing", want)
		}
	}
	if strings.Count(body, " autofocus") != 1 || strings.Contains(body, "<b>") {
		t.Fatal("error focus must be unique and input escaped")
	}
	body = renderHumanPage(t, "human-account", humanPageFixture())
	if !strings.Contains(body, `readonly rows="2"`) || !strings.Contains(body, ">synthetic-friend-code</textarea>") || strings.Contains(body, `name="code"`) {
		t.Fatal("account code must be a read-only, unsubmitted field")
	}
}

func TestGateNamesEndpointsWithinAuthority(t *testing.T) {
	v := map[string]any{"State": "pending", "From": "agent_x", "To": "agent_a", "Parent": "request-a", "Exp": "e", "Policy": "p", "Body": "{}", "Back": "/home",
		"FromAgent": humanAgent{ID: "agent_x", Name: "확인할 수 없는 에이전트", Short: "agent_x"}, "ToAgent": humanAgent{ID: "agent_a", Name: "내 도우미", Client: "Node 22", Short: "a001"}}
	body := renderHumanPage(t, "human-gate", v)
	if !strings.Contains(body, "확인할 수 없는 에이전트 / agent_x<br><code>agent_x</code>") || !strings.Contains(body, "내 도우미 / Node 22 / a001<br><code>agent_a</code>") {
		t.Fatal("gate must lead with permitted names and keep full IDs")
	}
}
