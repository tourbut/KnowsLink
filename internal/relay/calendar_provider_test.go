//go:build integration

// Synthetic Google provider for calendar checks: local OAuth code/refresh/revoke and a free/busy endpoint whose
// answers carry deliberately private extra fields that must never leave the server. No real account is involved.
package relay

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

const fixtureAccountCookie = "kl_fixture_account"

type syntheticGrant struct{ nonce, subject, scope, challenge, redirect string }

type syntheticProvider struct {
	t    *testing.T
	sign func(map[string]any) string
	key  []byte
	mu   sync.Mutex
	// auths keeps the relay's authorization requests by state until a synthetic consent decides them.
	auths    map[string]url.Values
	codes    map[string]syntheticGrant
	refresh  map[string]string // refresh token -> subject
	revoked  map[string]bool
	accesses int
	// Operator toggles for failure classes: provider-side revocation (invalid_grant) and a free/busy outage (503).
	grantRevoked, outage bool
	freeBusyCalls        int
	lastMin, lastMax     string
	busy                 func(start, end time.Time) []map[string]any
}

func newSyntheticProvider(t *testing.T, sign func(map[string]any) string) *syntheticProvider {
	key := make([]byte, calendarKeyLength)
	_, _ = rand.Read(key)
	p := &syntheticProvider{t: t, sign: sign, key: key, auths: map[string]url.Values{}, codes: map[string]syntheticGrant{}, refresh: map[string]string{}, revoked: map[string]bool{}}
	// Default synthetic schedule: two short unaligned busy items near each end of the window, each with private fields.
	p.busy = func(start, end time.Time) []map[string]any {
		item := func(s, e time.Time) map[string]any {
			return map[string]any{"start": s.Format(time.RFC3339), "end": e.Format(time.RFC3339), "summary": "secret-title 비밀 회의", "attendees": []string{"secret-attendee@example.test"}, "location": "secret-location"}
		}
		return []map[string]any{item(start.Add(5*time.Minute), start.Add(10*time.Minute)), item(end.Add(-25*time.Minute), end.Add(-20*time.Minute))}
	}
	return p
}

// install points the relay's Google client at this provider mounted under base; the login scopes stay openid email.
func (p *syntheticProvider) install(g *GoogleLogin, base string, client *http.Client) {
	g.client = client
	g.calendarKey = p.key
	g.freeBusyURL, g.revokeURL = base+"/fixture/google/freebusy", base+"/fixture/google/revoke"
	g.config = oauth2.Config{ClientID: "client", ClientSecret: "local-fixture", RedirectURL: base + "/auth/google/callback", Scopes: []string{"openid", "email"},
		Endpoint: oauth2.Endpoint{AuthURL: base + "/fixture/google/auth", TokenURL: base + "/fixture/google/token", AuthStyle: oauth2.AuthStyleInParams}}
}

// grant records the relay's authorization request and returns the callback URL a consent decision would produce.
// scope "" means the person declined; a scope without calendar.freebusy is a partial grant.
func (p *syntheticProvider) grant(location, subject, scope string) string {
	u, err := url.Parse(location)
	if err != nil {
		p.t.Fatal(err)
	}
	q := u.Query()
	p.mu.Lock()
	defer p.mu.Unlock()
	if scope == "" {
		return q.Get("redirect_uri") + "?state=" + url.QueryEscape(q.Get("state")) + "&error=access_denied"
	}
	code := "code-" + randomToken()[:16]
	p.codes[code] = syntheticGrant{nonce: q.Get("nonce"), subject: subject, scope: scope, challenge: q.Get("code_challenge"), redirect: q.Get("redirect_uri")}
	return q.Get("redirect_uri") + "?state=" + url.QueryEscape(q.Get("state")) + "&code=" + code + "&scope=" + url.QueryEscape(scope)
}

// sealedFor is a stored connection for a member as the relay would have saved it after consent.
func (p *syntheticProvider) sealedFor(subject, member string) []byte {
	token := "fixture-refresh-" + randomToken()[:16]
	p.mu.Lock()
	p.refresh[token] = subject
	p.mu.Unlock()
	sealed, err := sealToken(p.key, []byte(token), member)
	if err != nil {
		p.t.Fatal(err)
	}
	return sealed
}

func (p *syntheticProvider) handler() http.Handler {
	mux := http.NewServeMux()
	writeJSON := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}
	// The consent page of the browser fixture: the signed-in fixture account decides for the subject it belongs to.
	mux.HandleFunc("GET /fixture/google/auth", func(w http.ResponseWriter, r *http.Request) {
		account, _ := r.Cookie(fixtureAccountCookie)
		if account == nil || r.URL.Query().Get("state") == "" {
			http.Error(w, "fixture account cookie and state required", 400)
			return
		}
		subject := "sub-" + account.Value
		full := p.grant(r.URL.String(), subject, "openid email "+calendarScope)
		partial := p.grant(r.URL.String(), subject, "openid email")
		denied := p.grant(r.URL.String(), subject, "")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(w, `<!doctype html><html lang="ko"><meta charset="utf-8"><title>합성 Google 동의</title><h1>합성 Google 동의 화면</h1><p>실제 Google이 아닙니다. 계정 %s, 요청 범위: <code>%s</code></p><p><a href="%s">허용 (빈 시간 조회 포함)</a></p><p><a href="%s">일부만 허용 (캘린더 제외)</a></p><p><a href="%s">거절</a></p></html>`,
			html.EscapeString(account.Value), html.EscapeString(r.URL.Query().Get("scope")), html.EscapeString(full), html.EscapeString(partial), html.EscapeString(denied))
	})
	mux.HandleFunc("POST /fixture/google/token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		p.mu.Lock()
		defer p.mu.Unlock()
		switch r.PostForm.Get("grant_type") {
		case "authorization_code":
			g, ok := p.codes[r.PostForm.Get("code")]
			delete(p.codes, r.PostForm.Get("code"))
			if !ok || oauth2.S256ChallengeFromVerifier(r.PostForm.Get("code_verifier")) != g.challenge || r.PostForm.Get("redirect_uri") != g.redirect {
				writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
				return
			}
			now := time.Now()
			claims := map[string]any{"iss": googleIssuer, "sub": g.subject, "aud": "client", "iat": now.Unix(), "exp": now.Add(time.Hour).Unix(), "nonce": g.nonce, "email": strings.TrimPrefix(g.subject, "sub-") + "@example.test", "email_verified": true}
			token := map[string]any{"access_token": "fixture-access-code", "token_type": "Bearer", "expires_in": 3600, "id_token": p.sign(claims), "scope": g.scope}
			if strings.Contains(g.scope, calendarScope) {
				refresh := "fixture-refresh-" + randomToken()[:16]
				p.refresh[refresh] = g.subject
				token["refresh_token"] = refresh
			}
			writeJSON(w, 200, token)
		case "refresh_token":
			token := r.PostForm.Get("refresh_token")
			if _, ok := p.refresh[token]; !ok || p.revoked[token] || p.grantRevoked || r.PostForm.Get("client_id") != "client" {
				writeJSON(w, 400, map[string]string{"error": "invalid_grant"})
				return
			}
			p.accesses++
			writeJSON(w, 200, map[string]any{"access_token": fmt.Sprintf("fixture-access-%d", p.accesses), "token_type": "Bearer", "expires_in": 3600, "scope": calendarScope})
		default:
			writeJSON(w, 400, map[string]string{"error": "unsupported_grant_type"})
		}
	})
	mux.HandleFunc("POST /fixture/google/freebusy", func(w http.ResponseWriter, r *http.Request) {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.freeBusyCalls++
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer fixture-access-") {
			writeJSON(w, 401, map[string]any{"error": map[string]any{"code": 401}})
			return
		}
		if p.outage {
			writeJSON(w, 503, map[string]any{"error": map[string]any{"code": 503, "message": "synthetic outage"}})
			return
		}
		var req struct {
			TimeMin, TimeMax string
			Items            []struct{ ID string }
		}
		if json.NewDecoder(r.Body).Decode(&req) != nil || len(req.Items) != 1 || req.Items[0].ID != "primary" {
			writeJSON(w, 400, map[string]any{"error": map[string]any{"code": 400}})
			return
		}
		p.lastMin, p.lastMax = req.TimeMin, req.TimeMax
		start, _ := time.Parse(time.RFC3339, req.TimeMin)
		end, _ := time.Parse(time.RFC3339, req.TimeMax)
		writeJSON(w, 200, map[string]any{"kind": "calendar#freeBusy", "timeMin": req.TimeMin, "timeMax": req.TimeMax, "calendars": map[string]any{"primary": map[string]any{"busy": p.busy(start, end)}}})
	})
	mux.HandleFunc("POST /fixture/google/revoke", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		p.mu.Lock()
		defer p.mu.Unlock()
		token := r.PostForm.Get("token")
		if _, ok := p.refresh[token]; !ok {
			writeJSON(w, 400, map[string]string{"error": "invalid_token"})
			return
		}
		p.revoked[token] = true
		w.WriteHeader(200)
	})
	// Operator toggles for the browser fixture (same origin only); integration tests set the fields directly.
	mux.HandleFunc("POST /fixture/google/toggle/{name}", func(w http.ResponseWriter, r *http.Request) {
		p.mu.Lock()
		defer p.mu.Unlock()
		switch r.PathValue("name") {
		case "revoke-grant":
			p.grantRevoked = !p.grantRevoked
		case "outage":
			p.outage = !p.outage
		default:
			http.Error(w, "unknown toggle", 404)
			return
		}
		writeJSON(w, 200, map[string]bool{"grant_revoked": p.grantRevoked, "outage": p.outage})
	})
	return mux
}

func (p *syntheticProvider) snapshot() (calls int, min, max string, revoked []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for token := range p.revoked {
		revoked = append(revoked, token)
	}
	return p.freeBusyCalls, p.lastMin, p.lastMax, revoked
}

func (p *syntheticProvider) keyBase64() string { return base64.StdEncoding.EncodeToString(p.key) }
