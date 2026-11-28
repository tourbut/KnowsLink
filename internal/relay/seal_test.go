// SAR-EVENTS-SECRET-001: stored secrets round-trip only under the same key and record, and Events stay metadata-only.
package relay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"
)

func testSealer(t *testing.T, fill byte) *Sealer {
	t.Helper()
	s, err := NewSealer(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{fill}, 32)))
	if err != nil || s == nil {
		t.Fatal("sealer", err)
	}
	return s
}

func TestSealerRoundTripAndRefusals(t *testing.T) {
	if s, err := NewSealer(""); s != nil || err != nil {
		t.Fatal("absent key must turn the feature off, not fail startup")
	}
	for _, bad := range []string{"not base64!", base64.StdEncoding.EncodeToString(make([]byte, 16))} {
		if _, err := NewSealer(bad); err == nil || strings.Contains(err.Error(), bad) {
			t.Fatal("malformed key accepted or echoed", bad)
		}
	}
	s, other := testSealer(t, 1), testSealer(t, 2)
	secret := "whsec_" + base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, 32))
	sealed := s.Seal(secret, "sub_a")
	if !strings.HasPrefix(sealed, sealPrefix) || strings.Contains(sealed, secret) || sealed == s.Seal(secret, "sub_a") {
		t.Fatal("sealed form leaks plaintext or reuses a nonce")
	}
	if got, err := s.Open(sealed, "sub_a"); err != nil || got != secret {
		t.Fatal("round trip", err)
	}
	tampered := sealed[:len(sealed)-2] + "AA"
	var none *Sealer
	for name, open := range map[string]func() (string, error){
		"other record": func() (string, error) { return s.Open(sealed, "sub_b") },
		"other key":    func() (string, error) { return other.Open(sealed, "sub_a") },
		"tampered":     func() (string, error) { return s.Open(tampered, "sub_a") },
		"no key":       func() (string, error) { return none.Open(sealed, "sub_a") },
		"plain text":   func() (string, error) { return s.Open(secret, "sub_a") },
	} {
		if got, err := open(); err == nil || got != "" {
			t.Fatal("opened under", name)
		}
	}
}

func TestRemoteSecretMigrationReadsPlainTextUntilResealed(t *testing.T) {
	s := &Service{Sealer: testSealer(t, 1)}
	legacy := "whsec_" + base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{3}, 32))
	// A pre-sealing record still signs and matches its refresh; the refresh then writes Sealer.Seal.
	if got, err := s.remoteStoredSecret(legacy, "sub_a"); err != nil || got != legacy || !s.remoteSecretIs(legacy, "sub_a", legacy) {
		t.Fatal("legacy plain-text read fallback", err)
	}
	sealed := s.Sealer.Seal(legacy, "sub_a")
	if !s.remoteSecretIs(sealed, "sub_a", legacy) || s.remoteSecretIs(sealed, "sub_b", legacy) || s.remoteSecretIs(sealed, "sub_a", legacy+"x") {
		t.Fatal("sealed comparison")
	}
	// Losing the key keeps legacy records readable but never opens sealed ones; delivery then drops the subscription.
	lost := &Service{}
	if _, err := lost.remoteStoredSecret(sealed, "sub_a"); err == nil || lost.remoteSecretIs("", "sub_a", "") {
		t.Fatal("sealed secret opened without key")
	}
}

func TestRemoteSubscribeRefusedWithoutSealKey(t *testing.T) {
	s := &Service{}
	params := `{"name":"knowslink.text.arrived","arguments":{},"delivery":{"mode":"webhook","url":"https://callback.example/e","secret":"whsec_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="}}`
	w := httptest.NewRecorder()
	// Refused before any transaction, so a nil Pool proves no state (and no plain-text secret) is written.
	s.remoteSubscription(w, httptest.NewRequest("POST", "/mcp", nil), "install", remoteRequest{ID: json.RawMessage(`1`), Method: "events/subscribe", Params: json.RawMessage(params)}, false)
	if !strings.Contains(w.Body.String(), "events_webhook_unavailable") || strings.Contains(w.Body.String(), "whsec_") {
		t.Fatal("subscribe without seal key", w.Body.String())
	}
}

func TestRemoteEventsDualShapeAndMetadataOnlyBody(t *testing.T) {
	// Dots regression: an undeclared client keeps the OpenAI webhook-only definition and webhook-id = eventId.
	if d := remoteEventDefinition(false, true)["delivery"].([]string); !slices.Equal(d, []string{"webhook"}) {
		t.Fatal("legacy delivery list changed", d)
	}
	if d := remoteEventDefinition(true, true)["delivery"].([]string); !slices.Equal(d, []string{"webhook", "poll"}) {
		t.Fatal("beta delivery list", d)
	}
	if d := remoteEventDefinition(true, false)["delivery"].([]string); !slices.Equal(d, []string{"poll"}) {
		t.Fatal("beta without seal key must not offer webhook", d)
	}
	sub := remoteSubscription{EventID: "evt_1", Message: "msg_1", EventAt: time.Unix(0, 0)}
	if remoteHookID(sub) != "evt_1" {
		t.Fatal("OpenAI webhook-id must equal eventId")
	}
	sub.Beta = true
	if remoteHookID(sub) != "msg_1" {
		t.Fatal("Beta webhook-id derived from eventId")
	}
	modern := `{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{"extensions":{"io.modelcontextprotocol/events":{}}}}}`
	if !remoteDeclaresEvents(json.RawMessage(modern)) || remoteDeclaresEvents(json.RawMessage(`{"_meta":{"io.modelcontextprotocol/clientCapabilities":{}}}`)) {
		t.Fatal("per-request Events declaration")
	}
	var body map[string]any
	if err := json.Unmarshal(remoteEventBody(sub), &body); err != nil {
		t.Fatal(err)
	}
	keys := []string{}
	for k := range body {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	data, _ := json.Marshal(body["data"])
	if !slices.Equal(keys, []string{"cursor", "data", "eventId", "name", "timestamp"}) || string(data) != `{"arrival":true}` {
		t.Fatal("event body carries more than arrival metadata", keys, string(data))
	}
}

func TestRemotePollIsStatelessAndMetadataOnly(t *testing.T) {
	st, keys, now := textFixture()
	id, install, _, _ := remoteFixture(t, st, "agent_b", now)
	st.Pairs[pairID("agent_a", install.Agent)] = &Pair{A: "agent_a", B: install.Agent, State: "active", Generation: 1}
	poll := func(params string) map[string]any {
		t.Helper()
		v, err := st.remotePoll(id, install, json.RawMessage(params), now)
		if err != nil {
			t.Fatal(err)
		}
		return v.(map[string]any)
	}
	start := poll(`{"name":"knowslink.text.arrived","arguments":{},"cursor":null}`)
	cursor := start["cursor"].(string)
	if len(start["events"].([]any)) != 0 || poll(`{"name":"knowslink.text.arrived"}`)["cursor"] != cursor {
		t.Fatal("null or absent cursor must start from now with no events")
	}
	text := signedText(firstID, "agent_a", install.Agent, "remote-poll-key-01", "", "secret body never polled", now.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", text, now); err != nil {
		t.Fatal(err)
	}
	got := poll(`{"name":"knowslink.text.arrived","cursor":"` + cursor + `","maxEvents":10}`)
	events := got["events"].([]any)
	raw, _ := json.Marshal(got)
	if len(events) != 1 || got["cursor"] == cursor || strings.Contains(string(raw), "secret body") || strings.Contains(string(raw), "agent_a") || strings.Contains(string(raw), firstID) {
		t.Fatal("poll arrival", string(raw))
	}
	if len(poll(`{"name":"knowslink.text.arrived","cursor":"` + got["cursor"].(string) + `"}`)["events"].([]any)) != 0 {
		t.Fatal("same waiting set repeated an event")
	}
	for _, bad := range []string{`{"name":"other"}`, `{"name":"knowslink.text.arrived","arguments":{"x":1}}`, `{"name":"knowslink.text.arrived","maxEvents":0}`, `{"name":"knowslink.text.arrived","extra":1}`} {
		if _, err := st.remotePoll(id, install, json.RawMessage(bad), now); err == nil {
			t.Fatal("invalid poll accepted", bad)
		}
	}
}
