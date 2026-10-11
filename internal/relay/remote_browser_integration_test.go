//go:build integration && browser

// A real Chromium, isolated Postgres and synthetic Google and host origins exercise the remote consent return.
package relay

import (
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

// A real Chromium enforces form-action on redirects of a form submission, which an HTTP recorder cannot observe.
func TestRemoteConsentBrowserReturn(t *testing.T) {
	chrome := os.Getenv("TEST_BROWSER")
	if chrome == "" {
		t.Fatal("TEST_BROWSER must name a Chromium executable")
	}
	pool := messagePool(t)
	s, _ := identityService(t, pool)
	s.Mail = nil
	g, sign := googleFixture(t)
	s.Google = g
	var mu sync.Mutex
	var auth, callback url.Values
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		nonce := auth.Get("nonce")
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(googleClaims(time.Now(), nonce))})
	}))
	defer provider.Close()
	g.client = provider.Client()
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "test-only", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: provider.URL, AuthStyle: oauth2.AuthStyleInParams}}
	relay := s.Handler()
	// One TLS listener plays relay, synthetic Google and the registered host by Host header.
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Host {
		case "example.com":
			relay.ServeHTTP(w, r)
		case "accounts.google.com":
			mu.Lock()
			auth = r.URL.Query()
			mu.Unlock()
			http.Redirect(w, r, "https://example.com/auth/google/callback?code=fixture&state="+url.QueryEscape(r.URL.Query().Get("state")), 302)
		case "callback.test":
			if r.URL.Path != "/oauth" {
				http.NotFound(w, r) // favicon
				return
			}
			mu.Lock()
			callback = r.URL.Query()
			mu.Unlock()
			_, _ = w.Write([]byte("returned"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	do := func(method, path, body string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		relay.ServeHTTP(w, r)
		return w
	}
	w := do("POST", "/remote/oauth/register", `{"client_name":"browser host","redirect_uris":["https://callback.test/oauth"]}`)
	expect(t, w, 201)
	var client struct {
		ID string `json:"client_id"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &client)
	verifier := oauth2.GenerateVerifier()
	q := url.Values{"response_type": {"code"}, "client_id": {client.ID}, "redirect_uri": {"https://callback.test/oauth"}, "resource": {"https://example.com/mcp"}, "scope": {remoteScope}, "state": {"browser-state"}, "code_challenge": {oauth2.S256ChallengeFromVerifier(verifier)}, "code_challenge_method": {"S256"}}
	port := server.Listener.Addr().(*net.TCPAddr).Port
	rules := fmt.Sprintf("MAP example.com 127.0.0.1:%d,MAP accounts.google.com 127.0.0.1:%d,MAP callback.test 127.0.0.1:%d", port, port, port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	out, err := exec.CommandContext(ctx, "node", "testdata/consent_browser.mjs", chrome, rules, "https://example.com/remote/oauth/authorize?"+q.Encode(), t.TempDir()).Output()
	t.Logf("browser: %s", out)
	if err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	if callback == nil || callback.Get("state") != "browser-state" || callback.Get("iss") != "https://example.com" || callback.Get("code") == "" {
		t.Fatal("approved consent did not return to the registered callback")
	}
	form := url.Values{"grant_type": {"authorization_code"}, "client_id": {client.ID}, "redirect_uri": {"https://callback.test/oauth"}, "resource": {"https://example.com/mcp"}, "code": {callback.Get("code")}, "code_verifier": {verifier}}
	expect(t, do("POST", "/remote/oauth/token", form.Encode()), 200, "access_token")
	// The consumed request explains how to restart without echoing the callback or code.
	var browser struct{ Steps []string }
	if json.Unmarshal(out, &browser) != nil || len(browser.Steps) != 4 {
		t.Fatal("browser steps")
	}
	w = do("GET", strings.TrimPrefix(browser.Steps[2], "https://example.com"), "")
	expect(t, w, 401, "다시 시작")
	if strings.Contains(w.Body.String(), "callback.test") {
		t.Fatal("expired consent echoed callback")
	}
}
