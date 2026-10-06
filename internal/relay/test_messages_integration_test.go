//go:build integration

// Real Postgres machine-path checks prove no owner/business escalation and preserve trial bytes on rejected operations.
package relay

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTrialHTTP(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := setup(t, pool)
	f.call("GET", "/v1/test/registry", f.tokens["agent_a"], nil, 403)
	f.s.TestAgents = map[string]bool{"agent_a": true, "agent_b": true}
	f.call("GET", "/v1/test/registry", "", nil, 401)
	f.call("GET", "/v1/test/registry", f.tokens["agent_a_owner"], nil, 401)
	f.call("GET", "/v1/test/registry", f.tokens["agent_c"], nil, 403)
	f.call("GET", "/v1/test/registry", f.tokens["agent_a"], nil, 200)
	for _, path := range []string{"owners", "agents", "keys", "invites", "invite-decision", "unpair", "authorize", "gate-consume"} {
		f.call("POST", "/v1/test/"+path, f.tokens["agent_a"], map[string]any{}, 403)
	}
	if w := f.request("POST", "/v1/test/send", f.tokens["agent_a"], f.message(firstID, "trial-business-key"), ""); w.Code != 403 {
		t.Fatal("machine business send", w.Code)
	}
	f.send(f.message(secondID, "business-outside-key"), "agent_a", "", 200)
	f.call("GET", "/v1/test/receipts/"+secondID, f.tokens["agent_a"], nil, 403)
	f.call("POST", "/v1/test/claim", f.tokens["agent_b"], map[string]any{"id": secondID}, 403)
	raw := wire(t, f.private["agent_a"], firstID, "agent_a", "agent_b", "relay.test.message", "trial-http-key-01", "", map[string]any{"text": "ping"}, time.Now().Add(time.Minute))
	if w := f.request("POST", "/v1/test/send", f.tokens["agent_a"], raw, ""); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	// Rejected schema/authority requests must not clear the allowlisted queued trial (transaction rollback).
	f.call("POST", "/v1/test/claim", f.tokens["agent_b"], map[string]any{"id": secondID}, 403)
	lease := f.call("POST", "/v1/test/pull", f.tokens["agent_b"], map[string]any{}, 200)
	envelope := lease["envelope"].(map[string]any)
	if envelope["id"] != firstID {
		t.Fatal("machine pulled business intent")
	}
	delivery := map[string]any{"id": firstID, "token": lease["lease_token"]}
	f.call("POST", "/v1/test/persist", f.tokens["agent_a"], delivery, 409)
	f.call("POST", "/v1/test/ack", f.tokens["agent_b"], delivery, 409)
	f.call("POST", "/v1/test/persist", f.tokens["agent_b"], delivery, 200)
	f.call("POST", "/v1/test/ack", f.tokens["agent_b"], delivery, 200)
	f.call("POST", "/v1/test/claim", f.tokens["agent_b"], map[string]any{"id": firstID}, 200)
	f.call("POST", "/v1/test/claim", f.tokens["agent_b"], map[string]any{"id": firstID}, 403)
	f.call("GET", "/v1/test/receipts/"+firstID, f.tokens["agent_c"], nil, 403)
}

// A second relay with another allowlist on the same database sweeps foreign trial leases as revoked.
// make verify-mvp therefore stops the Compose relay (allowlist trial_codex,trial_grok) during Go integration tests.
func TestTrialForeignAllowlistRevokesLease(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := setup(t, pool)
	f.s.TestAgents = map[string]bool{"agent_a": true, "agent_b": true}
	raw := wire(t, f.private["agent_a"], firstID, "agent_a", "agent_b", "relay.test.message", "trial-foreign-key-01", "", map[string]any{"text": "ping"}, time.Now().Add(time.Minute))
	if w := f.request("POST", "/v1/test/send", f.tokens["agent_a"], raw, ""); w.Code != 200 {
		t.Fatal(w.Code, w.Body)
	}
	lease := f.call("POST", "/v1/test/pull", f.tokens["agent_b"], map[string]any{}, 200)
	// Same no-op transaction as Service.Cleanup in a relay started with the verify-mvp allowlist.
	foreign := &Service{Pool: pool, TestAgents: map[string]bool{"trial_codex": true, "trial_grok": true}}
	if _, err := foreign.transaction(context.Background(), func(*State, time.Time) (any, error) { return nil, nil }); err != nil {
		t.Fatal(err)
	}
	result := f.call("POST", "/v1/test/persist", f.tokens["agent_b"], map[string]any{"id": firstID, "token": lease["lease_token"]}, 409)
	if result["error"] != "invalid_lease" {
		t.Fatal("foreign sweep", result)
	}
}
