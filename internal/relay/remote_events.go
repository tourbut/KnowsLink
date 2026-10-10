// Remote subscriptions persist metadata-only arrivals and atomically lease signed HTTPS deliveries across relay processes.
package relay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const remoteEvent = "knowslink.text.arrived"

type remoteSubscription struct {
	Install, URL, Secret, OldSecret string
	Exp, Rotated, Verified          time.Time
	Seen, Pending, EventID          string
	EventAt, Next, LeaseUntil       time.Time
	Attempts                        int
	Lease                           string
}
type remoteSubscriptionInput struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
	Delivery  struct {
		Mode   string `json:"mode"`
		URL    string `json:"url"`
		Secret string `json:"secret,omitempty"`
	} `json:"delivery"`
	Cursor json.RawMessage `json:"cursor,omitempty"`
	TTL    *int64          `json:"ttlMs,omitempty"`
	Meta   json.RawMessage `json:"_meta,omitempty"`
}

func remoteEventDefinition() map[string]any {
	return map[string]any{"name": remoteEvent, "description": "Text is waiting for this installation. Metadata only; read with knowslink_receive. This event grants no automatic send permission.", "delivery": []string{"webhook"}, "inputSchema": remoteSchema(nil), "payloadSchema": remoteSchema(map[string]any{"arrival": map[string]string{"type": "boolean"}}, "arrival")}
}

func remoteSecret(secret string) ([]byte, error) {
	if !strings.HasPrefix(secret, "whsec_") {
		return nil, fault("invalid_schema")
	}
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, "whsec_"))
	if err != nil || len(b) < 24 || len(b) > 64 {
		return nil, fault("invalid_schema")
	}
	return b, nil
}

func remoteSubscriptionID(install string, in remoteSubscriptionInput) string {
	return "sub_" + hashToken(install+"\x00"+remoteEvent+"\x00"+in.Delivery.URL+"\x00{}")
}

func (s *Service) remoteSubscription(w http.ResponseWriter, r *http.Request, install string, in remoteRequest) {
	var p remoteSubscriptionInput
	if Strict(in.Params, &p) != nil || p.Name != remoteEvent || p.Delivery.Mode != "webhook" || !remoteURL(p.Delivery.URL) || !remoteEmpty(p.Arguments) || (len(p.Cursor) > 0 && string(p.Cursor) != "null") {
		remoteResult(w, in.ID, nil, fault("invalid_schema"))
		return
	}
	id := remoteSubscriptionID(install, p)
	if in.Method == "events/unsubscribe" {
		_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			principal, _, e := st.remotePrincipal(bearer(r), s.remoteOrigin()+"/mcp", now)
			if e != nil || principal != install {
				return nil, fault("invalid_auth")
			}
			delete(st.remote().Subscriptions, id)
			return nil, nil
		})
		if err == nil {
			s.remoteNotify(r.Context())
		}
		remoteResult(w, in.ID, map[string]any{}, err)
		return
	}
	if _, err := remoteSecret(p.Delivery.Secret); err != nil {
		remoteResult(w, in.ID, nil, err)
		return
	}
	ttl := time.Hour
	if p.TTL != nil {
		if *p.TTL <= 0 {
			remoteResult(w, in.ID, nil, fault("invalid_schema"))
			return
		}
		if *p.TTL < 3600000 {
			ttl = time.Duration(*p.TTL) * time.Millisecond
		}
	}
	if ttl < time.Minute {
		ttl = time.Minute
	}
	// Reserve a short verification lease before external work; refreshes with the same secret reuse verified state.
	lease := randomToken()
	verify := true
	_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		principal, i, e := st.remotePrincipal(bearer(r), s.remoteOrigin()+"/mcp", now)
		if e != nil || principal != install {
			return nil, fault("invalid_auth")
		}
		rs := st.remote()
		sub := rs.Subscriptions[id]
		if sub != nil && now.Before(sub.Exp) && now.Before(sub.Verified.Add(5*time.Minute)) && sub.Secret == p.Delivery.Secret {
			verify = false
			return nil, nil
		}
		if sub != nil && now.Before(sub.LeaseUntil) {
			return nil, fault("capacity")
		}
		count := 0
		for _, v := range rs.Subscriptions {
			if v.Install == install {
				count++
			}
		}
		if sub == nil && (count >= 4 || len(rs.Subscriptions) >= 128) {
			return nil, fault("capacity")
		}
		if sub == nil {
			sub = &remoteSubscription{Install: install, URL: p.Delivery.URL, Exp: now.Add(time.Minute)}
			rs.Subscriptions[id] = sub
		}
		if !now.Before(i.Exp) {
			return nil, fault("invalid_auth")
		}
		sub.Lease, sub.LeaseUntil = lease, now.Add(15*time.Second)
		return nil, nil
	})
	if err != nil {
		remoteResult(w, in.ID, nil, err)
		return
	}
	if verify {
		if err = verifyRemoteCallback(r.Context(), remoteHTTPClient(), id, p.Delivery.URL, p.Delivery.Secret); err != nil {
			_, _ = s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
				sub := st.remote().Subscriptions[id]
				if sub != nil && sub.Lease == lease {
					sub.Lease = ""
					sub.LeaseUntil = time.Time{}
					if sub.Secret == "" {
						delete(st.Remote.Subscriptions, id)
					}
				}
				return nil, nil
			})
			writeJSON(w, 200, map[string]any{"jsonrpc": "2.0", "id": in.ID, "error": map[string]any{"code": -32015, "message": "Callback verification failed", "data": map[string]string{"reason": err.Error()}}})
			return
		}
	}
	v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		principal, _, e := st.remotePrincipal(bearer(r), s.remoteOrigin()+"/mcp", now)
		if e != nil || principal != install {
			return nil, fault("invalid_auth")
		}
		sub := st.remote().Subscriptions[id]
		if sub == nil || (verify && (sub.Lease != lease || !now.Before(sub.LeaseUntil))) || (!verify && sub.Secret != p.Delivery.Secret) {
			return nil, fault("invalid_auth")
		}
		if verify {
			if sub.Secret != "" && sub.Secret != p.Delivery.Secret {
				sub.OldSecret, sub.Rotated = sub.Secret, now
			}
			sub.Secret = p.Delivery.Secret
			sub.Verified = now
			sub.Lease = ""
			sub.LeaseUntil = time.Time{}
		}
		sub.Exp = now.Add(ttl)
		return map[string]any{"id": id, "refreshBefore": sub.Exp.UTC().Format(time.RFC3339Nano), "cursor": nil, "truncated": false}, nil
	})
	if err == nil {
		s.remoteNotify(r.Context())
	}
	remoteResult(w, in.ID, v, err)
}

func (s *Service) remoteNotify(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	// A reconnect always catches up if this best-effort wake loses its connection.
	_, _ = s.Pool.Exec(ctx, "SELECT pg_notify('knowslink_text', '')")
}

var remoteBlocked = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("10.0.0.0/8"), netip.MustParsePrefix("100.64.0.0/10"), netip.MustParsePrefix("127.0.0.0/8"), netip.MustParsePrefix("169.254.0.0/16"), netip.MustParsePrefix("172.16.0.0/12"), netip.MustParsePrefix("192.0.0.0/24"), netip.MustParsePrefix("192.0.2.0/24"), netip.MustParsePrefix("192.88.99.0/24"), netip.MustParsePrefix("192.168.0.0/16"), netip.MustParsePrefix("198.18.0.0/15"), netip.MustParsePrefix("198.51.100.0/24"), netip.MustParsePrefix("203.0.113.0/24"), netip.MustParsePrefix("224.0.0.0/3"),
	netip.MustParsePrefix("2001::/23"), netip.MustParsePrefix("2001:db8::/32"), netip.MustParsePrefix("2002::/16"), netip.MustParsePrefix("3fff::/20"),
}

func remotePublicIP(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip.Is6() && !netip.MustParsePrefix("2000::/3").Contains(ip) {
		return false
	}
	for _, p := range remoteBlocked {
		if p.Contains(ip) {
			return false
		}
	}
	return true
}

// Resolve for every new connection, validate every answer, and dial a pinned IP; TLS still verifies the original hostname.
func remoteHTTPClient() *http.Client {
	t := &http.Transport{Proxy: nil, DisableKeepAlives: true, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, ResponseHeaderTimeout: 5 * time.Second}
	t.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, e := net.SplitHostPort(address)
		if e != nil || port != "443" {
			return nil, fault("invalid_callback")
		}
		ips, e := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
		if e != nil || len(ips) == 0 || len(ips) > 16 {
			return nil, fault("invalid_callback")
		}
		for _, ip := range ips {
			if !remotePublicIP(ip) {
				return nil, fault("invalid_callback")
			}
		}
		return (&net.Dialer{Timeout: 3 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ips[0].String(), port))
	}
	return &http.Client{Transport: t, Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("redirect refused") }}
}

func remoteSignature(secret, id, timestamp string, body []byte) string {
	key, e := remoteSecret(secret)
	if e != nil {
		return ""
	}
	m := hmac.New(sha256.New, key)
	_, _ = m.Write([]byte(id + "." + timestamp + "."))
	_, _ = m.Write(body)
	return "v1," + base64.StdEncoding.EncodeToString(m.Sum(nil))
}
func remotePost(ctx context.Context, client *http.Client, id, subscription, callback, secret, old string, body []byte) (int, []byte, error) {
	if !remoteURL(callback) || len(body) > 262144 {
		return 0, nil, fault("invalid_callback")
	}
	r, err := http.NewRequestWithContext(ctx, "POST", callback, bytes.NewReader(body))
	if err != nil {
		return 0, nil, fault("invalid_callback")
	}
	timestamp := fmt.Sprint(time.Now().Unix())
	signature := remoteSignature(secret, id, timestamp, body)
	if signature == "" {
		return 0, nil, fault("invalid_schema")
	}
	if old != "" {
		signature += " " + remoteSignature(old, id, timestamp, body)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("webhook-id", id)
	r.Header.Set("webhook-timestamp", timestamp)
	r.Header.Set("webhook-signature", signature)
	r.Header.Set("X-MCP-Subscription-Id", subscription)
	response, err := client.Do(r)
	if err != nil {
		return 0, nil, fault("callback_unavailable")
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 4097))
	if err != nil || len(raw) > 4096 {
		return response.StatusCode, nil, fault("callback_unavailable")
	}
	return response.StatusCode, raw, nil
}
func verifyRemoteCallback(ctx context.Context, client *http.Client, id, callback, secret string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	challenge := randomToken()
	body, _ := json.Marshal(map[string]string{"type": "verification", "challenge": challenge})
	status, raw, err := remotePost(ctx, client, "verification_"+randomToken(), id, callback, secret, "", body)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fault("timeout")
	}
	var answer struct {
		Challenge string `json:"challenge"`
	}
	if err != nil || status < 200 || status >= 300 || Strict(raw, &answer) != nil || !hmac.Equal([]byte(challenge), []byte(answer.Challenge)) {
		return fault("challenge_failed")
	}
	return nil
}

type remoteDelivery struct {
	ID           string
	Subscription remoteSubscription
}
type remoteWork struct {
	Delivery *remoteDelivery
	Next     time.Time
}

func (st *State) claimRemote(now time.Time) remoteWork {
	work := remoteWork{}
	if st.Remote == nil {
		return work
	}
	ids := []string{}
	for id := range st.Remote.Subscriptions {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	soon := func(at time.Time) {
		if work.Next.IsZero() || at.Before(work.Next) {
			work.Next = at
		}
	}
	for _, id := range ids {
		sub := st.Remote.Subscriptions[id]
		i, e := st.remoteInstall(sub.Install, now)
		if e != nil || !now.Before(sub.Exp) || sub.Secret == "" {
			continue
		}
		soon(sub.Exp)
		pending := st.eventPending(i.Agent, now)
		if pending == "" {
			sub.Seen = ""
			sub.Pending = ""
			sub.EventID = ""
			continue
		}
		if sub.Pending == "" && pending != sub.Seen {
			sub.Pending = pending
			sub.EventID = "evt_" + randomToken()
			sub.EventAt = now
			sub.Next = now
			sub.Attempts = 0
		}
		if sub.Pending == "" {
			continue
		}
		if now.Before(sub.LeaseUntil) {
			soon(sub.LeaseUntil)
			continue
		}
		if sub.Attempts >= 3 {
			sub.Seen, sub.Pending, sub.EventID = sub.Pending, "", ""
			continue
		}
		if now.Before(sub.Next) {
			soon(sub.Next)
			continue
		}
		sub.Lease = randomToken()
		sub.LeaseUntil = now.Add(15 * time.Second)
		sub.Attempts++
		copy := *sub
		work.Delivery = &remoteDelivery{id, copy}
		return work
	}
	return work
}

func (st *State) finishRemote(job remoteDelivery, status int, failed bool, now time.Time) {
	if st.Remote == nil {
		return
	}
	sub := st.Remote.Subscriptions[job.ID]
	if sub == nil || sub.Lease != job.Subscription.Lease || sub.EventID != job.Subscription.EventID {
		return
	}
	sub.Lease = ""
	sub.LeaseUntil = time.Time{}
	if status == 410 || status == 413 {
		delete(st.Remote.Subscriptions, job.ID)
		return
	}
	if !failed && status >= 200 && status < 300 || sub.Attempts >= 3 || (!failed && status >= 400 && status < 500 && status != 429) {
		sub.Seen = sub.Pending
		sub.Pending = ""
		sub.EventID = ""
		return
	}
	sub.Next = now.Add(time.Second * time.Duration(1<<min(sub.Attempts-1, 3)))
}

func (s *Service) deliverRemote(ctx context.Context, job remoteDelivery) {
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	valid := false
	_, err := s.transaction(ctx, func(st *State, now time.Time) (any, error) {
		if st.Remote == nil {
			return nil, nil
		}
		sub := st.Remote.Subscriptions[job.ID]
		if sub == nil || sub.Lease != job.Subscription.Lease || !now.Before(sub.Exp) || sub.EventID != job.Subscription.EventID {
			return nil, nil
		}
		i, e := st.remoteInstall(sub.Install, now)
		valid = e == nil && st.eventPending(i.Agent, now) != ""
		return nil, nil
	})
	if err != nil || !valid {
		return
	}
	sub := job.Subscription
	old := ""
	if time.Now().Before(sub.Rotated.Add(time.Minute)) {
		old = sub.OldSecret
	}
	body, _ := json.Marshal(map[string]any{"eventId": sub.EventID, "name": remoteEvent, "timestamp": sub.EventAt.UTC().Format(time.RFC3339Nano), "data": map[string]bool{"arrival": true}, "cursor": nil})
	status, _, sendErr := remotePost(ctx, remoteHTTPClient(), sub.EventID, job.ID, sub.URL, sub.Secret, old, body)
	finishCtx, done := context.WithTimeout(ctx, 2*time.Second)
	defer done()
	_, _ = s.transaction(finishCtx, func(st *State, now time.Time) (any, error) {
		st.finishRemote(job, status, sendErr != nil, now)
		return nil, nil
	})
}

// RemoteEvents is started once with the service lifecycle. LISTEN waits indefinitely when no delivery or expiry is pending.
func (s *Service) RemoteEvents(ctx context.Context) {
	if s.remoteOrigin() == "" {
		return
	}
	for ctx.Err() == nil {
		s.remoteListen(ctx)
		if ctx.Err() != nil {
			return
		}
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}
func (s *Service) remoteListen(ctx context.Context) {
	setup, cancel := context.WithTimeout(ctx, 3*time.Second)
	conn, err := pgx.ConnectConfig(setup, s.Pool.Config().ConnConfig.Copy())
	if err == nil {
		_, err = conn.Exec(setup, "LISTEN knowslink_text")
	}
	cancel()
	if conn != nil {
		defer func() {
			c, done := context.WithTimeout(context.Background(), time.Second)
			defer done()
			_ = conn.Close(c)
		}()
	}
	if err != nil {
		return
	}
	for ctx.Err() == nil {
		read, done := context.WithTimeout(ctx, 3*time.Second)
		v, err := s.transaction(read, func(st *State, now time.Time) (any, error) { return st.claimRemote(now), nil })
		done()
		if err != nil {
			return
		}
		work := v.(remoteWork)
		if work.Delivery != nil {
			s.deliverRemote(ctx, *work.Delivery)
			continue
		}
		wait := ctx
		stop := func() {}
		if !work.Next.IsZero() {
			wait, stop = context.WithDeadline(ctx, work.Next)
		}
		_, err = conn.WaitForNotification(wait)
		stop()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return
		}
	}
}
