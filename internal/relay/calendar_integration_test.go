//go:build integration

// SAR-SCHEDULE-001 SC-01..SC-08 on real Postgres with the synthetic provider: connection, one lookup per approval,
// field minimisation, scope enforcement, failure classes, revocation, fixed copy and the untouched commit stub.
package relay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

// calendarClient is a browser that accepts the Lax Google attempt cookie as well as the Strict session.
type calendarClient struct {
	t       *testing.T
	h       http.Handler
	cookies map[string]*http.Cookie
}

func (c *calendarClient) do(method, path string, form url.Values) *httptest.ResponseRecorder {
	c.t.Helper()
	r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(form.Encode()))
	r.RemoteAddr = "192.0.2.77:40000"
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	for _, k := range c.cookies {
		r.AddCookie(k)
	}
	w := httptest.NewRecorder()
	c.h.ServeHTTP(w, r)
	for _, k := range w.Result().Cookies() {
		if k.MaxAge < 0 {
			delete(c.cookies, k.Name)
		} else {
			c.cookies[k.Name] = k
		}
	}
	return w
}

// calendarSetup binds agent_b to a signed-in Google member and points the relay at a synthetic provider.
func calendarSetup(t *testing.T) (*fixture, *syntheticProvider, *calendarClient, string) {
	t.Helper()
	f := setup(t, messagePool(t))
	g, sign := googleFixture(t)
	p := newSyntheticProvider(t, sign)
	server := httptest.NewServer(p.handler())
	t.Cleanup(server.Close)
	p.install(g, server.URL, server.Client())
	g.config.RedirectURL = "https://example.com/auth/google/callback"
	f.s.Google = g
	bob := &calendarClient{t, f.handler, map[string]*http.Cookie{}}
	var member string
	f.mutate(func(st *State) {
		v := st.signIn(googleIssuer, "sub-bob", "bob@example.test", "", "", time.Now())
		member = st.Sessions[hashToken(v.Session)].Member
		st.Members[member].Name = "밥"
		st.Agents["agent_b"].Owner = st.Members[member].Owner
		bob.cookies[sessionCookie] = &http.Cookie{Name: sessionCookie, Value: v.Session}
	})
	return f, p, bob, member
}

func (f *fixture) gateFor(t *testing.T, body map[string]any, id, gate, key string) string {
	t.Helper()
	receipt := f.send(wire(t, f.private["agent_a"], id, "agent_a", "agent_b", "schedule.query", key, "", body, time.Now().Add(180*time.Second)), "agent_a", "", 200)
	claim := f.deliver(id)
	f.send(wire(t, f.private["agent_b"], gate, "agent_b", "agent_b", "relay.approval.request", key+"-g", id, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
	return claim
}

func csrfToken(t *testing.T, page string) string {
	t.Helper()
	i := strings.Index(page, `name="csrf" value="`)
	if i < 0 {
		t.Fatalf("no decision form: %s", page)
	}
	rest := page[i+len(`name="csrf" value="`):]
	return rest[:strings.Index(rest, `"`)]
}

func TestCalendarConnection(t *testing.T) {
	f, p, bob, member := calendarSetup(t)
	page := bob.do("GET", "/home/account", nil).Body.String()
	for _, want := range []string{"친구가 요청하고 내가 승인할 때만, 해당 기간의 바쁜 시간대를 알려 줍니다. 일정 제목·참석자·장소는 공유하지 않습니다.", `action="/home/calendar/connect"`} {
		if !strings.Contains(page, want) {
			t.Fatalf("account page lacks %q", want)
		}
	}
	if strings.Contains(page, "연결됨") || strings.Contains(page, `action="/home/calendar/disconnect"`) {
		t.Fatal("unconnected account shows connected controls")
	}
	// SC-01: the consent asks only for free/busy on top of the login scopes; the login request itself is unchanged.
	w := bob.do("POST", "/home/calendar/connect", nil)
	expect(t, w, 303, "/fixture/google/auth")
	q, _ := url.Parse(w.Header().Get("Location"))
	if q.Query().Get("scope") != "openid email "+calendarScope || q.Query().Get("access_type") != "offline" || q.Query().Get("include_granted_scopes") != "true" || q.Query().Get("code_challenge_method") != "S256" {
		t.Fatalf("calendar consent request: %s", w.Header().Get("Location"))
	}
	reauth := bob.do("POST", "/auth/reauth", nil)
	expect(t, reauth, 303)
	if l, _ := url.Parse(reauth.Header().Get("Location")); l.Query().Get("scope") != "openid email" || l.Query().Get("access_type") != "" {
		t.Fatalf("login consent widened: %s", reauth.Header().Get("Location"))
	}
	// Partial consent (no calendar scope) stores nothing.
	w = bob.do("POST", "/home/calendar/connect", nil)
	expect(t, bob.do("GET", strings.TrimPrefix(p.grant(w.Header().Get("Location"), "sub-bob", "openid email"), "https://example.com"), nil), 403, "캘린더 권한이 허용되지 않아 연결하지 않았습니다")
	// Another Google account cannot be attached to this member (SD-06 C is out of scope).
	w = bob.do("POST", "/home/calendar/connect", nil)
	expect(t, bob.do("GET", strings.TrimPrefix(p.grant(w.Header().Get("Location"), "sub-other", "openid email "+calendarScope), "https://example.com"), nil), 403)
	// The other account's unstored refresh token is revoked instead of leaving a free/busy grant behind.
	_, _, _, mismatched := p.snapshot()
	if len(mismatched) != 1 || p.refresh[mismatched[0]] != "sub-other" {
		t.Fatalf("identity mismatch revocation: %v", mismatched)
	}
	f.mutate(func(st *State) {
		if len(st.Calendars) != 0 {
			t.Fatal("refused consents stored a connection")
		}
	})
	w = bob.do("POST", "/home/calendar/connect", nil)
	expect(t, bob.do("GET", strings.TrimPrefix(p.grant(w.Header().Get("Location"), "sub-bob", "openid email "+calendarScope), "https://example.com"), nil), 200, "캘린더 연결 완료", `href="/home/account#calendar"`)
	var sealed []byte
	f.mutate(func(st *State) {
		c := st.Calendars[member]
		if c == nil || c.State != "connected" || len(c.Sealed) == 0 {
			t.Fatal("connection not stored")
		}
		sealed = c.Sealed
		raw, _ := json.Marshal(st)
		if bytes.Contains(raw, []byte("fixture-refresh-")) {
			t.Fatal("refresh token stored in clear")
		}
	})
	refresh, err := openToken(p.key, sealed, member)
	if err != nil || !strings.HasPrefix(string(refresh), "fixture-refresh-") {
		t.Fatal("sealed token does not open to the provider's refresh token", err)
	}
	page = bob.do("GET", "/home/account", nil).Body.String()
	if !strings.Contains(page, "연결됨 — 바쁜 시간만 확인합니다") || !strings.Contains(page, `action="/home/calendar/disconnect"`) || strings.Contains(page, `action="/home/calendar/connect"`) {
		t.Fatalf("connected account page: %s", page)
	}
	// SC-01: disconnecting leaves no usable token and asks the provider to revoke.
	expect(t, bob.do("POST", "/home/calendar/disconnect", nil), 303, "/home/account?n=calendar-off#calendar")
	f.mutate(func(st *State) {
		if len(st.Calendars) != 0 {
			t.Fatal("token kept after disconnect")
		}
	})
	if _, _, _, revoked := p.snapshot(); len(revoked) != 2 || !slices.Contains(revoked, string(refresh)) {
		t.Fatalf("provider revocation: %v", revoked)
	}
	expect(t, bob.do("GET", "/home/account?n=calendar-off", nil), 200, "캘린더 연결을 해제했습니다", `action="/home/calendar/connect"`)
	// An email-code member has no Google identity to extend: one line, no button.
	alice := &browser{t, f.handler, "192.0.2.78", map[string]string{}}
	mail := &inbox{codes: map[string]string{}}
	f.s.Mail = mail.send
	aliceMember := login(t, alice, mail, "alice@example.test")
	f.mutate(func(st *State) { st.Members[aliceMember].Name = "앨리스" })
	expect(t, alice.do("GET", "/home/account", nil), 200, "Google 계정으로 로그인한 회원만 캘린더를 연결할 수 있습니다.")
	if strings.Contains(alice.do("GET", "/home/account", nil).Body.String(), `action="/home/calendar/connect"`) {
		t.Fatal("email member offered a calendar connection")
	}
}

func TestCalendarApprovalRunsOneLookup(t *testing.T) {
	f, p, bob, member := calendarSetup(t)
	var logs bytes.Buffer
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	f.mutate(func(st *State) {
		st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected", Changed: time.Now()}}
	})
	start, _ := queryWindow()
	end := start.Add(2 * time.Hour)
	body := map[string]any{"window": map[string]any{"start": start.Format("2006-01-02T15:04:05Z"), "end": end.Format("2006-01-02T15:04:05Z")}, "granularity_min": 30}
	claim := f.gateFor(t, body, firstID, gateID, "calendar-key-0001")
	// SC-07: the connected member sees the fixed explanation, the disclosed-content row, and the decision form.
	page := bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
	for _, want := range []string{calendarGateScope, "공개되는 내용", "바쁜 시간대만", `value="approve"`} {
		if !strings.Contains(page, want) {
			t.Fatalf("connected gate page lacks %q: %s", want, page)
		}
	}
	for _, banned := range []string{gateScopes["schedule.query"], "정보 공개나 일정 실행을 허용하지 않습니다", "캘린더가 연결되지 않아"} {
		if strings.Contains(page, banned) {
			t.Fatalf("connected gate page shows %q", banned)
		}
	}
	// The agent sees pending until the approval's lookup has ended; a premature consume is refused.
	if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "pending" {
		t.Fatal("gate visible before decision", g)
	}
	token := csrfToken(t, page)
	expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 303, "/home/gates/"+gateID)
	calls, min, max, _ := p.snapshot()
	if calls != 1 || min != start.Format(time.RFC3339) || max != end.Format(time.RFC3339) {
		t.Fatalf("lookup calls=%d window=%s..%s", calls, min, max)
	}
	want := disclosureJSON(start, end, 30, []span{{start, start.Add(30 * time.Minute)}, {end.Add(-30 * time.Minute), end}})
	f.mutate(func(st *State) {
		g := st.Gates[gateID]
		if g.State != "approved" || g.Lookup == nil || g.Lookup.State != "done" || canonicalRaw(g.Lookup.Disclosure) != string(want) {
			t.Fatalf("lookup record: %+v", g.Lookup)
		}
		raw, _ := json.Marshal(st)
		for _, secret := range []string{"secret", "fixture-refresh-", "attendees", "location"} {
			if bytes.Contains(raw, []byte(secret)) {
				t.Fatalf("state carries provider field %q", secret)
			}
		}
	})
	// SC-02: a repeated submission, a re-opened page and a restarted process never run a second lookup.
	expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 409, "이번 제출은 적용되지 않았습니다")
	page = bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
	if strings.Contains(page, `value="approve"`) || !strings.Contains(page, "지금은 승인하거나 거절할 수 없습니다") {
		t.Fatalf("approved page: %s", page)
	}
	restarted := (&Service{Pool: f.s.Pool, Google: f.s.Google}).Handler()
	r := httptest.NewRequest("GET", "/v1/gates/"+gateID, nil)
	r.Header.Set("Authorization", "Bearer "+f.tokens["agent_b"])
	rec := httptest.NewRecorder()
	restarted.ServeHTTP(rec, r)
	gate := map[string]any{}
	_ = json.Unmarshal(rec.Body.Bytes(), &gate)
	if rec.Code != 200 || gate["state"] != "approved" || gate["lookup"] != "done" {
		t.Fatalf("restarted gate view: %d %s", rec.Code, rec.Body)
	}
	consumed := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 200)
	echoed, _ := json.Marshal(consumed["disclosure"])
	canonical, _ := Canonical(consumed["disclosure"])
	if consumed["executable"] != false || consumed["lookup"] != "done" || string(canonical) != string(want) {
		t.Fatalf("consume response: %s", echoed)
	}
	f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
	if calls, _, _, _ := p.snapshot(); calls != 1 {
		t.Fatal("consume or restart ran another lookup")
	}
	// SC-03: the done result must be the server's disclosure, unchanged; anything else is refused.
	var disclosure map[string]any
	_ = json.Unmarshal(want, &disclosure)
	resultID := "0199a3f2-4c10-7a11-8b22-3344556677bb"
	tampered := map[string]any{"window": disclosure["window"], "granularity_min": disclosure["granularity_min"], "busy": append(disclosure["busy"].([]any), map[string]any{"start": start.Add(30 * time.Minute).Format("2006-01-02T15:04:05Z"), "end": start.Add(time.Hour).Format("2006-01-02T15:04:05Z")})}
	f.send(wire(t, f.private["agent_b"], resultID, "agent_b", "agent_a", "relay.result", "result-key-00001", firstID, map[string]any{"status": "done", "disclosure": tampered}, time.Now().Add(time.Minute)), "agent_b", claim, 403)
	f.send(wire(t, f.private["agent_b"], resultID, "agent_b", "agent_a", "relay.result", "result-key-00002", firstID, map[string]any{"status": "done"}, time.Now().Add(time.Minute)), "agent_b", claim, 403)
	done := wire(t, f.private["agent_b"], resultID, "agent_b", "agent_a", "relay.result", "result-key-00003", firstID, map[string]any{"status": "done", "disclosure": disclosure}, time.Now().Add(time.Minute))
	accepted := f.send(done, "agent_b", claim, 200)
	// A same-key resend of the accepted done gets its receipt back, although the disclosure is already cleared.
	if again := f.send(done, "agent_b", claim, 200); again["digest"] != accepted["digest"] || again["id"] != accepted["id"] {
		t.Fatalf("done resend receipt: %v vs %v", again, accepted)
	}
	f.send(wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677cc", "agent_b", "agent_a", "relay.result", "result-key-00004", firstID, map[string]any{"status": "done", "disclosure": disclosure}, time.Now().Add(time.Minute)), "agent_b", claim, 403)
	if receipt := f.call("GET", "/v1/receipts/"+firstID, f.tokens["agent_a"], nil, 200); receipt["completion"] != "done" {
		t.Fatal("requester receipt", receipt)
	}
	lease := f.call("POST", "/v1/pull", f.tokens["agent_a"], map[string]any{}, 200)
	delivered, _ := json.Marshal(lease["envelope"])
	if !strings.Contains(string(delivered), `"status":"done"`) || !strings.Contains(string(delivered), start.Format("2006-01-02T15:04:05Z")) || strings.Contains(string(delivered), "secret") {
		t.Fatalf("requester result: %s", delivered)
	}
	f.mutate(func(st *State) {
		if len(st.Gates[gateID].Lookup.Disclosure) != 0 {
			t.Fatal("disclosure kept after the single result")
		}
	})
	if out := logs.String(); strings.Contains(out, "secret") || strings.Contains(out, "fixture-refresh-") || strings.Contains(out, "bob") {
		t.Fatalf("log leak: %s", out)
	}
}

func TestCalendarUnconnectedAndCommitStub(t *testing.T) {
	f, p, bob, member := calendarSetup(t)
	claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0002")
	// SC-07/05: no connected calendar means no approval button and no lookup; the U5 text stays.
	page := bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
	for _, want := range []string{"캘린더가 연결되지 않아 이 요청에 답할 수 없습니다.", `href="/home/account#calendar"`, "이 요청은 승인할 수 없습니다. 거절하거나 기한이 지나면 조회 없이 끝납니다.", "친구 동의나 대화 메시지는", "바쁜 시간대만"} {
		if !strings.Contains(page, want) {
			t.Fatalf("unconnected gate page lacks %q: %s", want, page)
		}
	}
	// V1: the two approval-scope lines contradict a page without an approve button.
	for _, banned := range []string{gateScopes["schedule.query"], "정보 공개나 일정 실행을 허용하지 않습니다"} {
		if strings.Contains(page, banned) {
			t.Fatalf("unconnected gate page shows %q", banned)
		}
	}
	if strings.Contains(page, `value="approve"`) || strings.Contains(page, calendarGateScope) || !strings.Contains(page, `value="deny"`) {
		t.Fatal("unconnected member must see deny only")
	}
	// A stale approve POST (page rendered before a disconnect) is refused and the request stays pending.
	f.mutate(func(st *State) {
		st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
	})
	token := csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())
	f.mutate(func(st *State) { delete(st.Calendars, member) })
	expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 409)
	if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "pending" {
		t.Fatal("stale approval decided the gate", g)
	}
	// Reconnection needed: the provider refuses the stored token, so approval records a failure, never a result.
	f.mutate(func(st *State) {
		st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "reconnect"}}
	})
	page = bob.do("GET", "/home/gates/"+gateID, nil).Body.String()
	if strings.Contains(page, `value="approve"`) || !strings.Contains(page, "캘린더가 연결되지 않아") {
		t.Fatal("reconnect state offered approval")
	}
	expect(t, bob.do("GET", "/home/account", nil), 200, "다시 연결 필요", `action="/home/calendar/connect"`, `action="/home/calendar/disconnect"`)
	if calls, _, _, _ := p.snapshot(); calls != 0 {
		t.Fatal("lookup ran without a connected calendar")
	}
	// Deny needs no calendar and runs no lookup; the agent then reports denied.
	expect(t, bob.do("POST", "/home/gates/"+gateID+"/deny", url.Values{"csrf": {token}, "decision": {"deny"}}), 303)
	f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
	f.send(wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677bb", "agent_b", "agent_a", "relay.result", "result-key-00001", firstID, map[string]any{"status": "denied"}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
	// SC-08: schedule.commit keeps the U5 copy, runs no lookup and still cannot report done.
	f.mutate(func(st *State) { st.Calendars[member].State = "connected" })
	commit := map[string]any{"slot": map[string]any{"start": "2026-12-01T01:00:00Z", "end": "2026-12-01T01:30:00Z"}, "timezone": "Asia/Seoul", "commitment": "확정해 주세요"}
	receipt := f.send(wire(t, f.private["agent_a"], secondID, "agent_a", "agent_b", "schedule.commit", "calendar-key-0003", "", commit, time.Now().Add(180*time.Second)), "agent_a", "", 200)
	commitClaim := f.deliver(secondID)
	commitGate := "0199a3f2-4c10-7a11-8b22-3344556677dd"
	f.send(wire(t, f.private["agent_b"], commitGate, "agent_b", "agent_b", "relay.approval.request", "calendar-key-0003-g", secondID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", commitClaim, 200)
	page = bob.do("GET", "/home/gates/"+commitGate, nil).Body.String()
	if !strings.Contains(page, gateScopes["schedule.commit"]) || !strings.Contains(page, "정보 공개나 일정 실행을 허용하지 않습니다") || strings.Contains(page, calendarGateScope) || strings.Contains(page, "캘린더") {
		t.Fatalf("commit gate copy changed: %s", page)
	}
	expect(t, bob.do("POST", "/home/gates/"+commitGate, url.Values{"csrf": {csrfToken(t, page)}, "decision": {"approve"}}), 303)
	consumed := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": commitGate, "claim": commitClaim}, 200)
	if consumed["executable"] != false || consumed["disclosure"] != false || consumed["lookup"] != nil {
		t.Fatal("commit gate gained an effect", consumed)
	}
	f.send(wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677ee", "agent_b", "agent_a", "relay.result", "result-key-00005", secondID, map[string]any{"status": "done"}, time.Now().Add(time.Minute)), "agent_b", commitClaim, 403)
	if calls, _, _, _ := p.snapshot(); calls != 0 {
		t.Fatal("commit approval ran a lookup")
	}
}

func TestCalendarScopeEnforcedAtAccept(t *testing.T) {
	f, _, _, _ := calendarSetup(t)
	start, _ := queryWindow()
	iso := func(t time.Time) string { return t.Format("2006-01-02T15:04:05Z") }
	body := func(start, end time.Time, minutes int) map[string]any {
		return map[string]any{"window": map[string]any{"start": iso(start), "end": iso(end)}, "granularity_min": minutes}
	}
	id := func(i int) string { return fmt.Sprintf("0199a3f2-4c10-7a11-8b22-3344556600%02d", i) }
	refused := func(name string, b map[string]any) {
		t.Helper()
		w := f.request("POST", "/v1/send", f.tokens["agent_a"], wire(t, f.private["agent_a"], id(0), "agent_a", "agent_b", "schedule.query", "scope-refused-"+name+"-0001", "", b, time.Now().Add(time.Minute)), "")
		if w.Code != 422 || !strings.Contains(w.Body.String(), "schedule_scope") {
			t.Fatalf("%s accepted: %d %s", name, w.Code, w.Body)
		}
	}
	refused("unit", body(start, start.Add(time.Hour), 20))
	refused("past", body(start.Add(-2*time.Hour), start.Add(time.Hour), 30))
	refused("long", body(start, time.Now().Add(scheduleScope.Window+time.Hour), 30))
	for i := 1; i <= scheduleScope.PerFriend; i++ {
		f.send(wire(t, f.private["agent_a"], id(i), "agent_a", "agent_b", "schedule.query", fmt.Sprintf("scope-count-%04d-0001", i), "", body(start, start.Add(time.Hour), 30), time.Now().Add(time.Minute)), "agent_a", "", 200)
	}
	refused("sixth", body(start, start.Add(time.Hour), 30))
	f.mutate(func(st *State) {
		for _, m := range st.Messages {
			if m.Receipt.Intent == "schedule.query" {
				m.Receipt.Accepted = m.Receipt.Accepted.Add(-scheduleScope.PerFriendWindow)
			}
		}
	})
	f.send(wire(t, f.private["agent_a"], id(6), "agent_a", "agent_b", "schedule.query", "scope-after-window-0001", "", body(start, start.Add(time.Hour), 30), time.Now().Add(time.Minute)), "agent_a", "", 200)
}

func TestCalendarRevocationAndFailureClasses(t *testing.T) {
	t.Run("disconnect_ends_pending_without_lookup", func(t *testing.T) {
		f, p, bob, member := calendarSetup(t)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
		})
		claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0010")
		token := csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())
		expect(t, bob.do("POST", "/home/calendar/disconnect", nil), 303)
		if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "denied" {
			t.Fatal("pending request survived disconnect", g)
		}
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {token}, "decision": {"approve"}}), 409)
		f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
		if calls, _, _, _ := p.snapshot(); calls != 0 {
			t.Fatal("lookup after disconnect")
		}
	})
	// D2: an approved lookup the agent has not consumed is dropped like a pending request; a consumed one is not recalled.
	t.Run("disconnect_drops_undelivered_result", func(t *testing.T) {
		f, p, bob, member := calendarSetup(t)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
		})
		claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0012")
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())}, "decision": {"approve"}}), 303)
		if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "approved" || g["lookup"] != "done" {
			t.Fatal("lookup not done before disconnect", g)
		}
		expect(t, bob.do("POST", "/home/calendar/disconnect", nil), 303)
		if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "denied" {
			t.Fatal("undelivered result survived disconnect", g)
		}
		f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
		f.mutate(func(st *State) {
			if len(st.Gates[gateID].Lookup.Disclosure) != 0 {
				t.Fatal("disclosure kept after disconnect")
			}
		})
		f.send(wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677bb", "agent_b", "agent_a", "relay.result", "result-key-00012", firstID, map[string]any{"status": "denied"}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
	})
	t.Run("disconnect_keeps_consumed_result", func(t *testing.T) {
		f, p, bob, member := calendarSetup(t)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
		})
		claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0013")
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())}, "decision": {"approve"}}), 303)
		consumed := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 200)
		expect(t, bob.do("POST", "/home/calendar/disconnect", nil), 303)
		if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "approved" || g["consumed"] != true {
			t.Fatal("consumed result recalled", g)
		}
		f.send(wire(t, f.private["agent_b"], "0199a3f2-4c10-7a11-8b22-3344556677bb", "agent_b", "agent_a", "relay.result", "result-key-00013", firstID, map[string]any{"status": "done", "disclosure": consumed["disclosure"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
	})
	t.Run("owner_revoke_disconnects_and_revokes", func(t *testing.T) {
		f, p, _, member := calendarSetup(t)
		sealed := p.sealedFor("sub-bob", member)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: sealed, State: "connected"}}
			st.Owners[st.Members[member].Owner].Credential = hashToken("owner-token-bob")
		})
		if v := f.call("POST", "/v1/owner-revoke", "owner-token-bob", map[string]any{}, 200); len(v) != 1 || v["state"] != "revoked" {
			t.Fatal("owner-revoke response", v)
		}
		f.mutate(func(st *State) {
			if len(st.Calendars) != 0 {
				t.Fatal("sealed token kept after owner-revoke")
			}
		})
		refresh, _ := openToken(p.key, sealed, member)
		if _, _, _, revoked := p.snapshot(); len(revoked) != 1 || revoked[0] != string(refresh) {
			t.Fatalf("owner-revoke provider revocation: %v", revoked)
		}
	})
	t.Run("provider_revocation_marks_reconnect", func(t *testing.T) {
		f, p, bob, member := calendarSetup(t)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
		})
		claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0011")
		p.mu.Lock()
		p.grantRevoked = true
		p.mu.Unlock()
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())}, "decision": {"approve"}}), 303)
		f.mutate(func(st *State) {
			if st.Calendars[member].State != "reconnect" || st.Gates[gateID].Lookup.State != "failed" || len(st.Gates[gateID].Lookup.Disclosure) != 0 {
				t.Fatalf("revoked grant: calendar=%s lookup=%+v", st.Calendars[member].State, st.Gates[gateID].Lookup)
			}
		})
		expect(t, bob.do("GET", "/home/account", nil), 200, "다시 연결 필요")
		consumed := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 200)
		if consumed["disclosure"] != false || consumed["lookup"] != "failed" {
			t.Fatal("failed lookup disclosed", consumed)
		}
		f.send(wire(t, f.private["agent_b"], secondID, "agent_b", "agent_a", "relay.result", "result-key-00011", firstID, map[string]any{"status": "done"}, time.Now().Add(time.Minute)), "agent_b", claim, 403)
		f.send(wire(t, f.private["agent_b"], secondID, "agent_b", "agent_a", "relay.result", "result-key-00012", firstID, map[string]any{"status": "failed"}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
		if calls, _, _, _ := p.snapshot(); calls != 0 {
			t.Fatal("free/busy called without an access token")
		}
	})
	t.Run("outage_fails_after_short_retry_and_keeps_connection", func(t *testing.T) {
		f, p, bob, member := calendarSetup(t)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
		})
		f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0012")
		p.mu.Lock()
		p.outage = true
		p.mu.Unlock()
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())}, "decision": {"approve"}}), 303)
		calls, _, _, _ := p.snapshot()
		f.mutate(func(st *State) {
			if st.Calendars[member].State != "connected" || st.Gates[gateID].Lookup.State != "failed" || calls != 3 {
				t.Fatalf("outage: calendar=%s lookup=%s calls=%d", st.Calendars[member].State, st.Gates[gateID].Lookup.State, calls)
			}
		})
		if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "approved" || g["lookup"] != "failed" {
			t.Fatal("failed lookup not reported", g)
		}
	})
	t.Run("relationship_revoked_after_approval_discloses_nothing", func(t *testing.T) {
		f, p, bob, member := calendarSetup(t)
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: p.sealedFor("sub-bob", member), State: "connected"}}
		})
		claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0013")
		expect(t, bob.do("POST", "/home/gates/"+gateID, url.Values{"csrf": {csrfToken(t, bob.do("GET", "/home/gates/"+gateID, nil).Body.String())}, "decision": {"approve"}}), 303)
		f.call("POST", "/v1/unpair", f.tokens["agent_a_owner"], map[string]any{"agent": "agent_a", "target": "agent_b"}, 200)
		f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
		f.mutate(func(st *State) {
			if st.Gates[gateID].State != "revoked" || len(st.Gates[gateID].Lookup.Disclosure) != 0 {
				t.Fatalf("revoked gate kept disclosure: %s %+v", st.Gates[gateID].State, st.Gates[gateID].Lookup)
			}
		})
		if calls, _, _, _ := p.snapshot(); calls != 1 {
			t.Fatal("lookup count", calls)
		}
	})
	t.Run("crashed_lookup_lease_fails_closed", func(t *testing.T) {
		f, _, _, member := calendarSetup(t)
		claim := f.gateFor(t, queryBody(), firstID, gateID, "calendar-key-0014")
		f.mutate(func(st *State) {
			st.Calendars = map[string]*Calendar{member: {Sealed: []byte("x"), State: "connected"}}
			st.Gates[gateID].State = "approved"
			st.Gates[gateID].Lookup = &Lookup{Started: time.Now().Add(-lookupLease)}
		})
		if g := f.call("GET", "/v1/gates/"+gateID, f.tokens["agent_b"], nil, 200); g["state"] != "approved" || g["lookup"] != "failed" {
			t.Fatal("expired lease", g)
		}
		consumed := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 200)
		if consumed["disclosure"] != false {
			t.Fatal("expired lease disclosed", consumed)
		}
	})
}
