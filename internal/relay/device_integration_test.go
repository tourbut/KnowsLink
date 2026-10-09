//go:build integration

// Google client HTTP flow uses a signed local provider and isolated Postgres; no real accounts or secrets are needed.
package relay

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
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

func TestGoogleDeviceHTTP(t *testing.T) {
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
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "local-fixture", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: provider.URL, AuthStyle: oauth2.AuthStyleInParams}}
	cookies := map[string]*http.Cookie{}
	do := func(method, path, body string, browser, cross bool) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		if browser {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Sec-Fetch-Site", "same-origin")
			if cross {
				r.Header.Set("Sec-Fetch-Site", "cross-site")
			}
			for _, c := range cookies {
				if !cross || c.SameSite != http.SameSiteStrictMode {
					r.AddCookie(c)
				}
			}
		} else {
			r.Header.Set("Content-Type", "application/json")
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
	var firstAgent, firstCredential string
	for i := range 2 {
		fixture, now := connectionFixture()
		token, c, private := deviceFixture(t, fixture, now)
		proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
		input, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "kid": c.Kid, "public": c.Public, "proof": proof})
		expect(t, do("POST", "/v1/connect/start", string(input), false, false), 200, "requested")
		id := hashToken(token)
		expect(t, do("GET", "/connect/"+id, "", true, false), 200, "Google로 로그인")
		w := do("POST", "/auth/google", "connection="+id, true, false)
		expect(t, w, 303, googleIssuer)
		u, _ := url.Parse(w.Header().Get("Location"))
		auth = u.Query()
		callback := "/auth/google/callback?code=fixture&state=" + auth.Get("state")
		expect(t, do("GET", callback, "", true, true), 200, "/connect/"+id)
		expect(t, do("GET", "/connect/"+id, "", true, false), 200, "/home/device-confirm", fingerprint(decodePublic(c.Public)))
		expect(t, do("POST", "/home/device-confirm", "connection="+id, true, true), 403)
		expect(t, do("POST", "/home/confirm", "connection="+id, true, false), 403)
		if i == 0 {
			s.mutateState(t, func(st *State) {
				st.Sessions[hashToken(cookies[sessionCookie].Value)].Verified = time.Now().Add(-6 * time.Minute)
			})
			expect(t, do("POST", "/home/device-confirm", "connection="+id, true, false), 422)
			s.mutateState(t, func(st *State) { st.Sessions[hashToken(cookies[sessionCookie].Value)].Verified = time.Now() })
		}
		expect(t, do("POST", "/home/device-confirm", "connection="+id, true, false), 303)
		expect(t, do("POST", "/home/device-confirm", "connection="+id, true, false), 401)
		poll, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "proof": proof})
		w = do("POST", "/v1/connect/poll", string(poll), false, false)
		expect(t, w, 200, "approved")
		var approved map[string]string
		if json.Unmarshal(w.Body.Bytes(), &approved) != nil {
			t.Fatal("poll response")
		}
		c.Owner, c.Agent = approved["owner"], approved["agent"]
		completion := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c)))
		payload, _ := json.Marshal(map[string]string{"token": token, "client": c.Client, "proof": completion})
		w = do("POST", "/v1/connect/complete", string(payload), false, false)
		expect(t, w, 200, "connected")
		var connected map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &connected)
		if i == 0 {
			firstAgent, firstCredential = connected["agent"], connected["credential"]
		} else if firstAgent == connected["agent"] || firstCredential == connected["credential"] {
			t.Fatal("clients share authority")
		}
		expect(t, do("POST", "/v1/connect/complete", string(payload), false, false), 401)
		expect(t, do("GET", callback, "", true, true), 401)
	}
	s.mutateState(t, func(st *State) {
		if len(st.Members) != 1 || len(st.Agents) != 2 || len(st.Pairs) != 0 {
			t.Fatal("identity/agent/pair isolation")
		}
	})
}

func decodePublic(value string) []byte {
	public, _ := base64.RawURLEncoding.DecodeString(value)
	return public
}
