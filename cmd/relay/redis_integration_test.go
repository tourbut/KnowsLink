//go:build integration

// The health test service uses the shared test Redis, like the relay package integration tests.
package main

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
