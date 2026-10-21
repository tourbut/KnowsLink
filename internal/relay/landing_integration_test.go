//go:build integration

// Public guidance runs against real synthetic sessions without changing member authorization or login recovery.
package relay

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestLandingSessionAndLoginRecovery(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated synthetic database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, mail := identityService(t, pool)
	b := &browser{t, s.Handler(), "192.0.2.70", map[string]string{}}
	w := b.do("GET", "/", nil)
	expect(t, w, 200, "친구와 연결하고, 에이전트로 대화하세요")
	if w.Header().Get("Cache-Control") != "no-store" || !strings.Contains(w.Header().Get("Content-Security-Policy"), "form-action 'self' https://accounts.google.com") {
		t.Fatal("anonymous guidance cache/CSP boundary changed")
	}
	expect(t, b.do("GET", "/home", nil), 303, "/?n=expired")
	expect(t, b.do("POST", "/home/profile", url.Values{"name": {"거부"}}), 303)
	expect(t, b.do("POST", "/auth/start", url.Values{"email": {"bad-address"}}), 422, "로그인 문제를 확인하세요", `aria-describedby="email-help problem"`)
	mail.fail = true
	expect(t, b.do("POST", "/auth/start", url.Values{"email": {"failed@example.test"}}), 503, "로그인되지 않았습니다", "1분 뒤")
	mail.fail = false
	id := login(t, b, mail, "landing-private@example.test")
	s.mutateState(t, func(st *State) { st.Members[id].Name = "" })
	w = b.do("GET", "/", nil)
	expect(t, w, 200, "로그인되어 있습니다.", `href="/home"`)
	if w.Header().Get("Cache-Control") != "no-store" || !strings.Contains(w.Header().Get("Content-Security-Policy"), "form-action 'self' https://accounts.google.com") {
		t.Fatal("authenticated guidance cache/CSP boundary changed")
	}
	if w.Header().Get("Location") != "" || strings.Contains(w.Body.String(), "landing-private") || strings.Contains(w.Body.String(), `action="/auth/start"`) {
		t.Fatal("authenticated guidance must not redirect, expose identity or ask to log in")
	}
	expect(t, b.do("GET", "/home", nil), 303, "/home/name")
	expect(t, b.do("POST", "/home/profile", url.Values{"name": {"비공개 합성 이름"}}), 303)
	expect(t, b.do("GET", "/home/account", nil), 200, `href="/">서비스 안내</a>`)
	w = b.do("GET", "/", nil)
	if strings.Contains(w.Body.String(), "비공개 합성 이름") || strings.Contains(w.Body.String(), id) {
		t.Fatal("public guidance leaked member data")
	}
	expect(t, b.do("POST", "/home/profile", url.Values{"name": {"변경 금지"}}, "Origin", "https://foreign.test"), 403)
	s.mutateState(t, func(st *State) {
		st.Sessions[hashToken(b.cookies[sessionCookie])].Created = time.Now().Add(-sessionAbsolute - time.Minute)
	})
	expect(t, b.do("GET", "/", nil), 200, notices["expired"], `action="/auth/start"`)
	if b.cookies[sessionCookie] != "" {
		t.Fatal("expired cookie survived public guidance")
	}
	s.Google, s.Mail = &GoogleLogin{}, nil
	expect(t, b.do("GET", "/auth/verify", nil), 410, "Google로 계속")
	w = b.do("GET", "/auth/google/callback?error=access_denied", nil)
	expect(t, w, 401, "로그인 문제를 확인하세요", "로그인 다시 시도")
	if strings.Contains(w.Body.String(), `action="/auth/start"`) {
		t.Fatal("disabled mail login appeared on error")
	}
	// The same configured view is used by invalid email requests, even if the endpoint is called directly.
	w = b.do("POST", "/auth/start", url.Values{"email": {"invalid"}})
	expect(t, w, 422, "Google로 계속")
	if strings.Contains(w.Body.String(), `action="/auth/start"`) {
		t.Fatal("direct POST re-enabled the disabled form")
	}
}

// SAR-PUBLISH-001 PC-01: both policy pages open without a session through the real handler, follow the browser's
// language choice, and follow a signed-in member's saved language as the first page does.
func TestPolicyPagesPublic(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated synthetic database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, mail := identityService(t, pool)
	b := &browser{t, s.Handler(), "192.0.2.72", map[string]string{}}
	expect(t, b.do("GET", "/privacy", nil), 200, "<h1>개인정보처리방침</h1>", policyEffective)
	expect(t, b.do("GET", "/terms", nil), 200, "<h1>서비스 약관</h1>")
	b.cookies[langCookie] = "en"
	w := b.do("GET", "/terms", nil)
	expect(t, w, 200, "<h1>Terms of Service</h1>")
	if w.Header().Get("Content-Language") != "en" || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("policy language or cache header", w.Header())
	}
	expect(t, b.do("GET", "/privacy/", nil), 404)
	delete(b.cookies, langCookie)
	id := login(t, b, mail, "policy-member@example.test")
	s.mutateState(t, func(st *State) { st.Members[id].Locale = "en" })
	expect(t, b.do("GET", "/privacy", nil), 200, "<h1>Privacy Policy</h1>", "Effective date: "+policyEffective)
}
