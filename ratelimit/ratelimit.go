package ratelimit

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client   *redis.Client
	duration time.Duration
	rate     int64
}

func New(client *redis.Client, rate int64, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:   client,
		duration: duration,
		rate:     rate,
	}
}

// AddAndCheckIfExceeds is used to determine whether or not the
// rate limit has been exceeded, whilst also adding another hit to it.
func (r *RateLimiter) AddAndCheckIfExceeds(ctx context.Context, ip net.IP) (bool, error) {
	// Start actions in a multi / pipeline tx
	p := r.client.Pipeline()

	// Incr the ip and capture the value
	res := p.Incr(ctx, ip.String())

	// Set an expiration on the rate only if no expiration already exists
	p.ExpireNX(ctx, ip.String(), r.duration)

	// Run the pipline
	if _, err := p.Exec(ctx); err != nil {
		return false, fmt.Errorf("failed to exec pipeline: %w", err)
	}

	// Check if rate has been exceeded
	hasExceeded := res.Val() > r.rate

	// Return the result
	return hasExceeded, nil
}
