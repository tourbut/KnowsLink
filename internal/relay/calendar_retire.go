// One-time removal of the retired server calendar connection (D018, SAR-SCHEDULE-001 '기존 연결 회원 처리' 1–3·5·6).
// shortcut: the legacy State/Gate/googleAttempt fields stay only so this cleanup can read them; drop the fields and this
// file once production state has no calendar connection left.
package relay

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const calendarRevokeEndpoint = "https://oauth2.googleapis.com/revoke"

// legacyCalendar is the stored shape of a retired connection; only the sealed refresh token is read, to revoke it.
type legacyCalendar struct{ Sealed []byte }

// RetireCalendars runs before the relay serves requests. key is the optional base64 KNOWSLINK_CALENDAR_KEY that
// sealed the tokens; without a valid key every connection is deleted without a revoke request.
func (s *Service) RetireCalendars(ctx context.Context, key string) error {
	return s.retireCalendars(ctx, key, calendarRevokeEndpoint, &http.Client{Timeout: 5 * time.Second})
}

func (s *Service) retireCalendars(ctx context.Context, key, endpoint string, client *http.Client) error {
	value, err := s.transaction(ctx, func(st *State, _ time.Time) (any, error) { return st.retireCalendars(), nil })
	if err != nil {
		return errors.New("calendar cleanup failed")
	}
	sealed := value.(map[string][]byte)
	if len(sealed) == 0 {
		return nil
	}
	// Deletion is committed first: a failed or skipped revoke never leaves a token behind (SAR-SCHEDULE-001 2).
	secret, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(secret) != 32 {
		secret = nil
	}
	revoked := 0
	for member, token := range sealed {
		if secret != nil && revokeSealed(ctx, client, endpoint, secret, token, member) {
			revoked++
		}
	}
	// Counts only: no member, token or provider detail (SAR-SCHEDULE-001 5).
	log.Printf("calendar cleanup: connections=%d revoked=%d failed=%d", len(sealed), revoked, len(sealed)-revoked)
	return nil
}

// retireCalendars deletes every connection and in-flight consent, and ends each gate whose lookup result or running
// lookup the agent has not consumed as denied; a consumed result was already handed over and is not recalled.
func (st *State) retireCalendars() map[string][]byte {
	sealed := map[string][]byte{}
	for member, c := range st.Calendars {
		if c != nil {
			sealed[member] = c.Sealed
		}
	}
	st.Calendars = nil
	for _, g := range st.Gates {
		if g.Lookup == nil {
			continue
		}
		if !g.Consumed && (g.State == "pending" || g.State == "approved") {
			g.State = "denied"
			if h := st.Messages[g.ID]; h != nil {
				h.Receipt.State, h.Persisted, h.Envelope, h.Inbox = "delivered", true, nil, nil
			}
		}
		g.Lookup = nil
	}
	for id, a := range st.GoogleAttempts {
		if a.Calendar {
			delete(st.GoogleAttempts, id)
		}
	}
	return sealed
}

// revokeSealed opens a token sealed with the member ID as GCM additional data and posts it once to the revoke endpoint.
func revokeSealed(ctx context.Context, client *http.Client, endpoint string, key, sealed []byte, member string) bool {
	block, err := aes.NewCipher(key)
	if err != nil {
		return false
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(sealed) < gcm.NonceSize() {
		return false
	}
	refresh, err := gcm.Open(nil, sealed[:gcm.NonceSize()], sealed[gcm.NonceSize():], []byte(member))
	if err != nil {
		return false
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(url.Values{"token": {string(refresh)}}.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
