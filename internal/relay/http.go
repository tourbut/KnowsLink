// HTTP API separates owner credentials from agent credentials and serves escaped approval forms.
package relay

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

type command struct {
	Agent    string `json:"agent,omitempty"`
	Kid      string `json:"kid,omitempty"`
	Public   string `json:"public,omitempty"`
	Proof    string `json:"proof,omitempty"`
	Target   string `json:"target,omitempty"`
	Decision string `json:"decision,omitempty"`
	ID       string `json:"id,omitempty"`
	Token    string `json:"token,omitempty"`
	Claim    string `json:"claim,omitempty"`
}

func bearer(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
func ownerToken(r *http.Request) string {
	if _, password, ok := r.BasicAuth(); ok {
		return password
	}
	return bearer(r)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func statusFor(err error) int {
	switch err.Error() {
	case "invalid_auth":
		return 401
	case "sender_not_allowed", "human_invite_required", "disclosure_denied":
		return 403
	case "idempotency_conflict", "id_collision", "duplicate_gate", "duplicate_result", "already_claimed", "invalid_lease", "invalid_gate", "key_exists":
		return 409
	case "unavailable", "abnormal_clock":
		return 503
	default:
		return 422
	}
}
func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/registry", func(w http.ResponseWriter, r *http.Request) {
		if !RegistryValid() {
			respond(w, nil, fault("unavailable"))
			return
		}
		writeJSON(w, 200, map[string]any{"sha256": RegistrySHA256, "manifest": string(registry)})
	})
	mux.HandleFunc("POST /v1/owners", func(w http.ResponseWriter, r *http.Request) {
		// Local synthetic signup is not identity verification; public deployments leave it off so it cannot mint members or owners.
		if !s.SyntheticSignup {
			respond(w, nil, fault("sender_not_allowed"))
			return
		}
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			token := randomToken()
			id := "owner_" + randomToken()
			st.Owners[id] = &Owner{hashToken(token), true}
			return map[string]string{"owner": id, "credential": token}, nil
		})
		respond(w, result, err)
	})
	for _, path := range []string{"agents", "keys", "key-revoke", "invites", "invite-decision", "unpair", "owner-revoke", "pull", "persist", "ack", "claim", "gate-consume", "authorize"} {
		mux.HandleFunc("POST /v1/"+path, s.operation(path))
	}
	mux.HandleFunc("POST /v1/send", func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64*1024))
		if err != nil {
			respond(w, nil, fault("invalid_json"))
			return
		}
		e, err := Parse(raw)
		if err != nil {
			respond(w, nil, err)
			return
		}
		result, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			return st.ingest(agent, e, r.Header.Get("X-Execution-Claim"), now)
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/gates/{id}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			g := st.Gates[r.PathValue("id")]
			if g == nil || g.To != agent {
				return nil, fault("sender_not_allowed")
			}
			return map[string]any{"state": g.State, "consumed": g.Consumed, "parent": g.Parent}, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			contacts := []string{}
			for _, p := range st.Pairs {
				if p.State == "active" {
					if p.A == agent {
						contacts = append(contacts, p.B)
					}
					if p.B == agent {
						contacts = append(contacts, p.A)
					}
				}
			}
			return contacts, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/keys/{agent}/{kid}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			target := r.PathValue("agent")
			p := st.Pairs[pairID(agent, target)]
			if agent != target && (p == nil || p.State != "active") {
				return nil, fault("sender_not_allowed")
			}
			a := st.Agents[target]
			if a == nil {
				return nil, fault("invalid_auth")
			}
			k := a.Keys[r.PathValue("kid")]
			if k == nil || k.Revoked {
				return nil, fault("invalid_auth")
			}
			return map[string]string{"public": base64.RawURLEncoding.EncodeToString(k.Public)}, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/receipts/{id}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			m := st.Messages[r.PathValue("id")]
			if m == nil || (m.Receipt.From != agent && m.Receipt.To != agent) || !st.current(m) {
				return nil, fault("sender_not_allowed")
			}
			return map[string]any{"receipt": m.Receipt, "completion": m.Completion}, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /owner", s.ownerPage)
	mux.HandleFunc("GET /owner/gates/{id}", s.gatePage(basicOwner))
	mux.HandleFunc("POST /owner/gates/{id}", s.gateDecision(basicOwner))
	s.memberRoutes(mux)
	// Stdlib Sec-Fetch-Site/Origin check rejects cross-origin browser writes; non-browser agent calls carry neither header.
	protected := http.NewCrossOriginProtection().Handler(mux)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if strings.HasPrefix(r.URL.Path, "/v1/test/") {
			s.testHandler(mux).ServeHTTP(w, r)
			return
		}
		protected.ServeHTTP(w, r)
	})
}
func respond(w http.ResponseWriter, result any, err error) {
	if err != nil {
		writeJSON(w, statusFor(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, result)
}
func (s *Service) operation(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c command
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64*1024))
		if err != nil {
			respond(w, nil, fault("invalid_json"))
			return
		}
		if len(raw) == 0 {
			raw = []byte("{}")
		}
		if err := Strict(raw, &c); err != nil {
			respond(w, nil, err)
			return
		}
		result, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) { return st.operate(path, bearer(r), c, now) })
		respond(w, result, err)
	}
}
func (st *State) operate(path, token string, c command, now time.Time) (any, error) {
	testOnly := strings.HasPrefix(path, "test-")
	if testOnly {
		path = strings.TrimPrefix(path, "test-")
	}
	kind := "agent"
	switch path {
	case "agents", "keys", "key-revoke", "invite-decision", "unpair", "owner-revoke":
		kind = "owner"
	}
	principal, err := st.principal(token, kind)
	if err != nil && path == "invites" {
		principal, err = st.principal(token, "owner")
		kind = "owner"
	}
	if err != nil {
		return nil, err
	}
	if testOnly {
		if !st.TestAgents[principal] {
			return nil, fault("sender_not_allowed")
		}
		if path != "pull" {
			m := st.Messages[c.ID]
			if m == nil || m.Receipt.Intent != "relay.test.message" {
				return nil, fault("sender_not_allowed")
			}
		}
	}
	switch path {
	case "agents", "keys":
		if !agentPattern.MatchString(c.Agent) || !kidPattern.MatchString(c.Kid) {
			return nil, fault("invalid_schema")
		}
		a := st.Agents[c.Agent]
		if path == "agents" && a != nil {
			return nil, fault("key_exists")
		}
		if path == "keys" && (a == nil || a.Owner != principal) {
			return nil, fault("sender_not_allowed")
		}
		public, e := base64.RawURLEncoding.DecodeString(c.Public)
		proof, e2 := base64.RawURLEncoding.DecodeString(c.Proof)
		if e != nil || e2 != nil || len(public) != 32 || !ed25519.Verify(public, popBytes(principal, c.Agent, c.Kid, c.Public), proof) {
			return nil, fault("invalid_signature")
		}
		credential := ""
		if a == nil {
			credential = randomToken()
			a = &Agent{principal, hashToken(credential), map[string]*Key{}}
			st.Agents[c.Agent] = a
		}
		if a.Keys[c.Kid] != nil {
			return nil, fault("key_exists")
		}
		if path == "keys" {
			for _, k := range a.Keys {
				k.Revoked = true
				k.Changed = now
			}
		}
		a.Keys[c.Kid] = &Key{public, false, now}
		st.sweep(now)
		return map[string]string{"agent": c.Agent, "credential": credential}, nil
	case "key-revoke":
		a := st.Agents[c.Agent]
		if a == nil || a.Owner != principal || a.Keys[c.Kid] == nil {
			return nil, fault("sender_not_allowed")
		}
		a.Keys[c.Kid].Revoked = true
		a.Keys[c.Kid].Changed = now
		st.sweep(now)
		return map[string]string{"state": "revoked"}, nil
	case "owner-revoke":
		st.Owners[principal].Active = false
		st.sweep(now)
		return map[string]string{"state": "revoked"}, nil
	case "invites":
		a, b := st.Agents[c.Agent], st.Agents[c.Target]
		if a == nil || b == nil || c.Agent == c.Target || (kind == "agent" && principal != c.Agent) || (kind == "owner" && a.Owner != principal) || !st.Owners[b.Owner].Active {
			return nil, fault("sender_not_allowed")
		}
		id := pairID(c.Agent, c.Target)
		p := st.Pairs[id]
		if p != nil && (p.State == "active" || p.State == "pending") {
			return p, nil
		}
		generation := int64(1)
		if p != nil {
			generation = p.Generation + 1
		}
		p = &Pair{c.Agent, c.Target, c.Agent, c.Target, "pending", generation}
		st.Pairs[id] = p
		return p, nil
	case "invite-decision":
		p := st.Pairs[pairID(c.Agent, c.Target)]
		if p == nil || st.Agents[p.Recipient].Owner != principal || (c.Decision != "accept" && c.Decision != "deny") {
			return nil, fault("sender_not_allowed")
		}
		if p.State == "active" && c.Decision == "accept" {
			return p, nil
		}
		if p.State != "pending" {
			return nil, fault("sender_not_allowed")
		}
		p.State = "denied"
		if c.Decision == "accept" {
			p.State = "active"
		}
		return p, nil
	case "unpair":
		p := st.Pairs[pairID(c.Agent, c.Target)]
		if p == nil || (st.Agents[p.A].Owner != principal && st.Agents[p.B].Owner != principal) {
			return nil, fault("sender_not_allowed")
		}
		p.State = "revoked"
		st.sweep(now)
		return p, nil
	case "pull":
		return st.leaseMessageFor(principal, now, testOnly)
	case "persist", "ack":
		m, err := st.leased(principal, c.ID, c.Token, now)
		if err != nil {
			return nil, err
		}
		if path == "persist" {
			m.Persisted = true
			m.Inbox = append(json.RawMessage(nil), m.Envelope...)
			return map[string]string{"state": "persisted"}, nil
		}
		if !m.Persisted {
			return nil, fault("invalid_lease")
		}
		m.Receipt.State = "delivered"
		m.LeaseToken = ""
		if m.Receipt.Intent == "relay.result" {
			m.Envelope = nil
			m.Inbox = nil
		}
		return m.Receipt, nil
	case "claim":
		m := st.Messages[c.ID]
		if m == nil || m.Receipt.To != principal || m.Receipt.State != "delivered" || !m.Persisted || len(m.Inbox) == 0 || !st.current(m) || !now.Before(m.Receipt.Exp) || m.Receipt.Intent == "relay.result" || m.Deliver != "agent" {
			return nil, fault("sender_not_allowed")
		}
		if m.Claimed {
			return nil, fault("already_claimed")
		}
		m.Claimed = true
		m.ClaimToken = randomToken()
		if m.Receipt.Intent == "relay.test.message" {
			m.Envelope = nil
			m.Inbox = nil
		}
		return map[string]string{"claim": m.ClaimToken, "policy": policy(m.Receipt.Intent)}, nil
	case "gate-consume":
		g := st.Gates[c.ID]
		if g == nil || g.To != principal || g.State != "approved" || g.Consumed || !now.Before(g.Exp) {
			return nil, fault("invalid_gate")
		}
		m, err := st.parentFor(principal, g.Parent, now)
		if err != nil {
			return nil, err
		}
		if m.ClaimToken != c.Claim || c.Claim == "" || m.Receipt.Digest != g.Digest || m.Generation != g.Generation {
			return nil, fault("invalid_gate")
		}
		g.Consumed = true
		return map[string]any{"gate_passed": true, "executable": false, "disclosure": false}, nil
	case "authorize":
		m, err := st.parentFor(principal, c.ID, now)
		if err != nil {
			return nil, err
		}
		if c.Claim == "" || m.ClaimToken != c.Claim {
			return nil, fault("sender_not_allowed")
		}
		return map[string]any{"executable": false, "disclosure": false, "policy": policy(m.Receipt.Intent)}, nil
	}
	return nil, fault("invalid_schema")
}
func csrf(token, id string) string {
	mac := hmac.New(sha256.New, []byte(token))
	_, _ = mac.Write([]byte("gate:" + id))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

var page = template.Must(template.New("gate").Parse(`<!doctype html><html lang="ko"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>KnowsLink 승인</title><style>body{font:18px system-ui;max-width:760px;margin:40px auto;padding:20px;color:#17232b;background:#f5f7f8}section{background:white;padding:24px;border:1px solid #bbc5cd}pre{white-space:pre-wrap;overflow-wrap:anywhere}button{font:inherit;padding:12px;margin:8px}dt{font-weight:bold}a{color:#164da0}</style><h1>KnowsLink 요청 승인</h1><section><h2>상태: {{.State}}</h2><dl><dt>발신 에이전트</dt><dd>{{.From}}</dd><dt>대상 에이전트</dt><dd>{{.To}}</dd><dt>원요청 ID</dt><dd>{{.Parent}}</dd><dt>Intent</dt><dd>{{.Intent}}</dd><dt>만료</dt><dd>{{.Exp}}</dd><dt>적용 정책</dt><dd>{{.Policy}}</dd></dl><h2>검증된 typed body</h2><pre>{{.Body}}</pre><p>Approve는 인간 게이트만 통과합니다. 정보 공개나 일정 실행을 허용하지 않습니다.</p>{{if .Active}}<form method="post"><input type="hidden" name="csrf" value="{{.CSRF}}"><button name="decision" value="approve">Approve 승인</button><button name="decision" value="deny">Deny 거절</button></form>{{else}}<p>승인·거절 버튼 비활성: {{.State}}</p>{{end}}<a href="{{.Back}}">작업 화면으로 돌아가기</a></section></html>`))

// ownerAuth returns the CSRF key, the owner resolver, and the page to return to.
type ownerAuth func(r *http.Request) (string, func(*State, time.Time) (string, error), string)

func basicOwner(r *http.Request) (string, func(*State, time.Time) (string, error), string) {
	token := ownerToken(r)
	return token, func(st *State, _ time.Time) (string, error) { return st.principal(token, "owner") }, "/owner"
}
func sessionOwner(r *http.Request) (string, func(*State, time.Time) (string, error), string) {
	token := readCookie(r, sessionCookie)
	return token, func(st *State, now time.Time) (string, error) {
		id, _, err := st.session(token, now)
		if err != nil {
			return "", err
		}
		return st.Members[id].Owner, nil
	}, "/home"
}

func (s *Service) gateView(r *http.Request, auth ownerAuth) (any, error) {
	token, resolve, back := auth(r)
	return s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		owner, err := resolve(st, now)
		if err != nil {
			return nil, err
		}
		g := st.Gates[r.PathValue("id")]
		if g == nil || g.Owner != owner {
			return nil, fault("sender_not_allowed")
		}
		view := map[string]any{"State": g.State, "From": g.From, "To": g.To, "Parent": g.Parent, "Exp": g.Exp, "Policy": g.Policy, "CSRF": csrf(token, g.ID), "Active": false, "Body": "원문 부재: 승인 불가", "Back": back}
		m := st.Messages[g.Parent]
		if m != nil && len(m.Envelope) > 0 {
			e, err := Parse(m.Envelope)
			if err != nil {
				return nil, fault("unavailable")
			}
			a := st.Agents[e.From]
			if a == nil || a.Keys[e.Sig.Kid] == nil || !e.Verify(a.Keys[e.Sig.Kid].Public) {
				return nil, fault("unavailable")
			}
			body, _ := json.MarshalIndent(e.Body, "", "  ")
			view["Body"] = string(body)
			view["Intent"] = e.Intent
			view["Active"] = g.State == "pending" && st.current(m) && now.Before(g.Exp)
		}
		return view, nil
	})
}
func ownerError(w http.ResponseWriter, err error) {
	if err.Error() == "invalid_auth" {
		w.Header().Set("WWW-Authenticate", `Basic realm="KnowsLink owner", charset="UTF-8"`)
	}
	http.Error(w, err.Error(), statusFor(err))
}

// gateError sends members back to email login instead of the owner Basic prompt.
func gateError(w http.ResponseWriter, r *http.Request, back string, err error) {
	if back == "/home" && err.Error() == "invalid_auth" {
		expired(w, r)
		return
	}
	ownerError(w, err)
}
func (s *Service) gatePage(auth ownerAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := s.gateView(r, auth)
		if err != nil {
			_, _, back := auth(r)
			gateError(w, r, back, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = page.Execute(w, view)
	}
}
func (s *Service) gateDecision(auth ownerAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, resolve, back := auth(r)
		r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid_schema", 422)
			return
		}
		if !hmac.Equal([]byte(r.PostForm.Get("csrf")), []byte(csrf(token, r.PathValue("id")))) {
			http.Error(w, "invalid_csrf", 403)
			return
		}
		_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			owner, err := resolve(st, now)
			if err != nil {
				return nil, err
			}
			g := st.Gates[r.PathValue("id")]
			if g == nil || g.Owner != owner || g.State != "pending" || !now.Before(g.Exp) {
				return nil, fault("invalid_gate")
			}
			m := st.Messages[g.Parent]
			h := st.Messages[g.ID]
			if m == nil || h == nil || len(m.Envelope) == 0 || !st.current(m) || !st.current(h) || m.Generation != g.Generation || m.Receipt.Digest != g.Digest {
				return nil, fault("invalid_gate")
			}
			decision := r.PostForm.Get("decision")
			if decision != "approve" && decision != "deny" {
				return nil, fault("invalid_schema")
			}
			h.Receipt.State = "delivered"
			h.Persisted = true
			h.Envelope = nil
			h.Inbox = nil
			g.State = "denied"
			if decision == "approve" {
				g.State = "approved"
			}
			return nil, nil
		})
		if err != nil {
			gateError(w, r, back, err)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
}

var dashboard = template.Must(template.New("owner").Parse(`<!doctype html><html lang="ko"><meta charset="utf-8"><title>KnowsLink Owner</title><h1>KnowsLink Owner 작업 화면</h1><p>가입·키·초대 관리는 인증된 /v1 API를 사용합니다. 원문은 요청 만료 또는 완료 후 지워집니다.</p><h2>등록 에이전트</h2><ul>{{range .Agents}}<li>{{.}}</li>{{end}}</ul><h2>승인 요청</h2><ul>{{range .Gates}}<li><a href="/owner/gates/{{.ID}}">{{.Parent}}</a> 상태: {{.State}}</li>{{end}}</ul></html>`))

func (s *Service) ownerPage(w http.ResponseWriter, r *http.Request) {
	view, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
		owner, err := st.principal(ownerToken(r), "owner")
		if err != nil {
			return nil, err
		}
		agents := []string{}
		gates := []*Gate{}
		for id, a := range st.Agents {
			if a.Owner == owner {
				agents = append(agents, id)
			}
		}
		for _, g := range st.Gates {
			if g.Owner == owner {
				gates = append(gates, g)
			}
		}
		return map[string]any{"Agents": agents, "Gates": gates}, nil
	})
	if err != nil {
		ownerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = dashboard.Execute(w, view)
}
