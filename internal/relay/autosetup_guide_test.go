// SAR-AUTOSETUP-001 AC-10: the setup command appears only for Claude Code and Codex, carries the page language and the
// released plugin version, and the guide text matches D04 in Korean and English.
package relay

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"os"
	"strings"
	"testing"
)

func TestAutosetupGuideOnlyForClaudeAndCodex(t *testing.T) {
	var plugin struct{ Version string }
	raw, err := os.ReadFile("../../adapters/.claude-plugin/plugin.json")
	if err != nil || json.Unmarshal(raw, &plugin) != nil || plugin.Version == "" {
		t.Fatal("plugin version", err)
	}
	command := "& ([scriptblock]::Create((irm https://link.knowslog.com/install/" + plugin.Version + "/knowslink-setup.ps1))) -App "
	want := map[string]map[string][]string{
		"ko": {
			"claude": {"자동 설정 · 시험 환경 확인 (Claude Code 2.1.302, Windows 11)", "이 컴퓨터의 PowerShell에서 실행할 명령", "메시지를 받을 세션에서 <code>/knowslink-session on</code>, 보낼 때마다 발송 승인.", "수동 설치 문서"},
			"codex":  {"자동 설정 · 시험 환경 확인 (Codex CLI 0.163.0, Windows 11) · 도착 알림 없음", "Codex 승인 모드는 바꾸지 않습니다.", "(예: <code>approval_policy = &#34;on-request&#34;</code>)로 Codex를 열고", "Codex 화면에는 도착 알림이 없습니다."},
			"other":  {"Hermes Agent, 그 밖의 호스트", "아직 연결 방법을 확인하는 중입니다."},
		},
		"en": {
			"claude": {"Automatic setup · Checked in a test environment (Claude Code 2.1.302, Windows 11)", "Command to run in PowerShell on this computer", "run <code>/knowslink-session on</code> in the session that should receive; approve every send.", "Manual install guide (in Korean)"},
			"codex":  {"Automatic setup · Checked in a test environment (Codex CLI 0.163.0, Windows 11) · No arrival notice", "It does not change Codex's approval mode.", "(for example <code>approval_policy = &#34;on-request&#34;</code>) and connect the session.", "Codex shows no arrival notice"},
			"other":  {"Hermes Agent and other hosts", "The connection method is still being confirmed."},
		},
	}
	for lang, pages := range map[string]*template.Template{"ko": memberPages, "en": memberPagesEN} {
		for _, use := range []string{"web", "grok", "dots", "claude", "codex", "other"} {
			for _, start := range []string{"first", "continue"} {
				v := humanPageFixture()
				v["Use"], v["Start"] = use, start
				var b bytes.Buffer
				if err := pages.ExecuteTemplate(&b, "human-agents", v); err != nil {
					t.Fatal(lang, use, err)
				}
				body := html.UnescapeString(b.String())
				has := strings.Contains(body, command)
				if (use == "claude" || use == "codex") != has {
					t.Fatalf("%s %s: command shown=%v", lang, use, has)
				}
				if has && !strings.Contains(body, command+use+" -Lang "+lang+"</textarea>") {
					t.Fatalf("%s %s: command must carry app and page language", lang, use)
				}
				if has && !strings.Contains(body, `id="setup-command" class="share" readonly rows="4"`) {
					t.Fatalf("%s %s: setup command box must show all 4 lines at 320px", lang, use)
				}
				texts := want[lang][use]
				if start == "continue" && has {
					texts = []string{map[string]string{"ko": "같은 명령을 실행하면 기존 설치를 갱신합니다. 에이전트·키·친구는 그대로입니다.", "en": "Running the same command updates your existing install. Your agent, key, and friends stay the same."}[lang]}
				}
				for _, text := range texts {
					if !strings.Contains(b.String(), text) && !strings.Contains(body, text) {
						t.Fatalf("%s %s: missing %q", lang, use, text)
					}
				}
				if lang == "en" && strings.Contains(body, "<h3>") && hangul.MatchString(body[strings.Index(body, `id="use-guide"`):strings.Index(body, "</section>")]) {
					t.Fatalf("%s %s: untranslated guide text", lang, use)
				}
			}
		}
	}
}
