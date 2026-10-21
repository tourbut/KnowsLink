// Language checks: every Korean template run and handler string has English, the English screens show no Korean,
// composed text translates through patterns, and the browser choice resolves in the U3 order.
package relay

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

var templateAction = regexp.MustCompile(`{{[^}]*}}`)

// dumpMissing writes a map skeleton for the translators when KNOWSLINK_I18N_MISSING names a file.
func dumpMissing(t *testing.T, label string, missing []string) {
	t.Helper()
	if len(missing) == 0 {
		return
	}
	sort.Strings(missing)
	var b strings.Builder
	for _, m := range missing {
		b.WriteString("\t" + strconv.Quote(m) + ": \"\",\n")
	}
	if path := os.Getenv("KNOWSLINK_I18N_MISSING"); path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = f.WriteString("// " + label + "\n" + b.String())
		_ = f.Close()
	}
	t.Errorf("%s: %d untranslated\n%s", label, len(missing), b.String())
}

func TestTranslationCoversTemplates(t *testing.T) {
	_, missing := translateTemplate(memberSource, en)
	dumpMissing(t, "template runs", missing)
	// Every English value keeps the same template actions as its Korean key and carries no Korean but the language name.
	for key, value := range en {
		if hangul.MatchString(strings.ReplaceAll(templateAction.ReplaceAllString(value, ""), "한국어", "")) {
			t.Errorf("English value still Korean: %q", key)
		}
		want, got := templateAction.FindAllString(key, -1), templateAction.FindAllString(value, -1)
		sort.Strings(want)
		sort.Strings(got)
		if strings.Join(want, "|") != strings.Join(got, "|") {
			t.Errorf("template actions differ for %q: %v vs %v", key, want, got)
		}
		if strings.TrimSpace(value) == "" {
			t.Errorf("empty English for %q", key)
		}
	}
}

// Handler text that never reaches a template run still needs English: every Korean string literal in the product
// files is a key, a fragment of a pattern key, or an intentional exception (mail, logs, template sources).
func TestTranslationCoversGoText(t *testing.T) {
	// Mail text is not a web screen; the short time words are formatted by humanClock per language; %d분 is the %v분 pattern.
	exceptions := []string{"KnowsLink 확인 코드", "이 메일을 무시하세요", "1월 2일 ", "오늘 ", "오전 ", "오후 ", "%d분"}
	skip := func(s string) bool {
		for _, e := range exceptions {
			if strings.Contains(s, e) {
				return true
			}
		}
		return false
	}
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var missing []string
	seen := map[string]bool{}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") || strings.HasPrefix(file, "i18n") {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			s, err := strconv.Unquote(lit.Value)
			if err != nil || !hangul.MatchString(s) || strings.Contains(s, "{{") || skip(s) || seen[s] {
				return true
			}
			seen[s] = true
			if _, ok := en[s]; ok {
				return true
			}
			for key := range en {
				if strings.Contains(key, strings.TrimSpace(s)) {
					return true
				}
			}
			missing = append(missing, s)
			return true
		})
	}
	dumpMissing(t, "handler text", missing)
}

// asciiData keeps every system value (translatable) and replaces person-made Korean with ASCII, so any Korean left in
// the English page is a missing translation and never a name or message. It walks maps, slices and struct fields.
func asciiData(v any) any {
	return asciiValue(reflect.ValueOf(v)).Interface()
}

func asciiValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		return asciiValue(v.Elem())
	case reflect.String:
		s := v.String()
		if hangul.MatchString(s) && (utf8.RuneCountInString(s) < 2 || tr("en", s) == s) {
			return reflect.ValueOf("ascii-text").Convert(v.Type())
		}
		return v
	case reflect.Map:
		out := reflect.MakeMap(v.Type())
		for _, k := range v.MapKeys() {
			out.SetMapIndex(k, asciiValue(v.MapIndex(k)))
		}
		return out
	case reflect.Slice:
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			out.Index(i).Set(asciiValue(v.Index(i)))
		}
		return out
	case reflect.Pointer:
		if v.IsNil() {
			return v
		}
		out := reflect.New(v.Type().Elem())
		out.Elem().Set(asciiValue(v.Elem()))
		return out
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		out.Set(v)
		for i := 0; i < v.NumField(); i++ {
			if out.Field(i).CanSet() {
				out.Field(i).Set(asciiValue(v.Field(i)))
			}
		}
		return out
	}
	return v
}

func renderEN(t *testing.T, name string, v map[string]any) string {
	t.Helper()
	w := httptest.NewRecorder()
	w.Header().Set("Content-Language", "en")
	render(w, 200, name, v)
	return w.Body.String()
}

// koreanLeft lists the Korean runs of an English page with context, ignoring the language name itself.
func koreanLeft(body string) []string {
	var out []string
	for _, m := range regexp.MustCompile(`[^<>"]*[\x{AC00}-\x{D7A3}][^<>"]*`).FindAllString(strings.ReplaceAll(body, "한국어", ""), -1) {
		if s := strings.TrimSpace(m); s != "" {
			out = append(out, s)
		}
	}
	return out
}

var (
	firstPerson = regexp.MustCompile(`(?i)\bmy\b([^-]|$)`) // ids such as my-code are not text
	// Menu and heading names stay first person (V-LANG-02); sentences, choices and errors address the reader.
	ownNames = regexp.MustCompile(`My (agents|name|friend code)\b`)
	docLink  = regexp.MustCompile(`<a [^>]*href="https://github\.com/[^"]*README[^"]*"[^>]*>[^<]*</a>`)
)

func TestEnglishScreensHaveNoKorean(t *testing.T) {
	check := func(name string, body string) {
		t.Helper()
		if left := koreanLeft(body); len(left) > 0 {
			t.Errorf("%s: Korean left in English page: %q", name, left)
		}
		if !strings.Contains(body, `<html lang="en">`) {
			t.Errorf("%s: html lang must be en", name)
		}
		if strings.Contains(body, "<no value>") {
			t.Errorf("%s: <no value>", name)
		}
		if mine := firstPerson.FindAllString(ownNames.ReplaceAllString(body, ""), -1); len(mine) > 0 {
			t.Errorf("%s: first person outside menu names: %d", name, len(mine))
		}
		// A link to the Korean README says so in English (U3).
		for _, a := range docLink.FindAllString(body, -1) {
			if !strings.Contains(a, `hreflang="ko"`) || !strings.Contains(a, "(in Korean)") {
				t.Errorf("%s: Korean document link without a language mark: %s", name, a)
			}
		}
		if dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR"); dir != "" {
			if err := os.WriteFile(filepath.Join(dir, "en-"+name+".html"), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, f := range memberFixtures() {
		data, _ := asciiData(f.data).(map[string]any)
		name := f.tmpl
		if name == "gate" {
			name = "human-gate"
		}
		if f.data["Member"] == nil && (name == "verify" || name == "refusal" || name == "device") {
			data["Here"] = "/auth/verify"
		}
		if name == "receipt" {
			data["Own"] = []humanAgent{{ID: "agent_a", Name: "Synthetic helper", Client: "Node 22", Status: "등록됨", Short: "a001"}} // the fixture's own agent is named like a heading
		}
		check(f.name, renderEN(t, name, data))
	}
	human := func() map[string]any {
		v := asciiData(humanPageFixture()).(map[string]any)
		v["Problem"], v["Notice"], v["Lang"] = "이름은 1–40자로 입력하세요", "요청을 보냈습니다. 가능한 요청은 상대가 확인할 수 있습니다.", "en"
		return v
	}
	for _, name := range []string{"human-home", "human-friends", "human-friend", "human-account", "human-name", "human-new", "human-conversation"} {
		check(name, renderEN(t, name, human()))
	}
	for _, calendar := range []string{"none", "connected", "reconnect"} {
		v := human()
		v["Title"], v["Section"], v["Calendar"], v["CalendarGoogle"], v["CalendarReady"] = "설정", "account", calendar, true, true
		check("human-account-calendar-"+calendar, renderEN(t, "human-account", v))
	}
	for _, choice := range append(useChoices, useChoice{ID: ""}) {
		for _, start := range []string{"first", "continue"} {
			v := human()
			v["Title"], v["Section"], v["Use"], v["Start"] = "내 에이전트", "agents", choice.ID, start
			v["Own"] = []map[string]any{{"ID": "agent_a", "Name": "", "Client": "", "Status": "연결 필요", "Short": "a001"}, {"ID": "agent_b", "Name": "Synthetic remote helper", "Client": "원격 호스트", "Status": "등록됨", "Short": "b002"}}
			check("human-agents-"+choice.ID+"-"+start, renderEN(t, "human-agents", v))
		}
	}
	check("remote-return", renderEN(t, "remote-return", map[string]any{"Target": "https://app.invalid/cb?code=x"}))
	check("remote-problem", renderEN(t, "remote", map[string]any{"Name": "client", "Redirect": "https://app.invalid/cb", "ID": "att", "CSRF": "c", "AgentName": "x", "Problem": "이름은 1–40자로 입력하세요"}))
	check("google-complete", renderEN(t, "google-complete", map[string]any{"Title": "Google 로그인 확인 완료", "Notice": "신원을 확인했습니다. 계속해서 연결을 확인하세요.", "Next": "/home"}))
}

// Every reachable connection resume state and the composed handler strings translate with their captures intact.
func TestComposedHandlerTextTranslates(t *testing.T) {
	st, now := newState(), time.Now().UTC()
	st.Owners["o"] = &Owner{Active: true}
	st.Agents["a"] = &Agent{Owner: "o", Keys: map[string]*Key{}}
	states := []string{}
	for _, withKey := range []bool{false, true} {
		if withKey {
			st.Agents["a"].Keys["k"] = &Key{Public: []byte{1}}
		}
		for _, state := range []string{"waiting", "prepared", "approved", "expired", "cancelled"} {
			st.Connections = map[string]*Connection{"c": {Owner: "o", Agent: "a", State: state, Exp: now.Add(time.Minute)}}
			if state == "expired" || state == "cancelled" {
				st.Connections["c"].Exp = now.Add(-time.Minute)
			}
			states = append(states, st.agentResume("a", now).State)
		}
	}
	st.Connections = map[string]*Connection{}
	states = append(states, st.agentResume("a", now).State, deviceUnsure.State)
	for _, s := range states {
		if got := tr("en", s); got == s || hangul.MatchString(got) {
			t.Errorf("resume state untranslated: %q -> %q", s, got)
		}
	}
	at := time.Date(2026, 10, 20, 3, 4, 5, 0, time.UTC)
	for _, s := range []string{limited(at), "코드가 맞지 않습니다. 남은 시도: 3회.", keyRecordsFull + " 지금은 새 agent도 만들 수 없습니다. " + agentProblem("active"), "입력한 코드 …ab12", "Alice — 친구", "상태를 확인할 수 없습니다. 잠시 뒤 같은 화면에서 다시 조회하세요." + restoreNote} {
		got := tr("en", s)
		if got == s || hangul.MatchString(got) {
			t.Errorf("composed text untranslated: %q -> %q", s, got)
		}
		for _, keep := range []string{clock(at), "3", "ab12", "Alice"} {
			if strings.Contains(s, keep) && !strings.Contains(got, keep) {
				t.Errorf("capture lost in %q", got)
			}
		}
	}
	if tr("ko", "아무 문장") != "아무 문장" || tr("en", "번역표에 없는 문장") != "번역표에 없는 문장" || tr("en", "plain ascii") != "plain ascii" {
		t.Fatal("unknown text must fall back unchanged")
	}
	if tr("en", "Alice") != "Alice" || agentFallback("en", "친구") != "친구" || agentFallback("en", "이름 없는 에이전트") == "이름 없는 에이전트" {
		t.Fatal("names pass through; only the server labels translate")
	}
}

func TestHumanClockKeepsInstantAcrossLanguages(t *testing.T) {
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, kst)
	for _, c := range []struct {
		at     time.Time
		ko, en string
	}{
		{time.Date(2026, 10, 20, 14, 10, 0, 0, kst), "오늘 오후 2:10", "Today 2:10 PM"},
		{time.Date(2026, 10, 2, 0, 5, 0, 0, kst), "10월 2일 오전 12:05", "Oct 2, 12:05 AM"},
		{time.Date(2025, 1, 2, 9, 30, 0, 0, kst), "2025년 1월 2일 오전 9:30", "Jan 2, 2025, 9:30 AM"},
	} {
		if got := humanClock(c.at, now, "ko"); got != c.ko {
			t.Fatal(got)
		}
		if got := humanClock(c.at, now, "en"); got != c.en {
			t.Fatal(got)
		}
	}
	if clock(now) != "2026-10-20 12:00:00 KST" {
		t.Fatal("exact deadlines keep KST in both languages")
	}
}

func TestLanguageChoiceOrderAndLink(t *testing.T) {
	req := func(cookie, query string) *http.Request {
		r := httptest.NewRequest("GET", "/"+query, nil)
		if cookie != "" {
			r.AddCookie(&http.Cookie{Name: langCookie, Value: cookie})
		}
		return r
	}
	for _, c := range []struct {
		member              *Member
		cookie, query, want string
	}{
		{nil, "", "", "ko"}, {nil, "en", "", "en"}, {nil, "xx", "", "ko"}, {nil, "", "?lang=en", "en"}, {nil, "ko", "?lang=en", "ko"},
		{&Member{Locale: "en"}, "ko", "", "en"}, {&Member{Locale: "ko"}, "en", "", "ko"}, {&Member{}, "en", "", "en"}, {&Member{Locale: "zz"}, "", "", "ko"}, {&Member{Locale: "zz"}, "en", "", "en"},
	} {
		if got := pageLang(req(c.cookie, c.query), c.member); got != c.want {
			t.Fatalf("%+v cookie %q query %q: %q", c.member, c.cookie, c.query, got)
		}
	}
	for raw, want := range map[string]string{"/auth/verify": "/auth/verify", "/remote/consent/x?y=1": "/remote/consent/x?y=1", "https://evil.invalid/": "/", "//evil.invalid": "/", "home": "/", "": "/", "/a\\b": "/"} {
		if got := safeBack(raw); got != want {
			t.Fatalf("safeBack(%q) = %q", raw, got)
		}
	}
	s := &Service{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /lang/{code}", s.setLang)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/lang/en?back=%2Fauth%2Fverify", nil))
	if w.Code != 303 || w.Header().Get("Location") != "/auth/verify?lang=en" {
		t.Fatal(w.Code, w.Header().Get("Location"))
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != langCookie || cookies[0].Value != "en" || !cookies[0].Secure || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode || cookies[0].MaxAge != langMaxAge {
		t.Fatalf("cookie %v", cookies)
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/lang/ko?back=https://evil.invalid/", nil))
	if w.Header().Get("Location") != "/?lang=ko" {
		t.Fatal(w.Header().Get("Location"))
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/lang/fr", nil))
	if w.Code != 404 || len(w.Result().Cookies()) != 0 {
		t.Fatal("unknown language must not set a cookie")
	}
	// Without a stored cookie the ?lang= hint applies to this response only and the page says so.
	w = httptest.NewRecorder()
	browserLang(w, httptest.NewRequest("GET", "/?lang=en", nil))
	render(w, 200, "start", (&Service{}).startView(""))
	if body := w.Body.String(); !strings.Contains(body, `<html lang="en">`) || !strings.Contains(body, "only") || w.Header().Get(langUnsavedHeader) != "" {
		t.Fatal("one-response language without the saved-cookie notice", w.Header())
	}
	w = httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/?lang=en", nil)
	r.AddCookie(&http.Cookie{Name: langCookie, Value: "en"})
	browserLang(w, r)
	render(w, 200, "start", (&Service{}).startView(""))
	if body := w.Body.String(); strings.Contains(body, `role="status"`) {
		t.Fatal("stored cookie must not show the unsaved notice")
	}
}

// A member whose saved value decides the language sees the settings link instead of browser links that would do
// nothing; members without a saved value and visitors keep the links back to the same page (V-LANG-01).
func TestSavedLanguageReplacesBrowserLinks(t *testing.T) {
	page := func(name string, v map[string]any, lang string, saved bool) (string, http.Header) {
		w := httptest.NewRecorder()
		w.Header().Set("Content-Language", lang)
		if saved {
			w.Header().Set(langSavedHeader, "1")
		}
		render(w, 200, name, v)
		return w.Body.String(), w.Header()
	}
	views := map[string]func() map[string]any{
		"start": func() map[string]any { v := (&Service{}).startView(""); v["Authenticated"] = true; return v },
		"remote": func() map[string]any {
			return map[string]any{"ID": "att", "Here": "/remote/consent/att", "Name": "client"}
		},
	}
	for name, view := range views {
		back := "%2f" // html/template escapes the path inside the query
		if name == "remote" {
			back = "%2fremote%2fconsent%2fatt"
		}
		for lang, text := range map[string]string{"ko": "설정에서 변경", "en": "Change in Settings"} {
			body, header := page(name, view(), lang, true)
			if !strings.Contains(body, `<a href="/home/account#language">`+text+`</a>`) || strings.Contains(body, "/lang/") || header.Get(langSavedHeader) != "" {
				t.Fatalf("%s %s saved member: %s", name, lang, header)
			}
			body, _ = page(name, view(), lang, false)
			if !strings.Contains(body, `/lang/ko?back=`+back+`"`) || !strings.Contains(body, `/lang/en?back=`+back+`"`) || strings.Contains(body, "#language") {
				t.Fatalf("%s %s without a saved value must keep the browser links", name, lang)
			}
		}
	}
}

// Admission refuses before any member lookup, so its 429 and 503 screens follow the browser choice (U3 LANG-VISUAL).
func TestAdmissionRefusalsFollowBrowserLanguage(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://x@127.0.0.1:1/x?connect_timeout=1") // nothing listens
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := &Service{Pool: pool}
	h := s.boundedHTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Fatal("refused request reached the handler") }))
	check := func(path, lang string, status int) {
		t.Helper()
		r := httptest.NewRequest("GET", path, nil)
		if lang == "en" {
			r.AddCookie(&http.Cookie{Name: langCookie, Value: "en"})
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		body := w.Body.String()
		if w.Code != status || w.Header().Get("Content-Language") != lang || !strings.Contains(body, `<html lang="`+lang+`">`) {
			t.Fatalf("%s %s: %d %q", path, lang, w.Code, w.Header().Get("Content-Language"))
		}
		if left := koreanLeft(body); lang == "en" && len(left) > 0 {
			t.Fatalf("%s: Korean left in English refusal: %q", path, left)
		}
	}
	for _, lang := range []string{"ko", "en"} {
		check("/home/receipts", lang, 503)
	}
	for range cap(s.httpNew) {
		s.httpNew <- struct{}{}
	}
	for _, path := range []string{"/", "/home/receipts"} {
		for _, lang := range []string{"ko", "en"} {
			check(path, lang, 429)
		}
	}
}

// O-LG-05: the other language's own name inside a problem sentence carries its language (WCAG 3.1.2), and the
// internal language markers leave no response, rendered or a plain redirect.
func TestLanguagePartsAndMarkers(t *testing.T) {
	problem := "지원하지 않는 언어입니다. 저장된 설정은 그대로입니다. 한국어 또는 English를 선택하세요."
	for lang, part := range map[string]string{"ko": `<span lang="en">English</span>를`, "en": `Choose <span lang="ko">한국어</span> or English.`} {
		w := httptest.NewRecorder()
		w.Header().Set("Content-Language", lang)
		render(w, 422, "refusal", map[string]any{"Title": "설정", "Problem": problem})
		if !strings.Contains(w.Body.String(), part) {
			t.Fatalf("%s problem lacks the language part %q", lang, part)
		}
	}
	w := httptest.NewRecorder()
	marked := langMarks{w}
	marked.Header().Set(langSavedHeader, "1")
	marked.Header().Set(langUnsavedHeader, "1")
	http.Redirect(marked, httptest.NewRequest("POST", "/home/language", nil), "/home/account", http.StatusSeeOther)
	if sent := w.Result().Header; w.Code != 303 || sent.Get(langSavedHeader) != "" || sent.Get(langUnsavedHeader) != "" {
		t.Fatal("redirect carried an internal language marker", sent)
	}
}
