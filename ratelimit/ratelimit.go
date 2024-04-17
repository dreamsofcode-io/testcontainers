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

func New(client *redis.Client, rate int64, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:   client,
		duration: duration,
		rate:     rate,
	}
}

func (r *RateLimiter) keyFunc(ip net.IP) string {
	return fmt.Sprintf("%s", ip.String())
}

// AddAndCheckIfExceeds is used to determine whether or not the
// rate limit has been exceeded, whilst also adding another hit to it.
func (r *RateLimiter) AddAndCheckIfExceeds(ctx context.Context, ip net.IP) (Info, error) {
	p := r.client.Pipeline()

	incr := p.Incr(ctx, r.keyFunc(ip))
	p.ExpireNX(ctx, r.keyFunc(ip), r.duration)
	expires := p.ExpireTime(ctx, r.keyFunc(ip)).Val()

	if _, err := p.Exec(ctx); err != nil {
		return Info{}, err
	}

	return Info{
		hits:    incr.Val(),
		limit:   r.rate,
		expires: time.Unix(0, 0).Add(expires),
	}, nil
}
