// Google login verifies OIDC codes and binds one-use browser attempts to existing member sessions.
package relay

import (
	"context"
	"crypto/hmac"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const googleIssuer = "https://accounts.google.com"
const googleCookie = "__Host-kl_google"

type GoogleLogin struct {
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
	client   *http.Client
}

type googleAttempt struct {
	Nonce, Verifier, CurrentHash, Member string
	Connection                           string
	Created, Exp                         time.Time
}

// ConfigureGoogle leaves login disabled when all values are absent; partial or unsafe settings fail closed.
func ConfigureGoogle(ctx context.Context, clientID, secret, redirect string) (*GoogleLogin, error) {
	if clientID == "" && secret == "" && redirect == "" {
		return nil, nil
	}
	u, err := url.Parse(redirect)
	if clientID == "" || secret == "" || err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "/auth/google/callback" || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("Google login requires client ID, client secret and an HTTPS /auth/google/callback redirect URL")
	}
	client := &http.Client{Timeout: 5 * time.Second}
	ctx = oidc.ClientContext(ctx, client)
	provider, err := oidc.NewProvider(ctx, googleIssuer)
	if err != nil {
		return nil, errors.New("Google discovery unavailable")
	}
	return &GoogleLogin{
		config:   oauth2.Config{ClientID: clientID, ClientSecret: secret, RedirectURL: redirect, Endpoint: provider.Endpoint(), Scopes: []string{oidc.ScopeOpenID, "email"}},
		verifier: provider.VerifierContext(context.WithoutCancel(ctx), &oidc.Config{ClientID: clientID, SupportedSigningAlgs: []string{oidc.RS256}}),
		client:   client,
	}, nil
}

func googleAttemptCookie(w http.ResponseWriter, value string, maxAge int) {
	// Only this temporary cookie is Lax so Google's top-level GET can return. Session cookie stays Strict.
	http.SetCookie(w, &http.Cookie{Name: googleCookie, Value: value, Path: "/", MaxAge: maxAge, Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

func (s *Service) googleFailure(w http.ResponseWriter, status int, problem string) {
	v := startView(problem)
	v["Google"], v["EmailLogin"] = s.Google != nil, s.Mail != nil
	render(w, status, "start", v)
}

func (s *Service) googleStart(w http.ResponseWriter, r *http.Request, reauth bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if r.ParseForm() != nil {
		s.googleFailure(w, 422, "로그인을 시작할 수 없습니다. 다시 시작하세요.")
		return
	}
	if s.Google == nil {
		s.googleFailure(w, 503, "Google 로그인이 아직 설정되지 않았습니다. 잠시 뒤 다시 시도하세요.")
		return
	}
	state := randomToken()
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, retry, sessionErr := s.memberHit(st, r, now, memberRate)
		if !retry.IsZero() {
			return nil, fault("rate_limited")
		}
		if reauth && (sessionErr != nil || st.Members[id].Issuer != googleIssuer) {
			return nil, fault("invalid_auth")
		}
		a := &googleAttempt{Nonce: randomToken(), Verifier: oauth2.GenerateVerifier(), CurrentHash: hashToken(readCookie(r, sessionCookie)), Created: now, Exp: now.Add(codeTTL)}
		if reauth {
			a.Member = id
		}
		if connection := r.FormValue("connection"); connection != "" {
			c := st.Connections[connection]
			if c == nil || !c.Device || !now.Before(c.Exp) || (c.State != "requested" && c.State != "prepared") {
				return nil, fault("invalid_auth")
			}
			if c.Owner != "" {
				if sessionErr != nil || st.Members[id].Owner != c.Owner || st.Members[id].Issuer != googleIssuer {
					return nil, fault("invalid_auth")
				}
				a.Member = id
			}
			a.Connection = connection
		}
		if st.GoogleAttempts == nil {
			st.GoogleAttempts = map[string]*googleAttempt{}
		}
		delete(st.GoogleAttempts, hashToken(readCookie(r, googleCookie)))
		st.GoogleAttempts[hashToken(state)] = a
		return a, nil
	})
	if err != nil {
		s.googleFailure(w, statusFor(err), "로그인을 시작할 수 없습니다. 현재 로그인과 요청 한도를 확인한 뒤 다시 시도하세요.")
		return
	}
	a := value.(*googleAttempt)
	googleAttemptCookie(w, state, int(codeTTL/time.Second))
	// An explicit account selection/consent plus a fresh nonce-bound ID token re-verifies identity, not the Google password.
	redirect := s.Google.config.AuthCodeURL(state, oidc.Nonce(a.Nonce), oauth2.S256ChallengeOption(a.Verifier), oauth2.SetAuthURLParam("prompt", "select_account consent"))
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (st *State) consumeGoogleAttempt(state, cookie string, now time.Time) (*googleAttempt, error) {
	if len(state) != 43 || !hmac.Equal([]byte(state), []byte(cookie)) {
		return nil, fault("invalid_auth")
	}
	a := st.GoogleAttempts[hashToken(state)]
	delete(st.GoogleAttempts, hashToken(state))
	if a == nil || !now.Before(a.Exp) {
		return nil, fault("invalid_auth")
	}
	return a, nil
}

func (g *GoogleLogin) verify(ctx context.Context, raw string, a *googleAttempt) (subject, email string, err error) {
	token, err := g.verifier.Verify(ctx, raw)
	if err != nil {
		return "", "", fault("invalid_auth")
	}
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	now := time.Now()
	if token.Claims(&claims) != nil || !claims.Verified || token.Subject == "" || len(token.Subject) > 255 || token.Nonce != a.Nonce ||
		token.IssuedAt.Before(a.Created.Add(-time.Minute)) || token.IssuedAt.After(now.Add(time.Minute)) || !now.Before(token.IssuedAt.Add(reauthWindow)) {
		return "", "", fault("invalid_auth")
	}
	email, valid := NormalizeEmail(claims.Email)
	if !valid {
		return "", "", fault("invalid_auth")
	}
	return token.Subject, email, nil
}

func (s *Service) googleCallback(w http.ResponseWriter, r *http.Request) {
	if s.Google == nil {
		s.googleFailure(w, 503, "Google 로그인이 아직 설정되지 않았습니다.")
		return
	}
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		// Return refusal as a value so consumption commits even for expiry/cancellation.
		a, _ := st.consumeGoogleAttempt(r.URL.Query().Get("state"), readCookie(r, googleCookie), now)
		return a, nil
	})
	googleAttemptCookie(w, "", -1)
	if err != nil {
		s.googleFailure(w, 503, "로그인을 처리할 수 없습니다. 다시 시작하세요.")
		return
	}
	a := value.(*googleAttempt)
	if a == nil || r.URL.Query().Get("error") != "" || r.URL.Query().Get("code") == "" {
		s.googleFailure(w, 401, "Google 로그인이 취소되었거나 요청이 만료되었습니다. 로그인되지 않았습니다. 다시 시작하세요.")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
	defer cancel()
	ctx = oidc.ClientContext(ctx, s.Google.client)
	exchanged, err := s.Google.config.Exchange(ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(a.Verifier))
	if err != nil {
		s.googleFailure(w, 401, "Google 인증에 실패했습니다. 다시 시작하세요.")
		return
	}
	raw, _ := exchanged.Extra("id_token").(string)
	subject, email, err := s.Google.verify(ctx, raw, a)
	if err != nil {
		s.googleFailure(w, 401, "Google 신원을 확인하지 못했습니다. 다시 시작하세요.")
		return
	}
	value, err = s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		return st.finishGoogleLogin(a, subject, email, now), nil
	})
	if err != nil {
		s.googleFailure(w, 503, "로그인을 완료하지 못했습니다. 다시 시작하세요.")
		return
	}
	result := value.(verifyResult)
	if result.Problem != "" {
		s.googleFailure(w, 403, "이 계정으로 로그인하거나 재확인할 수 없습니다. 기존 회원의 같은 Google 계정으로 다시 시도하세요. 신규 가입 한도도 적용됩니다.")
		return
	}
	setCookie(w, sessionCookie, result.Session, int(sessionAbsolute/time.Second))
	// A same-site click sends the Strict session; a redirect chain from Google may not.
	next := "/home"
	if a.Connection != "" {
		next = "/connect/" + a.Connection
	}
	render(w, 200, "google-complete", map[string]any{"Title": "Google 로그인 확인 완료", "Notice": "신원을 확인했습니다. 계속해서 연결을 확인하세요.", "Next": next})
}

// finishGoogleLogin binds a verified identity before browser consent; a different account cannot replace that binding.
func (st *State) finishGoogleLogin(a *googleAttempt, subject, email string, now time.Time) verifyResult {
	if !now.Before(a.Exp) {
		return verifyResult{Problem: "expired"}
	}
	if a.Connection != "" {
		c := st.Connections[a.Connection]
		if c == nil || !c.Device || !now.Before(c.Exp) || (c.State != "requested" && c.State != "prepared") {
			return verifyResult{Problem: "expired"}
		}
		// A competing account cannot bind an already claimed browser request or gain a session through it.
		if c.Owner != "" {
			m := st.Members[st.Identities[googleIssuer+"|"+subject]]
			if m == nil || m.Owner != c.Owner {
				return verifyResult{Problem: "invalid_auth"}
			}
		}
	}
	result := st.signIn(googleIssuer, subject, email, a.CurrentHash, a.Member, now)
	if result.Problem == "" && a.Connection != "" {
		id, _, e := st.session(result.Session, now)
		if e != nil {
			return verifyResult{Problem: "invalid_auth"}
		}
		c := st.Connections[a.Connection]
		c.Owner, c.State = st.Members[id].Owner, "prepared"
	}
	return result
}
