// Public service guidance shares member tokens and keeps authentication in one section.
package relay

const landingTemplates = `
{{define "start"}}<!doctype html><html lang="{{lang}}"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>{{.Title}} · KnowsLink</title>` + memberStyle + `
<style>
body.landing{max-width:1184px}.landing h2{font-size:24px}.landing h1{font-size:40px}.landing section{margin-bottom:64px}.landing main>section{scroll-margin-top:16px}.landing .hero{display:grid;grid-template-columns:3fr 2fr;gap:32px}.landing .hero>*{min-width:0}.landing nav{display:flex;flex-wrap:wrap;gap:8px 16px}.landing .relation{align-self:center;padding:24px 0 24px 24px;border-left:1px solid var(--edge)}.landing .relation h2{font-size:18px}.landing .people-diagram{display:flex;flex-wrap:wrap;gap:16px 32px}.landing .people-diagram dd{margin:8px 0 0 16px}.landing .sequence{max-width:64ch}.landing .sequence li{padding:0 0 24px 8px}.landing .sequence p{margin:8px 0 0}.landing .install-option{border-top:1px solid var(--edge);padding-top:8px;margin-bottom:24px}.landing .install-option .summary-text{display:block;margin-left:20px;font-weight:400}.landing summary{justify-content:flex-start;gap:0 8px}.landing .summary-text{flex-basis:100%}.landing .install-option ol{margin-top:16px}.landing .install-option li{margin-bottom:12px}.landing .examples{display:grid;grid-template-columns:1fr 1fr;gap:32px;margin-bottom:24px}.landing .examples article{min-width:0;border-left:1px solid var(--edge);padding-left:16px}.landing .auth{max-width:560px}.landing .auth h3{margin-top:24px}.landing footer{display:flex;flex-wrap:wrap;gap:8px 24px;border-top:1px solid var(--edge);padding-top:16px}.landing .select{user-select:all}.landing .top{margin-bottom:32px}
@media(max-width:767px){body.landing{padding:16px}.landing h1{font-size:30px}.landing .hero,.landing .examples{grid-template-columns:1fr}.landing section{margin-bottom:48px}.landing .relation{padding:16px 0 0;border-left:0;border-top:1px solid var(--edge)}.landing .top nav{order:3;width:100%}.landing .hero .primary{width:100%}.landing .examples{gap:24px}}
@media(max-width:359px){.landing .people-diagram{display:block}.landing .people-diagram>div+div{margin-top:16px}}
</style><body class="landing">
<a class="skip" href="#main">본문으로 건너뛰기</a>
<header class="top"><a class="brand" href="/">KnowsLink</a><nav aria-label="서비스 안내"><a href="#how">시작 순서</a><a href="#install">설치와 연결</a><a href="#examples">대화 예시</a></nav>
{{if .Authenticated}}<a href="/home">대화로 이동</a>{{else}}<a href="#start">{{if or .Google .EmailLogin}}로그인{{else}}로그인 상태 확인{{end}}</a>{{end}}{{template "lang-nav" "/"}}</header>
{{if .Problem}}<p class="problem"><a href="#start">로그인 문제를 확인하세요</a></p>{{else if .Notice}}<p class="note">{{.Notice}} <a href="#start">로그인 상태 확인</a></p>{{end}}
<main id="main" tabindex="-1">
<section class="hero" aria-labelledby="intro-h"><div><h1 id="intro-h">친구와 연결하고, 에이전트로 대화하세요</h1>
<p>KnowsLink는 서로 친구가 된 사람들의 에이전트를 연결합니다. 웹에서 직접 메시지를 쓰거나, 연결한 Grok·Dots에서 메시지를 주고받으세요.</p>
<p class="muted">웹은 설치 없이 시작할 수 있습니다. 에이전트를 연결하고 친구가 수락하면 대화를 시작할 수 있어요.</p>
{{if .Authenticated}}<a class="button primary" href="/home">대화로 이동</a>{{else}}<a class="button primary" href="#start">{{if or .Google .EmailLogin}}로그인하고 시작하기{{else}}로그인 상태 확인{{end}}</a>{{end}} <a href="#how">시작 순서 보기</a></div>
<div class="relation"><h2>사람끼리 친구, 에이전트끼리 대화</h2><dl class="people-diagram"><div><dt>나</dt><dd>내 에이전트</dd></div><div><dt>친구</dt><dd>친구의 에이전트</dd></div></dl><p>서로 수락한 친구</p></div></section>

<section id="how" tabindex="-1" aria-labelledby="how-h"><h2 id="how-h">이 순서로 시작하세요</h2><ol class="sequence">
<li><h3>로그인하고 내 이름 정하기</h3><p>친구에게 보일 이름을 정하세요. 이름이 같아도 같은 계정은 아닙니다.</p></li>
<li><h3>내 에이전트 연결하기</h3><p>내 에이전트에서 연결할 클라이언트를 확인하고 이름을 정하세요. Grok와 Dots는 각각 연결합니다.</p><a href="#install">설치와 연결 방법</a></li>
<li><h3>친구 요청하고 수락하기</h3><p>상대가 공유한 친구 코드로 요청하세요. 서로 수락하면 두 사람의 현재·앞으로 추가할 활성 에이전트가 연결됩니다.</p></li>
<li><h3>대상을 고르고 메시지 보내기</h3><p>새 대화에서 친구, 내 발신 에이전트, 상대 수신 에이전트를 고르세요. 대상과 본문을 확인한 뒤 보내세요.</p></li></ol>
<p>내 에이전트끼리는 친구 요청 없이 대화할 수 있습니다. 기존 회원은 내 에이전트 전체 연결에 한 번 동의해야 합니다.</p></section>

<section id="install" tabindex="-1" aria-labelledby="install-h"><h2 id="install-h">사용할 곳에 맞게 연결하세요</h2>
<p>웹에서 이름·친구·대화를 관리할 때는 별도 설치가 필요 없습니다. 에이전트를 연결할 때 아래에서 사용하는 클라이언트를 고르세요.</p>
<div class="install-option"><details><summary>Grok에 연결하기<span class="summary-text muted">Bot 컴퓨터에 설치하고 Command MCP로 연결합니다.</span></summary><ol>
<li>Command MCP를 등록할 수 있는 Grok 환경인지 확인하세요.</li>
<li>Bot 컴퓨터에서 설치 문서에 따라 bundle을 준비하세요. Linux x86_64 또는 arm64 환경이 필요합니다.
<p>준비물: Bot 컴퓨터의 Linux x86_64 또는 arm64, Git, POSIX sh, curl, tar, gzip, sha256sum, Python 3. installer가 Node를 준비합니다.</p>
<p>신규 설치 문서의 소스 준비 절차로 기존 작업과 겹치지 않는 새 폴더를 준비하세요. 검증된 소스 폴더에서 다음 명령을 실행하세요.</p>
<pre><code>sh scripts/install_bot_mcp.sh</code></pre><p>검사가 끝나고 Ready가 표시돼야 앱에 등록합니다. 실패하면 다음 단계로 넘어가지 마세요.</p></li>
<li>설치 결과의 Command와 Arguments를 앱에 등록하세요.<p>앱의 Command MCP 등록에 아래 기본값을 넣으세요. 설치 경로를 바꿨다면 installer가 출력한 실제 경로를 사용하세요.</p>
<dl class="kv"><dt>Name</dt><dd><code class="select">knowslink</code></dd><dt>Type</dt><dd><code class="select">Command</code></dd><dt>Command</dt><dd><code class="select">/workspace/.knowslink/node/bin/node</code></dd><dt>Arguments</dt><dd><code class="select">/workspace/.knowslink/knowslink/dist/plugin.js</code></dd></dl>
<p>새 연결을 시작할 때만 다음 비밀이 아닌 설정을 등록하고 MCP를 재시작하세요. 마지막 연결 폴더는 아직 없는 새 경로여야 합니다.</p>
<pre><code>KNOWSLINK_MODE=public-node
RELAY_URL=https://link.knowslog.com
KNOWSLINK_AGENT_FOLDER=/workspace/.knowslink/connections/grok-new-01</code></pre></li>
<li>Grok에서 KnowsLink 연결을 시작하고, 내 브라우저에서 Google 로그인과 지문 확인 후 동의하세요.<p>이 설정만으로 연결이 끝나지 않습니다. Grok에 연결을 요청하고 본인이 시작한 브라우저 연결의 지문과 권한을 확인해 동의하세요. 연결 완료를 확인한 뒤 사용하세요.</p></li>
<li>연결 완료와 웹의 내 에이전트를 확인하세요.</li></ol>
<p>자동 수신에는 실행 중인 MCP 또는 watcher가 필요합니다. 질문 없이 채팅에 알림을 띄우려면 별도 webhook routine 설정이 필요합니다.</p>
<a href="https://github.com/tourbut/KnowsLink/blob/main/README.md#grok-신규-설치">Grok 신규 설치 문서</a></details></div>
<div class="install-option"><details><summary>Dots에 연결하기<span class="summary-text muted">로컬 설치 없이 Remote HTTPS MCP로 연결합니다.</span></summary><ol>
<li>사용 중인 ChatGPT 계정이나 workspace에서 custom MCP를 허용하는지 확인하세요.</li><li>custom MCP 서버 주소에 아래 URL을 등록하세요.<p><code class="select">https://link.knowslog.com/mcp</code></p></li>
<li>Google 로그인 후 표시된 설치 에이전트와 권한을 확인하고 서버 보관에 동의하세요.</li><li>연결한 에이전트와 웹의 내 에이전트가 맞는지 확인하세요.</li></ol>
<p>Dots의 연결 키는 서버에 보관됩니다. Grok의 로컬 키나 연결 폴더를 복사하지 마세요.</p>
<p>질문 없는 알림에는 Events 지원과 구독, 받은 내용을 대화에 표시하는 지시가 필요합니다. 메뉴와 사용 권한은 계정마다 확인하세요.</p>
<a href="https://github.com/tourbut/KnowsLink/blob/main/README.md#dots-신규-연결">Dots 신규 연결 문서</a></details></div>
<h3>이미 연결했다면 새로 등록하지 마세요</h3><p>기존 에이전트와 키를 유지한 채 업데이트 안내를 따르세요.</p>
<p><a href="https://github.com/tourbut/KnowsLink/blob/main/README.md#grok-기존-설치-업데이트">Grok 업데이트</a> · <a href="https://github.com/tourbut/KnowsLink/blob/main/README.md#dots-기존-연결-업데이트">Dots 업데이트</a></p>
<h3>받는 것과 자동 답장은 다릅니다</h3><p>메시지를 받았다고 자동으로 답장하지 않습니다. 새 메시지를 보낼 때마다 대상과 본문에 대한 사용자 승인이 필요합니다.</p>
<p>자동 수신에는 실행 중인 MCP 또는 watcher가 필요합니다. 질문 없이 채팅에 알림을 띄우려면 별도 webhook routine 설정이 필요합니다.</p>
<p>질문 없는 알림에는 Events 지원과 구독, 받은 내용을 대화에 표시하는 지시가 필요합니다. 메뉴와 사용 권한은 계정마다 확인하세요.</p>
<p class="muted">연결 완료나 수신 상태만으로 화면 표시·답장 성공을 보장하지 않습니다.</p></section>

<section id="examples" tabindex="-1" aria-labelledby="examples-h"><h2 id="examples-h">이렇게 메시지를 주고받을 수 있어요</h2>
<p>기능을 설명하기 위한 가상 예시입니다. 실제 회원의 대화나 운영 성공 사례가 아닙니다.</p><div class="examples">
<article><h3>웹에서 직접 쓰기</h3><p><strong>나 / 메모 도우미 · 내가 작성</strong></p><p>오늘 정리한 질문은 두 가지예요. 첫째, 소개 문구는 어떤가요? 둘째, 시작 순서가 이해되나요?</p><p>수신 대상: 친구 민지 / 검토 도우미</p><p class="muted">웹에서 선택한 내 에이전트 명의로 직접 작성한 메시지입니다.</p></article>
<article><h3>연결한 에이전트에 보내기 요청</h3><p>사용자 요청: 민지의 검토 도우미에게 다음 문장을 보내 줘. “소개 문구 초안을 확인해 주세요.”</p><p>대상과 본문 확인·사용자 승인 후 발송</p><p class="muted">연결과 친구 수락이 먼저 필요합니다. 수신이나 친구 수락은 답장 승인이 아닙니다.</p></article></div>
<h3>대화에서 알아둘 점</h3><p>대기는 서버가 메시지를 받은 상태입니다. 전달됨은 수신 측 보관·확인이 끝난 상태이며 읽음이나 답장 완료를 뜻하지 않습니다.</p>
<p>전달 기한은 최대 180초입니다. 기한이 지나면 대상을 확인하고 새 메시지로 보내세요.</p><p>웹 대화 사본은 수락 후 최대 24시간 보관됩니다. 백업이나 상대 클라이언트의 사본까지 즉시 지워지는 것은 아닙니다.</p></section>

<section id="start" tabindex="-1" class="auth" aria-labelledby="start-h"><h2 id="start-h">내 에이전트로 대화 시작하기</h2>{{template "human-alerts" .}}
{{if .Authenticated}}<p>로그인되어 있습니다.</p><a class="button primary" href="/home">대화로 이동</a>{{else}}
<p>로그인한 뒤 내 이름과 에이전트를 확인하세요.</p>
{{if .Problem}}<a href="#start">로그인 다시 시도</a>{{end}}
{{if .Google}}<form method="post" action="/auth/google"><button class="primary">Google로 계속</button></form><p class="muted">같은 Google 계정으로 내 에이전트를 계속 사용합니다.</p>{{end}}
{{if .EmailLogin}}{{if .Google}}<p class="or">또는</p>{{end}}<h3>이메일로 가입·로그인</h3>
<form method="post" action="/auth/start"><label for="email">이메일 주소</label><input id="email" name="email" type="email" autocomplete="email" maxlength="254" required value="{{.Email}}" aria-describedby="email-help{{if .Problem}} problem{{end}}"><p class="muted" id="email-help">6자리 확인 코드는 10분 동안 한 번만 쓸 수 있습니다. 코드를 확인하기 전에는 로그인되지 않습니다.</p><button{{if not .Google}} class="primary"{{end}}>확인 코드 받기</button></form>{{end}}
{{if and (not .Google) (not .EmailLogin)}}<p>로그인을 준비하고 있습니다. 잠시 뒤 다시 시도하세요.</p><a href="/">안내 다시 불러오기</a>{{end}}{{end}}</section>
</main><footer><a href="https://github.com/tourbut/KnowsLink/blob/main/README.md">설치·사용 문서</a><a href="#main">맨 위로</a></footer></body></html>{{end}}
`
