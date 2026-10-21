// Calendar connection (SAR-SCHEDULE-001): one Google free/busy lookup per approved schedule.query, sealed refresh
// tokens that never leave this package, and the single D010 scope setting every request and fixture must obey.
package relay

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	calendarScope     = "https://www.googleapis.com/auth/calendar.freebusy" // E2: availability only, no event fields
	freeBusyEndpoint  = "https://www.googleapis.com/calendar/v3/freeBusy"   // E1
	revokeEndpoint    = "https://oauth2.googleapis.com/revoke"              // E7
	lookupLease       = time.Minute                                         // an approval whose lookup never records a result fails after this
	lookupBudget      = 8 * time.Second                                     // whole approval request, including short retries (E8 backoff)
	calendarKeyLength = 32
)

// scheduleScope is the D010 setting in one place: window from the request time, allowed units, per-friend count.
var scheduleScope = struct {
	Window          time.Duration
	Units           []int64
	Default         int64
	PerFriend       int
	PerFriendWindow time.Duration
}{14 * 24 * time.Hour, []int64{15, 30, 60}, 30, 5, 24 * time.Hour}

// Calendar is a member's connected Google calendar: only the sealed refresh token and a state, never busy data.
type Calendar struct {
	Sealed  []byte // AES-256-GCM nonce||ciphertext of the refresh token bound to the member ID; the key lives only in process memory
	State   string // "connected" or "reconnect" (provider refused the token; the member must consent again)
	Changed time.Time
}

// Lookup is the audit record of one approval's single lookup: timestamps and outcome, plus the rounded disclosure
// until the result has been sent once.
type Lookup struct {
	Started, Finished time.Time
	State             string          `json:",omitempty"` // "" running, "done", "failed"
	Disclosure        json.RawMessage `json:",omitempty"` // canonical {window, granularity_min, busy}; cleared after the result
}

type span struct{ Start, End time.Time }

// canonicalRaw re-canonicalises a stored disclosure: the JSONB column reorders keys and adds spaces on the way back.
func canonicalRaw(raw json.RawMessage) string {
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return ""
	}
	out, err := Canonical(v)
	if err != nil {
		return ""
	}
	return string(out)
}

var errReconnect = errors.New("calendar_reconnect")

func (st *State) calendar(member string) *Calendar {
	if st.Calendars == nil {
		st.Calendars = map[string]*Calendar{}
	}
	return st.Calendars[member]
}

// calendarState is the member's connection state for pages: "connected", "reconnect" or "none".
func (st *State) calendarState(owner string) string {
	id, m := st.humanMember(owner)
	if m == nil {
		return "none"
	}
	if c := st.calendar(id); c != nil {
		return c.State
	}
	return "none"
}

// sealToken binds the member ID as GCM additional data, so a sealed value moved to another member's row never opens.
func sealToken(key, plain []byte, member string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plain, []byte(member)), nil
}

func openToken(key, sealed []byte, member string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(sealed) < gcm.NonceSize() {
		return nil, errors.New("sealed token too short")
	}
	return gcm.Open(nil, sealed[:gcm.NonceSize()], sealed[gcm.NonceSize():], []byte(member))
}

// scheduleAllowed enforces the D010 scope on an accepted schedule.query before any lookup can happen (SC-04).
func (st *State) scheduleAllowed(e *Envelope, now time.Time) error {
	minutes := int64(e.Body["granularity_min"].(float64))
	unit := false
	for _, u := range scheduleScope.Units {
		unit = unit || u == minutes
	}
	window := e.Body["window"].(map[string]any)
	start, _ := instant(window["start"])
	end, _ := instant(window["end"])
	gran := time.Duration(minutes) * time.Minute
	if !unit || start.Before(now.Truncate(gran)) || end.After(now.Add(scheduleScope.Window)) {
		return fault("schedule_scope")
	}
	from, to := st.Agents[e.From], st.Agents[e.To]
	count := 0
	for _, m := range st.Messages {
		a, b := st.Agents[m.Receipt.From], st.Agents[m.Receipt.To]
		if m.Receipt.Intent == "schedule.query" && a != nil && b != nil && a.Owner == from.Owner && b.Owner == to.Owner && now.Sub(m.Receipt.Accepted) < scheduleScope.PerFriendWindow {
			count++
		}
	}
	if count >= scheduleScope.PerFriend {
		return fault("schedule_scope")
	}
	return nil
}

// roundBusy widens each busy span outward to the unit, clips it to the window and merges overlaps (SD-04 A).
func roundBusy(busy []span, start, end time.Time, gran time.Duration) []span {
	rounded := []span{}
	for _, b := range busy {
		s, e := b.Start.Truncate(gran), b.End.Truncate(gran)
		if e.Before(b.End) {
			e = e.Add(gran)
		}
		if s.Before(start) {
			s = start
		}
		if e.After(end) {
			e = end
		}
		if e.After(s) {
			rounded = append(rounded, span{s, e})
		}
	}
	sort.Slice(rounded, func(i, j int) bool { return rounded[i].Start.Before(rounded[j].Start) })
	merged := []span{}
	for _, r := range rounded {
		if n := len(merged); n > 0 && !r.Start.After(merged[n-1].End) {
			if r.End.After(merged[n-1].End) {
				merged[n-1].End = r.End
			}
			continue
		}
		merged = append(merged, r)
	}
	return merged
}

func disclosureJSON(start, end time.Time, minutes int64, busy []span) json.RawMessage {
	const layout = "2006-01-02T15:04:05Z"
	list := []any{}
	for _, b := range busy {
		list = append(list, map[string]any{"start": b.Start.UTC().Format(layout), "end": b.End.UTC().Format(layout)})
	}
	raw, _ := Canonical(map[string]any{"window": map[string]any{"start": start.UTC().Format(layout), "end": end.UTC().Format(layout)}, "granularity_min": minutes, "busy": list})
	return raw
}

// lookupJob is what the approval transaction hands to the provider call that runs outside the state lock.
type lookupJob struct {
	gate, member string
	sealed       []byte
	start, end   time.Time
	minutes      int64
}

// startLookup records the one lookup an approval may run; the gate stays invisible to the agent until it ends.
func (st *State) startLookup(g *Gate, e *Envelope, now time.Time) (*lookupJob, error) {
	if e.Intent != "schedule.query" {
		return nil, nil
	}
	id, m := st.humanMember(g.Owner)
	c := st.calendar(id)
	if m == nil || c == nil || c.State != "connected" {
		return nil, fault("invalid_gate")
	}
	window := e.Body["window"].(map[string]any)
	start, _ := instant(window["start"])
	end, _ := instant(window["end"])
	g.Lookup = &Lookup{Started: now}
	return &lookupJob{g.ID, id, c.Sealed, start, end, int64(e.Body["granularity_min"].(float64))}, nil
}

// freeBusy refreshes the access token and queries the primary calendar once per attempt; only start/end pairs survive.
func (g *GoogleLogin) freeBusy(ctx context.Context, refresh string, start, end time.Time) ([]span, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, g.client)
	token, err := g.config.TokenSource(ctx, &oauth2.Token{RefreshToken: refresh}).Token()
	if err != nil {
		var retrieve *oauth2.RetrieveError
		if errors.As(err, &retrieve) && (retrieve.Response.StatusCode == 400 || retrieve.Response.StatusCode == 401) {
			return nil, errReconnect
		}
		return nil, err
	}
	body, _ := json.Marshal(map[string]any{"timeMin": start.UTC().Format(time.RFC3339), "timeMax": end.UTC().Format(time.RFC3339), "items": []map[string]string{{"id": "primary"}}})
	req, err := http.NewRequestWithContext(ctx, "POST", g.freeBusyURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	token.SetAuthHeader(req)
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, errReconnect
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("free/busy status " + resp.Status)
	}
	var result struct {
		Calendars map[string]struct {
			Busy   []struct{ Start, End string }
			Errors []json.RawMessage
		}
	}
	if json.Unmarshal(raw, &result) != nil || len(result.Calendars) != 1 {
		return nil, errors.New("free/busy response shape")
	}
	busy := []span{}
	for _, c := range result.Calendars {
		if len(c.Errors) > 0 {
			return nil, errors.New("free/busy calendar error")
		}
		for _, b := range c.Busy {
			s, e1 := time.Parse(time.RFC3339, b.Start)
			e, e2 := time.Parse(time.RFC3339, b.End)
			if e1 != nil || e2 != nil {
				return nil, errors.New("free/busy time shape")
			}
			busy = append(busy, span{s, e})
		}
	}
	return busy, nil
}

// runLookup performs the approval's single lookup outside the lock and records exactly one outcome.
// Transient provider failures retry briefly inside the budget; a dead process leaves the lease to expire as failed.
func (s *Service) runLookup(ctx context.Context, job *lookupJob) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), lookupBudget)
	defer cancel()
	var busy []span
	err := errors.New("calendar not configured")
	if s.Google != nil && s.Google.calendarKey != nil {
		var refresh []byte
		if refresh, err = openToken(s.Google.calendarKey, job.sealed, job.member); err == nil {
			for attempt := 0; ; attempt++ {
				busy, err = s.Google.freeBusy(ctx, string(refresh), job.start, job.end)
				if err == nil || errors.Is(err, errReconnect) || attempt >= 2 || ctx.Err() != nil {
					break
				}
				select {
				case <-time.After(500 * time.Millisecond):
				case <-ctx.Done():
				}
			}
		}
	}
	if err != nil {
		log.Print("calendar lookup failed") // no identifiers or provider detail in logs
	}
	_, txErr := s.transaction(context.WithoutCancel(ctx), func(st *State, now time.Time) (any, error) {
		g := st.Gates[job.gate]
		if g == nil || g.Lookup == nil || g.Lookup.State != "" {
			return nil, nil
		}
		g.Lookup.Finished, g.Lookup.State = now, "failed"
		if errors.Is(err, errReconnect) {
			if c := st.calendar(job.member); c != nil {
				c.State, c.Changed = "reconnect", now
			}
		}
		// A gate revoked or expired while the provider answered discloses nothing (SC-06).
		if err == nil && g.State == "approved" && now.Before(g.Exp) {
			g.Lookup.State = "done"
			g.Lookup.Disclosure = disclosureJSON(job.start, job.end, job.minutes, roundBusy(busy, job.start, job.end, time.Duration(job.minutes)*time.Minute))
		}
		return nil, nil
	})
	if txErr != nil {
		log.Print("calendar lookup outcome not recorded")
	}
}

// lookupVisible hides an approved gate from the agent until its lookup has an outcome.
func lookupVisible(g *Gate) string {
	if g.State == "approved" && g.Lookup != nil && g.Lookup.State == "" {
		return "pending"
	}
	return g.State
}

// sweepLookups ends leases a crashed process never finished and drops disclosures that can no longer be sent.
func (st *State) sweepLookups(now time.Time) {
	for _, g := range st.Gates {
		if g.Lookup == nil {
			continue
		}
		if g.Lookup.State == "" && !now.Before(g.Lookup.Started.Add(lookupLease)) {
			g.Lookup.State, g.Lookup.Finished = "failed", now
		}
		parent := st.Messages[g.Parent]
		if g.State != "approved" || parent == nil || parent.Completion != "" || !now.Before(parent.Receipt.Exp) {
			g.Lookup.Disclosure = nil
		}
	}
}

// disconnectCalendar drops the sealed token and ends every pending or not-yet-looked-up schedule.query approval as denied.
func (st *State) disconnectCalendar(member string, now time.Time) []byte {
	c := st.calendar(member)
	if c == nil {
		return nil
	}
	delete(st.Calendars, member)
	m := st.Members[member]
	for _, g := range st.Gates {
		intent, _, _ := strings.Cut(g.Policy, ":")
		if g.Owner != m.Owner || intent != "schedule.query" {
			continue
		}
		if g.State == "pending" || (g.State == "approved" && (g.Lookup == nil || g.Lookup.State != "done")) {
			g.State = "denied"
			if g.Lookup != nil {
				g.Lookup.State, g.Lookup.Disclosure = "failed", nil
			}
			if h := st.Messages[g.ID]; h != nil {
				h.Receipt.State, h.Persisted, h.Envelope, h.Inbox = "delivered", true, nil, nil
			}
		}
	}
	return c.Sealed
}

// revoke asks the provider to end the grant; failure is logged without detail and never blocks the local removal.
func (s *Service) revoke(ctx context.Context, member string, sealed []byte) {
	if s.Google == nil || s.Google.calendarKey == nil || sealed == nil {
		return
	}
	refresh, err := openToken(s.Google.calendarKey, sealed, member)
	if err != nil {
		return
	}
	s.revokeToken(ctx, string(refresh))
}

// revokeToken posts one refresh token to the provider's revoke endpoint.
func (s *Service) revokeToken(ctx context.Context, refresh string) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", s.Google.revokeURL, strings.NewReader(url.Values{"token": {refresh}}.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.Google.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	if err != nil || resp.StatusCode != 200 {
		log.Print("calendar revoke not confirmed by provider")
	}
}

func (s *Service) calendarRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /home/calendar/connect", func(w http.ResponseWriter, r *http.Request) { s.googleStart(w, r, true, true) })
	mux.HandleFunc("POST /home/calendar/disconnect", s.calendarDisconnect)
}

func (s *Service) calendarDisconnect(w http.ResponseWriter, r *http.Request) {
	var member string
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, retry, err := s.memberHit(w, st, r, now, cleanupRate)
		if !retry.IsZero() {
			return refusal{status: 429, problem: limited(retry)}, nil
		}
		if err != nil {
			return nil, err
		}
		member = id
		return st.disconnectCalendar(id, now), nil
	})
	if err != nil {
		if err.Error() == "invalid_auth" {
			expired(w, r)
			return
		}
		refusedTo(w, "설정", refusal{status: 503, problem: "캘린더 연결을 해제하지 못했습니다. 잠시 뒤 다시 시도하세요."}, "/home/account#calendar")
		return
	}
	if v, ok := value.(refusal); ok {
		refusedTo(w, "설정", v, "/home/account#calendar")
		return
	}
	s.revoke(r.Context(), member, value.([]byte))
	http.Redirect(w, r, "/home/account?n=calendar-off#calendar", http.StatusSeeOther)
}
