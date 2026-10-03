// Frozen wire tests cover strict parsing, JCS, signed fields, and semantic replay boundaries.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

const firstID = "0199a3f2-4c10-7a11-8b22-334455667788"
const secondID = "0199a3f2-4c10-7a11-8b22-334455667799"
const gateID = "0199a3f2-4c10-7a11-8b22-3344556677aa"

func wire(t *testing.T, private ed25519.PrivateKey, id, from, to, intent, key, reply string, body map[string]any, exp time.Time) []byte {
	t.Helper()
	deliver := "agent"
	if intent == "relay.approval.request" {
		deliver = "human"
	}
	raw := map[string]any{"v": "relay.v1", "id": id, "from": from, "to": to, "intent": intent, "body": body, "deliver": deliver, "exp": exp.UTC().Format("2006-01-02T15:04:05Z"), "idempotency_key": key, "sig": map[string]any{"alg": "Ed25519", "kid": "key1", "value": ""}}
	if reply != "" {
		raw["reply_to"] = reply
	}
	return sign(t, private, raw)
}
func sign(t *testing.T, private ed25519.PrivateKey, raw map[string]any) []byte {
	t.Helper()
	data, _ := json.Marshal(raw)
	e, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	raw["sig"].(map[string]any)["value"] = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, e.SigningBytes()))
	data, _ = json.Marshal(raw)
	return data
}
func queryBody() map[string]any {
	return map[string]any{"window": map[string]any{"start": "2026-10-03T10:00:00Z", "end": "2026-10-03T11:00:00Z"}, "granularity_min": 30}
}
func TestFrozenParsingAndSigning(t *testing.T) {
	if !RegistryValid() {
		t.Fatal("registry integrity failure")
	}
	public, private, _ := ed25519.GenerateKey(rand.Reader)
	raw := wire(t, private, "0199a3f2-4c10-7a11-8b22-334455667788", "agent_a", "agent_b", "schedule.query", "idempotent-key-01", "", queryBody(), time.Now().Add(time.Minute))
	e, err := Parse(raw)
	if err != nil || !e.Verify(public) {
		t.Fatal("valid envelope rejected", err)
	}
	for _, bad := range []string{strings.Replace(string(raw), `"v":"relay.v1"`, `"v":"relay.v1","v":"relay.v1"`, 1), strings.Replace(string(raw), `"from":"agent_a"`, `"from":"AGENT_A"`, 1), strings.Replace(string(raw), `"deliver":"agent"`, `"deliver":"both"`, 1), strings.Replace(string(raw), `"body":{`, `"unknown":1,"body":{`, 1), strings.Replace(string(raw), `"granularity_min":30`, `"granularity_min":30,"text":"execute"`, 1), strings.Replace(string(raw), `"body":{`, `"render":{"hint":"\ud800"},"body":{`, 1), strings.Replace(string(raw), `"body":{`, `"ext":null,"body":{`, 1)} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Fatalf("accepted invalid JSON: %s", bad)
		}
	}
	changed := strings.Replace(string(raw), `"granularity_min":30`, `"granularity_min":31`, 1)
	other, err := Parse([]byte(changed))
	if err != nil || other.Verify(public) {
		t.Fatal("tamper not rejected")
	}
	bytes, err := Canonical(map[string]any{"n": 1e30, "a": 0.000001})
	if err != nil || string(bytes) != `{"a":0.000001,"n":1e+30}` {
		t.Fatalf("JCS mismatch %s %v", bytes, err)
	}
	replay := wire(t, private, "0199a3f2-4c10-7a11-8b22-334455667799", "agent_a", "agent_b", "schedule.query", "idempotent-key-01", "", queryBody(), time.Now().Add(2*time.Minute))
	second, _ := Parse(replay)
	if second.Digest() != e.Digest() {
		t.Fatal("id/exp affect semantic digest")
	}
}

func TestFrozenFieldLimits(t *testing.T) {
	_, private, _ := ed25519.GenerateKey(rand.Reader)
	raw := wire(t, private, firstID, "agent_a", "agent_b", "schedule.query", "idempotency-key-01", "", queryBody(), time.Now().Add(time.Minute))
	check := func(change func(map[string]any), valid bool) {
		t.Helper()
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
		change(m)
		bytes, _ := json.Marshal(m)
		_, err := Parse(bytes)
		if (err == nil) != valid {
			t.Fatalf("field limit valid=%v err=%v", valid, err)
		}
	}
	check(func(m map[string]any) { m["render"] = map[string]any{"hint": strings.Repeat("x", 1024)} }, true)
	check(func(m map[string]any) { m["render"] = map[string]any{"hint": strings.Repeat("x", 1025)} }, false)
	check(func(m map[string]any) { m["render"] = map[string]any{"hint": "e\u0301"} }, false)
	check(func(m map[string]any) { m["idempotency_key"] = strings.Repeat("k", 128) }, true)
	check(func(m map[string]any) { m["idempotency_key"] = strings.Repeat("k", 129) }, false)
	check(func(m map[string]any) { m["idempotency_key"] = strings.Repeat("k", 15) }, false)
	check(func(m map[string]any) { m["sig"].(map[string]any)["kid"] = "https://key.example/key" }, false)
	evidence := func(n int) []any {
		items := []any{}
		for i := 0; i < n; i++ {
			items = append(items, map[string]any{"ref": "urn:synthetic:" + string(rune('a'+i))})
		}
		return items
	}
	check(func(m map[string]any) { m["evidence"] = evidence(8) }, true)
	check(func(m map[string]any) { m["evidence"] = evidence(9) }, false)
	check(func(m map[string]any) { m["evidence"] = []any{map[string]any{"ref": "https:"}} }, false)
	check(func(m map[string]any) {
		m["evidence"] = []any{map[string]any{"ref": "urn:a"}, map[string]any{"ref": "urn:a"}}
	}, false)
	check(func(m map[string]any) {
		m["intent"] = "schedule.commit"
		m["body"] = map[string]any{"slot": queryBody()["window"], "timezone": "UTC", "commitment": strings.Repeat("x", 16*1024)}
	}, false)
	check(func(m map[string]any) { m["ext"] = map[string]any{"com.example": map[string]any{"approve": true}} }, true)
	check(func(m map[string]any) {
		m["trace"] = map[string]any{"traceparent": "00-00000000000000000000000000000000-0000000000000000-01"}
	}, false)
}
