//go:build integration && browser

// Tester-owned independent QA (SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER). Copied beside the fixed-SHA clone in the isolated Linux
// scratch tree; never part of the product SHA. Real Chromium + isolated Postgres + synthetic Google and registered-host origins.
package relay

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

const qaCSP = "default-src 'none'; style-src 'unsafe-inline'; form-action 'self' https://accounts.google.com; frame-ancestors 'none'; base-uri 'none'"

type qaDriverResult struct {
	Name     string `json:"name"`
	Stage    string `json:"stage"`
	Error    string `json:"error"`
	CSP      []string
	Blocked  []string
	Steps    []string
	Docs     []qaDoc
	Arrivals []struct {
		Host string
		Path string
		Keys []string
		Held bool
	}
	ReturnPage *struct {
		H1, LinkText string
		Rect         []int
		InnerWidth   int
		InnerHeight  int
		ScrollWidth  int
		RefreshZero  bool
		Origin       string
	}
	ArrivedAfterClick *bool
	FinalLocation     string
	Refusal           *struct {
		Text  string
		Links []string
	}
	Revisit *struct{ Text string }
	Shots   []string
	History []string
}

type qaDoc struct {
	Method  string
	Status  int
	URL     string
	Headers map[string]string
}

type qaEnv struct {
	t      *testing.T
	s      *Service
	relay  http.Handler
	rules  string
	chrome string
	out    string
	mu     sync.Mutex
	hits   []url.Values
}

func newQAEnv(t *testing.T) *qaEnv {
	chrome := os.Getenv("TEST_BROWSER")
	if chrome == "" || os.Getenv("QA_DRIVER") == "" || os.Getenv("QA_OUT") == "" {
		t.Fatal("TEST_BROWSER, QA_DRIVER and QA_OUT are required")
	}
	pool := messagePool(t)
	s, _ := identityService(t, pool)
	s.Mail = nil
	g, sign := googleFixture(t)
	s.Google = g
	e := &qaEnv{t: t, s: s, relay: s.Handler(), chrome: chrome, out: os.Getenv("QA_OUT")}
	var auth url.Values
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.mu.Lock()
		nonce := auth.Get("nonce")
		e.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(googleClaims(time.Now(), nonce))})
	}))
	t.Cleanup(provider.Close)
	g.client = provider.Client()
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "test-only", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: provider.URL, AuthStyle: oauth2.AuthStyleInParams}}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Host {
		case "example.com":
			e.relay.ServeHTTP(w, r)
		case "accounts.google.com":
			e.mu.Lock()
			auth = r.URL.Query()
			e.mu.Unlock()
			http.Redirect(w, r, "https://example.com/auth/google/callback?code=fixture&state="+url.QueryEscape(r.URL.Query().Get("state")), 302)
		case "callback.test":
			if r.URL.Path != "/oauth" {
				http.NotFound(w, r)
				return
			}
			e.mu.Lock()
			e.hits = append(e.hits, r.URL.Query())
			e.mu.Unlock()
			_, _ = w.Write([]byte("returned"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	port := server.Listener.Addr().(*net.TCPAddr).Port
	e.rules = fmt.Sprintf("MAP example.com 127.0.0.1:%d,MAP accounts.google.com 127.0.0.1:%d,MAP callback.test 127.0.0.1:%d", port, port, port)
	return e
}

func (e *qaEnv) do(method, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.relay.ServeHTTP(w, r)
	return w
}

func (e *qaEnv) register() string {
	w := e.do("POST", "/remote/oauth/register", `{"client_name":"qa host","redirect_uris":["https://callback.test/oauth"]}`)
	expect(e.t, w, 201)
	var c struct {
		ID string `json:"client_id"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &c)
	return c.ID
}

func (e *qaEnv) authorize(client, verifier string) string {
	q := url.Values{"response_type": {"code"}, "client_id": {client}, "redirect_uri": {"https://callback.test/oauth"}, "resource": {"https://example.com/mcp"}, "scope": {remoteScope}, "state": {"qa-state"}, "code_challenge": {oauth2.S256ChallengeFromVerifier(verifier)}, "code_challenge_method": {"S256"}}
	return "https://example.com/remote/oauth/authorize?" + q.Encode()
}

func (e *qaEnv) installs() int {
	n := 0
	e.s.mutateState(e.t, func(st *State) { n = len(st.remote().Installs) })
	return n
}

// drive runs one Chromium scenario. gate (optional) runs while the browser waits on the consent page, before the decision click.
func (e *qaEnv) drive(sc map[string]any, client, verifier string, gate func()) qaDriverResult {
	e.t.Helper()
	b, _ := json.Marshal(sc)
	if gate != nil {
		sc["gate"] = true
		b, _ = json.Marshal(sc)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "node", os.Getenv("QA_DRIVER"), e.chrome, e.rules, e.authorize(client, verifier), e.t.TempDir(), string(b), e.out)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		e.t.Fatal(err)
	}
	var res qaDriverResult
	sc2 := bufio.NewScanner(out)
	sc2.Buffer(make([]byte, 1<<20), 1<<20)
	for sc2.Scan() {
		line := sc2.Text()
		if line == "READY" {
			gate()
			_, _ = fmt.Fprintln(in, "GO")
			continue
		}
		if strings.HasPrefix(line, "{") {
			if err := json.Unmarshal([]byte(line), &res); err != nil {
				e.t.Fatal(err)
			}
		}
	}
	_ = cmd.Wait()
	// Redacted evidence only: no URL query, code, token or cookie is printed.
	e.t.Logf("driver %s stage=%s error=%q steps=%v csp=%d blocked=%v final=%s arrivals=%v shots=%d", res.Name, res.Stage, res.Error, res.Steps, len(res.CSP), res.Blocked, res.FinalLocation, res.Arrivals, len(res.Shots))
	e.t.Logf("  history entries: %v", res.History)
	for _, d := range res.Docs {
		e.t.Logf("  doc %s %d %s", d.Method, d.Status, d.URL)
	}
	if res.Error != "" || res.Stage != "done" {
		e.t.Fatalf("driver did not finish: stage=%s error=%s", res.Stage, res.Error)
	}
	return res
}

func (e *qaEnv) takeHits() []url.Values {
	e.mu.Lock()
	defer e.mu.Unlock()
	h := e.hits
	e.hits = nil
	return h
}

func (r qaDriverResult) post() *qaDoc {
	for i := len(r.Docs) - 1; i >= 0; i-- {
		if r.Docs[i].Method == "POST" && strings.Contains(r.Docs[i].URL, "/remote/consent/") {
			return &r.Docs[i]
		}
	}
	return nil
}

func (e *qaEnv) postStatuses(r qaDriverResult) []int {
	var s []int
	for _, d := range r.Docs {
		if d.Method == "POST" && strings.Contains(d.URL, "/remote/consent/") {
			s = append(s, d.Status)
		}
	}
	return s
}

func (e *qaEnv) mustHeaders(d *qaDoc) {
	e.t.Helper()
	if d == nil {
		e.t.Fatal("no consent POST response recorded")
	}
	h := d.Headers
	if h["cache-control"] != "no-store" || h["referrer-policy"] != "no-referrer" || h["content-security-policy"] != qaCSP || h["x-content-type-options"] != "nosniff" || !strings.HasPrefix(h["content-type"], "text/html") {
		e.t.Fatalf("security headers changed on consent POST response: %v", h)
	}
}

func (e *qaEnv) cleanBrowser(r qaDriverResult) {
	e.t.Helper()
	if len(r.CSP) != 0 || len(r.Blocked) != 0 {
		e.t.Fatalf("CSP violation or blocked request in Chromium: csp=%v blocked=%v", r.CSP, r.Blocked)
	}
}

func (e *qaEnv) exchange(client, verifier, code string, mutate func(url.Values)) *httptest.ResponseRecorder {
	form := url.Values{"grant_type": {"authorization_code"}, "client_id": {client}, "redirect_uri": {"https://callback.test/oauth"}, "resource": {"https://example.com/mcp"}, "code": {code}, "code_verifier": {verifier}}
	if mutate != nil {
		mutate(form)
	}
	return e.do("POST", "/remote/oauth/token", form.Encode())
}

func (e *qaEnv) oneHit() url.Values {
	e.t.Helper()
	h := e.takeHits()
	if len(h) != 1 {
		e.t.Fatalf("expected exactly one callback arrival, got %d", len(h))
	}
	return h[0]
}

func (e *qaEnv) noHit() {
	e.t.Helper()
	if h := e.takeHits(); len(h) != 0 {
		e.t.Fatalf("callback must not be reached, got %d arrivals", len(h))
	}
}

func (e *qaEnv) restartNotice(text string) {
	e.t.Helper()
	if !strings.Contains(text, "다시 시작") || strings.Contains(text, "callback.test") || strings.Contains(text, "code=") {
		e.t.Fatalf("refusal must carry the restart notice and no callback/code: %q", text)
	}
}

func TestQAOAuthReturnAutoApproveAndTokenBoundaries(t *testing.T) {
	e := newQAEnv(t)
	client, other := e.register(), e.register()
	verifier := oauth2.GenerateVerifier()
	before := e.installs()
	r := e.drive(map[string]any{"name": "approve-auto-desktop", "revisit": true}, client, verifier, nil)
	e.cleanBrowser(r)
	e.mustHeaders(r.post())
	if r.post().Status != 200 || r.FinalLocation != "https://callback.test/oauth" {
		t.Fatalf("approval did not reach callback automatically: %d %s", r.post().Status, r.FinalLocation)
	}
	hit := e.oneHit()
	code := hit.Get("code")
	if hit.Get("state") != "qa-state" || hit.Get("iss") != "https://example.com" || code == "" || hit.Get("error") != "" {
		t.Fatal("callback parameters state/iss/code")
	}
	t.Logf("callback: state=%v iss=%v code_present=%v error_absent=%v", hit.Get("state") == "qa-state", hit.Get("iss") == "https://example.com", code != "", hit.Get("error") == "")
	if e.installs() != before+1 {
		t.Fatal("approval must create exactly one install")
	}
	// Consumed consent: restart notice, no callback echo.
	if r.Revisit == nil {
		t.Fatal("no revisit capture")
	}
	e.restartNotice(r.Revisit.Text)
	bad := map[string]func(url.Values){
		"wrong_verifier":     func(f url.Values) { f.Set("code_verifier", oauth2.GenerateVerifier()) },
		"wrong_redirect_uri": func(f url.Values) { f.Set("redirect_uri", "https://callback.test/other") },
		"wrong_resource":     func(f url.Values) { f.Set("resource", "https://example.com/other") },
		"other_client":       func(f url.Values) { f.Set("client_id", other) },
		"short_verifier":     func(f url.Values) { f.Set("code_verifier", "short") },
		"missing_verifier":   func(f url.Values) { f.Del("code_verifier") },
		"wrong_scope":        func(f url.Values) { f.Set("scope", "admin") },
	}
	for name, m := range bad {
		w := e.exchange(client, verifier, code, m)
		t.Logf("token boundary %s -> %d", name, w.Code)
		expect(t, w, 400, "invalid_grant")
	}
	w := e.exchange(client, verifier, code, nil)
	t.Logf("token exchange with correct PKCE/redirect/resource -> %d no-store=%v", w.Code, w.Header().Get("Cache-Control") == "no-store")
	expect(t, w, 200, "access_token", "refresh_token")
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("token response cacheable")
	}
	expect(t, e.exchange(client, verifier, code, nil), 400, "invalid_grant")
	t.Log("code replay -> 400")
}

func TestQAOAuthReturnExpiredCodeRefused(t *testing.T) {
	e := newQAEnv(t)
	client := e.register()
	verifier := oauth2.GenerateVerifier()
	r := e.drive(map[string]any{"name": "expired-code"}, client, verifier, nil)
	e.cleanBrowser(r)
	code := e.oneHit().Get("code")
	e.s.mutateState(t, func(st *State) {
		for _, c := range st.remote().Codes {
			c.Exp = time.Now().Add(-time.Second)
		}
	})
	w := e.exchange(client, verifier, code, nil)
	t.Logf("expired code exchange -> %d", w.Code)
	expect(t, w, 400, "invalid_grant")
}

func TestQAOAuthReturnLinkDesktopAndMobile(t *testing.T) {
	for _, mobile := range []bool{false, true} {
		name := "approve-link-desktop"
		if mobile {
			name = "approve-link-mobile"
		}
		t.Run(name, func(t *testing.T) {
			e := newQAEnv(t)
			client := e.register()
			verifier := oauth2.GenerateVerifier()
			r := e.drive(map[string]any{"name": name, "mobile": mobile, "hold": true}, client, verifier, nil)
			e.cleanBrowser(r)
			e.mustHeaders(r.post())
			p := r.ReturnPage
			if p == nil || !p.RefreshZero || p.Origin != "https://example.com" || !strings.Contains(p.H1, "돌아가는 중") || !strings.Contains(p.LinkText, "돌아가기") {
				t.Fatalf("return page content: %+v", p)
			}
			if p.ScrollWidth > p.InnerWidth || p.Rect[0] < 0 || p.Rect[0]+p.Rect[2] > p.InnerWidth || p.Rect[1]+p.Rect[3] > p.InnerHeight || p.Rect[3] < 20 {
				t.Fatalf("return link not usable in viewport: %+v", p)
			}
			t.Logf("return page %s: viewport=%dx%d link rect=%v scrollWidth=%d meta-refresh-0=%v", name, p.InnerWidth, p.InnerHeight, p.Rect, p.ScrollWidth, p.RefreshZero)
			if r.ArrivedAfterClick == nil || !*r.ArrivedAfterClick || len(r.Arrivals) != 2 || !r.Arrivals[0].Held {
				t.Fatalf("held navigation then real link click must arrive: %+v", r.Arrivals)
			}
			hit := e.oneHit()
			if hit.Get("state") != "qa-state" || hit.Get("iss") != "https://example.com" || hit.Get("code") == "" {
				t.Fatal("callback parameters after link click")
			}
			expect(t, e.exchange(client, verifier, hit.Get("code"), nil), 200, "access_token")
			t.Log("link click callback arrival and PKCE token exchange -> 200")
		})
	}
}

func TestQAOAuthReturnDeny(t *testing.T) {
	e := newQAEnv(t)
	client := e.register()
	before := e.installs()
	r := e.drive(map[string]any{"name": "deny-desktop", "decision": "deny", "revisit": true}, client, oauth2.GenerateVerifier(), nil)
	e.cleanBrowser(r)
	e.mustHeaders(r.post())
	hit := e.oneHit()
	if hit.Get("error") != "access_denied" || hit.Get("state") != "qa-state" || hit.Get("iss") != "https://example.com" || hit.Get("code") != "" {
		t.Fatal("deny callback must carry error/state/iss and no code")
	}
	if e.installs() != before || r.Revisit == nil {
		t.Fatal("deny created an install or capture missing")
	}
	e.restartNotice(r.Revisit.Text)
}

func TestQAOAuthReturnRefusalBoundaries(t *testing.T) {
	type tc struct {
		name   string
		mobile bool
		sc     map[string]any
		gate   func(e *qaEnv)
		status int
		retry  bool
	}
	past := func(e *qaEnv) {
		e.s.mutateState(e.t, func(st *State) {
			for _, a := range st.remote().Attempts {
				a.Exp = time.Now().Add(-time.Second)
			}
		})
	}
	rebind := func(e *qaEnv) {
		e.s.mutateState(e.t, func(st *State) {
			for _, a := range st.remote().Attempts {
				a.Owner = "another-google-account-owner"
			}
		})
	}
	cases := []tc{
		{name: "expired-attempt-desktop", sc: map[string]any{}, gate: past, status: 401},
		{name: "owner-binding-mobile", mobile: true, sc: map[string]any{}, gate: rebind, status: 401},
		{name: "no-login-desktop", sc: map[string]any{"clearCookies": true}, status: 401},
		{name: "csrf-tamper-mobile", mobile: true, sc: map[string]any{"tamperCsrf": true, "retry": true}, status: 403, retry: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := newQAEnv(t)
			client := e.register()
			before := e.installs()
			c.sc["name"], c.sc["mobile"] = c.name, c.mobile
			var gate func()
			if c.gate != nil {
				gate = func() { c.gate(e) }
			}
			r := e.drive(c.sc, client, oauth2.GenerateVerifier(), gate)
			e.cleanBrowser(r)
			st := e.postStatuses(r)
			t.Logf("consent POST statuses: %v", st)
			if len(st) == 0 || st[0] != c.status {
				t.Fatalf("expected first consent POST status %d, got %v", c.status, st)
			}
			e.mustHeaders(&r.Docs[docIndex(r, "POST")])
			if r.Refusal == nil {
				t.Fatal("refusal page not shown")
			}
			e.restartNotice(r.Refusal.Text)
			t.Logf("refusal page text: %q links=%v", r.Refusal.Text, r.Refusal.Links)
			if c.retry {
				// A CSRF refusal keeps the attempt alive: the genuine form then succeeds exactly once.
				if len(st) != 2 || st[1] != 200 {
					t.Fatalf("retry after CSRF refusal: %v", st)
				}
				e.oneHit()
				if e.installs() != before+1 {
					t.Fatal("retry must create exactly one install")
				}
				return
			}
			e.noHit()
			if e.installs() != before {
				t.Fatal("refused consent created an install")
			}
		})
	}
}

func docIndex(r qaDriverResult, method string) int {
	for i, d := range r.Docs {
		if d.Method == method && strings.Contains(d.URL, "/remote/consent/") {
			return i
		}
	}
	return 0
}

func TestQAOAuthReturnAuthorizeBoundaries(t *testing.T) {
	e := newQAEnv(t)
	client := e.register()
	verifier := oauth2.GenerateVerifier()
	base := func() url.Values {
		u, _ := url.Parse(e.authorize(client, verifier))
		return u.Query()
	}
	attempts := func() int {
		n := 0
		e.s.mutateState(t, func(st *State) { n = len(st.remote().Attempts) })
		return n
	}
	before := attempts()
	for name, m := range map[string]func(url.Values){
		"unregistered_redirect": func(q url.Values) { q.Set("redirect_uri", "https://evil.test/oauth") },
		"wrong_resource":        func(q url.Values) { q.Set("resource", "https://evil.test/mcp") },
		"plain_pkce":            func(q url.Values) { q.Set("code_challenge_method", "plain") },
		"no_challenge":          func(q url.Values) { q.Del("code_challenge") },
		"unknown_client":        func(q url.Values) { q.Set("client_id", "nope") },
		"wrong_scope":           func(q url.Values) { q.Set("scope", "admin") },
		"token_response_type":   func(q url.Values) { q.Set("response_type", "token") },
	} {
		q := base()
		m(q)
		w := e.do("GET", "/remote/oauth/authorize?"+q.Encode(), "")
		t.Logf("authorize boundary %s -> %d", name, w.Code)
		expect(t, w, 400, "invalid_request")
	}
	if attempts() != before {
		t.Fatal("rejected authorize requests created attempts")
	}
	w := e.do("GET", "/remote/oauth/authorize?"+base().Encode(), "")
	expect(t, w, 303)
	t.Logf("valid authorize -> %d to consent page", w.Code)
}

// Capacity refusal: an owner at the agent limit approves; the page must still be a safe HTML refusal. Observation only for the message wording.
func TestQAOAuthReturnCapacityRefusal(t *testing.T) {
	e := newQAEnv(t)
	client := e.register()
	before := e.installs()
	fill := func() {
		e.s.mutateState(t, func(st *State) {
			for _, a := range st.remote().Attempts {
				for i := 0; i < ownerAgentRecords; i++ {
					st.Agents[fmt.Sprintf("agent_qa_capacity_%d", i)] = &Agent{Owner: a.Owner}
				}
			}
		})
	}
	r := e.drive(map[string]any{"name": "capacity-desktop"}, client, oauth2.GenerateVerifier(), fill)
	e.cleanBrowser(r)
	st := e.postStatuses(r)
	t.Logf("capacity consent POST statuses: %v", st)
	if len(st) != 1 || st[0] != 409 || r.Refusal == nil {
		t.Fatalf("expected HTML refusal with 409, got %v", st)
	}
	t.Logf("capacity refusal page text: %q", r.Refusal.Text)
	e.noHit()
	if e.installs() != before {
		t.Fatal("capacity refusal created an install")
	}
}
