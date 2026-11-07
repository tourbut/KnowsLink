//go:build integration

// Readiness stays answerable while every allowed transaction waits on a held row lock.
package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tourbut/KnowsLink/internal/relay"
)

func TestHealthzAnswersDuringLockWait(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || url == "" {
		t.Fatal("isolated database required")
	}
	ctx := context.Background()
	cfg, err := relay.PoolConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	holder, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = holder.Close(ctx) }()
	tx, err := holder.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err = tx.Exec(ctx, "SELECT 1 FROM relay_state WHERE singleton = true FOR UPDATE"); err != nil {
		t.Fatal(err)
	}

	api := (&relay.Service{Pool: pool}).Handler()
	var requests sync.WaitGroup
	for range 12 {
		requests.Add(1)
		go func() {
			defer requests.Done()
			r := httptest.NewRequest(http.MethodPost, "/v1/pull", strings.NewReader("{}"))
			r.Header.Set("Authorization", "Bearer none")
			w := httptest.NewRecorder()
			api.ServeHTTP(w, r)
			if w.Code != http.StatusServiceUnavailable && w.Code != http.StatusTooManyRequests {
				t.Errorf("request behind held lock: %d %s", w.Code, w.Body)
			}
		}()
	}
	observer, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = observer.Close(ctx) }()
	var waiting int
	for deadline := time.Now().Add(time.Second); waiting < 3 && time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		if err = observer.QueryRow(ctx, "SELECT count(*) FROM pg_stat_activity WHERE wait_event_type = 'Lock' AND datname = current_database()").Scan(&waiting); err != nil {
			t.Fatal(err)
		}
	}
	if waiting != 3 {
		t.Fatalf("%d transactions in lock wait, want 3", waiting)
	}

	health := healthHandler(pool.Ping)
	for range 5 {
		start := time.Now()
		w := httptest.NewRecorder()
		health.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if w.Code != http.StatusOK || time.Since(start) > 2*time.Second {
			t.Fatalf("healthz during lock wait: %d after %v", w.Code, time.Since(start))
		}
		time.Sleep(100 * time.Millisecond)
	}
	requests.Wait()
}
