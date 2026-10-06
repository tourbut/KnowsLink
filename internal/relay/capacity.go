// Shared message and HTTP capacity keep new work bounded while reserving independent cleanup admission.
package relay

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

type admission struct {
	Clean bool
	Exp   time.Time
}

type requestBudgetKey struct{}

func (s *Service) requestHit(st *State, r *http.Request, now time.Time, buckets ...bucket) (bool, time.Time) {
	if r.Context().Value(requestBudgetKey{}) == true {
		return true, time.Time{}
	}
	return st.hit(now, buckets...)
}

// Classify the authenticated request once, before schema, CSRF or route rejection, without trusting the request's claimed agent ID.
func (s *Service) requestBuckets(st *State, r *http.Request, now time.Time) []bucket {
	principal := ""
	switch {
	case strings.HasPrefix(r.URL.Path, "/home") || r.URL.Path == "/auth/reauth" || r.URL.Path == "/auth/logout" || r.URL.Path == "/auth/logout-all":
		principal, _, _ = st.session(readCookie(r, sessionCookie), now)
	case strings.HasPrefix(r.URL.Path, "/v1/connect/"):
		var c struct{ Token, Client, Kid, Public, Proof string }
		if Strict(httpBody(r), &c) == nil {
			if conn, err := st.connection(c.Token, now); err == nil {
				for id, m := range st.Members {
					if m.Owner == conn.Owner {
						principal = id
						break
					}
				}
			}
		}
	case strings.HasPrefix(r.URL.Path, "/v1/"):
		var err error
		principal, err = st.principal(bearer(r), "agent")
		if err != nil {
			principal, _ = st.principal(ownerToken(r), "owner")
		}
	case strings.HasPrefix(r.URL.Path, "/owner"):
		principal, _ = st.principal(ownerToken(r), "owner")
	}
	if principal == "" {
		return anonymousRate(s.clientIP(r))
	}
	if cleanupRequest(r) {
		return cleanupRate(principal)
	}
	return memberRate(principal)
}

type requestBodyKey struct{}

func httpBody(r *http.Request) []byte {
	raw, _ := r.Context().Value(requestBodyKey{}).([]byte)
	return raw
}

func (st *State) messageCapacity(gate bool) bool {
	queue, pending := 0, 0
	for _, m := range st.Messages {
		if m.Receipt.State == "queued" || m.Receipt.State == "leased" {
			queue++
		}
	}
	for _, g := range st.Gates {
		if g.State == "pending" {
			pending++
		}
	}
	return len(st.Messages) < 20000 && queue < 100 && (!gate || pending < 100)
}
func (st *State) claimCapacity() bool {
	n := 0
	for _, m := range st.Messages {
		if m.ClaimToken != "" {
			n++
		}
	}
	return n < 4
}
func cleanupRequest(r *http.Request) bool {
	if r.Method != "POST" {
		return false
	}
	switch r.URL.Path {
	case "/v1/ack", "/v1/test/ack", "/v1/text/ack", "/v1/key-revoke", "/v1/unpair", "/v1/owner-revoke", "/home/key-revoke", "/home/agent-revoke", "/home/unpair", "/home/cancel", "/auth/logout", "/auth/logout-all", "/home/invite-deny":
		return true
	}
	if r.URL.Path == "/home/invite-decision" || strings.HasPrefix(r.URL.Path, "/home/gates/") || strings.HasPrefix(r.URL.Path, "/owner/gates/") {
		return strings.HasSuffix(r.URL.Path, "/deny") || r.PostForm.Get("decision") == "deny"
	}
	return false
}
func (st *State) enterHTTP(token string, clean bool, now time.Time) bool {
	n, limit := 0, 16
	if clean {
		limit = 4
	}
	for _, a := range st.HTTP {
		if a.Clean == clean && now.Before(a.Exp) {
			n++
		}
	}
	if n >= limit {
		return false
	}
	st.HTTP[token] = admission{clean, now.Add(30 * time.Second)}
	return true
}

// Local channels also bound requests waiting for DB admission; shared records enforce the limit across relay processes.
// Request contexts and socket body reads finish within 10s; 30s admission expiry recovers crashed processes conservatively.
func (s *Service) boundedHTTP(next http.Handler) http.Handler {
	s.httpOnce.Do(func() { s.httpNew = make(chan struct{}, 16); s.httpClean = make(chan struct{}, 4) })
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		r = r.WithContext(ctx)
		_ = http.NewResponseController(w).SetReadDeadline(time.Now().Add(10 * time.Second))
		// Explicit deny paths need no body reads to reserve cleanup. Legacy decision forms start in new-work admission.
		limit := int64(rawEnvelopeLimit)
		if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") || strings.HasPrefix(r.URL.Path, "/v1/connect/") {
			limit = 8192
		}
		r.Body = http.MaxBytesReader(w, r.Body, limit)
		clean := cleanupRequest(r)
		queue := s.httpNew
		if clean {
			queue = s.httpClean
		}
		select {
		case queue <- struct{}{}:
			defer func() { <-queue }()
		default:
			capacityResponse(w, r)
			return
		}
		// Read only after local admission. Even malformed bodies then spend the authenticated/anonymous request budget.
		raw, readErr := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(raw))
		r = r.WithContext(context.WithValue(r.Context(), requestBodyKey{}, raw))
		if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			_ = r.ParseForm()
			r.Body = io.NopCloser(bytes.NewReader(raw))
		}
		clean = cleanupRequest(r)
		token := randomToken()
		v, err := s.transaction(ctx, func(st *State, now time.Time) (any, error) {
			if ok, retry := st.hit(now, s.requestBuckets(st, r, now)...); !ok {
				return retry, nil
			}
			return st.enterHTTP(token, clean, now), nil
		})
		if err != nil {
			respond(w, nil, err)
			return
		}
		if retry, ok := v.(time.Time); ok {
			if strings.HasPrefix(r.URL.Path, "/v1/") {
				rateLimited(w, retry)
			} else {
				refused(w, "요청 제한", refusal{status: 429, problem: limited(retry)})
			}
			return
		}
		if !v.(bool) {
			capacityResponse(w, r)
			return
		}
		defer func() {
			finishCtx, done := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
			defer done()
			_, _ = s.transaction(finishCtx, func(st *State, _ time.Time) (any, error) { delete(st.HTTP, token); return nil, nil })
		}()
		if readErr != nil {
			respond(w, nil, fault("invalid_json"))
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), requestBudgetKey{}, true))
		next.ServeHTTP(w, r)
	})
}
func capacityResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", "1")
	if strings.HasPrefix(r.URL.Path, "/v1/") {
		writeJSON(w, 429, map[string]string{"error": "capacity", "retry_at": time.Now().Add(time.Second).UTC().Format(time.RFC3339)})
		return
	}
	refused(w, "요청 대기", refusal{status: 429, problem: "동시 처리 한도에 도달했습니다. 1초 뒤 수동으로 다시 시도하세요. 성공 처리되지 않았습니다."})
}
