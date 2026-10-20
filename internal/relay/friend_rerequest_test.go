package relay

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
)

// U4 CONNECT-VISUAL 결정: re-sending an ended request with the same code reuses one requester row, whether the code
// is a member's, unknown, the requester's own, blocked by its owner, or its owner is at the receive limit.
func TestFriendRerequestReusesOneRow(t *testing.T) {
	type outcome struct {
		rows                    int
		state, label            string
		sameRow, generationUp   bool
		freshDeadline, incoming bool
	}
	view := func(t *testing.T, st *State, id string, at time.Time) []humanFriend {
		t.Helper()
		st.sweep(at)
		v, err := st.humanView(id, &Session{Verified: at}, "friends", httptest.NewRequest("GET", "/home/friends", nil), at)
		if err != nil {
			t.Fatal(err)
		}
		return v["Requests"].([]humanFriend)
	}
	cases := map[string]func(st *State, ids map[string]string, now time.Time) string{
		"existing_expired": func(st *State, ids map[string]string, now time.Time) string {
			return st.Members[ids["b"]].FriendCode
		},
		"existing_cancelled": func(st *State, ids map[string]string, now time.Time) string {
			code := st.Members[ids["b"]].FriendCode
			st.requestFriend(ids["a"], code, now)
			f := st.Friends[pairID(ids["a"], ids["b"])]
			if e := st.decideFriend(ids["a"], f.ID, "cancel", f.Generation, now); e != nil {
				t.Fatal(e)
			}
			return code
		},
		"unknown":  func(st *State, ids map[string]string, now time.Time) string { return "synthetic-unknown-zz99" },
		"own_code": func(st *State, ids map[string]string, now time.Time) string { return st.Members[ids["a"]].FriendCode },
		"blocked_expired": func(st *State, ids map[string]string, now time.Time) string {
			code := st.Members[ids["c"]].FriendCode
			st.requestFriend(ids["a"], code, now)
			f := st.Friends[pairID(ids["a"], ids["c"])]
			if e := st.decideFriend(ids["c"], f.ID, "block", f.Generation, now); e != nil {
				t.Fatal(e)
			}
			return code
		},
		"blocked_cancelled": func(st *State, ids map[string]string, now time.Time) string {
			code := st.Members[ids["c"]].FriendCode
			st.requestFriend(ids["a"], code, now)
			f := st.Friends[pairID(ids["a"], ids["c"])]
			if e := st.decideFriend(ids["c"], f.ID, "block", f.Generation, now); e != nil {
				t.Fatal(e)
			}
			if e := st.decideFriend(ids["a"], f.ID, "cancel", f.Generation, now); e != nil {
				t.Fatal(e)
			}
			return code
		},
		"blocked_then_unblocked": func(st *State, ids map[string]string, now time.Time) string {
			code := st.Members[ids["c"]].FriendCode
			st.requestFriend(ids["a"], code, now)
			f := st.Friends[pairID(ids["a"], ids["c"])]
			if e := st.decideFriend(ids["c"], f.ID, "block", f.Generation, now); e != nil {
				t.Fatal(e)
			}
			if e := st.decideFriend(ids["c"], f.ID, "unblock", f.Generation, now); e != nil {
				t.Fatal(e)
			}
			return code
		},
		"receive_limit": func(st *State, ids map[string]string, now time.Time) string {
			code := st.Members[ids["b"]].FriendCode
			st.requestFriend(ids["a"], code, now)
			for i := 0; i < 10; i++ {
				other := fmt.Sprintf("synthetic-sender-%d", i)
				st.Friends[pairID(other, ids["b"])] = &Friend{ID: fmt.Sprintf("held%02d", i), A: other, B: ids["b"], Inviter: other, State: "pending", Generation: 1, Exp: now.Add(72 * time.Hour)}
			}
			return code
		},
	}
	var want *outcome
	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			st, ids, _, now := humanFixture(t)
			code := setup(st, ids, now)
			if len(view(t, st, ids["a"], now)) == 0 {
				st.requestFriend(ids["a"], code, now)
			}
			later := now.Add(25 * time.Hour)
			before := view(t, st, ids["a"], later)
			if len(before) != 1 || before[0].State != "expired" && before[0].State != "cancelled" {
				t.Fatalf("setup left %+v", before)
			}
			st.requestFriend(ids["a"], code, later)
			after := view(t, st, ids["a"], later)
			got := outcome{rows: len(after), state: after[0].State, label: after[0].CodeLabel(), sameRow: after[0].ID == before[0].ID,
				generationUp: after[0].Generation == before[0].Generation+1, freshDeadline: after[0].Exp == clock(later.Add(24*time.Hour)),
				incoming: after[0].Incoming}
			got.label = got.label[:len(got.label)-4] // the tail differs by code; its shape is the same
			if want == nil {
				want = &got
			}
			if got != *want || got.rows != 1 || got.state != "pending" || !got.sameRow || !got.generationUp || !got.freshDeadline {
				t.Fatalf("requester row differs: %+v, want %+v", got, *want)
			}
			// The blocker and the full recipient see no new request; a repeat changes nothing.
			for _, person := range []string{"b", "c"} {
				for _, row := range view(t, st, ids[person], later) {
					held := name == "receive_limit" || name == "blocked_expired" || name == "blocked_cancelled"
					if held && row.ID == after[0].ID && row.State == "pending" {
						t.Fatal("held request shown to its target")
					}
				}
			}
			st.requestFriend(ids["a"], code, later.Add(time.Hour))
			if again := view(t, st, ids["a"], later.Add(time.Hour)); len(again) != 1 || again[0].Generation != after[0].Generation || again[0].Exp != after[0].Exp {
				t.Fatal("repeat changed a pending request")
			}
		})
	}
}

// A request held by a block still counts toward the requester's ten open requests, like any unanswered one.
func TestFriendBlockedRequestCountsTowardSentLimit(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	st.requestFriend(ids["a"], st.Members[ids["c"]].FriendCode, now)
	f := st.Friends[pairID(ids["a"], ids["c"])]
	if e := st.decideFriend(ids["c"], f.ID, "block", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 10; i++ {
		st.requestFriend(ids["a"], fmt.Sprintf("synthetic-unknown-%04d", i), now)
	}
	if len(st.Friends) != 10 {
		t.Fatalf("sent limit ignored the blocked request: %d rows", len(st.Friends))
	}
}
