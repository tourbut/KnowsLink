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
		w = do(s, "POST", consent, url.Values{"decision": {"approve"}, "csrf": {csrf(cookies[sessionCookie].Value, id)}}.Encode(), formType, "", true)
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
		other := &Service{Pool: pool, Google: g}
		w = do(other, "POST", "/mcp", `{"jsonrpc":"2.0","id":1,"method":"server/discover"}`, "application/json", tokens.Access, false)
		expect(t, w, 200, "2026-07-28", "events")
		w = do(other, "POST", "/mcp", `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"knowslink_status","arguments":{}}}`, "application/json", tokens.Access, false)
		expect(t, w, 200, "automatic_send")
		var out struct {
			Result struct {
				Value struct {
					Agent string `json:"agent"`
				} `json:"structuredContent"`
			} `json:"result"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		agents = append(agents, out.Result.Value.Agent)
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
