package ratelimit_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/dreamsofcode-io/testcontainers/ratelimit"
)

func loadClient() (*redis.Client, error) {
	opts, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}

func TestRateLimiter(t *testing.T) {
	ctx := context.Background()

	client, err := loadClient()
	assert.NoError(t, err)

	limiter := ratelimit.New(client, 3, time.Minute)

	ip := "192.168.1.54"

	t.Run("happy path flow", func(t *testing.T) {
		res, err := limiter.AddAndCheckIfExceeds(ctx, net.ParseIP(ip))
		assert.NoError(t, err)

		// Rate should not be exceeded
		assert.False(t, res.IsExceeded())

		// Check key exists
		assert.Equal(t, client.Get(ctx, ip).Val(), "1")

		client.FlushAll(ctx)
	})

	t.Run("should expire after three times", func(t *testing.T) {
		client.Set(ctx, ip, "3", 0)

		res, err := limiter.AddAndCheckIfExceeds(ctx, net.ParseIP(ip))
		assert.NoError(t, err)

		// Rate should be exceeded
		assert.True(t, res.IsExceeded())

		// Check expire time is set
		assert.Greater(t, client.ExpireTime(ctx, ip).Val(), time.Duration(0))
	})
}
