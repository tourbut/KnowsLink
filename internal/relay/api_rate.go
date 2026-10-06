// Agent API rates are keyed by stable agent identity; cleanup work has its own persisted budget.
package relay

import (
	"net/http"
	"strings"
	"time"
)

func (s *Service) rateAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/") || strings.HasPrefix(r.URL.Path, "/v1/connect/") {
			next.ServeHTTP(w, r)
			return
		}
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			principal, e := st.principal(bearer(r), "agent")
			if e != nil {
				principal, e = st.principal(ownerToken(r), "owner")
			}
			rates := anonymousRate(s.clientIP(r))
			if e == nil {
				rates = memberRate(principal)
				switch r.URL.Path {
				case "/v1/ack", "/v1/test/ack", "/v1/key-revoke", "/v1/unpair", "/v1/owner-revoke":
					rates = cleanupRate(principal)
				}
			}
			if ok, retry := st.hit(now, rates...); !ok {
				return retry, nil
			}
			return time.Time{}, nil
		})
		if err != nil {
			respond(w, nil, err)
			return
		}
		if retry := value.(time.Time); !retry.IsZero() {
			rateLimited(w, retry)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// rateLimited gives every JSON API the same actual retry instant, rounded up to the whole second both formats carry
// so a client retrying at it is not refused (and counted) again.
func rateLimited(w http.ResponseWriter, retry time.Time) {
	retry = retry.Add(time.Second - 1).Truncate(time.Second)
	w.Header().Set("Retry-After", retry.Format(http.TimeFormat))
	writeJSON(w, 429, map[string]string{"error": "rate_limited", "retry_at": retry.Format(time.RFC3339)})
}
