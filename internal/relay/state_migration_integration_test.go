//go:build integration

// Migration 00003 keeps an existing jsonb row readable as the same State and makes older writers fail.
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
	raw, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, "UPDATE relay_state SET data=$1", raw); err != nil {
		t.Fatal(err)
	}
	normalized := func(data []byte) []byte {
		t.Helper()
		st := newState()
		if err := json.Unmarshal(data, st); err != nil {
			t.Fatal(err)
		}
		for _, m := range st.Messages {
			if len(m.Envelope) > 0 {
				var value any
				if err := json.Unmarshal(m.Envelope, &value); err != nil {
					t.Fatal(err)
				}
				if m.Envelope, err = Canonical(value); err != nil {
					t.Fatal(err)
				}
			}
		}
		out, err := json.Marshal(st)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	var before []byte
	if err = pool.QueryRow(ctx, "SELECT data FROM relay_state").Scan(&before); err != nil {
		t.Fatal(err)
	}

	if _, err = provider.Up(ctx); err != nil {
		t.Fatal(err)
	}
	pool.Reset() // cached statements still describe the old column

	var after []byte
	var humanModel int
	var compression string
	if err = pool.QueryRow(ctx, "SELECT state, human_model, pg_column_compression(state) FROM relay_state").Scan(&after, &humanModel, &compression); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(normalized(before), normalized(after)) {
		t.Fatal("migrated row reads as a different State")
	}
	if humanModel != 1 || compression != "lz4" {
		t.Fatalf("human_model %d, compression %q", humanModel, compression)
	}
	verify := func(stage string, data []byte) {
		t.Helper()
		st := newState()
		if err := json.Unmarshal(data, st); err != nil {
			t.Fatal(stage, err)
		}
		e, err := Parse(st.Messages[id].Envelope)
		if err != nil || e.Digest() != signed.Digest() || !e.Verify(public) {
			t.Fatalf("%s: envelope lost its digest or signature: %v", stage, err)
		}
	}
	verify("migrated", after)

	// A write and a read through the product path keep the signed envelope and the compression.
	s := &Service{Pool: pool}
	if _, err = s.transaction(ctx, func(st *State, _ time.Time) (any, error) {
		if len(st.Messages) != 1001 {
			t.Fatalf("messages %d", len(st.Messages))
		}
		return nil, nil
	}); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, "SELECT state, human_model, pg_column_compression(state) FROM relay_state").Scan(&after, &humanModel, &compression); err != nil {
		t.Fatal(err)
	}
	verify("rewritten", after)
	if humanModel != 1 || compression != "lz4" {
		t.Fatalf("after write: human_model %d, compression %q", humanModel, compression)
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
