// Google checks use real RSA/JWK verification and deterministic member/session state, without external credentials.
package relay

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-jose/go-jose/v4"
)

func googleFixture(t *testing.T) (*GoogleLogin, func(map[string]any) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", "test-key"))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "test-key", Algorithm: "RS256", Use: "sig"}}})
	}))
	t.Cleanup(server.Close)
	return &GoogleLogin{verifier: oidc.NewVerifier(googleIssuer, oidc.NewRemoteKeySet(context.Background(), server.URL), &oidc.Config{ClientID: "client", SupportedSigningAlgs: []string{oidc.RS256}})}, func(claims map[string]any) string {
		t.Helper()
		raw, err := json.Marshal(claims)
		if err != nil {
			t.Fatal(err)
		}
		signed, err := signer.Sign(raw)
		if err != nil {
			t.Fatal(err)
		}
		token, err := signed.CompactSerialize()
		if err != nil {
			t.Fatal(err)
		}
		return token
	}
}

func googleClaims(now time.Time, nonce string) map[string]any {
	return map[string]any{"iss": googleIssuer, "sub": "google-subject", "aud": "client", "iat": now.Unix(), "exp": now.Add(time.Hour).Unix(), "nonce": nonce, "email": "member@example.com", "email_verified": true}
}

func TestGoogleTokenVerification(t *testing.T) {
	g, sign := googleFixture(t)
	now := time.Now().UTC()
	a := &googleAttempt{Nonce: "nonce", Created: now}
	for _, field := range []string{"valid", "iss", "aud", "exp", "nonce", "email_verified", "email", "sub", "iat", "future", "stale", "unsigned", "signature"} {
		t.Run(field, func(t *testing.T) {
			c := googleClaims(now, a.Nonce)
			switch field {
			case "iss":
				c[field] = "https://attacker.example"
			case "aud":
				c[field] = "another-client"
			case "exp":
				c[field] = now.Add(-time.Minute).Unix()
			case "nonce":
				c[field] = "other"
			case "email_verified":
				c[field] = false
			case "email":
				c[field] = "invalid"
			case "sub":
				c[field] = ""
			case "iat":
				delete(c, field)
			case "future":
				c["iat"] = now.Add(2 * time.Minute).Unix()
			case "stale":
				c["iat"] = now.Add(-6 * time.Minute).Unix()
			}
			raw := sign(c)
			if field == "unsigned" {
				raw = "eyJhbGciOiJub25lIn0.e30."
			}
			if field == "signature" {
				raw = raw[:len(raw)-10] + "AAAAAAAAAA"
			}
			subject, email, err := g.verify(context.Background(), raw, a)
			if field == "valid" {
				if err != nil || subject != "google-subject" || email != "member@example.com" {
					t.Fatal("valid identity refused", err)
				}
			} else if err == nil {
				t.Fatal("invalid identity accepted", field)
			}
		})
	}
}

func TestGoogleIdentityContinuityAndReauth(t *testing.T) {
	st, now := newState(), time.Now().UTC()
	email := "member@example.com"
	old := st.signIn(emailIssuer, email, email, "", "", now)
	first := st.signIn(googleIssuer, "subject", email, hashToken(old.Session), "", now)
	id := st.Sessions[hashToken(first.Session)].Member
	owner := st.Members[id].Owner
	if len(st.Members) != 2 || len(st.Owners) != 2 || st.Sessions[hashToken(old.Session)] != nil || st.Owners[owner].Credential != "" {
		t.Fatal("email merge or stale session authority")
	}
	st.Agents["kept"] = &Agent{Owner: owner}
	second := st.signIn(googleIssuer, "subject", "changed@example.com", hashToken(first.Session), id, now.Add(time.Minute))
	if second.Session == "" || len(st.Members) != 2 || st.Sessions[hashToken(second.Session)].Member != id || st.Sessions[hashToken(first.Session)] != nil || st.Agents["kept"].Owner != owner {
		t.Fatal("continuity or rotation")
	}
	for _, subject := range []string{"stranger", "subject"} {
		hash := hashToken(second.Session)
		if subject == "subject" {
			delete(st.Sessions, hash)
		}
		if r := st.signIn(googleIssuer, subject, email, hash, id, now.Add(2*time.Minute)); r.Session != "" || len(st.Members) != 2 {
			t.Fatal("foreign identity or logged-out session reauthenticated")
		}
	}
	for i := len(st.Members); i < activeMembers; i++ {
		st.Members[randomToken()] = &Member{Active: true}
	}
	if r := st.signIn(googleIssuer, "new", email, "", "", now); r.Problem != "capacity" {
		t.Fatal("capacity")
	}
	if r := st.signIn(googleIssuer, "subject", email, "", "", now); r.Session == "" {
		t.Fatal("existing login blocked")
	}
	st.Owners[owner].Active = false
	if r := st.signIn(googleIssuer, "subject", email, "", "", now); r.Problem != "inactive" {
		t.Fatal("inactive owner")
	}
}

func TestGoogleAttemptAndConfiguration(t *testing.T) {
	st, now, state := newState(), time.Now(), randomToken()
	st.GoogleAttempts[hashToken(state)] = &googleAttempt{Exp: now.Add(time.Minute)}
	if _, err := st.consumeGoogleAttempt(state, "foreign-browser", now); err == nil {
		t.Fatal("state without browser cookie")
	}
	if _, err := st.consumeGoogleAttempt(state, state, now); err != nil {
		t.Fatal(err)
	}
	if _, err := st.consumeGoogleAttempt(state, state, now); err == nil {
		t.Fatal("state replay")
	}
	st.GoogleAttempts[hashToken(state)] = &googleAttempt{Exp: now}
	if _, err := st.consumeGoogleAttempt(state, state, now); err == nil {
		t.Fatal("expired attempt")
	}
	st.GoogleAttempts[hashToken(state)] = &googleAttempt{Exp: now}
	st.sweepIdentity(now)
	if len(st.GoogleAttempts) != 0 {
		t.Fatal("expired state retained")
	}
	if g, err := ConfigureGoogle(context.Background(), "", "", ""); g != nil || err != nil {
		t.Fatal("disabled configuration")
	}
	for _, u := range []string{"", "http://localhost/auth/google/callback", "https://example.com/wrong", "https://user@example.com/auth/google/callback", "https://example.com/auth/google/callback?x=1"} {
		if _, err := ConfigureGoogle(context.Background(), "client", "secret", u); err == nil {
			t.Fatal("invalid redirect accepted")
		}
	}
	w := httptest.NewRecorder()
	googleAttemptCookie(w, state, 600)
	cookie := w.Result().Cookies()[0]
	if !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" {
		t.Fatal("weak attempt cookie")
	}
}
