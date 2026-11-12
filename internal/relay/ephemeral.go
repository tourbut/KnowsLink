// Messages, their idempotency records and Dots Inbox copies live only in volatile Redis, never in Postgres (D049-D052).
package relay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// messageRetention bounds every message record kept in Redis, including the web conversation copy (D049).
const messageRetention = time.Hour

// ephemeral is the Redis value. The Postgres row keeps only its random key name (State.Ephemeral), so a row commit
// switches to the copy written under that commit's row lock, and a rolled-back commit leaves the previous copy current.
type ephemeral struct {
	Messages    map[string]*Message
	Idempotency map[string]string
	Inboxes     map[string]map[string]json.RawMessage `json:",omitempty"`
}

func ephemeralKey(name string) string { return "knowslink:messages:" + name }

func (st *State) ephemeralPart() (ephemeral, bool) {
	e := ephemeral{Messages: st.Messages, Idempotency: st.Idempotency}
	if st.Remote != nil {
		for id, i := range st.Remote.Installs {
			if len(i.Inbox) > 0 {
				if e.Inboxes == nil {
					e.Inboxes = map[string]map[string]json.RawMessage{}
				}
				e.Inboxes[id] = i.Inbox
			}
		}
	}
	return e, len(e.Messages) == 0 && len(e.Idempotency) == 0 && len(e.Inboxes) == 0
}

// readEphemeral returns the copy a row points to. A missing key (TTL or a Redis restart) is an empty store; an
// unreachable or unconfigured Redis is unavailable, never a fallback to Postgres.
func readEphemeral(ctx context.Context, rdb *redis.Client, name string) ([]byte, error) {
	if rdb == nil {
		return nil, fault("unavailable")
	}
	if name == "" {
		return []byte{}, nil
	}
	raw, err := rdb.Get(ctx, ephemeralKey(name)).Bytes()
	if errors.Is(err, redis.Nil) {
		return []byte{}, nil
	}
	if err != nil {
		return nil, fault("unavailable")
	}
	return raw, nil
}

func (st *State) attachEphemeral(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	var e ephemeral
	if json.Unmarshal(raw, &e) != nil {
		return fault("unavailable")
	}
	if e.Messages != nil {
		st.Messages = e.Messages
	}
	if e.Idempotency != nil {
		st.Idempotency = e.Idempotency
	}
	for id, inbox := range e.Inboxes {
		if st.Remote != nil && st.Remote.Installs[id] != nil {
			st.Remote.Installs[id].Inbox = inbox
		}
	}
	return nil
}

// writeEphemeral stores a changed copy under a new key and points st at it; it returns the key the commit replaces.
// The write happens under the row lock, so concurrent relay processes see the same serial order as before.
func writeEphemeral(ctx context.Context, rdb *redis.Client, st *State, loaded []byte) (string, error) {
	e, empty := st.ephemeralPart()
	raw, err := json.Marshal(e)
	if err != nil {
		return "", fault("unavailable")
	}
	old := st.Ephemeral
	if empty {
		st.Ephemeral = ""
		return old, nil
	}
	if bytes.Equal(raw, loaded) {
		return "", nil
	}
	name := randomToken()
	// A full Redis (noeviction) refuses the write, so the request fails as unavailable instead of dropping messages.
	if rdb.Set(ctx, ephemeralKey(name), raw, messageRetention).Err() != nil {
		return "", fault("unavailable")
	}
	st.Ephemeral = name
	return old, nil
}

// retire deletes the replaced copy after its row commit, so Redis holds at most two copies during one commit.
// Failure leaves the retention TTL as the bound.
func retire(ctx context.Context, rdb *redis.Client, name string) {
	if name != "" {
		_ = rdb.Del(ctx, ephemeralKey(name)).Err()
	}
}
