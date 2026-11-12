//go:build integration

// Migration 00003 keeps an existing jsonb row readable and makes older writers fail; 00004 deletes its message copies.
package relay

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func TestStateFormatMigration(t *testing.T) {
	ctx := context.Background()
	url := os.Getenv("TEST_DATABASE_URL")
	admin, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = admin.Close(ctx) }()
	name := fmt.Sprintf("kl_state_migration_%d", time.Now().UnixNano())
	if _, err = admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = admin.Exec(ctx, "DROP DATABASE "+name+" WITH (FORCE)") })
	cfg, err := PoolConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.Database = name
	sqlDB := stdlib.OpenDB(*cfg.ConnConfig)
	defer func() { _ = sqlDB.Close() }()
	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, os.DirFS("../../db/migrations"), goose.WithDisableGlobalRegistry(true))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = provider.UpTo(ctx, 2); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	// The legacy row: retained history big enough to compress, and a signed envelope whose text JSON may escape.
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	legacy, _ := retentionFixture(1000)
	const id = "0199a3f2-4c10-7a11-8b22-334455667788"
	signed, err := Parse(wire(t, private, id, "agent_000_0", "agent_000_1", "relay.test.message", "migration-key-01", "", map[string]any{"text": "<b>& 한글</b>"}, time.Now().Add(time.Hour)))
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := Canonical(signed.Raw)
	if err != nil {
		t.Fatal(err)
	}
	// Same sender, recipient and authority as a fixture message, so the first sweep keeps the envelope as current.
	sibling := legacy.Messages["message_00000"]
	legacy.Messages[id] = &Message{Receipt: Receipt{ID: id, From: signed.From, To: signed.To, Digest: signed.Digest()}, Envelope: envelope, Kid: sibling.Kid, Scope: sibling.Scope, Generation: sibling.Generation}
	for _, m := range legacy.Messages { // the fixture's fixed dates would let the first sweep purge the history
		m.Receipt.Accepted, m.Receipt.Exp = time.Now().Add(-time.Minute), time.Now().Add(time.Hour)
	}
	// The legacy writer kept messages, idempotency and a Dots Inbox inside the row; State no longer serializes them.
	base, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	row := map[string]any{}
	if err = json.Unmarshal(base, &row); err != nil {
		t.Fatal(err)
	}
	row["Messages"], row["Idempotency"] = legacy.Messages, legacy.Idempotency
	row["Remote"] = map[string]any{"Installs": map[string]any{"install_1": map[string]any{
		"Owner": "member_000", "Agent": "agent_000_0", "Kid": "key1", "Exp": time.Now().Add(time.Hour),
		"Inbox": map[string]json.RawMessage{id: envelope}}}}
	raw, err := json.Marshal(row)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, "UPDATE relay_state SET data=$1", raw); err != nil {
		t.Fatal(err)
	}
	decoded := func(data []byte) map[string]any {
		t.Helper()
		var value map[string]any
		if err := json.Unmarshal(data, &value); err != nil {
			t.Fatal(err)
		}
		return value
	}
	var before []byte
	if err = pool.QueryRow(ctx, "SELECT data FROM relay_state").Scan(&before); err != nil {
		t.Fatal(err)
	}

	if _, err = provider.UpTo(ctx, 3); err != nil {
		t.Fatal(err)
	}
	pool.Reset() // cached statements still describe the old column

	var after []byte
	var humanModel int
	var compression string
	read := func() {
		t.Helper()
		if err = pool.QueryRow(ctx, "SELECT state, human_model, pg_column_compression(state) FROM relay_state").Scan(&after, &humanModel, &compression); err != nil {
			t.Fatal(err)
		}
	}
	read()
	if !reflect.DeepEqual(decoded(before), decoded(after)) {
		t.Fatal("migrated row reads as a different value")
	}
	if humanModel != 1 || compression != "lz4" {
		t.Fatalf("human_model %d, compression %q", humanModel, compression)
	}
	var legacyView struct{ Messages map[string]*Message }
	if err = json.Unmarshal(after, &legacyView); err != nil {
		t.Fatal(err)
	}
	if e, err := Parse(legacyView.Messages[id].Envelope); err != nil || e.Digest() != signed.Digest() || !e.Verify(public) {
		t.Fatalf("00003 envelope lost its digest or signature: %v", err)
	}

	// The queries of the previous writer fail with an SQL error instead of bypassing the guard.
	for _, old := range []string{
		"SELECT epoch, clock, data, clock_timestamp()::timestamptz AS now FROM relay_state WHERE singleton = true FOR UPDATE",
		"UPDATE relay_state SET data = '{}', epoch = epoch + 1, clock = clock_timestamp() WHERE singleton = true",
	} {
		var pgErr *pgconn.PgError
		if _, err = pool.Exec(ctx, old); !errors.As(err, &pgErr) || pgErr.Code != "42703" {
			t.Fatalf("old query %q: %v", old, err)
		}
	}

	// 00004 deletes every message copy and keeps the rest of the row (D052).
	if _, err = provider.Up(ctx); err != nil {
		t.Fatal(err)
	}
	read()
	want, got := decoded(before), decoded(after)
	delete(want, "Messages")
	delete(want, "Idempotency")
	delete(want["Remote"].(map[string]any)["Installs"].(map[string]any)["install_1"].(map[string]any), "Inbox")
	if !reflect.DeepEqual(want, got) || humanModel != 1 || compression != "lz4" {
		t.Fatalf("00004 changed more than message copies: human_model %d, compression %q", humanModel, compression)
	}
	if bytes.Contains(after, []byte(id)) || bytes.Contains(after, []byte("synthetic retained conversation")) {
		t.Fatal("00004 left a message id or conversation copy in Postgres")
	}

	// The product path reads the cleaned row, keeps members and agents, and writes no message copy back.
	s := &Service{Redis: testRedis(), Pool: pool}
	if _, err = s.transaction(ctx, func(st *State, _ time.Time) (any, error) {
		if len(st.Messages) != 0 || len(st.Idempotency) != 0 || len(st.Members) != 100 || len(st.Agents) != 200 {
			t.Fatalf("messages %d, members %d, agents %d", len(st.Messages), len(st.Members), len(st.Agents))
		}
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}
	read()
	if humanModel != 1 || compression != "lz4" || bytes.Contains(after, []byte(`"Messages"`)) {
		t.Fatalf("after write: human_model %d, compression %q", humanModel, compression)
	}
	if _, err = pool.Exec(ctx, "UPDATE relay_state SET human_model = 0"); err == nil {
		t.Fatal("human model downgrade accepted")
	}
	if _, err = pool.Exec(ctx, "UPDATE relay_state SET state = '[]'"); err == nil {
		t.Fatal("non-object state accepted")
	}
	if _, err = provider.Down(ctx); err == nil {
		t.Fatal("down migration ran")
	}
}

// The jsonb column refused NUL at save. The json column would store it, so the transaction keeps refusing it.
func TestEscapedNULIsRefusedAtSave(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	var before int64
	if err := pool.QueryRow(context.Background(), "SELECT epoch FROM relay_state").Scan(&before); err != nil {
		t.Fatal(err)
	}
	_, err := f.s.transaction(context.Background(), func(st *State, _ time.Time) (any, error) {
		st.Agents["agent_a"].Name = "nul\x00name"
		return nil, nil
	})
	if err == nil || err.Error() != "unavailable" {
		t.Fatalf("NUL state saved: %v", err)
	}
	var after int64
	if err := pool.QueryRow(context.Background(), "SELECT epoch FROM relay_state").Scan(&after); err != nil || after != before {
		t.Fatalf("epoch moved from %d to %d: %v", before, after, err)
	}
}
