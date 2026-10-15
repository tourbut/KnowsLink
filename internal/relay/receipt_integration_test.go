//go:build integration

// Isolated SQL and real login handlers verify receipt recovery without exposing another member's context.
package relay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReceiptRecoveryHTTP(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, mail := identityService(t, pool)
	a := &browser{t, s.Handler(), "192.0.2.71", map[string]string{}}
	b := &browser{t, s.Handler(), "192.0.2.72", map[string]string{}}
	c := &browser{t, s.Handler(), "192.0.2.73", map[string]string{}}
	ma, mb, mc := login(t, a, mail, "receipt-a@example.test"), login(t, b, mail, "receipt-b@example.test"), login(t, c, mail, "receipt-c@example.test")
	s.mutateState(t, func(st *State) {
		for member, agent := range map[string]string{ma: "agent_a", mb: "agent_b", mc: "agent_c"} {
			st.Agents[agent] = &Agent{Owner: st.Members[member].Owner, Name: "이름-" + agent, Client: "Node 22"}
		}
		st.Messages[firstID] = &Message{Receipt: Receipt{ID: firstID, From: "agent_a", To: "agent_b", State: "delivered", Accepted: time.Now(), Exp: time.Now().Add(time.Minute)}, Completion: "received"}
	})
	path := "/home/receipts?agent=agent_a&id=" + firstID
	expect(t, a.do("GET", "/home/receipts", nil), 200, "요청 상태 조회", `name="id"`)
	expect(t, a.do("GET", path, nil), 200, `value="`+firstID+`"`, `value="agent_a" selected`, "수신 클라이언트가 보관했습니다", `id="receipt-result"`)
	expect(t, a.do("GET", "/home/receipts", nil), 200, `value="`+firstID+`"`, "조회 결과")
	for _, id := range []string{secondID, gateID} {
		if id == gateID {
			s.mutateState(t, func(st *State) {
				st.Messages[gateID] = &Message{Receipt: Receipt{ID: gateID, From: "agent_b", To: "agent_c", State: "delivered", Accepted: time.Now(), Exp: time.Now().Add(time.Minute)}}
			})
		}
		expect(t, a.do("GET", "/home/receipts?agent=agent_a&id="+id, nil), 403, receiptUnavailable, `value="`+id+`"`, `value="agent_a" selected`)
	}
	expect(t, a.do("GET", "/home/receipts?agent=agent_a&id="+url.QueryEscape(`<script>bad()</script>`), nil), 422, "&lt;script&gt;", `value="agent_a" selected`)
	s.mutateState(t, func(st *State) {
		st.Messages[firstID].Receipt.Exp = time.Now().Add(-time.Minute)
		st.Messages[firstID].Receipt.State = "queued"
	})
	expect(t, a.do("GET", path, nil), 200, "전달 기한 만료", "failed:expired")
	s.mutateState(t, func(st *State) {
		for range 40 {
			st.Rates["http:member:"+ma] = append(st.Rates["http:member:"+ma], time.Now())
		}
	})
	expect(t, a.do("GET", path, nil), 429, `value="`+firstID+`"`, `value="agent_a" selected`, "이후 다시 시도")
	s.mutateState(t, func(st *State) {
		st.Rates = map[string][]time.Time{}
		st.Sessions[hashToken(a.cookies[sessionCookie])].Seen = time.Now().Add(-sessionIdle)
	})
	expect(t, a.do("GET", path, nil), 303, "/?n=expired")
	if login(t, a, mail, "receipt-a@example.test") != ma {
		t.Fatal("member changed on re-login")
	}
	expect(t, a.do("GET", "/home/receipts", nil), 200, `value="`+firstID+`"`, `value="agent_a" selected`)

	// A DB outage cannot authenticate even a saved cookie. Preserve it without displaying private input.
	draft := a.cookies[receiptCookie]
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequest("GET", path, nil).WithContext(ctx)
	for name, value := range a.cookies {
		r.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	w := httptest.NewRecorder()
	s.memberReceipt(w, r)
	expect(t, w, 503, "로그인 확인이 가능해지면 복원")
	if strings.Contains(w.Body.String(), firstID) || len(w.Result().Cookies()) != 0 {
		t.Fatal("outage exposed or erased private context")
	}
	expect(t, a.do("GET", "/home/receipts", nil), 200, `value="`+firstID+`"`)

	// Re-login as B with A's draft and old URL must not reveal A's selection, name or request ID.
	a.cookies[sessionCookie] = b.cookies[sessionCookie]
	w = a.do("GET", path, nil)
	expect(t, w, 403, receiptUnavailable)
	for _, private := range []string{firstID, "이름-agent_a", `value="agent_a"`} {
		if strings.Contains(w.Body.String(), private) {
			t.Fatal("account switch leaked", private)
		}
	}
	a.cookies[receiptCookie] = draft
	expect(t, a.do("GET", "/home/receipts", nil), 200)
	if strings.Contains(a.do("GET", "/home/receipts", nil).Body.String(), firstID) {
		t.Fatal("other account restored draft")
	}
	expect(t, a.do("POST", "/auth/logout", nil), 303)
	if a.cookies[receiptCookie] != "" {
		t.Fatal("logout retained draft")
	}
	expect(t, a.do("GET", path, nil), 303)
	// Receipt retention is still 24 hours; it has the same refusal as an unknown ID.
	s.mutateState(t, func(st *State) { st.Messages[firstID].Receipt.Accepted = time.Now().Add(-24 * time.Hour) })
	expect(t, b.do("GET", "/home/receipts?agent=agent_b&id="+firstID, nil), 403, receiptUnavailable)
}
