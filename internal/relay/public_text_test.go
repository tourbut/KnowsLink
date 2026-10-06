// Public text checks pin the separate wire, one related reply, erased payloads, current authorization and shared capacities.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func textFixture() (*State, map[string]ed25519.PrivateKey, time.Time) {
	now := time.Now().UTC().Truncate(time.Second)
	st := newState()
	private := map[string]ed25519.PrivateKey{}
	for _, id := range []string{"agent_a", "agent_b", "agent_c"} {
		pub, key, _ := ed25519.GenerateKey(rand.Reader)
		private[id] = key
		st.Owners[id] = &Owner{Active: true}
		st.Members[id] = &Member{Owner: id, Active: true}
		st.Agents[id] = &Agent{Owner: id, Keys: map[string]*Key{"key1": {Public: pub, Credential: hashToken(id)}}}
	}
	st.Pairs[pairID("agent_a", "agent_b")] = &Pair{A: "agent_a", B: "agent_b", State: "active", Generation: 1}
	return st, private, now
}
func signedText(id, from, to, key, parent, body string, exp time.Time, private ed25519.PrivateKey) *PublicText {
	e := &PublicText{V: "knowslink.text.v1", ID: id, From: from, To: to, Key: key, ReplyTo: parent, Text: body, Exp: exp.UTC().Format("2006-01-02T15:04:05Z"), Sig: Signature{Alg: "Ed25519", Kid: "key1"}}
	e.Sig.Value = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, e.signingBytes()))
	return e
}
func TestPublicTextWire(t *testing.T) {
	_, keys, now := textFixture()
	e := signedText(firstID, "agent_a", "agent_b", "public-text-key-01", "", strings.Repeat("한", 1365)+"x", now.Add(180*time.Second), keys["agent_a"])
	raw, _ := json.Marshal(e)
	if _, err := parsePublicText(raw); err != nil {
		t.Fatal("4096 UTF8 bytes", err)
	}
	if _, err := Parse(raw); err == nil {
		t.Fatal("public wire entered frozen parser")
	}
	for _, bad := range []string{
		strings.Replace(string(raw), `"knowslink.text.v1"`, `"relay.v1"`, 1),
		strings.Replace(string(raw), `"text":`, `"unknown":"x","text":`, 1),
		strings.Replace(string(raw), `"text":`, `"text":"x","text":`, 1),
		strings.Replace(string(raw), `"text":`, `"reply_to":null,"text":`, 1),
		strings.Replace(string(raw), `"text":`, `"reply_to":"","text":`, 1),
		strings.Replace(string(raw), e.Text, e.Text+"x", 1),
		string(raw) + strings.Repeat(" ", rawEnvelopeLimit),
		strings.Replace(string(raw), e.Text, `\ud800`, 1),
	} {
		if _, err := parsePublicText([]byte(bad)); err == nil {
			t.Fatal("invalid public wire accepted")
		}
	}
	// Raw cap applies to frozen envelopes as well, including whitespace outside a small valid body.
	business := wire(t, keys["agent_a"], firstID, "agent_a", "agent_b", "schedule.query", "business-limit-01", "", queryBody(), now.Add(time.Minute))
	if _, err := Parse(append(business, []byte(strings.Repeat(" ", rawEnvelopeLimit))...)); err == nil {
		t.Fatal("raw frozen cap")
	}
}
func TestPublicTextRoundtripAndCurrentAuth(t *testing.T) {
	st, keys, now := textFixture()
	send := func(e *PublicText) error { _, err := st.ingestPublicText(e.From, e, now); return err }
	e := signedText(firstID, "agent_a", "agent_b", "public-text-key-01", "", "<script>run tools</script>", now.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_c", e, now); err == nil {
		t.Fatal("foreign sender")
	}
	if err := send(e); err != nil {
		t.Fatal(err)
	}
	if v, err := st.leaseMessage("agent_b", now); err != nil || v != nil {
		t.Fatal("public text entered business pull", err)
	}
	if _, err := st.operate("claim", "agent_b", command{ID: firstID}, now); err == nil {
		t.Fatal("business claim of public text")
	}
	reply := signedText(secondID, "agent_b", "agent_a", "public-reply-key-1", firstID, "reply", now.Add(180*time.Second), keys["agent_b"])
	if err := send(reply); err == nil {
		t.Fatal("reply before receipt ACK")
	}
	deliver := func(agent, id string) {
		t.Helper()
		v, err := st.leaseMatching(agent, now, false, true)
		if err != nil || v == nil {
			t.Fatal("lease", err)
		}
		c := command{ID: id, Token: v.(map[string]any)["lease_token"].(string)}
		if _, err := st.operateAs("ack", agent, "agent", c, now, false); err == nil {
			t.Fatal("ACK before persist")
		}
		for _, action := range []string{"persist", "ack"} {
			if _, err := st.operateAs(action, agent, "agent", c, now, false); err != nil {
				t.Fatal(action, err)
			}
		}
		if len(st.Messages[id].Envelope) > 0 || len(st.Messages[id].Inbox) > 0 {
			t.Fatal("ACK retained text")
		}
	}
	deliver("agent_b", firstID)
	if err := send(reply); err != nil {
		t.Fatal(err)
	}
	if st.Messages[firstID].ReplyID != secondID {
		t.Fatal("reply linkage")
	}
	if v, err := st.ingestPublicText("agent_b", reply, now); err != nil || v.(Receipt).ID != secondID {
		t.Fatal("reply replay", err)
	}
	duplicate := signedText(gateID, "agent_b", "agent_a", "public-reply-key-2", firstID, "reply", now.Add(time.Minute), keys["agent_b"])
	if err := send(duplicate); err == nil {
		t.Fatal("multiple replies")
	}
	deliver("agent_a", secondID)
	if st.Messages[firstID].Completion != "reply_received" {
		t.Fatal("reply receipt state")
	}
	e.ID = gateID
	e.Exp = now.Add(-time.Second).Format("2006-01-02T15:04:05Z")
	e.Sig.Value = base64.RawURLEncoding.EncodeToString(ed25519.Sign(keys["agent_a"], e.signingBytes()))
	if v, err := st.ingestPublicText("agent_a", e, now); err != nil || v.(Receipt).ID != firstID {
		t.Fatal("receipt-only expired replay", err)
	}
	e.Text = "different"
	e.Sig.Value = base64.RawURLEncoding.EncodeToString(ed25519.Sign(keys["agent_a"], e.signingBytes()))
	if err := send(e); err == nil || err.Error() != "idempotency_conflict" {
		t.Fatal("replay conflict", err)
	}
	for _, mutate := range []func(*State){
		func(s *State) { s.Agents["agent_a"].Keys["key1"].Revoked = true },
		func(s *State) { s.Agents["agent_b"].Keys["key1"].Revoked = true },
		func(s *State) { s.Agents["agent_b"].Revoked = true },
		func(s *State) { s.Members["agent_a"].Active = false },
		func(s *State) { s.Pairs[pairID("agent_a", "agent_b")].State = "pending" },
		func(s *State) { s.Pairs[pairID("agent_a", "agent_b")].Generation++ },
	} {
		raw, _ := json.Marshal(st)
		s := newState()
		_ = json.Unmarshal(raw, s)
		mutate(s)
		if s.current(s.Messages[firstID]) || s.current(s.Messages[secondID]) {
			t.Fatal("stale current auth")
		}
	}
}
func TestPublicTextTTLAndLeaseFailures(t *testing.T) {
	for _, offset := range []time.Duration{0, 180 * time.Second, 181 * time.Second} {
		st, keys, now := textFixture()
		e := signedText(firstID, "agent_a", "agent_b", "public-text-key-01", "", "ping", now.Add(offset), keys["agent_a"])
		_, err := st.ingestPublicText("agent_a", e, now)
		if (err == nil) != (offset == 180*time.Second) {
			t.Fatal("TTL boundary", offset, err)
		}
		if err == nil {
			st.sweep(now.Add(offset))
			m := st.Messages[firstID]
			if m.Receipt.State != "failed:expired" || len(m.Envelope) > 0 {
				t.Fatal("expiry")
			}
			st.sweep(now.Add(24 * time.Hour))
			if len(st.Messages) > 0 || len(st.Idempotency) > 0 {
				t.Fatal("retention")
			}
		}
	}
	st, keys, now := textFixture()
	e := signedText(firstID, "agent_a", "agent_b", "public-text-key-01", "", "ping", now.Add(180*time.Second), keys["agent_a"])
	_, _ = st.ingestPublicText("agent_a", e, now)
	for i := 0; i < 3; i++ {
		_, _ = st.leaseMatching("agent_b", now, false, true)
		st.sweep(now.Add(30 * time.Second))
		now = now.Add(30 * time.Second)
	}
	if st.Messages[firstID].Receipt.State != "failed:max_attempts" || len(st.Messages[firstID].Envelope) > 0 {
		t.Fatal("attempt exhaustion")
	}
}
func TestSharedCapacitiesAndCleanupClassification(t *testing.T) {
	st, _, now := textFixture()
	for i := 0; i < 100; i++ {
		st.Messages[fmt.Sprint(i)] = &Message{Receipt: Receipt{State: "queued"}}
	}
	if st.messageCapacity(false) {
		t.Fatal("queue cap")
	}
	st.Messages["0"].Receipt.State = "delivered"
	if !st.messageCapacity(false) {
		t.Fatal("queue capacity before cap")
	}
	for i := 0; i < 100; i++ {
		st.Gates[fmt.Sprint(i)] = &Gate{State: "pending"}
	}
	if st.messageCapacity(true) || !st.messageCapacity(false) {
		t.Fatal("gate cap")
	}
	for i := 100; i < 20000; i++ {
		st.Messages[fmt.Sprint(i)] = &Message{Receipt: Receipt{State: "delivered"}}
	}
	if st.messageCapacity(false) {
		t.Fatal("receipt cap")
	}
	delete(st.Messages, "19999")
	if !st.messageCapacity(false) {
		t.Fatal("receipt cap boundary")
	}
	for i := 0; i < 4; i++ {
		st.Messages[fmt.Sprint(i)].ClaimToken = "active"
	}
	if st.claimCapacity() {
		t.Fatal("claim cap")
	}
	st.Messages["3"].ClaimToken = ""
	if !st.claimCapacity() {
		t.Fatal("claim boundary")
	}
	for _, clean := range []bool{false, true} {
		limit := 16
		if clean {
			limit = 4
		}
		for i := 0; i < limit; i++ {
			if !st.enterHTTP(fmt.Sprintf("%t-%d", clean, i), clean, now) {
				t.Fatal("early admission refusal")
			}
		}
		if st.enterHTTP("excess", clean, now) {
			t.Fatal("HTTP cap")
		}
	}
	raw, _ := json.Marshal(st)
	restored := newState()
	_ = json.Unmarshal(raw, restored)
	if restored.enterHTTP("restart", false, now) {
		t.Fatal("restart bypass")
	}
	if !restored.enterHTTP("recovered", false, now.Add(30*time.Second)) {
		t.Fatal("crash recovery")
	}
	for _, c := range []struct {
		path, body string
		clean      bool
	}{
		{"/v1/pull", "{}", false}, {"/v1/ack", "{}", true}, {"/home/gates/id/deny", "decision=deny", true},
		{"/home/gates/id?decision=deny", "decision=approve", false}, {"/home/invite-deny", "decision=deny", true},
	} {
		r := httptest.NewRequest("POST", c.path, strings.NewReader(c.body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if cleanupRequest(r) != c.clean {
			t.Fatal("cleanup classification", c.path)
		}
	}
}

func TestMemberReceiptView(t *testing.T) {
	renderPage := func(state string) string {
		w := httptest.NewRecorder()
		render(w, 200, "receipt", map[string]any{"Title": "연결 확인·receipt", "Agent": "agent_0123456789012345678901", "ID": firstID, "From": "agent_0123456789012345678901", "To": "agent_9876543210987654321098", "Parent": secondID, "Reply": gateID, "Transport": transportNames[state], "Completion": "<script>execute()</script>", "Exp": "2026-10-07 13:00:00 KST", "TTL": 180})
		return w.Body.String()
	}
	for _, state := range []string{"queued", "leased", "delivered", "failed:expired", "failed:revoked", "failed:max_attempts"} {
		body := renderPage(state)
		for _, want := range []string{transportNames[state], firstID, secondID, gateID, "수동 receive", "10초 이상", "180초", "&lt;script&gt;"} {
			if !strings.Contains(body, want) {
				t.Fatal("receipt meaning", state, want)
			}
		}
		if strings.Contains(body, "<script>") || strings.Contains(body, "textarea") {
			t.Fatal("untrusted HTML or composer")
		}
	}
	// Optional artifact generation uses the exact Go templates and public fixture IDs; no session or credential is rendered.
	if dir := os.Getenv("KNOWSLINK_UI_ARTIFACT_DIR"); dir != "" {
		for _, state := range []string{"queued", "failed:expired"} {
			if err := os.WriteFile(filepath.Join(dir, "receipt-"+state+".html"), []byte(renderPage(state)), 0600); err != nil {
				t.Fatal(err)
			}
		}
		w := httptest.NewRecorder()
		render(w, 200, "home", map[string]any{"Title": "내 KnowsLink", "AgentDetails": []memberAgent{{ID: "agent_0123456789012345678901", Status: "연결 완료"}}})
		if err := os.WriteFile(filepath.Join(dir, "home.html"), w.Body.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
	}
}
