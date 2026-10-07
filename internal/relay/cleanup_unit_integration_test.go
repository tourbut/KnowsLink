//go:build integration

// Valid self-cleanup admission across stale snapshots and an owner's agents' ACKs (FIX-2 review, Sol observations).
package relay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Review FIX-2 (Sol): a cleanup that another process's commit or a restart made provable is judged on current records
// before admission, not refused as new work until the 1s sweep, while invented or foreign cleanup stays new work.
func TestCleanupProvenWhileSnapshotStaleAndNewFull(t *testing.T) {
	for _, mode := range []string{"stale", "restart"} {
		t.Run(mode, func(t *testing.T) {
			pool := messagePool(t)
			f := setup(t, pool)
			s := f.s
			other := &Service{Pool: pool, SyntheticSignup: true}
			f.s, f.handler = other, other.Handler()
			f.addAgent("agent_d")
			if mode == "restart" {
				s = &Service{Pool: pool}
			}
			h := s.boundedHTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
			for range 16 {
				s.httpNew <- struct{}{}
			}
			defer func() {
				for range 16 {
					<-s.httpNew
				}
			}()
			serve := func(token, body string) (int, string) {
				r := httptest.NewRequest("POST", "/v1/key-revoke", strings.NewReader(body))
				r.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				return w.Code, w.Body.String()
			}
			if c, b := serve(f.tokens["agent_d_owner"], `{"agent":"agent_d","kid":"key1"}`); c != 200 {
				t.Fatalf("own revoke before the sweep: %d %s", c, b)
			}
			for _, c := range []struct{ token, body string }{
				{f.tokens["agent_d_owner"], `{"agent":"agent_a","kid":"key1"}`},
				{"invented", `{"agent":"agent_d","kid":"key1"}`},
			} {
				if code, b := serve(c.token, c.body); code != 429 || !strings.Contains(b, "capacity") {
					t.Fatalf("unproven cleanup took the cleanup channel: %d %s", code, b)
				}
			}
		})
	}
}

// Review FIX-2 (Sol): the owner's own control (key/agent revoke, unpair, deny, logout) never waits behind its agents' ACKs,
// so an agent repeating ACKs cannot keep its owner from revoking it. Each unit still holds one cleanup slot.
func TestOwnerControlNotHeldByOwnAgentAcks(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	f.send(f.message("0199a3f2-4c10-7a11-8b22-3344556677aa", "idempotency-key-ack"), "agent_a", "", 200)
	lease := f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]any{}, 200)["lease_token"].(string)
	ack := `{"id":"0199a3f2-4c10-7a11-8b22-3344556677aa","token":"` + lease + `"}`
	f.call("POST", "/v1/persist", f.tokens["agent_b"], json.RawMessage(ack), 200)
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	h := f.s.boundedHTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	serve := func(path, token, body string) (int, string) {
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code, w.Body.String()
	}
	revoke := func() (int, string) {
		return serve("/v1/key-revoke", f.tokens["agent_b_owner"], `{"agent":"agent_b","kid":"key1"}`)
	}
	wait := func(n int) {
		t.Helper()
		for deadline := time.Now().Add(5 * time.Second); len(f.s.httpClean) != n; time.Sleep(time.Millisecond) {
			if time.Now().After(deadline) {
				t.Fatalf("cleanup slots %d want %d", len(f.s.httpClean), n)
			}
		}
	}

	// Deterministic: while the row lock blocks admission, the agent's ACK and the owner's revoke each hold one slot.
	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	if _, err = tx.Exec(context.Background(), "SELECT 1 FROM relay_state WHERE singleton = true FOR UPDATE"); err != nil {
		t.Fatal(err)
	}
	codes := make(chan int, 2)
	go func() { c, _ := serve("/v1/ack", f.tokens["agent_b"], ack); codes <- c }()
	wait(1)
	go func() { c, _ := revoke(); codes <- c }()
	wait(2)
	if c, b := serve("/v1/ack", f.tokens["agent_b"], ack); c != 429 || !strings.Contains(b, "capacity") {
		t.Fatalf("second agent ACK of the same owner: %d %s", c, b)
	}
	if c, b := revoke(); c != 429 || !strings.Contains(b, "capacity") {
		t.Fatalf("second owner control: %d %s", c, b)
	}
	_ = tx.Rollback(context.Background())
	for range 2 {
		if c := <-codes; c != 200 {
			t.Fatalf("after the lock: %d", c)
		}
	}

	// No lock: 48 connections repeat the agent's ACK while the owner revokes 18 times.
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	var stop atomic.Bool
	var flood sync.WaitGroup
	for range 48 {
		flood.Add(1)
		go func() {
			defer flood.Done()
			for !stop.Load() {
				serve("/v1/ack", f.tokens["agent_b"], ack)
			}
		}()
	}
	time.Sleep(200 * time.Millisecond)
	refused := 0
	for range 18 {
		if c, b := revoke(); c != 200 {
			refused++
			t.Logf("owner revoke during agent ACK flood: %d %s", c, strings.TrimSpace(b))
		}
		time.Sleep(50 * time.Millisecond)
	}
	stop.Store(true)
	flood.Wait()
	t.Logf("owner revoke refused %d/18 during its agent's ACK flood", refused)
	if refused != 0 {
		t.Fatalf("agent ACK flood refused its owner's revoke %d/18", refused)
	}
	f.mutate(func(st *State) {
		if len(st.HTTP) != 0 {
			t.Fatalf("HTTP admission leaked: %d", len(st.HTTP))
		}
	})
}
