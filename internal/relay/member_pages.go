// Member page views: the shared design tokens, status badges, and every member HTML template.
package relay

import "html/template"

// memberStyle is the single source of the member web tokens (D04 service-member-ui.md §8). CSP allows only inline
// styles, so there are no fonts, images, or scripts; color, type, spacing, buttons, and badges all live here.
const memberStyle = `<style>
:root{--canvas:#FFFFFF;--ink:#182B31;--muted:#52646B;--accent:#005A62;--selected:#E8F3F3;--edge:#75868B;--danger:#A01616;--bg:var(--canvas);--surface:var(--canvas);--text:var(--ink);--line:var(--edge);--line-strong:var(--edge);--accent-soft:var(--selected);--danger-soft:var(--canvas);--mono:ui-monospace,SFMono-Regular,Consolas,monospace;color-scheme:light;accent-color:var(--accent)}
*,::before,::after{box-sizing:border-box}::selection{background:var(--selected);color:var(--ink)}
body{max-width:1440px;margin:0 auto;padding:24px 32px 32px;background:var(--canvas);color:var(--ink);font:16px/1.6 "Malgun Gothic","Apple SD Gothic Neo","Segoe UI",sans-serif;word-break:keep-all;overflow-wrap:break-word;font-variant-numeric:tabular-nums}
h1{font-size:28px;line-height:1.3;font-weight:700;margin:0 0 24px}h2{font-size:20px;line-height:1.3;margin:0 0 16px}h3,h4{font-size:18px;line-height:1.3;font-weight:600;margin:0 0 8px}h4{font-size:16px}p{max-width:64ch;margin:0 0 16px}ul,ol{margin:0 0 16px;padding-left:24px}
a{display:inline-flex;align-items:center;min-height:44px;color:var(--accent);text-underline-offset:4px}
code,pre{font:14px/1.6 var(--mono)}code{word-break:normal;overflow-wrap:anywhere}pre{white-space:pre-wrap;overflow-wrap:anywhere;border-left:1px solid var(--edge);padding:12px 16px;margin:8px 0 24px}
.skip{position:absolute;left:16px;top:-80px;background:var(--canvas);padding:8px 12px}.skip:focus{top:8px;z-index:1}
.top{display:flex;flex-wrap:wrap;align-items:center;gap:8px 16px;padding-bottom:16px;margin-bottom:24px;border-bottom:1px solid var(--edge)}.brand{font-weight:700;font-size:20px;margin-right:auto;color:var(--ink);text-decoration:none}.who,.muted{color:var(--muted);font-size:14px}.page{max-width:800px;margin:0 auto}.page.narrow{max-width:560px}.top form,.inline{display:inline}.nowrap{white-space:nowrap}
section{padding:0;margin:0 0 32px}section>:last-child{margin-bottom:0}.sechead{display:flex;flex-wrap:wrap;align-items:center;gap:8px 16px;margin-bottom:16px}.sechead h1,.sechead h2{margin:0 auto 0 0}
label{display:block;font-weight:600;margin:16px 0 4px}input,select,textarea{display:block;width:100%;min-width:0;min-height:44px;margin:4px 0 12px;padding:8px 12px;border:1px solid var(--edge);border-radius:4px;background:var(--canvas);color:inherit;font:inherit}textarea{min-height:128px;resize:vertical;line-height:1.6}input.code{font:600 28px/1.3 var(--mono);letter-spacing:.3em;max-width:9em}
button,.button{display:inline-flex;align-items:center;justify-content:center;min-height:44px;max-width:100%;margin:4px 8px 4px 0;padding:8px 16px;border:1px solid var(--edge);border-radius:4px;background:var(--canvas);color:var(--ink);font:inherit;text-align:center;text-decoration:none;cursor:pointer}.primary{background:var(--accent);border-color:var(--accent);color:var(--canvas);font-weight:600}.danger-button{border-color:var(--danger);color:var(--danger)}.link{border:0;background:none;color:var(--accent);text-decoration:underline;padding:8px 4px;margin:0}button:disabled{cursor:default;color:var(--muted);background:var(--canvas)}
.badge{display:inline;font-size:14px;font-weight:400;color:var(--muted);white-space:normal}.ok,.warn,.off{color:var(--muted)}.note,.problem,.todo{margin:0 0 24px;padding:12px 16px;border-left:3px solid var(--accent);background:var(--selected)}.problem,.caution{border-left-color:var(--danger);background:var(--canvas)}.problem strong{color:var(--danger)}.todo ul{margin:4px 0 0}.note p:last-child,.problem p:last-child{margin-bottom:0}
.rows,.keys{list-style:none;padding:0;margin:0}.row{padding:16px 0;border-top:1px solid var(--edge)}summary{display:flex;flex-wrap:wrap;align-items:center;gap:8px 12px;min-height:44px;cursor:pointer;list-style:none}summary::-webkit-details-marker{display:none}summary::before{content:"▸";color:var(--muted);width:12px}details[open]>summary::before{content:"▾"}summary h3{flex:1 1 12em;min-width:0;margin:0}.agent>div{padding:16px 0 16px 16px}.keys li{padding:8px 0}.fp{display:block;letter-spacing:.02em;user-select:all}.fp.big{font-size:16px;font-weight:600;padding:12px 16px;border:1px solid var(--edge);border-radius:4px;margin:8px 0 16px}.select{user-select:all}.value{display:block;padding:8px 12px;margin:4px 0 16px;border:1px solid var(--edge);border-radius:4px}.danger{margin-top:24px}.danger>summary{color:var(--danger);font-weight:600}.danger :focus-visible{outline-color:var(--danger)}.danger>p:first-of-type{margin-top:12px}
.steps{list-style:none;padding:0;display:grid;gap:8px;margin:0 0 24px}.steps li{padding:8px 12px;border-left:1px solid var(--edge)}.steps li[aria-current]{border-left:3px solid var(--accent);font-weight:600}.steps .n{display:inline-block;min-width:24px;font-weight:700}.kv{display:grid;grid-template-columns:minmax(8em,14em) minmax(0,1fr);gap:8px 16px;margin:0 0 24px}.kv dt{font-weight:600}.kv dd{margin:0;min-width:0}.or{display:flex;align-items:center;gap:12px;color:var(--muted);margin:24px 0}.or::before,.or::after{content:"";flex:1;border-top:1px solid var(--edge)}
.shell{display:grid;grid-template-columns:184px minmax(0,1fr);gap:24px}.workspace{min-width:0}.nav,.toc{display:flex;flex-direction:column;align-self:start;gap:4px}.nav a,.toc a{padding:8px 12px;text-decoration:none;color:var(--ink)}.nav a[aria-current],.conversation-list a[aria-current]{border-left:3px solid var(--accent);background:var(--selected);font-weight:600}.chat-layout{display:grid;grid-template-columns:280px minmax(0,1fr);gap:24px}.conversation-list{border-right:1px solid var(--edge);padding-right:24px;min-width:0}.conversation-list ul,.people,.agents,.messages{list-style:none;padding:0;margin:0}.conversation-list a{display:block;padding:12px 8px;text-decoration:none;border-left:3px solid transparent;color:var(--ink)}.person-name{font-size:18px;font-weight:600}.name-preview{display:-webkit-box;-webkit-box-orient:vertical;-webkit-line-clamp:2;overflow:hidden}.agent-line{display:block;margin-left:16px;font-size:16px}.short{color:var(--muted);font-size:14px;overflow-wrap:anywhere}.people>li,.agents>li{padding:16px 0}.people .agents,.friend-agents{margin-left:16px}.messages>li{padding:16px 0 16px 16px;border-left:3px solid transparent}.messages>li.owner-message{border-left-color:var(--accent)}.message-meta{display:flex;flex-wrap:wrap;gap:4px 12px;margin-bottom:8px}.message-meta>*{min-width:0;max-width:100%;overflow-wrap:anywhere}.conversation-body{min-width:0}.message-meta time{color:var(--muted);font-size:14px}.message-text{white-space:pre-wrap;max-width:64ch;overflow-wrap:anywhere}.authorship{font-size:14px;background:var(--selected);padding:0 4px}.conversation-header{margin-bottom:24px}.conversation-header h1{margin-bottom:8px}.composer{margin-top:24px;max-width:64ch}.target{margin:12px 0;overflow-wrap:anywhere}.consent{display:flex;align-items:flex-start;gap:8px;font-weight:400}.consent input{flex:0 0 24px;width:24px;min-height:44px;margin:0}.empty{padding:24px 0}.back-mobile{display:none}.retention{margin-bottom:24px}.retention summary{color:var(--accent)}.layout{display:grid;grid-template-columns:184px minmax(0,1fr);gap:24px}
:focus-visible{outline:3px solid var(--accent);outline-offset:2px}@media(prefers-reduced-motion:reduce){*{scroll-behavior:auto}}
@media(max-width:1199px){.shell,.layout{display:block}.nav,.toc{flex-direction:row;flex-wrap:wrap;margin-bottom:24px;gap:4px 8px}.chat-layout{grid-template-columns:280px minmax(0,1fr)}}
@media(max-width:767px){body{padding:16px}.top{margin-bottom:16px}.nav{gap:4px}.nav a{padding:8px}.chat-layout{display:block}.conversation-list{border:0;padding:0}.chat-layout.has-conversation>.conversation-list{display:none}.chat-layout:not(.has-conversation)>.conversation-body{display:none}.back-mobile{display:inline-flex}.primary{width:100%;margin-right:0}.kv{grid-template-columns:minmax(0,1fr);gap:0}.kv dd{margin-bottom:12px}.agent>div{padding-left:0}.conversation-header h1{font-size:28px}.composer{max-width:100%}}
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
	"queued": {"warn", "…", "대기"}, "leased": {"warn", "…", "대기"}, "delivered": {"ok", "✓", "전달됨"}, "completed": {"ok", "✓", "처리 완료"}, "failed": {"off", "!", "실패"},
	"cancelled": {"off", "—", "취소"}, "removed": {"off", "—", "해제"}, "blocked": {"off", "—", "차단"}, "registered": {"ok", "✓", "등록됨"}, "unconnected": {"off", "—", "연결 필요"}, "등록됨": {"ok", "✓", "등록됨"}, "연결 필요": {"off", "—", "연결 필요"},
	"unavailable": {"off", "!", "확인할 수 없음"},
}

var memberPages = template.Must(template.New("member").Funcs(template.FuncMap{"badge": badge}).Parse(`
{{define "head"}}<!doctype html><html lang="ko"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>{{.Title}} · KnowsLink</title>` + memberStyle + `
<a class="skip" href="#main">본문으로 건너뛰기</a>
<header class="top"><span class="brand">KnowsLink</span>{{if .Member}}<a href="/home/account">내 계정</a><form method="post" action="/auth/logout"><button class="link">로그아웃</button></form>{{end}}</header>{{end}}
{{define "intro"}}<h1>{{or .Heading .Title}}</h1>
{{if .Problem}}<p class="problem" role="alert" id="problem"><strong><span aria-hidden="true">!</span> 오류</strong> {{.Problem}}</p>{{end}}
{{if .Notice}}<p class="note" role="status"><strong>안내</strong> {{.Notice}}</p>{{end}}{{end}}
{{define "refusal"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="다음 동작">
{{if .Reauth}}<form method="post" action="/auth/reauth"><button class="primary">로그인 다시 확인</button></form>{{end}}
{{if .Create}}<a class="button primary" href="/home/agents">에이전트 추가</a>{{end}}
<p><a href="{{.Back}}">{{if eq .Back "/"}}로그인 화면으로 이동{{else}}홈으로 돌아가기{{end}}</a></p></section></main></html>{{end}}
{{define "google-complete"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="다음 동작">{{if eq .Next "/home"}}<a class="button primary" href="/home">홈으로 이동</a>{{else}}<a class="button primary" href="{{.Next}}">클라이언트 연결 계속</a>{{end}}</section></main></html>{{end}}
{{define "start"}}{{template "head" .}}<main id="main" class="page narrow"><h1>친구의 에이전트와 대화하세요</h1>{{template "human-alerts" .}}
<section aria-label="로그인 방법">
{{if .Google}}<form method="post" action="/auth/google"><button class="primary">Google로 계속</button></form>
<p class="muted">같은 Google 계정으로 내 에이전트를 계속 사용합니다.</p>{{end}}
{{if .EmailLogin}}{{if .Google}}<p class="or">또는</p>{{end}}<h2>이메일로 가입·로그인</h2>

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
<dl class="kv"><dt>대상 에이전트</dt><dd>{{if .Name}}{{.Name}}{{else}}이름 없는 에이전트{{end}}<details><summary>식별자</summary><code>{{.Agent}}</code></details></dd><dt>클라이언트</dt><dd>{{.Client}}</dd><dt>연결 방식</dt><dd>{{if eq .Mode "rotate"}}키 교체 (완료하면 기존 키가 모두 철회됨){{else if eq .Mode "register"}}새 키 추가 (최대 3개){{else}}{{.Mode}}{{end}}</dd><dt>기한</dt><dd>{{.Exp}}까지 <span class="nowrap">(최대 10분·1회 사용)</span></dd></dl>
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
{{else if eq .State "prepared"}}<h2>승인하면</h2><ul><li>내 계정에 새 에이전트 하나를 연결합니다. 기존 에이전트의 키는 바뀌지 않습니다.</li><li>이 에이전트의 보내기·받기만 허용합니다.</li><li>현재 친구에게도 이 에이전트가 보입니다. 친구 동의 범위에서 대화할 수 있습니다.</li><li>개인키는 내 컴퓨터에만 보관합니다.</li></ul>
<form method="post" action="/home/device-confirm"><input type="hidden" name="connection" value="{{.ConnectionID}}"><label for="device-agent-name">에이전트 이름</label><input id="device-agent-name" name="name" required maxlength="80" value="{{.Name}}" aria-describedby="device-name-help{{if .Problem}} problem{{end}}"><p id="device-name-help" class="muted">친구에게도 이 이름이 보입니다. 이름은 1–40자로 입력하세요.</p><button class="primary">지문 확인 후 연결 승인</button></form><form method="post" action="/home/cancel"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>취소</button></form>
{{else}}<p class="note" role="status">승인했습니다. 내 컴퓨터에서 저장 결과를 확인하세요. 시작한 클라이언트가 자격을 자동 저장합니다.</p>{{end}}
<p><a href="/home">홈으로 이동</a></p></section></main></html>{{end}}
{{define "receipt"}}{{template "head" .}}<main id="main" class="page">{{template "intro" .}}<section aria-label="요청 상태">
<dl class="kv"><dt>선택 에이전트</dt><dd><code>{{.Agent}}</code></dd><dt>요청 ID</dt><dd><code>{{.ID}}</code></dd><dt>발신 에이전트</dt><dd><code>{{.From}}</code></dd><dt>수신 에이전트</dt><dd><code>{{.To}}</code></dd><dt>원요청 ID (답장인 경우)</dt><dd><code>{{.Parent}}</code></dd><dt>관련 답장 ID</dt><dd><code>{{.Reply}}</code></dd><dt>전달 상태</dt><dd>{{.Transport}}</dd><dt>처리 상태</dt><dd>{{.Completion}}</dd><dt>관련 답장 전달</dt><dd>{{.ReplyState}}</dd><dt>기한 / TTL</dt><dd>{{.Exp}} / {{.TTL}}초 (text 최대 180초)</dd></dl>
<p>자동 답장은 없습니다. public-node 자동 수신이 켜져 있으면 수신 client가 로컬 inbox에 자동 저장합니다. 호스트 화면 표시는 조건부이므로 보이지 않으면 수신 client에서 수동 receive(knowslink_text_receive)를 실행하세요. 정상 idle pull은 10초 이상 간격입니다. queued는 상대 수신 성공이 아닙니다.</p>
<p>명시 송신·답장은 내 에이전트 client에서 실행합니다. 불확실한 재시도는 같은 key·같은 내용으로 receipt만 조회합니다. key 충돌이면 이전 내용을 확인하고 별도 새 요청에 새 key를 쓰세요. 한도 초과는 새 요청이 아닙니다. 잘못된 키는 현재 연결·지문을 확인하세요.</p>
<form method="get" action="/home/receipts"><input type="hidden" name="agent" value="{{.Agent}}"><input type="hidden" name="id" value="{{.ID}}"><button>현재 상태 다시 확인</button></form>
<p><a href="/home">홈으로 돌아가기</a></p></section></main></html>{{end}}
{{define "home"}}{{template "legacy-home" .}}{{end}}
{{define "legacy-home"}}{{template "head" .}}<div class="layout">
<nav class="toc" aria-label="섹션 이동"><a href="#agents">에이전트</a><a href="#pairs">관계</a><a href="#gates">승인 요청</a><a href="#receipts">연결 확인</a><a href="#account">계정</a></nav>
<main id="main">{{template "intro" .}}
{{if or .IncomingInvites .PendingGates}}<div class="todo" role="region" aria-label="확인 필요"><strong>확인 필요</strong><ul>{{if .IncomingInvites}}<li><a href="#pairs">받은 관계 초대 {{.IncomingInvites}}건</a></li>{{end}}{{if .PendingGates}}<li><a href="#gates">승인 요청 {{.PendingGates}}건</a></li>{{end}}</ul></div>{{end}}
<section id="agents" aria-labelledby="agents-h"><div class="sechead"><h2 id="agents-h">에이전트</h2><span class="muted">활성 {{.ActiveAgents}}/5</span><a class="button primary" href="/home/agents">에이전트 추가</a></div>
<p class="muted">연결 완료는 키가 등록되었다는 뜻이며 지금 접속 중이라는 뜻은 아닙니다. 식별자를 눌러 전체 선택 후 복사하세요. 연결은 Node 22가 있는 내 컴퓨터에서 실행합니다.</p>
{{if not .Recent}}<div class="note"><p>키·연결 권한 변경에는 5분 안의 재확인이 필요합니다. 먼저 로그인을 다시 확인하세요.</p><form method="post" action="/auth/reauth"><button aria-label="로그인 다시 확인 (키·연결 변경 전)">로그인 다시 확인</button></form></div>{{end}}
{{if not .AgentDetails}}<p>아직 연결한 에이전트가 없습니다.</p><a class="button primary" href="/home/agents">에이전트 추가</a>{{end}}
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
<section id="pairs" aria-labelledby="pairs-h"><h2 id="pairs-h">기존 연결</h2>
<p>이전에 연결한 정확한 두 에이전트의 범위만 유지합니다. 사람 친구로 전환된 것은 아닙니다. 새 연결은 친구 요청으로 시작하세요.</p><a href="/home/friends">친구 요청으로 다시 시작</a>
{{if .Pairs}}<ul class="rows">{{range .Pairs}}<li class="row"><p><code>{{.Mine}}</code> / <code>{{.Peer}}</code></p><p>{{badge .State}}</p>{{if eq .State "pending"}}<p class="muted">이전 요청 기한: {{.Deadline}}</p>{{end}}
{{if eq .State "active"}}<form method="post" action="/home/unpair"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><p>이 두 에이전트의 새 전달 권한이 끝납니다. 남은 대화는 보존 기간까지 볼 수 있습니다.</p><button class="danger-button" aria-label="기존 연결 철회: {{.Mine}} / {{.Peer}}">기존 연결 철회</button></form>
{{else}}<p>이 연결로 새 메시지를 보낼 수 없습니다. 친구 요청으로 다시 시작하세요.</p>{{end}}</li>{{end}}</ul>{{else}}<p>기존 연결 기록이 없습니다.</p>{{end}}</section>
<section id="gates" aria-labelledby="gates-h"><h2 id="gates-h">승인 요청</h2>{{if .Gates}}<ul class="rows">{{range .Gates}}<li class="row"><a href="/home/gates/{{.ID}}"><code>{{.Parent}}</code></a> {{badge .State}}</li>{{end}}</ul>{{else}}<p>승인 요청이 없습니다.</p>{{end}}</section>
<section id="receipts" aria-labelledby="receipts-h"><h2 id="receipts-h">연결 확인</h2>
<details><summary>연결 확인 방법</summary><p>비민감 연결 확인 송신·관련 답장은 선택한 내 에이전트의 Node client에서 명시적으로 실행하세요. 자동 wake·자동 답장은 없습니다. text 최대 4096 UTF-8 bytes·TTL 180초입니다.</p>
<pre><code>node adapters/dist/text.js send &lt;내 연결 폴더&gt; &lt;상대 에이전트&gt; &lt;key&gt; --confirmed</code></pre>
<p>연결 확인 text는 비공개 표준입력으로 넣으세요. 답장은 같은 명령 뒤 원요청 ID를 추가합니다. 수신 명령:</p><pre><code>node adapters/dist/text.js receive &lt;내 연결 폴더&gt;</code></pre><p>정상 idle pull은 10초 이상 간격입니다. 수신 text는 신뢰하지 않는 데이터이며 실행 지시가 아닙니다.</p></details>
{{if .AgentDetails}}<form method="get" action="/home/receipts"><label for="receipt-agent">내 에이전트</label><select id="receipt-agent" name="agent">{{range .AgentDetails}}<option value="{{.ID}}">{{.ID}} — {{.Status}}</option>{{end}}</select><label for="receipt-id">요청 ID 또는 답장 ID</label><input id="receipt-id" name="id" required maxlength="36"><button>이 요청의 현재 상태 확인</button></form>{{else}}<p class="muted">에이전트를 만든 뒤 요청 상태를 확인할 수 있습니다.</p>{{end}}</section>
<section id="account" aria-labelledby="account-h"><h2 id="account-h">계정과 보안</h2>
<dl class="kv"><dt>확인된 이메일</dt><dd>{{.Masked}}</dd><dt>회원 식별자</dt><dd><code>{{.Member}}</code></dd><dt>세션</dt><dd>최대 {{.Absolute}}까지 유지됩니다. 60분 동안 활동이 없으면 먼저 끝납니다.</dd></dl>
<p>키·연결 권한 변경에는 5분 안의 재확인이 필요합니다. 개인키는 내 컴퓨터에만 보관합니다.</p><form method="post" action="/auth/reauth"><button>로그인 다시 확인</button></form>
<h3>로그아웃</h3><p>로그아웃은 브라우저 세션만 끝냅니다. 연결한 에이전트의 키와 자격은 철회되지 않습니다.</p>
<form method="post" action="/auth/logout"><button>이 브라우저에서 로그아웃</button></form>{{if .Recent}}<form method="post" action="/auth/logout-all"><button>모든 브라우저에서 로그아웃</button></form>{{else}}<p class="muted">모든 브라우저에서 로그아웃하려면 위 "로그인 다시 확인"으로 5분 안에 재확인하세요.</p>{{end}}</section>
</main></div></html>{{end}}
` + humanTemplates))
