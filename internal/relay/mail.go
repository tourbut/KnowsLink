// SMTP submission sends verification codes through any provider; credentials stay in the URL and are never logged.
package relay

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"
	"time"
)

// SMTPMailer accepts smtps://user:pass@host:465 (implicit TLS) or smtp://[user:pass@]host:port (STARTTLS when offered).
// net/smtp refuses PLAIN auth without TLS except on localhost.
func SMTPMailer(rawURL, from string) (Mailer, error) {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "smtp" && u.Scheme != "smtps") || u.Hostname() == "" || u.Port() == "" || u.Path != "" {
		return nil, errors.New("KNOWSLINK_SMTP_URL must be smtp:// or smtps:// with host and port")
	}
	sender, err := mail.ParseAddress(from)
	if err != nil || sender.Name != "" {
		return nil, errors.New("KNOWSLINK_MAIL_FROM must be one bare address")
	}
	host := u.Hostname()
	return func(ctx context.Context, to, subject, body string) error {
		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		dialer := &net.Dialer{}
		var conn net.Conn
		var err error
		if u.Scheme == "smtps" {
			conn, err = (&tls.Dialer{NetDialer: dialer, Config: &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}}).DialContext(ctx, "tcp", u.Host)
		} else {
			conn, err = dialer.DialContext(ctx, "tcp", u.Host)
		}
		if err != nil {
			return errors.New("smtp connect failed")
		}
		deadline, _ := ctx.Deadline()
		_ = conn.SetDeadline(deadline)
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			conn.Close()
			return errors.New("smtp greeting failed")
		}
		defer c.Close()
		if u.Scheme == "smtp" {
			if ok, _ := c.Extension("STARTTLS"); ok {
				if err := c.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
					return errors.New("smtp starttls failed")
				}
			} else if ip := net.ParseIP(host); host != "localhost" && (ip == nil || !ip.IsLoopback()) {
				// Codes never travel in cleartext beyond this host.
				return errors.New("smtp server without TLS refused")
			}
		}
		if u.User != nil {
			password, _ := u.User.Password()
			if err := c.Auth(smtp.PlainAuth("", u.User.Username(), password, host)); err != nil {
				return errors.New("smtp auth failed")
			}
		}
		if err := c.Mail(sender.Address); err != nil {
			return errors.New("smtp sender rejected")
		}
		if err := c.Rcpt(to); err != nil {
			return errors.New("smtp recipient rejected")
		}
		w, err := c.Data()
		if err != nil {
			return errors.New("smtp data failed")
		}
		domain := sender.Address[strings.LastIndex(sender.Address, "@")+1:]
		fmt.Fprintf(w, "From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMessage-ID: <%s@%s>\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n",
			sender.Address, to, mime.BEncoding.Encode("UTF-8", subject), time.Now().UTC().Format(time.RFC1123Z), randomToken()[:24], domain)
		qp := quotedprintable.NewWriter(w)
		_, _ = qp.Write([]byte(strings.ReplaceAll(body, "\n", "\r\n")))
		if qp.Close() != nil || w.Close() != nil {
			return errors.New("smtp message rejected")
		}
		_ = c.Quit()
		return nil
	}, nil
}
