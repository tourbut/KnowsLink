// Language choice (U3): a member's saved locale, else the browser's explicit choice, else Korean. The Korean text in
// templates and handlers is the key; the English table maps it and anything missing falls back to the Korean original.
package relay

import (
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

const (
	langCookie        = "__Host-kl_lang"
	langMaxAge        = 365 * 24 * 60 * 60
	langUnsavedHeader = "X-KnowsLink-Lang-Unsaved" // set for one response when ?lang= applied without a stored cookie; render removes it
)

var hangul = regexp.MustCompile(`[\x{AC00}-\x{D7A3}\x{3131}-\x{318E}]`)

func validLang(l string) bool { return l == "ko" || l == "en" }

// pageLang follows the U3 priority table: member value, explicit browser cookie, one-response ?lang= hint, Korean.
// An unknown or damaged value is never shown; it simply loses to the next source.
func pageLang(r *http.Request, m *Member) string {
	if m != nil && validLang(m.Locale) {
		return m.Locale
	}
	if l := readCookie(r, langCookie); validLang(l) {
		return l
	}
	if l := r.URL.Query().Get("lang"); validLang(l) {
		return l
	}
	return "ko"
}

// browserLang records the pre-login language on the response so every page and refusal of this request can read it.
func browserLang(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/v1/") || r.URL.Path == "/mcp" {
		return
	}
	w.Header().Set("Content-Language", pageLang(r, nil))
	if !validLang(readCookie(r, langCookie)) && validLang(r.URL.Query().Get("lang")) {
		w.Header().Set(langUnsavedHeader, "1")
	}
}

func responseLang(w http.ResponseWriter) string {
	if l := w.Header().Get("Content-Language"); validLang(l) {
		return l
	}
	return "ko"
}

// safeBack keeps the language link inside this site: a path only, never a scheme, host or protocol-relative URL.
func safeBack(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "" || u.Host != "" || !strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") || strings.ContainsAny(raw, "\\\r\n") {
		return "/"
	}
	return raw
}

// setLang stores the anonymous choice in one cookie and returns to the page it came from. The ?lang= hint on the
// return lets that page apply the choice even when the browser refuses the cookie, and say so.
func (s *Service) setLang(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if !validLang(code) {
		http.NotFound(w, r)
		return
	}
	setCookie(w, langCookie, code, langMaxAge)
	back, _ := url.Parse(safeBack(r.URL.Query().Get("back")))
	q := back.Query()
	q.Set("lang", code)
	back.RawQuery = q.Encode()
	http.Redirect(w, r, back.String(), http.StatusSeeOther)
}

// tr translates one handler string. A key with %s or %v is a pattern for composed text: %s captures are translated
// again, %v captures (names, times, numbers) stay verbatim. Unknown text returns unchanged.
func tr(lang, s string) string {
	if lang != "en" || s == "" || !hangul.MatchString(s) {
		return s
	}
	if v, ok := en[s]; ok {
		return v
	}
	for _, p := range enPatterns() {
		m := p.re.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		out := p.value
		for i, verb := range p.verbs {
			capture := m[i+1]
			if verb == 's' {
				capture = tr(lang, capture)
			}
			out = strings.Replace(out, "%"+string(verb), capture, 1)
		}
		return out
	}
	return s
}

type langPattern struct {
	re    *regexp.Regexp
	value string
	verbs []byte
}

var (
	patternOnce sync.Once
	patternList []langPattern
)

func enPatterns() []langPattern {
	patternOnce.Do(func() {
		for key, value := range en {
			if !strings.Contains(key, "%s") && !strings.Contains(key, "%v") {
				continue
			}
			expr, verbs := "^", []byte{}
			for len(key) > 0 {
				i := strings.IndexAny(key, "%")
				if i < 0 || i+1 >= len(key) || (key[i+1] != 's' && key[i+1] != 'v') {
					expr += regexp.QuoteMeta(key)
					break
				}
				expr += regexp.QuoteMeta(key[:i]) + "(.*?)"
				verbs = append(verbs, key[i+1])
				key = key[i+2:]
			}
			patternList = append(patternList, langPattern{regexp.MustCompile(expr + "$"), value, verbs})
		}
	})
	return patternList
}

// badges map a stored state to its tone, mark and label in both languages; the same Korean word can be a button
// elsewhere (거절 = Decline), so badge labels keep their own English column instead of the shared table.
var badges = map[string][4]string{
	"연결 완료": {"ok", "✓", "연결 완료", "Connected"}, "활성": {"ok", "✓", "활성", "Active"}, "미연결": {"warn", "…", "미연결", "Not connected"}, "연결 대기": {"warn", "…", "연결 대기", "Awaiting connection"}, "철회": {"off", "—", "철회", "Revoked"},
	"pending": {"warn", "…", "대기", "Pending"}, "active": {"ok", "✓", "연결됨", "Connected"}, "approved": {"ok", "✓", "승인", "Approved"}, "denied": {"off", "—", "거절", "Denied"}, "expired": {"off", "—", "만료", "Expired"}, "revoked": {"off", "—", "철회", "Revoked"},
	"queued": {"warn", "…", "대기", "Pending"}, "leased": {"warn", "…", "대기", "Pending"}, "delivered": {"ok", "✓", "전달됨", "Delivered"}, "completed": {"ok", "✓", "처리 완료", "Completed"}, "failed": {"off", "!", "실패", "Failed"},
	"cancelled": {"off", "—", "취소", "Cancelled"}, "removed": {"off", "—", "해제", "Removed"}, "blocked": {"off", "—", "차단", "Blocked"}, "registered": {"ok", "✓", "등록됨", "Registered"}, "unconnected": {"off", "—", "연결 필요", "Needs connection"}, "등록됨": {"ok", "✓", "등록됨", "Registered"}, "연결 필요": {"off", "—", "연결 필요", "Needs connection"},
	"unavailable": {"off", "!", "확인할 수 없음", "Cannot be checked"},
}

func badge(lang, state string) template.HTML {
	b, ok := badges[state]
	if !ok {
		b = [4]string{"off", "—", state, state}
	}
	label := b[2]
	if lang == "en" {
		label = b[3]
	}
	return template.HTML(`<span class="badge ` + b[0] + `"><span aria-hidden="true">` + b[1] + `</span> ` + template.HTMLEscapeString(label) + `</span>`)
}

// agentFallback translates only the two server-made agent labels; a person's own name is never translated.
func agentFallback(lang, name string) string {
	if name == "이름 없는 에이전트" || name == "확인할 수 없는 에이전트" {
		return tr(lang, name)
	}
	return name
}

// pageFuncs take any value: a view map without the key renders as empty text instead of stopping the page.
func pageFuncs(lang string) template.FuncMap {
	text := func(v any) string { s, _ := v.(string); return s }
	return template.FuncMap{
		"badge": func(state string) template.HTML { return badge(lang, state) },
		"tr":    func(v any) string { return tr(lang, text(v)) },
		"name":  func(v any) string { return agentFallback(lang, text(v)) },
		"lang":  func() string { return lang },
	}
}

// segment is one run of template text: content between tags, or one attribute value. Inline actions such as
// {{.Name}} stay inside the run so the English order can move around them; control actions split runs.
type segment struct {
	start, end int
	text       string
}

func controlAction(action string) bool {
	body := strings.TrimLeft(strings.TrimPrefix(action, "{{"), "- ")
	for _, word := range []string{"if ", "else", "end", "range ", "with ", "define ", "block "} {
		if strings.HasPrefix(body, word) {
			return true
		}
	}
	return strings.HasPrefix(body, "$") && strings.Contains(body, ":=")
}

func templateSegments(src string) []segment {
	var out []segment
	start, inTag, inAttr := 0, false, false
	flush := func(end int) {
		text := src[start:end]
		if hangul.MatchString(text) {
			lead, trail := len(text)-len(strings.TrimLeft(text, " \t\n")), len(text)-len(strings.TrimRight(text, " \t\n"))
			out = append(out, segment{start + lead, end - trail, strings.TrimSpace(text)})
		}
		start = end
	}
	for i := 0; i < len(src); {
		switch {
		case strings.HasPrefix(src[i:], "{{"):
			end := strings.Index(src[i:], "}}")
			if end < 0 {
				end = len(src) - i - 2
			}
			if controlAction(src[i : i+end+2]) {
				flush(i)
				start = i + end + 2
			}
			i += end + 2
		case !inTag && src[i] == '<':
			flush(i)
			inTag = true
			i++
			start = i
		case inTag && !inAttr && src[i] == '>':
			inTag = false
			i++
			start = i
		case inTag && src[i] == '"':
			if inAttr {
				flush(i)
			}
			inAttr = !inAttr
			i++
			start = i
		default:
			i++
		}
	}
	flush(len(src))
	return out
}

// translateTemplate rewrites every Korean run of a template source from the table; untranslated runs stay Korean
// and are reported so the coverage test can list them.
func translateTemplate(src string, table map[string]string) (string, []string) {
	segs := templateSegments(src)
	var missing []string
	var b strings.Builder
	last := 0
	for _, seg := range segs {
		b.WriteString(src[last:seg.start])
		if v, ok := table[seg.text]; ok {
			b.WriteString(v)
		} else {
			b.WriteString(seg.text)
			missing = append(missing, seg.text)
		}
		last = seg.end
	}
	b.WriteString(src[last:])
	return b.String(), missing
}

func parsePages(lang string) *template.Template {
	src := memberSource
	if lang == "en" {
		src, _ = translateTemplate(memberSource, en)
	}
	return template.Must(template.New("member").Funcs(pageFuncs(lang)).Parse(src))
}

var en = func() map[string]string {
	m := map[string]string{"한국어": "한국어"}
	for _, part := range []map[string]string{enMember, enHuman} {
		for k, v := range part {
			m[k] = v
		}
	}
	return m
}()
