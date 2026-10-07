//go:build integration

// Cleanup admission under valid-credential floods, across processes and after failed finishes (review H-2/L-2).
package relay

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Review H2: valid credentials sending foreign, lease-less or repeated own cleanup cannot hold the cleanup channel
// while DB admission waits, so other owners' valid cleanup keeps its four reserved slots.
func TestValidCredentialCleanupFlood(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	f.addAgent("agent_d")
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	h := f.s.boundedHTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	own := func(id string) string { return `{"agent":"` + id + `","kid":"key1"}` }
	serve := func(path, token, body string) (int, string) {
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code, w.Body.String()
	}
	attacks := []struct{ name, path, token, body string }{
		{"foreign owner revoke", "/v1/key-revoke", f.tokens["agent_b_owner"], own("agent_a")},
		{"agent no-lease ack", "/v1/ack", f.tokens["agent_b"], `{"id":"none","token":"none"}`},
		{"own repeated revoke", "/v1/key-revoke", f.tokens["agent_b_owner"], own("agent_b")},
		{"new work control", "/v1/pull", f.tokens["agent_b"], "{}"},
	}
	wait := func(what string, cond func() bool) {
		t.Helper()
		for deadline := time.Now().Add(5 * time.Second); !cond(); time.Sleep(10 * time.Millisecond) {
			if time.Now().After(deadline) {
				t.Fatalf("%s: new %d clean %d", what, len(f.s.httpNew), len(f.s.httpClean))
			}
		}
	}

	// A. Deterministic: while the global row lock blocks DB admission, 4 foreign and 4 lease-less requests wait as new
	// work, the owner's 4 repeated own revokes keep only one cleanup slot, and 3 other owners still take the other 3.
	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	if _, err = tx.Exec(context.Background(), "SELECT 1 FROM relay_state WHERE singleton = true FOR UPDATE"); err != nil {
		t.Fatal(err)
	}
	codes := make(chan string, 16)
	for _, a := range attacks[:3] {
		for range 4 {
			go func() {
				c, b := serve(a.path, a.token, a.body)
				codes <- fmt.Sprint(a.name, " ", c, " ", strings.Contains(b, "capacity"))
			}()
		}
	}
	wait("attack admission", func() bool { return len(f.s.httpNew) == 8 && len(f.s.httpClean) == 1 && len(codes) == 3 })
	for range 3 {
		if got := <-codes; got != "own repeated revoke 429 true" {
			t.Fatalf("repeat beyond the owner's slot: %s", got)
		}
	}
	for _, id := range []string{"agent_a", "agent_c", "agent_d"} {
		go func() {
			c, _ := serve("/v1/key-revoke", f.tokens[id+"_owner"], own(id))
			codes <- fmt.Sprint(id, " ", c)
		}()
	}
	wait("valid cleanup behind flood", func() bool { return len(f.s.httpClean) == 4 })
	_ = tx.Rollback(context.Background())
	for range 12 {
		if got := <-codes; !strings.HasSuffix(got, " 200") && !strings.HasSuffix(got, " 200 false") {
			t.Fatalf("after the lock: %s", got)
		}
	}
	f.mutate(func(st *State) {
		b, agent := f.owners["agent_b"], "agent_b"
		if len(st.Rates["cleanup:member:"+b]) != 1 || len(st.Rates["http:member:"+b]) != 4 || len(st.Rates["http:member:"+agent]) != 4 || len(st.Rates["cleanup:member:"+agent]) != 0 {
			t.Fatalf("A rates: %v", st.Rates)
		}
		for _, id := range []string{"agent_a", "agent_c", "agent_d"} {
			if len(st.Rates["cleanup:member:"+f.owners[id]]) != 1 {
				t.Fatalf("A own cleanup %s: %v", id, st.Rates)
			}
		}
		st.Rates = map[string][]time.Time{}
	})

	// B. No lock: 48 connections flood each attack while another owner's valid cleanup probes 18 times (review B/C shape).
	for _, a := range attacks {
		f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
		var stop atomic.Bool
		var flood sync.WaitGroup
		for range 48 {
			flood.Add(1)
			go func() {
				defer flood.Done()
				for !stop.Load() {
					serve(a.path, a.token, a.body)
				}
			}()
		}
		time.Sleep(200 * time.Millisecond)
		refused, slowest := 0, time.Duration(0)
		for range 18 {
			start := time.Now()
			if c, b := serve("/v1/key-revoke", f.tokens["agent_a_owner"], own("agent_a")); c != 200 {
				refused++
				t.Logf("%s: probe %d %s", a.name, c, strings.TrimSpace(b))
			}
			slowest = max(slowest, time.Since(start))
			time.Sleep(150 * time.Millisecond)
		}
		stop.Store(true)
		flood.Wait()
		// Latency is the shared global row-lock queue (L-1), the same for the new-work control.
		t.Logf("B %s: valid cleanup refused %d/18, slowest %s", a.name, refused, slowest.Round(time.Millisecond))
		if refused != 0 {
			t.Fatalf("%s starved valid cleanup %d/18", a.name, refused)
		}
	}
	f.mutate(func(st *State) {
		if len(st.HTTP) != 0 {
			t.Fatalf("HTTP admission leaked: %d", len(st.HTTP))
		}
	})
	if len(f.s.httpNew) != 0 || len(f.s.httpClean) != 0 {
		t.Fatal("local channel leaked")
	}
	if _, busy := f.s.cleaning.Load(f.owners["agent_b"]); busy {
		t.Fatal("owner cleanup slot leaked")
	}
}

// Review L2: the local snapshot learns another process's new credential and a restarted process's state at its next
// commit; until then cleanupOwner alone never grants a cleanup slot. boundedHTTP rereads before refusing it as new
// work (TestCleanupProvenWhileSnapshotStaleAndNewFull).
func TestCleanupSnapshotAcrossProcessesAndRestart(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	other := &Service{Pool: pool, SyntheticSignup: true}
	before := f.s
	f.s, f.handler = other, other.Handler()
	f.addAgent("agent_d")
	f.s = before
	r := httptest.NewRequest("POST", "/v1/key-revoke", strings.NewReader(`{"agent":"agent_d","kid":"key1"}`))
	r.Header.Set("Authorization", "Bearer "+f.tokens["agent_d_owner"])
	r = r.WithContext(context.WithValue(r.Context(), requestBodyKey{}, []byte(`{"agent":"agent_d","kid":"key1"}`)))
	if before.cleanupOwner(r, f.tokens["agent_d_owner"]) != "" {
		t.Fatal("stale snapshot knew another process's credential")
	}
	restarted := &Service{Pool: pool}
	if restarted.cleanupOwner(r, f.tokens["agent_d_owner"]) != "" {
		t.Fatal("restart without snapshot took cleanup")
	}
	for _, s := range []*Service{before, restarted} {
		if _, err := s.transaction(context.Background(), func(*State, time.Time) (any, error) { return nil, nil }); err != nil {
			t.Fatal(err)
		}
		if s.cleanupOwner(r, f.tokens["agent_d_owner"]) != f.owners["agent_d"] {
			t.Fatal("commit did not refresh the cleanup snapshot")
		}
	}
}

// TESTER probe and H-2: a finish transaction that times out (2s) behind the row lock, as with a large relay state, leaves
// its shared admission to the process's next commit instead of the 30s expiry, so neither the owner's cleanup slot nor a
// new-work slot stays occupied.
func TestFailedFinishReclaimedAtNextCommit(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	ctx := context.Background()
	rollback := func() {}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := pool.Begin(ctx)
		if err == nil {
			_, err = tx.Exec(ctx, "SELECT 1 FROM relay_state WHERE singleton = true FOR UPDATE")
			rollback = func() { _ = tx.Rollback(ctx) }
		}
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(200)
	})
	orphans := func() (n int) { f.s.orphans.Range(func(any, any) bool { n++; return true }); return n }
	for _, c := range []struct{ path, token, body string }{
		{"/v1/key-revoke", f.tokens["agent_a_owner"], `{"agent":"agent_a","kid":"key1"}`},
		{"/v1/pull", f.tokens["agent_b"], "{}"},
	} {
		r := httptest.NewRequest("POST", c.path, strings.NewReader(c.body))
		r.Header.Set("Authorization", "Bearer "+c.token)
		w := httptest.NewRecorder()
		f.s.boundedHTTP(next).ServeHTTP(w, r) // the finish waits 2s on the held lock and fails
		rollback()
		if w.Code != 200 || orphans() != 1 {
			t.Fatalf("%s: %d, %d orphaned admissions", c.path, w.Code, orphans())
		}
		f.mutate(func(st *State) {
			if len(st.HTTP) != 0 || !st.enterHTTP("next", f.owners["agent_a"], time.Now()) {
				t.Fatalf("%s: failed finish kept %d shared admissions", c.path, len(st.HTTP))
			}
			delete(st.HTTP, "next")
		})
		if orphans() != 0 {
			t.Fatal("reclaimed token kept")
		}
	}
}
