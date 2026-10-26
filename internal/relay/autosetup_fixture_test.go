//go:build integration && autosetupfixture

// Opt-in loopback fixture for the Windows auto-setup installer (SAR-AUTOSETUP-001-DEV): synthetic Google member A with
// accepted friend B, the relay API and /install files, plus control endpoints that stand in for the person's browser
// steps (approve or cancel a device attempt, drop one completion answer). No real account, mail or counterpart.
package relay

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutosetupFixture(t *testing.T) {
	addr, dir := os.Getenv("AUTOSETUP_FIXTURE_ADDR"), os.Getenv("AUTOSETUP_FIXTURE_INSTALL_DIR")
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || addr == "" || dir == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	s, _ := identityService(t, messagePool(t))
	s.InstallDir = dir
	var owner, memberA string
	s.mutateState(t, func(st *State) {
		now := time.Now()
		memberA = st.Sessions[hashToken(st.signIn(googleIssuer, "autosetup-a", "autosetup-a@example.test", "", "", now).Session)].Member
		memberB := st.Sessions[hashToken(st.signIn(emailIssuer, "autosetup-b@example.test", "autosetup-b@example.test", "", "", now).Session)].Member
		st.Members[memberA].Name, st.Members[memberB].Name = "합성 회원 A", "합성 회원 B"
		st.Members[memberB].FriendCode = "fixture-code-autosetup-b"
		st.requestFriend(memberA, "fixture-code-autosetup-b", now)
		f := st.Friends[pairID(memberA, memberB)]
		if e := st.decideFriend(memberB, f.ID, "accept", f.Generation, now); e != nil {
			t.Fatal(e)
		}
		owner = st.Members[memberA].Owner
	})
	mutate := func(fn func(*State, time.Time) error) error {
		_, err := s.transaction(context.Background(), func(st *State, now time.Time) (any, error) { return nil, fn(st, now) })
		return err
	}
	var lose atomic.Bool
	handler := s.Handler()
	server := &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fixture/approve", "/fixture/cancel":
			id, name := r.FormValue("connection"), r.FormValue("name")
			err := mutate(func(st *State, now time.Time) error {
				c := st.Connections[id]
				if c == nil || !c.Device {
					return errors.New("unknown device attempt")
				}
				if r.URL.Path == "/fixture/cancel" {
					c.State = "cancelled"
					return nil
				}
				c.Owner, c.State = owner, "prepared"
				if e := st.approveDevice(id, owner, "합성 설치", now); e != nil {
					return e
				}
				st.Agents[c.Agent].Name = name
				return nil
			})
			if err != nil {
				http.Error(w, err.Error(), 409)
			}
			return
		case "/fixture/revoke-all":
			// Separate failure scenarios stay under the member's active-agent and agent-record limits (synthetic data only).
			_ = mutate(func(st *State, _ time.Time) error {
				for id, a := range st.Agents {
					if a.Owner == owner {
						delete(st.Agents, id)
					}
				}
				return nil
			})
			return
		case "/fixture/lose-complete":
			lose.Store(r.FormValue("on") == "1")
			return
		case "/fixture/state":
			counts := map[string]int{}
			_ = mutate(func(st *State, _ time.Time) error {
				for _, a := range st.Agents {
					if a.Owner == owner && !a.Revoked {
						counts["agents"]++
					}
				}
				for _, f := range st.Friends {
					if (f.A == memberA || f.B == memberA) && f.State == "active" {
						counts["friends"]++
					}
				}
				return nil
			})
			_ = json.NewEncoder(w).Encode(counts)
			return
		}
		if r.URL.Path == "/v1/connect/complete" && lose.Swap(false) {
			// One answer is lost: the relay consumes the attempt and the client never receives the credential (F-07).
			handler.ServeHTTP(httptest.NewRecorder(), r)
			if conn, _, err := http.NewResponseController(w).Hijack(); err == nil {
				_ = conn.Close()
			}
			return
		}
		handler.ServeHTTP(w, r)
	})}
	go func() { _ = server.ListenAndServe() }()
	t.Log("autosetup fixture ready on " + addr + "; synthetic accounts only")
	<-time.After(6 * time.Hour)
}
