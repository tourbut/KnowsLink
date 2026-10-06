// Package relay validates frozen relay.v1 and preserves one parsed object for signing and processing.
package relay

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	jcs "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
	"golang.org/x/text/unicode/norm"
)

var agentPattern = regexp.MustCompile(`^[a-z][a-z0-9_:-]*$`)
var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
var kidPattern = regexp.MustCompile(`^[A-Za-z0-9_:-]+$`)
var extensionPattern = regexp.MustCompile(`^[a-z0-9-]+(\.[a-z0-9-]+)+$`)
var tracePattern = regexp.MustCompile(`^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$`)

type Envelope struct {
	V        string           `json:"v"`
	ID       string           `json:"id"`
	From     string           `json:"from"`
	To       string           `json:"to"`
	Intent   string           `json:"intent"`
	Body     map[string]any   `json:"body"`
	Deliver  string           `json:"deliver"`
	Exp      string           `json:"exp"`
	Sig      Signature        `json:"sig"`
	Key      string           `json:"idempotency_key"`
	Priority string           `json:"priority,omitempty"`
	Evidence []map[string]any `json:"evidence,omitempty"`
	ReplyTo  string           `json:"reply_to,omitempty"`
	Trace    map[string]any   `json:"trace,omitempty"`
	Render   map[string]any   `json:"render,omitempty"`
	Ext      map[string]any   `json:"ext,omitempty"`
	Raw      map[string]any   `json:"-"`
}
type Signature struct {
	Alg   string `json:"alg"`
	Kid   string `json:"kid"`
	Value string `json:"value"`
}

func fault(code string) error { return errors.New(code) }
func Canonical(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return jcs.Transform(raw)
}
func Strict(raw []byte, target any) error {
	if !utf8.Valid(raw) {
		return fault("invalid_json")
	}
	// Transform rejects duplicate properties and lone UTF-16 surrogates before encoding/json can replace them.
	if _, err := jcs.Transform(raw); err != nil {
		return fault("invalid_json")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fault("invalid_schema")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fault("invalid_json")
	}
	return nil
}
func object(value any, required []string, optional ...string) bool {
	m, ok := value.(map[string]any)
	if !ok {
		return false
	}
	allowed := map[string]bool{}
	for _, k := range required {
		allowed[k] = true
		if _, ok := m[k]; !ok {
			return false
		}
	}
	for _, k := range optional {
		allowed[k] = true
	}
	for k := range m {
		if !allowed[k] {
			return false
		}
	}
	return true
}
func text(value any) (string, bool) { s, ok := value.(string); return s, ok }
func instant(value any) (time.Time, bool) {
	s, ok := text(value)
	t, e := time.Parse("2006-01-02T15:04:05Z", s)
	return t, ok && e == nil && t.Format("2006-01-02T15:04:05Z") == s
}
func interval(value any) bool {
	if !object(value, []string{"start", "end"}) {
		return false
	}
	m := value.(map[string]any)
	a, ok := instant(m["start"])
	b, ok2 := instant(m["end"])
	return ok && ok2 && b.After(a)
}
func Parse(raw []byte) (*Envelope, error) {
	if len(raw) > rawEnvelopeLimit {
		return nil, fault("invalid_schema")
	}
	e := new(Envelope)
	if err := Strict(raw, e); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &e.Raw); err != nil {
		return nil, fault("invalid_json")
	}
	if !object(e.Raw, []string{"v", "id", "from", "to", "intent", "body", "deliver", "exp", "sig", "idempotency_key"}, "priority", "evidence", "reply_to", "trace", "render", "ext") {
		return nil, fault("invalid_schema")
	}
	if e.V != "relay.v1" || !uuidPattern.MatchString(e.ID) || !agentPattern.MatchString(e.From) || !agentPattern.MatchString(e.To) || len(e.Intent) > 64 || !kidPattern.MatchString(e.Sig.Kid) || e.Sig.Alg != "Ed25519" {
		return nil, fault("invalid_schema")
	}
	if e.Deliver != "agent" && e.Deliver != "human" {
		return nil, fault("invalid_deliver")
	}
	if _, ok := instant(e.Exp); !ok {
		return nil, fault("invalid_schema")
	}
	if len(e.Key) < 16 || len(e.Key) > 128 {
		return nil, fault("invalid_schema")
	}
	for _, c := range e.Key {
		if c > 127 || c < 32 {
			return nil, fault("invalid_schema")
		}
	}
	if e.Priority != "" && e.Priority != "low" && e.Priority != "normal" && e.Priority != "high" {
		return nil, fault("invalid_schema")
	}
	if e.ReplyTo != "" && !uuidPattern.MatchString(e.ReplyTo) {
		return nil, fault("invalid_schema")
	}
	if _, ok := e.Raw["sig"].(map[string]any); !ok || !object(e.Raw["sig"], []string{"alg", "kid", "value"}) {
		return nil, fault("invalid_schema")
	}
	body, err := Canonical(e.Body)
	if err != nil || len(body) > 16*1024 {
		return nil, fault("invalid_schema")
	}
	valid := false
	switch e.Intent {
	case "relay.test.message":
		value, ok := text(e.Body["text"])
		valid = object(e.Body, []string{"text"}) && ok && strings.TrimSpace(value) != "" && len(value) <= 4096 && e.Deliver == "agent" && e.ReplyTo == "" && e.Evidence == nil && e.Ext == nil && e.Render == nil
	case "schedule.query":
		n, ok := e.Body["granularity_min"].(float64)
		valid = object(e.Body, []string{"window", "granularity_min"}) && interval(e.Body["window"]) && ok && n > 0 && n == float64(int64(n))
	case "schedule.commit":
		zone, z := text(e.Body["timezone"])
		commit, c := text(e.Body["commitment"])
		_, err := time.LoadLocation(zone)
		valid = object(e.Body, []string{"slot", "timezone", "commitment"}) && interval(e.Body["slot"]) && z && err == nil && c && commit != ""
	case "relay.approval.request":
		digest, d := text(e.Body["request_digest"])
		reason, r := text(e.Body["reason"])
		decoded, err := base64.RawURLEncoding.DecodeString(digest)
		valid = object(e.Body, []string{"reason", "request_digest"}) && d && err == nil && len(decoded) == 32 && r && (reason == "permission_required" || reason == "judgment_required") && e.Deliver == "human" && e.ReplyTo != ""
	case "relay.result":
		status, ok := text(e.Body["status"])
		valid = object(e.Body, []string{"status"}) && ok && (status == "denied" || status == "failed" || status == "done") && e.Deliver == "agent" && e.ReplyTo != ""
	default:
		return nil, fault("unsupported_intent")
	}
	if !valid {
		return nil, fault("invalid_schema")
	}
	if e.Render != nil {
		hint, ok := text(e.Render["hint"])
		if !object(e.Render, []string{"hint"}, "lang") || !ok || len(hint) > 1024 || !norm.NFC.IsNormalString(hint) {
			return nil, fault("invalid_schema")
		}
		if lang, exists := e.Render["lang"]; exists {
			if _, ok := text(lang); !ok {
				return nil, fault("invalid_schema")
			}
		}
	}
	if len(e.Evidence) > 8 {
		return nil, fault("invalid_schema")
	}
	previous := ""
	for _, item := range e.Evidence {
		ref, ok := text(item["ref"])
		parsed, err := url.Parse(ref)
		if !object(item, []string{"ref"}, "media_type", "sha256") || !ok || err != nil || (parsed.Scheme != "https" && parsed.Scheme != "urn") || ref <= previous || (parsed.Scheme == "https" && parsed.Hostname() == "") || (parsed.Scheme == "urn" && parsed.Opaque == "") {
			return nil, fault("invalid_schema")
		}
		previous = ref
		for k, v := range item {
			if _, ok := text(v); !ok || k == "sha256" && !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(v.(string)) {
				return nil, fault("invalid_schema")
			}
		}
	}
	if e.Trace != nil {
		parent, ok := text(e.Trace["traceparent"])
		if !object(e.Trace, []string{"traceparent"}, "tracestate") || !ok || !tracePattern.MatchString(parent) || strings.HasPrefix(parent, "ff-") || strings.Split(parent, "-")[1] == strings.Repeat("0", 32) || strings.Split(parent, "-")[2] == strings.Repeat("0", 16) {
			return nil, fault("invalid_schema")
		}
		if value, exists := e.Trace["tracestate"]; exists {
			if _, ok := text(value); !ok {
				return nil, fault("invalid_schema")
			}
		}
	}
	for k := range e.Ext {
		if !extensionPattern.MatchString(k) {
			return nil, fault("invalid_schema")
		}
	}
	// Reject explicit nulls and wrong container types, including optional fields.
	for _, k := range []string{"body", "trace", "render", "ext"} {
		if v, ok := e.Raw[k]; ok {
			if _, ok := v.(map[string]any); !ok {
				return nil, fault("invalid_schema")
			}
		}
	}
	if v, ok := e.Raw["evidence"]; ok {
		if _, ok := v.([]any); !ok {
			return nil, fault("invalid_schema")
		}
	}
	for _, k := range []string{"priority", "reply_to"} {
		if v, ok := e.Raw[k]; ok {
			s, ok := text(v)
			if !ok || s == "" {
				return nil, fault("invalid_schema")
			}
		}
	}
	return e, nil
}
func (e *Envelope) SigningBytes() []byte {
	unsigned := map[string]any{}
	for k, v := range e.Raw {
		if k != "sig" {
			unsigned[k] = v
		}
	}
	canonical, _ := Canonical(unsigned)
	return append([]byte("SILENT-AGENT-RELAY\x00relay.v1\x00Ed25519\x00"+e.Sig.Kid+"\x00"), canonical...)
}
func (e *Envelope) Verify(public []byte) bool {
	sig, err := base64.RawURLEncoding.DecodeString(e.Sig.Value)
	return err == nil && len(public) == ed25519.PublicKeySize && ed25519.Verify(public, e.SigningBytes(), sig)
}
func (e *Envelope) Digest() string {
	priority := e.Priority
	if priority == "" {
		priority = "normal"
	}
	evidence := any(e.Evidence)
	if e.Evidence == nil {
		evidence = []any{}
	}
	ext := e.Ext
	if ext == nil {
		ext = map[string]any{}
	}
	reply := any(e.ReplyTo)
	if e.ReplyTo == "" {
		reply = nil
	}
	raw, _ := Canonical(map[string]any{"v": e.V, "from": e.From, "to": e.To, "intent": e.Intent, "body": e.Body, "deliver": e.Deliver, "priority": priority, "evidence": evidence, "reply_to": reply, "ext": ext})
	sum := sha256.Sum256(raw)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
func popBytes(owner, agent, kid, public string) []byte {
	return []byte(strings.Join([]string{"KNOWSLINK-KEY-POP", owner, agent, kid, public}, "\x00"))
}
