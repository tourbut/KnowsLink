package relay

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// One owner's retained Redis bytes refuse only that owner's new text; another owner still sends, and the shared
// budget refuses everyone (SAR-REDIS-CAP-001).
func TestOwnerByteLimit(t *testing.T) {
	st, keys, now := textFixture()
	filler := json.RawMessage(`"` + strings.Repeat("x", 64<<10) + `"`)
	held := &Message{Receipt: Receipt{From: "agent_a", State: "delivered"}, Envelope: filler, Inbox: filler}
	per := messageBytes(held)
	n := 0
	for ; (n+1)*per <= ownerMessageBytesLimit-16<<10; n++ { // leave less room than one 4 KiB text needs
		st.Messages[fmt.Sprint("held", n)] = held
	}
	if _, mine := st.retainedBytes("agent_a"); mine != n*per {
		t.Fatal("owner bytes", mine, n*per)
	}
	big := strings.Repeat("<", 4096) // escapes to 6 bytes each in every Redis copy
	if _, err := st.ingestPublicText("agent_a", signedText(firstID, "agent_a", "agent_b", "byte-limit-key-01", "", big, now.Add(time.Minute), keys["agent_a"]), now); err == nil || err.Error() != "capacity" {
		t.Fatal("owner over its byte share", err)
	}
	if _, err := st.ingestPublicText("agent_b", signedText(secondID, "agent_b", "agent_a", "byte-limit-key-02", "", big, now.Add(time.Minute), keys["agent_b"]), now); err != nil {
		t.Fatal("other owner must keep sending", err)
	}
	delete(st.Messages, "held0")
	if _, err := st.ingestPublicText("agent_a", signedText(firstID, "agent_a", "agent_b", "byte-limit-key-01", "", big, now.Add(time.Minute), keys["agent_a"]), now); err != nil {
		t.Fatal("owner byte boundary", err)
	}
	// Records of no current owner count only against the shared budget, which then refuses every owner.
	for i := 0; i*per <= messageBytesLimit; i++ {
		st.Messages[fmt.Sprint("orphan", i)] = &Message{Receipt: Receipt{From: "agent_gone", State: "delivered"}, Envelope: filler, Inbox: filler}
	}
	if st.messageCapacity(false, "agent_c", 0) || st.messageCapacity(false, "", 0) {
		t.Fatal("shared byte budget")
	}
}

// messageBytes never undercounts the Redis JSON, for escaping-heavy bodies and the largest metadata, and a queued
// envelope reserves its later Inbox copy.
func TestMessageBytesBoundStoredJSON(t *testing.T) {
	st, keys, now := textFixture()
	st.Agents["agent_a"].Name = strings.Repeat("한", 40)
	st.Agents["agent_b"].Name = strings.Repeat(" ", 40)
	for _, m := range st.Members {
		m.Name, m.Issuer = strings.Repeat(" ", 40), emailIssuer
	}
	bodies := []string{strings.Repeat("<", 4096), strings.Repeat(" ", 1365), strings.Repeat(`"`, 4096), strings.Repeat("\n", 4096), "hi"}
	stamp := now.Add(time.Hour).UTC()
	for i, body := range bodies {
		id := fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012d", i)
		e := signedText(id, "agent_a", "agent_b", strings.Repeat("k", 120)+fmt.Sprint(i), "", body, now.Add(time.Minute), keys["agent_a"])
		if _, err := st.ingestPublicText("agent_a", e, now); err != nil {
			t.Fatal(i, err)
		}
		m := st.Messages[id]
		m.Kid, m.Parent, m.ReplyID, m.Completion = strings.Repeat("K", 128), id, id, "reply_received"
		m.LeaseToken, m.ClaimToken, m.LeaseUntil, m.SessionSeen = randomToken(), randomToken(), stamp, stamp
		m.Generation, m.LeaseGeneration, m.Attempts = 1<<62, 1<<62, 3
	}
	queued, _ := st.retainedBytes("")
	for _, m := range st.Messages { // persist: the reserved copy becomes real
		m.Inbox = append(json.RawMessage(nil), m.Envelope...)
	}
	persisted, _ := st.retainedBytes("")
	e, _ := st.ephemeralPart()
	raw, _ := json.Marshal(e)
	if queued != persisted || len(raw) > persisted {
		t.Fatal("bytes bound", queued, persisted, len(raw))
	}
	if persisted > 2*len(raw) {
		t.Fatal("bound too loose", persisted, len(raw))
	}
}
