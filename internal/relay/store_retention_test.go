// Retained-history benchmarks and boundary checks preserve transport expiry and 24-hour conversation retention.
package relay

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func retentionFixture(count int) (*State, time.Time) {
	now := time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC)
	st := newState()
	st.HumanModel = 1
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("member_%03d", i)
		st.Owners[id] = &Owner{Active: true}
		st.Members[id] = &Member{Owner: id, Issuer: emailIssuer, Active: true, SelfConnect: true}
		for n := 0; n < 2; n++ {
			aid := fmt.Sprintf("agent_%03d_%d", i, n)
			st.Agents[aid] = &Agent{Owner: id, Keys: map[string]*Key{"key1": {}}}
		}
	}
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("message_%05d", i)
		from, to := fmt.Sprintf("agent_%03d_0", i%100), fmt.Sprintf("agent_%03d_1", i%100)
		st.Messages[id] = &Message{
			Receipt: Receipt{ID: id, From: from, To: to, Intent: publicTextIntent, Accepted: now.Add(-time.Hour), Exp: now.Add(-time.Minute), State: "delivered"},
			Key:     id, Kid: "key1", Scope: fmt.Sprintf("self:member_%03d", i%100), Generation: 1, Deliver: "agent",
			Conversation: &ConversationCopy{Text: "synthetic retained conversation", Authorship: "agent"},
		}
		st.Idempotency[from+"/"+id] = id
	}
	return st, now
}

func TestSweepExpiryPreservesRetentionAndRevocation(t *testing.T) {
	for _, state := range []string{"queued", "leased", "delivered", "failed:revoked"} {
		for _, revoked := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/revoked=%t", state, revoked), func(t *testing.T) {
				st, now := retentionFixture(1)
				m := st.Messages["message_00000"]
				m.Receipt.State, m.Receipt.Exp = state, now
				m.Envelope, m.Inbox = json.RawMessage(`{"synthetic":true}`), json.RawMessage(`{}`)
				m.LeaseToken, m.ClaimToken, m.Claimed = "synthetic-lease", "synthetic-claim", true
				st.Agents[m.Receipt.From].Keys[m.Kid].Revoked = revoked
				st.sweep(now)
				want := "failed:expired"
				if state == "delivered" {
					want = state
				}
				if m.Receipt.State != want || m.Completion != "failed:expired" || len(m.Envelope)+len(m.Inbox) != 0 || m.LeaseToken != "" || m.ClaimToken != "" {
					t.Fatalf("expiry did not clear transport/claim: %+v", m)
				}
				if m.Conversation == nil || st.Idempotency[m.Receipt.From+"/"+m.Key] != m.Receipt.ID {
					t.Fatal("expiry discarded retained conversation or replay metadata")
				}
				st.sweep(m.Receipt.Accepted.Add(24 * time.Hour))
				if len(st.Messages)+len(st.Idempotency) != 0 {
					t.Fatal("24-hour boundary retained message or replay metadata")
				}
			})
		}
	}
	st, now := retentionFixture(1)
	m := st.Messages["message_00000"]
	m.Receipt.State, m.Receipt.Exp = "queued", now.Add(time.Nanosecond)
	m.Envelope = json.RawMessage(`{}`)
	st.Agents[m.Receipt.From].Keys[m.Kid].Revoked = true
	st.sweep(now)
	if m.Receipt.State != "failed:revoked" || m.Envelope != nil {
		t.Fatal("unexpired revoked transport survived sweep")
	}
}

func BenchmarkRetainedHistory(b *testing.B) {
	for _, count := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			st, now := retentionFixture(count)
			raw, err := json.Marshal(st)
			if err != nil {
				b.Fatal(err)
			}
			b.Run("sweep", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					st.sweep(now)
				}
			})
			b.Run("transactionCPU", func(b *testing.B) {
				b.ReportAllocs()
				b.ReportMetric(float64(len(raw)), "state_bytes")
				for i := 0; i < b.N; i++ {
					next := newState()
					if err := json.Unmarshal(raw, next); err != nil {
						b.Fatal(err)
					}
					before := next.eventVersion(now)
					next.sweep(now)
					if _, err := json.Marshal(next); err != nil {
						b.Fatal(err)
					}
					if before != next.eventVersion(now) {
						b.Fatal("retained history changed arrival signal")
					}
				}
			})
		})
	}
}
