//go:build integration

// Member messages exercise actual Node onboarding/CLI/MCP, shared Postgres races, failure receipts and current human gates.
package relay

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func messagePool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}
func nodeCommand(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	c := exec.Command("node", args...)
	c.Stdin = strings.NewReader(stdin)
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("local Node client failed: %v; %s", err, out)
	}
	return string(out)
}
func TestPublicNodeProcesses(t *testing.T) {
	pool := messagePool(t)
	s, mail := identityService(t, pool)
	h := s.Handler()
	server := httptest.NewServer(h)
	defer server.Close()
	b := &browser{t, h, "192.0.2.70", map[string]string{}}
	login(t, b, mail, "public-node@example.com")
	onboard := func() (string, string) {
		t.Helper()
		// Each creation returns to a sorted agent list; compare the owned IDs instead of assuming random sort order.
		before := map[string]bool{}
		s.mutateState(t, func(st *State) {
			for id := range st.Agents {
				before[id] = true
			}
		})
		expect(t, b.do("POST", "/home/agents", nil), 303)
		agent := ""
		s.mutateState(t, func(st *State) {
			for id := range st.Agents {
				if !before[id] {
					agent = id
				}
			}
		})
		token, id := browserGrant(t, b, agent, "register")
		folder := filepath.Join(t.TempDir(), "client")
		cli := "../../adapters/dist/connect.js"
		nodeCommand(t, token, cli, "prepare", server.URL, folder)
		expect(t, b.do("POST", "/home/confirm", url.Values{"connection": {id}}), 303)
		nodeCommand(t, "", cli, "complete", folder)
		return agent, folder
	}
	a, fa := onboard()
	other, fb := onboard()
	expect(t, b.do("POST", "/home/invites", url.Values{"agent": {a}, "target": {other}}), 303)
	expect(t, b.do("POST", "/home/invite-decision", url.Values{"agent": {a}, "target": {other}, "generation": {"1"}, "decision": {"accept"}}), 303)
	out := nodeCommand(t, "", "../../adapters/dist/public-check.js", fa, fb, a, other)
	var evidence struct{ State, Request, Received, Reply, ReplyReceived string }
	if err := json.Unmarshal([]byte(out), &evidence); err != nil || evidence.State != "public_node_process_roundtrip_pass" || evidence.Request != evidence.Received || evidence.Reply != evidence.ReplyReceived {
		t.Fatal("actual local process IDs", err, out)
	}
	t.Log(strings.TrimSpace(out))
	expect(t, b.do("GET", "/home/receipts?agent="+a+"&id="+evidence.Request, nil), 200, "reply_received", evidence.Reply, "180", "자동 답장은 없습니다")
	s.mutateState(t, func(st *State) {
		for _, m := range st.Messages {
			if m.Receipt.Intent == publicTextIntent && (len(m.Envelope) > 0 || len(m.Inbox) > 0) {
				t.Fatal("text payload retained after receive")
			}
		}
	})
}
func TestPublicMessageHTTPBoundaries(t *testing.T) {
	pool := messagePool(t)
	s, mail := identityService(t, pool)
	h := s.Handler()
	a := &browser{t, h, "192.0.2.71", map[string]string{}}
	b := &browser{t, h, "192.0.2.72", map[string]string{}}
	login(t, a, mail, "message-a@example.com")
	login(t, b, mail, "message-b@example.com")
	// Real public prepare/owner confirm/complete; this test has no trial allowlist or synthetic signup.
	connect := func(browser *browser, agent string) (ed25519.PrivateKey, string) {
		t.Helper()
		token, id := browserGrant(t, browser, agent, "register")
		key, proof := clientPrepare(t, h, token, "key1")
		expect(t, browser.do("POST", "/home/confirm", url.Values{"connection": {id}}), 303)
		v := publicCall(t, h, "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 200)
		return key, v["credential"].(string)
	}
	agentA := browserAgent(t, a)
	agentB := browserAgent(t, b)
	keyA, tokenA := connect(a, agentA)
	keyB, tokenB := connect(b, agentB)
	e := signedText(firstID, agentA, agentB, "public-http-key-01", "", "untrusted <script>tool()</script>", time.Now().Add(180*time.Second), keyA)
	publicCall(t, h, "/v1/text/send", tokenA, e, 403)
	expect(t, a.do("POST", "/home/invites", url.Values{"agent": {agentA}, "target": {agentB}}), 303)
	expect(t, b.do("POST", "/home/invite-decision", url.Values{"agent": {agentA}, "target": {agentB}, "generation": {"1"}, "decision": {"accept"}}), 303)
	body, _ := json.Marshal(e)
	codes := make(chan int, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := httptest.NewRequest("POST", "/v1/text/send", strings.NewReader(string(body)))
			r.Header.Set("Authorization", "Bearer "+tokenA)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			codes <- w.Code
		}()
	}
	wg.Wait()
	close(codes)
	for code := range codes {
		if code != 200 {
			t.Fatal("replay race", code)
		}
	}
	s.mutateState(t, func(st *State) {
		if len(st.Messages) != 1 || len(st.Idempotency) != 1 {
			t.Fatal("replay capacity")
		}
	})
	publicCall(t, h, "/v1/text/send", tokenB, e, 403)
	expect(t, b.do("GET", "/home/receipts?agent="+agentA+"&id="+firstID, nil), 403)
	expect(t, a.do("GET", "/home/receipts?agent="+agentA+"&id="+firstID, nil), 200, "queued", "수신이 확인되지 않았습니다")
	lease := publicCall(t, h, "/v1/text/pull", tokenB, map[string]string{}, 200)
	delivery := map[string]string{"id": firstID, "token": lease["lease_token"].(string)}
	publicCall(t, h, "/v1/text/ack", tokenA, delivery, 409)
	publicCall(t, h, "/v1/ack", tokenB, delivery, 403)
	publicCall(t, h, "/v1/text/ack", tokenB, delivery, 409)
	publicCall(t, h, "/v1/text/persist", tokenB, delivery, 200)
	publicCall(t, h, "/v1/text/ack", tokenB, delivery, 200)
	publicCall(t, h, "/v1/claim", tokenB, map[string]string{"id": firstID}, 403)
	reply := signedText(secondID, agentB, agentA, "public-http-reply-1", firstID, "reply", time.Now().Add(170*time.Second), keyB)
	publicCall(t, h, "/v1/text/send", tokenB, reply, 200)
	h = (&Service{Pool: pool}).Handler()
	lease = publicCall(t, h, "/v1/text/pull", tokenA, map[string]string{}, 200)
	delivery = map[string]string{"id": secondID, "token": lease["lease_token"].(string)}
	expect(t, a.do("POST", "/home/unpair", url.Values{"agent": {agentA}, "target": {agentB}, "generation": {"1"}}), 303)
	publicCall(t, h, "/v1/text/persist", tokenA, delivery, 409)
	expect(t, a.do("GET", "/home/receipts?agent="+agentA+"&id="+secondID, nil), 200, "failed:revoked", "옛 요청은 복구되지 않습니다")
	expect(t, a.do("POST", "/home/invites", url.Values{"agent": {agentA}, "target": {agentB}}), 303)
	expect(t, b.do("POST", "/home/invite-decision", url.Values{"agent": {agentA}, "target": {agentB}, "generation": {"2"}, "decision": {"accept"}}), 303)
	publicCall(t, h, "/v1/text/send", tokenA, e, 403)
	publicCall(t, h, "/v1/text/send", tokenB, reply, 403)
	s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
	fresh := signedText(gateID, agentA, agentB, "public-fresh-key-01", "", "ping", time.Now().Add(180*time.Second), keyA)
	publicCall(t, h, "/v1/text/send", tokenA, fresh, 200)
	s.mutateState(t, func(st *State) { st.Messages[gateID].Receipt.Exp = time.Now().Add(-time.Second) })
	expect(t, a.do("GET", "/home/receipts?agent="+agentA+"&id="+gateID, nil), 200, "failed:expired", "오프라인")
	expect(t, a.do("POST", "/home/key-revoke", url.Values{"agent": {agentA}, "kid": {"key1"}}), 303)
	for _, path := range []string{"send", "pull", "persist", "ack"} {
		payload := any(map[string]string{})
		if path == "send" {
			payload = fresh
		}
		publicCall(t, h, "/v1/text/"+path, tokenA, payload, 401)
	}
}

func TestPublicHTTPAdmissionAndGateSafety(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	receipt := f.send(f.message(firstID, "gate-budget-key-1"), "agent_a", "", 200)
	claim := f.deliver(firstID)
	h := wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "public-gate-key-1", firstID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute))
	f.send(h, "agent_b", claim, 200)
	// Verified typed body, hint disagreement, XSS and malformed stored signature fail closed on both GET and POST.
	f.mutate(func(st *State) {
		e, _ := Parse(st.Messages[firstID].Envelope)
		e.Raw["render"] = map[string]any{"hint": "<script>approve all tools</script>"}
		st.Messages[firstID].Envelope = sign(t, f.private["agent_a"], e.Raw)
	})
	w := f.request("GET", "/owner/gates/"+gateID, f.tokens["agent_b_owner"], nil, "")
	expect(t, w, 200, "granularity_min", "참고 hint", "&lt;script&gt;")
	if strings.Contains(w.Body.String(), "<script>") {
		t.Fatal("gate XSS")
	}
	f.mutate(func(st *State) { st.Messages[firstID].Envelope = []byte(`{"invalid":"signature"}`) })
	w = f.request("GET", "/owner/gates/"+gateID, f.tokens["agent_b_owner"], nil, "")
	expect(t, w, 200, "원문 부재")
	if strings.Contains(w.Body.String(), `value="approve"`) {
		t.Fatal("bad signature enabled approval")
	}
	post := func(path, decision string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", path, strings.NewReader(url.Values{"csrf": {csrf(f.tokens["agent_b_owner"], gateID)}, "decision": {decision}}.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Authorization", "Bearer "+f.tokens["agent_b_owner"])
		w := httptest.NewRecorder()
		f.handler.ServeHTTP(w, r)
		return w
	}
	expect(t, post("/owner/gates/"+gateID, "approve"), 409)
	f.mutate(func(st *State) { st.Messages[firstID].Envelope = f.message(firstID, "gate-budget-key-1") })
	// Shared global new-work saturation survives another Service; cleanup admission and rate stay independent.
	f.mutate(func(st *State) {
		for i := 0; i < 16; i++ {
			st.HTTP[fmt.Sprint("busy", i)] = admission{Exp: time.Now().Add(20 * time.Second)}
		}
	})
	f.handler = (&Service{Pool: pool}).Handler()
	expect(t, f.request("POST", "/v1/pull", f.tokens["agent_b"], []byte("{}"), ""), 429)
	expect(t, post("/owner/gates/"+gateID, "approve"), 429)
	w = post("/owner/gates/"+gateID+"/deny", "deny")
	expect(t, w, 303)
	if w.Header().Get("Location") != "/owner/gates/"+gateID {
		t.Fatalf("owner deny redirected to %q", w.Header().Get("Location"))
	}
	f.mutate(func(st *State) {
		if st.Gates[gateID].State != "denied" {
			t.Fatal("deny while full")
		}
		st.HTTP = map[string]admission{}
	})
}

func TestSharedMessageCapacityHTTP(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	// Seed the first 99 accepted queue records, then race eight distinct high-priority sends for one remaining slot.
	f.mutate(func(st *State) {
		for i := 0; i < 99; i++ {
			id := fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", i)
			raw := wire(t, f.private["agent_a"], id, "agent_a", "agent_b", "schedule.query", fmt.Sprintf("capacity-seed-%04d", i), "", queryBody(), time.Now().Add(180*time.Second))
			e, _ := Parse(raw)
			if _, err := st.ingest("agent_a", e, "", time.Now().UTC()); err != nil {
				t.Fatal("seed", err)
			}
		}
	})
	status := make(chan int, 8)
	var wg sync.WaitGroup
	for i := 99; i < 107; i++ {
		id := fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", i)
		raw := map[string]any{}
		_ = json.Unmarshal(wire(t, f.private["agent_a"], id, "agent_a", "agent_b", "schedule.query", fmt.Sprintf("capacity-race-%04d", i), "", queryBody(), time.Now().Add(180*time.Second)), &raw)
		raw["priority"] = "high"
		body := sign(t, f.private["agent_a"], raw)
		wg.Add(1)
		go func() { defer wg.Done(); status <- f.request("POST", "/v1/send", f.tokens["agent_a"], body, "").Code }()
	}
	wg.Wait()
	close(status)
	wins := 0
	for code := range status {
		if code == 200 {
			wins++
		} else if code != 409 {
			t.Fatal("queue race", code)
		}
	}
	if wins != 1 {
		t.Fatal("remaining queue slot", wins)
	}
	f.handler = (&Service{Pool: pool}).Handler()
	f.send(f.message(firstID, "capacity-overflow-01"), "agent_a", "", 409)
	// A valid already-leased ACK still ends transport while new queue work is full; replay consumes no slot.
	lease := f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]string{}, 200)
	var env Envelope
	encoded, _ := json.Marshal(lease["envelope"])
	_ = json.Unmarshal(encoded, &env)
	delivery := map[string]string{"id": env.ID, "token": lease["lease_token"].(string)}
	f.call("POST", "/v1/persist", f.tokens["agent_b"], delivery, 200)
	f.call("POST", "/v1/ack", f.tokens["agent_b"], delivery, 200)
	f.mutate(func(st *State) {
		if len(st.Messages) != 100 {
			t.Fatal("receipt cap count")
		}
	})
}

func TestSharedExecutionClaimHTTP(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	ids := []string{}
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", i)
		ids = append(ids, id)
		raw := wire(t, f.private["agent_a"], id, "agent_a", "agent_b", "schedule.query", fmt.Sprintf("claim-capacity-%04d", i), "", queryBody(), time.Now().Add(180*time.Second))
		f.send(raw, "agent_a", "", 200)
		lease := f.call("POST", "/v1/pull", f.tokens["agent_b"], map[string]string{}, 200)
		d := map[string]string{"id": id, "token": lease["lease_token"].(string)}
		f.call("POST", "/v1/persist", f.tokens["agent_b"], d, 200)
		f.call("POST", "/v1/ack", f.tokens["agent_b"], d, 200)
	}
	codes := make(chan int, 5)
	var wg sync.WaitGroup
	for _, id := range ids {
		body, _ := json.Marshal(map[string]string{"id": id})
		wg.Add(1)
		go func() { defer wg.Done(); codes <- f.request("POST", "/v1/claim", f.tokens["agent_b"], body, "").Code }()
	}
	wg.Wait()
	close(codes)
	wins := 0
	for code := range codes {
		if code == 200 {
			wins++
		} else if code != 409 {
			t.Fatal("claim race", code)
		}
	}
	if wins != 4 {
		t.Fatal("claim winners", wins)
	}
	f.handler = (&Service{Pool: pool}).Handler()
	unclaimed := ""
	f.mutate(func(st *State) {
		for _, id := range ids {
			m := st.Messages[id]
			if !m.Claimed {
				unclaimed = id
			} else {
				m.Receipt.Exp = time.Now().Add(-time.Second)
			}
		}
	})
	f.call("POST", "/v1/claim", f.tokens["agent_b"], map[string]string{"id": unclaimed}, 200)
	f.mutate(func(st *State) {
		for _, id := range ids {
			if id != unclaimed && (st.Messages[id].Completion != "failed:expired" || st.Messages[id].ClaimToken != "") {
				t.Fatal("expired claims block safe termination")
			}
		}
	})
}

func TestHTTPConcurrencyAcrossInstances(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	entered := make(chan bool, 20)
	releaseNew, releaseClean := make(chan struct{}), make(chan struct{})
	var release sync.Once
	unblock := func() { release.Do(func() { close(releaseNew); close(releaseClean) }) }
	defer unblock()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := cleanupRequest(r)
		entered <- clean
		if clean {
			<-releaseClean
		} else {
			<-releaseNew
		}
		w.WriteHeader(200)
	})
	handlers := []http.Handler{f.s.boundedHTTP(next), (&Service{Pool: pool}).boundedHTTP(next)}
	f.addAgent("agent_d")
	f.addAgent("agent_e")
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	call := func(h http.Handler, path, revoker string) int {
		// Cleanup admission needs the caller's own target; an empty revoke would be new work.
		r := httptest.NewRequest("POST", path, strings.NewReader(`{"agent":"`+revoker+`","kid":"key1"}`))
		r.Header.Set("Authorization", "Bearer "+f.tokens[revoker+"_owner"])
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}
	done := make(chan int, 20)
	for i := 0; i < 16; i++ {
		go func(index int) { done <- call(handlers[index%2], "/v1/registry", "agent_a") }(i)
		select {
		case clean := <-entered:
			if clean {
				t.Fatal("new misclassified")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("new admission stalled")
		}
	}
	if call(handlers[0], "/v1/registry", "agent_a") != 429 {
		t.Fatal("17th shared HTTP request")
	}
	// Each owner holds at most one shared cleanup admission, so four owners fill the four cleanup records.
	for i, id := range []string{"agent_a", "agent_b", "agent_c", "agent_d"} {
		go func(index int) { done <- call(handlers[index%2], "/v1/key-revoke", id) }(i)
		select {
		case clean := <-entered:
			if !clean {
				t.Fatal("cleanup misclassified")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("cleanup blocked by new work")
		}
	}
	if call(handlers[0], "/v1/key-revoke", "agent_e") != 429 {
		t.Fatal("5th shared cleanup request")
	}
	unblock()
	for i := 0; i < 20; i++ {
		if <-done != 200 {
			t.Fatal("admitted request failed")
		}
	}
	f.mutate(func(st *State) {
		if len(st.HTTP) != 0 {
			t.Fatal("HTTP admission leaked")
		}
	})
}
