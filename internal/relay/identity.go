// Email identity binds one verified address to one member and owner; browser sessions never carry agent authority.
package relay

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

// Product defaults (D02 SAR-PUBLIC-SERVICE operation defaults).
const (
	codeTTL         = 10 * time.Minute
	codeAttempts    = 5
	sessionAbsolute = 12 * time.Hour
	sessionIdle     = 60 * time.Minute
	reauthWindow    = 5 * time.Minute
	activeMembers   = 100
	emailIssuer     = "knowslink-email-otp"
	sessionCookie   = "__Host-kl_session"
	pendingCookie   = "__Host-kl_pending"
)

type Member struct {
	Owner, Email, Issuer string
	Active               bool
	Created              time.Time
}
type Session struct {
	Member                  string
	Created, Seen, Verified time.Time
}
type Challenge struct {
	Email, Code string
	Exp         time.Time
	Attempts    int
}

// Mailer sends one plain-text message; tests inject a capture and production uses SMTP.
type Mailer func(ctx context.Context, to, subject, body string) error

type bucket struct {
	key    string
	limit  int
	window time.Duration
}

func within(hits []time.Time, window time.Duration, now time.Time) (n int, oldest time.Time) {
	for _, t := range hits {
		if t.After(now.Add(-window)) {
			if n == 0 || t.Before(oldest) {
				oldest = t
			}
			n++
		}
	}
	return n, oldest
}

// take records one use in every bucket only when all allow it; rolling window is (now-window, now].
func (st *State) take(now time.Time, buckets ...bucket) (bool, time.Time) {
	ok, retry := true, time.Time{}
	for _, b := range buckets {
		if n, oldest := within(st.Rates[b.key], b.window, now); n >= b.limit {
			ok = false
			if at := oldest.Add(b.window); at.After(retry) {
				retry = at
			}
		}
	}
	if ok {
		for _, b := range buckets {
			st.Rates[b.key] = append(st.Rates[b.key], now)
		}
	}
	return ok, retry
}

// hit counts every request, including rejected ones, as the HTTP rate rows require.
// Callers list the principal bucket first. A request its principal refuses never spends the shared budget.
// A request a later bucket refuses is recorded only where earlier buckets already hold hits in the window,
// so rotating sources cannot open new keys or keep a full shared budget saturated by themselves.
// It stops at the first full bucket and keeps at most limit+1 hits, so a flood cannot grow state without bound.
func (st *State) hit(now time.Time, buckets ...bucket) (bool, time.Time) {
	kept := make([][]time.Time, len(buckets))
	for i, b := range buckets {
		for _, t := range st.Rates[b.key] {
			if t.After(now.Add(-b.window)) {
				kept[i] = append(kept[i], t)
			}
		}
		if len(kept[i]) < b.limit {
			continue
		}
		for _, earlier := range kept[:i] {
			if len(earlier) == 0 {
				return false, kept[i][len(kept[i])-b.limit].Add(b.window)
			}
		}
		for j := range i {
			st.Rates[buckets[j].key] = append(kept[j], now)
		}
		full := append(kept[i], now)
		full = full[len(full)-b.limit-1:]
		st.Rates[b.key] = full
		// The next request is counted too, so two of the kept hits must leave the window.
		return false, full[1].Add(b.window)
	}
	for i, b := range buckets {
		st.Rates[b.key] = append(kept[i], now)
	}
	return true, time.Time{}
}

func anonymousRate(ip string) []bucket {
	return []bucket{{"http:ip:" + ip, 30, time.Minute}, {"http:new", 200, time.Minute}}
}
func memberRate(member string) []bucket {
	return []bucket{{"http:member:" + member, 40, time.Minute}, {"http:new", 200, time.Minute}}
}
func cleanupRate(member string) []bucket {
	return []bucket{{"cleanup:member:" + member, 20, time.Minute}, {"cleanup", 100, time.Minute}}
}

// NormalizeEmail accepts one bare address and lowercases it. Plus/dot aliases stay distinct identities.
func NormalizeEmail(value string) (string, bool) {
	value = strings.TrimSpace(value)
	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Name != "" || !strings.EqualFold(addr.Address, value) || len(value) > 254 {
		return "", false
	}
	at := strings.LastIndex(value, "@")
	domain := value[at+1:]
	if at < 1 || strings.HasPrefix(domain, "[") || !strings.Contains(domain, ".") || strings.HasSuffix(domain, ".") {
		return "", false
	}
	return strings.ToLower(value), true
}

func maskEmail(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 1 {
		return "***"
	}
	return email[:1] + "***" + email[at:]
}

func sixDigits() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *Service) clientIP(r *http.Request) string {
	if s.ClientIPHeader != "" {
		if ip := net.ParseIP(strings.TrimSpace(r.Header.Get(s.ClientIPHeader))); ip != nil {
			return ip.String()
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// startChallenge reserves send budget and replaces earlier codes for the address. The caller sends outside the lock.
func (st *State) startChallenge(email, ip string, now time.Time) (pending, code string, retry time.Time, ok bool) {
	if ok, retry := st.take(now,
		bucket{"send:space:" + hashToken(email), 1, time.Minute},
		bucket{"send:email:" + hashToken(email), 5, time.Hour},
		bucket{"send:ip:" + ip, 20, time.Hour},
		bucket{"send", 100, time.Hour}); !ok {
		return "", "", retry, false
	}
	for id, c := range st.Challenges {
		if c.Email == email {
			delete(st.Challenges, id)
		}
	}
	pending, code = randomToken(), sixDigits()
	st.Challenges[hashToken(pending)] = &Challenge{email, hashToken(pending + ":" + code), now.Add(codeTTL), 0}
	return pending, code, time.Time{}, true
}

// verify outcomes are returned as values so failed attempts still commit their counters.
type verifyResult struct {
	Session, Problem string
	Remaining        int
}

func (st *State) verifyChallenge(pending, code, current string, now time.Time) verifyResult {
	id := hashToken(pending)
	c := st.Challenges[id]
	if pending == "" || c == nil || !now.Before(c.Exp) || c.Attempts >= codeAttempts {
		return verifyResult{Problem: "expired"}
	}
	if !hmac.Equal([]byte(c.Code), []byte(hashToken(pending+":"+code))) {
		c.Attempts++
		if c.Attempts >= codeAttempts {
			delete(st.Challenges, id)
			return verifyResult{Problem: "expired"}
		}
		return verifyResult{Problem: "wrong", Remaining: codeAttempts - c.Attempts}
	}
	delete(st.Challenges, id)
	key := emailIssuer + "|" + c.Email
	member := st.Identities[key]
	if member == "" {
		active := 0
		for _, m := range st.Members {
			if m.Active {
				active++
			}
		}
		if active >= activeMembers {
			return verifyResult{Problem: "capacity"}
		}
		member = "mem_" + randomToken()[:22]
		owner := "owner_" + randomToken()
		// Member owners have no bearer credential; only a verified browser session reaches them.
		st.Owners[owner] = &Owner{"", true}
		st.Members[member] = &Member{owner, c.Email, emailIssuer, true, now}
		st.Identities[key] = member
	}
	if m := st.Members[member]; !m.Active || st.Owners[m.Owner] == nil || !st.Owners[m.Owner].Active {
		return verifyResult{Problem: "inactive"}
	}
	if current != "" {
		delete(st.Sessions, hashToken(current))
	}
	token := randomToken()
	st.Sessions[hashToken(token)] = &Session{member, now, now, now}
	return verifyResult{Session: token}
}

// session validates absolute and idle lifetime and slides the idle deadline.
func (st *State) session(token string, now time.Time) (string, *Session, error) {
	id := hashToken(token)
	s := st.Sessions[id]
	if token == "" || s == nil {
		return "", nil, fault("invalid_auth")
	}
	m := st.Members[s.Member]
	if !now.Before(s.Created.Add(sessionAbsolute)) || !now.Before(s.Seen.Add(sessionIdle)) || m == nil || !m.Active || st.Owners[m.Owner] == nil || !st.Owners[m.Owner].Active {
		delete(st.Sessions, id)
		return "", nil, fault("invalid_auth")
	}
	s.Seen = now
	return s.Member, s, nil
}

func (st *State) sweepIdentity(now time.Time) {
	for id, c := range st.Challenges {
		if !now.Before(c.Exp) {
			delete(st.Challenges, id)
		}
	}
	for id, s := range st.Sessions {
		if !now.Before(s.Created.Add(sessionAbsolute)) || !now.Before(s.Seen.Add(sessionIdle)) {
			delete(st.Sessions, id)
		}
	}
	for key, hits := range st.Rates {
		kept := hits[:0]
		for _, t := range hits {
			if t.After(now.Add(-time.Hour)) {
				kept = append(kept, t)
			}
		}
		if len(kept) == 0 {
			delete(st.Rates, key)
		} else {
			st.Rates[key] = kept
		}
	}
}
