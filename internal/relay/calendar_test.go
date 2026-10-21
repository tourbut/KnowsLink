// Calendar unit checks: D010 scope boundaries, outward rounding, sealed tokens, and the disclosure result shape.
package relay

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRoundBusyWidensClipsAndMerges(t *testing.T) {
	start := time.Date(2026, 11, 2, 1, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	gran := 30 * time.Minute
	cases := []struct {
		name string
		busy []span
		want []span
	}{
		{"none", nil, []span{}},
		{"inside_one_unit", []span{{start.Add(5 * time.Minute), start.Add(10 * time.Minute)}}, []span{{start, start.Add(gran)}}},
		{"exact_unit_not_widened", []span{{start.Add(gran), start.Add(2 * gran)}}, []span{{start.Add(gran), start.Add(2 * gran)}}},
		{"crosses_units_and_window_end", []span{{start.Add(95 * time.Minute), end.Add(time.Hour)}}, []span{{start.Add(90 * time.Minute), end}}},
		{"before_window_start", []span{{start.Add(-time.Hour), start.Add(time.Minute)}}, []span{{start, start.Add(gran)}}},
		{"outside_window_dropped", []span{{end.Add(time.Minute), end.Add(time.Hour)}}, []span{}},
		{"overlap_merged_and_sorted", []span{{start.Add(50 * time.Minute), start.Add(70 * time.Minute)}, {start.Add(35 * time.Minute), start.Add(55 * time.Minute)}}, []span{{start.Add(gran), start.Add(90 * time.Minute)}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := roundBusy(c.busy, start, end, gran)
			if len(got) != len(c.want) {
				t.Fatalf("got %v want %v", got, c.want)
			}
			for i := range got {
				if !got[i].Start.Equal(c.want[i].Start) || !got[i].End.Equal(c.want[i].End) {
					t.Fatalf("got %v want %v", got, c.want)
				}
			}
		})
	}
	raw := disclosureJSON(start, end, 30, roundBusy(cases[1].busy, start, end, gran))
	if string(raw) != `{"busy":[{"end":"2026-11-02T01:30:00Z","start":"2026-11-02T01:00:00Z"}],"granularity_min":30,"window":{"end":"2026-11-02T03:00:00Z","start":"2026-11-02T01:00:00Z"}}` {
		t.Fatalf("disclosure JSON: %s", raw)
	}
	var shape map[string]any
	_ = json.Unmarshal(raw, &shape)
	if !disclosureShape(shape) {
		t.Fatal("server disclosure must satisfy the relay.result shape")
	}
}

func TestSealedTokenRoundTripAndTamper(t *testing.T) {
	key := make([]byte, calendarKeyLength)
	_, _ = rand.Read(key)
	sealed, err := sealToken(key, []byte("1//refresh-secret"), "mem_a")
	if err != nil || bytes.Contains(sealed, []byte("refresh-secret")) {
		t.Fatal("plaintext must not appear in the sealed token", err)
	}
	plain, err := openToken(key, sealed, "mem_a")
	if err != nil || string(plain) != "1//refresh-secret" {
		t.Fatal("round trip failed", err)
	}
	other := make([]byte, calendarKeyLength)
	_, _ = rand.Read(other)
	tampered := append([]byte(nil), sealed...)
	tampered[len(tampered)-1] ^= 1
	for name, c := range map[string]struct {
		key, sealed []byte
	}{"wrong_key": {other, sealed}, "tampered": {key, tampered}, "short": {key, sealed[:5]}, "bad_key": {key[:7], sealed}} {
		if _, err := openToken(c.key, c.sealed, "mem_a"); err == nil {
			t.Fatal(name, "opened")
		}
	}
	// A sealed value copied into another member's row must not open there (member ID is the GCM additional data).
	if _, err := openToken(key, sealed, "mem_b"); err == nil {
		t.Fatal("sealed token opened for another member")
	}
	if _, err := ConfigureGoogle(context.Background(), "client", "secret", "https://example.com/auth/google/callback", "not-base64*"); err == nil || !strings.Contains(err.Error(), "KNOWSLINK_CALENDAR_KEY") || strings.Contains(err.Error(), "secret") {
		t.Fatal("invalid key must fail closed without echoing values", err)
	}
}

// scheduleAllowed is the single D010 gate for accepted requests: unit set, window from now up to 14 days, 5 per friend per 24h.
func TestScheduleScopeBoundaries(t *testing.T) {
	now := time.Date(2026, 11, 2, 1, 7, 0, 0, time.UTC)
	st := newState()
	st.Owners["oa"], st.Owners["ob"], st.Owners["oc"] = &Owner{Active: true}, &Owner{Active: true}, &Owner{Active: true}
	st.Agents["agent_a"], st.Agents["agent_b"], st.Agents["agent_c"] = &Agent{Owner: "oa"}, &Agent{Owner: "ob"}, &Agent{Owner: "oc"}
	_, private, _ := ed25519.GenerateKey(rand.Reader)
	iso := func(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05Z") }
	query := func(start, end time.Time, minutes int) *Envelope {
		body := map[string]any{"window": map[string]any{"start": iso(start), "end": iso(end)}, "granularity_min": minutes}
		e, err := Parse(wire(t, private, "0199a3f2-4c10-7a11-8b22-334455667788", "agent_a", "agent_b", "schedule.query", "scope-key-000001", "", body, now.Add(time.Minute)))
		if err != nil {
			t.Fatal(err)
		}
		return e
	}
	slot := now.Truncate(30 * time.Minute) // 01:00, the unit the request time falls in
	for name, c := range map[string]struct {
		e  *Envelope
		ok bool
	}{
		"current_unit_ok":    {query(slot, slot.Add(time.Hour), 30), true},
		"fourteen_days_ok":   {query(slot, now.Add(scheduleScope.Window), 60), true},
		"unit_15_ok":         {query(slot.Add(15*time.Minute), slot.Add(time.Hour), 15), true},
		"unit_20_rejected":   {query(slot, slot.Add(time.Hour), 20), false},
		"past_start":         {query(slot.Add(-30*time.Minute), slot.Add(time.Hour), 30), false},
		"beyond_fourteen":    {query(slot, now.Add(scheduleScope.Window+time.Second), 30), false},
		"unit_60_past_slot":  {query(now.Truncate(time.Hour).Add(-time.Hour), slot.Add(time.Hour), 60), false},
		"unit_60_current_ok": {query(now.Truncate(time.Hour), slot.Add(2*time.Hour), 60), true},
	} {
		if err := st.scheduleAllowed(c.e, now); (err == nil) != c.ok {
			t.Fatalf("%s: err=%v want ok=%v", name, err, c.ok)
		}
	}
	// Five accepted requests from this friend in 24h fill the budget; an older one and another friend do not count.
	e := query(slot, slot.Add(time.Hour), 30)
	for i := 0; i < scheduleScope.PerFriend-1; i++ {
		st.Messages[string(rune('a'+i))] = &Message{Receipt: Receipt{From: "agent_a", To: "agent_b", Intent: "schedule.query", Accepted: now.Add(-time.Duration(i+1) * time.Hour)}}
	}
	st.Messages["old"] = &Message{Receipt: Receipt{From: "agent_a", To: "agent_b", Intent: "schedule.query", Accepted: now.Add(-scheduleScope.PerFriendWindow)}}
	st.Messages["other"] = &Message{Receipt: Receipt{From: "agent_c", To: "agent_b", Intent: "schedule.query", Accepted: now}}
	st.Messages["text"] = &Message{Receipt: Receipt{From: "agent_a", To: "agent_b", Intent: "relay.test.message", Accepted: now}}
	if err := st.scheduleAllowed(e, now); err != nil {
		t.Fatal("fifth request must pass", err)
	}
	st.Messages["fifth"] = &Message{Receipt: Receipt{From: "agent_a", To: "agent_b", Intent: "schedule.query", Accepted: now.Add(-time.Minute)}}
	if err := st.scheduleAllowed(e, now); err == nil || err.Error() != "schedule_scope" {
		t.Fatal("sixth request within 24h must be refused", err)
	}
}

// relay.result carries a disclosure only with status done and only in the server's window/unit/busy shape.
func TestResultDisclosureShape(t *testing.T) {
	_, private, _ := ed25519.GenerateKey(rand.Reader)
	disclosure := map[string]any{"window": map[string]any{"start": "2026-11-02T01:00:00Z", "end": "2026-11-02T03:00:00Z"}, "granularity_min": 30, "busy": []any{map[string]any{"start": "2026-11-02T01:00:00Z", "end": "2026-11-02T01:30:00Z"}}}
	result := func(body map[string]any) error {
		raw := map[string]any{"v": "relay.v1", "id": "0199a3f2-4c10-7a11-8b22-334455667788", "from": "agent_b", "to": "agent_a", "intent": "relay.result", "body": body, "deliver": "agent", "exp": "2026-11-02T01:00:00Z", "idempotency_key": "result-key-000001", "reply_to": "0199a3f2-4c10-7a11-8b22-334455667799", "sig": map[string]any{"alg": "Ed25519", "kid": "key1", "value": "x"}}
		data, _ := json.Marshal(raw)
		_, err := Parse(data)
		_ = private
		return err
	}
	withBusy := func(busy any) map[string]any {
		d := map[string]any{}
		for k, v := range disclosure {
			d[k] = v
		}
		d["busy"] = busy
		return d
	}
	if err := result(map[string]any{"status": "done", "disclosure": disclosure}); err != nil {
		t.Fatal("valid disclosure rejected", err)
	}
	if err := result(map[string]any{"status": "done"}); err != nil {
		t.Fatal("done without disclosure is a schema-valid result; ingest decides", err)
	}
	for name, body := range map[string]map[string]any{
		"denied_with_disclosure": {"status": "denied", "disclosure": disclosure},
		"title_in_busy":          {"status": "done", "disclosure": withBusy([]any{map[string]any{"start": "2026-11-02T01:00:00Z", "end": "2026-11-02T01:30:00Z", "summary": "비밀"}})},
		"busy_not_list":          {"status": "done", "disclosure": withBusy("none")},
		"busy_end_before_start":  {"status": "done", "disclosure": withBusy([]any{map[string]any{"start": "2026-11-02T01:30:00Z", "end": "2026-11-02T01:00:00Z"}})},
		"extra_field":            {"status": "done", "disclosure": map[string]any{"window": disclosure["window"], "granularity_min": 30, "busy": []any{}, "email": "x@example.test"}},
		"other_key":              {"status": "done", "result": disclosure},
	} {
		if err := result(body); err == nil {
			t.Fatal(name, "accepted")
		}
	}
}

// policy and the agent-visible gate state: an approved gate stays pending to the agent until its single lookup ends.
func TestLookupVisibilityAndSweep(t *testing.T) {
	now := time.Now()
	g := &Gate{State: "approved", Lookup: &Lookup{Started: now}}
	if lookupVisible(g) != "pending" {
		t.Fatal("running lookup must hide approval")
	}
	st := newState()
	st.Gates["g"] = g
	st.Messages["p"] = &Message{Receipt: Receipt{Exp: now.Add(time.Hour)}}
	g.Parent = "p"
	st.sweepLookups(now.Add(lookupLease - time.Second))
	if g.Lookup.State != "" {
		t.Fatal("lease ended early")
	}
	st.sweepLookups(now.Add(lookupLease))
	if g.Lookup.State != "failed" || lookupVisible(g) != "approved" {
		t.Fatal("expired lease must fail and become visible", g.Lookup.State)
	}
	g.Lookup.State, g.Lookup.Disclosure = "done", json.RawMessage(`{}`)
	st.sweepLookups(now)
	if len(g.Lookup.Disclosure) == 0 {
		t.Fatal("approved, unsent disclosure must stay")
	}
	st.Messages["p"].Completion = "done"
	st.sweepLookups(now)
	if len(g.Lookup.Disclosure) != 0 {
		t.Fatal("sent disclosure must be dropped")
	}
	if !strings.HasPrefix(policy("schedule.query"), "schedule.query:") || !strings.HasPrefix(policy("schedule.commit"), "schedule.commit:") {
		t.Fatal("policy labels keep their intent prefix for gate rows")
	}
}

type closeBody struct {
	io.Reader
	closed bool
}

func (b *closeBody) Close() error { b.closed = true; return nil }

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// The revoke response body is closed whatever the provider answers.
func TestRevokeClosesBodyOnEveryStatus(t *testing.T) {
	for _, status := range []int{200, 400, 503} {
		body := &closeBody{Reader: strings.NewReader(`{"error":"invalid_token"}`)}
		client := &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: status, Body: body, Header: http.Header{}, Request: r}, nil
		})}
		s := &Service{Google: &GoogleLogin{client: client, revokeURL: "https://revoke.example.test/revoke"}}
		s.revokeToken(context.Background(), "refresh")
		if !body.closed {
			t.Fatalf("status %d: revoke response body left open", status)
		}
	}
}
