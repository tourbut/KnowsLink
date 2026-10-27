// SAR-CONNECT-SIMPLE-001: keyless attempts are not agents, the 24h rate is the only per-member refusal, default
// names, and connection failures that say why without naming another account.
package relay

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// approvedDevice runs one client attempt up to browser consent and returns what completing it needs.
func approvedDevice(t *testing.T, st *State, now time.Time, name string) (token, proof string, c *Connection) {
	t.Helper()
	token, c, private := deviceFixture(t, st, now)
	c.Owner, c.State = "owner", "prepared"
	if err := st.approveDevice(hashToken(token), "owner", name, now); err != nil {
		t.Fatal(err)
	}
	return token, base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c))), c
}

func TestAttemptsAreNotAgents(t *testing.T) {
	st, now := connectionFixture()
	before := len(st.humanAgents("owner"))
	// Five keyless attempts: waiting, approved without a saved key, cancelled, expired, removed (CS-AC-04/05).
	_, waiting, _ := deviceFixture(t, st, now)
	waiting.Owner, waiting.State = "owner", "prepared"
	token, proof, approved := approvedDevice(t, st, now, "Claude Code")
	_, _, cancelled := approvedDevice(t, st, now, "취소")
	cancelled.State = "cancelled"
	_, expired, _ := deviceFixture(t, st, now)
	expired.Owner, expired.State, expired.Exp = "owner", "prepared", now.Add(-time.Minute)
	removed, _, _ := approvedDevice(t, st, now, "정리")
	delete(st.Connections, hashToken(removed))
	st.sweepConnections(now)
	if len(st.humanAgents("owner")) != before || st.Agents[approved.Agent] != nil {
		t.Fatal("an attempt became an agent before its key was saved")
	}
	states := map[string]bool{}
	for _, a := range st.ownAttempts("owner") {
		states[a.State] = true
	}
	if len(st.ownAttempts("owner")) != 3 || !states["prepared"] || !states["approved"] || !states["expired"] {
		t.Fatal("in-progress list", st.ownAttempts("owner"))
	}
	for i := range agentRate - 1 {
		st.Agents[fmt.Sprint("agent_n", i)] = &Agent{Owner: "owner", Keys: map[string]*Key{}, Created: now}
	}
	if err := st.agentLimit("owner", now); err != nil {
		t.Fatal("attempts counted toward the rate", err)
	}
	// Saving the key makes the agent with the approved name; it then counts and leaves the in-progress list.
	if _, err := st.completeConnection(token, supportedClient, proof, now); err != nil {
		t.Fatal(err)
	}
	if a := st.Agents[approved.Agent]; a == nil || a.Name != "Claude Code" || a.Created != now || activeKeys(a) != 1 {
		t.Fatal("completed agent", a)
	}
	if err := st.agentLimit("owner", now); err == nil {
		t.Fatal("completed agent not counted")
	}
	for _, a := range st.ownAttempts("owner") {
		if a.State == "approved" {
			t.Fatal("completed attempt still listed")
		}
	}
	// A-02: an approved attempt without a saved key is gone a day after its deadline.
	_, _, late := approvedDevice(t, newStateLike(st), now, "late")
	st.Connections["late"] = late
	st.sweepConnections(late.Exp.Add(24*time.Hour - time.Second))
	if st.Connections["late"] == nil || st.Connections["late"].State != "approved" {
		t.Fatal("approved attempt expired before cleanup")
	}
	st.sweepConnections(late.Exp.Add(24 * time.Hour))
	if st.Connections["late"] != nil {
		t.Fatal("keyless attempt kept past cleanup")
	}
}

// C1: an attempt approved by the old flow already had its agent. Once that agent is revoked and swept, a late
// completion is refused instead of making the same agent ID again.
func TestOldAttemptNeverRemakesAgent(t *testing.T) {
	st, now := connectionFixture()
	token, proof, c := approvedDevice(t, st, now, "old")
	c.Name = "" // the old flow made the agent at approval and kept no name on the attempt
	if _, err := st.completeConnection(token, supportedClient, proof, now); err == nil || err.Error() != "invalid_auth" || st.Agents[c.Agent] != nil {
		t.Fatal("swept agent ID made again", err)
	}
}

// newStateLike gives a fixture a fresh connection map with the same owner, so a helper's attempt stays separate.
func newStateLike(st *State) *State {
	other := newState()
	other.Owners = st.Owners
	return other
}

func TestAgentRate(t *testing.T) {
	st, now := connectionFixture()
	st.Agents["agent_a"].Keys["k"] = &Key{Changed: now, Credential: hashToken("live-credential")}
	oldest := now.Add(-time.Hour)
	for i := range agentRate {
		st.Agents[fmt.Sprint("agent_", i)] = &Agent{Owner: "owner", Keys: map[string]*Key{}, Created: oldest.Add(time.Duration(i) * time.Minute), Revoked: i == 0}
	}
	// The 21st within 24h is refused with the time the oldest leaves the window; the revoked one still counts.
	err := st.agentLimit("owner", now)
	if r, ok := err.(agentRefused); !ok || r.Code != "agent_rate" || !r.Retry.Equal(oldest.Add(24*time.Hour)) {
		t.Fatal("21st agent", err)
	}
	if body := errorBody(err); body["error"] != "agent_rate" || body["retry_at"] != oldest.Add(24*time.Hour).UTC().Add(time.Second-1).Truncate(time.Second).Format(time.RFC3339) {
		t.Fatal("client refusal", body)
	}
	if statusFor(err) != 429 || !strings.Contains(safeProblem(err), clock(oldest.Add(24*time.Hour))) {
		t.Fatal("refusal text", safeProblem(err))
	}
	// Approval is refused and recorded so the polling client can report it; nothing else changes.
	token, c, _ := deviceFixture(t, st, now)
	c.Owner, c.State = "owner", "prepared"
	if err := st.approveDevice(hashToken(token), "owner", "21번째", now); err == nil || c.State != "refused" || c.Refused != "agent_rate" || c.Agent != "" || c.open(now) {
		t.Fatal("refused approval", err, c.State)
	}
	if a, err := st.principal("live-credential", "agent"); err != nil || a != "agent_a" {
		t.Fatal("existing agent changed", err)
	}
	if err := st.agentLimit("owner", oldest.Add(24*time.Hour)); err != nil {
		t.Fatal("slot not freed after 24h", err)
	}
}

func TestDefaultNames(t *testing.T) {
	st, now := connectionFixture()
	st.Agents["agent_a"].Name = "Claude Code"
	if got := st.defaultName("owner", "claude", "ko"); got != "Claude Code 2" {
		t.Fatal(got)
	}
	_, _, _ = approvedDevice(t, st, now, "Claude Code 2")
	if got := st.defaultName("owner", "claude", "ko"); got != "Claude Code 3" {
		t.Fatal("approved attempt holds its name", got)
	}
	st.Agents["agent_a"].Revoked = true
	if got := st.defaultName("owner", "claude", "en"); got != "Claude Code" {
		t.Fatal("revoked agent holds its name", got)
	}
	for host, want := range map[string]string{"codex": "Codex", "grok": "Grok Bot", "dots": "Dots", "": "내 에이전트 1", "unknown": "내 에이전트 1"} {
		if got := st.defaultName("owner", host, "ko"); got != want {
			t.Fatal(host, got)
		}
	}
	st.Agents["agent_m"] = &Agent{Owner: "owner", Name: "My agent 1", Keys: map[string]*Key{}}
	if got := st.defaultName("owner", "", "en"); got != "My agent 2" {
		t.Fatal(got)
	}
	// Another member's names never number mine.
	st.Agents["agent_o"] = &Agent{Owner: "other", Name: "Codex", Keys: map[string]*Key{}}
	if got := st.defaultName("owner", "codex", "ko"); got != "Codex" {
		t.Fatal(got)
	}
	// Only a known host is kept as the hint.
	token, c, private := deviceFixture(t, newState(), now)
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
	other := newState()
	if got, err := other.startDevice(token, c.Client, c.Kid, c.Public, proof, "../../evil", now); err != nil || got.Host != "" {
		t.Fatal("unknown host kept", err)
	}
}

func TestConnectFailureMessages(t *testing.T) {
	connection := &googleAttempt{Connection: "c"}
	seen := map[string]bool{}
	for problem, retry := range map[string]bool{"invalid_auth": true, "capacity": true, "expired": false} {
		text, again := connectFailure(connection, problem)
		if text == "" || again != retry || seen[text] || strings.Contains(text, "@") {
			t.Fatal(problem, text, again)
		}
		seen[text] = true
		if tr("en", text) == text {
			t.Fatal("untranslated", text)
		}
	}
	if text, _ := connectFailure(&googleAttempt{}, "invalid_auth"); text != "" {
		t.Fatal("plain login got a connection message")
	}
	if text, _ := connectFailure(connection, "identity_mismatch"); text != "" {
		t.Fatal("other failures keep the current text")
	}
	w := httptest.NewRecorder()
	render(w, 403, "refusal", map[string]any{"Title": "Google 로그인", "Problem": "p", "Back": "/", "Retry": "conn-id"})
	if body := w.Body.String(); !strings.Contains(body, `action="/auth/google"`) || !strings.Contains(body, `value="conn-id"`) || !strings.Contains(body, "다른 계정 선택") {
		t.Fatal("account chooser missing")
	}
}

func TestConnectScreens(t *testing.T) {
	v := humanPageFixture()
	v["Attempts"] = []ownAttempt{{ID: "attempt-1", Client: "Claude Code", Started: "10:42", State: "approved"}}
	v["Unnamed"] = []humanAgent{{ID: "agent_u", Short: "a1b2"}}
	body := renderHumanPage(t, "human-agents", v)
	attempts, unnamed, list := strings.Index(body, "진행 중인 연결"), strings.Index(body, "이름을 정할 에이전트"), strings.Index(body, "내 에이전트와 연결 상태")
	if attempts < 0 || unnamed < attempts || list < unnamed {
		t.Fatal("M-01 order")
	}
	for _, want := range []string{"Claude Code · 10:42 시작", "동의 완료 · 키 저장 확인 필요", `action="/home/attempt-remove"`, `value="attempt-1"`, "이름 없는 에이전트 · a1b2", `action="/home/agent-name"`, `value="agent_u"`} {
		if !strings.Contains(body, want) {
			t.Fatal("missing", want)
		}
	}
	v["Attempts"], v["Unnamed"] = []ownAttempt{}, []humanAgent{}
	body = renderHumanPage(t, "human-agents", v)
	for _, banned := range []string{"진행 중인 연결", "이름을 정할 에이전트", "회원당", "먼저 철회"} {
		if strings.Contains(body, banned) {
			t.Fatal("empty section or limit text shown", banned)
		}
	}
}
