// Friend re-request rows (U4 CONNECT-VISUAL 결정): one requester row per code whatever the target, held requests unseen.
package relay

import (
	"fmt"
	"net/http/httptest"
	"reflect"
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
		// F-R1: the full receiver once sent the earlier request, so it must keep seeing its own sent record.
		"receive_limit_reversed": func(st *State, ids map[string]string, now time.Time) string {
			st.requestFriend(ids["b"], st.Members[ids["a"]].FriendCode, now)
			for i := 0; i < 10; i++ {
				other := fmt.Sprintf("synthetic-sender-%d", i)
				st.Friends[pairID(other, ids["b"])] = &Friend{ID: fmt.Sprintf("held%02d", i), A: other, B: ids["b"], Inviter: other, State: "pending", Generation: 1, Exp: now.Add(72 * time.Hour)}
			}
			return st.Members[ids["b"]].FriendCode
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
			held := name == "receive_limit" || name == "receive_limit_reversed" || name == "blocked_expired" || name == "blocked_cancelled"
			// The rendered friends pages are what a target observes; the per-view form attempt is random and dropped.
			targets := func(at time.Time) []string {
				pages := []string{}
				for _, person := range []string{"b", "c"} {
					st.sweep(at)
					v, err := st.humanView(ids[person], &Session{Verified: at}, "friends", httptest.NewRequest("GET", "/home/friends", nil), at)
					if err != nil {
						t.Fatal(err)
					}
					v["Attempt"] = ""
					pages = append(pages, renderHumanPage(t, "human-friends", v))
				}
				return pages
			}
			unseen := targets(later)
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
			// The blocker and the full receiver see exactly what they saw before, generation included, and keep seeing it
			// after the requester cancels; a repeat changes nothing.
			if held && !reflect.DeepEqual(targets(later), unseen) {
				t.Fatal("held request changed its target's page")
			}
			st.requestFriend(ids["a"], code, later.Add(time.Hour))
			if again := view(t, st, ids["a"], later.Add(time.Hour)); len(again) != 1 || again[0].Generation != after[0].Generation || again[0].Exp != after[0].Exp {
				t.Fatal("repeat changed a pending request")
			}
			if e := st.decideFriend(ids["a"], after[0].ID, "cancel", after[0].Generation, later); e != nil {
				t.Fatal(e)
			}
			if held && !reflect.DeepEqual(targets(later), unseen) {
				t.Fatal("cancelling a held request showed it to its target")
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

// A placeholder row stored by the old receive-limit path stays the code's only requester row, and a stale blocker
// page can still unblock after a held re-request.
func TestFriendLegacyPlaceholderAndStaleUnblock(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	code := st.Members[ids["b"]].FriendCode
	legacy := &Friend{ID: "legacy-row", A: ids["a"], B: "undisclosed:" + hashToken(code), Inviter: ids["a"], State: "expired", Generation: 3, Exp: now, CodeTail: codeTail(code)}
	st.Friends[pairID(ids["a"], legacy.B)] = legacy
	st.requestFriend(ids["a"], code, now)
	if f := st.Friends[pairID(ids["a"], ids["b"])]; len(st.Friends) != 1 || f == nil || f.ID != "legacy-row" || f.Generation != 4 || f.State != "pending" {
		t.Fatalf("legacy placeholder not reused: %+v", st.Friends)
	}
	f := st.Friends[pairID(ids["a"], ids["b"])]
	if e := st.decideFriend(ids["b"], f.ID, "block", f.Generation, now); e != nil {
		t.Fatal(e)
	}
	stale := f.Generation
	st.requestFriend(ids["a"], code, now.Add(25*time.Hour))
	if f.Generation == stale {
		t.Fatal("held re-request did not advance the requester's generation")
	}
	if e := st.decideFriend(ids["b"], f.ID, "unblock", 0, now.Add(25*time.Hour)); e != nil || f.BlockB {
		t.Fatal("blocker could not unblock from a page without the generation", e)
	}
	if e := st.decideFriend(ids["a"], f.ID, "unblock", 0, now.Add(25*time.Hour)); e == nil {
		t.Fatal("requester unblocked a block they did not set")
	}
}

// F-R2: a receiver who posts block on a request held from them gets the same result as blocking a visible pending
// request: the requester's live request stays pending until its deadline instead of turning into an ended row.
func TestFriendBlockOnHeldRequestKeepsRequesterPending(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	code := st.Members[ids["b"]].FriendCode
	st.requestFriend(ids["a"], code, now)
	for i := 0; i < 10; i++ {
		other := fmt.Sprintf("synthetic-sender-%d", i)
		st.Friends[pairID(other, ids["b"])] = &Friend{ID: fmt.Sprintf("held%02d", i), A: other, B: ids["b"], Inviter: other, State: "pending", Generation: 1, Exp: now.Add(72 * time.Hour)}
	}
	later := now.Add(25 * time.Hour)
	st.sweep(later)
	st.requestFriend(ids["a"], code, later)
	f := st.Friends[pairID(ids["a"], ids["b"])]
	if f.State != "blocked-pending" || f.Held != "expired" {
		t.Fatalf("setup did not hold the request: %+v", f)
	}
	if e := st.decideFriend(ids["b"], f.ID, "block", f.Generation, later); e != nil {
		t.Fatal(e)
	}
	v, err := st.humanView(ids["a"], &Session{Verified: later}, "friends", httptest.NewRequest("GET", "/home/friends", nil), later)
	if err != nil {
		t.Fatal(err)
	}
	if rows := v["Requests"].([]humanFriend); len(rows) != 1 || rows[0].State != "pending" {
		t.Fatalf("requester saw the receiver's block: %+v", rows)
	}
}

// F-R3: data stored before the one-row rule can hold an ended pair row blocked by the target and an open placeholder
// for the same code; a repeat then changes nothing, like a repeat of an unknown code.
func TestFriendRepeatWithOpenLegacyPlaceholderIsNoop(t *testing.T) {
	st, ids, _, now := humanFixture(t)
	code := st.Members[ids["c"]].FriendCode
	pair := &Friend{ID: "pair-row", A: ids["a"], B: ids["c"], Inviter: ids["a"], State: "removed", Generation: 2, Exp: now, BlockB: true, SuppressLegacy: true}
	st.Friends[pairID(ids["a"], ids["c"])] = pair
	placeholder := "undisclosed:" + hashToken(code)
	st.Friends[pairID(ids["a"], placeholder)] = &Friend{ID: "legacy-open", A: ids["a"], B: placeholder, Inviter: ids["a"], State: "pending", Generation: 1, Exp: now.Add(time.Hour), CodeTail: codeTail(code)}
	st.requestFriend(ids["a"], code, now)
	if pair.State != "removed" || pair.Generation != 2 {
		t.Fatalf("repeat reopened the pair row beside the open placeholder: %+v", pair)
	}
}
