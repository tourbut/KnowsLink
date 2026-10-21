//go:build integration && calendarfixture

// Opt-in loopback SAR-SCHEDULE-001 fixture for tester/designer: synthetic Google provider, a connected and an
// unconnected member, friend pairs, pending/expired/processed schedule.query gates and synthetic busy data.
// No real account, token or Google endpoint is touched.
package relay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCalendarBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("CALENDAR_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	const life = 12 * time.Hour
	f := setup(t, messagePool(t))
	g, sign := googleFixture(t)
	p := newSyntheticProvider(t, sign)
	f.s.Google = g
	// Accounts: a requests (agent_a); b answers with a connected calendar (agent_b); c answers without one (agent_c).
	members := map[string]string{}
	f.mutate(func(st *State) {
		now := time.Now()
		for _, name := range []string{"a", "b", "c"} {
			v := st.signIn(googleIssuer, "sub-"+name, name+"@example.test", "", "", now)
			members[name] = st.Sessions[hashToken(v.Session)].Member
			st.Members[members[name]].Name = "합성 회원 " + strings.ToUpper(name)
		}
		st.Agents["agent_a"].Owner, st.Agents["agent_a"].Name = st.Members[members["a"]].Owner, "합성 요청 도우미"
		st.Agents["agent_b"].Owner, st.Agents["agent_b"].Name = st.Members[members["b"]].Owner, "합성 캘린더 도우미"
		st.Agents["agent_c"].Owner, st.Agents["agent_c"].Name = st.Members[members["c"]].Owner, "합성 미연결 도우미"
		// Friend pairs a–b and a–c decide the parents' scope, so they exist before any request is sent.
		for _, peer := range []string{"b", "c"} {
			st.requestFriend(members["a"], st.Members[members[peer]].FriendCode, now)
			for _, fr := range st.Friends {
				if fr.State == "pending" && (fr.A == members[peer] || fr.B == members[peer]) {
					if err := st.decideFriend(members[peer], fr.ID, "accept", fr.Generation, now); err != nil {
						t.Fatal(err)
					}
				}
			}
		}
		st.Calendars = map[string]*Calendar{members["b"]: {Sealed: p.sealedFor("sub-b", members["b"]), State: "connected", Changed: now}}
	})
	start, _ := queryWindow()
	window := func(hours int, minutes int64) map[string]any {
		return map[string]any{"window": map[string]any{"start": start.Format("2006-01-02T15:04:05Z"), "end": start.Add(time.Duration(hours) * time.Hour).Format("2006-01-02T15:04:05Z")}, "granularity_min": minutes}
	}
	specs := []struct {
		name, to string
		body     map[string]any
		state    string
	}{
		{"connected_pending_2h_30", "agent_b", window(2, 30), "pending"},
		{"connected_pending_4h_15", "agent_b", window(4, 15), "pending"},
		{"connected_pending_8h_60", "agent_b", window(8, 60), "pending"},
		{"connected_expired", "agent_b", window(2, 30), "expired"},
		{"connected_processed_denied", "agent_b", window(2, 30), "denied"},
		{"unconnected_pending", "agent_c", window(2, 30), "pending"},
	}
	gates, claims := map[string]string{}, map[string]string{}
	for i, spec := range specs {
		parent, gate := fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", 0x300+i), fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", 0x400+i)
		receipt := f.send(wire(t, f.private["agent_a"], parent, "agent_a", spec.to, "schedule.query", fmt.Sprintf("calendar-fixture-p%02d", i), "", spec.body, time.Now().Add(180*time.Second)), "agent_a", "", 200)
		lease := f.call("POST", "/v1/pull", f.tokens[spec.to], map[string]any{}, 200)
		payload := map[string]any{"id": parent, "token": lease["lease_token"]}
		f.call("POST", "/v1/persist", f.tokens[spec.to], payload, 200)
		f.call("POST", "/v1/ack", f.tokens[spec.to], payload, 200)
		claim := f.call("POST", "/v1/claim", f.tokens[spec.to], map[string]any{"id": parent}, 200)["claim"].(string)
		f.send(wire(t, f.private[spec.to], gate, spec.to, spec.to, "relay.approval.request", fmt.Sprintf("calendar-fixture-g%02d", i), parent, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), spec.to, claim, 200)
		gates[spec.name] = gate
		// Claim capacity is global (4 live claims): park each claim and restore all of them after the loop.
		f.mutate(func(st *State) { st.Messages[parent].ClaimToken = "" })
		claims[gate] = claim
	}
	f.mutate(func(st *State) {
		now := time.Now()
		// Deadlines stretch for the fixture lifetime only; signatures, digests and generations stay as sent.
		for name, gate := range gates {
			g := st.Gates[gate]
			g.Exp, st.Messages[g.Parent].Receipt.Exp, st.Messages[gate].Receipt.Exp = now.Add(life), now.Add(life), now.Add(life)
			st.Messages[g.Parent].ClaimToken = claims[gate]
			switch {
			case strings.HasSuffix(name, "_expired"):
				g.Exp = now.Add(-time.Second)
			case strings.HasSuffix(name, "_denied"):
				g.State = "denied"
			}
		}
	})
	handler := f.s.Handler()
	provider := p.handler()
	allowedGET := map[string]bool{"/": true, "/home": true, "/home/legacy": true, "/home/friends": true, "/home/account": true, "/auth/google/callback": true, "/fixture/google/auth": true}
	allowedPOST := map[string]bool{"/home/calendar/connect": true, "/home/calendar/disconnect": true, "/auth/reauth": true}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		sameOrigin := r.Header.Get("Sec-Fetch-Site") != "cross-site" && (r.Header.Get("Origin") == "" || r.Header.Get("Origin") == "http://"+r.Host)
		if name, ok := strings.CutPrefix(r.URL.Path, "/fixture/as/"); ok {
			if r.Method != "POST" || !sameOrigin || members[name] == "" {
				http.Error(w, "fixture origin, method or account refused", 403)
				return
			}
			var token string
			f.mutate(func(st *State) {
				token = st.signIn(googleIssuer, "sub-"+name, name+"@example.test", "", "", time.Now()).Session
			})
			setCookie(w, sessionCookie, token, int(sessionAbsolute/time.Second))
			// The synthetic consent page decides for this account's Google subject.
			http.SetCookie(w, &http.Cookie{Name: fixtureAccountCookie, Value: name, Path: "/", MaxAge: int(life / time.Second), HttpOnly: true, SameSite: http.SameSiteLaxMode})
			w.WriteHeader(204)
			return
		}
		// Read-only view of what the requester's agent would receive: gate state, lookup outcome and disclosure.
		if gate, ok := strings.CutPrefix(r.URL.Path, "/fixture/gate/"); ok && r.Method == "GET" {
			var view map[string]any
			f.mutate(func(st *State) {
				g := st.Gates[gate]
				if g == nil {
					return
				}
				view = map[string]any{"state": g.State, "agent_sees": lookupVisible(g), "consumed": g.Consumed}
				if g.Lookup != nil {
					view["lookup"] = g.Lookup.State
					if len(g.Lookup.Disclosure) > 0 {
						view["disclosure"] = g.Lookup.Disclosure
					}
				}
			})
			if view == nil {
				http.Error(w, "unknown gate", 404)
				return
			}
			writeJSON(w, 200, view)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/fixture/google/") {
			if !sameOrigin {
				http.Error(w, "fixture origin refused", 403)
				return
			}
			provider.ServeHTTP(w, r)
			return
		}
		gatePath := strings.HasPrefix(r.URL.Path, "/home/gates/")
		if !(r.Method == "GET" && (allowedGET[r.URL.Path] || gatePath)) && !(r.Method == "POST" && (gatePath || allowedPOST[r.URL.Path])) {
			http.Error(w, "outside calendar fixture scope", 403)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	// The relay reaches the provider over loopback HTTP; the callback stays on the fixture origin so cookies return.
	p.install(g, server.URL, server.Client())
	data, _ := json.Marshal(map[string]any{"origin": server.URL, "gates": gates, "accounts": map[string]string{"a": "requester", "b": "connected", "c": "unconnected"},
		"provider_toggles": []string{"/fixture/google/toggle/revoke-grant", "/fixture/google/toggle/outage"}, "gate_view": "/fixture/gate/{id}", "expires_at": time.Now().Add(life).UTC().Format(time.RFC3339)})
	if err := os.WriteFile(os.Getenv("CALENDAR_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("calendar fixture ready; synthetic provider and accounts only")
	<-time.After(life)
}
