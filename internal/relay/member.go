// Member pages run email verification, sessions, logout, and the self-owner home without owner Basic auth or SSH.
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

const memberStyle = `<style>body{font:18px system-ui;max-width:760px;margin:40px auto;padding:20px;color:#17232b;background:#f5f7f8}section{background:white;padding:24px;border:1px solid #bbc5cd;margin-bottom:16px}label{display:block;font-weight:bold}input,button{font:inherit;padding:12px;margin:8px 0}input{width:100%;box-sizing:border-box}dt{font-weight:bold}code{overflow-wrap:anywhere}.note{border-left:6px solid #164da0;padding-left:12px}.problem{border-left:6px solid #a01616;padding-left:12px}a{color:#164da0}</style>`

var memberPages = template.Must(template.New("member").Parse(`
{{define "head"}}<!doctype html><html lang="ko"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>{{.Title}}</title>` + memberStyle + `<h1>{{.Title}}</h1>{{if .Problem}}<p class="problem" role="alert">문제: {{.Problem}}</p>{{end}}{{if .Notice}}<p class="note" role="status">안내: {{.Notice}}</p>{{end}}{{end}}
{{define "start"}}{{template "head" .}}<section><h2>이메일로 가입·로그인</h2><p>처음이면 가입하고, 이미 가입했다면 같은 회원으로 로그인합니다. 관리자 아이디, 서버 접속, 별도 초대는 필요하지 않습니다.</p><form method="post" action="/auth/start"><label for="email">이메일 주소</label><input id="email" name="email" type="email" autocomplete="email" maxlength="254" required value="{{.Email}}"><button>확인 코드 받기</button></form><p>6자리 확인 코드는 10분 동안 한 번만 쓸 수 있습니다. 코드를 확인하기 전에는 로그인되지 않습니다.</p></section></html>{{end}}
{{define "verify"}}{{template "head" .}}<section><h2>확인 대기</h2><p>{{.Masked}} 주소로 확인 코드를 요청했습니다. 메일 발송은 신원 확인이 아닙니다. 아직 로그인되지 않았습니다.</p><form method="post" action="/auth/verify"><label for="code">6자리 확인 코드</label><input id="code" name="code" inputmode="numeric" autocomplete="one-time-code" pattern="[0-9]{6}" maxlength="6" required><button>확인하고 로그인</button></form><p>코드는 {{.Exp}}까지 유효합니다. 5번 틀리면 새 코드를 받아야 합니다.</p><form method="post" action="/auth/start"><input type="hidden" name="email" value="{{.Email}}"><button>코드 다시 받기</button></form><p>다시 받기는 이전 요청 60초 뒤부터 가능합니다. 새 코드를 받으면 이전 코드는 쓸 수 없습니다.</p><p><a href="/">다른 주소로 시작</a></p></section></html>{{end}}
{{define "connection"}}{{template "head" .}}<section><p>대상 agent: <code>{{.Agent}}</code></p><p>클라이언트: {{.Client}} / 연결 방식: {{.Mode}}</p><p>개인키는 자기 클라이언트에만 보관합니다. relay에 보내지 마세요.</p><p>권한: 이 agent의 송신·수신만 허용합니다. owner 승인·운영 관리자 권한은 없습니다.</p><p>기한: {{.Exp}} (최대 10분·1회 사용)</p>{{if .Token}}<h2>상태: 연결 대기</h2><p>자기 클라이언트 컴퓨터에서 저장소의 adapters를 설치한 뒤 실행하세요.</p><code>node adapters/dist/connect.js prepare &lt;서비스 URL&gt; &lt;새 비공개 폴더&gt;</code><p>아래 수단을 클라이언트의 비공개 표준입력에 넣으세요. 대화·로그·캡처에 넣지 마세요. 발급만으로 연결되지 않습니다.</p><label>연결 수단<input readonly value="{{.Token}}"></label><p><a href="/home/connections/{{.ConnectionID}}">공개키 준비 상태 확인</a></p>{{else}}<p>상태: {{.State}}</p>{{if eq .State "prepared"}}<h2>공개키 확인</h2><p>지문: <code>{{.Fingerprint}}</code></p><p>클라이언트에 표시된 지문과 같은지 확인하세요. 다르면 취소하세요.</p><form method="post" action="/home/confirm"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>대상·지문 확인 후 연결 승인</button></form>{{else if eq .State "approved"}}<p>승인 대기 완료. 아직 새 키가 연결되지 않았습니다. 자기 클라이언트에서 다음 명령을 실행하세요.</p><code>node adapters/dist/connect.js complete &lt;새 비공개 폴더&gt;</code>{{else if eq .State "consumed"}}<p>연결 완료. 클라이언트가 수신한 자격을 자기 환경에 보관합니다.</p>{{else if eq .State "expired"}}<p>만료. 새 키가 연결되지 않았습니다. 홈에서 다시 시작하세요.</p>{{else if eq .State "cancelled"}}<p>취소. 새 키가 연결되지 않았습니다.</p>{{else}}<p>공개키 준비 전입니다. 아직 연결되지 않았습니다.</p>{{end}}{{end}}{{if ne .State "consumed"}}<form method="post" action="/home/cancel"><input type="hidden" name="connection" value="{{.ConnectionID}}"><button>연결 취소</button></form>{{end}}<a href="/home">자기 홈으로 돌아가기</a></section></html>{{end}}
{{define "home"}}{{template "head" .}}<section><h2>내 신원</h2><dl><dt>확인된 이메일</dt><dd>{{.Masked}}</dd><dt>회원 식별자</dt><dd><code>{{.Member}}</code></dd><dt>세션</dt><dd>최대 {{.Absolute}}까지 유지됩니다. 60분 동안 활동이 없으면 먼저 끝납니다.</dd></dl></section><section><h2>내 agent</h2><p>활성 agent는 회원당 5개입니다. 연결은 Node 22가 있는 자기 클라이언트 컴퓨터에서 실행합니다. 다닷·외부 앱 실제 연결은 후속 검증입니다.</p><form method="post" action="/home/agents"><button>새 agent 만들기</button></form>{{if not .AgentDetails}}<p>아직 연결한 agent가 없습니다.</p>{{end}}{{range .AgentDetails}}{{$agent := .ID}}<section><h3><code>{{.ID}}</code></h3><p>상태: {{.Status}}</p><label>복사용 agent 식별자<input readonly value="{{.ID}}"></label>{{range .Keys}}<p>키 <code>{{.Kid}}</code> / {{.Fingerprint}} / 상태: {{.Status}}</p>{{if eq .Status "활성"}}<form method="post" action="/home/key-revoke"><input type="hidden" name="agent" value="{{$agent}}"><input type="hidden" name="kid" value="{{.Kid}}"><button>선택한 키 철회</button></form>{{end}}{{end}}{{if ne .Status "철회"}}<form method="post" action="/home/connect"><input type="hidden" name="agent" value="{{.ID}}"><label>지원 클라이언트<select name="client"><option value="node-local">Node 22 로컬 클라이언트</option></select></label><label>연결 방식<select name="mode"><option value="register">새 키 등록 (활성 최대 3개)</option><option value="rotate">회전 (완료 시 기존 키 전체 철회)</option></select></label><button>연결 수단 발급</button></form><form method="post" action="/home/agent-revoke"><input type="hidden" name="agent" value="{{.ID}}"><button>이 agent와 모든 키·관계 철회</button></form>{{end}}</section>{{end}}<p>키·연결 권한 변경에는 5분 안의 이메일 재확인이 필요합니다. 개인키는 클라이언트에만 보관합니다.</p><form method="post" action="/auth/reauth"><button>이메일 다시 확인</button></form></section><section><h2>관계</h2><p>가입 초대가 아닙니다. 상대 agent 식별자로 초대하고, 수신 owner가 명시적으로 수락해야 메시지를 허용합니다. 같은 owner의 두 agent도 수락해야 합니다.</p><form method="post" action="/home/invites"><label for="from-agent">내 발신 agent 식별자</label><input id="from-agent" name="agent" required maxlength="128"><label for="target-agent">상대 agent 식별자</label><input id="target-agent" name="target" required maxlength="128"><button>관계 초대</button></form>{{if .Pairs}}{{range .Pairs}}<section><p>발신 <code>{{.Inviter}}</code> / 수신 <code>{{.Recipient}}</code></p><p>상태: {{.State}} / 관계 세대: {{.Generation}} / 초대 기한: {{.Exp}}</p>{{if and (eq .State "pending") .Incoming}}<form method="post" action="/home/invite-decision"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><button name="decision" value="accept">수신 owner로 수락</button><button name="decision" value="deny">수신 owner로 거절</button></form>{{end}}<form method="post" action="/home/unpair"><input type="hidden" name="agent" value="{{.A}}"><input type="hidden" name="target" value="{{.B}}"><input type="hidden" name="generation" value="{{.Generation}}"><button>관계 철회</button></form></section>{{end}}{{else}}<p>관계가 없습니다.</p>{{end}}<p>수락 전·철회 후 메시지는 거부됩니다. pending은 24시간 후 만료됩니다.</p></section><section><h2>승인 요청</h2>{{if .Gates}}<ul>{{range .Gates}}<li><a href="/home/gates/{{.ID}}">{{.Parent}}</a> 상태: {{.State}}</li>{{end}}</ul>{{else}}<p>승인 요청이 없습니다.</p>{{end}}</section><section><h2>로그아웃</h2><p>로그아웃과 전체 로그아웃은 브라우저 세션만 끝냅니다. 별도로 연결한 agent의 키와 자격은 철회하지 않습니다.</p><form method="post" action="/auth/logout"><button>이 브라우저에서 로그아웃</button></form>{{if .Recent}}<form method="post" action="/auth/logout-all"><button>모든 브라우저에서 로그아웃</button></form>{{else}}<p>모든 브라우저에서 로그아웃하려면 5분 안에 이메일을 다시 확인해야 합니다.</p><form method="post" action="/auth/reauth"><button>이메일 다시 확인</button></form>{{end}}<h3>계정 비활성화</h3><p>상태: 준비 중. 아직 요청할 수 없습니다.</p></section></html>{{end}}
`))

var notices = map[string]string{
	"logout":  "로그아웃했습니다. 연결한 agent의 키와 자격은 그대로입니다.",
	"all":     "모든 브라우저에서 로그아웃했습니다. 연결한 agent의 키와 자격은 그대로입니다.",
	"expired": "세션이 끝났습니다. 다시 로그인하세요.",
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
	return map[string]any{"Title": "KnowsLink 시작", "Problem": problem}
}
func limited(retry time.Time) string {
	return "요청이 많습니다. " + clock(retry) + " 이후 다시 시도하세요."
}

func (s *Service) memberRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", s.startPage)
	mux.HandleFunc("POST /auth/start", s.authStart)
	mux.HandleFunc("GET /auth/verify", s.verifyPage)
	mux.HandleFunc("POST /auth/verify", s.authVerify)
	mux.HandleFunc("POST /auth/reauth", s.authReauth)
	mux.HandleFunc("POST /auth/logout", s.logout(false))
	mux.HandleFunc("POST /auth/logout-all", s.logout(true))
	mux.HandleFunc("GET /home", s.homePage)
	mux.HandleFunc("GET /home/gates/{id}", s.gatePage(sessionOwner))
	mux.HandleFunc("POST /home/gates/{id}", s.gateDecision(sessionOwner))
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
		if ok, retry := st.hit(now, anonymousRate(ip)...); !ok {
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
		if ok, retry := st.hit(now, anonymousRate(ip)...); !ok {
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
	_, retry := st.hit(now, buckets...)
	return id, session, retry, err
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, session, retry, err := s.memberHit(st, r, now, memberRate)
		if !retry.IsZero() {
			return map[string]any{"Title": "내 KnowsLink", "Problem": limited(retry), "limited": true}, nil
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
				pairs = append(pairs, memberPair{p, st.Agents[p.Recipient] != nil && st.Agents[p.Recipient].Owner == m.Owner})
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
	v := value.(map[string]any)
	if v["limited"] == true {
		render(w, 429, "head", v)
		return
	}
	render(w, 200, "home", v)
}

func (s *Service) authReauth(w http.ResponseWriter, r *http.Request) {
	ip := s.clientIP(r)
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
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
		render(w, 429, "head", map[string]any{"Title": "내 KnowsLink", "Problem": limited(result.Retry)})
		return
	}
	s.deliverCode(w, r, result)
}

type refusal struct {
	status  int
	problem string
}

// logout uses the separate cleanup budget so saturated new-work limits never block ending a session.
func (s *Service) logout(all bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := readCookie(r, sessionCookie)
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			id, session, retry, err := s.memberHit(st, r, now, cleanupRate)
			if !retry.IsZero() {
				return refusal{429, limited(retry)}, nil
			}
			if err != nil {
				return nil, err
			}
			if all && now.Sub(session.Verified) >= reauthWindow {
				return refusal{403, "모든 브라우저에서 로그아웃하려면 5분 안에 이메일을 다시 확인해야 합니다."}, nil
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
		if refused := value.(refusal); refused.status != 0 {
			render(w, refused.status, "head", map[string]any{"Title": "내 KnowsLink", "Problem": refused.problem})
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
