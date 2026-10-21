// SAR-PUBLISH-001 PC-01–05: both policy pages in both languages, their links, and no calendar wording left.
package relay

import (
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
	if !strings.Contains(account, `<p><a href="/">서비스 안내</a> <a href="/privacy">개인정보처리방침</a> <a href="/terms">서비스 약관</a></p>`) {
		t.Fatal("settings policy links")
	}
}
