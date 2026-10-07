//go:build integration

// Google HTTP regression exercises real Postgres state and signed local provider responses, without SMTP.
package relay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

func TestGoogleHTTP(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" {
		t.Fatal("isolated synthetic database required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s, _ := identityService(t, pool)
	s.Mail = nil
	g, sign := googleFixture(t)
	s.Google = g
	var auth url.Values
	subject := "google-subject"
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("code") != "one-use-code" || oauth2.S256ChallengeFromVerifier(r.PostFormValue("code_verifier")) != auth.Get("code_challenge") || r.PostFormValue("redirect_uri") != "https://example.com/auth/google/callback" {
			http.Error(w, "invalid exchange", 400)
			return
		}
		claims := googleClaims(time.Now(), auth.Get("nonce"))
		claims["sub"] = subject
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(claims)})
	}))
	defer tokenServer.Close()
	g.client = tokenServer.Client()
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "local-fixture", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: tokenServer.URL, AuthStyle: oauth2.AuthStyleInParams}}
	cookies := map[string]*http.Cookie{}
	do := func(method, path string, crossSite bool) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, "https://example.com"+path, nil)
		r.Header.Set("Sec-Fetch-Site", "same-origin")
		if crossSite {
			r.Header.Set("Sec-Fetch-Site", "cross-site")
		}
		for _, c := range cookies {
			if crossSite && c.SameSite == http.SameSiteStrictMode {
				continue
			}
			r.AddCookie(c)
		}
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		for _, c := range w.Result().Cookies() {
			if c.MaxAge < 0 {
				delete(cookies, c.Name)
			} else {
				cookies[c.Name] = c
			}
		}
		return w
	}
	begin := func(path string) string {
		t.Helper()
		w := do("POST", path, false)
		expect(t, w, 303, googleIssuer)
		u, err := url.Parse(w.Header().Get("Location"))
		if err != nil {
			t.Fatal(err)
		}
		auth = u.Query()
		if auth.Get("prompt") != "select_account consent" || auth.Get("nonce") == "" || auth.Get("code_challenge_method") != "S256" {
			t.Fatal("missing nonce/PKCE/explicit confirmation")
		}
		return "/auth/google/callback?code=one-use-code&state=" + auth.Get("state")
	}
	expect(t, do("GET", "/", false), 200, "Google로 계속")
	firstCallback := begin("/auth/google")
	expect(t, do("GET", firstCallback, true), 200, "자기 홈으로 이동")
	first := cookies[sessionCookie].Value
	if cookies[sessionCookie].SameSite != http.SameSiteStrictMode {
		t.Fatal("session weakened")
	}
	expect(t, do("GET", "/home", false), 200, "회원 식별자")
	expect(t, do("POST", "/home/agents", false), 303)
	var member, owner string
	s.mutateState(t, func(st *State) {
		member = st.Sessions[hashToken(first)].Member
		owner = st.Members[member].Owner
		st.Sessions[hashToken(first)].Verified = time.Now().Add(-6 * time.Minute)
	})
	expect(t, do("POST", "/auth/logout-all", false), 403)
	callback := begin("/auth/reauth")
	subject = "foreign-subject"
	expect(t, do("GET", callback, true), 403)
	if cookies[sessionCookie].Value != first {
		t.Fatal("foreign identity replaced session")
	}
	subject = "google-subject"
	callback = begin("/auth/reauth")
	expect(t, do("GET", callback, true), 200)
	second := cookies[sessionCookie].Value
	if second == first {
		t.Fatal("session not rotated")
	}
	s.mutateState(t, func(st *State) {
		if len(st.Members) != 1 || len(st.Owners) != 1 || len(st.Agents) != 1 || st.Sessions[hashToken(first)] != nil || st.Sessions[hashToken(second)].Member != member || st.Members[member].Owner != owner {
			t.Fatal("lost identity/owner/agent continuity")
		}
	})
	expect(t, do("POST", "/auth/logout-all", false), 303)
	expect(t, do("GET", firstCallback, true), 401)
	callback = begin("/auth/google")
	expect(t, do("GET", callback+"&error=access_denied", true), 401)
	callback = begin("/auth/google")
	expect(t, do("POST", "/auth/google", true), 403)
	expect(t, do("GET", strings.Replace(callback, "state=", "state=wrong", 1), true), 401)
	callback = begin("/auth/google")
	expect(t, do("GET", callback, true), 200)
	expect(t, do("GET", "/home", false), 200, "회원 식별자")
	expect(t, do("POST", "/auth/logout", false), 303)
}
