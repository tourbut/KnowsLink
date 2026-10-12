package relay

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func humanPageFixture() map[string]any {
	agent := func(id, short string) map[string]any {
		return map[string]any{"ID": id, "Name": "같은 이름", "Client": "Node 로컬", "Status": "연결 완료", "Short": short}
	}
	return map[string]any{
		"Title": "합성 대화", "Section": "conversations", "Name": "합성 회원", "Masked": "s***@example.test",
		"Own":    []map[string]any{agent("agent_a", "a001"), agent("agent_b", "b002")},
		"Agents": []map[string]any{agent("agent_c", "c003")}, "From": agent("agent_a", "a001"), "To": agent("agent_c", "c003"),
		"Peer":    map[string]any{"ID": "friend-a", "Name": "합성 친구", "State": "active", "Generation": 7},
		"Friends": []map[string]any{{"ID": "friend-a", "Name": "합성 친구", "State": "active", "Generation": 7}},
		"Requests": []map[string]any{
			{"ID": "incoming", "Name": "받은 이름", "State": "pending", "Generation": 8, "Incoming": true, "Exp": "오늘 오후 2:10"},
			{"ID": "outgoing", "Name": "노출하지 않을 보낸 이름", "State": "pending", "Generation": 9, "Incoming": false, "Exp": "오늘 오후 2:10"},
		},
		"Conversations": []map[string]any{{"ID": "conversation-a", "Name": "합성 친구", "FromName": "노우", "ToName": "다닷", "State": "queued", "URL": "/home/conversation?from=agent_a&to=agent_c"}},
		"Messages": []map[string]any{
			{"ID": "message-a", "Author": "나", "Agent": "노우", "Short": "a001", "Text": "합성 메시지\n두 번째 줄", "When": "오늘 오후 2:10", "State": "queued", "Authorship": "owner"},
			{"ID": "message-b", "Author": "합성 친구", "Agent": "다닷", "Short": "c003", "Text": "합성 답장", "When": "오늘 오후 2:11", "State": "delivered", "Authorship": "agent"},
		},
		"CanSend": true, "Attempt": "synthetic-attempt", "Code": "synthetic-friend-code", "Recent": true,
		"PendingGates": 1, "ConversationID": "conversation-a", "Legacy": true,
	}
}

func renderHumanPage(t *testing.T, name string, data map[string]any) string {
	t.Helper()
	var b bytes.Buffer
	if err := memberPages.ExecuteTemplate(&b, name, data); err != nil {
		t.Fatal(name, err)
	}
	return b.String()
}

func TestHumanPagesRenderAndEscape(t *testing.T) {
	for _, name := range []string{"human-home", "human-friends", "human-friend", "human-agents", "human-account", "human-name", "human-new", "human-conversation"} {
		t.Run(name, func(t *testing.T) {
			v := humanPageFixture()
			v["Problem"] = "이름은 1–40자로 입력하세요"
			v["Notice"] = "합성 안내"
			body := renderHumanPage(t, name, v)
			for _, required := range []string{`<html lang="ko">`, `<title>합성 대화 · KnowsLink</title>`, `href="#main"`, `id="main"`, `role="alert"`, `role="status"`, `min-height:44px`, `outline:3px solid var(--accent)`} {
				if !strings.Contains(body, required) {
					t.Fatalf("missing %q", required)
				}
			}
			if strings.Count(body, "<h1") != 1 {
				t.Fatal("page must have one h1")
			}
			for _, banned := range []string{"<script", "<link", "<img", "url(", "@import", "<no value>"} {
				if strings.Contains(body, banned) {
					t.Fatalf("unexpected %q", banned)
				}
			}
			if dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR"); dir != "" {
				if err := os.WriteFile(filepath.Join(dir, name+".html"), []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	v := humanPageFixture()
	payload := `<script>alert("x")</script>`
	v["Name"], v["Text"] = payload, payload
	v["Messages"] = []map[string]any{{"Author": payload, "Agent": payload, "Text": payload, "Authorship": "owner", "State": "queued"}}
	for _, name := range []string{"human-name", "human-conversation"} {
		body := renderHumanPage(t, name, v)
		if strings.Contains(body, payload) || !strings.Contains(body, "&lt;script&gt;") {
			t.Fatal(name, "untrusted text must stay text")
		}
	}
}

func TestHumanFriendActionsMatchRequestDirection(t *testing.T) {
	body := renderHumanPage(t, "human-friends", humanPageFixture())
	for _, required := range []string{`action="/home/friends/decision"`, `name="generation" value="8"`, `name="decision" value="accept"`, `name="decision" value="deny"`, `name="decision" value="cancel"`, `type="hidden" name="consent" value="yes"`, `name="consent" value="yes" required`, "서로의 현재·앞으로 추가할 에이전트가 연결됩니다"} {
		if !strings.Contains(body, required) {
			t.Fatal("missing", required)
		}
	}
	if strings.Contains(body, "노출하지 않을 보낸 이름") {
		t.Fatal("outgoing request must not reveal code target")
	}
	for _, incoming := range []bool{true, false} {
		v := humanPageFixture()
		v["Requests"] = []map[string]any{{"ID": "request-a", "Name": "합성 이름", "State": "pending", "Generation": 3, "Incoming": incoming}}
		body := renderHumanPage(t, "human-friends", v)
		if strings.Contains(body, `value="accept"`) != incoming || strings.Contains(body, `value="deny"`) != incoming || strings.Contains(body, `value="cancel"`) == incoming {
			t.Fatal("request direction exposes wrong action")
		}
	}
	for _, state := range []string{"denied", "cancelled", "expired", "removed", "blocked"} {
		v := humanPageFixture()
		v["Requests"] = []map[string]any{{"ID": "request-a", "State": state, "Generation": 3, "Incoming": true}}
		body := renderHumanPage(t, "human-friends", v)
		if strings.Contains(body, `value="accept"`) || strings.Contains(body, `value="deny"`) || strings.Contains(body, `value="cancel"`) || strings.Contains(body, `value="unblock"`) != (state == "blocked") {
			t.Fatal(state, "finished requests must not offer pending actions")
		}
	}
}

func TestHumanConversationPinsTargetAndReadOnlyBoundary(t *testing.T) {
	v := humanPageFixture()
	body := renderHumanPage(t, "human-conversation", v)
	for _, required := range []string{`name="from" value="agent_a"`, `name="to" value="agent_c"`, `name="attempt" value="synthetic-attempt"`, `action="/home/conversation/send"`, `<textarea id="message" name="text" required`, "a001", "c003", "내가 작성", "에이전트가 작성", "메시지는 24시간 안에 사라집니다", "대기", "전달됨", "새로고침"} {
		if !strings.Contains(body, required) {
			t.Fatal("missing", required)
		}
	}
	if strings.Contains(body, "<select") || strings.Contains(body, "읽음") {
		t.Fatal("conversation must fix endpoints and avoid inferred read state")
	}
	v["CanSend"] = false
	body = renderHumanPage(t, "human-conversation", v)
	if strings.Contains(body, "<textarea") || strings.Contains(body, `action="/home/conversation/send"`) || !strings.Contains(body, "지금은 메시지를 보낼 수 없습니다") || !strings.Contains(body, "합성 메시지") {
		t.Fatal("revocation must leave only retained messages")
	}
	v["Messages"] = nil
	body = renderHumanPage(t, "human-conversation", v)
	if !strings.Contains(body, "보존 기간이 지난 대화입니다") {
		t.Fatal("expired history needs recovery text")
	}
	v["CanSend"], v["Problem"], v["Text"] = true, "전송 결과를 확인하세요", "합성 재시도 본문"
	body = renderHumanPage(t, "human-conversation", v)
	if !strings.Contains(body, "합성 재시도 본문</textarea>") || !strings.Contains(body, "attempt=synthetic-attempt") {
		t.Fatal("failed send must preserve text and attempt")
	}
}

func TestHumanEmptyChoicesOfferRecovery(t *testing.T) {
	for _, name := range []string{"human-home", "human-new", "human-agents"} {
		body := renderHumanPage(t, name, map[string]any{"Title": "빈 계정", "Section": "conversations"})
		if strings.Contains(body, "<select") || strings.Contains(body, "<textarea") || !strings.Contains(body, "첫 에이전트를 추가하세요") {
			t.Fatal(name, "empty account must lead to agent setup")
		}
	}
	v := humanPageFixture()
	v["Agents"] = nil
	for _, name := range []string{"human-new", "human-friend"} {
		body := renderHumanPage(t, name, v)
		if strings.Contains(body, "<select") || strings.Contains(body, "<textarea") || !strings.Contains(body, "연결된 에이전트가 없습니다") {
			t.Fatal(name, "no peer agent must hide composer and selects")
		}
	}
	v = humanPageFixture()
	body := renderHumanPage(t, "human-new", v)
	for _, required := range []string{`method="get" action="/home/conversation"`, `name="from"`, `name="to"`, "a001", "b002", "c003"} {
		if !strings.Contains(body, required) {
			t.Fatal("duplicate names need stable distinctions", required)
		}
	}
	v["Peer"], v["Agents"], v["To"] = nil, v["Own"], map[string]any{"ID": ""}
	if !strings.Contains(renderHumanPage(t, "human-new", v), `<option value="agent_b" selected>`) {
		t.Fatal("new self conversation must start with different endpoints")
	}
}

func TestHumanGateKeepsVerifiedBodyAndCSRF(t *testing.T) {
	v := map[string]any{"State": "pending", "From": "agent_a", "To": "agent_b", "Parent": "request-a", "Intent": "schedule.query", "Exp": "오늘 오후 2:10", "Policy": "owner_gate", "Hint": "합성 참고", "Body": `<script>untrusted</script>`, "Active": true, "CSRF": "synthetic-csrf", "Deny": "/home/gates/gate-a/deny", "Back": "/home"}
	body := renderHumanPage(t, "human-gate", v)
	for _, required := range []string{`name="csrf" value="synthetic-csrf"`, `name="decision" value="approve"`, `name="decision" value="deny"`, `formaction="/home/gates/gate-a/deny"`, "검증된 typed body", "&lt;script&gt;", "친구 동의나 대화 메시지는 이 요청의 승인 증명이 아닙니다"} {
		if !strings.Contains(body, required) {
			t.Fatal("missing", required)
		}
	}
	v["Active"] = false
	body = renderHumanPage(t, "human-gate", v)
	if strings.Contains(body, "<form") || !strings.Contains(body, "지금은 승인하거나 거절할 수 없습니다") {
		t.Fatal("inactive gate must not offer decision forms")
	}
}

func TestHumanAgentNameRecoveryAndSelfConsent(t *testing.T) {
	v := humanPageFixture()
	v["NewName"], v["EditID"], v["EditName"] = `새 "이름"`, "agent_a", `수정 <이름>`
	v["Legacy"], v["SelfConnected"] = false, false
	body := renderHumanPage(t, "human-agents", v)
	for _, want := range []string{`id="new-agent-name" name="name" required maxlength="40" value="새 &#34;이름&#34;"`, `id="agent-name-agent_a" name="name" required maxlength="40" value="수정 &lt;이름&gt;"`, `id="agent-name-agent_b" name="name" required maxlength="40" value="같은 이름"`, `action="/home/self-connect"`} {
		if !strings.Contains(body, want) {
			t.Fatal("missing recovery/consent", want)
		}
	}
	v["SelfConnected"] = true
	if strings.Contains(renderHumanPage(t, "human-agents", v), `action="/home/self-connect"`) {
		t.Fatal("self-connect consent must end after explicit opt-in")
	}
}
