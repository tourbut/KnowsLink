// SAR-PUBLISH-001 PC-01–05: both policy pages in both languages, their links, and no calendar wording left.
package relay

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func policyRender(t *testing.T, page, lang string) string {
	t.Helper()
	w := httptest.NewRecorder()
	w.Header().Set("Content-Language", lang)
	(&Service{}).policyPage(page)(w, httptest.NewRequest("GET", "/"+page, nil))
	if w.Code != 200 {
		t.Fatal(page, lang, w.Code)
	}
	// KNOWSLINK_UI_FIXTURE_DIR also writes each page for the 1440/390/320 render probe.
	if dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "policy-"+page+"-"+lang+".html"), w.Body.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return w.Body.String()
}

func TestPolicyPages(t *testing.T) {
	titles := map[string][2]string{"privacy": {"개인정보처리방침", "Privacy Policy"}, "terms": {"서비스 약관", "Terms of Service"}}
	other := map[string]string{"privacy": "/terms", "terms": "/privacy"}
	for page, title := range titles {
		for i, lang := range []string{"ko", "en"} {
			body := policyRender(t, page, lang)
			for _, want := range []string{`<html lang="` + lang + `">`, "<title>" + title[i] + " · KnowsLink</title>", "<h1>" + title[i] + "</h1>", policyEffective, `href="` + other[page] + `"`,
				`href="/` + page + `" aria-current="page"`, `<a href="/">`, `href="mailto:knowslog01@gmail.com"`, `href="#main"`, `class="brand" href="/"`, `href="/lang/en?back=%2f` + page + `"`} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s %s missing %q", page, lang, want)
				}
			}
			// D051: both pages state in both languages that messages are never stored permanently and leave within 1 hour.
			retention := map[string][]string{"ko": {"영구 저장하지 않습니다", "최대 1시간", "재시작"}, "en": {"never stored permanently", "at most 1 hour", "restart"}}[lang]
			// D054: the privacy policy also names the bodiless gate and delivery metadata kept in Postgres and its backups.
			if page == "privacy" {
				retention = append(retention, map[string][]string{"ko": {"서버 메모리(Redis)", "최대 180초", "백업에는 메시지를 넣지 않습니다", "본문 없는 승인 요청·전달 정보", "최대 약 1시간", "백업에 들어갑니다"}, "en": {"server memory (Redis)", "at most 180 seconds", "Backups do not include messages", "approval and delivery information without message text", "at most about 1 hour", "backups include the approval"}}[lang]...)
			}
			for _, want := range retention {
				if !strings.Contains(strings.ToLower(body), strings.ToLower(want)) {
					t.Fatalf("%s %s retention statement missing %q", page, lang, want)
				}
			}
			if strings.Contains(body, "받은 때부터 24시간") || strings.Contains(body, "24 hours after the Service receives") {
				t.Fatalf("%s %s still promises 24-hour message storage", page, lang)
			}
			if strings.Count(body, `aria-current="page"`) != 1 || strings.Count(body, "<h1") != 1 || strings.Contains(body, `href="#how"`) {
				t.Fatalf("%s %s: one current link, one h1 and no landing anchor menu", page, lang)
			}
			for _, banned := range []string{"{", "}", "<script", "<link", "<img", "<details", "<code", "<no value>"} {
				if strings.Contains(regexp.MustCompile(`(?s)<style>.*?</style>`).ReplaceAllString(body, ""), banned) {
					t.Fatalf("%s %s contains %q", page, lang, banned)
				}
			}
			if lang == "en" {
				// The only Korean left is the language switch name; the docs link says "(in Korean)".
				if left := koreanLeft(body); len(left) != 0 {
					t.Fatalf("%s en keeps Korean: %v", page, left)
				}
			}
			// PC-05: the privacy policy names calendars once, only to say they are not collected; the terms never do.
			want := map[string]int{"privacy": 1, "terms": 0}[page]
			if n := len(regexp.MustCompile(`(?i)calendar|busy|캘린더|바쁜 시간`).FindAllString(body, -1)); n != want {
				t.Fatalf("%s %s calendar mentions %d", page, lang, n)
			}
		}
	}
}

// The first page footer and sign-in section, and member settings, link to both pages (PC-03).
func TestPolicyLinks(t *testing.T) {
	s := &Service{Google: &GoogleLogin{}}
	w := httptest.NewRecorder()
	s.startPage(w, httptest.NewRequest("GET", "/", nil))
	body := w.Body.String()
	footer := body[strings.Index(body, "<footer>"):]
	if !regexp.MustCompile(`<footer><a href="/privacy">개인정보처리방침</a><a href="/terms">서비스 약관</a><a hreflang="ko" [^>]*>설치·사용 문서</a><a href="#main">맨 위로</a></footer>`).MatchString(footer) {
		t.Fatal("first page footer order", footer)
	}
	start := body[strings.Index(body, `id="start"`):strings.Index(body, "<footer>")]
	if !strings.Contains(start, `<p class="muted">로그인하기 전에 <a href="/terms">서비스 약관</a>과 <a href="/privacy">개인정보처리방침</a>을 확인하세요.</p>`) {
		t.Fatal("sign-in policy line", start)
	}
	en := renderEN(t, "start", s.startView(""))
	if dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR"); dir != "" {
		for name, page := range map[string]string{"start-ko": body, "start-en": en} {
			if err := os.WriteFile(filepath.Join(dir, "policy-"+name+".html"), []byte(page), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if !strings.Contains(en, `Before you sign in, read the <a href="/terms">Terms of Service</a> and the <a href="/privacy">Privacy Policy</a>.`) {
		t.Fatal("English sign-in policy line")
	}
	v := humanPageFixture()
	v["Title"], v["Section"] = "설정", "account"
	account := renderHumanPage(t, "human-account", v)
	if !strings.Contains(account, `<p class="links"><a href="/">서비스 안내</a> <a href="/privacy">개인정보처리방침</a> <a href="/terms">서비스 약관</a></p>`) {
		t.Fatal("settings policy links")
	}
}

// A language switch the browser would not store still explains itself on the policy page, as on every other page.
func TestPolicyPageUnsavedLanguageNotice(t *testing.T) {
	for cookie, want := range map[string]bool{"": true, "en": false} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/privacy?lang=en", nil)
		if cookie != "" {
			r.AddCookie(&http.Cookie{Name: langCookie, Value: cookie})
		}
		browserLang(w, r)
		(&Service{}).policyPage("privacy")(w, r)
		if body := w.Body.String(); !strings.Contains(body, `<html lang="en">`) || strings.Contains(body, "applies to this page only") != want {
			t.Fatalf("cookie %q: notice shown must be %v", cookie, want)
		}
	}
}
