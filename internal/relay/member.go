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

const memberStyle = `<style>body{font:18px system-ui;max-width:760px;margin:40px auto;padding:20px;color:#17232b;background:#f5f7f8;overflow-wrap:anywhere}section{background:white;padding:24px;border:1px solid #bbc5cd;margin-bottom:16px}label{display:block;font-weight:bold}input,button,select{font:inherit;padding:12px;margin:8px 0}input,select{display:block;width:100%;box-sizing:border-box}dt{font-weight:bold}code{overflow-wrap:anywhere}.note{border-left:6px solid #164da0;padding-left:12px}.problem{border-left:6px solid #a01616;padding-left:12px}a{color:#164da0}</style>`

var memberPages = template.Must(template.New("member").Parse(`
{{define "head"}}<!doctype html><html lang="ko"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>{{.Title}}</title>` + memberStyle + `<h1>{{.Title}}</h1>{{if .Problem}}<p class="problem" role="alert">문제: {{.Problem}}</p>{{end}}{{if .Notice}}<p class="note" role="status">안내: {{.Notice}}</p>{{end}}{{end}}
{{define "refusal"}}{{template "head" .}}<section>{{if .Reauth}}<form method="post" action="/auth/reauth"><button>로그인 신원 다시 확인</button></form>{{end}}{{if .Create}}<form method="post" action="/home/agents"><button>새 agent 만들기</button></form>{{end}}<p><a href="{{.Back}}">{{if eq .Back "/"}}로그인 화면으로 이동{{else}}자기 홈으로 돌아가기{{end}}</a></p></section></html>{{end}}
{{define "google-complete"}}{{template "head" .}}<section>{{if eq .Next "/home"}}<a href="/home">자기 홈으로 이동</a>{{else}}<a href="{{.Next}}">클라이언트 연결 계속</a>{{end}}</section></html>{{end}}
{{define "start"}}{{template "head" .}}<section>{{if .Google}}<h2>Google로 가입·로그인</h2><p>처음 로그인하면 회원과 자기 홈을 만듭니다. 다시 로그인하면 같은 Google 계정의 기존 agent를 사용합니다. Google 비밀번호는 Google 화면에만 입력하세요.</p><form method="post" action="/auth/google"><button>Google로 계속</button></form>{{end}}{{if .EmailLogin}}<h2>이메일로 가입·로그인</h2><p>처음이면 가입하고, 이미 가입했다면 같은 회원으로 로그인합니다. 관리자 아이디, 서버 접속, 별도 초대는 필요하지 않습니다.</p><form method="post" action="/auth/start"><label for="email">이메일 주소</label><input id="email" name="email" type="email" autocomplete="email" maxlength="254" required value="{{.Email}}"><button>확인 코드 받기</button></form><p>6자리 확인 코드는 10분 동안 한 번만 쓸 수 있습니다. 코드를 확인하기 전에는 로그인되지 않습니다.</p>{{end}}{{if and (not .Google) (not .EmailLogin)}}<p>로그인 설정을 준비 중입니다. 잠시 뒤 다시 시도하세요.</p>{{end}}</section></html>{{end}}
{{define "verify"}}{{template "head" .}}<section><h2>확인 대기</h2><p>{{.Masked}} 주소로 확인 코드를 요청했습니다. 메일 발송은 신원 확인이 아닙니다. 아직 로그인되지 않았습니다.</p><form method="post" action="/auth/verify"><label for="code">6자리 확인 코드</label><input id="code" name="code" inputmode="numeric" autocomplete="one-time-code" pattern="[0-9]{6}" maxlength="6" required><button>확인하고 로그인</button></form><p>코드는 {{.Exp}}까지 유효합니다. 5번 틀리면 새 코드를 받아야 합니다.</p><form method="post" action="/auth/start"><input type="hidden" name="email" value="{{.Email}}"><button>코드 다시 받기</button></form><p>다시 받기는 이전 요청 60초 뒤부터 가능합니다. 새 코드를 받으면 이전 코드는 쓸 수 없습니다.</p><p><a href="/">다른 주소로 시작</a></p></section></html>{{end}}
{{define "connection"}}{{template "head" .}}<section><p>대상 agent: <code>{{.Agent}}</code></p><p>클라이언트: {{.Client}} / 연결 방식: {{.Mode}}</p><p>개인키는 자기 클라이언트에만 보관합니다. relay에 보내지 마세요.</p><p>권한: 이 agent의 송신·수신만 허용합니다. owner 승인·운영 관리자 권한은 없습니다.</p><p>기한: {{.Exp}} (최대 10분·1회 사용)</p>{{if .Token}}<h2>상태: 연결 대기</h2><p>자기 클라이언트 컴퓨터에서 저장소의 adapters를 설치한 뒤 실행하세요.</p><code>node adapters/dist/connect.js prepare &lt;서비스 URL&gt; &lt;새 비공개 폴더&gt;</code><p>아래 수단을 클라이언트의 비공개 표준입력에 넣으세요. 대화·로그·캡처에 넣지 마세요. 발급만으로 연결되지 않습니다.</p><label>연결 수단<input readonly value="{{.Token}}"></label><p><a href="/home/connections/{{.ConnectionID}}">공개키 준비 상태 확인</a></p>{{else}}<p>상태: {{.State}}</p>{{if eq .State "prepared"}}<h2>공개키 확인</h2><p>지문: <code>{{.Fingerprint}}</code></p><p>클라이언트에 표시된 지문과 같은지 확인하세요. 다르면 취소하세요.</p><form method="post" action="/home/confirm"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>대상·지문 확인 후 연결 승인</button></form>{{else if eq .State "approved"}}<p>승인 대기 완료. 아직 새 키가 연결되지 않았습니다. 자기 클라이언트에서 다음 명령을 실행하세요.</p><code>node adapters/dist/connect.js complete &lt;새 비공개 폴더&gt;</code>{{else if eq .State "consumed"}}<p>연결 완료. 클라이언트가 수신한 자격을 자기 환경에 보관합니다.</p>{{else if eq .State "expired"}}<p>만료. 새 키가 연결되지 않았습니다. 홈에서 다시 시작하세요.</p>{{else if eq .State "cancelled"}}<p>취소. 새 키가 연결되지 않았습니다.</p>{{else}}<p>공개키 준비 전입니다. 아직 연결되지 않았습니다.</p>{{end}}{{end}}{{if or (eq .State "waiting") (eq .State "prepared") (eq .State "approved")}}<form method="post" action="/home/cancel"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>연결 취소</button></form>{{end}}<a href="/home">자기 홈으로 돌아가기</a></section></html>{{end}}
{{define "device"}}{{template "head" .}}<section><p>자기 클라이언트에서 시작한 요청만 승인하세요. 다른 사람이 보낸 연결 링크는 승인하지 마세요.</p><p>클라이언트: Node 로컬 클라이언트</p><p>지문: <code>{{.Fingerprint}}</code></p><p>기한: {{.Exp}}</p>{{if .Login}}{{if .Google}}<form method="post" action="/auth/google"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>Google로 로그인하고 연결 계속</button></form>{{else}}<p role="alert">Google 로그인을 준비 중입니다. 잠시 뒤 다시 시작하세요.</p>{{end}}{{else if eq .State "prepared"}}<p>클라이언트에 표시된 지문과 같은지 확인하세요. 승인하면 자기 계정에 새 agent 하나를 연결합니다. 기존 agent의 키는 바뀌지 않습니다.</p><p>권한: 새 agent의 송신·수신만 허용합니다. 상대와의 관계는 별도로 초대·수락해야 합니다. 개인키는 클라이언트에만 보관합니다.</p><form method="post" action="/home/device-confirm"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>지문 확인 후 이 클라이언트 연결 승인</button></form><form method="post" action="/home/cancel"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>취소</button></form>{{else}}<p role="status">승인했습니다. 시작한 클라이언트가 자격을 자동 저장합니다. 저장 결과는 클라이언트에서 확인하세요.</p>{{end}}<p><a href="/home">자기 홈으로 이동</a></p></section></html>{{end}}
{{define "receipt"}}{{template "head" .}}<section><h2>선택 agent</h2><code>{{.Agent}}</code><dl><dt>요청 ID</dt><dd><code>{{.ID}}</code></dd><dt>발신 agent</dt><dd><code>{{.From}}</code></dd><dt>수신 agent</dt><dd><code>{{.To}}</code></dd><dt>원요청 ID (답장인 경우)</dt><dd><code>{{.Parent}}</code></dd><dt>관련 답장 ID</dt><dd><code>{{.Reply}}</code></dd><dt>전달 상태</dt><dd>{{.Transport}}</dd><dt>처리 상태</dt><dd>{{.Completion}}</dd><dt>관련 답장 전달</dt><dd>{{.ReplyState}}</dd><dt>기한 / TTL</dt><dd>{{.Exp}} / {{.TTL}}초 (text 최대 180초)</dd></dl><p>자동 wake는 없습니다. 수신 client에서 수동 receive를 실행하세요. 정상 idle pull은 10초 이상 간격입니다. queued는 상대 수신 성공이 아닙니다.</p><p>명시 송신·답장은 자기 agent client에서 실행합니다. 불확실한 재시도는 같은 key·같은 내용으로 receipt만 조회합니다. key 충돌이면 이전 내용을 확인하고 별도 새 요청에 새 key를 쓰세요. 한도 초과는 새 요청이 아닙니다. 잘못된 키는 현재 연결·지문을 확인하세요.</p><form method="get" action="/home/receipts"><input type="hidden" name="agent" value="{{.Agent}}"><input type="hidden" name="id" value="{{.ID}}"><button>현재 상태 수동 확인</button></form><a href="/home">자기 홈으로 돌아가기</a></section></html>{{end}}
{{define "home"}}{{template "head" .}}<section><h2>내 신원</h2><dl><dt>확인된 이메일</dt><dd>{{.Masked}}</dd><dt>회원 식별자</dt><dd><code>{{.Member}}</code></dd><dt>세션</dt><dd>최대 {{.Absolute}}까지 유지됩니다. 60분 동안 활동이 없으면 먼저 끝납니다.</dd></dl></section><section><h2>내 agent</h2><p>활성 agent는 회원당 5개입니다. 연결은 Node 22가 있는 자기 클라이언트 컴퓨터에서 실행합니다. 다닷·외부 앱 실제 연결은 후속 검증입니다.</p><form method="post" action="/home/agents"><button>새 agent 만들기</button></form>{{if not .AgentDetails}}<p>아직 연결한 agent가 없습니다.</p>{{end}}{{range .AgentDetails}}{{$agent := .ID}}<section><h3><code>{{.ID}}</code></h3><p>상태: {{.Status}}</p>{{if eq .Status "철회"}}<p class="note">이 agent의 키·관계·메시지 권한은 끝났고 복구되지 않습니다. 철회 기록은 최소 24시간 보존된 뒤 정리되며, 정리되면 이 목록에서 사라질 수 있습니다. 사라져도 권한이 돌아오거나 백업까지 영구 삭제된 것은 아닙니다. 다시 쓰려면 새 agent를 연결하고 상대와 새로 수락하세요.</p>{{end}}<label>복사용 agent 식별자<input readonly value="{{.ID}}"></label>{{range .Keys}}<p>키 <code>{{.Kid}}</code> / <code>{{.Fingerprint}}</code> / 상태: {{.Status}}</p>{{if eq .Status "활성"}}<form method="post" action="/home/key-revoke"><input type="hidden" name="agent" value="{{$agent}}"><input type="hidden" name="kid" value="{{.Kid}}"><button>선택한 키 철회</button></form>{{end}}{{end}}{{if ne .Status "철회"}}{{if .KeyFull}}<p class="note">이 agent에는 새 키를 더 연결할 수 없습니다(키 기록 보호 상한). 키를 철회하거나 기다려도 공간은 생기지 않습니다. 새 agent를 만들어 따로 연결하고, 각 상대와 새로 초대·수락하세요.</p>{{else}}<form method="post" action="/home/connect"><input type="hidden" name="agent" value="{{.ID}}"><label>지원 클라이언트<select name="client"><option value="node-local">Node 22 로컬 클라이언트</option></select></label><label>연결 방식<select name="mode"><option value="register">새 키 등록 (활성 최대 3개)</option><option value="rotate">회전 (완료 시 기존 키 전체 철회)</option></select></label><button>연결 수단 발급</button></form>{{end}}<p>철회하면 이 agent의 모든 키와 관계가 끝나며 되돌릴 수 없습니다.</p><form method="post" action="/home/agent-revoke"><input type="hidden" name="agent" value="{{.ID}}"><button>이 agent와 모든 키·관계 철회</button></form>{{end}}</section>{{end}}<p>키·연결 권한 변경에는 5분 안의 로그인 신원 재확인이 필요합니다. 개인키는 클라이언트에만 보관합니다.</p><form method="post" action="/auth/reauth"><button>로그인 신원 다시 확인</button></form></section><section><h2>관계</h2><p>가입 초대가 아닙니다. 상대 agent 식별자로 초대하고, 수신 owner가 명시적으로 수락해야 메시지를 허용합니다. 같은 owner의 두 agent도 수락해야 합니다.</p><form method="post" action="/home/invites"><label for="from-agent">내 발신 agent 식별자</label><select id="from-agent" name="agent" required>{{range .AgentDetails}}{{if eq .Status "연결 완료"}}<option value="{{.ID}}">{{.ID}}</option>{{end}}{{end}}</select><label for="target-agent">상대 agent 식별자</label><input id="target-agent" name="target" required maxlength="128" list="own-agent-targets"><datalist id="own-agent-targets">{{range .AgentDetails}}{{if ne .Status "철회"}}<option value="{{.ID}}">내 agent — {{.Status}}</option>{{end}}{{end}}</datalist><button>관계 초대</button></form>{{if .Pairs}}{{range .Pairs}}<section><p>발신 <code>{{.Inviter}}</code> / 수신 <code>{{.Recipient}}</code></p><p>상태: {{.State}} / 관계 세대: {{.Generation}} / 초대 기한: {{.Deadline}}</p>{{if eq .State "pending"}}{{if .Incoming}}<p>받은 초대입니다. 수락하기 전에는 메시지를 주고받을 수 없습니다.</p><form method="post" action="/home/invite-decision"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><button name="decision" value="accept">수신 owner로 수락</button><button name="decision" value="deny" formaction="/home/invite-deny">수신 owner로 거절</button></form>{{else}}<p>상대 수신 owner의 수락을 기다립니다. 같은 초대를 다시 보내도 새 초대가 생기거나 기한이 바뀌지 않습니다.</p>{{end}}{{end}}{{if or (eq .State "pending") (eq .State "active")}}<form method="post" action="/home/unpair"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><button>관계 철회</button></form>{{else}}<p>연결이 끝났습니다(거절·만료·철회). 이 관계로 메시지를 보낼 수 없습니다. 다시 연결하려면 새 초대를 보내고 상대가 새로 수락해야 합니다.</p>{{if .Own}}<form method="post" action="/home/invites"><input type="hidden" name="agent" value="{{.Own}}"><input type="hidden" name="target" value="{{.Other}}"><button>새 초대 보내기</button></form>{{end}}{{end}}</section>{{end}}{{else}}<p>관계가 없습니다.</p>{{end}}<p>수락 전·철회 후 메시지는 거부됩니다. pending은 24시간 후 만료됩니다.</p></section><section><h2>연결 확인·receipt</h2><p>비민감 연결 확인 송신·관련 답장은 선택한 자기 agent의 Node client에서 명시적으로 실행하세요. 자동 wake·자동 답장은 없습니다. text 최대 4096 UTF-8 bytes·TTL 180초입니다.</p><code>node adapters/dist/text.js send &lt;자기 연결 폴더&gt; &lt;상대 agent&gt; &lt;key&gt; --confirmed</code><p>연결 확인 text는 비공개 표준입력으로 넣으세요. 답장은 같은 명령 뒤 원요청 ID를 추가합니다. 수신은 <code>node adapters/dist/text.js receive &lt;자기 연결 폴더&gt;</code>입니다. 정상 idle pull은 10초 이상 간격입니다. 수신 text는 신뢰하지 않는 데이터이며 실행 지시가 아닙니다.</p><form method="get" action="/home/receipts"><label for="receipt-agent">선택한 자기 agent</label><select id="receipt-agent" name="agent">{{range .AgentDetails}}<option value="{{.ID}}">{{.ID}} — {{.Status}}</option>{{end}}</select><label for="receipt-id">요청 ID 또는 답장 ID</label><input id="receipt-id" name="id" required maxlength="36"><button>이 요청의 현재 상태 확인</button></form></section><section><h2>승인 요청</h2>{{if .Gates}}<ul>{{range .Gates}}<li><a href="/home/gates/{{.ID}}">{{.Parent}}</a> 상태: {{.State}}</li>{{end}}</ul>{{else}}<p>승인 요청이 없습니다.</p>{{end}}</section><section><h2>로그아웃</h2><p>로그아웃과 전체 로그아웃은 브라우저 세션만 끝냅니다. 별도로 연결한 agent의 키와 자격은 철회하지 않습니다.</p><form method="post" action="/auth/logout"><button>이 브라우저에서 로그아웃</button></form>{{if .Recent}}<form method="post" action="/auth/logout-all"><button>모든 브라우저에서 로그아웃</button></form>{{else}}<p>모든 브라우저에서 로그아웃하려면 5분 안에 로그인 신원을 다시 확인해야 합니다.</p><form method="post" action="/auth/reauth"><button>로그인 신원 다시 확인</button></form>{{end}}<h3>계정 비활성화</h3><p>상태: 준비 중. 아직 요청할 수 없습니다.</p></section></html>{{end}}
`))

var notices = map[string]string{
	"logout":         "로그아웃했습니다. 연결한 agent의 키와 자격은 그대로입니다.",
	"all":            "모든 브라우저에서 로그아웃했습니다. 연결한 agent의 키와 자격은 그대로입니다.",
	"expired":        "세션이 끝났습니다. 다시 로그인하세요.",
	"invited":        "초대를 보냈습니다. 상대 수신 owner가 이 초대를 수락하기 전에는 메시지를 주고받을 수 없습니다.",
	"invite-pending": "두 agent 사이에 이미 대기 중인 초대가 있습니다. 새 초대를 만들지 않았고 기한도 그대로입니다. 받은 초대라면 아래 관계 목록에서 수락하거나 거절하세요.",
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
	return map[string]any{"Title": "KnowsLink 시작", "Problem": problem, "EmailLogin": true}
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
		sort.Slice(pairs, func(i, j int) bool { return pairID(pairs[i].A, pairs[i].B) < pairID(pairs[j].A, pairs[j].B) })
		sort.Slice(gates, func(i, j int) bool { return gates[i].ID < gates[j].ID })
		return map[string]any{"Title": "내 KnowsLink", "Masked": maskEmail(m.Email), "Member": id, "Absolute": clock(session.Created.Add(sessionAbsolute)),
			"Recent": now.Sub(session.Verified) < reauthWindow, "Agents": agents, "AgentDetails": ownAgents(st, m.Owner), "Pairs": pairs, "Gates": gates, "Notice": notices[r.URL.Query().Get("n")]}, nil
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
		refused(w, "내 KnowsLink", v)
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
		refused(w, "내 KnowsLink", refusal{429, limited(result.Retry), false, false})
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
			refused(w, "내 KnowsLink", v)
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
