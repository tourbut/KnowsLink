// Arrival selection and shared stream budgets preserve current authorization and the ordinary pull path.
package relay

import (
	"fmt"
	"testing"
	"time"
)

func TestTextEventSelectionAndAuthorization(t *testing.T) {
	st, keys, now := textFixture()
	before := st.eventVersion(now)
	e := signedText(firstID, "agent_a", "agent_b", "events-selection-01", "", "never in event data", now.Add(time.Minute), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", e, now); err != nil {
		t.Fatal(err)
	}
	if st.eventVersion(now) == before || st.eventPending("agent_b", now) != firstID || st.eventPending("agent_c", now) != "" {
		t.Fatal("missing arrival or cross-agent disclosure")
	}
	queued := st.eventVersion(now)
	st.Rates["idle"] = []time.Time{now}
	if st.eventVersion(now) != queued {
		t.Fatal("idle bookkeeping triggered an event")
	}
	if st.eventPending("agent_b", now.Add(time.Minute)) != "" {
		t.Fatal("expired arrival")
	}
	st.Pairs[pairID("agent_a", "agent_b")].State = "revoked"
	if st.eventPending("agent_b", now) != "" || st.eventVersion(now) == queued {
		t.Fatal("revoked pair arrival")
	}
	before = st.eventVersion(now)
	st.Agents["agent_b"].Keys["key1"].Revoked = true
	if st.eventVersion(now) == before {
		t.Fatal("key revocation must wake listeners for reauthorization")
	}
}

func TestTextEventSharedBudgetAndExpiry(t *testing.T) {
	st := newState()
	now := time.Now()
	for i := 0; i < 16; i++ {
		if err := st.enterStream(fmt.Sprint(i), fmt.Sprint(i/2), now); err != nil {
			t.Fatal(err)
		}
	}
	if st.enterStream("overflow", "other", now) == nil {
		t.Fatal("global stream cap")
	}
	delete(st.Streams, "15")
	if st.enterStream("third", "0", now) == nil {
		t.Fatal("per-agent stream cap")
	}
	st.sweep(now.Add(eventLifetime + 10*time.Second))
	if len(st.Streams) != 0 || st.enterStream("recovered", "0", now) != nil {
		t.Fatal("crashed stream expiry")
	}
}
