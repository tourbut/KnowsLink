// Identity unit checks cover address normalization, rolling limits, one-use codes, member continuity, sessions, and SMTP.
package relay

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNormalizeEmail(t *testing.T) {
	for in, want := range map[string]string{"User.Name+tag@Example.COM": "user.name+tag@example.com", " a@b.co ": "a@b.co"} {
		if got, ok := NormalizeEmail(in); !ok || got != want {
			t.Fatalf("%q: %q %v", in, got, ok)
		}
	}
	for _, in := range []string{"", "a", "a@b", "a@b.", "Name <a@b.co>", "a@[127.0.0.1]", "a@b.co\r\nBcc: x@y.co", "a@b.co, c@d.co", strings.Repeat("a", 250) + "@b.co"} {
		if _, ok := NormalizeEmail(in); ok {
			t.Fatalf("accepted %q", in)
		}
	}
	if maskEmail("someone@example.com") != "s***@example.com" {
		t.Fatal("mask")
	}
}

func TestRollingWindowBoundary(t *testing.T) {
	st, now := newState(), time.Now().UTC()
	b := bucket{"k", 2, time.Minute}
	if ok, _ := st.take(now, b); !ok {
		t.Fatal("first")
	}
	if ok, _ := st.take(now.Add(time.Second), b); !ok {
		t.Fatal("at limit")
	}
	if ok, retry := st.take(now.Add(2*time.Second), b); ok || !retry.Equal(now.Add(time.Minute)) {
		t.Fatal("over limit", retry)
	}
	// (t-window, t]: the first hit leaves exactly at now+window.
	if ok, _ := st.take(now.Add(time.Minute), b); !ok {
		t.Fatal("window edge")
	}
	// take is all-or-nothing; a refused bucket does not spend another.
	other := bucket{"other", 5, time.Minute}
	if ok, _ := st.take(now.Add(time.Minute), b, other); ok || len(st.Rates["other"]) != 0 {
		t.Fatal("partial spend")
	}
	// hit counts refused requests too, keeps limit+1 hits, and skips buckets after the first refusal.
	h, after := bucket{"h", 1, time.Minute}, bucket{"after", 5, time.Minute}
	st.hit(now, h)
	for i := 0; i < 50; i++ {
		if ok, retry := st.hit(now.Add(time.Duration(i)*time.Second), h, after); ok || len(st.Rates["h"]) != 2 || !retry.Equal(st.Rates["h"][1].Add(time.Minute)) {
			t.Fatal("hit", i, len(st.Rates["h"]))
		}
	}
	if len(st.Rates["after"]) != 0 {
		t.Fatal("bucket after refusal spent")
	}
	if ok, _ := st.hit(st.Rates["h"][1].Add(time.Minute), h); !ok {
		t.Fatal("refused at the announced retry time")
	}
	st.sweepIdentity(now.Add(2 * time.Hour))
	if len(st.Rates) != 0 {
		t.Fatal("rates retained past an hour")
	}
}

func TestSendLimits(t *testing.T) {
	st, now := newState(), time.Now().UTC()
	if _, _, _, ok := st.startChallenge("a@b.co", "ip1", now); !ok {
		t.Fatal("first send")
	}
	if _, _, retry, ok := st.startChallenge("a@b.co", "ip1", now.Add(59*time.Second)); ok || !retry.Equal(now.Add(time.Minute)) {
		t.Fatal("60s spacing")
	}
	for i := 1; i < 5; i++ {
		if _, _, _, ok := st.startChallenge("a@b.co", "ip1", now.Add(time.Duration(i)*time.Minute)); !ok {
			t.Fatal("email budget", i)
		}
	}
	if _, _, _, ok := st.startChallenge("a@b.co", "ip1", now.Add(10*time.Minute)); ok {
		t.Fatal("email hourly limit")
	}
	if len(st.Challenges) != 1 {
		t.Fatal("new code must replace older codes", len(st.Challenges))
	}
	// Source IP: 20 per hour across different addresses; a second source stays independent (NAT boundary).
	st = newState()
	for i := 0; i < 20; i++ {
		if _, _, _, ok := st.startChallenge(string(rune('a'+i))+"@b.co", "nat", now); !ok {
			t.Fatal("ip budget", i)
		}
	}
	if _, _, _, ok := st.startChallenge("z@b.co", "nat", now); ok {
		t.Fatal("ip hourly limit")
	}
	if _, _, _, ok := st.startChallenge("z@b.co", "other", now); !ok {
		t.Fatal("other source blocked")
	}
	st = newState()
	for i := 0; i < 100; i++ {
		if _, _, _, ok := st.startChallenge(strings.Repeat("x", i+1)+"@b.co", "ip"+strings.Repeat("y", i/20), now); !ok {
			t.Fatal("global budget", i)
		}
	}
	if _, _, _, ok := st.startChallenge("new@b.co", "fresh", now); ok {
		t.Fatal("global hourly limit")
	}
	if strings.Contains(strings.Join(keys(st.Rates), ","), "@") {
		t.Fatal("raw email kept in rate keys")
	}
}

func keys(m map[string][]time.Time) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestCodeVerificationAndMembers(t *testing.T) {
	st, now := newState(), time.Now().UTC()
	pending, code, _, _ := st.startChallenge("a@b.co", "ip", now)
	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}
	for i := 1; i < codeAttempts; i++ {
		if r := st.verifyChallenge(pending, wrong, "", now); r.Problem != "wrong" || r.Remaining != codeAttempts-i {
			t.Fatal("wrong", r)
		}
	}
	if r := st.verifyChallenge(pending, wrong, "", now); r.Problem != "expired" {
		t.Fatal("fifth wrong must void", r)
	}
	if r := st.verifyChallenge(pending, code, "", now); r.Problem != "expired" || len(st.Owners) != 0 {
		t.Fatal("voided code accepted or owner created before verification")
	}
	pending, code, _, _ = st.startChallenge("a@b.co", "ip", now.Add(time.Minute))
	if r := st.verifyChallenge(pending, code, "", now.Add(time.Minute+codeTTL)); r.Problem != "expired" {
		t.Fatal("expired code accepted")
	}
	pending, code, _, _ = st.startChallenge("a@b.co", "ip", now.Add(2*time.Minute))
	first := st.verifyChallenge(pending, code, "", now.Add(2*time.Minute))
	if first.Session == "" || len(st.Members) != 1 {
		t.Fatal("verified login failed", first)
	}
	if r := st.verifyChallenge(pending, code, "", now.Add(2*time.Minute)); r.Problem != "expired" {
		t.Fatal("code reused")
	}
	// Re-login maps to the same member and owner and replaces the presented session.
	pending, code, _, _ = st.startChallenge("a@b.co", "ip", now.Add(3*time.Minute))
	second := st.verifyChallenge(pending, code, first.Session, now.Add(3*time.Minute))
	if len(st.Members) != 1 || len(st.Owners) != 1 || st.Sessions[hashToken(first.Session)] != nil || st.Sessions[hashToken(second.Session)] == nil {
		t.Fatal("continuity")
	}
	// Same address from another issuer is not merged.
	st.Members["mem_other"] = &Member{"owner_other", "c@d.co", "other-issuer", true, now}
	st.Identities["other-issuer|c@d.co"] = "mem_other"
	st.Owners["owner_other"] = &Owner{"", true}
	pending, code, _, _ = st.startChallenge("c@d.co", "ip", now.Add(4*time.Minute))
	r := st.verifyChallenge(pending, code, "", now.Add(4*time.Minute))
	if st.Sessions[hashToken(r.Session)].Member == "mem_other" || len(st.Members) != 3 {
		t.Fatal("cross-issuer merge")
	}
	// Capacity refuses new members but keeps existing logins.
	for i := len(st.Members); i < activeMembers; i++ {
		st.Members["mem_fill"+string(rune('a'+i%26))+strings.Repeat("z", i)] = &Member{Active: true}
	}
	pending, code, _, _ = st.startChallenge("new@b.co", "ip2", now.Add(5*time.Minute))
	if r := st.verifyChallenge(pending, code, "", now.Add(5*time.Minute)); r.Problem != "capacity" {
		t.Fatal("capacity", r)
	}
	pending, code, _, _ = st.startChallenge("a@b.co", "ip2", now.Add(6*time.Minute))
	if r := st.verifyChallenge(pending, code, "", now.Add(6*time.Minute)); r.Session == "" {
		t.Fatal("existing member blocked by capacity")
	}
	for _, m := range st.Members {
		if m.Email == "" {
			continue
		}
		if st.Owners[m.Owner].Credential != "" {
			t.Fatal("member owner received a bearer credential")
		}
	}
}

func TestSessionLifetime(t *testing.T) {
	st, now := newState(), time.Now().UTC()
	pending, code, _, _ := st.startChallenge("a@b.co", "ip", now)
	token := st.verifyChallenge(pending, code, "", now).Session
	at := now
	for at.Add(50 * time.Minute).Before(now.Add(sessionAbsolute)) {
		at = at.Add(50 * time.Minute)
		if _, _, err := st.session(token, at); err != nil {
			t.Fatal("active session ended early", at.Sub(now))
		}
	}
	if _, _, err := st.session(token, now.Add(sessionAbsolute)); err == nil {
		t.Fatal("absolute lifetime exceeded")
	}
	pending, code, _, _ = st.startChallenge("b@b.co", "ip", now)
	idle := st.verifyChallenge(pending, code, "", now).Session
	if _, _, err := st.session(idle, now.Add(sessionIdle)); err == nil || st.Sessions[hashToken(idle)] != nil {
		t.Fatal("idle lifetime exceeded")
	}
	if _, _, err := st.session("", now); err == nil {
		t.Fatal("empty session")
	}
}

// fakeSMTP accepts one plain SMTP transaction and returns the DATA section.
func fakeSMTP(t *testing.T) (string, chan string) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	got := make(chan string, 1)
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		r := bufio.NewReader(conn)
		reply := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
		reply("220 fake")
		data := ""
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			switch cmd := strings.ToUpper(strings.TrimSpace(line)); {
			case strings.HasPrefix(cmd, "EHLO"):
				reply("250 fake")
			case cmd == "DATA":
				reply("354 go")
				for {
					l, _ := r.ReadString('\n')
					if l == ".\r\n" || l == "" {
						break
					}
					data += l
				}
				got <- data
				reply("250 ok")
			case cmd == "QUIT":
				reply("221 bye")
				return
			default:
				reply("250 ok")
			}
		}
	}()
	return listener.Addr().String(), got
}

func TestSMTPMailer(t *testing.T) {
	for _, bad := range []string{"https://example-secret@host:1", "smtp://host", "smtp://user:example-secret@host:25/path", "::"} {
		if _, err := SMTPMailer(bad, "from@b.co"); err == nil || strings.Contains(err.Error(), "example-secret") {
			t.Fatal("bad url", bad, err)
		}
	}
	if _, err := SMTPMailer("smtp://h:25", "Name <from@b.co>"); err == nil {
		t.Fatal("display-name sender")
	}
	addr, got := fakeSMTP(t)
	send, err := SMTPMailer("smtp://"+addr, "noreply@knowslog.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := send(context.Background(), "a@b.co", "KnowsLink 확인 코드", "코드: 123456\n"); err != nil {
		t.Fatal(err)
	}
	data := <-got
	for _, want := range []string{"From: noreply@knowslog.com", "To: a@b.co", "Subject: =?UTF-8?b?", "charset=UTF-8", "123456"} {
		if !strings.Contains(data, want) {
			t.Fatalf("missing %q in %q", want, data)
		}
	}
	closed, _ := SMTPMailer("smtp://127.0.0.1:1", "noreply@knowslog.com")
	if err := closed(context.Background(), "a@b.co", "s", "b"); err == nil || strings.Contains(err.Error(), "127.0.0.1") {
		t.Fatal("unreachable server", err)
	}
}

// A principal refused by its own bucket must not spend the shared budget; only independent principals together fill it.
func TestRatePrincipalIsolation(t *testing.T) {
	st, now := newState(), time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	allowed := func(b []bucket) bool { ok, _ := st.hit(now, b...); return ok }
	for i := 1; i <= 1000; i++ {
		if allowed(anonymousRate("198.51.100.1")) != (i <= 30) {
			t.Fatal("own IP limit", i)
		}
		if allowed(memberRate("mem_flood")) != (i <= 40) {
			t.Fatal("own member limit", i)
		}
		if allowed(cleanupRate("mem_flood")) != (i <= 20) {
			t.Fatal("own cleanup limit", i)
		}
	}
	if !allowed(anonymousRate("203.0.113.9")) || !allowed(memberRate("mem_other")) || !allowed(cleanupRate("mem_other")) {
		t.Fatal("one refused principal blocked others")
	}
	// Shared new budget: 30+40+1+1 used; 128 more independent hits reach exactly 200, the next one is refused.
	for i := 0; i < 128; i++ {
		if !allowed(anonymousRate(fmt.Sprintf("192.0.2.%d", i))) {
			t.Fatal("below shared limit", i)
		}
	}
	if allowed(anonymousRate("192.0.2.200")) || allowed(memberRate("mem_late")) {
		t.Fatal("shared new limit not enforced")
	}
	// Shared cleanup: 20+1 used; 79 more reach 100.
	for i := 0; i < 79; i++ {
		if !allowed(cleanupRate(fmt.Sprintf("mem_%d", i))) {
			t.Fatal("below cleanup limit", i)
		}
	}
	if allowed(cleanupRate("mem_late")) {
		t.Fatal("shared cleanup limit not enforced")
	}
	// Rolling window: the minute-old hits leave at now+window and everyone recovers.
	now = now.Add(time.Minute)
	if !allowed(anonymousRate("198.51.100.1")) || !allowed(memberRate("mem_late")) || !allowed(cleanupRate("mem_late")) {
		t.Fatal("no recovery after window")
	}
}
