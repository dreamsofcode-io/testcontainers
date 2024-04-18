package ratelimit_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/dreamsofcode-io/testcontainers/ratelimit"
)

func TestRateLimiter(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7.2",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	endpoint, err := container.Endpoint(ctx, "")
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	limiter := ratelimit.New(client, 3, time.Minute)

	ip := "192.168.1.54"

	t.Run("happy path flow", func(t *testing.T) {
		res, err := limiter.AddAndCheckIfExceeds(ctx, net.ParseIP(ip))
		assert.NoError(t, err)

		// Rate should not be exceeded
		assert.False(t, res)

		// Check key exists
		assert.Equal(t, client.Get(ctx, ip).Val(), "1")

		client.FlushAll(ctx)
	})

	t.Run("should expire after three times", func(t *testing.T) {
		client.Set(ctx, ip, "3", 0)

		res, err := limiter.AddAndCheckIfExceeds(ctx, net.ParseIP(ip))
		assert.NoError(t, err)

		// Rate should be exceeded
		assert.True(t, res)

		// Check expire time is set
		assert.Greater(t, client.ExpireTime(ctx, ip).Val(), time.Duration(0))
	})
}
