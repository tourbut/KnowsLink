//go:build integration

// Integration services share one test Redis: each commit writes a random key name, so parallel test databases never collide.
package relay

import (
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var testRedis = sync.OnceValue(func() *redis.Client {
	options, err := redis.ParseURL(os.Getenv("TEST_REDIS_URL"))
	if err != nil {
		panic("isolated TEST_REDIS_URL required")
	}
	return redis.NewClient(options)
})
