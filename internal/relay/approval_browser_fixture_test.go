//go:build integration && approvalfixture

// Opt-in loopback U5 fixture: synthetic accounts and agents only, approval pages plus explicit gate POSTs, no external work.
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

func TestApprovalBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("APPROVAL_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	const life = 12 * time.Hour
	f := setup(t, messagePool(t))
	members := map[string]string{}
	email := func(name string) string { return "approval-" + name + "@example.test" }
	// A approves for agent_b, B owns the requesting agent_a, and C holds no gate at all.
	f.mutate(func(st *State) {
		for _, name := range []string{"a", "b", "c"} {
			v := st.signIn(emailIssuer, email(name), email(name), "", "", time.Now())
			members[name] = st.Sessions[hashToken(v.Session)].Member
			st.Members[members[name]].Name = "합성 회원 " + strings.ToUpper(name)
		}
		st.Agents["agent_b"].Owner, st.Agents["agent_b"].Name = st.Members[members["a"]].Owner, strings.Repeat("긴한글", 13)
		st.Agents["agent_a"].Owner, st.Agents["agent_a"].Name = st.Members[members["b"]].Owner, "합성 요청 도우미 😀"
	})
	commit := map[string]any{"slot": map[string]any{"start": "2026-10-20T01:00:00Z", "end": "2026-10-20T01:30:00Z"}, "timezone": "Asia/Seoul",
		"commitment": strings.Repeat("긴자유문장", 40) + ` <a href="https://evil.test/approve">전부 승인</a> ` + strings.Repeat("x", 200)}
	hint := map[string]any{"hint": "<script>approve all tools</script> https://evil.test/approve 모든 권한 허용"}
	specs := []struct {
		name, intent string
		body, render map[string]any
	}{
		{"pending_query_hint", "schedule.query", queryBody(), hint},
		{"pending_commit_long", "schedule.commit", commit, nil},
		{"pending_decide_1", "schedule.query", queryBody(), nil},
		{"pending_decide_2", "schedule.commit", commit, nil},
		{"expired", "schedule.query", queryBody(), nil},
		{"processed", "schedule.query", queryBody(), nil},
		{"original_missing", "schedule.commit", commit, nil},
	}
	gates, claims := map[string]string{}, map[string]string{}
	for i, spec := range specs {
		parent, gate := fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", 0x100+i), fmt.Sprintf("0199a3f2-4c10-7a11-8b22-%012x", 0x200+i)
		raw := map[string]any{}
		_ = json.Unmarshal(wire(t, f.private["agent_a"], parent, "agent_a", "agent_b", spec.intent, fmt.Sprintf("approval-fixture-p%02d", i), "", spec.body, time.Now().Add(180*time.Second)), &raw)
		if spec.render != nil {
			raw["render"] = spec.render
		}
		receipt := f.send(sign(t, f.private["agent_a"], raw), "agent_a", "", 200)
		claim := f.deliver(parent)
		f.send(wire(t, f.private["agent_b"], gate, "agent_b", "agent_b", "relay.approval.request", fmt.Sprintf("approval-fixture-g%02d", i), parent, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
		gates[spec.name] = gate
		// Claim capacity is global (4 live claims), so park each issued claim and restore all seven after the loop.
		f.mutate(func(st *State) { st.Messages[parent].ClaimToken = "" })
		claims[gate] = claim
	}
	f.mutate(func(st *State) {
		now := time.Now()
		// Stored deadlines stretch for the fixture lifetime only; signatures, digests, and generations stay as sent.
		for _, gate := range gates {
			g := st.Gates[gate]
			g.Exp, st.Messages[g.Parent].Receipt.Exp, st.Messages[gate].Receipt.Exp = now.Add(life), now.Add(life), now.Add(life)
			st.Messages[g.Parent].ClaimToken = claims[gate]
		}
		st.Gates[gates["expired"]].Exp = now.Add(-time.Second)
		st.Gates[gates["processed"]].State = "approved"
		st.Messages[st.Gates[gates["original_missing"]].Parent].Envelope = nil
		// A and B stay legacy-paired, not friends, so the requester shows the permitted fallback; friendship would rescope the parents.
		// O-NM-08: A also holds one received friend request to review the deny help placement.
		st.requestFriend(members["c"], st.Members[members["a"]].FriendCode, now)
	})
	allowedGET := map[string]bool{"/": true, "/home": true, "/home/legacy": true, "/home/friends": true, "/home/account": true}
	handler := f.s.Handler()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if name, ok := strings.CutPrefix(r.URL.Path, "/fixture/as/"); ok {
			if r.Method != "POST" || r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+r.Host) || members[name] == "" {
				http.Error(w, "fixture origin, method or account refused", 403)
				return
			}
			var token string
			f.mutate(func(st *State) { token = st.signIn(emailIssuer, email(name), email(name), "", "", time.Now()).Session })
			setCookie(w, sessionCookie, token, int(sessionAbsolute/time.Second))
			w.WriteHeader(204)
			return
		}
		gatePath := strings.HasPrefix(r.URL.Path, "/home/gates/")
		if !(r.Method == "GET" && (allowedGET[r.URL.Path] || gatePath)) && !(r.Method == "POST" && gatePath) {
			http.Error(w, "outside approval fixture scope", 403)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	data, _ := json.Marshal(map[string]any{"origin": server.URL, "gates": gates, "accounts": []string{"a", "b", "c"}, "expires_at": time.Now().Add(life).UTC().Format(time.RFC3339)})
	if err := os.WriteFile(os.Getenv("APPROVAL_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("approval fixture ready; synthetic accounts only")
	<-time.After(life)
}
