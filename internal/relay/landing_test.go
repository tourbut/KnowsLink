// Landing regressions pin public content, configured authentication and safe error rendering.
package relay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// A database failure while checking the session is not an expired session: the cookie stays and no expiry is shown.
func TestSessionLookupFailureKeepsCookie(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://x@127.0.0.1:1/x?connect_timeout=1") // nothing listens
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := &Service{Pool: pool}
	for _, tc := range []struct {
		method, path string
		handler      http.HandlerFunc
		status       int
	}{{"GET", "/", s.startPage, 200}, {"POST", "/auth/reauth", s.authReauth, 503}} {
		r := httptest.NewRequest(tc.method, tc.path, nil)
		r.AddCookie(&http.Cookie{Name: sessionCookie, Value: "synthetic-session"})
		w := httptest.NewRecorder()
		tc.handler(w, r)
		if w.Code != tc.status || len(w.Result().Cookies()) != 0 || strings.Contains(w.Body.String(), "세션이 끝났습니다") {
			t.Fatalf("%s %s: %d %v", tc.method, tc.path, w.Code, w.Result().Cookies())
		}
	}
}

func TestLandingConfiguredLoginAndPublicContent(t *testing.T) {
	for _, mode := range []string{"google", "both", "email", "none"} {
		t.Run(mode, func(t *testing.T) {
			s := &Service{}
			google, email := mode == "google" || mode == "both", mode == "email" || mode == "both"
			if google {
				s.Google = &GoogleLogin{}
			}
			if email {
				s.Mail = func(context.Context, string, string, string) error { return nil }
			}
			w := httptest.NewRecorder()
			s.startPage(w, httptest.NewRequest("GET", "/", nil))
			body := w.Body.String()
			if w.Code != 200 || w.Header().Get("Location") != "" {
				t.Fatal("public start handler must render without a session lookup", w.Code, w.Header())
			}
			for action, enabled := range map[string]bool{"/auth/google": google, "/auth/start": email} {
				count := strings.Count(body, `action="`+action+`"`)
				if count > 1 || (count == 1) != enabled {
					t.Fatal("login form must exist exactly once iff configured", action)
				}
			}
			for _, want := range []string{`id="main"`, `id="how"`, `id="install"`, `id="examples"`, `id="start"`, "기능을 설명하기 위한 가상 예시", "내가 작성", "사용자 승인 후 발송", "최대 180초", "최대 1시간 서버 메모리(Redis)", "sh scripts/install_bot_mcp.sh", "/workspace/.knowslink/node/bin/node", "/workspace/.knowslink/knowslink/dist/plugin.js", "KNOWSLINK_AGENT_FOLDER=/workspace/.knowslink/connections/grok-new-01", "이미 연결했다면 새로 등록하지 마세요"} {
				if !strings.Contains(body, want) {
					t.Fatal("missing public requirement", want)
				}
			}
			for _, forbidden := range []string{"<script", "<img", "<link", "url(", ".fullops-squad", "SAR-LANDING", "C:/Users", "<details open"} {
				if strings.Contains(body, forbidden) {
					t.Fatal("public guidance leaked implementation or added resources", forbidden)
				}
			}
			v := s.startView("synthetic <unsafe>")
			v["Email"] = `"><script>alert(1)</script>`
			body = renderFixture(t, uiFixture{"landing-error", "start", v})
			if strings.Contains(body, "<script>") || strings.Contains(body, "<unsafe>") || !strings.Contains(body, "로그인 문제를 확인하세요") || !strings.Contains(body, `role="alert"`) {
				t.Fatal("error must escape input and remain discoverable")
			}
			if email && !strings.Contains(body, `aria-describedby="email-help problem"`) {
				t.Fatal("email error lost field association")
			}
		})
	}
}
