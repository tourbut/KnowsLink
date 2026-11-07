//go:build integration

// A stalled session holding the row lock fails new work within the limit, ends by itself, and writes then resume.
package relay

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestStalledLockHolderFailsThenRecovers(t *testing.T) {
	ctx := context.Background()
	cfg, err := PoolConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	f := setup(t, pool)

	// The holder uses the relay's own settings: it stands for a relay process that stalled inside a transaction.
	holder, err := pgx.ConnectConfig(ctx, cfg.ConnConfig.Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = holder.Close(ctx) }()
	tx, err := holder.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.Exec(ctx, "SELECT 1 FROM relay_state WHERE singleton = true FOR UPDATE"); err != nil {
		t.Fatal(err)
	}
	stalled := time.Now()

	w := f.request("POST", "/v1/pull", f.tokens["agent_b"], []byte("{}"), "")
	waited := time.Since(stalled)
	if w.Code != 503 || !strings.Contains(w.Body.String(), "unavailable") {
		t.Fatalf("request behind stalled lock: %d %s", w.Code, w.Body)
	}
	if waited < lockTimeout || waited > lockTimeout+time.Second {
		t.Fatalf("request failed after %v, limit %v", waited, lockTimeout)
	}

	for {
		if _, err = f.s.transaction(ctx, func(*State, time.Time) (any, error) { return nil, nil }); err == nil {
			break
		}
		if time.Since(stalled) > idleTransactionTimeout+lockTimeout+time.Second {
			t.Fatalf("writes did not resume after %v: %v", time.Since(stalled), err)
		}
	}
	if resumed := time.Since(stalled); resumed < idleTransactionTimeout-time.Second {
		t.Fatalf("writes resumed after %v, before the holder could be ended", resumed)
	}
	_, err = tx.Exec(ctx, "SELECT 1")
	var pgErr *pgconn.PgError
	if err == nil || !holder.IsClosed() && !(errors.As(err, &pgErr) && pgErr.Code == "25P03") {
		t.Fatalf("holder session survived: %v", err)
	}
	f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]any{}, 200)
}
