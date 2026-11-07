// DB limits bound how long a stalled session or a queue of writers can hold relay requests.
package relay

import (
	"bytes"
	"context"
	"net"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Every write takes the same row lock, so extra connections only queue inside Postgres.
	poolConnections = 4
	// The connection left over keeps /healthz and the lock-free snapshot read answering while writers wait.
	maxTransactions = poolConnections - 1
	// A request waits for at most two transactions (admission, handler) and then spends finishWait releasing its
	// admission. The sum stays under the 10s HTTP deadline, so the failure reaches the client (TestDBLimitsFitHTTPDeadline).
	transactionSlotWait = 2 * time.Second
	lockTimeout         = 1500 * time.Millisecond
	finishWait          = 2 * time.Second
	// Between the row read and the save the session is idle in transaction while the relay computes. This exceeds that
	// computation at the 20,000 message cap under the production CPU limit.
	idleTransactionTimeout = 5 * time.Second
	// A dead peer is noticed after keepAliveIdle + keepAliveInterval*keepAliveCount.
	keepAliveIdle     = 5 * time.Second
	keepAliveInterval = 2 * time.Second
	keepAliveCount    = 3
)

// PoolConfig applies the relay's connection limits to every connection, including the LISTEN connections copied from it.
func PoolConfig(databaseURL string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = poolConnections
	ms := func(d time.Duration) string { return strconv.FormatInt(d.Milliseconds(), 10) }
	seconds := func(d time.Duration) string { return strconv.Itoa(int(d.Seconds())) }
	for name, value := range map[string]string{
		"lock_timeout":                        ms(lockTimeout),
		"idle_in_transaction_session_timeout": ms(idleTransactionTimeout),
		"tcp_keepalives_idle":                 seconds(keepAliveIdle),
		"tcp_keepalives_interval":             seconds(keepAliveInterval),
		"tcp_keepalives_count":                strconv.Itoa(keepAliveCount),
		// Keepalive is silent while sent data is unacknowledged; a backend stuck writing to a dead relay would keep the row lock.
		"tcp_user_timeout": ms(keepAliveIdle + keepAliveCount*keepAliveInterval),
	} {
		cfg.ConnConfig.RuntimeParams[name] = value
	}
	dialer := &net.Dialer{KeepAliveConfig: net.KeepAliveConfig{Enable: true, Idle: keepAliveIdle, Interval: keepAliveInterval, Count: keepAliveCount}}
	cfg.ConnConfig.DialFunc = dialer.DialContext
	return cfg, nil
}

// encodedNUL reports an escaped NUL in encoded JSON. The json column would store it, but the state has never held one:
// jsonb refused it at save, and readers and operator jsonb casts rely on that. A preceding escaped backslash makes the
// text a literal "\u0000", which is allowed.
func encodedNUL(raw []byte) bool {
	escape := []byte(`\u0000`)
	for {
		i := bytes.Index(raw, escape)
		if i < 0 {
			return false
		}
		backslashes := 0
		for i-1-backslashes >= 0 && raw[i-1-backslashes] == '\\' {
			backslashes++
		}
		if backslashes%2 == 0 {
			return true
		}
		raw = raw[i+len(escape):]
	}
}

// enterTransaction takes a transaction slot. Waiters are served in arrival order, so the wait stays bounded by the queue
// ahead instead of depending on which backend Postgres wakes next.
func (s *Service) enterTransaction(ctx context.Context) error {
	s.transactionOnce.Do(func() { s.transactionSlots = make(chan struct{}, maxTransactions) })
	timer := time.NewTimer(transactionSlotWait)
	defer timer.Stop()
	select {
	case s.transactionSlots <- struct{}{}:
		return nil
	case <-ctx.Done():
	case <-timer.C:
	}
	return fault("unavailable")
}
