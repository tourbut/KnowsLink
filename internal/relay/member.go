// Member pages run verified login, sessions, logout, and the self-owner home without owner Basic auth or SSH.
package relay

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var kst = time.FixedZone("KST", 9*60*60)

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
	mux.HandleFunc("GET /home", s.humanPage("home"))
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

func (s *Service) legacyHomePage(w http.ResponseWriter, r *http.Request) {
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
	render(w, 200, "legacy-home", value.(map[string]any))
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
