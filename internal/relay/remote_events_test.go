// Remote claim checks cover crash recovery, current authority and repeated subscription scan costs.
package relay

import (
	"fmt"
	"testing"
	"time"
)

func TestRemoteExhaustedLeaseWakesNewArrival(t *testing.T) {
	st, keys, now := textFixture()
	id, install, _, _ := remoteFixture(t, st, "agent_b", now)
	st.Pairs[pairID("agent_a", install.Agent)] = &Pair{A: "agent_a", B: install.Agent, State: "active", Generation: 1}
	first := signedText(firstID, "agent_a", install.Agent, "remote-crash-first", "", "first", now.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", first, now); err != nil {
		t.Fatal(err)
	}
	sub := &remoteSubscription{Install: id, Secret: "test-secret", Exp: now.Add(time.Hour)}
	st.Remote.Subscriptions["sub"] = sub
	for attempt := range 3 {
		if st.claimRemote(now.Add(time.Duration(attempt)*16*time.Second)).Delivery == nil {
			t.Fatal("missing delivery before sender crash")
		}
	}
	oldEvent := sub.EventID
	later := now.Add(40 * time.Second)
	second := signedText(secondID, "agent_a", install.Agent, "remote-crash-second", "", "second", later.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", second, later); err != nil {
		t.Fatal(err)
	}
	if work := st.claimRemote(later); work.Delivery != nil || !work.Next.Equal(sub.LeaseUntil) {
		t.Fatal("new arrival bypassed the active lease")
	}
	later = sub.LeaseUntil
	work := st.claimRemote(later)
	if work.Delivery != nil || !work.Next.Equal(later) {
		t.Fatalf("exhausted event must wake new arrival immediately: next=%v want=%v", work.Next, later)
	}
	work = st.claimRemote(later)
	if work.Delivery == nil || work.Delivery.Subscription.EventID == oldEvent || sub.Attempts != 1 {
		t.Fatal("new arrival did not receive a separate event and attempt budget")
	}
	st.finishRemote(*work.Delivery, 200, false, later)
	if work = st.claimRemote(later); work.Delivery != nil || !work.Next.Equal(sub.Exp) {
		t.Fatal("acknowledged arrival should wait without immediate polling")
	}
}

func TestRemoteClaimsRecheckCurrentAuthority(t *testing.T) {
	st, keys, now := textFixture()
	id, install, _, _ := remoteFixture(t, st, "agent_b", now)
	pair := &Pair{A: "agent_a", B: install.Agent, State: "active", Generation: 1}
	st.Pairs[pairID("agent_a", install.Agent)] = pair
	for _, name := range []string{"first", "second"} {
		st.Remote.Subscriptions[name] = &remoteSubscription{Install: id, Secret: "test-secret", Exp: now.Add(time.Hour)}
	}
	if st.claimRemote(now).Delivery != nil {
		t.Fatal("empty queue produced a delivery")
	}
	message := signedText(firstID, "agent_a", install.Agent, "remote-recheck-key", "", "new", now.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", message, now); err != nil {
		t.Fatal(err)
	}
	work := st.claimRemote(now)
	if work.Delivery == nil || work.Delivery.ID != "first" {
		t.Fatal("new arrival reused an empty pending result")
	}
	st.finishRemote(*work.Delivery, 200, false, now)
	pair.State = "revoked"
	if st.claimRemote(now).Delivery != nil {
		t.Fatal("second subscription reused authority from the previous claim")
	}
}

func BenchmarkRemoteClaimSubscriptions(b *testing.B) {
	for _, messages := range []int{100, 10000} {
		for _, subscriptions := range []int{1, 4} {
			b.Run(fmt.Sprintf("messages=%d/subscriptions=%d", messages, subscriptions), func(b *testing.B) {
				st, keys, now := textFixture()
				st.remote().Installs["install"] = &remoteInstall{Owner: "agent_b", Agent: "agent_b", Kid: "key1", Private: keys["agent_b"], Exp: now.Add(time.Hour)}
				for n := range messages {
					st.Messages[fmt.Sprint(n)] = &Message{Receipt: Receipt{To: "agent_b", State: "delivered"}}
				}
				for n := range subscriptions {
					st.Remote.Subscriptions[fmt.Sprint(n)] = &remoteSubscription{Install: "install", Secret: "test-secret", Exp: now.Add(time.Hour)}
				}
				b.ResetTimer()
				for b.Loop() {
					if work := st.claimRemote(now); work.Delivery != nil || work.Next.IsZero() {
						b.Fatal("idle subscription state changed")
					}
				}
			})
		}
	}
}
