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

type Info struct {
	hits    int64
	limit   int64
	expires time.Time
}

func (i Info) IsExceeded() bool {
	return i.hits > i.limit
}

func (i Info) Remaining() int64 {
	return max(i.limit-i.hits, 0)
}

func (i Info) Resets() time.Duration {
	return i.expires.Sub(time.Now())
}

func (i Info) Limit() int64 {
	return i.limit
}

func New(url string, rate int64, duration time.Duration) (*RateLimiter, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	client := redis.NewClient(opts)

	return &RateLimiter{
		client:   client,
		duration: duration,
		rate:     rate,
	}, nil
}

func (r *RateLimiter) keyFunc(ip net.IP) string {
	return fmt.Sprintf("rate:%s", ip.String())
}

// AddAndCheckIfExceeds is used to determine whether or not the
// rate limit has been exceeded, whilst also adding another hit to it.
func (r *RateLimiter) AddAndCheckIfExceeds(ctx context.Context, ip net.IP) (Info, error) {
	expires := time.Now().Add(r.duration)

	pipeline := r.client.TxPipeline()
	incr := pipeline.Incr(ctx, r.keyFunc(ip))
	if incr.Val() == 1 {
		pipeline.ExpireAt(ctx, r.keyFunc(ip), expires)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return Info{}, fmt.Errorf("failed to exec pipeline: %w", err)
	}

	return Info{
		hits:    incr.Val(),
		limit:   r.rate,
		expires: time.Now().Add(r.duration).Add(time.Second),
	}, nil
}
