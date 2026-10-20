// Connection guide pages: host choices, resume rows and the requester's code tail render safely in every variant.
package relay

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodeTailKeepsOnlyPlainTails(t *testing.T) {
	for input, want := range map[string]string{"fixture-code-ab12": "ab12", "x-_Z9": "-_Z9", "abc": "?", "abc!def!": "?", "코드가나다라": "?", "abcdé123": "?", "": "?"} {
		if got := codeTail(input); got != want {
			t.Fatalf("codeTail(%q) = %q, want %q", input, got, want)
		}
	}
	labels := map[string]bool{}
	for _, tail := range []string{"", "?", "ab12"} {
		labels[humanFriend{CodeTail: tail}.CodeLabel()] = true
	}
	if len(labels) != 3 {
		t.Fatal("legacy, malformed and plain tails must read differently")
	}
}

// KNOWSLINK_UI_FIXTURE_DIR also writes each variant for the 1440/390/320 render probe.
func TestConnectGuidePagesRender(t *testing.T) {
	dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR")
	write := func(name, body string) {
		if dir != "" {
			if err := os.WriteFile(filepath.Join(dir, name+".html"), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	long := strings.Repeat("긴한글이름", 8)
	resume := map[string]agentResume{
		"agent_a": {"키 등록 전 · 지문 확인 필요", "내 컴퓨터에 표시된 지문과 같은지 확인하세요. 다르면 취소하세요.", "/home/connections/c", "지문 확인"},
		"agent_b": {"확인 필요 · 앱의 토큰 교환을 확인하지 못함", "시작한 앱에서 연결 상태를 먼저 확인하세요.", "", ""},
	}
	for _, choice := range append(useChoices, useChoice{ID: ""}) {
		for _, start := range []string{"first", "continue"} {
			v := humanPageFixture()
			v["Title"], v["Section"], v["Use"], v["Start"], v["Resume"] = "내 에이전트", "agents", choice.ID, start, resume
			v["Own"] = []map[string]any{{"ID": "agent_a", "Name": long, "Client": "Node 로컬", "Status": "연결 필요", "Short": "a001"}, {"ID": "agent_b", "Name": "Synthetic remote helper with a long English name", "Client": "원격 호스트", "Status": "등록됨", "Short": "b002"}}
			body := renderHumanPage(t, "human-agents", v)
			for _, banned := range []string{"OpenClaw", "Meta Muse", "온라인</span>", "<no value>"} {
				if strings.Contains(body, banned) {
					t.Fatal(choice.ID, start, "renders", banned)
				}
			}
			if strings.Contains(body, `id="use-guide"`) != (choice.ID != "") {
				t.Fatal("guide shown without a valid choice", choice.ID)
			}
			if choice.ID == "other" {
				guide := body[strings.Index(body, `id="use-guide"`):]
				if guide = guide[:strings.Index(guide, "</section>")]; strings.Contains(guide, "<code") || strings.Contains(guide, "github.com") || strings.Contains(guide, "<form") {
					t.Fatal("unsupported host guide offers an install path")
				}
			}
			if start == "first" || choice.ID == "" {
				write("agents-"+map[string]string{"": "none"}[choice.ID]+choice.ID, body)
			}
		}
	}
	v := humanPageFixture()
	v["Title"], v["Section"] = "친구", "friends"
	v["Requests"] = []humanFriend{
		{ID: "r1", State: "pending", Generation: 1, Exp: "2026-10-21 13:19:59 KST", CodeTail: "ab12"},
		{ID: "r2", State: "pending", Generation: 1, Exp: "2026-10-21 13:20:01 KST", CodeTail: "?"},
		{ID: "r3", State: "pending", Generation: 1, Exp: "2026-10-21 13:20:03 KST"},
		{ID: "r4", State: "cancelled", Generation: 2, Exp: "2026-10-21 13:20:05 KST", CodeTail: "zz99"},
		{ID: "r5", Name: "합성 친구 요청자", State: "pending", Generation: 1, Incoming: true, Exp: "2026-10-21 13:20:07 KST"},
	}
	body := renderHumanPage(t, "human-friends", v)
	for _, want := range []string{"입력한 코드 …ab12", "입력한 코드 형식 다름", "이전 요청", "입력한 코드 …zz99", `aria-label="입력한 코드 …ab12 친구 요청 취소, 기한 2026-10-21 13:19:59 KST"`} {
		if !strings.Contains(body, want) {
			t.Fatal("friends page lacks", want)
		}
	}
	if strings.Count(body, "입력한 코드") != 5 { // r1 row+button, r2 row+button, r4 row; never the incoming r5
		t.Fatal("code tail outside the requester rows", strings.Count(body, "입력한 코드"))
	}
	write("friends-code-tail", body)
	conn := renderHumanPage(t, "connection", map[string]any{"Title": "에이전트 연결", "ConnectionID": "c", "Agent": "agent_a", "Name": long, "Short": "a001", "Client": supportedClient, "Mode": "register", "State": "consumed", "Exp": "2026-10-20 13:19:59 KST"})
	write("connection-consumed", conn)
	var b bytes.Buffer
	if err := remoteReturn.Execute(&b, "https://app.invalid/cb?code=synthetic&iss=https%3A%2F%2Flink.invalid&state=s"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), `href="https://app.invalid/cb?code=synthetic`) {
		t.Fatal("remote return lost the verified app link")
	}
	write("remote-return", b.String())
}
