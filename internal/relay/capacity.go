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
	Owner string `json:",omitempty"`
}

type requestBudgetKey struct{}

func (s *Service) requestHit(st *State, r *http.Request, now time.Time, buckets ...bucket) (bool, time.Time) {
	if r.Context().Value(requestBudgetKey{}) == true {
		return true, time.Time{}
	}
	return st.hit(now, buckets...)
}

// Classify the authenticated request once, before schema, CSRF or route rejection, without trusting the request's claimed agent ID.
// Only a verified principal ending its own record spends cleanup budget; anonymous, foreign or malformed cleanup is new work.
// cleanOwner names the owner that cleanup acts for, the unit of cleanup fairness, and is "" for new work.
func (s *Service) requestBuckets(st *State, r *http.Request, now time.Time, cleanable bool) (buckets []bucket, cleanOwner string) {
	principal, owner := "", ""
	switch {
	case strings.HasPrefix(r.URL.Path, "/home") || r.URL.Path == "/auth/reauth" || r.URL.Path == "/auth/logout" || r.URL.Path == "/auth/logout-all":
		principal, _, _ = st.session(readCookie(r, sessionCookie), now)
		if principal != "" {
			owner = st.Members[principal].Owner
		}
	case strings.HasPrefix(r.URL.Path, "/v1/connect/"):
		var c struct{ Token, Client, Kid, Public, Proof string }
		if Strict(httpBody(r), &c) == nil {
			if conn, err := st.connection(c.Token, now); err == nil {
				for id, m := range st.Members {
					if m.Owner == conn.Owner {
						principal, owner = id, conn.Owner
						break
					}
				}
			}
		}
	case strings.HasPrefix(r.URL.Path, "/v1/"):
		var err error
		principal, err = st.principal(bearer(r), "agent")
		if err == nil {
			owner = st.Agents[principal].Owner
		} else {
			principal, _ = st.principal(ownerToken(r), "owner")
			owner = principal
		}
	case strings.HasPrefix(r.URL.Path, "/owner"):
		principal, _ = st.principal(ownerToken(r), "owner")
		owner = principal
	}
	if principal == "" {
		return anonymousRate(s.clientIP(r)), ""
	}
	if cleanable && st.cleanupTarget(r, principal, now) {
		return cleanupRate(principal), owner
	}
	return memberRate(principal), ""
}

type requestBodyKey struct{}

func httpBody(r *http.Request) []byte {
	raw, _ := r.Context().Value(requestBodyKey{}).([]byte)
	return raw
}

// commandBody reads an empty body as {} like the /v1 operation handlers.
func commandBody(r *http.Request) []byte {
	if raw := httpBody(r); len(raw) > 0 {
		return raw
	}
	return []byte("{}")
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
	if r.URL.Path == "/v1/invite-decision" {
		var c command
		return Strict(commandBody(r), &c) == nil && c.Decision == "deny"
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

// enterHTTP records one shared admission; cleanOwner != "" takes the cleanup budget, at most one per owner across processes.
func (st *State) enterHTTP(token, cleanOwner string, now time.Time) bool {
	clean := cleanOwner != ""
	n, limit := 0, 16
	if clean {
		limit = 4
	}
	for _, a := range st.HTTP {
		if a.Clean == clean && now.Before(a.Exp) {
			if clean && a.Owner == cleanOwner {
				return false
			}
			n++
		}
	}
	if n >= limit {
		return false
	}
	st.HTTP[token] = admission{clean, now.Add(30 * time.Second), cleanOwner}
	return true
}

// Local channels also bound requests waiting for DB admission; shared records enforce the limit across relay processes.
// Only a cleanup the last committed state proves to be the caller's own record waits on the cleanup channel, and each
// owner holds at most one cleanup slot locally and in shared records, so neither a foreign, lease-less or rate-limited
// cleanup nor one member repeating its own valid cleanup can starve other owners' cleanup. The capped body is received
// before any slot, under the same 10s deadline as the request context, so a slow or anonymous sender holds only its own
// connection, like header reception. 30s admission expiry recovers crashed processes.
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
		limit := int64(rawEnvelopeLimit)
		form := strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		if form || strings.HasPrefix(r.URL.Path, "/v1/connect/") {
			limit = 8192
		}
		raw, readErr := io.ReadAll(http.MaxBytesReader(w, r.Body, limit))
		r.Body = io.NopCloser(bytes.NewReader(raw))
		r = r.WithContext(context.WithValue(r.Context(), requestBodyKey{}, raw))
		var formErr error
		if form {
			formErr = r.ParseForm()
			r.Body = io.NopCloser(bytes.NewReader(raw))
		}
		credential := ""
		if readErr == nil && formErr == nil {
			credential = cleanupCredential(r)
		}
		queue := s.httpNew
		if owner := s.cleanupOwner(r, credential); owner != "" {
			if _, busy := s.cleaning.LoadOrStore(owner, nil); busy {
				capacityResponse(w, r)
				return
			}
			defer s.cleaning.Delete(owner)
			queue = s.httpClean
		}
		select {
		case queue <- struct{}{}:
			defer func() { <-queue }()
		default:
			capacityResponse(w, r)
			return
		}
		// Even malformed bodies spend the authenticated/anonymous request budget once they reach admission.
		token, clean := randomToken(), false
		v, err := s.transaction(ctx, func(st *State, now time.Time) (any, error) {
			buckets, owner := s.requestBuckets(st, r, now, credential != "")
			if ok, retry := st.hit(now, buckets...); !ok {
				return retry, nil
			}
			clean = owner != ""
			return st.enterHTTP(token, owner, now), nil
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
			if _, err := s.transaction(finishCtx, func(st *State, _ time.Time) (any, error) { delete(st.HTTP, token); return nil, nil }); err != nil {
				s.orphans.Store(token, nil) // the next commit reclaims it, so one owner's slot does not wait for the 30s expiry
			}
		}()
		// A live credential that no longer proves its own cleanup target hands the reserved local slot back.
		if queue == s.httpClean && !clean {
			select {
			case s.httpNew <- struct{}{}:
				<-s.httpClean
				queue = s.httpNew
			default:
				capacityResponse(w, r)
				return
			}
		}
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
