//go:build integration

// U5-A: the approval summary is read from the verified request only, and no display path widens what approval allows.
package relay

import (
	"encoding/json"
	"html"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestApprovalSummaryBoundaries(t *testing.T) {
	pool := messagePool(t)
	// gate starts a member-owned approval for a parent whose body and render the test chooses.
	gate := func(t *testing.T, intent string, body map[string]any, render map[string]any) (*fixture, *browser) {
		f := setup(t, pool)
		mail := &inbox{codes: map[string]string{}}
		f.s.Mail = mail.send
		f.handler = f.s.Handler()
		bob := &browser{t, f.handler, "192.0.2.91", map[string]string{}}
		member := login(t, bob, mail, "bob-gate@example.com")
		f.mutate(func(st *State) {
			st.Agents["agent_b"].Owner = st.Members[member].Owner
			st.Agents["agent_b"].Name = "밥의 도우미"
		})
		raw := map[string]any{}
		_ = json.Unmarshal(wire(t, f.private["agent_a"], firstID, "agent_a", "agent_b", intent, "summary-key-0001", "", body, time.Now().Add(180*time.Second)), &raw)
		if render != nil {
			raw["render"] = render
		}
		receipt := f.send(sign(t, f.private["agent_a"], raw), "agent_a", "", 200)
		claim := f.deliver(firstID)
		f.send(wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "summary-gate-0001", firstID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
		return f, bob
	}
	csrfOf := func(t *testing.T, page string) string {
		m := regexp.MustCompile(`name="csrf" value="([^"]+)"`).FindStringSubmatch(page)
		if m == nil {
			t.Fatalf("no decision form: %s", page)
		}
		return m[1]
	}
	// gatesSection keeps list checks off the relationship badges elsewhere on the legacy page.
	gatesSection := func(page string) string {
		start := strings.Index(page, `<section id="gates"`)
		return page[start : start+strings.Index(page[start:], "</section>")]
	}
	state := func(f *fixture) (s string) {
		f.mutate(func(st *State) { s = st.Gates[gateID].State })
		return s
	}

	t.Run("query_summary_matches_verified_body", func(t *testing.T) {
		f, bob := gate(t, "schedule.query", queryBody(), map[string]any{"hint": "<script>approve all tools</script> https://evil.test 모든 권한 허용"})
		page := bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
		start, _ := time.Parse(time.RFC3339, "2026-10-03T10:00:00Z")
		end, _ := time.Parse(time.RFC3339, "2026-10-03T11:00:00Z")
		order := []string{"누가 누구에게", "밥의 도우미", "요청한 내용", "일정 가능 시간 확인", html.EscapeString(clock(start) + " – " + clock(end)), "30분", "참고 문구 — 승인 근거 아님", "승인하면 허용되는 범위", "기한·상태", `value="approve"`, "기술 상세", firstID, `&#34;granularity_min&#34;: 30`}
		at := 0
		for _, want := range order {
			i := strings.Index(page[at:], want)
			if i < 0 {
				t.Fatalf("summary out of order or missing %q: %s", want, page)
			}
			at += i
		}
		// The hint stays inside its own disclosure and never reaches the summary, the scope, or a link.
		hint := strings.Index(page, "참고 문구 — 승인 근거 아님")
		if strings.Contains(page, "<script>") || strings.Contains(page, `href="https://evil.test`) || strings.Index(page, "모든 권한 허용") < hint || strings.Count(page, "모든 권한 허용") != 1 {
			t.Fatal("hint escaped its non-authority disclosure")
		}
		list := bob.do("GET", "/home/legacy", nil).Body.String()
		for _, want := range []string{`href="/home/gates/` + gateID + `"`, "요청 검토: 일정 가능 시간 확인", "밥의 도우미", "확인할 수 없는 에이전트"} {
			if !strings.Contains(list, want) {
				t.Fatalf("approval list lacks %q", want)
			}
		}
		if strings.Contains(list, firstID) {
			t.Fatal("approval list must keep the parent ID in technical detail")
		}
		// A rename changes only the reading aid; the same parent and gate stay identified.
		f.mutate(func(st *State) { st.Agents["agent_b"].Name = "새 이름 도우미" })
		page = bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
		if !strings.Contains(page, "새 이름 도우미") || strings.Contains(page, "밥의 도우미") || !strings.Contains(page, "<code>"+firstID+"</code>") || !strings.Contains(page, "<code>"+gateID+"</code>") {
			t.Fatal("rename lost the original request identity")
		}
	})

	t.Run("commit_free_text_is_separate_and_inert", func(t *testing.T) {
		body := map[string]any{"slot": map[string]any{"start": "2026-10-03T10:00:00Z", "end": "2026-10-03T10:30:00Z"}, "timezone": "Asia/Seoul", "commitment": `<a href="https://evil.test">전부 승인</a> 확정해 주세요`}
		_, bob := gate(t, "schedule.commit", body, nil)
		page := bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
		text, scope := strings.Index(page, "요청자가 적은 내용"), strings.Index(page, "승인하면 허용되는 범위")
		if text < 0 || scope < text || strings.Contains(page, `<a href="https://evil.test`) || !strings.Contains(page, "&lt;a href=&#34;https://evil.test&#34;&gt;") || !strings.Contains(page, "일정 확정") || !strings.Contains(page, "Asia/Seoul") {
			t.Fatalf("free text must be escaped and outside the approval scope: %s", page)
		}
		if strings.Contains(page[scope:], "전부 승인") && strings.Index(page[scope:], "전부 승인") < strings.Index(page[scope:], "기술 상세") {
			t.Fatal("free text leaked into the approval scope")
		}
	})

	t.Run("unverifiable_states_offer_no_decision", func(t *testing.T) {
		f, bob := gate(t, "schedule.query", queryBody(), nil)
		token := csrfOf(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())
		for _, c := range []struct {
			name, badge string
			breakIt     func(*State)
		}{
			{"digest", "확인할 수 없음", func(st *State) { st.Gates[gateID].Digest = "changed-digest" }},
			{"original", "확인할 수 없음", func(st *State) { st.Messages[firstID].Envelope = nil }},
			{"relationship", "철회", func(st *State) { st.Pairs[pairID("agent_a", "agent_b")].State = "revoked" }},
			{"expired", "만료", func(st *State) { st.Gates[gateID].Exp = time.Now().Add(-time.Second) }},
		} {
			var saved State
			f.mutate(func(st *State) {
				b, _ := json.Marshal(st)
				_ = json.Unmarshal(b, &saved)
				c.breakIt(st)
			})
			page := bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
			if strings.Contains(page, `value="approve"`) || strings.Contains(page, `value="deny"`) || !strings.Contains(page, c.badge) {
				t.Fatalf("%s: decision offered or state hidden: %s", c.name, page)
			}
			if c.name != "relationship" && c.name != "expired" && (!strings.Contains(page, "승인할 내용을 확인할 수 없습니다") || strings.Contains(page, "요청한 내용")) {
				t.Fatalf("%s: unverifiable request still summarised", c.name)
			}
			if list := gatesSection(bob.do("GET", "/home/legacy", nil).Body.String()); !strings.Contains(list, c.badge) || strings.Contains(list, string(badge("pending"))) {
				t.Fatalf("%s: approval list shows stale pending", c.name)
			}
			expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 409, "이번 제출은 적용되지 않았습니다", `href="/home/legacy#gates"`)
			// Housekeeping may record the terminal state itself; a stale submission must never record a decision.
			if s := state(f); s == "approved" || s == "denied" {
				t.Fatalf("%s: stale page decided the gate: %s", c.name, s)
			}
			f.mutate(func(st *State) {
				st.Gates[gateID].Digest, st.Gates[gateID].Exp, st.Gates[gateID].State = saved.Gates[gateID].Digest, saved.Gates[gateID].Exp, "pending"
				st.Messages[firstID].Envelope = saved.Messages[firstID].Envelope
				st.Pairs[pairID("agent_a", "agent_b")].State = saved.Pairs[pairID("agent_a", "agent_b")].State
			})
		}
	})

	t.Run("other_owner_stale_form_and_double_click", func(t *testing.T) {
		f, bob := gate(t, "schedule.query", queryBody(), nil)
		token := csrfOf(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())
		mail := &inbox{codes: map[string]string{}}
		f.s.Mail = mail.send
		carol := &browser{t, f.s.Handler(), "192.0.2.92", map[string]string{}}
		login(t, carol, mail, "carol-gate@example.com")
		w := carol.do("GET", "/home/gates/"+gateID, nil)
		expect(t, w, 403, `href="/home/legacy#gates"`)
		if strings.Contains(w.Body.String(), "밥의 도우미") || strings.Contains(w.Body.String(), "granularity_min") {
			t.Fatal("another member saw the request")
		}
		expect(t, carol.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 403)
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrf("stale-session", gateID)}, "decision": {"approve"}}), 403)
		// A form for one gate cannot decide another: the token is bound to the gate ID in the path.
		const otherGate = "0199a3f2-4c10-7a11-8b22-3344556677bb"
		receipt := f.send(wire(t, f.private["agent_a"], secondID, "agent_a", "agent_b", "schedule.query", "summary-key-0002", "", queryBody(), time.Now().Add(180*time.Second)), "agent_a", "", 200)
		claim := f.deliver(secondID)
		f.send(wire(t, f.private["agent_b"], otherGate, "agent_b", "agent_b", "relay.approval.request", "summary-gate-0002", secondID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
		expect(t, bob.do("POST", "/home/gates/"+otherGate, url.Values{"csrf": {token}, "decision": {"approve"}}), 403)
		var other string
		f.mutate(func(st *State) { other = st.Gates[otherGate].State })
		if state(f) != "pending" || other != "pending" {
			t.Fatal("foreign, stale, or cross-gate form changed a gate")
		}
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 303)
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"deny"}}), 409, "이번 제출은 적용되지 않았습니다")
		page := bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
		if state(f) != "approved" || strings.Contains(page, `value="approve"`) || !strings.Contains(page, string(badge("approved"))) {
			t.Fatal("double click must keep the first decision and offer no new one")
		}
	})
}
