// Cleanup removes expired payloads and retained metadata even when no client polls.
package relay

import (
	"context"
	"log"
	"time"
)

func (s *Service) Cleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cleanupCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			_, err := s.transaction(cleanupCtx, func(_ *State, _ time.Time) (any, error) { return nil, nil })
			cancel()
			if err != nil && ctx.Err() == nil {
				log.Print("relay retention cleanup unavailable")
			}
		}
	}
}
