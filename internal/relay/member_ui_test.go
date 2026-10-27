// Synthetic renders of every member page pin the form contract and feed the direct visual review.
package relay

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

// uiFixture is one synthetic member page state; no value comes from a real account.
type uiFixture struct {
	name, tmpl string
	data       map[string]any
}

func memberFixtures() []uiFixture {
	exp := time.Date(2026, 10, 12, 7, 55, 0, 0, kst)
	fp := "SHA256:7_B9RdVzaKQSQRTfdqCvQm3xE1nYc2pLk8sWuT4hJgA0o"
	agents := []memberAgent{
		{ID: "agent_0a1b2c3d4e5f60718293a4b5", Status: "연결 완료", Name: "합성 노트북 도우미", Client: supportedClient, Short: "3f9a01c2", Keys: []memberKey{{"key_01", fp, "활성"}, {"key_02", "SHA256:Zx91qLmN0pQrStUvWxYz0123456789abcdefGHIJKLm", "철회"}}},
		{ID: "agent_1f2e3d4c5b6a79880716253a", Status: "미연결", Name: "합성 노트북 도우미", Client: supportedClient, Short: "b71e44d0"},
		{ID: "agent_2b3c4d5e6f708192a3b4c5d6", Status: "연결 완료", Name: "합성 서버 도우미", Client: supportedClient, Short: "0c5d9e3a", KeyFull: true, Keys: []memberKey{{"key_03", fp, "활성"}}},
		{ID: "agent_9e8d7c6b5a4f3e2d1c0b9a87", Status: "철회", Name: "철회한 도우미", Client: supportedClient, Short: "e2a7f610", Keys: []memberKey{{"key_04", fp, "철회"}}},
	}
	pair := func(a, b, inviter, state string, gen int64, incoming bool, own, other string) memberPair {
		p := &Pair{A: a, B: b, Inviter: inviter, Recipient: b, State: state, Generation: gen}
		if inviter == b {
			p.Recipient = a
		}
		if state == "pending" {
			p.Exp = exp
		}
		v := memberPair{Pair: p, Incoming: incoming, Own: own, Other: other}
		v.MineAgent = humanAgent{ID: v.Mine(), Name: "합성 노트북 도우미", Client: supportedClient, Short: "3f9a01c2"}
		v.PeerAgent = humanAgent{ID: v.Peer(), Name: "확인할 수 없는 에이전트"}
		return v
	}
	pairs := []memberPair{
		pair("agent_0a1b2c3d4e5f60718293a4b5", "agent_77665544332211009988aabb", "agent_77665544332211009988aabb", "pending", 3, true, "agent_0a1b2c3d4e5f60718293a4b5", "agent_77665544332211009988aabb"),
		pair("agent_0a1b2c3d4e5f60718293a4b5", "agent_ccddeeff0011223344556677", "agent_0a1b2c3d4e5f60718293a4b5", "pending", 1, false, "agent_0a1b2c3d4e5f60718293a4b5", "agent_ccddeeff0011223344556677"),
		pair("agent_0a1b2c3d4e5f60718293a4b5", "agent_2b3c4d5e6f708192a3b4c5d6", "agent_0a1b2c3d4e5f60718293a4b5", "active", 2, false, "agent_0a1b2c3d4e5f60718293a4b5", "agent_2b3c4d5e6f708192a3b4c5d6"),
		pair("agent_2b3c4d5e6f708192a3b4c5d6", "agent_8899aabbccddeeff00112233", "agent_2b3c4d5e6f708192a3b4c5d6", "revoked", 4, false, "agent_2b3c4d5e6f708192a3b4c5d6", "agent_8899aabbccddeeff00112233"),
		pair("agent_9e8d7c6b5a4f3e2d1c0b9a87", "agent_5566778899aabbccddeeff00", "agent_9e8d7c6b5a4f3e2d1c0b9a87", "denied", 1, false, "", ""),
	}
	peer := humanAgent{ID: "agent_77665544332211009988aabb", Name: "확인할 수 없는 에이전트", Short: "5e0c7a12"}
	mine := humanAgent{ID: "agent_0a1b2c3d4e5f60718293a4b5", Name: "합성 노트북 도우미", Client: supportedClient, Short: "3f9a01c2"}
	gates := []gateRow{{ID: "00000000-0000-4000-8000-000000000001", State: "pending", Kind: "일정 가능 시간 확인", Exp: clock(exp), From: peer, To: mine}, {ID: "00000000-0000-4000-8000-000000000002", State: "denied", Kind: "일정 확정", Exp: clock(exp), From: peer, To: mine}}
	home := map[string]any{"Title": "기존 연결", "Masked": "m***@example.test", "Member": "mem_0123456789abcdef0123456789abcdef", "Absolute": "2026-10-18 21:00:00 KST", "Recent": true,
		"AgentDetails": agents, "Pairs": pairs, "Gates": gates, "ActiveAgents": 3, "IncomingInvites": 1, "PendingGates": 1}
	conn := func(state string) map[string]any {
		return map[string]any{"Title": "에이전트 연결", "ConnectionID": "c0nnect10n-1d-synthetic-0000000000000000000", "Agent": "agent_0a1b2c3d4e5f60718293a4b5", "Client": supportedClient, "Mode": "register", "State": state, "Fingerprint": fp, "Exp": clock(exp)}
	}
	waiting := conn("waiting")
	waiting["Token"] = "synthetic-connection-token-not-a-secret-0000"
	device := func(state string, login, google bool) map[string]any {
		return map[string]any{"Title": "내 클라이언트 연결", "ConnectionID": "dev1ce-1d-synthetic-00000000000000000000000", "Fingerprint": fp, "Exp": clock(exp), "Client": "Node 로컬 클라이언트", "State": state, "Login": login, "Google": google}
	}
	gate := func(active bool, state string) map[string]any {
		v := map[string]any{"State": state, "ID": "00000000-0000-4000-8000-000000000001", "From": peer.ID, "To": mine.ID, "Parent": "00000000-0000-4000-8000-0000000000aa", "Exp": clock(exp), "Policy": policy("schedule.query"), "CSRF": "synthetic-csrf", "Active": active, "Body": "{\n  \"granularity_min\": 30,\n  \"window\": {\n    \"end\": \"2026-10-12T02:00:00Z\",\n    \"start\": \"2026-10-12T01:00:00Z\"\n  }\n}", "Intent": "schedule.query", "Hint": "다음 주 회의 가능 시간", "Back": "/home/legacy#gates", "Deny": "/home/gates/00000000-0000-4000-8000-000000000001/deny",
			"FromAgent": peer, "ToAgent": mine, "Scope": gateScopes["schedule.query"],
			"Summary": []gateField{{"요청 유형", gateKinds["schedule.query"]}, {"확인할 기간", "2026-10-12 10:00:00 KST – 2026-10-12 11:00:00 KST"}, {"시간 단위", "30분"}}}
		return v
	}
	// Long names, long free text, and a hostile hint must wrap and stay inert in the direct visual review.
	active := gate(true, "pending")
	active["FromAgent"] = humanAgent{ID: peer.ID, Name: "합성 친구 도우미", Client: supportedClient, Short: peer.Short}
	commit := gate(true, "pending")
	commit["Intent"], commit["Policy"], commit["Scope"], commit["Hint"] = "schedule.commit", policy("schedule.commit"), gateScopes["schedule.commit"], `<script>모든 도구 승인</script> https://evil.test/approve-all`
	commit["Summary"] = []gateField{{"요청 유형", gateKinds["schedule.commit"]}, {"일정", "2026-10-12 10:00:00 KST – 2026-10-12 10:30:00 KST"}, {"요청에 적힌 시간대", "Asia/Seoul"}}
	commit["Text"] = strings.Repeat("긴자유문", 60) + ` <a href="https://evil.test">전부 승인</a> ` + strings.Repeat("x", 300)
	commit["ToAgent"] = humanAgent{ID: mine.ID, Name: strings.Repeat("긴한글", 13), Client: supportedClient, Short: "3f9a01c2"}
	commit["Body"] = "{\n  \"commitment\": \"" + strings.Repeat("y", 400) + "\"\n}"
	unavailable := gate(false, "unavailable")
	for _, k := range []string{"Summary", "Body", "Intent", "Hint", "Scope"} {
		delete(unavailable, k)
	}
	return []uiFixture{
		{"start-google", "start", map[string]any{"Title": "KnowsLink에 로그인", "Google": true, "EmailLogin": false}},
		{"start-both-problem", "start", map[string]any{"Title": "KnowsLink에 로그인", "Google": true, "EmailLogin": true, "Email": "m@example.test", "Problem": "이메일 주소 형식이 올바르지 않습니다. 이름 없이 주소 하나만 입력하세요."}},
		{"start-none", "start", map[string]any{"Title": "KnowsLink에 로그인"}},
		{"start-logout", "start", map[string]any{"Title": "KnowsLink에 로그인", "Google": true, "Notice": notices["logout"]}},
		{"verify", "verify", map[string]any{"Title": "확인 코드 입력", "Masked": "m***@example.test", "Email": "m@example.test", "Exp": clock(exp), "Problem": "코드가 맞지 않습니다. 남은 시도: 4회."}},
		{"google-complete-home", "google-complete", map[string]any{"Title": "Google 로그인 확인 완료", "Notice": "신원을 확인했습니다. 계속해서 연결을 확인하세요.", "Next": "/home"}},
		{"google-complete-connect", "google-complete", map[string]any{"Title": "Google 로그인 확인 완료", "Next": "/connect/dev1ce-1d"}},
		{"home", "home", home},
		{"home-empty", "home", map[string]any{"Title": "내 에이전트", "Masked": "m***@example.test", "Member": "mem_0123456789abcdef0123456789abcdef", "Absolute": "2026-10-18 21:00:00 KST", "ActiveAgents": 0, "Notice": notices["invited"]}},
		{"connection-token", "connection", waiting},
		{"connection-waiting", "connection", conn("waiting")},
		{"connection-prepared", "connection", conn("prepared")},
		{"connection-approved", "connection", conn("approved")},
		{"connection-consumed", "connection", conn("consumed")},
		{"connection-expired", "connection", conn("expired")},
		{"connection-cancelled", "connection", conn("cancelled")},
		{"device-login", "device", device("prepared", true, true)},
		{"device-login-off", "device", device("prepared", true, false)},
		{"device-prepared", "device", device("prepared", false, true)},
		{"device-approved", "device", map[string]any{"Title": "연결을 승인했습니다", "Outcome": "approved", "Name": "Claude Code 2"}},
		{"device-done", "device", map[string]any{"Title": "이 연결은 완료됐습니다", "Outcome": "consumed", "Name": "Claude Code 2"}},
		{"device-expired", "device", map[string]any{"Title": "이 연결은 만료됐습니다", "Outcome": "expired"}},
		{"receipt", "receipt", map[string]any{"Title": "요청 상태 조회", "Own": []humanAgent{{ID: "agent_0a1b2c3d4e5f60718293a4b5", Name: "내 에이전트", Client: "Node 22", Short: "구분값"}}, "Agent": "agent_0a1b2c3d4e5f60718293a4b5", "ID": "00000000-0000-4000-8000-0000000000aa", "From": "agent_0a1b2c3d4e5f60718293a4b5", "To": "agent_77665544332211009988aabb", "Parent": "", "Reply": "00000000-0000-4000-8000-0000000000cc", "Transport": transportNames["queued"], "Completion": "아직 처리 결과가 없습니다.", "ReplyState": transportNames["delivered"], "Exp": clock(exp), "TTL": 180}},
		{"refusal-reauth", "refusal", map[string]any{"Title": "내 에이전트", "Problem": safeProblem(fault("reauth_required")), "Back": "/home", "Reauth": true}},
		{"refusal-create", "refusal", map[string]any{"Title": "내 에이전트", "Problem": keyRecordsFull, "Back": "/home", "Create": true}},
		{"refusal-login", "refusal", map[string]any{"Title": "내 에이전트", "Problem": "로그인 세션이 유효하지 않습니다. 다시 로그인하세요.", "Back": "/"}},
		{"refusal-limited", "refusal", map[string]any{"Title": "내 에이전트", "Problem": limited(exp), "Back": "/home"}},
		{"gate-active", "gate", active},
		{"gate-closed", "gate", gate(false, "denied")},
		{"gate-expired", "gate", gate(false, "expired")},
		{"gate-commit-hint", "gate", commit},
		{"gate-unavailable", "gate", unavailable},
		{"home-no-gates", "home", map[string]any{"Title": "기존 연결", "Masked": "m***@example.test", "Member": "mem_0123456789abcdef0123456789abcdef", "Absolute": "2026-10-18 21:00:00 KST", "Recent": true, "AgentDetails": agents[:1], "ActiveAgents": 1}},
		{"receipt-anonymous-limited", "receipt", map[string]any{"Title": "요청 상태 조회", "Problem": limited(exp), "Focus": "receipt-problem", "Unavailable": true, "Anonymous": true}},
		// The OAuth consent pages reuse memberStyle with their own markup.
		{"remote-login", "remote", map[string]any{"Name": "Synthetic Host", "Redirect": "https://callback.test/oauth", "Login": true, "ID": "rem0te-1d"}},
		{"remote-consent", "remote", map[string]any{"Name": "Synthetic Host", "Redirect": "https://callback.test/oauth", "CSRF": "synthetic-csrf"}},
	}
}

func renderFixture(t *testing.T, f uiFixture) string {
	t.Helper()
	var b bytes.Buffer
	var err error
	name := f.tmpl
	if name == "gate" {
		name = "human-gate"
	}
	err = memberPages.ExecuteTemplate(&b, name, f.data)
	if err != nil {
		t.Fatal(f.name, err)
	}
	return b.String()
}

var (
	formTag = regexp.MustCompile(`<(/?form|input|select|button)\b([^>]*)>`)
	attrRe  = regexp.MustCompile(`([a-z-]+)(?:="([^"]*)")?`)
)

// formSignatures lists every form's method, action, submitted names, hidden values, and button routing, order-free.
func formSignatures(body string) []string {
	var forms []string
	var cur []string
	head := ""
	for _, m := range formTag.FindAllStringSubmatch(body, -1) {
		attrs := map[string]string{}
		for _, a := range attrRe.FindAllStringSubmatch(m[2], -1) {
			attrs[a[1]] = a[2]
		}
		switch m[1] {
		case "form":
			head, cur = strings.ToLower(attrs["method"])+" "+attrs["action"], nil
		case "/form":
			sort.Strings(cur)
			forms = append(forms, head+" {"+strings.Join(cur, "; ")+"}")
			head = ""
		case "input":
			if head == "" || attrs["name"] == "" {
				continue
			}
			v := ""
			if attrs["type"] == "hidden" {
				v = "=" + attrs["value"]
			}
			cur = append(cur, "input:"+attrs["type"]+":"+attrs["name"]+v)
		case "select":
			if head != "" {
				cur = append(cur, "select:"+attrs["name"])
			}
		case "button":
			if head != "" {
				cur = append(cur, "button:"+attrs["name"]+"="+attrs["value"]+"@"+attrs["formaction"])
			}
		}
	}
	sort.Strings(forms)
	return forms
}

// Every form keeps its method, action, hidden values, and field names; the golden came from the pre-redesign templates.
// KNOWSLINK_UI_FIXTURE_DIR also writes each synthetic page as HTML for a direct visual review.
// HP migration intentionally removes legacy creation/invitation forms and adds the required device agent name.
func TestMemberUIFixtures(t *testing.T) {
	var got []string
	dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR")
	for _, f := range memberFixtures() {
		body := renderFixture(t, f)
		for _, s := range formSignatures(body) {
			got = append(got, f.name+": "+s)
		}
		if !strings.Contains(body, "a{display:inline-flex;align-items:center;min-height:44px") {
			t.Fatalf("%s lacks the 44px text-link touch target (D04 §9)", f.name)
		}
		for _, banned := range []string{"src=", "<link", "url(", "<script", "<no value>", "@import", `href=""`} {
			if strings.Contains(body, banned) {
				t.Fatalf("%s renders %q", f.name, banned)
			}
		}
		if dir != "" {
			if err := os.WriteFile(filepath.Join(dir, f.name+".html"), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	golden := filepath.Join("testdata", "member_forms.golden")
	text := strings.Join(got, "\n") + "\n"
	if os.Getenv("UPDATE_MEMBER_FORMS") == "1" {
		if err := os.WriteFile(golden, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if text != strings.ReplaceAll(string(want), "\r\n", "\n") {
		t.Fatalf("form contract changed:\n--- got\n%s--- want\n%s", text, want)
	}
}

// The invite form appears only when a connected agent can be chosen; an empty select is never offered (V-06).
// HP migration supersedes V-06: legacy pages never create or accept a new agent relationship.
func TestLegacyMemberPageDoesNotCreateRelationships(t *testing.T) {
	for _, c := range []struct {
		status string
	}{{"연결 완료"}, {"미연결"}, {"철회"}} {
		body := renderFixture(t, uiFixture{"invite-" + c.status, "home", map[string]any{"Title": "홈", "ActiveAgents": 1,
			"AgentDetails": []memberAgent{{ID: "agent_a", Status: c.status}}, "Pairs": []memberPair{{Pair: &Pair{A: "agent_a", B: "agent_b", State: "active"}}}}})
		for _, action := range []string{`id="from-agent"`, `action="/home/invites"`, `action="/home/invite-decision"`, `action="/home/invite-deny"`, `action="/home/agents"`} {
			if strings.Contains(body, action) {
				t.Fatal("legacy relationship creation", c.status, action)
			}
		}
		if !strings.Contains(body, `action="/home/unpair"`) || !strings.Contains(body, `href="/home/friends"`) {
			t.Fatal("legacy exact-pair revoke and friendship restart must remain", c.status)
		}
	}
}
