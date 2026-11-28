// Pre-approval checks pin S1: every send outside one live owner grant is approval_required, and nothing revives a grant.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// grantFixture: one Google member with two self agents (grok, orca) and a friend's agent (peer).
func grantFixture(t *testing.T) (*State, map[string]ed25519.PrivateKey, time.Time) {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	st := newState()
	keys := map[string]ed25519.PrivateKey{}
	for _, o := range []string{"owner_1", "owner_2"} {
		st.Owners[o] = &Owner{Active: true}
		st.Members["m_"+o] = &Member{Owner: o, Active: true, Issuer: googleIssuer, SelfConnect: true, Name: o}
	}
	for id, owner := range map[string]string{"grok": "owner_1", "orca": "owner_1", "peer": "owner_2"} {
		pub, key, _ := ed25519.GenerateKey(rand.Reader)
		keys[id] = key
		st.Agents[id] = &Agent{Name: id, Owner: owner, Keys: map[string]*Key{"key1": {Public: pub, Credential: hashToken(id)}}}
	}
	st.Friends[pairID("m_owner_1", "m_owner_2")] = &Friend{ID: "f1", A: "m_owner_1", B: "m_owner_2", State: "active", Generation: 1}
	return st, keys, now
}

func kindText(t *testing.T, n int, from, to, kind, body string, now time.Time, key ed25519.PrivateKey) *PublicText {
	t.Helper()
	e := &PublicText{V: "knowslink.text.v1", ID: fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012d", n), From: from, To: to, Key: fmt.Sprintf("grant-test-key-%04d", n), Text: body, Kind: kind, Exp: now.Add(170 * time.Second).Format("2006-01-02T15:04:05Z"), Sig: Signature{Alg: "Ed25519", Kid: "key1"}}
	e.Sig.Value = base64.RawURLEncoding.EncodeToString(ed25519.Sign(key, e.signingBytes()))
	return e
}

const answerBody = `{"answer":"1","by":"owner","item":3}`

func TestGrantAdmitsOnlyLiveOwnerForm(t *testing.T) {
	st, keys, now := grantFixture(t)
	n := 0
	send := func(from, to, kind, body string) error {
		n++
		_, err := st.ingestPublicText(from, kindText(t, n, from, to, kind, body, now, keys[from]), now)
		return err
	}
	refused := func(why string, err error) {
		t.Helper()
		if err == nil || err.Error() != "approval_required" {
			t.Fatalf("%s: %v", why, err)
		}
	}
	refused("no grant", send("grok", "orca", "orca.answer.v1", answerBody))
	if err := st.createGrant("owner_1", "grok", "orca", "orca.answer.v1", now); err != nil {
		t.Fatal(err)
	}
	if err := send("grok", "orca", "orca.answer.v1", answerBody); err != nil {
		t.Fatal("granted send", err)
	}
	for id, m := range st.Messages {
		if m.Conversation != nil {
			t.Fatal("grant send kept a conversation copy", id)
		}
	}
	refused("other direction", send("orca", "grok", "orca.answer.v1", answerBody))
	refused("other kind", send("grok", "orca", "orca.task.v1", `{"text":"x"}`))
	refused("unknown kind", send("grok", "orca", "orca.other.v1", answerBody))
	for _, bad := range []string{
		`{"answer":"1","by":"owner","extra":1}`,
		`{"by":"owner","answer":"1"}`, // not canonical
		`{"answer":" ","by":"owner"}`,
		`{"answer":"1","by":"owner","item":0}`,
		`{"answer":"1","by":"owner","item":1.5}`,
		`{"answer":"` + strings.Repeat("a", grantBodyLimit) + `","by":"o"}`,
		`not json`,
	} {
		refused("schema "+bad[:min(len(bad), 40)], send("grok", "orca", "orca.answer.v1", bad))
	}
	// The hourly limit counts accepted sends only and frees after an hour.
	for len(st.grantsOf("grok", "orca.answer.v1")[0].Sent) < grantHourly {
		for _, m := range st.Messages { // orca receives them, so the owner queue limit (ownerQueueLimit) never binds here
			m.Receipt.State = "delivered"
		}
		if err := send("grok", "orca", "orca.answer.v1", answerBody); err != nil {
			t.Fatal("under hourly limit", err)
		}
	}
	refused("hourly", send("grok", "orca", "orca.answer.v1", answerBody))
	now = now.Add(time.Hour)
	st.sweep(now)
	if err := send("grok", "orca", "orca.answer.v1", answerBody); err != nil {
		t.Fatal("next hour", err)
	}
	// Friends and other owners never get a grant or a pre-approved send.
	if err := st.createGrant("owner_1", "grok", "peer", "orca.task.v1", now); err == nil {
		t.Fatal("grant to another owner's agent")
	}
	if err := st.createGrant("owner_2", "grok", "orca", "orca.task.v1", now); err == nil {
		t.Fatal("grant by another owner")
	}
	refused("friend", send("grok", "peer", "orca.task.v1", `{"text":"x"}`))
	// Expiry ends the grant; renewal is a fresh owner grant.
	now = now.Add(grantTTL)
	st.sweep(now)
	refused("expired", send("grok", "orca", "orca.answer.v1", answerBody))
}

func TestGrantEndsOnRevokeAndNeverRevives(t *testing.T) {
	for name, end := range map[string]func(*State, time.Time){
		"grant revoke":  func(st *State, _ time.Time) { _ = st.revokeGrant("owner_1", st.grantsOf("grok", "orca.task.v1")[0].ID) },
		"sender key":    func(st *State, _ time.Time) { st.Agents["grok"].Keys["key1"].Revoked = true },
		"recipient key": func(st *State, _ time.Time) { st.Agents["orca"].Keys["key1"].Revoked = true },
		"agent revoke":  func(st *State, _ time.Time) { st.Agents["orca"].Revoked = true },
		"self connect":  func(st *State, _ time.Time) { st.Members["m_owner_1"].SelfConnect = false },
	} {
		t.Run(name, func(t *testing.T) {
			st, keys, now := grantFixture(t)
			if err := st.createGrant("owner_1", "grok", "orca", "orca.task.v1", now); err != nil {
				t.Fatal(err)
			}
			if st.revokeGrant("owner_2", st.grantsOf("grok", "orca.task.v1")[0].ID) == nil {
				t.Fatal("another owner ended the grant")
			}
			end(st, now)
			st.sweep(now)
			// Undo the cause: a new key, a restored flag or a restart must not bring the old grant back.
			pub, key, _ := ed25519.GenerateKey(rand.Reader)
			st.Agents["grok"].Keys["key1"] = &Key{Public: pub}
			st.Members["m_owner_1"].SelfConnect = true
			keys["grok"] = key
			st = restart(t, st)
			st.sweep(now)
			if len(st.Grants) != 0 {
				t.Fatal("grant survived")
			}
			e := kindText(t, 1, "grok", "orca", "orca.task.v1", `{"text":"x"}`, now, keys["grok"])
			if _, err := st.ingestPublicText("grok", e, now); err == nil || (err.Error() != "approval_required" && err.Error() != "sender_not_allowed") {
				t.Fatal("revived", err)
			}
		})
	}
}

func TestGrantWireKeepsOldSignatures(t *testing.T) {
	st, keys, now := grantFixture(t)
	e := kindText(t, 1, "grok", "orca", "", "plain", now, keys["grok"])
	old, _ := Canonical(map[string]any{"v": e.V, "id": e.ID, "from": e.From, "to": e.To, "text": e.Text, "exp": e.Exp, "idempotency_key": e.Key, "reply_to": e.ReplyTo})
	if string(e.signingBytes()) != "KNOWSLINK-TEXT\x00knowslink.text.v1\x00Ed25519\x00key1\x00"+string(old) {
		t.Fatal("signing bytes without kind changed")
	}
	if strings.Contains(string(e.signingBytes()), "kind") {
		t.Fatal("empty kind signed")
	}
	k := kindText(t, 2, "grok", "orca", "orca.task.v1", `{"text":"x"}`, now, keys["grok"])
	raw, _ := json.Marshal(k)
	parsed, err := parsePublicText(raw)
	if err != nil || parsed.Kind != "orca.task.v1" || parsed.digest() == e.digest() {
		t.Fatal("kind wire", err)
	}
	// The kind is signed: a relabelled text fails verification.
	k.Kind = "orca.answer.v1"
	if err := st.createGrant("owner_1", "grok", "orca", "orca.answer.v1", now); err != nil {
		t.Fatal(err)
	}
	if _, err := st.ingestPublicText("grok", k, now); err == nil || err.Error() != "invalid_signature" {
		t.Fatal("relabelled kind", err)
	}
	for _, bad := range []string{
		strings.Replace(string(raw), `"kind":"orca.task.v1"`, `"kind":""`, 1),
		strings.Replace(string(raw), `"kind":"orca.task.v1"`, `"kind":1`, 1),
		strings.Replace(string(raw), `"kind":"orca.task.v1"`, `"kind":"orca.task.v1","reply_to":"`+firstID+`"`, 1),
	} {
		if _, err := parsePublicText([]byte(bad)); err == nil {
			t.Fatal("invalid kind wire accepted", bad)
		}
	}
}

func (st *State) grantsOf(from, kind string) []*Grant {
	var out []*Grant
	for _, g := range st.Grants {
		if g.From == from && g.Kind == kind {
			out = append(out, g)
		}
	}
	return out
}

// D076: the relay, not the sender, gives 1 hour only between one owner's active self agents; others keep 180 s/300 s.
func TestDeliveryCeilingBySelfAgents(t *testing.T) {
	st, keys, now := grantFixture(t)
	n := 0
	send := func(to string, ttl time.Duration) (time.Duration, error) {
		n++
		e := signedText(fmt.Sprintf("0199a3f2-4c10-7a11-8b33-%012d", n), "grok", to, fmt.Sprintf("ceiling-test-key-%04d", n), "", "ping", now.Add(ttl), keys["grok"])
		r, err := st.ingestPublicText("grok", e, now)
		if err != nil {
			return 0, err
		}
		return r.(Receipt).Exp.Sub(now), nil
	}
	for _, c := range []struct {
		to       string
		ask, got time.Duration
	}{
		{"orca", messageRetention, messageRetention},
		{"peer", messageRetention, peerTextTTL},        // a 1-hour adapter request to a friend is cut, not refused
		{"peer", 170 * time.Second, 170 * time.Second}, // an older client's 170 s stays as sent
	} {
		if got, err := send(c.to, c.ask); err != nil || got != c.got {
			t.Fatalf("%s %v: got %v %v, want %v", c.to, c.ask, got, err, c.got)
		}
	}
	for _, to := range []string{"orca", "peer"} {
		if _, err := send(to, messageRetention+time.Second); err == nil || err.Error() != "ttl_too_long" {
			t.Fatalf("%s past 1 hour: %v", to, err)
		}
	}
	envelope := func(to, intent string) time.Duration {
		return st.envelopeTTL(&Envelope{From: "grok", To: to, Intent: intent})
	}
	if envelope("orca", "relay.request") != messageRetention || envelope("peer", "relay.request") != peerEnvelopeTTL || envelope("orca", "relay.approval.request") != peerEnvelopeTTL {
		t.Fatal("envelope ceilings")
	}
	st.Members["m_owner_1"].SelfConnect = false // no self authority, no 1 hour, whatever the sender asks
	if st.selfAgents("grok", "orca") || st.deliveryTTL("grok", "orca", peerTextTTL) != peerTextTTL {
		t.Fatal("self test without self authority")
	}
}
