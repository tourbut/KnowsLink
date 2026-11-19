//go:build integration

// Independent tester QA for SAR-GOOGLE-CONNECT-001: HTTP-level trust boundaries with a signed local OAuth provider and a synthetic Postgres.
// Copied into a temporary server clone only; it is not product code and uses no real Google account or credential.
package relay

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
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

func TestTesterGoogleDeviceBoundaries(t *testing.T) {
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
	auths := map[string]url.Values{} // OAuth code -> authorization request values
	subjects := map[string]string{"A": "subject-a", "B": "subject-b", "D": "subject-d"}
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.PostFormValue("code")
		auth := auths[code]
		if auth == nil || oauth2.S256ChallengeFromVerifier(r.PostFormValue("code_verifier")) != auth.Get("code_challenge") {
			http.Error(w, "PKCE", 400)
			return
		}
		claims := googleClaims(time.Now(), auth.Get("nonce"))
		claims["sub"], claims["email"] = subjects[code], strings.ToLower(subjects[code])+"@example.com"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "unused", "token_type": "Bearer", "id_token": sign(claims)})
	}))
	defer provider.Close()
	g.client = provider.Client()
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "local-fixture", RedirectURL: "https://example.com/auth/google/callback", Scopes: []string{"openid", "email"}, Endpoint: oauth2.Endpoint{AuthURL: googleIssuer + "/auth", TokenURL: provider.URL, AuthStyle: oauth2.AuthStyleInParams}}

	type browser struct {
		ip      string
		cookies map[string]*http.Cookie
	}
	newBrowser := func(ip string) *browser { return &browser{ip, map[string]*http.Cookie{}} }
	do := func(b *browser, method, path, body string, cross bool) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		r.RemoteAddr = b.ip + ":4000"
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Sec-Fetch-Site", "same-origin")
		if cross {
			r.Header.Set("Sec-Fetch-Site", "cross-site")
		}
		for _, c := range b.cookies {
			if !cross || c.SameSite != http.SameSiteStrictMode {
				r.AddCookie(c)
			}
		}
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		for _, c := range w.Result().Cookies() {
			if c.MaxAge < 0 {
				delete(b.cookies, c.Name)
			} else {
				b.cookies[c.Name] = c
			}
		}
		return w
	}
	api := func(ip, path string, v map[string]string) *httptest.ResponseRecorder {
		t.Helper()
		raw, _ := json.Marshal(v)
		r := httptest.NewRequest("POST", "https://example.com"+path, strings.NewReader(string(raw)))
		r.RemoteAddr = ip + ":4000"
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		return w
	}
	start := func(ip string) (token, id string, c *Connection, private ed25519.PrivateKey, proof string) {
		t.Helper()
		fixture, now := connectionFixture()
		token, c, private = deviceFixture(t, fixture, now)
		proof = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
		expect(t, api(ip, "/v1/connect/start", map[string]string{"token": token, "client": c.Client, "kid": c.Kid, "public": c.Public, "proof": proof}), 200, "requested")
		return token, hashToken(token), c, private, proof
	}
	signin := func(b *browser, code string) *httptest.ResponseRecorder {
		t.Helper()
		w := do(b, "POST", "/auth/google", "", false)
		expect(t, w, 303, googleIssuer)
		u, _ := url.Parse(w.Header().Get("Location"))
		auths[code] = u.Query()
		return do(b, "GET", "/auth/google/callback?code="+code+"&state="+u.Query().Get("state"), "", true)
	}
	state := func(id string) (owner, st string) {
		s.mutateState(t, func(x *State) {
			if c := x.Connections[id]; c != nil {
				owner, st = c.Owner, c.State
			}
		})
		return
	}
	members := func() (n int) {
		s.mutateState(t, func(x *State) { n = len(x.Members) })
		return
	}

	A, B, D, anon := newBrowser("192.0.2.10"), newBrowser("192.0.2.11"), newBrowser("192.0.2.13"), newBrowser("192.0.2.12")
	token, id, c, _, proof := start("192.0.2.20")
	expect(t, signin(B, "B"), 200, "/home") // B is an existing, signed-in member before it meets the request.
	sessionB := B.cookies[sessionCookie].Value

	// S1: an existing member and a new account both begin Google login for one unclaimed request; only the first callback binds it.
	wA := do(A, "POST", "/auth/google", "connection="+id, false)
	wB := do(B, "POST", "/auth/google", "connection="+id, false)
	wD := do(D, "POST", "/auth/google", "connection="+id, false)
	expect(t, wA, 303, googleIssuer)
	expect(t, wB, 303, googleIssuer)
	expect(t, wD, 303, googleIssuer)
	uA, _ := url.Parse(wA.Header().Get("Location"))
	uB, _ := url.Parse(wB.Header().Get("Location"))
	uD, _ := url.Parse(wD.Header().Get("Location"))
	auths["A"], auths["B"], auths["D"] = uA.Query(), uB.Query(), uD.Query()
	expect(t, do(A, "GET", "/auth/google/callback?code=A&state="+uA.Query().Get("state"), "", true), 200, "/connect/"+id)
	ownerA, st := state(id)
	if ownerA == "" || st != "prepared" || members() != 2 {
		t.Fatalf("A binding owner=%q state=%q members=%d", ownerA, st, members())
	}
	w := do(B, "GET", "/auth/google/callback?code=B&state="+uB.Query().Get("state"), "", true)
	if w.Code != 403 || B.cookies[sessionCookie].Value != sessionB || members() != 2 {
		t.Fatalf("existing competing member got %d session-kept=%v members=%d", w.Code, B.cookies[sessionCookie].Value == sessionB, members())
	}
	w = do(D, "GET", "/auth/google/callback?code=D&state="+uD.Query().Get("state"), "", true)
	if w.Code != 403 || D.cookies[sessionCookie] != nil || members() != 2 {
		t.Fatalf("new competing account got %d session=%v members=%d", w.Code, D.cookies[sessionCookie] != nil, members())
	}
	if o, _ := state(id); o != ownerA {
		t.Fatal("binding replaced")
	}

	// S2: a foreign anonymous browser or a foreign signed-in member cannot view, restart, cancel, approve or confirm the bound request.
	expect(t, do(anon, "GET", "/connect/"+id, "", false), 403)
	if w := do(anon, "POST", "/auth/google", "connection="+id, false); w.Code == 303 {
		t.Fatal("foreign anonymous browser restarted a bound request")
	}
	for _, p := range []string{"/home/cancel", "/home/device-confirm", "/home/confirm"} {
		if w := do(B, "POST", p, "connection="+id, false); w.Code == 303 {
			t.Fatalf("foreign member %s accepted", p)
		}
	}
	expect(t, do(B, "GET", "/connect/"+id, "", false), 403)
	expect(t, do(B, "GET", "/home/connections/"+id, "", false), 403)
	if w := do(B, "POST", "/auth/google", "connection="+id, false); w.Code == 303 {
		t.Fatal("foreign signed-in member restarted a bound request")
	}
	if o, st := state(id); o != ownerA || st != "prepared" {
		t.Fatalf("foreign actions changed request: %q %q", o, st)
	}

	// S3: legacy token APIs cannot drive an unapproved client-started request; polling needs the private key.
	expect(t, api("192.0.2.21", "/v1/connect/info", map[string]string{"token": token, "client": c.Client}), 401)
	expect(t, api("192.0.2.21", "/v1/connect/prepare", map[string]string{"token": token, "client": c.Client, "kid": "k_other", "public": c.Public, "proof": proof}), 401)
	_, other, _ := ed25519.GenerateKey(rand.Reader)
	forged := base64.RawURLEncoding.EncodeToString(ed25519.Sign(other, deviceBytes(token, c)))
	expect(t, api("192.0.2.22", "/v1/connect/poll", map[string]string{"token": token, "client": c.Client, "proof": forged}), 401)
	expect(t, api("192.0.2.22", "/v1/connect/poll", map[string]string{"token": token, "client": "other-client", "proof": proof}), 401)
	w = api("192.0.2.22", "/v1/connect/poll", map[string]string{"token": token, "client": c.Client, "proof": proof})
	expect(t, w, 200, "prepared")
	if strings.Contains(w.Body.String(), "agent") || strings.Contains(w.Body.String(), "credential") {
		t.Fatal("poll leaked authority before consent")
	}
	if w := api("192.0.2.23", "/v1/connect/complete", map[string]string{"token": token, "client": c.Client, "proof": proof}); w.Code != 401 {
		t.Fatalf("credential before consent: %d", w.Code)
	}

	// S4: owner cancels; consent and poll are refused afterwards.
	expect(t, do(A, "POST", "/home/cancel", "connection="+id, false), 303)
	if w := do(A, "POST", "/home/device-confirm", "connection="+id, false); w.Code == 303 {
		t.Fatal("approved a cancelled request")
	}
	expect(t, api("192.0.2.22", "/v1/connect/poll", map[string]string{"token": token, "client": c.Client, "proof": proof}), 401)

	// S5: expiry refuses page, restart and poll; no agent is created.
	token2, id2, c2, _, proof2 := start("192.0.2.24")
	s.mutateState(t, func(x *State) { x.Connections[id2].Exp = time.Now().Add(-time.Second) })
	expect(t, do(A, "GET", "/connect/"+id2, "", false), 401)
	if w := do(A, "POST", "/auth/google", "connection="+id2, false); w.Code == 303 {
		t.Fatal("restart of an expired request")
	}
	expect(t, api("192.0.2.25", "/v1/connect/poll", map[string]string{"token": token2, "client": c2.Client, "proof": proof2}), 401)

	// S6: malformed connection values never produce a redirect.
	for _, bad := range []string{"https://evil.example", "../home", "x"} {
		if w := do(A, "POST", "/auth/google", "connection="+url.QueryEscape(bad), false); w.Code == 303 {
			t.Fatalf("accepted connection %q", bad)
		}
	}

	// S7: cross-site browser writes are refused for the new binding and consent entry points.
	_, id3, _, _, _ := start("192.0.2.26")
	expect(t, do(anon, "POST", "/auth/google", "connection="+id3, true), 403)
	expect(t, do(A, "POST", "/home/device-confirm", "connection="+id3, true), 403)

	// S8: no refused path created an agent or pair; only the two Google members exist.
	s.mutateState(t, func(x *State) {
		if len(x.Agents) != 0 || len(x.Pairs) != 0 || len(x.Members) != 2 {
			t.Fatalf("unexpected authority: agents=%d pairs=%d members=%d", len(x.Agents), len(x.Pairs), len(x.Members))
		}
	})
}
