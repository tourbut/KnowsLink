// The Postgres encoding of State omits every message record; the Redis encoding carries all of them.
package relay

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

// restart mirrors a relay restart that keeps Redis: the Postgres row and the Redis copy are decoded into a new State.
func restart(t *testing.T, st *State) *State {
	t.Helper()
	row, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	records, _ := st.ephemeralPart()
	stored, err := json.Marshal(records)
	if err != nil {
		t.Fatal(err)
	}
	next := newState()
	if err = json.Unmarshal(row, next); err != nil || next.attachEphemeral(stored) != nil {
		t.Fatal("restart decode", err)
	}
	return next
}

func TestStateEncodingSplitsMessages(t *testing.T) {
	st, _ := retentionFixture(3)
	const secret = "synthetic retained conversation"
	st.remote().Installs["install_1"] = &remoteInstall{Agent: "agent_000_0", Exp: time.Now(), Inbox: map[string]json.RawMessage{"message_00000": json.RawMessage(`{"text":"inbox only"}`)}}
	row, _ := json.Marshal(st)
	for _, leak := range []string{secret, "inbox only", "message_00000", `"Messages"`, `"Idempotency"`, `"Inbox"`} {
		if bytes.Contains(row, []byte(leak)) {
			t.Fatalf("Postgres encoding holds %q", leak)
		}
	}
	next := restart(t, st)
	if len(next.Messages) != 3 || len(next.Idempotency) != 3 || next.Messages["message_00000"].Conversation.Text != secret ||
		string(next.Remote.Installs["install_1"].Inbox["message_00000"]) != `{"text":"inbox only"}` {
		t.Fatal("Redis encoding lost records")
	}
	if _, empty := newState().ephemeralPart(); !empty {
		t.Fatal("new state is not an empty store")
	}
}
