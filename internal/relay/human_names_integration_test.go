//go:build integration

// Real forms and Postgres pin U2/U4: name-first labels, friend-code sharing and request boundaries.
package relay

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestHumanNamesAndFriendCodePostgres(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, e := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	s, mail := identityService(t, pool)
	h := s.Handler()
	browsers, members := map[string]*browser{}, map[string]string{}
	for i, person := range []string{"a", "b", "c"} {
		b := &browser{t, h, fmt.Sprintf("192.0.2.%d", i+41), map[string]string{}}
		browsers[person] = b
		email := person + "@names.invalid"
		expect(t, b.do("POST", "/auth/start", url.Values{"email": {email}}), 303)
		expect(t, b.do("POST", "/auth/verify", url.Values{"code": {mail.code(email)}}), 303)
		expect(t, b.do("POST", "/home/profile", url.Values{"name": {"합성 이름 " + person}}), 303)
		s.mutateState(t, func(st *State) { members[person], _, _ = st.session(b.cookies[sessionCookie], time.Now()) })
	}
	a, b, c := browsers["a"], browsers["b"], browsers["c"]
	codes := map[string]string{}
	friends := func() int {
		n := 0
		s.mutateState(t, func(st *State) {
			n = len(st.Friends)
			for person, id := range members {
				codes[person] = st.Members[id].FriendCode
			}
		})
		return n
	}
	friends()

	// The own code is a labelled read-only field on both screens; GETs with request-like queries create nothing.
	expect(t, a.do("GET", "/home/account", nil), 200, `<label for="my-code">내 친구 코드</label>`, `readonly`, codes["a"]+"</textarea>")
	expect(t, a.do("GET", "/home/friends?code="+url.QueryEscape(codes["b"])+"&consent=yes", nil), 200, codes["a"]+"</textarea>", "보내는 사람", "합성 이름 a", "연결 범위")
	if friends() != 0 {
		t.Fatal("GET created a relationship")
	}

	// Local format errors keep the typed code, never search for a target, and the agentless member can still request.
	expect(t, a.do("POST", "/home/friends/request", url.Values{"code": {" \t "}, "consent": {"yes"}}), 422, "친구 코드를 입력하세요", `aria-invalid="true" autofocus`)
	expect(t, a.do("POST", "/home/friends/request", url.Values{"code": {codes["b"]}}), 422, "연결 범위에 동의해야", `value="`+codes["b"]+`"`)
	if friends() != 0 {
		t.Fatal("invalid request created a relationship")
	}
	expect(t, a.do("POST", "/home/friends/request", url.Values{"code": {"  " + codes["b"] + "\n"}, "consent": {"yes"}}), 303, "n=requested")
	var f *Friend
	s.mutateState(t, func(st *State) { f = st.Friends[pairID(members["a"], members["b"])] })
	if f == nil || f.State != "pending" {
		t.Fatal("trimmed code did not reach the owner")
	}
	expect(t, b.do("POST", "/home/friends/decision", url.Values{"friend": {f.ID}, "generation": {fmt.Sprint(f.Generation)}, "decision": {"accept"}, "consent": {"yes"}}), 303)

	// Names: 40 emoji pass, a control character is refused without a second agent; a rename keeps ID and relationship.
	expect(t, b.do("POST", "/home/agents", url.Values{"name": {strings.Repeat("😀", 40)}}), 303)
	expect(t, b.do("POST", "/home/agents", url.Values{"name": {"bad\u0007name"}}), 422, "1–40자")
	var agent string
	agentCount := 0
	s.mutateState(t, func(st *State) {
		for id, ag := range st.Agents {
			if ag.Owner == st.Members[members["b"]].Owner {
				agent, agentCount = id, agentCount+1
			}
		}
	})
	if agentCount != 1 {
		t.Fatal("name error created an agent")
	}
	view := a.do("GET", "/home/friends/"+f.ID, nil)
	expect(t, view, 200, strings.Repeat("😀", 40), "첫 에이전트를 추가하세요")
	if strings.Contains(view.Body.String(), "/home/new?peer=") {
		t.Fatal("agentless member offered an unusable conversation")
	}
	renamed := strings.Repeat("가", 40)
	expect(t, b.do("POST", "/home/agent-name", url.Values{"agent": {agent}, "name": {renamed}}), 303)
	expect(t, a.do("GET", "/home/friends/"+f.ID, nil), 200, renamed, shortID(agent))
	s.mutateState(t, func(st *State) {
		g := st.Friends[pairID(members["a"], members["b"])]
		if st.Agents[agent] == nil || st.Agents[agent].Name != renamed || g.State != "active" || g.Generation != f.Generation {
			t.Fatal("rename changed identity or relationship")
		}
	})

	// A replaced code finds nobody new; the existing friendship stays.
	old := codes["b"]
	expect(t, b.do("POST", "/home/friend-code", nil), 303)
	expect(t, c.do("POST", "/home/friends/request", url.Values{"code": {old}, "consent": {"yes"}}), 303, "n=requested")
	s.mutateState(t, func(st *State) {
		if st.Friends[pairID(members["c"], members["b"])] != nil || st.Friends[pairID(members["a"], members["b"])].State != "active" || st.Members[members["b"]].FriendCode == old {
			t.Fatal("old code still discovers its owner or rotation ended a friendship")
		}
	})
}
