//go:build integration && landingfixture

// Opt-in loopback fixture serves real handlers and a synthetic DB; mail and Google cancellation stay local.
package relay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

func TestLandingBrowserFixture(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("LANDING_FIXTURE_MANIFEST") == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, mail := identityService(t, pool)
	stop := make(chan struct{})
	var once sync.Once
	manifest := map[string]string{}
	for _, mode := range []string{"google", "both", "email", "none"} {
		s := &Service{Pool: pool}
		if mode == "both" || mode == "email" {
			s.Mail = mail.send
		}
		if mode == "both" || mode == "google" {
			s.Google = &GoogleLogin{config: oauth2.Config{ClientID: "synthetic", Endpoint: oauth2.Endpoint{AuthURL: "/fixture/google-cancel"}}}
		}
		handler := s.Handler()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			if strings.HasPrefix(r.URL.Path, "/fixture/") {
				if r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Header.Get("Origin") != "" && r.Header.Get("Origin") != "http://"+r.Host) {
					http.Error(w, "fixture origin refused", 403)
					return
				}
				switch r.Method + " " + r.URL.Path {
				case "GET /fixture/code":
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					_, _ = w.Write([]byte(mail.code("landing@example.test")))
				case "GET /fixture/google-cancel":
					http.Redirect(w, r, "/auth/google/callback?error=access_denied&state="+r.URL.Query().Get("state"), 303)
				case "POST /fixture/expire":
					_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
						if session := st.Sessions[hashToken(readCookie(r, sessionCookie))]; session != nil {
							session.Created = now.Add(-sessionAbsolute - time.Minute)
						}
						return nil, nil
					})
					if err != nil {
						http.Error(w, "fixture transaction failed", 503)
						return
					}
					http.Redirect(w, r, "/", 303)
				case "POST /fixture/stop":
					w.WriteHeader(204)
					once.Do(func() { close(stop) })
				default:
					http.NotFound(w, r)
				}
				return
			}
			allowed := map[string]bool{
				"GET /": true, "GET /auth/verify": true, "GET /auth/google/callback": true,
				"GET /home": true, "GET /home/name": true, "GET /home/account": true,
				"POST /auth/start": true, "POST /auth/verify": true, "POST /auth/google": true,
				"POST /auth/logout": true, "POST /home/profile": true,
			}
			if !allowed[r.Method+" "+r.URL.Path] {
				http.Error(w, "outside landing fixture scope", 403)
				return
			}
			handler.ServeHTTP(w, r)
		}))
		defer server.Close()
		manifest[mode] = server.URL
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(os.Getenv("LANDING_FIXTURE_MANIFEST"), data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("synthetic landing fixture ready; no external auth, mail or messages")
	select {
	case <-stop:
	case <-time.After(24 * time.Hour):
	}
}
