// Landing regressions pin public content, configured authentication and safe error rendering.
package relay

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			for _, want := range []string{`id="main"`, `id="how"`, `id="install"`, `id="examples"`, `id="start"`, "기능을 설명하기 위한 가상 예시", "내가 작성", "사용자 승인 후 발송", "최대 180초", "최대 24시간", "sh scripts/install_bot_mcp.sh", "/workspace/.knowslink/node/bin/node", "/workspace/.knowslink/knowslink/dist/plugin.js", "KNOWSLINK_AGENT_FOLDER=/workspace/.knowslink/connections/grok-new-01", "이미 연결했다면 새로 등록하지 마세요"} {
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
