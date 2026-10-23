//go:build integration

// U8-A and O-NM-03 failure boundaries over real forms and Postgres: resume rows follow stored state only, repeats never
// mint agents or replace keys, and the sent-request code tail is the same shape whether or not the code exists.
package relay

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestConnectionResumeAndCodeTail(t *testing.T) {
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
	people, members, owners := map[string]*browser{}, map[string]string{}, map[string]string{}
	for i, p := range []string{"a", "b", "c"} {
		b := &browser{t, h, fmt.Sprintf("192.0.2.%d", 61+i), map[string]string{}}
		members[p] = login(t, b, mail, p+"@connect-guide.invalid")
		expect(t, b.do("POST", "/home/profile", url.Values{"name": {"합성 회원 " + p}}), 303)
		expect(t, b.do("POST", "/home/friend-code", nil), 303)
		people[p] = b
	}
	a := people["a"]
	codes := map[string]string{}
	s.mutateState(t, func(st *State) {
		for p, id := range members {
			owners[p], codes[p] = st.Members[id].Owner, st.Members[id].FriendCode
		}
	})
	reset := func() { s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} }) }
	agentsOf := func(p string) (ids []string) {
		s.mutateState(t, func(st *State) {
			for id, ag := range st.Agents {
				if ag.Owner == owners[p] {
					ids = append(ids, id)
				}
			}
		})
		return ids
	}
	resumeOf := func(svc *Service, agent string) (r agentResume) {
		svc.mutateState(t, func(st *State) { r = st.agentResume(agent, time.Now()) })
		return r
	}
	row := func(body, agent string) string {
		start := strings.Index(body, `<li id="agent-`+agent+`">`)
		if start < 0 {
			t.Fatal("agent row missing", agent)
		}
		end := strings.Index(body[start:], "</li>")
		return body[start : start+end]
	}

	t.Run("choices_follow_host_table_and_reading_creates_nothing", func(t *testing.T) {
		body := a.do("GET", "/home/agents", nil).Body.String()
		last := -1
		for _, value := range []string{"web", "grok", "dots", "claude", "codex", "other"} {
			i := strings.Index(body, `name="use" value="`+value+`"`)
			if i <= last {
				t.Fatal("choice order", value)
			}
			last = i
		}
		for _, banned := range []string{"OpenClaw", "Meta Muse", `id="use-guide"`, "온라인</span>"} {
			if strings.Contains(body, banned) {
				t.Fatal("agents page shows", banned)
			}
		}
		guide := a.do("GET", "/home/agents?use=other&start=first", nil).Body.String()
		guide = guide[strings.Index(guide, `id="use-guide"`):]
		guide = guide[:strings.Index(guide, "</section>")]
		for _, banned := range []string{"<code", "<form", "<button", "github.com", "OpenClaw", "Meta Muse"} {
			if strings.Contains(guide, banned) {
				t.Fatal("unsupported host guide offers", banned)
			}
		}
		if !strings.Contains(guide, `href="/home/agents?use=web#use-guide"`) {
			t.Fatal("unsupported host guide must point to web use")
		}
		if strings.Contains(a.do("GET", "/home/agents?use=%3Cscript%3E", nil).Body.String(), `id="use-guide"`) {
			t.Fatal("unknown choice rendered a guide")
		}
		if claude := a.do("GET", "/home/agents?use=claude&start=continue", nil).Body.String(); !strings.Contains(claude, "<code>/knowslink-session on</code>") {
			t.Fatal("session guide lacks the documented command")
		}
		if len(agentsOf("a")) != 0 {
			t.Fatal("reading or refreshing the guide created an agent")
		}
	})

	attempt := strings.Repeat("r", 32)
	var x string
	t.Run("repeated_create_and_rate_limit_recovery_make_one_agent", func(t *testing.T) {
		reset()
		s.mutateState(t, func(st *State) {
			for i := 0; i < 40; i++ {
				st.Rates["http:member:"+members["a"]] = append(st.Rates["http:member:"+members["a"]], time.Now())
			}
		})
		// The shared admission budget refuses before the handler; only the retry instant and "nothing made" matter here.
		expect(t, a.do("POST", "/home/agents", url.Values{"name": {"합성 재개"}, "attempt": {attempt}}), 429, "이후 다시 시도")
		if len(agentsOf("a")) != 0 {
			t.Fatal("rate-limited create made an agent")
		}
		reset()
		first := a.do("POST", "/home/agents", url.Values{"name": {"합성 재개"}, "attempt": {attempt}})
		second := a.do("POST", "/home/agents", url.Values{"name": {"합성 재개"}, "attempt": {attempt}})
		expect(t, first, 303)
		expect(t, second, 303)
		ids := agentsOf("a")
		if len(ids) != 1 || first.Header().Get("Location") != "/home/agents#agent-"+ids[0] || second.Header().Get("Location") != first.Header().Get("Location") {
			t.Fatal("repeat submit", ids, first.Header().Get("Location"), second.Header().Get("Location"))
		}
		x = ids[0]
		if r := resumeOf(s, x); r.Link != keyManage {
			t.Fatal("unused agent next step", r)
		}
	})

	t.Run("local_stages_cancel_expiry_and_lost_response_keep_keys", func(t *testing.T) {
		reset()
		token, cid := browserGrant(t, a, x, "register")
		if !strings.Contains(row(a.do("GET", "/home/agents", nil).Body.String(), x), `href="/home/connections/`+cid+`"`) {
			t.Fatal("waiting attempt not resumable")
		}
		waiting := resumeOf(s, x)
		_, proof := clientPrepare(t, h, token, "key1")
		prepared := resumeOf(s, x)
		expect(t, a.do("POST", "/home/confirm", url.Values{"connection": {cid}}), 303, "/home/connections/"+cid)
		approved := resumeOf(s, x)
		page := a.do("GET", "/home/connections/"+cid, nil).Body.String()
		if strings.Count(page, `aria-current="step"`) != 1 || !strings.Contains(page, `href="/home/agents#agent-`+x+`"`) {
			t.Fatal("approved connection page step/return")
		}
		if waiting.State == prepared.State || prepared.State == approved.State || approved.Link != "/home/connections/"+cid {
			t.Fatal("stages not distinguished", waiting, prepared, approved)
		}
		// The client's completion response is lost: the key exists, and retrying the used token neither adds nor replaces one.
		publicCall(t, h, "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 200)
		publicCall(t, h, "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 401)
		registered := resumeOf(s, x)
		if registered.Link != keyManage || registered.State == approved.State {
			t.Fatal("registered key resume", registered)
		}
		// Cancelled and expired follow-up attempts leave the registered key untouched, even in rotate mode.
		token2, cid2 := browserGrant(t, a, x, "register")
		expect(t, a.do("POST", "/home/cancel", url.Values{"connection": {cid2}}), 303, "/home/agents#agent-"+x)
		publicCall(t, h, "/v1/connect/info", "", map[string]string{"token": token2, "client": supportedClient}, 401)
		token3, _ := browserGrant(t, a, x, "rotate")
		s.mutateState(t, func(st *State) { st.Connections[hashToken(token3)].Exp = time.Now().Add(-time.Second) })
		publicCall(t, h, "/v1/connect/info", "", map[string]string{"token": token3, "client": supportedClient}, 401)
		s.mutateState(t, func(st *State) {
			if k := st.Agents[x].Keys["key1"]; k == nil || k.Revoked || activeKeys(st.Agents[x]) != 1 {
				t.Fatal("registered key replaced or revoked")
			}
		})
		if after := resumeOf(s, x); after != registered {
			t.Fatal("ended attempts changed the registered resume", after)
		}
		// A keyless agent whose attempt expired differs from an unused one; a restart reads the same stored state.
		reset()
		expect(t, a.do("POST", "/home/agents", url.Values{"name": {"합성 만료"}, "attempt": {strings.Repeat("y", 32)}}), 303)
		expect(t, a.do("POST", "/home/agents", url.Values{"name": {"합성 미사용"}, "attempt": {strings.Repeat("z", 32)}}), 303)
		var y, z string
		s.mutateState(t, func(st *State) {
			y, z = st.Idempotency[agentCreateKey(owners["a"], strings.Repeat("y", 32))], st.Idempotency[agentCreateKey(owners["a"], strings.Repeat("z", 32))]
		})
		token4, _ := browserGrant(t, a, y, "register")
		s.mutateState(t, func(st *State) { st.Connections[hashToken(token4)].Exp = time.Now().Add(-time.Second) })
		if resumeOf(s, y).State == resumeOf(s, z).State {
			t.Fatal("expired attempt looks unused")
		}
		restarted := &Service{Pool: pool}
		for _, id := range []string{x, y, z} {
			if resumeOf(restarted, id) != resumeOf(s, id) {
				t.Fatal("restart changed resume", id)
			}
		}
		if len(agentsOf("a")) != 3 {
			t.Fatal("attempts created agents", agentsOf("a"))
		}
	})

	t.Run("device_and_remote_uncertain_results_stay_check_needed", func(t *testing.T) {
		reset()
		device := "agent_device_lost"
		s.mutateState(t, func(st *State) {
			st.Agents[device] = &Agent{Owner: owners["a"], Keys: map[string]*Key{}}
			st.Connections["device-lost"] = &Connection{Device: true, Owner: owners["a"], Agent: device, Client: supportedClient, Mode: "register", State: "approved", Exp: time.Now().Add(-time.Second)}
		})
		if r := resumeOf(s, device); !strings.Contains(r.State, "확인 필요") {
			t.Fatal("device approval without saved key", r)
		}
		// D-CT-01: the connection page asks for the same check as the card instead of declaring no key.
		page := people["a"].do("GET", "/home/connections/device-lost", nil).Body.String()
		if !strings.Contains(page, deviceUnsure.State) || !strings.Contains(page, deviceUnsure.Next) || strings.Contains(page, "새 키가 연결되지 않았습니다") || strings.Contains(page, "새로 시작하세요") {
			t.Fatal("device-lost connection page contradicts the card", page)
		}
		// V-CN-02: an unknown connection returns to the agent list.
		expect(t, people["a"].do("GET", "/home/connections/no-such-connection", nil), 403, `href="/home/agents"`, "내 에이전트로 돌아가기")
		// C consents in the browser: deny makes nothing and cannot be replayed; approve leaves a code until the app returns.
		s.mutateState(t, func(st *State) {
			st.Members[members["c"]].Issuer = googleIssuer
			rs := st.remote()
			rs.Clients["guide-client"] = &remoteClient{Name: "합성 앱", Redirects: []string{"https://app.invalid/cb"}, Exp: time.Now().Add(time.Hour)}
			for _, id := range []string{"guide-deny", "guide-approve"} {
				rs.Attempts[id] = &remoteAttempt{Client: "guide-client", Redirect: "https://app.invalid/cb", State: "synthetic", Owner: owners["c"], Exp: time.Now().Add(5 * time.Minute)}
			}
		})
		c := people["c"]
		deny := url.Values{"csrf": {csrf(c.cookies[sessionCookie], "guide-deny")}, "decision": {"deny"}}
		expect(t, c.do("POST", "/remote/consent/guide-deny", deny), 200, "error=access_denied")
		if w := c.do("POST", "/remote/consent/guide-deny", deny); w.Code == 200 || !strings.Contains(w.Body.String(), `href="/home/agents"`) {
			t.Fatal("replayed deny", w.Code)
		}
		if len(agentsOf("c")) != 0 {
			t.Fatal("deny created an agent")
		}
		approve := url.Values{"csrf": {csrf(c.cookies[sessionCookie], "guide-approve")}, "decision": {"approve"}, "name": {"합성 원격"}}
		expect(t, c.do("POST", "/remote/consent/guide-approve", approve), 200, "https://app.invalid/cb?code=")
		remote := agentsOf("c")
		if len(remote) != 1 {
			t.Fatal("remote approve agents", remote)
		}
		returning := resumeOf(s, remote[0])
		s.mutateState(t, func(st *State) {
			for id, code := range st.Remote.Codes {
				if st.Remote.Installs[code.Install].Agent == remote[0] {
					delete(st.Remote.Codes, id)
				}
			}
		})
		lost := resumeOf(s, remote[0])
		s.mutateState(t, func(st *State) {
			for id, in := range st.Remote.Installs {
				if in.Agent == remote[0] {
					st.Remote.Tokens["guide-token"] = &remoteToken{id, "guide-client", "synthetic", "refresh", time.Now().Add(time.Hour)}
				}
			}
		})
		exchanged := resumeOf(s, remote[0])
		if returning.State == lost.State || lost.State == exchanged.State || !strings.Contains(lost.State, "확인 필요") {
			t.Fatal("remote stages", returning, lost, exchanged)
		}
		if body := c.do("GET", "/home/agents", nil).Body.String(); strings.Contains(body, "온라인</span>") {
			t.Fatal("online claim")
		}
	})

	t.Run("sent_request_tail_is_existence_blind", func(t *testing.T) {
		reset()
		request := func(code string) {
			expect(t, a.do("POST", "/home/friends/request", url.Values{"code": {"  " + code + " "}, "consent": {"yes"}}), 303, "/home/friends?n=requested")
		}
		tail := func(code string) string { return "입력한 코드 …" + code[len(code)-4:] }
		unknown := "synthetic-unknown-zz99"
		request(codes["c"])
		request(unknown)
		// C blocks A from C's incoming row; a later request with C's code lands on the same placeholder shape.
		var friend string
		var generation int64
		s.mutateState(t, func(st *State) {
			f := st.Friends[pairID(members["a"], members["c"])]
			friend, generation = f.ID, f.Generation
		})
		expect(t, people["c"].do("POST", "/home/friends/decision", url.Values{"friend": {friend}, "generation": {fmt.Sprint(generation)}, "decision": {"block"}}), 303)
		reset()
		request(codes["c"])
		request(codes["a"])
		request("ab")
		request("abc!def!")
		request(codes["b"])
		body := a.do("GET", "/home/friends", nil).Body.String()
		for _, want := range []string{tail(unknown), tail(codes["c"]), tail(codes["a"]), tail(codes["b"]), "입력한 코드 형식 다름", "…" + codes["b"][len(codes["b"])-4:] + " 친구 요청 취소"} {
			if !strings.Contains(body, want) {
				t.Fatal("sent row lacks", want)
			}
		}
		// The block on A's first C request stays invisible: one unanswered row, and the repeat is a no-op like any pending repeat.
		if strings.Count(body, tail(codes["c"])) != 2 || strings.Count(body, tail(unknown)) != 2 || strings.Contains(body, "해제") {
			t.Fatal("blocked request differs from an unanswered one")
		}
		for _, secret := range []string{codes["b"], codes["c"], unknown} {
			if strings.Contains(body, secret) {
				t.Fatal("full input kept on the page")
			}
		}
		// Existing, unknown and blocked targets differ only in tail, ID, generation and deadline.
		normal := regexp.MustCompile(`…[A-Za-z0-9_-]{4}|value="[^"]*"|\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} KST`)
		shapes := map[string]bool{}
		for _, li := range regexp.MustCompile(`(?s)<li><h3>보낸 친구 요청</h3>.*?</li>`).FindAllString(body, -1) {
			if !strings.Contains(li, "형식 다름") {
				shapes[normal.ReplaceAllString(li, "_")] = true
			}
		}
		if len(shapes) != 1 {
			t.Fatal("sent rows differ in shape", shapes)
		}
		// B's code change does not recompute A's row; the counterpart never sees any part of an entered code.
		expect(t, people["b"].do("POST", "/home/friend-code", nil), 303)
		if !strings.Contains(a.do("GET", "/home/friends", nil).Body.String(), tail(codes["b"])) {
			t.Fatal("tail recomputed from the current code")
		}
		for _, p := range []string{"b", "c"} {
			if strings.Contains(people[p].do("GET", "/home/friends", nil).Body.String(), "입력한 코드") {
				t.Fatal("counterpart sees the code tail", p)
			}
		}
		// Ended rows keep the requester's tail; a record from before the tail existed says so instead of guessing.
		s.mutateState(t, func(st *State) {
			f := st.Friends[pairID(members["a"], "undisclosed:"+hashToken(unknown))]
			friend, generation = f.ID, f.Generation
			st.Friends[pairID(members["a"], "undisclosed:"+hashToken("ab"))].CodeTail = ""
		})
		expect(t, a.do("POST", "/home/friends/decision", url.Values{"friend": {friend}, "generation": {fmt.Sprint(generation)}, "decision": {"cancel"}}), 303)
		body = a.do("GET", "/home/friends", nil).Body.String()
		if !strings.Contains(body, tail(unknown)) || !strings.Contains(body, "이전 요청") {
			t.Fatal("ended or legacy row label")
		}
	})

	t.Run("existing_agent_key_and_friend_survive", func(t *testing.T) {
		reset()
		var friend string
		var generation int64
		s.mutateState(t, func(st *State) {
			f := st.Friends[pairID(members["a"], members["b"])]
			friend, generation = f.ID, f.Generation
		})
		expect(t, people["b"].do("POST", "/home/friends/decision", url.Values{"friend": {friend}, "generation": {fmt.Sprint(generation)}, "decision": {"accept"}, "consent": {"yes"}}), 303)
		reset()
		token, _ := browserGrant(t, a, x, "register")
		expect(t, a.do("POST", "/home/agents", url.Values{"name": {"합성 재개"}, "attempt": {attempt}}), 303, "/home/agents#agent-"+x)
		s.mutateState(t, func(st *State) {
			st.Connections[hashToken(token)].Exp = time.Now().Add(-time.Second)
			if f := st.Friends[pairID(members["a"], members["b"])]; f.State != "active" || f.Generation != generation {
				t.Fatal("friendship changed")
			}
			if activeKeys(st.Agents[x]) != 1 || st.Agents[x].Revoked {
				t.Fatal("existing key changed")
			}
		})
	})
}
