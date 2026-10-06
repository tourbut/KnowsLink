// Rate refusals expose the same rounded retry instant to JSON clients; shared admission spends the stable principal budget.
package relay

import (
	"net/http"
	"time"
)

// rateLimited gives every JSON API the same actual retry instant, rounded up to the whole second both formats carry
// so a client retrying at it is not refused (and counted) again.
func rateLimited(w http.ResponseWriter, retry time.Time) {
	retry = retry.Add(time.Second - 1).Truncate(time.Second)
	w.Header().Set("Retry-After", retry.Format(http.TimeFormat))
	writeJSON(w, 429, map[string]string{"error": "rate_limited", "retry_at": retry.Format(time.RFC3339)})
}
