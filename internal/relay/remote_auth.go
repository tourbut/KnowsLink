// Remote OAuth binds each explicitly approved host installation to its own Google member agent and key.
package relay

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const remoteScope = "knowslink.text"
const remoteLifetime = 30 * 24 * time.Hour

type remoteClient struct {
	Name      string
	Redirects []string
	Exp       time.Time
}
type remoteAttempt struct {
	Client, Redirect, Challenge, State, Resource, Owner string
	Exp                                                 time.Time
}
type remoteCode struct {
	Install, Client, Redirect, Challenge, Resource string
	Exp                                            time.Time
}
type remoteToken struct {
	Install, Client, Resource, Kind string
	Exp                             time.Time
}
type remoteInstall struct {
	Owner, Agent, Kid string
	Private           []byte
	Exp               time.Time
	Inbox             map[string]json.RawMessage
}
type remoteState struct {
	Clients       map[string]*remoteClient
	Attempts      map[string]*remoteAttempt
	Codes         map[string]*remoteCode
	Tokens        map[string]*remoteToken
	Installs      map[string]*remoteInstall
	Subscriptions map[string]*remoteSubscription
}

func (st *State) remote() *remoteState {
	if st.Remote == nil {
		st.Remote = &remoteState{map[string]*remoteClient{}, map[string]*remoteAttempt{}, map[string]*remoteCode{}, map[string]*remoteToken{}, map[string]*remoteInstall{}, map[string]*remoteSubscription{}}
	}
	return st.Remote
}

func (st *State) remoteInstall(id string, now time.Time) (*remoteInstall, error) {
	if st.Remote == nil {
		return nil, fault("invalid_auth")
	}
	i := st.Remote.Installs[id]
	if i == nil || !now.Before(i.Exp) || !st.publicOwner(i.Owner) {
		return nil, fault("invalid_auth")
	}
	a := st.Agents[i.Agent]
	if a == nil || a.Revoked || a.Owner != i.Owner {
		return nil, fault("invalid_auth")
	}
	k := a.Keys[i.Kid]
	if k == nil || k.Revoked || len(i.Private) != ed25519.PrivateKeySize || !hmac.Equal(k.Public, ed25519.PrivateKey(i.Private).Public().(ed25519.PublicKey)) {
		return nil, fault("invalid_auth")
	}
	return i, nil
}

func (st *State) remotePrincipal(token, resource string, now time.Time) (string, *remoteInstall, error) {
	if st.Remote == nil || len(token) != 43 {
		return "", nil, fault("invalid_auth")
	}
	t := st.Remote.Tokens[hashToken(token)]
	if t == nil || t.Kind != "access" || t.Resource != resource || !now.Before(t.Exp) {
		return "", nil, fault("invalid_auth")
	}
	i, err := st.remoteInstall(t.Install, now)
	return t.Install, i, err
}

func (st *State) sweepRemote(now time.Time) {
	if st.Remote == nil {
		return
	}
	r := st.Remote
	for id, c := range r.Clients {
		if !now.Before(c.Exp) {
			delete(r.Clients, id)
		}
	}
	for id, a := range r.Attempts {
		if !now.Before(a.Exp) {
			delete(r.Attempts, id)
		}
	}
	for id, c := range r.Codes {
		if !now.Before(c.Exp) {
			delete(r.Codes, id)
		}
	}
	for id, t := range r.Tokens {
		_, err := st.remoteInstall(t.Install, now)
		if !now.Before(t.Exp) || err != nil {
			delete(r.Tokens, id)
		}
	}
	for id, sub := range r.Subscriptions {
		if !now.Before(sub.Rotated.Add(time.Minute)) {
			sub.OldSecret = ""
		}
		_, err := st.remoteInstall(sub.Install, now)
		if !now.Before(sub.Exp) || err != nil {
			delete(r.Subscriptions, id)
		}
	}
	for id, i := range r.Installs {
		for mid := range i.Inbox {
			m := st.Messages[mid]
			if m == nil || !st.current(m) || !now.Before(m.Receipt.Exp) {
				delete(i.Inbox, mid)
			}
		}
		if _, err := st.remoteInstall(id, now); err != nil {
			if a := st.Agents[i.Agent]; a != nil {
				a.Revoked = true
				a.Changed = now
				revokeKeys(a, now)
			}
			delete(r.Installs, id)
		}
	}
}

func (s *Service) remoteOrigin() string {
	if s.Google == nil {
		return ""
	}
	u, err := url.Parse(s.Google.config.RedirectURL)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return ""
	}
	return "https://" + u.Host
}

// Remote tokens are recognized only by /mcp admission, never by the local credential API.
func (st *State) remoteAdmission(token string, now time.Time) string {
	if st.Remote == nil {
		return ""
	}
	t := st.Remote.Tokens[hashToken(token)]
	if t == nil || t.Kind != "access" || !now.Before(t.Exp) {
		return ""
	}
	i, e := st.remoteInstall(t.Install, now)
	if e != nil {
		return ""
	}
	return i.Agent
}
func (st *State) remoteCredentials(live map[string]honored) {
	if st.Remote == nil {
		return
	}
	for id, t := range st.Remote.Tokens {
		if t.Kind == "access" {
			if i, e := st.remoteInstall(t.Install, time.Time{}); e == nil {
				live[id] = honored{kind: "remote", id: i.Agent, exp: t.Exp}
			}
		}
	}
}

func remoteURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && len(raw) <= 2048 && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && u.Fragment == "" && (u.Port() == "" || u.Port() == "443") && !strings.ContainsAny(u.Host, "\\\r\n")
}

func (s *Service) remoteRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/oauth-protected-resource", s.remoteMetadata(false))
	mux.HandleFunc("GET /.well-known/oauth-protected-resource/mcp", s.remoteMetadata(false))
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", s.remoteMetadata(true))
	mux.HandleFunc("POST /remote/oauth/register", s.remoteRegister)
	mux.HandleFunc("GET /remote/oauth/authorize", s.remoteAuthorize)
	mux.HandleFunc("POST /remote/oauth/token", s.remoteToken)
	mux.HandleFunc("POST /remote/oauth/revoke", s.remoteRevoke)
	mux.HandleFunc("GET /remote/consent/{id}", s.remoteConsent)
	mux.HandleFunc("POST /remote/consent/{id}", s.remoteConsent)
	mux.HandleFunc("POST /mcp", s.remoteMCP)
	// Stateless Streamable HTTP has no server-to-client SSE session or persistent DELETE session.
	mux.HandleFunc("GET /mcp", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusMethodNotAllowed) })
}

func (s *Service) remoteMetadata(auth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := s.remoteOrigin()
		if o == "" {
			writeJSON(w, 503, map[string]string{"error": "temporarily_unavailable"})
			return
		}
		v := map[string]any{"resource": o + "/mcp", "authorization_servers": []string{o}, "scopes_supported": []string{remoteScope}, "bearer_methods_supported": []string{"header"}}
		if auth {
			v = map[string]any{"issuer": o, "authorization_response_iss_parameter_supported": true, "authorization_endpoint": o + "/remote/oauth/authorize", "token_endpoint": o + "/remote/oauth/token", "registration_endpoint": o + "/remote/oauth/register", "revocation_endpoint": o + "/remote/oauth/revoke", "token_endpoint_auth_methods_supported": []string{"none"}, "revocation_endpoint_auth_methods_supported": []string{"none"}, "response_types_supported": []string{"code"}, "grant_types_supported": []string{"authorization_code", "refresh_token"}, "code_challenge_methods_supported": []string{"S256"}, "scopes_supported": []string{remoteScope}}
		}
		writeJSON(w, 200, v)
	}
}

func (s *Service) remoteRegister(w http.ResponseWriter, r *http.Request) {
	if s.remoteOrigin() == "" {
		writeJSON(w, 503, map[string]string{"error": "temporarily_unavailable"})
		return
	}
	var in struct {
		Name            string   `json:"client_name"`
		Redirects       []string `json:"redirect_uris"`
		Method          string   `json:"token_endpoint_auth_method"`
		Grants          []string `json:"grant_types"`
		Responses       []string `json:"response_types"`
		Scope           string   `json:"scope"`
		ClientURI       string   `json:"client_uri,omitempty"`
		LogoURI         string   `json:"logo_uri,omitempty"`
		Contacts        []string `json:"contacts,omitempty"`
		SoftwareID      string   `json:"software_id,omitempty"`
		SoftwareVersion string   `json:"software_version,omitempty"`
	}
	if decodeCommand(w, r, &in) != nil || len(in.Name) > 128 || len(in.Redirects) == 0 || len(in.Redirects) > 5 || (in.Method != "" && in.Method != "none") || (in.Scope != "" && in.Scope != remoteScope) {
		writeJSON(w, 400, map[string]string{"error": "invalid_client_metadata"})
		return
	}
	for _, u := range in.Redirects {
		if !remoteURL(u) {
			writeJSON(w, 400, map[string]string{"error": "invalid_redirect_uri"})
			return
		}
	}
	for _, g := range in.Grants {
		if g != "authorization_code" && g != "refresh_token" {
			writeJSON(w, 400, map[string]string{"error": "invalid_client_metadata"})
			return
		}
	}
	for _, v := range in.Responses {
		if v != "code" {
			writeJSON(w, 400, map[string]string{"error": "invalid_client_metadata"})
			return
		}
	}
	id := randomToken()
	_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		rs := st.remote()
		if len(rs.Clients) >= 1000 {
			return nil, fault("capacity")
		}
		rs.Clients[id] = &remoteClient{in.Name, in.Redirects, now.Add(remoteLifetime)}
		return nil, nil
	})
	if err != nil {
		writeJSON(w, statusFor(err), map[string]string{"error": "temporarily_unavailable"})
		return
	}
	writeJSON(w, 201, map[string]any{"client_id": id, "client_name": in.Name, "redirect_uris": in.Redirects, "token_endpoint_auth_method": "none", "grant_types": []string{"authorization_code", "refresh_token"}, "response_types": []string{"code"}, "scope": remoteScope})
}

func (s *Service) remoteAuthorize(w http.ResponseWriter, r *http.Request) {
	o := s.remoteOrigin()
	q := r.URL.Query()
	challenge, e := base64.RawURLEncoding.DecodeString(q.Get("code_challenge"))
	if o == "" || q.Get("response_type") != "code" || q.Get("code_challenge_method") != "S256" || e != nil || len(challenge) != 32 || q.Get("resource") != o+"/mcp" || q.Get("scope") != remoteScope || len(q.Get("state")) > 2048 {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	id := randomToken()
	_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		rs := st.remote()
		c := rs.Clients[q.Get("client_id")]
		if c == nil || !now.Before(c.Exp) {
			return nil, fault("invalid_auth")
		}
		allowed := false
		for _, u := range c.Redirects {
			if u == q.Get("redirect_uri") {
				allowed = true
			}
		}
		if !allowed {
			return nil, fault("invalid_auth")
		}
		if len(rs.Attempts) >= 1000 {
			return nil, fault("capacity")
		}
		rs.Attempts[id] = &remoteAttempt{Client: q.Get("client_id"), Redirect: q.Get("redirect_uri"), Challenge: q.Get("code_challenge"), State: q.Get("state"), Resource: o + "/mcp", Exp: now.Add(codeTTL)}
		return nil, nil
	})
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	http.Redirect(w, r, "/remote/consent/"+id, 303)
}

var remotePage = template.Must(template.New("remote").Parse(`<!doctype html><html lang="ko"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>KnowsLink 원격 연결</title>` + memberStyle + `<h1>원격 호스트 연결</h1><p>클라이언트: {{.Name}}</p><p>반환 주소: {{.Redirect}}</p><p>이 설치에 새 agent와 키를 만듭니다. 원격 서명용 개인 키와 구독 webhook 비밀값은 KnowsLink 서버의 비공개 DB에 보관합니다. 기존 로컬 개인 키는 가져오지 않습니다. 기존 agent와 관계는 유지됩니다. 연결한 상대에게 명시적으로 텍스트를 보내고 받을 수 있습니다. 이벤트는 도착 신호만 전송합니다. 자동 답장 권한은 부여하지 않습니다.</p>{{if .Login}}<form method="post" action="/auth/google"><input type="hidden" name="remote" value="{{.ID}}"><button>Google 계정 확인</button></form>{{else}}<form method="post"><input type="hidden" name="csrf" value="{{.CSRF}}"><button name="decision" value="approve">이 설치 연결 승인</button><button name="decision" value="deny">취소</button></form>{{end}}<a href="/home">내 agent 관리</a></html>`))

func (s *Service) remoteConsent(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if r.ParseForm() != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		rs := st.remote()
		a := rs.Attempts[r.PathValue("id")]
		if a == nil || !now.Before(a.Exp) {
			return nil, fault("invalid_auth")
		}
		c := rs.Clients[a.Client]
		if c == nil || !now.Before(c.Exp) {
			return nil, fault("invalid_auth")
		}
		member, session, e := st.session(readCookie(r, sessionCookie), now)
		login := e != nil || a.Owner == "" || st.Members[member].Owner != a.Owner || st.Members[member].Issuer != googleIssuer || now.Sub(session.Verified) >= reauthWindow
		if r.Method == "GET" {
			return map[string]any{"Name": c.Name, "Redirect": a.Redirect, "ID": r.PathValue("id"), "Login": login, "CSRF": csrf(readCookie(r, sessionCookie), r.PathValue("id"))}, nil
		}
		if login {
			return nil, fault("invalid_auth")
		}
		if !hmac.Equal([]byte(r.PostForm.Get("csrf")), []byte(csrf(readCookie(r, sessionCookie), r.PathValue("id")))) {
			return nil, fault("invalid_csrf")
		}
		u, _ := url.Parse(a.Redirect)
		q := u.Query()
		q.Set("state", a.State)
		q.Set("iss", s.remoteOrigin())
		if r.PostForm.Get("decision") == "deny" {
			delete(rs.Attempts, r.PathValue("id"))
			q.Set("error", "access_denied")
			u.RawQuery = q.Encode()
			return u.String(), nil
		}
		if r.PostForm.Get("decision") != "approve" {
			return nil, fault("invalid_schema")
		}
		code, e := st.approveRemote(a, now)
		if e != nil {
			return nil, e
		}
		delete(rs.Attempts, r.PathValue("id"))
		q.Set("code", code)
		u.RawQuery = q.Encode()
		return u.String(), nil
	})
	if err != nil {
		writeJSON(w, statusFor(err), map[string]string{"error": "authorization_failed"})
		return
	}
	if target, ok := v.(string); ok {
		http.Redirect(w, r, target, 303)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = remotePage.Execute(w, v)
}

func (st *State) approveRemote(a *remoteAttempt, now time.Time) (string, error) {
	if a.Owner == "" || !now.Before(a.Exp) || !st.publicOwner(a.Owner) {
		return "", fault("invalid_auth")
	}
	if !st.agentCapacity(a.Owner) {
		return "", fault("capacity")
	}
	r := st.remote()
	if len(r.Installs) >= 500 {
		return "", fault("capacity")
	}
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fault("unavailable")
	}
	id, agent, kid := randomToken(), newMemberAgentID(), "remote_"+hex.EncodeToString(public)
	st.Agents[agent] = &Agent{Owner: a.Owner, Keys: map[string]*Key{kid: {Public: public, Changed: now}}}
	r.Installs[id] = &remoteInstall{Owner: a.Owner, Agent: agent, Kid: kid, Private: private, Exp: now.Add(remoteLifetime)}
	code := randomToken()
	r.Codes[hashToken(code)] = &remoteCode{id, a.Client, a.Redirect, a.Challenge, a.Resource, now.Add(time.Minute)}
	return code, nil
}

func validVerifier(v string) bool {
	if len(v) < 43 || len(v) > 128 {
		return false
	}
	for _, c := range v {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.ContainsRune("-._~", c)) {
			return false
		}
	}
	return true
}

func (st *State) exchangeRemote(form url.Values, resource string, now time.Time) (any, error) {
	r := st.remote()
	client := form.Get("client_id")
	c := r.Clients[client]
	if c == nil || !now.Before(c.Exp) || form.Get("resource") != resource || (form.Get("scope") != "" && form.Get("scope") != remoteScope) {
		return nil, fault("invalid_auth")
	}
	install := ""
	switch form.Get("grant_type") {
	case "authorization_code":
		code := r.Codes[hashToken(form.Get("code"))]
		if code == nil || !now.Before(code.Exp) || code.Client != client || code.Redirect != form.Get("redirect_uri") || code.Resource != resource || !validVerifier(form.Get("code_verifier")) {
			return nil, fault("invalid_auth")
		}
		sum := sha256.Sum256([]byte(form.Get("code_verifier")))
		if !hmac.Equal([]byte(code.Challenge), []byte(base64.RawURLEncoding.EncodeToString(sum[:]))) {
			return nil, fault("invalid_auth")
		}
		install = code.Install
		delete(r.Codes, hashToken(form.Get("code")))
	case "refresh_token":
		t := r.Tokens[hashToken(form.Get("refresh_token"))]
		if t == nil || t.Kind != "refresh" || !now.Before(t.Exp) || t.Client != client || t.Resource != resource {
			return nil, fault("invalid_auth")
		}
		install = t.Install
		// One live token pair per installation; refresh rotates both credentials atomically.
		for id, old := range r.Tokens {
			if old.Install == install {
				delete(r.Tokens, id)
			}
		}
	default:
		return nil, fault("invalid_auth")
	}
	i, e := st.remoteInstall(install, now)
	if e != nil {
		return nil, e
	}
	i.Exp = now.Add(remoteLifetime)
	c.Exp = now.Add(remoteLifetime)
	access, refresh := randomToken(), randomToken()
	r.Tokens[hashToken(access)] = &remoteToken{install, client, resource, "access", now.Add(time.Hour)}
	r.Tokens[hashToken(refresh)] = &remoteToken{install, client, resource, "refresh", now.Add(remoteLifetime)}
	return map[string]any{"access_token": access, "token_type": "Bearer", "expires_in": 3600, "refresh_token": refresh, "scope": remoteScope, "resource": resource}, nil
}

func (s *Service) remoteToken(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if s.remoteOrigin() == "" || r.ParseForm() != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		return st.exchangeRemote(r.PostForm, s.remoteOrigin()+"/mcp", now)
	})
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
		return
	}
	writeJSON(w, 200, v)
}

func (s *Service) remoteRevoke(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if r.ParseForm() != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		rs := st.remote()
		t := rs.Tokens[hashToken(r.PostForm.Get("token"))]
		if t == nil || t.Client != r.PostForm.Get("client_id") {
			return nil, nil
		}
		if i := rs.Installs[t.Install]; i != nil {
			i.Exp = now
		}
		st.sweepRemote(now)
		return nil, nil
	})
	if err != nil {
		writeJSON(w, 503, map[string]string{"error": "temporarily_unavailable"})
		return
	}
	w.WriteHeader(200)
}
