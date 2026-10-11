// Member pages run verified login, sessions, logout, and the self-owner home without owner Basic auth or SSH.
package relay

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var kst = time.FixedZone("KST", 9*60*60)

// memberStyle is the single source of the member web tokens (D04 service-member-ui.md §8). CSP allows only inline
// styles, so there are no fonts, images, or scripts; color, type, spacing, buttons, and badges all live here.
const memberStyle = `<style>
:root{--bg:#f4f6f8;--surface:#fff;--text:#17232b;--muted:#44525c;--line:#c9d2d9;--line-strong:#6b7a84;--accent:#164da0;--accent-soft:#e8effb;--danger:#a01616;--danger-soft:#fdecec;--ok:#146c43;--ok-soft:#e7f5ee;--warn:#7a4b00;--warn-soft:#fff4d6;--mono:ui-monospace,SFMono-Regular,Consolas,"Courier New",monospace;color-scheme:light}
*,::before,::after{box-sizing:border-box}
body{max-width:1072px;margin:0 auto;padding:16px 16px 48px;background:var(--bg);color:var(--text);font:16px/1.65 system-ui,-apple-system,"Segoe UI","Apple SD Gothic Neo","Malgun Gothic","Noto Sans KR",sans-serif;word-break:keep-all;overflow-wrap:break-word;font-variant-numeric:tabular-nums}
h1{font-size:26px;line-height:1.3;font-weight:700;margin:8px 0 16px}h2{font-size:20px;line-height:1.3;font-weight:600;margin:0 0 12px}h3,h4{font-size:16px;line-height:1.3;font-weight:600;margin:0}h4,section>h3,section>h2{margin:24px 0 8px}section>h2:first-child{margin-top:0}
p{margin:0 0 12px}ul,ol{margin:0 0 12px;padding-left:24px}a{color:var(--accent)}
code,pre{font:14px/1.6 var(--mono)}code{word-break:normal;overflow-wrap:anywhere}
pre{white-space:pre-wrap;overflow-wrap:anywhere;background:var(--bg);border:1px solid var(--line);border-radius:6px;padding:12px 16px;margin:4px 0 16px}
.skip{position:absolute;left:16px;top:-80px;background:var(--surface);padding:8px 12px;border-radius:6px}.skip:focus{top:8px}
.top{display:flex;flex-wrap:wrap;align-items:center;gap:4px 16px;padding-bottom:12px;margin-bottom:16px;border-bottom:1px solid var(--line)}
.brand{font-weight:700;font-size:18px;margin-right:auto}.who{color:var(--muted);font-size:14px}
.page{max-width:800px;margin:0 auto}
.toc{display:flex;flex-wrap:wrap;gap:0 16px;margin:0 0 16px}.toc a{display:inline-flex;align-items:center;min-height:44px}
@media(min-width:960px){.layout{display:grid;grid-template-columns:200px minmax(0,1fr);gap:32px}.toc{position:sticky;top:16px;align-self:start;flex-direction:column}}
section{background:var(--surface);border:1px solid var(--line);border-radius:8px;padding:24px;margin:0 0 24px}
@media(max-width:600px){section{padding:16px}}
.sechead{display:flex;flex-wrap:wrap;align-items:center;gap:8px 16px;margin-bottom:8px}.sechead h2{margin:0 auto 0 0}
.muted{color:var(--muted);font-size:14px}
label{display:block;font-weight:600;margin:16px 0 4px}
input,select{display:block;width:100%;min-height:44px;margin:4px 0 12px;padding:8px 12px;border:1px solid var(--line-strong);border-radius:6px;background:var(--surface);color:inherit;font:inherit;font-weight:400}
input.code{font:20px var(--mono);letter-spacing:.2em;max-width:12em}
button,.button{display:inline-flex;align-items:center;justify-content:center;min-height:44px;margin:4px 8px 4px 0;padding:8px 16px;border:1px solid var(--line-strong);border-radius:6px;background:var(--surface);color:var(--text);font:inherit;text-align:center;text-decoration:none;cursor:pointer}
.primary{background:var(--accent);border-color:var(--accent);color:#fff;font-weight:600;min-width:240px}
.danger-button{border-color:var(--danger);color:var(--danger)}
.link{border:0;background:none;color:var(--accent);text-decoration:underline;padding:8px 4px;margin:0}
@media(max-width:600px){.primary{width:100%;min-width:0}}
.top form,.inline{display:inline}
.badge{display:inline-block;padding:0 8px;border-radius:4px;font-size:14px;font-weight:600;white-space:nowrap}
.ok{color:var(--ok);background:var(--ok-soft)}.warn{color:var(--warn);background:var(--warn-soft)}.off{color:var(--muted);background:var(--bg);border:1px solid var(--line)}
.note,.problem,.todo{margin:0 0 16px;padding:12px 16px;border-left:6px solid var(--accent);background:var(--accent-soft);border-radius:0 6px 6px 0}
.problem,.caution{border-left-color:var(--danger);background:var(--danger-soft)}
.todo ul{margin:4px 0 0}.note p:last-child,.problem p:last-child{margin-bottom:0}
.row{border-top:1px solid var(--line);padding:12px 0}.rows{list-style:none;padding:0;margin:0}
summary{display:flex;flex-wrap:wrap;align-items:center;gap:4px 12px;min-height:44px;cursor:pointer;list-style:none}
summary::-webkit-details-marker{display:none}summary::before{content:"▸";color:var(--muted);width:12px}details[open]>summary::before{content:"▾"}
summary h3{flex:1 1 16em}.agent{padding:0}.agent>summary{padding:12px 0}.agent>div{padding:0 0 16px 24px}
@media(max-width:600px){.agent>div{padding-left:0}}
.keys{list-style:none;padding:0}.keys li{padding:8px 0;border-top:1px solid var(--line)}
.fp{display:block;font-size:14px;letter-spacing:.02em;user-select:all}
.fp.big{font-size:16px;font-weight:600;padding:12px 16px;margin:8px 0;background:var(--bg);border:1px solid var(--line-strong);border-radius:6px}
.select{user-select:all}.value{display:block;padding:8px 12px;margin:4px 0 12px;background:var(--bg);border:1px solid var(--line-strong);border-radius:6px}
.danger{margin-top:16px;border:1px solid var(--danger);border-radius:6px;padding:0 16px}.danger>summary{color:var(--danger);font-weight:600}.danger>p:first-of-type{margin-top:4px}
.danger :focus-visible{outline-color:var(--danger)}
.steps{list-style:none;padding:0;display:grid;gap:8px;margin:0 0 24px}.steps li{padding:8px 12px;border:1px solid var(--line);border-radius:6px;background:var(--surface)}
.steps li[aria-current]{border:2px solid var(--accent);font-weight:600}.steps .n{display:inline-block;min-width:24px;font-weight:700}
.kv{display:grid;grid-template-columns:minmax(8em,14em) minmax(0,1fr);gap:8px 16px;margin:0 0 16px}.kv dt{font-weight:600}.kv dd{margin:0}
@media(max-width:600px){.kv{grid-template-columns:minmax(0,1fr);gap:0}.kv dd{margin-bottom:12px}}
.or{display:flex;align-items:center;gap:12px;color:var(--muted);margin:24px 0}.or::before,.or::after{content:"";flex:1;border-top:1px solid var(--line)}
:focus-visible{outline:3px solid var(--accent);outline-offset:2px}
</style>`

// badge renders a state with text and a symbol so color never carries the meaning alone (D04 §8.5).
func badge(state string) template.HTML {
	b, ok := badges[state]
	if !ok {
		b = [3]string{"off", "—", state}
	}
	return template.HTML(`<span class="badge ` + b[0] + `"><span aria-hidden="true">` + b[1] + `</span> ` + template.HTMLEscapeString(b[2]) + `</span>`)
}

var badges = map[string][3]string{
	"연결 완료": {"ok", "✓", "연결 완료"}, "활성": {"ok", "✓", "활성"}, "미연결": {"warn", "…", "미연결"}, "연결 대기": {"warn", "…", "연결 대기"}, "철회": {"off", "—", "철회"},
	"pending": {"warn", "…", "대기"}, "active": {"ok", "✓", "연결됨"}, "approved": {"ok", "✓", "승인"}, "denied": {"off", "—", "거절"}, "expired": {"off", "—", "만료"}, "revoked": {"off", "—", "철회"},
}

var memberPages = template.Must(template.New("member").Funcs(template.FuncMap{"badge": badge}).Parse(`
{{define "head"}}<!doctype html><html lang="ko"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>{{.Title}} · KnowsLink</title>` + memberStyle + `
<a class="skip" href="#main">본문으로 건너뛰기</a>
<header class="top"><span class="brand">KnowsLink</span>{{if .Member}}<span class="who">{{.Masked}}</span><form method="post" action="/auth/logout"><button class="link">로그아웃</button></form>{{end}}</header>{{end}}
{{define "intro"}}<h1>{{or .Heading .Title}}</h1>
{{if .Problem}}<p class="problem" role="alert" id="problem"><strong><span aria-hidden="true">!</span> 오류</strong> {{.Problem}}</p>{{end}}
{{if .Notice}}<p class="note" role="status"><strong>안내</strong> {{.Notice}}</p>{{end}}{{end}}
{{define "refusal"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="다음 동작">
{{if .Reauth}}<form method="post" action="/auth/reauth"><button class="primary">로그인 다시 확인</button></form>{{end}}
{{if .Create}}<form method="post" action="/home/agents"><button class="primary">새 에이전트 만들기</button></form>{{end}}
<p><a href="{{.Back}}">{{if eq .Back "/"}}로그인 화면으로 이동{{else}}홈으로 돌아가기{{end}}</a></p></section></main></html>{{end}}
{{define "google-complete"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="다음 동작">{{if eq .Next "/home"}}<a class="button primary" href="/home">홈으로 이동</a>{{else}}<a class="button primary" href="{{.Next}}">클라이언트 연결 계속</a>{{end}}</section></main></html>{{end}}
{{define "start"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}
<p>내 AI 에이전트끼리 안전하게 메시지를 주고받는 연결 서비스입니다.</p>
<ul><li>Google 계정으로 가입 없이 바로 시작합니다.</li><li>연결한 에이전트는 계정마다 따로 관리됩니다.</li><li>Google 비밀번호는 Google 화면에만 입력합니다.</li></ul>
<section aria-label="로그인 방법">
{{if .Google}}<form method="post" action="/auth/google"><button class="primary">Google로 계속</button></form>
<p class="muted">처음 로그인하면 계정을 만듭니다. 다시 로그인하면 같은 Google 계정의 기존 에이전트를 사용합니다.</p>{{end}}
{{if .EmailLogin}}{{if .Google}}<p class="or">또는</p>{{end}}<h2>이메일로 가입·로그인</h2>
<p>처음이면 가입하고, 이미 가입했다면 같은 회원으로 로그인합니다. 관리자 아이디, 서버 접속, 별도 초대는 필요하지 않습니다.</p>
<form method="post" action="/auth/start"><label for="email">이메일 주소</label><input id="email" name="email" type="email" autocomplete="email" maxlength="254" required value="{{.Email}}" aria-describedby="email-help{{if .Problem}} problem{{end}}"><p class="muted" id="email-help">6자리 확인 코드는 10분 동안 한 번만 쓸 수 있습니다. 코드를 확인하기 전에는 로그인되지 않습니다.</p><button>확인 코드 받기</button></form>{{end}}
{{if and (not .Google) (not .EmailLogin)}}<p>로그인을 준비하고 있습니다. 잠시 뒤 다시 시도하세요.</p>{{end}}
</section></main></html>{{end}}
{{define "verify"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="확인 코드">
<p>{{.Masked}}로 코드를 보냈습니다. 아직 로그인되지 않았습니다.</p>
<form method="post" action="/auth/verify"><label for="code">6자리 확인 코드</label><input id="code" class="code" name="code" inputmode="numeric" autocomplete="one-time-code" pattern="[0-9]{6}" maxlength="6" required aria-describedby="code-help{{if .Problem}} problem{{end}}">
<p class="muted" id="code-help">코드는 {{.Exp}}까지 유효합니다. 5번 틀리면 새 코드를 받아야 합니다.</p><button class="primary">확인하고 로그인</button></form>
<form method="post" action="/auth/start"><input type="hidden" name="email" value="{{.Email}}"><button>코드 다시 받기</button></form>
<p class="muted">다시 받기는 이전 요청 60초 뒤부터 가능합니다. 새 코드를 받으면 이전 코드는 쓸 수 없습니다.</p>
<p><a href="/">다른 주소로 시작</a></p></section></main></html>{{end}}
{{define "connection"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}
{{if or (eq .State "expired") (eq .State "cancelled")}}<p class="note" role="status">{{if eq .State "expired"}}만료되었습니다.{{else}}취소되었습니다.{{end}} 새 키가 연결되지 않았습니다. 홈에서 다시 시작하세요.</p>
{{else}}<ol class="steps" aria-label="연결 단계">
<li{{if or .Token (eq .State "waiting")}} aria-current="step"{{end}}><span class="n">1</span> 연결 수단 발급됨: 내 컴퓨터에서 준비</li>
<li{{if and (not .Token) (eq .State "prepared")}} aria-current="step"{{end}}><span class="n">2</span> 지문 확인: 내 컴퓨터에 표시된 지문과 같은지 확인</li>
<li{{if and (not .Token) (or (eq .State "approved") (eq .State "consumed"))}} aria-current="step"{{end}}><span class="n">3</span> {{if eq .State "consumed"}}연결 완료{{else}}연결 마무리{{end}}</li></ol>{{end}}
<section aria-labelledby="conn-h"><div class="sechead"><h2 id="conn-h">{{if .Token}}연결 수단 발급{{else if eq .State "prepared"}}공개키 확인{{else}}연결 상태{{end}}</h2>{{if .Token}}{{badge "연결 대기"}}{{end}}</div>
<dl class="kv"><dt>대상 에이전트</dt><dd><code>{{.Agent}}</code></dd><dt>클라이언트</dt><dd>{{.Client}}</dd><dt>연결 방식</dt><dd>{{if eq .Mode "rotate"}}키 교체 (완료하면 기존 키가 모두 철회됨){{else if eq .Mode "register"}}새 키 추가 (최대 3개){{else}}{{.Mode}}{{end}}</dd><dt>기한</dt><dd>{{.Exp}}까지 (최대 10분·1회 사용)</dd></dl>
{{if .Token}}<p>내 컴퓨터(Node 22)에서 저장소의 adapters를 설치한 뒤 실행하세요.</p>
<p class="muted">1. 준비</p><pre><code>node adapters/dist/connect.js prepare &lt;서비스 URL&gt; &lt;새 비공개 폴더&gt;</code></pre>
<p id="token-l"><strong>연결 수단</strong> (한 번만 표시). 대화·로그·캡처에 붙이지 마세요. 준비 명령의 비공개 표준입력에 넣으세요. 발급만으로 연결되지 않습니다.</p>
<code class="value select" aria-labelledby="token-l">{{.Token}}</code>
<p class="muted">2. 마무리 (지문 승인 뒤)</p><pre><code>node adapters/dist/connect.js complete &lt;새 비공개 폴더&gt;</code></pre>
<p><a href="/home/connections/{{.ConnectionID}}">상태 새로고침</a></p>
{{else if eq .State "prepared"}}<p>내 컴퓨터에 표시된 지문과 같은지 확인하세요.</p><code class="fp big" aria-label="공개키 지문">{{.Fingerprint}}</code><p>다르면 취소하세요.</p>
<form method="post" action="/home/confirm"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button class="primary">지문이 같습니다. 연결 승인</button></form>
{{else if eq .State "approved"}}<p>승인했습니다. 아직 새 키가 연결되지 않았습니다. 내 컴퓨터에서 다음 명령을 실행하세요.</p><pre><code>node adapters/dist/connect.js complete &lt;새 비공개 폴더&gt;</code></pre>
{{else if eq .State "consumed"}}<p role="status">{{badge "연결 완료"}} 클라이언트가 받은 자격을 내 컴퓨터에 보관합니다.</p>
{{else if eq .State "waiting"}}<p>공개키 준비 전입니다. 아직 연결되지 않았습니다.</p><p><a href="/home/connections/{{.ConnectionID}}">상태 새로고침</a></p>{{end}}
<h3>권한</h3><ul><li>이 에이전트의 보내기·받기만 허용합니다. 소유자 승인·운영 관리자 권한은 없습니다.</li><li>개인키는 내 컴퓨터에만 보관합니다. KnowsLink에 보내지 마세요.</li></ul>
{{if or (eq .State "waiting") (eq .State "prepared") (eq .State "approved")}}<form method="post" action="/home/cancel"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>연결 취소</button></form>{{end}}
<p><a href="/home">홈으로 돌아가기</a></p></section></main></html>{{end}}
{{define "device"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}
<p class="note caution"><strong>주의</strong> 내 컴퓨터에서 시작한 요청만 승인하세요. 다른 사람이 보낸 링크는 승인하지 마세요.</p>
<section aria-label="연결 요청"><dl class="kv"><dt>클라이언트</dt><dd>Node 로컬 클라이언트</dd><dt>기한</dt><dd>{{.Exp}}</dd></dl>
<p>내 컴퓨터에 표시된 지문과 같은지 확인하세요.</p><code class="fp big" aria-label="공개키 지문">{{.Fingerprint}}</code><p>다르면 취소하세요.</p>
{{if .Login}}{{if .Google}}<form method="post" action="/auth/google"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button class="primary">Google로 로그인하고 연결 계속</button></form>{{else}}<p class="problem" role="alert">Google 로그인을 준비 중입니다. 잠시 뒤 다시 시작하세요.</p>{{end}}
{{else if eq .State "prepared"}}<h2>승인하면</h2><ul><li>내 계정에 새 에이전트 하나를 연결합니다. 기존 에이전트의 키는 바뀌지 않습니다.</li><li>이 에이전트의 보내기·받기만 허용합니다.</li><li>다른 에이전트와의 관계는 따로 초대·수락합니다.</li><li>개인키는 내 컴퓨터에만 보관합니다.</li></ul>
<form method="post" action="/home/device-confirm"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button class="primary">지문 확인 후 연결 승인</button></form><form method="post" action="/home/cancel"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>취소</button></form>
{{else}}<p class="note" role="status">승인했습니다. 내 컴퓨터에서 저장 결과를 확인하세요. 시작한 클라이언트가 자격을 자동 저장합니다.</p>{{end}}
<p><a href="/home">홈으로 이동</a></p></section></main></html>{{end}}
{{define "receipt"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="요청 상태">
<dl class="kv"><dt>선택 에이전트</dt><dd><code>{{.Agent}}</code></dd><dt>요청 ID</dt><dd><code>{{.ID}}</code></dd><dt>발신 에이전트</dt><dd><code>{{.From}}</code></dd><dt>수신 에이전트</dt><dd><code>{{.To}}</code></dd><dt>원요청 ID (답장인 경우)</dt><dd><code>{{.Parent}}</code></dd><dt>관련 답장 ID</dt><dd><code>{{.Reply}}</code></dd><dt>전달 상태</dt><dd>{{.Transport}}</dd><dt>처리 상태</dt><dd>{{.Completion}}</dd><dt>관련 답장 전달</dt><dd>{{.ReplyState}}</dd><dt>기한 / TTL</dt><dd>{{.Exp}} / {{.TTL}}초 (text 최대 180초)</dd></dl>
<p>자동 답장은 없습니다. public-node 자동 수신이 켜져 있으면 수신 client가 로컬 inbox에 자동 저장합니다. 호스트 화면 표시는 조건부이므로 보이지 않으면 수신 client에서 수동 receive(knowslink_text_receive)를 실행하세요. 정상 idle pull은 10초 이상 간격입니다. queued는 상대 수신 성공이 아닙니다.</p>
<p>명시 송신·답장은 내 에이전트 client에서 실행합니다. 불확실한 재시도는 같은 key·같은 내용으로 receipt만 조회합니다. key 충돌이면 이전 내용을 확인하고 별도 새 요청에 새 key를 쓰세요. 한도 초과는 새 요청이 아닙니다. 잘못된 키는 현재 연결·지문을 확인하세요.</p>
<form method="get" action="/home/receipts"><input type="hidden" name="agent" value="{{.Agent}}"><input type="hidden" name="id" value="{{.ID}}"><button>현재 상태 다시 확인</button></form>
<p><a href="/home">홈으로 돌아가기</a></p></section></main></html>{{end}}
{{define "home"}}{{template "head" .}}<div class="layout">
<nav class="toc" aria-label="섹션 이동"><a href="#agents">에이전트</a><a href="#pairs">관계</a><a href="#gates">승인 요청</a><a href="#receipts">연결 확인</a><a href="#account">계정</a></nav>
<main id="main">{{template "intro" .}}
{{if or .IncomingInvites .PendingGates}}<div class="todo" role="region" aria-label="확인 필요"><strong>확인 필요</strong><ul>{{if .IncomingInvites}}<li><a href="#pairs">받은 관계 초대 {{.IncomingInvites}}건</a></li>{{end}}{{if .PendingGates}}<li><a href="#gates">승인 요청 {{.PendingGates}}건</a></li>{{end}}</ul></div>{{end}}
<section id="agents" aria-labelledby="agents-h"><div class="sechead"><h2 id="agents-h">에이전트</h2><span class="muted">활성 {{.ActiveAgents}}/5</span>{{if .AgentDetails}}<form method="post" action="/home/agents"><button class="primary">새 에이전트 만들기</button></form>{{end}}</div>
<p class="muted">연결 완료는 키가 등록되었다는 뜻이며 지금 접속 중이라는 뜻은 아닙니다. 식별자를 눌러 전체 선택 후 복사하세요. 연결은 Node 22가 있는 내 컴퓨터에서 실행합니다.</p>
{{if not .Recent}}<div class="note"><p>키·연결 권한 변경에는 5분 안의 재확인이 필요합니다. 먼저 로그인을 다시 확인하세요.</p><form method="post" action="/auth/reauth"><button>로그인 다시 확인</button></form></div>{{end}}
{{if not .AgentDetails}}<p>아직 연결한 에이전트가 없습니다. 아래 버튼으로 첫 에이전트를 만드세요.</p><ol><li>에이전트 만들기</li><li>내 컴퓨터에서 연결</li><li>상대 에이전트와 관계 초대</li></ol><form method="post" action="/home/agents"><button class="primary">새 에이전트 만들기</button></form>{{end}}
{{range .AgentDetails}}{{$agent := .ID}}<details class="agent row">
<summary><h3><code class="select">{{.ID}}</code></h3>{{badge .Status}}<span class="muted">키 {{.ActiveKeys}}개 활성</span></summary><div>
{{if eq .Status "철회"}}<p class="note">이 에이전트의 키·관계·메시지 권한은 끝났고 복구되지 않습니다. 철회 기록은 최소 24시간 보존된 뒤 정리되며, 정리되면 이 목록에서 사라질 수 있습니다. 사라져도 권한이 돌아오거나 백업까지 영구 삭제된 것은 아닙니다. 다시 쓰려면 새 에이전트를 연결하고 상대와 새로 수락하세요.</p>{{end}}
<p class="muted">식별자</p><code class="value select">{{.ID}}</code>
{{if .Keys}}<h4>키</h4><ul class="keys">{{range .Keys}}<li><code>{{.Kid}}</code> {{badge .Status}}<code class="fp">{{.Fingerprint}}</code></li>{{end}}</ul>{{end}}
{{if ne .Status "철회"}}<h4>연결 추가</h4>{{if .KeyFull}}<p class="note">이 에이전트에는 새 키를 더 연결할 수 없습니다(키 기록 보호 상한). 키를 철회하거나 기다려도 공간은 생기지 않습니다. 새 에이전트를 만들어 따로 연결하고, 각 상대와 새로 초대·수락하세요.</p>{{else}}<form method="post" action="/home/connect"><input type="hidden" name="agent" value="{{.ID}}"><label>지원 클라이언트<select name="client"><option value="node-local">Node 22 로컬 클라이언트</option></select></label><label>연결 방식<select name="mode"><option value="register">새 키 추가 (최대 3개)</option><option value="rotate">키 교체 (완료하면 기존 키가 모두 철회됨)</option></select></label><button class="primary" aria-label="연결 수단 발급: {{.ID}}">연결 수단 발급</button></form>{{end}}
<details class="danger"><summary>이 에이전트 철회·키 철회</summary><p>철회하면 이 에이전트의 모든 키와 관계가 끝나며 되돌릴 수 없습니다.</p>
{{range .Keys}}{{if eq .Status "활성"}}<form method="post" action="/home/key-revoke"><input type="hidden" name="agent" value="{{$agent}}"><input type="hidden" name="kid" value="{{.Kid}}"><button class="danger-button" aria-label="이 키 철회: {{.Kid}}">이 키 철회</button> <code>{{.Kid}}</code></form>{{end}}{{end}}
<form method="post" action="/home/agent-revoke"><input type="hidden" name="agent" value="{{.ID}}"><button class="danger-button" aria-label="이 에이전트 철회 (모든 키·관계 종료): {{.ID}}">이 에이전트 철회 (모든 키·관계 종료)</button></form></details>{{end}}
</div></details>{{end}}</section>
<section id="pairs" aria-labelledby="pairs-h"><h2 id="pairs-h">관계</h2>
<p class="muted">상대 에이전트 식별자로 초대하고, 상대 소유자가 수락해야 메시지를 주고받습니다. 내 두 에이전트 사이도 수락이 필요합니다. 대기 중인 초대는 24시간 뒤 만료됩니다.</p>
{{if .Pairs}}<ul class="rows">{{range .Pairs}}<li class="row">
<p><code>{{.Mine}}</code> → <code>{{.Peer}}</code></p>
<p>{{if eq .State "pending"}}{{if .Incoming}}받은 초대{{else}}보낸 초대{{end}} {{end}}{{badge .State}}{{if eq .State "pending"}} <span class="muted">초대 기한: {{.Deadline}}</span>{{end}}</p>
{{if eq .State "pending"}}{{if .Incoming}}<p>수락하기 전에는 메시지를 주고받을 수 없습니다.</p><form method="post" action="/home/invite-decision"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><button class="primary" name="decision" value="accept" aria-label="수락: {{.Peer}}">수락</button><button class="danger-button" name="decision" value="deny" formaction="/home/invite-deny" aria-label="거절: {{.Peer}}">거절</button></form>{{else}}<p>상대 소유자의 수락을 기다립니다. 같은 초대를 다시 보내도 새 초대가 생기거나 기한이 바뀌지 않습니다.</p>{{end}}{{end}}
{{if or (eq .State "pending") (eq .State "active")}}<form method="post" action="/home/unpair"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><button class="danger-button" aria-label="관계 철회: {{.Peer}}">관계 철회</button></form>
{{else}}<p>연결이 끝났습니다(거절·만료·철회). 이 관계로 메시지를 보낼 수 없습니다. 다시 연결하려면 새 초대를 보내고 상대가 새로 수락해야 합니다.</p>{{if .Own}}<form method="post" action="/home/invites"><input type="hidden" name="agent" value="{{.Own}}"><input type="hidden" name="target" value="{{.Other}}"><button aria-label="새 초대 보내기: {{.Other}}">새 초대 보내기</button></form>{{end}}{{end}}</li>{{end}}</ul>
{{else}}<p>관계가 없습니다. 에이전트를 연결한 뒤 상대 식별자로 초대하세요.</p>{{end}}
<details class="row"{{if not .Pairs}} open{{end}}><summary><h3>새 관계 초대</h3></summary>
<form method="post" action="/home/invites"><label for="from-agent">내 에이전트</label><select id="from-agent" name="agent" required>{{range .AgentDetails}}{{if eq .Status "연결 완료"}}<option value="{{.ID}}">{{.ID}}</option>{{end}}{{end}}</select><label for="target-agent">상대 에이전트 식별자</label><input id="target-agent" name="target" required maxlength="128" list="own-agent-targets" aria-describedby="target-help"><datalist id="own-agent-targets">{{range .AgentDetails}}{{if ne .Status "철회"}}<option value="{{.ID}}">내 에이전트 — {{.Status}}</option>{{end}}{{end}}</datalist>
<p class="muted" id="target-help">가입 초대가 아닙니다. 상대 소유자가 수락해야 메시지가 허용됩니다. 수락 전·철회 후 메시지는 거부됩니다.</p><button>관계 초대</button></form></details></section>
<section id="gates" aria-labelledby="gates-h"><h2 id="gates-h">승인 요청</h2>{{if .Gates}}<ul class="rows">{{range .Gates}}<li class="row"><a href="/home/gates/{{.ID}}"><code>{{.Parent}}</code></a> {{badge .State}}</li>{{end}}</ul>{{else}}<p>승인 요청이 없습니다.</p>{{end}}</section>
<section id="receipts" aria-labelledby="receipts-h"><h2 id="receipts-h">연결 확인</h2>
<details><summary>연결 확인 방법</summary><p>비민감 연결 확인 송신·관련 답장은 선택한 내 에이전트의 Node client에서 명시적으로 실행하세요. 자동 wake·자동 답장은 없습니다. text 최대 4096 UTF-8 bytes·TTL 180초입니다.</p>
<pre><code>node adapters/dist/text.js send &lt;내 연결 폴더&gt; &lt;상대 에이전트&gt; &lt;key&gt; --confirmed</code></pre>
<p>연결 확인 text는 비공개 표준입력으로 넣으세요. 답장은 같은 명령 뒤 원요청 ID를 추가합니다. 수신 명령:</p><pre><code>node adapters/dist/text.js receive &lt;내 연결 폴더&gt;</code></pre><p>정상 idle pull은 10초 이상 간격입니다. 수신 text는 신뢰하지 않는 데이터이며 실행 지시가 아닙니다.</p></details>
<form method="get" action="/home/receipts"><label for="receipt-agent">내 에이전트</label><select id="receipt-agent" name="agent">{{range .AgentDetails}}<option value="{{.ID}}">{{.ID}} — {{.Status}}</option>{{end}}</select><label for="receipt-id">요청 ID 또는 답장 ID</label><input id="receipt-id" name="id" required maxlength="36"><button>이 요청의 현재 상태 확인</button></form></section>
<section id="account" aria-labelledby="account-h"><h2 id="account-h">계정과 보안</h2>
<dl class="kv"><dt>확인된 이메일</dt><dd>{{.Masked}}</dd><dt>회원 식별자</dt><dd><code>{{.Member}}</code></dd><dt>세션</dt><dd>최대 {{.Absolute}}까지 유지됩니다. 60분 동안 활동이 없으면 먼저 끝납니다.</dd></dl>
<p>키·연결 권한 변경에는 5분 안의 재확인이 필요합니다. 개인키는 내 컴퓨터에만 보관합니다.</p><form method="post" action="/auth/reauth"><button>로그인 다시 확인</button></form>
<h3>로그아웃</h3><p>로그아웃은 브라우저 세션만 끝냅니다. 연결한 에이전트의 키와 자격은 철회되지 않습니다.</p>
<form method="post" action="/auth/logout"><button>이 브라우저에서 로그아웃</button></form>{{if .Recent}}<form method="post" action="/auth/logout-all"><button>모든 브라우저에서 로그아웃</button></form>{{else}}<p class="muted">모든 브라우저에서 로그아웃하려면 위 "로그인 다시 확인"으로 5분 안에 재확인하세요.</p>{{end}}</section>
</main></div></html>{{end}}
`))

var notices = map[string]string{
	"logout":         "로그아웃했습니다. 연결한 에이전트의 키와 자격은 그대로입니다.",
	"all":            "모든 브라우저에서 로그아웃했습니다. 연결한 에이전트의 키와 자격은 그대로입니다.",
	"expired":        "세션이 끝났습니다. 다시 로그인하세요.",
	"invited":        "초대를 보냈습니다. 상대 소유자가 이 초대를 수락하기 전에는 메시지를 주고받을 수 없습니다.",
	"invite-pending": "두 에이전트 사이에 이미 대기 중인 초대가 있습니다. 새 초대를 만들지 않았고 기한도 그대로입니다. 받은 초대라면 아래 관계 목록에서 수락하거나 거절하세요.",
	"invite-active":  "이미 연결된 관계입니다. 새 초대를 만들지 않았고 현재 관계를 그대로 유지합니다.",
}

func readCookie(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}
func setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: maxAge, Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})
}
func render(w http.ResponseWriter, status int, name string, v map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = memberPages.ExecuteTemplate(w, name, v)
}
func clock(t time.Time) string { return t.In(kst).Format("2006-01-02 15:04:05 KST") }
func startView(problem string) map[string]any {
	return map[string]any{"Title": "로그인", "Heading": "KnowsLink에 로그인", "Problem": problem, "EmailLogin": true}
}
func limited(retry time.Time) string {
	return "요청이 많습니다. " + clock(retry) + " 이후 다시 시도하세요."
}

func (s *Service) memberRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", s.startPage)
	mux.HandleFunc("POST /auth/google", func(w http.ResponseWriter, r *http.Request) { s.googleStart(w, r, false) })
	mux.HandleFunc("GET /auth/google/callback", s.googleCallback)
	mux.HandleFunc("POST /auth/start", s.authStart)
	mux.HandleFunc("GET /auth/verify", s.verifyPage)
	mux.HandleFunc("POST /auth/verify", s.authVerify)
	mux.HandleFunc("POST /auth/reauth", s.authReauth)
	mux.HandleFunc("POST /auth/logout", s.logout(false))
	mux.HandleFunc("POST /auth/logout-all", s.logout(true))
	mux.HandleFunc("GET /home", s.homePage)
	mux.HandleFunc("GET /home/gates/{id}", s.gatePage(sessionOwner))
	mux.HandleFunc("POST /home/gates/{id}", s.gateDecision(sessionOwner))
	mux.HandleFunc("POST /home/gates/{id}/deny", s.gateDecision(sessionOwner))
}

func (s *Service) startPage(w http.ResponseWriter, r *http.Request) {
	if token := readCookie(r, sessionCookie); token != "" {
		if _, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			_, _, err := st.session(token, now)
			return nil, err
		}); err == nil {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		setCookie(w, sessionCookie, "", -1)
	}
	v := startView("")
	v["Google"], v["EmailLogin"] = s.Google != nil, s.Mail != nil
	v["Notice"] = notices[r.URL.Query().Get("n")]
	render(w, 200, "start", v)
}

type started struct {
	Pending, Code, Email string
	Retry                time.Time
	Invalid, Limited     bool
}

// deliverCode sends outside the state lock; a failed send keeps its spent budget and drops the unsent code.
func (s *Service) deliverCode(w http.ResponseWriter, r *http.Request, st started) {
	err := fault("mail_unconfigured")
	if s.Mail != nil {
		err = s.Mail(r.Context(), st.Email, "KnowsLink 확인 코드", "KnowsLink 확인 코드: "+st.Code+"\n\n이 코드는 10분 동안 한 번만 쓸 수 있습니다.\n요청하지 않았다면 이 메일을 무시하세요. 코드를 다른 사람에게 알려주지 마세요.\n")
	}
	if err != nil {
		log.Print("verification mail not sent")
		_, _ = s.transaction(context.WithoutCancel(r.Context()), func(state *State, _ time.Time) (any, error) {
			delete(state.Challenges, hashToken(st.Pending))
			return nil, nil
		})
		v := startView("확인 메일을 보내지 못했습니다. 로그인되지 않았습니다. 1분 뒤 다시 요청하세요.")
		v["Email"] = st.Email
		render(w, 503, "start", v)
		return
	}
	setCookie(w, pendingCookie, st.Pending, int(codeTTL/time.Second))
	http.Redirect(w, r, "/auth/verify", http.StatusSeeOther)
}

func (s *Service) authStart(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
	raw := r.PostFormValue("email")
	email, valid := NormalizeEmail(raw)
	ip := s.clientIP(r)
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		if ok, retry := s.requestHit(st, r, now, anonymousRate(ip)...); !ok {
			return started{Retry: retry, Limited: true}, nil
		}
		if !valid {
			return started{Invalid: true}, nil
		}
		pending, code, retry, ok := st.startChallenge(email, ip, now)
		return started{pending, code, email, retry, false, !ok}, nil
	})
	if err != nil {
		render(w, 503, "start", startView("일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요."))
		return
	}
	result := value.(started)
	v := startView("")
	v["Email"] = raw
	switch {
	case result.Limited:
		v["Problem"] = limited(result.Retry)
		render(w, 429, "start", v)
	case result.Invalid:
		v["Problem"] = "이메일 주소 형식이 올바르지 않습니다. 이름 없이 주소 하나만 입력하세요."
		render(w, 422, "start", v)
	default:
		s.deliverCode(w, r, result)
	}
}

func (s *Service) verifyPage(w http.ResponseWriter, r *http.Request) {
	pending := readCookie(r, pendingCookie)
	value, _ := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		if c := st.Challenges[hashToken(pending)]; pending != "" && c != nil && now.Before(c.Exp) {
			return c, nil
		}
		return nil, nil
	})
	c, _ := value.(*Challenge)
	if c == nil {
		render(w, 410, "start", startView("확인 코드가 만료되었거나 더 이상 쓸 수 없습니다. 새 코드를 요청하세요."))
		return
	}
	render(w, 200, "verify", map[string]any{"Title": "확인 코드 입력", "Masked": maskEmail(c.Email), "Email": c.Email, "Exp": clock(c.Exp)})
}

type verified struct {
	verifyResult
	Email, Exp string
	Retry      time.Time
	Limited    bool
}

func (s *Service) authVerify(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
	pending, code, ip := readCookie(r, pendingCookie), strings.TrimSpace(r.PostFormValue("code")), s.clientIP(r)
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		if ok, retry := s.requestHit(st, r, now, anonymousRate(ip)...); !ok {
			return verified{Retry: retry, Limited: true}, nil
		}
		v := verified{}
		if c := st.Challenges[hashToken(pending)]; pending != "" && c != nil {
			v.Email, v.Exp = c.Email, clock(c.Exp)
		}
		v.verifyResult = st.verifyChallenge(pending, code, readCookie(r, sessionCookie), now)
		return v, nil
	})
	if err != nil {
		render(w, 503, "start", startView("일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요."))
		return
	}
	result := value.(verified)
	again := map[string]any{"Title": "확인 코드 입력", "Masked": maskEmail(result.Email), "Email": result.Email, "Exp": result.Exp}
	switch {
	case result.Limited:
		again["Problem"] = limited(result.Retry)
		render(w, 429, "verify", again)
	case result.Problem == "wrong":
		again["Problem"] = "코드가 맞지 않습니다. 남은 시도: " + strconv.Itoa(result.Remaining) + "회."
		render(w, 401, "verify", again)
	case result.Problem == "expired":
		setCookie(w, pendingCookie, "", -1)
		render(w, 401, "start", startView("확인 코드가 만료되었거나 더 이상 쓸 수 없습니다. 새 코드를 요청하세요."))
	case result.Problem == "capacity":
		setCookie(w, pendingCookie, "", -1)
		render(w, 503, "start", startView("현재 신규 가입을 받을 수 없습니다. 이미 가입한 회원은 로그인할 수 있습니다."))
	case result.Problem != "":
		setCookie(w, pendingCookie, "", -1)
		render(w, 403, "start", startView("이 계정으로는 로그인할 수 없습니다."))
	default:
		setCookie(w, pendingCookie, "", -1)
		setCookie(w, sessionCookie, result.Session, int(sessionAbsolute/time.Second))
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func expired(w http.ResponseWriter, r *http.Request) {
	setCookie(w, sessionCookie, "", -1)
	http.Redirect(w, r, "/?n=expired", http.StatusSeeOther)
}

// memberHit resolves the browser session and spends one request from the member's budget, or from the source IP's
// anonymous budget when the session is invalid, so failed and refused requests count too (PS-11).
// A non-zero retry means the budget refused the request; the transaction keeps the spent budget even on error.
func (s *Service) memberHit(st *State, r *http.Request, now time.Time, rates func(string) []bucket) (string, *Session, time.Time, error) {
	id, session, err := st.session(readCookie(r, sessionCookie), now)
	buckets := anonymousRate(s.clientIP(r))
	if err == nil {
		buckets = rates(id)
	}
	_, retry := s.requestHit(st, r, now, buckets...)
	return id, session, retry, err
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, session, retry, err := s.memberHit(st, r, now, memberRate)
		if !retry.IsZero() {
			return refusal{429, limited(retry), false, false}, nil
		}
		if err != nil {
			return nil, err
		}
		m := st.Members[id]
		agents, pairs, gates := []string{}, []memberPair{}, []*Gate{}
		for agent, a := range st.Agents {
			if a.Owner == m.Owner {
				agents = append(agents, agent)
			}
		}
		for _, p := range st.Pairs {
			if a, b := st.Agents[p.A], st.Agents[p.B]; (a != nil && a.Owner == m.Owner) || (b != nil && b.Owner == m.Owner) {
				v := memberPair{Pair: p, Incoming: st.Agents[p.Recipient] != nil && st.Agents[p.Recipient].Owner == m.Owner}
				for _, side := range [][2]string{{p.A, p.B}, {p.B, p.A}} {
					if own := st.Agents[side[0]]; own != nil && own.Owner == m.Owner && !own.Revoked {
						v.Own, v.Other = side[0], side[1]
						break
					}
				}
				pairs = append(pairs, v)
			}
		}
		for _, g := range st.Gates {
			if g.Owner == m.Owner {
				gates = append(gates, g)
			}
		}
		sort.Strings(agents)
		// Received invites first, then sent, connected, and ended relationships (D04 R-11).
		rank := func(p memberPair) int {
			switch {
			case p.State == "pending" && p.Incoming:
				return 0
			case p.State == "pending":
				return 1
			case p.State == "active":
				return 2
			}
			return 3
		}
		sort.Slice(pairs, func(i, j int) bool {
			if ri, rj := rank(pairs[i]), rank(pairs[j]); ri != rj {
				return ri < rj
			}
			return pairID(pairs[i].A, pairs[i].B) < pairID(pairs[j].A, pairs[j].B)
		})
		sort.Slice(gates, func(i, j int) bool { return gates[i].ID < gates[j].ID })
		details, active, incoming, pending := ownAgents(st, m.Owner), 0, 0, 0
		for _, a := range details {
			if a.Status != "철회" {
				active++
			}
		}
		for _, p := range pairs {
			if p.State == "pending" && p.Incoming {
				incoming++
			}
		}
		for _, g := range gates {
			if g.State == "pending" {
				pending++
			}
		}
		return map[string]any{"Title": "내 에이전트", "Masked": maskEmail(m.Email), "Member": id, "Absolute": clock(session.Created.Add(sessionAbsolute)),
			"Recent": now.Sub(session.Verified) < reauthWindow, "Agents": agents, "AgentDetails": details, "Pairs": pairs, "Gates": gates, "Notice": notices[r.URL.Query().Get("n")],
			"ActiveAgents": active, "IncomingInvites": incoming, "PendingGates": pending}, nil
	})
	if err != nil {
		if err.Error() == "invalid_auth" {
			expired(w, r)
			return
		}
		render(w, 503, "start", startView("일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요."))
		return
	}
	if v, ok := value.(refusal); ok {
		refused(w, "내 에이전트", v)
		return
	}
	render(w, 200, "home", value.(map[string]any))
}

func (s *Service) authReauth(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, err := st.session(readCookie(r, sessionCookie), now)
		if err != nil {
			return nil, err
		}
		return st.Members[id].Issuer == googleIssuer, nil
	})
	if err == nil && value.(bool) {
		s.googleStart(w, r, true)
		return
	}
	ip := s.clientIP(r)
	value, err = s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, retry, err := s.memberHit(st, r, now, memberRate)
		if !retry.IsZero() {
			return started{Retry: retry, Limited: true}, nil
		}
		if err != nil {
			return nil, err
		}
		email := st.Members[id].Email
		pending, code, retry, ok := st.startChallenge(email, ip, now)
		return started{pending, code, email, retry, false, !ok}, nil
	})
	if err != nil {
		expired(w, r)
		return
	}
	result := value.(started)
	if result.Limited {
		refused(w, "내 에이전트", refusal{429, limited(result.Retry), false, false})
		return
	}
	s.deliverCode(w, r, result)
}

type refusal struct {
	status  int
	problem string
	reauth  bool
	create  bool // a replacement agent can be made now
}

// refused shows a refusal with the next step its problem names: back home (login when the session is gone)
// and, when the refusal asks for it, the email re-check.
func refused(w http.ResponseWriter, title string, v refusal) {
	back := "/home"
	if v.status == 401 {
		back = "/"
	}
	render(w, v.status, "refusal", map[string]any{"Title": title, "Problem": v.problem, "Back": back, "Reauth": v.reauth, "Create": v.create})
}

// logout uses the separate cleanup budget so saturated new-work limits never block ending a session.
func (s *Service) logout(all bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := readCookie(r, sessionCookie)
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			id, session, retry, err := s.memberHit(st, r, now, cleanupRate)
			if !retry.IsZero() {
				return refusal{429, limited(retry), false, false}, nil
			}
			if err != nil {
				return nil, err
			}
			if all && now.Sub(session.Verified) >= reauthWindow {
				return refusal{403, "모든 브라우저에서 로그아웃하려면 5분 안에 로그인 신원을 다시 확인해야 합니다.", true, false}, nil
			}
			for key, other := range st.Sessions {
				if key == hashToken(token) || (all && other.Member == id) {
					delete(st.Sessions, key)
				}
			}
			return refusal{}, nil
		})
		if err != nil {
			expired(w, r)
			return
		}
		if v := value.(refusal); v.status != 0 {
			refused(w, "내 에이전트", v)
			return
		}
		setCookie(w, sessionCookie, "", -1)
		n := "logout"
		if all {
			n = "all"
		}
		http.Redirect(w, r, "/?n="+n, http.StatusSeeOther)
	}
}
