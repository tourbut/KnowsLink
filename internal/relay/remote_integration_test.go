//go:build integration

// Isolated Postgres and a signed local Google provider exercise actual remote OAuth HTTP and persistent MCP authorization.
package relay

import (
	"context"
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestRemoteOAuthHTTP(t *testing.T) {
	pool := messagePool(t)
	s, _ := identityService(t, pool)
	s.Mail = nil
	g, sign := googleFixture(t)
	s.Google = g
	var auth url.Values
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if oauth2.S256ChallengeFromVerifier(r.PostFormValue("code_verifier")) != auth.Get("code_challenge") {
			http.Error(w, "PKCE", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(googleClaims(time.Now(), auth.Get("nonce")))})
	}))
	defer provider.Close()
	g.client = provider.Client()
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "test-only", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: provider.URL, AuthStyle: oauth2.AuthStyleInParams}}
	cookies := map[string]*http.Cookie{}
	do := func(service *Service, method, path, body, content, token string, browser bool) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		r.Header.Set("Content-Type", content)
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		if browser {
			for _, c := range cookies {
				r.AddCookie(c)
			}
			if r.Method == "POST" {
				r.Header.Set("Sec-Fetch-Site", "same-origin")
			}
		}
		w := httptest.NewRecorder()
		service.Handler().ServeHTTP(w, r)
		if browser {
			for _, c := range w.Result().Cookies() {
				if c.MaxAge < 0 {
					delete(cookies, c.Name)
				} else {
					cookies[c.Name] = c
				}
			}
		}
		return w
	}
	formType := "application/x-www-form-urlencoded"
	metadata := do(s, "GET", "/.well-known/oauth-protected-resource/mcp", "", "", "", false)
	expect(t, metadata, 200, "https://example.com/mcp")
	expect(t, do(s, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"server/discover"}`, "application/json", "wrong", false), 401, "invalid_token")
	// Refusal probes spend their own client-IP bucket so the main flow keeps its anonymous rate budget.
	probe := func(method, path, body string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.RemoteAddr = "198.51.100.7:1234"
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		return w
	}
	expect(t, probe("POST", "/remote/oauth/register", `{"client_name":"x","redirect_uris":["http://chatgpt.com/callback"]}`), 400, "invalid_redirect_uri")
	expect(t, probe("POST", "/remote/oauth/register", `{"client_name":"x","redirect_uris":["https://chatgpt.com/callback"],"grant_types":["implicit"]}`), 400, "invalid_client_metadata")
	expect(t, probe("POST", "/remote/oauth/register", `{"client_name":"x","redirect_uris":["https://chatgpt.com/callback"],"response_types":["token"]}`), 400, "invalid_client_metadata")
	var agents []string
	for range 2 {
		w := do(s, "POST", "/remote/oauth/register", `{"client_name":"test host","redirect_uris":["https://chatgpt.com/callback"],"token_endpoint_auth_method":"none","grant_types":["authorization_code","refresh_token"],"response_types":["code"]}`, "application/json", "", false)
		expect(t, w, 201)
		var client struct {
			ID string `json:"client_id"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &client)
		verifier := oauth2.GenerateVerifier()
		q := url.Values{"response_type": {"code"}, "client_id": {client.ID}, "redirect_uri": {"https://chatgpt.com/callback"}, "resource": {"https://example.com/mcp"}, "scope": {remoteScope}, "state": {"original-state"}, "code_challenge": {oauth2.S256ChallengeFromVerifier(verifier)}, "code_challenge_method": {"S256"}}
		// The real client_id keeps the unregistered_redirect case on the redirect allowlist, not the client lookup.
		for name, mutate := range map[string]func(url.Values){
			"unregistered_redirect": func(v url.Values) { v.Set("redirect_uri", "https://evil.example/callback") },
			"unknown_client":        func(v url.Values) { v.Set("client_id", "nope") },
			"plain_pkce":            func(v url.Values) { v.Set("code_challenge_method", "plain") },
			"wrong_resource":        func(v url.Values) { v.Set("resource", "https://other.example/mcp") },
		} {
			bad := url.Values{}
			for k, v := range q {
				bad[k] = append([]string(nil), v...)
			}
			mutate(bad)
			if w := probe("GET", "/remote/oauth/authorize?"+bad.Encode(), ""); w.Code != 400 || w.Header().Get("Location") != "" {
				t.Fatal("authorize accepted", name, w.Code)
			}
		}
		w = do(s, "GET", "/remote/oauth/authorize?"+q.Encode(), "", "", "", true)
		expect(t, w, 303)
		consent := w.Header().Get("Location")
		id := strings.TrimPrefix(consent, "/remote/consent/")
		expect(t, do(s, "GET", consent, "", "", "", true), 200, "/auth/google", "Google")
		w = do(s, "POST", "/auth/google", url.Values{"remote": {id}}.Encode(), formType, "", true)
		expect(t, w, 303)
		u, _ := url.Parse(w.Header().Get("Location"))
		auth = u.Query()
		w = do(s, "GET", "/auth/google/callback?code=fixture&state="+auth.Get("state"), "", "", "", true)
		expect(t, w, 200, consent)
		expect(t, do(s, "GET", consent, "", "", "", true), 200, "이 설치 연결 승인", "비공개 DB")
		expect(t, do(s, "POST", consent, "decision=approve&csrf=wrong", formType, "", true), 403)
		w = do(s, "POST", consent, url.Values{"decision": {"approve"}, "name": {"합성 원격 도우미"}, "csrf": {csrf(cookies[sessionCookie].Value, id)}}.Encode(), formType, "", true)
		expect(t, w, 200, `http-equiv="refresh"`)
		u = remoteReturnURL(t, w)
		if u.Query().Get("iss") != "https://example.com" || u.Query().Get("state") != "original-state" {
			t.Fatal("issuer/state response binding")
		}
		form := url.Values{"grant_type": {"authorization_code"}, "client_id": {client.ID}, "redirect_uri": {"https://chatgpt.com/callback"}, "resource": {"https://example.com/mcp"}, "code": {u.Query().Get("code")}, "code_verifier": {verifier}}
		w = do(s, "POST", "/remote/oauth/token", form.Encode(), formType, "", false)
		expect(t, w, 200)
		if w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("token response cacheable")
		}
		var tokens struct {
			Access  string `json:"access_token"`
			Refresh string `json:"refresh_token"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &tokens)
		expect(t, do(s, "POST", "/remote/oauth/token", form.Encode(), formType, "", false), 400, "invalid_grant")
		other := &Service{Redis: testRedis(), Pool: pool, Google: g}
		w = do(other, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"server/discover"}`, "application/json", tokens.Access, false)
		expect(t, w, 200, "2026-07-28", "events")
		remoteEraChecks(t, other, tokens.Access)
		w = do(other, "POST", "/mcp", `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"knowslink_status","arguments":{}}}`, "application/json", tokens.Access, false)
		expect(t, w, 200, "automatic_send")
		// A 4096-byte text JSON-escaped as \uXXXX exceeds 8 KiB but stays within the raw envelope budget.
		long := `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"knowslink_send","arguments":{"to":"agent_unknown","text":"` + strings.Repeat(`한`, 1365) + `x","idempotency_key":"remote-long-key-01"}}}`
		expect(t, do(other, "POST", "/mcp", long, "application/json", tokens.Access, false), 200, `"jsonrpc"`)
		var out struct {
			Result struct {
				Value struct {
					Agent string `json:"agent"`
				} `json:"structuredContent"`
			} `json:"result"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		agents = append(agents, out.Result.Value.Agent)
		// U2-A: the consent-screen name labels the same agent the installed client reports.
		s.mutateState(t, func(st *State) {
			if a := st.Agents[out.Result.Value.Agent]; a == nil || a.Name != "합성 원격 도우미" {
				t.Fatal("remote consent name is not on the client's agent")
			}
		})
		expect(t, do(other, "POST", "/v1/text/pull", `{}`, "application/json", tokens.Access, false), 401)
		expect(t, do(other, "POST", "/mcp", `{"jsonrpc":"2.0","id":3,"method":"events/list"}`, "application/json", tokens.Access, false), 200, remoteEvent)
		// A callback to a non-public address fails challenge verification, never activates a subscription.
		payload := `{"jsonrpc":"2.0","id":4,"method":"events/subscribe","params":{"name":"knowslink.text.arrived","arguments":{},"delivery":{"mode":"webhook","url":"https://127.0.0.1/callback","secret":"whsec_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="}}}`
		expect(t, do(other, "POST", "/mcp", payload, "application/json", tokens.Access, false), 200, "-32015")
		var count int
		_, err := s.transaction(context.Background(), func(st *State, _ time.Time) (any, error) { count = len(st.remote().Subscriptions); return nil, nil })
		if err != nil || count != 0 {
			t.Fatal("failed callback persisted active subscription", err)
		}
		if len(agents) == 1 {
			form = url.Values{"grant_type": {"refresh_token"}, "client_id": {client.ID}, "resource": {"https://example.com/mcp"}, "refresh_token": {tokens.Refresh}}
			w = do(s, "POST", "/remote/oauth/token", form.Encode(), formType, "", false)
			expect(t, w, 200)
			var rotated struct {
				Access string `json:"access_token"`
			}
			_ = json.Unmarshal(w.Body.Bytes(), &rotated)
			expect(t, do(s, "POST", "/remote/oauth/token", form.Encode(), formType, "", false), 400, "invalid_grant")
			// This must remain revoked in a new transaction, rather than rolling the replay refusal back.
			tokens.Access = rotated.Access
		} else {
			form = url.Values{"client_id": {client.ID}, "token": {tokens.Refresh}}
			expect(t, do(s, "POST", "/remote/oauth/revoke", form.Encode(), formType, "", false), 200)
		}
		expect(t, do(other, "POST", "/mcp", `{"jsonrpc":"2.0","id":5,"method":"server/discover"}`, "application/json", tokens.Access, false), 401)
	}
	if len(agents) != 2 || agents[0] == "" || agents[0] == agents[1] {
		t.Fatal("same Google account did not create distinct installation agents")
	}
}

func remoteReturnURL(t *testing.T, w *httptest.ResponseRecorder) *url.URL {
	t.Helper()
	m := regexp.MustCompile(`<a href="([^"]+)"`).FindStringSubmatch(w.Body.String())
	if m == nil {
		t.Fatal("return page has no callback link")
	}
	u, err := url.Parse(html.UnescapeString(m[1]))
	if err != nil {
		t.Fatal(err)
	}
	return u
}

// SAR-MCP-UPGRADE-001: a legacy (initialize) client sees the same answers as before, and a 2026-07-28 request on the same
// endpoint gets the stateless result fields, 404 for an absent method and a tool refusal (not 404) for an unknown tool.
func remoteEraChecks(t *testing.T, s *Service, token string) {
	t.Helper()
	call := func(body string, headers ...string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", "https://example.com/mcp", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer "+token)
		for i := 0; i+1 < len(headers); i += 2 {
			r.Header.Set(headers[i], headers[i+1])
		}
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		t.Logf("%s %s -> %d %.160s", r.Header.Get("MCP-Protocol-Version"), body[:min(len(body), 72)], w.Code, w.Body.String())
		return w
	}
	refuse := func(w *httptest.ResponseRecorder, absent ...string) {
		t.Helper()
		for _, s := range absent {
			if strings.Contains(w.Body.String(), s) {
				t.Fatalf("legacy answer gained %q: %s", s, w.Body.String())
			}
		}
	}
	for _, v := range []string{"2025-03-26", "2025-06-18", "2025-11-25"} {
		w := call(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"`+v+`","capabilities":{},"clientInfo":{"name":"legacy","version":"1"}}}`, "MCP-Protocol-Version", "")
		expect(t, w, 200, `"protocolVersion":"`+v+`"`, `"tools"`)
		refuse(w, "resultType", "ttlMs")
	}
	w := call(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"old","version":"1"}}}`)
	expect(t, w, 200, `"protocolVersion":"2025-11-25"`)
	expect(t, call(`{"jsonrpc":"2.0","method":"notifications/initialized"}`, "MCP-Protocol-Version", "2025-06-18"), 202)
	w = call(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`, "MCP-Protocol-Version", "2025-06-18")
	expect(t, w, 200, "knowslink_send", "knowslink_receive", "knowslink_ack")
	refuse(w, "resultType", "ttlMs", "cacheScope")
	expect(t, call(`{"jsonrpc":"2.0","id":3,"method":"resources/list"}`), 200, "unsupported_method")

	meta := `"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}`
	modern := func(method, name, params string) *httptest.ResponseRecorder {
		return call(`{"jsonrpc":"2.0","id":9,"method":"`+method+`","params":{`+params+meta+`}}`, "MCP-Protocol-Version", "2026-07-28", "Mcp-Method", method, "Mcp-Name", name)
	}
	expect(t, modern("server/discover", "", ""), 200, `"resultType":"complete"`, `"ttlMs":300000`, `"cacheScope":"public"`, "untrusted data", "io.modelcontextprotocol/serverInfo")
	expect(t, modern("tools/list", "", ""), 200, `"resultType":"complete"`, `"ttlMs":300000`, "knowslink_send")
	w = modern("tools/call", "knowslink_status", `"name":"knowslink_status","arguments":{},`)
	expect(t, w, 200, `"resultType":"complete"`, "automatic_send")
	refuse(w, "ttlMs")
	expect(t, modern("tools/call", "knowslink_nope", `"name":"knowslink_nope","arguments":{},`), 200, "unsupported_method")
	expect(t, modern("resources/list", "", ""), 404, `"code":-32601`)
	expect(t, modern("initialize", "", ""), 404, `"code":-32601`)
}
