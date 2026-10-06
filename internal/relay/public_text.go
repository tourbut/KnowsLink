// Public connection-check text uses a separate signed wire, current member authorization and the shared short-lived inbox.
package relay

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

const publicTextIntent = "knowslink.text"
const rawEnvelopeLimit = 32768

type PublicText struct {
	V       string    `json:"v"`
	ID      string    `json:"id"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Text    string    `json:"text"`
	Exp     string    `json:"exp"`
	Key     string    `json:"idempotency_key"`
	ReplyTo string    `json:"reply_to,omitempty"`
	Sig     Signature `json:"sig"`
}

func parsePublicText(raw []byte) (*PublicText, error) {
	e := new(PublicText)
	if len(raw) > rawEnvelopeLimit {
		return nil, fault("invalid_schema")
	}
	if err := Strict(raw, e); err != nil {
		return nil, err
	}
	var fields map[string]any
	_ = json.Unmarshal(raw, &fields)
	if !object(fields, []string{"v", "id", "from", "to", "text", "exp", "idempotency_key", "sig"}, "reply_to") || !object(fields["sig"], []string{"alg", "kid", "value"}) {
		return nil, fault("invalid_schema")
	}
	if e.V != "knowslink.text.v1" || !uuidPattern.MatchString(e.ID) || !agentPattern.MatchString(e.From) || !agentPattern.MatchString(e.To) || e.From == e.To || !kidPattern.MatchString(e.Sig.Kid) || len(e.Sig.Kid) > 128 || e.Sig.Alg != "Ed25519" || strings.TrimSpace(e.Text) == "" || len(e.Text) > 4096 || len(e.Key) < 16 || len(e.Key) > 128 {
		return nil, fault("invalid_schema")
	}
	if _, ok := instant(e.Exp); !ok {
		return nil, fault("invalid_schema")
	}
	for _, c := range e.Key {
		if c < 32 || c > 126 {
			return nil, fault("invalid_schema")
		}
	}
	if v, exists := fields["reply_to"]; exists {
		if s, ok := text(v); !ok || !uuidPattern.MatchString(s) {
			return nil, fault("invalid_schema")
		}
	}
	return e, nil
}
func (e *PublicText) signingBytes() []byte {
	raw, _ := Canonical(map[string]any{"v": e.V, "id": e.ID, "from": e.From, "to": e.To, "text": e.Text, "exp": e.Exp, "idempotency_key": e.Key, "reply_to": e.ReplyTo})
	return append([]byte("KNOWSLINK-TEXT\x00knowslink.text.v1\x00Ed25519\x00"+e.Sig.Kid+"\x00"), raw...)
}
func (e *PublicText) digest() string {
	raw, _ := Canonical(map[string]any{"v": e.V, "from": e.From, "to": e.To, "text": e.Text, "reply_to": e.ReplyTo})
	sum := sha256.Sum256(raw)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
func (st *State) ingestPublicText(agent string, e *PublicText, now time.Time) (any, error) {
	a, b := st.Agents[e.From], st.Agents[e.To]
	if agent != e.From || a == nil || b == nil || a.Revoked || b.Revoked || !st.publicOwner(a.Owner) || !st.publicOwner(b.Owner) {
		return nil, fault("sender_not_allowed")
	}
	k := a.Keys[e.Sig.Kid]
	if k == nil || k.Revoked {
		return nil, fault("invalid_auth")
	}
	sig, err := base64.RawURLEncoding.DecodeString(e.Sig.Value)
	if err != nil || len(k.Public) != ed25519.PublicKeySize || !ed25519.Verify(k.Public, e.signingBytes(), sig) {
		return nil, fault("invalid_signature")
	}
	p := st.Pairs[pairID(e.From, e.To)]
	if p == nil || p.State != "active" {
		return nil, fault("human_invite_required")
	}
	m := &Message{Receipt: Receipt{e.ID, e.From, e.To, publicTextIntent, e.digest(), time.Time{}, now, "queued"}, Key: e.Key, Kid: e.Sig.Kid, Generation: p.Generation, Parent: e.ReplyTo, Deliver: "agent"}
	if !st.current(m) {
		return nil, fault("sender_not_allowed")
	}
	if id := st.Idempotency[e.From+"/"+e.Key]; id != "" {
		old := st.Messages[id]
		if old == nil || old.Receipt.Digest != m.Receipt.Digest {
			return nil, fault("idempotency_conflict")
		}
		if old.Generation != p.Generation || !st.current(old) {
			return nil, fault("sender_not_allowed")
		}
		return old.Receipt, nil
	}
	exp, _ := instant(e.Exp)
	if !now.Before(exp) {
		return nil, fault("expired")
	}
	if exp.After(now.Add(180 * time.Second)) {
		return nil, fault("ttl_too_long")
	}
	if st.Messages[e.ID] != nil {
		return nil, fault("id_collision")
	}
	var parent *Message
	if e.ReplyTo != "" {
		parent = st.Messages[e.ReplyTo]
		if parent == nil || parent.Receipt.Intent != publicTextIntent || parent.Parent != "" || parent.Receipt.From != e.To || parent.Receipt.To != e.From || parent.Receipt.State != "delivered" || !parent.Persisted || !st.current(parent) || parent.Generation != p.Generation || !now.Before(parent.Receipt.Exp) {
			return nil, fault("sender_not_allowed")
		}
		if parent.ReplyID != "" {
			return nil, fault("duplicate_result")
		}
		if exp.After(parent.Receipt.Exp) {
			return nil, fault("expired")
		}
	}
	if !st.messageCapacity(false) {
		return nil, fault("capacity")
	}
	m.Receipt.Exp = exp
	m.Envelope, _ = json.Marshal(e)
	st.Messages[e.ID] = m
	st.Idempotency[e.From+"/"+e.Key] = e.ID
	if parent != nil {
		parent.ReplyID = e.ID
	}
	return m.Receipt, nil
}
func (st *State) publicOwner(owner string) bool {
	if st.Owners[owner] == nil || !st.Owners[owner].Active {
		return false
	}
	for _, m := range st.Members {
		if m.Owner == owner && m.Active {
			return true
		}
	}
	return false
}
func (s *Service) textRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/text/send", func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, rawEnvelopeLimit))
		if err != nil {
			respond(w, nil, fault("invalid_json"))
			return
		}
		e, err := parsePublicText(raw)
		if err != nil {
			respond(w, nil, err)
			return
		}
		v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			a, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			return st.ingestPublicText(a, e, now)
		})
		respond(w, v, err)
	})
	for _, action := range []string{"pull", "persist", "ack"} {
		mux.HandleFunc("POST /v1/text/"+action, func(w http.ResponseWriter, r *http.Request) {
			var c command
			if err := decodeCommand(w, r, &c); err != nil {
				respond(w, nil, err)
				return
			}
			v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
				a, err := st.principal(bearer(r), "agent")
				if err != nil {
					return nil, err
				}
				if !st.publicOwner(st.Agents[a].Owner) {
					return nil, fault("sender_not_allowed")
				}
				if action == "pull" {
					return st.leaseMatching(a, now, false, true)
				}
				m := st.Messages[c.ID]
				if m == nil || m.Receipt.Intent != publicTextIntent {
					return nil, fault("sender_not_allowed")
				}
				return st.operateAs(action, a, "agent", c, now, false)
			})
			respond(w, v, err)
		})
	}
}
