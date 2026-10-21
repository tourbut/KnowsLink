//go:build integration

// D018 cleanup on real Postgres: legacy calendar state written by the SAR-SCHEDULE-001 build is revoked once (key
// present), deleted regardless of the revoke result, and an unconsumed lookup result ends denied and undisclosed.
package relay

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRetireCalendars(t *testing.T) {
	pool := messagePool(t)
	key := bytes.Repeat([]byte{7}, 32)
	seal := func(plain, member string) string {
		block, _ := aes.NewCipher(key)
		gcm, _ := cipher.NewGCM(block)
		nonce := make([]byte, gcm.NonceSize())
		return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(plain), []byte(member)))
	}
	for _, c := range []struct {
		name, key string
		status    int
		revoked   []string
		logLine   string
	}{
		{"key_revoke_ok", base64.StdEncoding.EncodeToString(key), 200, []string{"refresh-a", "refresh-b"}, "connections=2 revoked=2 failed=0"},
		{"key_revoke_refused", base64.StdEncoding.EncodeToString(key), 500, []string{"refresh-a", "refresh-b"}, "connections=2 revoked=0 failed=2"},
		{"no_key", "", 200, nil, "connections=2 revoked=0 failed=2"},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := setup(t, pool)
			receipt := f.send(f.message(firstID, "retire-request-0001"), "agent_a", "", 200)
			claim := f.deliver(firstID)
			f.send(wire(t, f.private["agent_b"], gateID, "agent_b", "agent_b", "relay.approval.request", "retire-gate-0001", firstID, map[string]any{"reason": "judgment_required", "request_digest": receipt["digest"]}, time.Now().Add(time.Minute)), "agent_b", claim, 200)
			f.mutate(func(st *State) {
				st.Gates[gateID].State = "approved"
				st.GoogleAttempts["calendar-attempt"] = &googleAttempt{Calendar: true, Exp: time.Now().Add(time.Minute)}
				st.GoogleAttempts["login-attempt"] = &googleAttempt{Exp: time.Now().Add(time.Minute)}
			})
			// The exact JSON the SAR-SCHEDULE-001 build stored: two sealed connections and an undelivered busy result.
			calendars := `{"mem_a":{"Sealed":"` + seal("refresh-a", "mem_a") + `","State":"connected","Changed":"2026-10-21T00:00:00Z"},"mem_b":{"Sealed":"` + seal("refresh-b", "mem_b") + `","State":"reconnect","Changed":"2026-10-21T00:00:00Z"}}`
			lookup := `{"Started":"2026-10-21T00:00:00Z","Finished":"2026-10-21T00:00:01Z","State":"done","Disclosure":{"busy":[{"end":"2026-10-21T02:00:00Z","start":"2026-10-21T01:00:00Z"}],"granularity_min":30,"window":{"end":"2026-10-21T03:00:00Z","start":"2026-10-21T00:00:00Z"}}}`
			if _, err := pool.Exec(context.Background(), `UPDATE relay_state SET data=jsonb_set(jsonb_set(data,'{Calendars}',$1::jsonb),ARRAY['Gates',$2::text,'Lookup'],$3::jsonb),epoch=epoch+1`, calendars, gateID, lookup); err != nil {
				t.Fatal(err)
			}

			var mu sync.Mutex
			var posted []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = r.ParseForm()
				mu.Lock()
				posted = append(posted, r.PostForm.Get("token"))
				mu.Unlock()
				w.WriteHeader(c.status)
			}))
			defer server.Close()
			var logs bytes.Buffer
			prev := log.Writer()
			log.SetOutput(&logs)
			err := f.s.retireCalendars(context.Background(), c.key, server.URL, server.Client())
			log.SetOutput(prev)
			if err != nil {
				t.Fatal(err)
			}
			slices.Sort(posted)
			if !slices.Equal(posted, c.revoked) {
				t.Fatalf("revoke requests %v want %v", posted, c.revoked)
			}
			if !strings.Contains(logs.String(), c.logLine) || strings.Contains(logs.String(), "refresh-") || strings.Contains(logs.String(), "mem_") {
				t.Fatalf("cleanup record %q", logs.String())
			}
			var raw []byte
			if err := pool.QueryRow(context.Background(), `SELECT data::text FROM relay_state`).Scan(&raw); err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(raw, []byte("Calendars")) || bytes.Contains(raw, []byte("Lookup")) || bytes.Contains(raw, []byte("Disclosure")) || bytes.Contains(raw, []byte(`"Calendar"`)) {
				t.Fatalf("calendar state left behind: %s", raw)
			}
			f.mutate(func(st *State) {
				if st.Gates[gateID].State != "denied" || st.GoogleAttempts["calendar-attempt"] != nil || st.GoogleAttempts["login-attempt"] == nil {
					t.Fatalf("gate %s attempts %v", st.Gates[gateID].State, st.GoogleAttempts)
				}
			})
			// The undelivered result is never handed out: the denied gate cannot be consumed.
			result := f.call("POST", "/v1/gate-consume", f.tokens["agent_b"], map[string]any{"id": gateID, "claim": claim}, 409)
			if _, leaked := result["disclosure"]; leaked {
				t.Fatal("disclosure after cleanup", result)
			}
			// A second start finds nothing to clean and sends nothing.
			posted = nil
			if err := f.s.retireCalendars(context.Background(), c.key, server.URL, server.Client()); err != nil || len(posted) != 0 {
				t.Fatal("second cleanup", err, posted)
			}
		})
	}
}

// RC-01: the calendar routes are gone; a request to them reaches no handler.
func TestCalendarRoutesRetired(t *testing.T) {
	f := setup(t, messagePool(t))
	for _, path := range []string{"/home/calendar/connect", "/home/calendar/disconnect"} {
		if w := f.request("POST", path, "", nil, ""); w.Code != http.StatusNotFound {
			t.Fatal(path, w.Code)
		}
	}
}

// A consumed result was already handed over: the gate keeps its decision and only the record is dropped.
func TestRetireCalendarsKeepsConsumedGate(t *testing.T) {
	st := newState()
	st.Gates["consumed"] = &Gate{ID: "consumed", State: "approved", Consumed: true, Lookup: json.RawMessage(`{"State":"done"}`)}
	st.Gates["plain"] = &Gate{ID: "plain", State: "pending"}
	if sealed := st.retireCalendars(); len(sealed) != 0 {
		t.Fatal(sealed)
	}
	if g := st.Gates["consumed"]; g.State != "approved" || g.Lookup != nil || st.Gates["plain"].State != "pending" {
		t.Fatal("gates changed", g.State, st.Gates["plain"].State)
	}
}
