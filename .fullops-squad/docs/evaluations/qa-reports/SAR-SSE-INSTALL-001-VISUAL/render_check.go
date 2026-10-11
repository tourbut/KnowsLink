// Verify that each saved synthetic fixture matches Go's real html/template output.
package main

import (
	"bytes"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	const sha = "febd1282df2e84a29fa0a53834569e7b9cc028e2"
	const out = ".fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL"
	source := func(path string) string {
		data, err := exec.Command("git", "show", sha+":"+path).Output()
		if err != nil {
			panic(err)
		}
		return string(data)
	}
	style := regexp.MustCompile("const memberStyle = `([^`]+)`").FindStringSubmatch(source("internal/relay/member.go"))[1]
	parts := regexp.MustCompile("var remotePage = .*?Parse\\(`([^`]+)` \\+ memberStyle \\+ `([^`]+)`\\)\\)").FindStringSubmatch(source("internal/relay/remote_auth.go"))
	page := template.Must(template.New("remote").Parse(parts[1] + style + parts[2]))
	for _, longURI := range []bool{false, true} {
		redirect := "https://client.example/oauth/callback"
		if longURI {
			redirect += "/" + strings.Repeat("syntheticcallbacksegment", 20) + "?fixture=long&review=visual"
		}
		for _, login := range []bool{true, false} {
			name := "consent"
			if login {
				name = "login"
			}
			if longURI {
				name += "-long"
			}
			var rendered bytes.Buffer
			err := page.Execute(&rendered, map[string]any{"Name": "Synthetic Work client", "Redirect": redirect,
				"ID": "synthetic-visual-attempt", "CSRF": "synthetic-not-a-credential", "Login": login})
			if err != nil {
				panic(err)
			}
			fixture, err := os.ReadFile(filepath.Join(out, name+".html"))
			if err != nil || !bytes.Equal(rendered.Bytes(), fixture) {
				panic("fixture differs from Go html/template: " + name)
			}
		}
	}
}
