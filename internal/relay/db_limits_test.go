// DB limit values keep failures inside the HTTP deadline and transaction waits bounded without a database.
package relay

import (
	"context"
	"testing"
	"time"
)

func TestPoolConfigAppliesLimits(t *testing.T) {
	cfg, err := PoolConfig("postgres://user:pass@db.example:5432/app?sslmode=disable&lock_timeout=0")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MaxConns != poolConnections || maxTransactions >= poolConnections {
		t.Fatalf("max conns %d, transactions %d", cfg.MaxConns, maxTransactions)
	}
	for name, want := range map[string]string{
		"lock_timeout":                        "1500",
		"idle_in_transaction_session_timeout": "5000",
		"tcp_keepalives_idle":                 "5",
		"tcp_keepalives_interval":             "2",
		"tcp_keepalives_count":                "3",
		"tcp_user_timeout":                    "11000",
	} {
		if got := cfg.ConnConfig.RuntimeParams[name]; got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
	// The LISTEN connection copies this config, so it inherits every limit.
	if copied := cfg.ConnConfig.Copy(); copied.RuntimeParams["lock_timeout"] != "1500" || copied.DialFunc == nil {
		t.Fatal("LISTEN connection config lost the limits")
	}
}

func TestDBLimitsFitHTTPDeadline(t *testing.T) {
	const httpDeadline = 10 * time.Second
	if worst := 2*(transactionSlotWait+lockTimeout) + finishWait; worst >= httpDeadline {
		t.Fatalf("worst failure path %v does not fit %v", worst, httpDeadline)
	}
}

func TestTransactionSlotWaitIsBounded(t *testing.T) {
	s := &Service{}
	for range maxTransactions {
		if err := s.enterTransaction(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	if err := s.enterTransaction(cancelled); err == nil || err.Error() != "unavailable" || time.Since(start) > time.Second {
		t.Fatalf("cancelled wait: %v after %v", err, time.Since(start))
	}
	start = time.Now()
	if err := s.enterTransaction(context.Background()); err == nil || err.Error() != "unavailable" || time.Since(start) < transactionSlotWait {
		t.Fatalf("bounded wait: %v after %v", err, time.Since(start))
	}
	<-s.transactionSlots
	if err := s.enterTransaction(context.Background()); err != nil {
		t.Fatalf("freed slot refused: %v", err)
	}
}

func TestEncodedNUL(t *testing.T) {
	back := string(rune(92))
	for _, test := range []struct {
		raw  string
		want bool
	}{
		{`{"a":"plain"}`, false},
		{`{"a":"` + back + `u0000"}`, true},
		{`{"a":"` + back + back + `u0000"}`, false},
		{`{"a":"` + back + back + back + `u0000"}`, true},
		{`{"a":"` + back + back + `u0000` + back + `u0000"}`, true},
		{`{"a":"` + back + `u00001"}`, true},
	} {
		if got := encodedNUL([]byte(test.raw)); got != test.want {
			t.Errorf("%s: got %t", test.raw, got)
		}
	}
}
