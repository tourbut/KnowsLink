// Public privacy policy and terms (SAR-PUBLISH-001, D017). Each language's body is its own static document, because
// the two texts are written per language rather than translated run by run; the page chrome uses the shared table.
package relay

import (
	"embed"
	"html/template"
	"net/http"
	"time"
)

// policyEffective is the effective date both languages show (P-01); ops changes only this value if the deploy moves.
const policyEffective = "2026-10-22"

//go:embed policy/*.html
var policyFiles embed.FS

// policyDocs holds the trusted static bodies by page and language, read once at start.
var policyDocs = func() map[string]template.HTML {
	docs := map[string]template.HTML{}
	for _, page := range []string{"privacy", "terms"} {
		for _, lang := range []string{"ko", "en"} {
			raw, err := policyFiles.ReadFile("policy/" + page + "." + lang + ".html")
			if err != nil {
				panic(err)
			}
			docs[page+"."+lang] = template.HTML(raw)
		}
	}
	return docs
}()

var policyTitles = map[string]string{"privacy": "개인정보처리방침", "terms": "서비스 약관"}

const policyTemplates = `
{{define "policy"}}<!doctype html><html lang="{{lang}}"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>{{.Title}} · KnowsLink</title>` + memberStyle + `{{template "landing-style"}}
<style>.policy{max-width:72ch}.policy section{margin-bottom:32px}.policy h2{margin-bottom:12px}.policy li{margin-bottom:8px}.policy .table-scroll{max-width:100%;overflow-x:auto;margin:0 0 16px}.policy table{border-collapse:collapse;min-width:100%}.policy th,.policy td{padding:8px 16px 8px 0;border-top:1px solid var(--edge);text-align:left;vertical-align:top}.policy th{font-weight:600}</style><body class="landing">
<a class="skip" href="#main">본문으로 건너뛰기</a>
<header class="top"><a class="brand" href="/">KnowsLink</a>{{template "lang-nav" .}}</header>
<main id="main" tabindex="-1" class="policy"><h1>{{.Title}}</h1><p class="muted">시행일: {{.Effective}}</p>
{{.Doc}}</main>{{template "landing-footer" .}}</body></html>{{end}}
`

// policyPage serves one policy page to anyone; a signed-in member's saved language applies as on the first page.
func (s *Service) policyPage(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token := readCookie(r, sessionCookie); token != "" {
			// A failed lookup only loses the saved language; the page itself never depends on the database.
			_, _ = s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
				id, _, err := st.session(token, now)
				if err == nil {
					memberLang(w, r, st.Members[id])
				}
				return nil, err
			})
		}
		lang := responseLang(w)
		render(w, 200, "policy", map[string]any{"Title": policyTitles[page], "Page": page, "Here": "/" + page, "Effective": policyEffective, "Doc": policyDocs[page+"."+lang]})
	}
}
