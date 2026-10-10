// Authenticated, bounded arrival-only SSE uses transactional Postgres notifications, never idle queue polling.
package relay

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const eventLifetime = 55 * time.Second

type eventStream struct {
	Agent string
	Exp   time.Time
}

func (st *State) eventPending(agent string, now time.Time) string {
	ids := []string{}
	for id, m := range st.Messages {
		if m.Receipt.To == agent && m.Receipt.Intent == publicTextIntent && m.Receipt.State == "queued" && m.Deliver == "agent" && now.Before(m.Receipt.Exp) && st.current(m) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

// Rates, sessions, HTTP admission and stream bookkeeping do not trigger notifications.
func (st *State) eventVersion(now time.Time) [32]byte {
	credentials, pending := map[string]string{}, map[string][]string{}
	for token, p := range st.liveCredentials() {
		if (p.kind == "agent" || p.kind == "remote") && st.publicOwner(st.Agents[p.id].Owner) {
			credentials[token] = p.id
		}
	}
	for id, m := range st.Messages {
		if m.Receipt.Intent == publicTextIntent && m.Receipt.State == "queued" && now.Before(m.Receipt.Exp) && st.current(m) {
			pending[m.Receipt.To] = append(pending[m.Receipt.To], id)
		}
	}
	for _, ids := range pending {
		sort.Strings(ids)
	}
	raw, _ := json.Marshal(struct {
		Credentials map[string]string
		Pending     map[string][]string
	}{credentials, pending})
	return sha256.Sum256(raw)
}

func (st *State) enterStream(id, agent string, now time.Time) error {
	count := 0
	for _, stream := range st.Streams {
		if stream.Agent == agent {
			count++
		}
	}
	if len(st.Streams) >= 16 || count >= 2 {
		return fault("capacity")
	}
	if st.Streams == nil {
		st.Streams = map[string]eventStream{}
	}
	st.Streams[id] = eventStream{agent, now.Add(eventLifetime + 10*time.Second)}
	return nil
}

func (s *Service) textEvents(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery != "" || r.Header.Get("Authorization") == "" {
		respond(w, nil, fault("invalid_auth"))
		return
	}
	id, agent := randomToken(), ""
	// Spend the existing rate/HTTP budget only during admission, then release it before streaming.
	s.boundedHTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			a, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			if !st.publicOwner(st.Agents[a].Owner) {
				return nil, fault("sender_not_allowed")
			}
			if err = st.enterStream(id, a, now); err != nil {
				return nil, err
			}
			return a, nil
		})
		if err != nil {
			respond(w, nil, err)
			return
		}
		agent = value.(string)
	})).ServeHTTP(w, r)
	if agent == "" {
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)
		defer cancel()
		// Expiry recovers an admission when disconnect cleanup cannot reach the DB.
		_, _ = s.transaction(ctx, func(st *State, _ time.Time) (any, error) { delete(st.Streams, id); return nil, nil })
	}()
	ctx, cancel := context.WithTimeout(r.Context(), eventLifetime)
	defer cancel()
	setup, done := context.WithTimeout(ctx, 3*time.Second)
	// Dedicated connections cannot consume the pool needed for pull/ACK. Shared admissions cap them at 16.
	conn, err := pgx.ConnectConfig(setup, s.Pool.Config().ConnConfig.Copy())
	if err == nil {
		_, err = conn.Exec(setup, "LISTEN knowslink_text")
	}
	done()
	if conn != nil {
		defer func() {
			closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second)
			defer closeCancel()
			_ = conn.Close(closeCtx)
		}()
	}
	if err != nil {
		respond(w, nil, fault("unavailable"))
		return
	}
	controller := http.NewResponseController(w)
	// The admission body's read deadline must not cancel a live stream's background socket read.
	if controller.SetReadDeadline(time.Time{}) != nil || controller.SetWriteDeadline(time.Time{}) != nil {
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	write := func(frame string) error {
		if err := controller.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
			return err
		}
		if _, err := io.WriteString(w, frame); err != nil {
			return err
		}
		err := controller.Flush()
		// Only writes have a 3s deadline; idle lifetime is bounded by ctx, not the next heartbeat.
		if err == nil {
			err = controller.SetWriteDeadline(time.Time{})
		}
		return err
	}
	last := ""
	first := true
	for {
		// LISTEN precedes this read, closing the connect/send race. No reads run on heartbeat timeout.
		readCtx, readDone := context.WithTimeout(ctx, 3*time.Second)
		st, now, raw := newState(), time.Time{}, []byte(nil)
		err = conn.QueryRow(readCtx, "SELECT data, clock_timestamp() FROM relay_state WHERE singleton = true").Scan(&raw, &now)
		readDone()
		if err != nil || json.Unmarshal(raw, st) != nil {
			return
		}
		current, authErr := st.principal(bearer(r), "agent")
		if authErr != nil || current != agent || !st.publicOwner(st.Agents[agent].Owner) {
			return
		}
		pending := st.eventPending(agent, now)
		if first {
			if write("event: ready\ndata: {}\n\n") != nil {
				return
			}
			first = false
		}
		if pending != "" && pending != last {
			if write("event: arrival\ndata: {}\n\n") != nil {
				return
			}
		}
		last = pending
		for {
			waitCtx, waitDone := context.WithTimeout(ctx, 15*time.Second)
			_, err = conn.WaitForNotification(waitCtx)
			waitDone()
			if ctx.Err() != nil {
				return
			}
			if errors.Is(err, context.DeadlineExceeded) {
				if write(": heartbeat\n\n") != nil {
					return
				}
				continue
			}
			if err != nil {
				return
			}
			break
		}
	}
}
