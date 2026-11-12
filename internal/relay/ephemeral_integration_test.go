//go:build integration

// Message records reach only Redis: Postgres tables never hold them, the Redis copy expires, and a Redis outage is unavailable.
package relay

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// postgresHolds reports whether any row of any public table contains needle as text.
func postgresHolds(t *testing.T, pool *pgxpool.Pool, needle string) bool {
	t.Helper()
	ctx := context.Background()
	rows, err := pool.Query(ctx, "SELECT quote_ident(table_name) FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		tables = append(tables, name)
	}
	if rows.Err() != nil || len(tables) == 0 {
		t.Fatal("no tables to scan", rows.Err())
	}
	for _, table := range tables {
		var dump string
		if err = pool.QueryRow(ctx, "SELECT COALESCE(string_agg(t::text, ''), '') FROM "+table+" t").Scan(&dump); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(dump, needle) {
			return true
		}
	}
	return false
}

func TestMessagesLiveOnlyInRedis(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	ctx := context.Background()
	const id, needle = "0199a3f2-4c10-7a11-8b22-ephemeral0001", "합성 휘발 본문 ephemeral-needle"
	envelope := json.RawMessage(`{"text":"` + needle + `"}`)
	if _, err := f.s.transaction(ctx, func(st *State, now time.Time) (any, error) {
		st.Messages[id] = &Message{Receipt: Receipt{ID: id, From: "agent_a", To: "agent_b", Intent: publicTextIntent, Accepted: now, Exp: now.Add(publicTextTTL), State: "queued"},
			Key: "ephemeral-key-01", Kid: "key1", Deliver: "agent", Envelope: envelope, Inbox: envelope, Conversation: &ConversationCopy{Text: needle, Authorship: "agent"}}
		st.Idempotency["agent_a/ephemeral-key-01"] = id
		st.remote().Installs["install_ephemeral"] = &remoteInstall{Owner: "x", Agent: "agent_none", Exp: now.Add(time.Hour), Inbox: map[string]json.RawMessage{id: envelope}}
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}
	if postgresHolds(t, pool, needle) || postgresHolds(t, pool, id) {
		t.Fatal("a message, its copy or its id reached Postgres")
	}
	var name string
	if err := pool.QueryRow(ctx, "SELECT state->>'Ephemeral' FROM relay_state").Scan(&name); err != nil || name == "" {
		t.Fatal("row does not point at a Redis copy", err)
	}
	raw, err := testRedis().Get(ctx, ephemeralKey(name)).Result()
	if err != nil || strings.Count(raw, needle) < 4 {
		t.Fatal("Redis copy lacks the envelope, Inbox, conversation or Dots Inbox copy", err)
	}
	if ttl := testRedis().TTL(ctx, ephemeralKey(name)).Val(); ttl <= 59*time.Minute || ttl > messageRetention {
		t.Fatalf("Redis copy TTL %v, want at most %v", ttl, messageRetention)
	}

	// An unchanged store writes no new copy; a later process reads the same records from Redis.
	restarted := &Service{Redis: testRedis(), Pool: pool}
	var seen bool
	if _, err = restarted.transaction(ctx, func(st *State, _ time.Time) (any, error) {
		seen = st.Messages[id] != nil && st.Idempotency["agent_a/ephemeral-key-01"] == id
		return nil, nil
	}); err != nil || !seen {
		t.Fatal("records not read back from Redis", err)
	}

	// An unreachable Redis makes the request unavailable and commits nothing, including no Postgres fallback.
	var before, after int64
	_ = pool.QueryRow(ctx, "SELECT epoch FROM relay_state").Scan(&before)
	down := &Service{Redis: redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", DialTimeout: 200 * time.Millisecond, MaxRetries: -1}), Pool: pool}
	if _, err = down.transaction(ctx, func(*State, time.Time) (any, error) { return nil, nil }); err == nil || err.Error() != "unavailable" {
		t.Fatalf("Redis outage: %v", err)
	}
	if _, err = (&Service{Pool: pool}).transaction(ctx, func(*State, time.Time) (any, error) { return nil, nil }); err == nil || err.Error() != "unavailable" {
		t.Fatalf("unconfigured Redis: %v", err)
	}
	_ = pool.QueryRow(ctx, "SELECT epoch FROM relay_state").Scan(&after)
	if after != before {
		t.Fatal("a transaction committed without Redis")
	}

	// A Redis restart loses the copy: the next transaction sees an empty store and new records work again.
	if err = pool.QueryRow(ctx, "SELECT state->>'Ephemeral' FROM relay_state").Scan(&name); err != nil || name == "" {
		t.Fatal("row lost its Redis pointer", err)
	}
	if err = testRedis().Del(ctx, ephemeralKey(name)).Err(); err != nil {
		t.Fatal(err)
	}
	if _, err = f.s.transaction(ctx, func(st *State, now time.Time) (any, error) {
		if len(st.Messages) != 0 || len(st.Idempotency) != 0 {
			t.Fatal("records survived the lost Redis copy")
		}
		st.Messages[id] = &Message{Receipt: Receipt{ID: id, From: "agent_a", To: "agent_b", Accepted: now.Add(-messageRetention), Exp: now}}
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}
	// The retention sweep drops a record one hour after acceptance and clears the pointer once the store is empty.
	if _, err = f.s.transaction(ctx, func(st *State, _ time.Time) (any, error) {
		if len(st.Messages) != 0 {
			t.Fatal("record kept past the one-hour retention")
		}
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, "SELECT COALESCE(state->>'Ephemeral', '') FROM relay_state").Scan(&name); err != nil || name != "" {
		t.Fatalf("empty store still points at %q: %v", name, err)
	}
}
