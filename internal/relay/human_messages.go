// Owner-written text is authenticated by the browser session, never represented as an agent signature.
package relay

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
	"unicode/utf8"
)

type OwnerText struct {
	V          string `json:"v"`
	ID         string `json:"id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Text       string `json:"text"`
	Exp        string `json:"exp"`
	Key        string `json:"idempotency_key"`
	Authorship string `json:"authorship"`
}

func (st *State) sendOwnerText(member, from, to, body, attempt string, now time.Time) (*Message, error) {
	m := st.Members[member]
	a := st.Agents[from]
	if m == nil || a == nil || a.Owner != m.Owner || m.Name == "" {
		return nil, fault("sender_not_allowed")
	}
	if _, eligible := st.humanMember(m.Owner); eligible == nil {
		return nil, fault("sender_not_allowed")
	}
	if !utf8.ValidString(body) || strings.TrimSpace(body) == "" || len(body) > 4096 || len(attempt) < 16 || len(attempt) > 96 {
		return nil, fault("invalid_schema")
	}
	for _, r := range attempt {
		if r < 33 || r > 126 {
			return nil, fault("invalid_schema")
		}
	}
	scope, generation, ok := st.authority(from, to)
	if !ok {
		return nil, fault("human_invite_required")
	}
	raw, _ := json.Marshal([]string{"owner", from, to, body})
	sum := sha256.Sum256(raw)
	digest := base64.RawURLEncoding.EncodeToString(sum[:])
	key := "web:" + attempt
	if id := st.Idempotency[from+"/"+key]; id != "" {
		old := st.Messages[id]
		if old == nil || old.Receipt.Digest != digest {
			return nil, fault("idempotency_conflict")
		}
		if old.Scope != scope || old.Generation != generation || !st.current(old) {
			return nil, fault("sender_not_allowed")
		}
		return old, nil
	}
	if !st.messageCapacity(false) {
		return nil, fault("capacity")
	}
	message := &Message{Receipt: Receipt{ID: remoteUUID(now), From: from, To: to, Intent: publicTextIntent, Digest: digest, Exp: now.Add(180 * time.Second), Accepted: now, State: "queued"}, Key: key, Scope: scope, Generation: generation, Deliver: "agent"}
	st.copyConversation(message, body, "owner")
	if message.Conversation == nil || !st.current(message) {
		return nil, fault("sender_not_allowed")
	}
	message.Envelope, _ = json.Marshal(OwnerText{V: "knowslink.owner-text.v1", ID: message.Receipt.ID, From: from, To: to, Text: body, Exp: message.Receipt.Exp.UTC().Format(time.RFC3339), Key: key, Authorship: "owner"})
	st.Messages[message.Receipt.ID] = message
	st.Idempotency[from+"/"+key] = message.Receipt.ID
	return message, nil
}

func (st *State) historical(owner, from, to string, now time.Time) []*Message {
	var messages []*Message
	for _, m := range st.Messages {
		c := m.Conversation
		if c != nil && (c.OwnerA == owner || c.OwnerB == owner) && pairID(m.Receipt.From, m.Receipt.To) == pairID(from, to) && now.Before(m.Receipt.Accepted.Add(24*time.Hour)) {
			messages = append(messages, m)
		}
	}
	return messages
}
