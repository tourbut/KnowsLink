//go:build integration

// Isolated Postgres tests prove atomic replay, current authorization, transport races, and owner gate security.
package relay

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type fixture struct {
	s       *Service
	handler http.Handler
	tokens  map[string]string
	owners  map[string]string
	private map[string]ed25519.PrivateKey
	t       *testing.T
}

func setup(t *testing.T, pool *pgxpool.Pool) *fixture {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `UPDATE relay_state SET data='{}',epoch=epoch+1,clock=clock_timestamp()`); err != nil {
		t.Fatal(err)
	}
	f := &fixture{&Service{pool}, nil, map[string]string{}, map[string]string{}, map[string]ed25519.PrivateKey{}, t}
	f.handler = f.s.Handler()
	for _, id := range []string{"agent_a", "agent_b", "agent_c"} {
		owner := f.call("POST", "/v1/owners", "", map[string]any{}, 200)
		f.tokens[id+"_owner"] = owner["credential"].(string)
		f.owners[id] = owner["owner"].(string)
		public, private, _ := ed25519.GenerateKey(rand.Reader)
		f.private[id] = private
		encoded := base64.RawURLEncoding.EncodeToString(public)
		proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, popBytes(f.owners[id], id, "key1", encoded)))
		agent := f.call("POST", "/v1/agents", f.tokens[id+"_owner"], map[string]any{"agent": id, "kid": "key1", "public": encoded, "proof": proof}, 200)
		f.tokens[id] = agent["credential"].(string)
	}
	f.call("POST", "/v1/invites", f.tokens["agent_a"], map[string]any{"agent": "agent_a", "target": "agent_b"}, 200)
	f.call("POST", "/v1/invite-decision", f.tokens["agent_b"], map[string]any{"agent": "agent_a", "target": "agent_b", "decision": "accept"}, 401)
	f.call("POST", "/v1/invite-decision", f.tokens["agent_b_owner"], map[string]any{"agent": "agent_a", "target": "agent_b", "decision": "accept"}, 200)
	return f
}
func (f *fixture) request(method, path, token string, body []byte, claim string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(string(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Execution-Claim", claim)
	w := httptest.NewRecorder()
	f.handler.ServeHTTP(w, req)
	return w
}
func (f *fixture) call(method, path, token string, body any, status int) map[string]any {
	f.t.Helper()
	raw, _ := json.Marshal(body)
	w := f.request(method, path, token, raw, "")
	if w.Code != status {
		f.t.Fatalf("%s %s got %d %s want %d", method, path, w.Code, w.Body, status)
	}
	result := map[string]any{}
	_ = json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func (f *fixture) message(id, key string) []byte {
	return wire(f.t, f.private["agent_a"], id, "agent_a", "agent_b", "schedule.query", key, "", queryBody(), time.Now().Add(180*time.Second))
}
func (f *fixture) send(raw []byte, agent, claim string, status int) map[string]any {
	f.t.Helper()
	w := f.request("POST", "/v1/send", f.tokens[agent], raw, claim)
	if w.Code != status {
		f.t.Fatalf("send got %d %s want %d", w.Code, w.Body, status)
	}
	r := map[string]any{}
	_ = json.Unmarshal(w.Body.Bytes(), &r)
	return r
}
func (f *fixture) deliver(id string) string {
	f.t.Helper()
	lease := f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]any{}, 200)
	token := lease["lease_token"].(string)
	payload := map[string]any{"id": id, "token": token}
	f.call("POST", "/v1/ack", f.tokens["agent_b"], payload, 409)
	f.call("POST", "/v1/persist", f.tokens["agent_b"], payload, 200)
	f.call("POST", "/v1/ack", f.tokens["agent_b"], payload, 200)
	claim := f.call("POST", "/v1/claim", f.tokens["agent_b"], map[string]any{"id": id}, 200)
	return claim["claim"].(string)
}
func (f *fixture) mutate(fn func(*State)) {
	f.t.Helper()
	_, err := f.s.transaction(context.Background(), func(st *State, _ time.Time) (any, error) { fn(st); return nil, nil })
	if err != nil {
		f.t.Fatal(err)
	}
}
func TestPostgresSafety(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is required; run make verify-mvp")
	}
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" {
		t.Fatal("isolated database marker required")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	t.Run("auth_first_atomic_replay_and_rollback", func(t *testing.T) {
		f := setup(t, pool)
		raw := f.message(firstID, "idempotency-key-01")
		f.send(raw, "agent_c", "", 403)
		var wg sync.WaitGroup
		codes := make(chan int, 12)
		for i := 0; i < 12; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); codes <- f.request("POST", "/v1/send", f.tokens["agent_a"], raw, "").Code }()
		}
		wg.Wait()
		close(codes)
		for c := range codes {
			if c != 200 {
				t.Fatalf("concurrent replay %d", c)
			}
		}
		body := queryBody()
		body["granularity_min"] = 31
		conflict := wire(t, f.private["agent_a"], secondID, "agent_a", "agent_b", "schedule.query", "idempotency-key-01", "", body, time.Now().Add(time.Minute))
		f.send(conflict, "agent_a", "", 409)
		expired := wire(t, f.private["agent_a"], secondID, "agent_a", "agent_b", "schedule.query", "idempotency-key-01", "", queryBody(), time.Now().Add(-time.Minute))
		r := f.send(expired, "agent_a", "", 200)
		if r["id"] != firstID || r["envelope"] != nil {
			t.Fatal("replay not receipt-only")
		}
		long := wire(t, f.private["agent_a"], secondID, "agent_a", "agent_b", "schedule.query", "idempotency-key-02", "", queryBody(), time.Now().Add(301*time.Second))
		f.send(long, "agent_a", "", 422)
		f.send(f.message(secondID, "idempotency-key-02"), "agent_a", "", 200)
		f.send(f.message(firstID, "idempotency-key-03"), "agent_a", "", 409)
		f.send(f.message(gateID, "idempotency-key-03"), "agent_a", "", 200)
		f.call("POST", "/v1/key-revoke", f.tokens["agent_a_owner"], map[string]any{"agent": "agent_a", "kid": "key1"}, 200)
		f.send(raw, "agent_a", "", 401)
	})
	t.Run("last_lease_ack_and_shared_claim_restart", func(t *testing.T) {
		f := setup(t, pool)
		f.send(f.message(firstID, "idempotency-key-01"), "agent_a", "", 200)
		old := ""
		for attempt := 1; attempt <= 3; attempt++ {
			lease := f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]any{}, 200)
			token := lease["lease_token"].(string)
			if int(lease["attempts"].(float64)) != attempt {
				t.Fatal("wrong attempt count")
			}
			if old != "" {
				f.call("POST", "/v1/ack", f.tokens["agent_b"], map[string]any{"id": firstID, "token": old}, 409)
			}
			if attempt < 3 {
				f.mutate(func(st *State) { st.Messages[firstID].LeaseUntil = time.Now().Add(-time.Second) })
				old = token
				continue
			}
			f.call("POST", "/v1/persist", f.tokens["agent_b"], map[string]any{"id": firstID, "token": token}, 200)
			f.call("POST", "/v1/ack", f.tokens["agent_b"], map[string]any{"id": firstID, "token": token}, 200)
		}
		codes := make(chan int, 8)
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				body, _ := json.Marshal(map[string]string{"id": firstID})
				codes <- f.request("POST", "/v1/claim", f.tokens["agent_b"], body, "").Code
			}()
		}
		wg.Wait()
		close(codes)
		wins := 0
		for code := range codes {
			if code == 200 {
				wins++
			} else if code != 409 {
				t.Fatalf("claim failed %d", code)
			}
		}
		if wins != 1 {
			t.Fatalf("claim winners %d", wins)
		}
		restarted, err := pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			t.Fatal(err)
		}
		defer restarted.Close()
		f.s = &Service{restarted}
		f.handler = f.s.Handler()
		f.call("POST", "/v1/claim", f.tokens["agent_b"], map[string]any{"id": firstID}, 409)
	})
	t.Run("owner_gate_csrf_binding_and_no_effect", func(t *testing.T) {
		f := setup(t, pool)
		receipt := f.send(f.message(firstID, "idempotency-key-01"), "agent_a", "", 200)
		claim := f.deliver(firstID)
		h := wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "approval-key-0001", firstID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute))
		f.send(h, "agent_b", "", 403)
		f.send(h, "agent_b", claim, 200)
		duplicate := wire(t, f.private["agent_b"], secondID, "agent_b", "agent_b", "relay.approval.request", "approval-key-0002", firstID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute))
		f.send(duplicate, "agent_b", claim, 409)
		w := f.request("GET", "/owner/gates/"+gateID, f.tokens["agent_b_owner"], nil, "")
		if w.Code != 200 || !strings.Contains(w.Body.String(), "granularity_min") || !strings.Contains(w.Body.String(), "pending") {
			t.Fatal("verified UI missing")
		}
		f.call("GET", "/owner/gates/"+gateID, f.tokens["agent_b"], nil, 401)
		post := func(token, csrfValue, decision string) int {
			req := httptest.NewRequest("POST", "/owner/gates/"+gateID, strings.NewReader(url.Values{"csrf": {csrfValue}, "decision": {decision}}.Encode()))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			f.handler.ServeHTTP(w, req)
			return w.Code
		}
		owner := f.tokens["agent_b_owner"]
		if post(owner, "bad", "approve") != 403 || post(owner, csrf(owner, gateID), "approve") != 303 || post(owner, csrf(owner, gateID), "deny") != 409 {
			t.Fatal("CSRF or atomic decision failed")
		}
		f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": "wrong"}, 409)
		result := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 200)
		if result["executable"] != false || result["disclosure"] != false {
			t.Fatal("gate enabled effect")
		}
		f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
		f.call("POST", "/v1/authorize", f.tokens["agent_b"], map[string]any{"id": firstID, "claim": claim}, 200)
		done := wire(t, f.private["agent_b"], secondID, "agent_b", "agent_a", "relay.result", "result-key-00001", firstID, map[string]any{"status": "done"}, time.Now().Add(time.Minute))
		f.send(done, "agent_b", claim, 403)
		denied := wire(t, f.private["agent_b"], secondID, "agent_b", "agent_a", "relay.result", "result-key-00001", firstID, map[string]any{"status": "denied"}, time.Now().Add(time.Minute))
		f.send(denied, "agent_b", claim, 200)
		f.mutate(func(st *State) {
			if len(st.Messages[firstID].Envelope) != 0 || len(st.Messages[firstID].Inbox) != 0 {
				t.Fatal("completed original retained")
			}
		})
	})
	t.Run("unpair_generation_expiry_max_attempts_and_clock", func(t *testing.T) {
		f := setup(t, pool)
		raw := f.message(firstID, "idempotency-key-01")
		f.send(raw, "agent_a", "", 200)
		f.call("POST", "/v1/unpair", f.tokens["agent_b_owner"], map[string]any{"agent": "agent_a", "target": "agent_b"}, 200)
		f.call("POST", "/v1/invites", f.tokens["agent_a"], map[string]any{"agent": "agent_a", "target": "agent_b"}, 200)
		f.call("POST", "/v1/invite-decision", f.tokens["agent_b_owner"], map[string]any{"agent": "agent_a", "target": "agent_b", "decision": "accept"}, 200)
		f.send(raw, "agent_a", "", 403)
		f.send(f.message(secondID, "idempotency-key-02"), "agent_a", "", 200)
		f.mutate(func(st *State) { st.Messages[secondID].Receipt.Exp = time.Now().Add(-time.Second) })
		r := f.call("GET", "/v1/receipts/"+secondID, f.tokens["agent_a"], nil, 200)
		if r["receipt"].(map[string]any)["transport"] != "failed:expired" {
			t.Fatal("expiry not transport failure")
		}
		f.send(f.message(gateID, "idempotency-key-03"), "agent_a", "", 200)
		for i := 0; i < 3; i++ {
			f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]any{}, 200)
			f.mutate(func(st *State) { st.Messages[gateID].LeaseUntil = time.Now().Add(-time.Second) })
		}
		r = f.call("GET", "/v1/receipts/"+gateID, f.tokens["agent_a"], nil, 200)
		if r["receipt"].(map[string]any)["transport"] != "failed:max_attempts" {
			t.Fatal("attempt exhaustion wrong")
		}
		_, err := pool.Exec(context.Background(), `UPDATE relay_state SET clock=clock_timestamp()+interval '1 hour'`)
		if err != nil {
			t.Fatal(err)
		}
		f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]any{}, 503)
		_, _ = pool.Exec(context.Background(), `UPDATE relay_state SET clock=clock_timestamp()`)
	})
	t.Run("rotation_pop_and_key_reassignment", func(t *testing.T) {
		f := setup(t, pool)
		public, private, _ := ed25519.GenerateKey(rand.Reader)
		encoded := base64.RawURLEncoding.EncodeToString(public)
		body := map[string]any{"agent": "agent_a", "kid": "key2", "public": encoded, "proof": base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, popBytes(f.owners["agent_b"], "agent_a", "key2", encoded)))}
		f.call("POST", "/v1/keys", f.tokens["agent_a_owner"], body, 422)
		body["proof"] = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, popBytes(f.owners["agent_a"], "agent_a", "key2", encoded)))
		f.call("POST", "/v1/keys", f.tokens["agent_a_owner"], body, 200)
		f.call("POST", "/v1/keys", f.tokens["agent_a_owner"], body, 409)
		f.send(f.message(firstID, "idempotency-key-01"), "agent_a", "", 401)
	})
	if !t.Failed() {
		fmt.Println("PASS: isolated Postgres safety invariants")
	}
}

func TestGateFailureStates(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is required; run make verify-mvp")
	}
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" {
		t.Fatal("isolated database marker required")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, state := range []string{"denied", "expired", "revoked", "unavailable"} {
		t.Run(state, func(t *testing.T) {
			f := setup(t, pool)
			receipt := f.send(f.message(firstID, "idempotency-key-01"), "agent_a", "", 200)
			claim := f.deliver(firstID)
			h := wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "approval-key-0001", firstID, map[string]any{"reason": "permission_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute))
			f.send(h, "agent_b", claim, 200)
			f.mutate(func(st *State) {
				switch state {
				case "denied":
					st.Gates[gateID].State = "denied"
				case "expired":
					st.Gates[gateID].Exp = time.Now().Add(-time.Second)
				case "revoked":
					st.Pairs[pairID("agent_a", "agent_b")].State = "revoked"
				case "unavailable":
					st.Messages[firstID].Envelope = nil
				}
			})
			w := f.request("GET", "/owner/gates/"+gateID, f.tokens["agent_b_owner"], nil, "")
			if w.Code != 200 || !strings.Contains(w.Body.String(), "상태: "+state) || strings.Contains(w.Body.String(), `value="approve"`) {
				t.Fatalf("unsafe %s UI: %d %s", state, w.Code, w.Body)
			}
			f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
		})
	}
}
func TestApprovalAndResultInstanceBinding(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is required; run make verify-mvp")
	}
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" {
		t.Fatal("isolated database marker required")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := setup(t, pool)
	receipt := f.send(f.message(firstID, "idempotency-key-01"), "agent_a", "", 200)
	claim := f.deliver(firstID)
	f.send(f.message(secondID, "idempotency-key-02"), "agent_a", "", 200)
	secondClaim := f.deliver(secondID)
	h := wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "approval-key-0001", firstID, map[string]any{"reason": "permission_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute))
	f.send(h, "agent_b", claim, 200)
	f.mutate(func(st *State) { st.Gates[gateID].State = "approved" })
	f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": secondClaim}, 409)
	wrong := wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677bb", "agent_b", "agent_b", "relay.result", "result-key-00001", firstID, map[string]any{"status": "denied"}, time.Now().Add(time.Minute))
	f.send(wrong, "agent_b", claim, 403)
	valid := wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677bb", "agent_b", "agent_a", "relay.result", "result-key-00001", firstID, map[string]any{"status": "denied"}, time.Now().Add(time.Minute))
	bad := strings.Replace(string(valid), `"status":"denied"`, `"status":"denied","result":{"title":"private"}`, 1)
	f.send([]byte(bad), "agent_b", claim, 422)
	f.call("POST", "/v1/authorize", f.tokens["agent_b"], map[string]any{"id": secondID, "claim": secondClaim}, 200)
	f.send(valid, "agent_b", claim, 200)
	public, private, _ := ed25519.GenerateKey(rand.Reader)
	encoded := base64.RawURLEncoding.EncodeToString(public)
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, popBytes(f.owners["agent_a"], "agent_a", "key2", encoded)))
	f.call("POST", "/v1/keys", f.tokens["agent_a_owner"], map[string]any{"agent": "agent_a", "kid": "key2", "public": encoded, "proof": proof}, 200)
	f.call("GET", "/v1/receipts/0199a3f2-4c10-7a11-8b22-3344556677bb", f.tokens["agent_b"], nil, 403)
	f.call("POST", "/v1/owner-revoke", f.tokens["agent_b_owner"], map[string]any{}, 200)
	f.call("POST", "/v1/authorize", f.tokens["agent_b"], map[string]any{"id": secondID, "claim": secondClaim}, 401)
}

func TestPendingAndConcurrentAccept(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is required; run make verify-mvp")
	}
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" {
		t.Fatal("isolated database marker required")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := setup(t, pool)
	f.call("POST", "/v1/invites", f.tokens["agent_a"], map[string]any{"agent": "agent_a", "target": "agent_c"}, 200)
	f.mutate(func(st *State) {
		active := 0
		for _, p := range st.Pairs {
			if p.State == "active" {
				active++
			}
		}
		if active != 1 {
			t.Fatal("pending invite consumed active relation")
		}
	})
	body, _ := json.Marshal(map[string]string{"agent": "agent_a", "target": "agent_c", "decision": "accept"})
	codes := make(chan int, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes <- f.request("POST", "/v1/invite-decision", f.tokens["agent_c_owner"], body, "").Code
		}()
	}
	wg.Wait()
	close(codes)
	for code := range codes {
		if code != 200 {
			t.Fatalf("concurrent accept %d", code)
		}
	}
	w := f.request("GET", "/v1/contacts", f.tokens["agent_c"], nil, "")
	var contacts []string
	_ = json.Unmarshal(w.Body.Bytes(), &contacts)
	if len(contacts) != 1 || contacts[0] != "agent_a" {
		t.Fatal("duplicate active relation")
	}
}
