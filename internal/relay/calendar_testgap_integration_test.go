//go:build integration

// SAR-SCHEDULE-001-TESTGAP: the three paths the independent QA could not run — schedule.query outside a relationship
// (SC-04), friend removal and block right after a request (SC-06), and a failed re-consent of the member's own
// connected account (L3).
package relay

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

func connectCalendar(f *fixture, p *syntheticProvider, member string) []byte {
	sealed := p.sealedFor("sub-bob", member)
	f.mutate(func(st *State) { st.Calendars = map[string]*Calendar{member: {Sealed: sealed, State: "connected"}} })
	return sealed
}

// SC-04: a stranger's query and a query of a revoked generation are refused before any lookup.
func TestCalendarQueryOutsideRelationship(t *testing.T) {
	t.Run("stranger_refused_at_accept", func(t *testing.T) {
		f, p, _, member := calendarSetup(t)
		connectCalendar(f, p, member)
		f.send(wire(t, f.private["agent_c"], firstID, "agent_c", "agent_b", "schedule.query", "testgap-stranger-01", "", queryBody(), time.Now().Add(time.Minute)), "agent_c", "", 403)
		f.mutate(func(st *State) {
			if st.Messages[firstID] != nil || len(st.Gates) != 0 {
				t.Fatal("stranger query stored")
			}
		})
		if calls, _, _, _ := p.snapshot(); calls != 0 {
			t.Fatal("lookup for a stranger", calls)
		}
	})
	t.Run("revoked_generation_not_delivered", func(t *testing.T) {
		f, p, _, member := calendarSetup(t)
		connectCalendar(f, p, member)
		f.send(wire(t, f.private["agent_a"], firstID, "agent_a", "agent_b", "schedule.query", "testgap-old-gen-01", "", queryBody(), time.Now().Add(time.Minute)), "agent_a", "", 200)
		f.call("POST", "/v1/unpair", f.tokens["agent_a_owner"], map[string]any{"agent": "agent_a", "target": "agent_b"}, 200)
		// A human member cannot re-pair through /v1/invites; the next generation is what a new acceptance leaves behind.
		f.mutate(func(st *State) {
			pair := st.Pairs[pairID("agent_a", "agent_b")]
			pair.State, pair.Generation = "active", pair.Generation+1
		})
		f.mutate(func(st *State) {
			m := st.Messages[firstID]
			if pair := st.Pairs[pairID("agent_a", "agent_b")]; pair.State != "active" || pair.Generation == m.Generation || m.Receipt.State != "failed:revoked" || len(m.Envelope) != 0 {
				t.Fatalf("old generation query: pair=%+v receipt=%+v", pair, m.Receipt)
			}
		})
		f.call("POST", "/v1/claim", f.tokens["agent_b"], map[string]any{"id": firstID}, 403)
		f.call("GET", "/v1/receipts/"+firstID, f.tokens["agent_a"], nil, 403)
		if calls, _, _, _ := p.snapshot(); calls != 0 {
			t.Fatal("lookup for a revoked generation", calls)
		}
	})
}

// SC-06: removing or blocking the friend stops a pending approval and hides an approved result not yet consumed.
func TestCalendarFriendEndedStopsQuery(t *testing.T) {
	for _, decision := range []string{"remove", "block"} {
		for _, approved := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s_approved_%v", decision, approved), func(t *testing.T) {
				f, p, bob, member := calendarSetup(t)
				connectCalendar(f, p, member)
				var friend Friend
				f.mutate(func(st *State) {
					v := st.signIn(googleIssuer, "sub-alice", "alice@example.test", "", "", time.Now())
					alice := st.Sessions[hashToken(v.Session)].Member
					st.Members[alice].Name = "앨리스"
					st.Agents["agent_a"].Owner = st.Members[alice].Owner
					friend = Friend{ID: randomToken()[:22], A: alice, B: member, Inviter: alice, NameA: "앨리스", NameB: "밥", State: "active", Generation: 1, Exp: time.Now().Add(time.Hour)}
					if st.Friends == nil {
						st.Friends = map[string]*Friend{}
					}
					st.Friends[pairID(alice, member)] = &friend
				})
				claim := f.gateFor(t, queryBody(), firstID, gateID, "testgap-friend-"+decision)
				token := csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())
				lookups := 0
				if approved {
					expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 303)
					if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["lookup"] != "done" {
						t.Fatal("lookup not done before the friendship ended", g)
					}
					lookups = 1
				}
				expect(t, bob.do("POST", "/home/friends/decision", url.Values{"friend": {friend.ID}, "generation": {"1"}, "decision": {decision}}), 303)
				if !approved {
					expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 409)
				}
				f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
				f.mutate(func(st *State) {
					if g := st.Gates[gateID]; g.State != "revoked" || (g.Lookup != nil && len(g.Lookup.Disclosure) != 0) {
						t.Fatalf("gate after %s: %s %+v", decision, g.State, g.Lookup)
					}
				})
				// The requester's agent gets no result and can no longer read the request's receipt.
				f.send(wire(t, f.private["agent_b"], secondID, "agent_b", "agent_a", "relay.result", "testgap-result-"+decision, firstID, map[string]any{"status": "denied"}, time.Now().Add(time.Minute)), "agent_b", claim, 403)
				f.call("GET", "/v1/receipts/"+firstID, f.tokens["agent_a"], nil, 403)
				if calls, _, _, _ := p.snapshot(); calls != lookups {
					t.Fatal("lookup count", calls)
				}
			})
		}
	}
}

// L3: a refused re-consent of the member's own connected account shares that grant, so it is not revoked; the
// other-account refusal that is revoked is covered by TestCalendarConnection.
func TestCalendarOwnAccountRefusalKeepsGrant(t *testing.T) {
	f, p, bob, member := calendarSetup(t)
	sealed := connectCalendar(f, p, member)
	w := bob.do("POST", "/home/calendar/connect", nil)
	expect(t, w, 303, "/fixture/google/auth")
	// The session the consent started from ends before Google returns, so signIn refuses the same account.
	f.mutate(func(st *State) { delete(st.Sessions, hashToken(bob.cookies[sessionCookie].Value)) })
	expect(t, bob.do("GET", strings.TrimPrefix(p.grant(w.Header().Get("Location"), "sub-bob", "openid email "+calendarScope), "https://example.com"), nil), 403, "이 계정으로 로그인하거나 재확인할 수 없습니다")
	if _, _, _, revoked := p.snapshot(); len(revoked) != 0 {
		t.Fatalf("own account grant revoked: %v", revoked)
	}
	f.mutate(func(st *State) {
		if c := st.Calendars[member]; c == nil || c.State != "connected" || string(c.Sealed) != string(sealed) {
			t.Fatal("existing connection changed by a refused re-consent")
		}
	})
}
