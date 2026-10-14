//go:build integration

// Benchmark the actual singleton transaction only against an explicitly marked synthetic Postgres database.
package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BenchmarkRetainedPostgres(b *testing.B) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		b.Fatal("isolated TEST_DATABASE_URL and TEST_SYNTHETIC_DATABASE=1 required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		b.Fatal(err)
	}
	defer pool.Close()
	for _, count := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			st, _ := retentionFixture(count)
			for _, m := range st.Messages {
				m.Receipt.Accepted, m.Receipt.Exp = time.Now().Add(-time.Hour), time.Now().Add(-time.Minute)
			}
			raw, err := json.Marshal(st)
			if err != nil {
				b.Fatal(err)
			}
			if _, err := pool.Exec(ctx, "UPDATE relay_state SET data=$1, epoch=epoch+1, clock=clock_timestamp()", raw); err != nil {
				b.Fatal(err)
			}
			s := &Service{Pool: pool}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := s.transaction(ctx, func(current *State, _ time.Time) (any, error) {
					if len(current.Messages) != count || len(current.Idempotency) != count {
						return nil, fmt.Errorf("retained history changed")
					}
					return nil, nil
				}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
