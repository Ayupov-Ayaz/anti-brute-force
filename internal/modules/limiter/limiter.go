package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthRateLimiter struct {
	client              *redis.Client
	loginRateLimiter    *LeakyBucketLimiter
	passwordRateLimiter *LeakyBucketLimiter
	ipRateLimiter       *LeakyBucketLimiter
}

type Config func(*AuthRateLimiter)

func NewAuthRateLimiter(configs ...Config) *AuthRateLimiter {
	rt := &AuthRateLimiter{}
	for _, config := range configs {
		config(rt)
	}

	return rt
}

func WithRedisClient(client *redis.Client) Config {
	return func(a *AuthRateLimiter) {
		a.client = client
	}
}

func WithLoginLimiter(maxRequests int64, refillInterval time.Duration) Config {
	return func(a *AuthRateLimiter) {
		a.loginRateLimiter = NewLeakyBucketLimiter(maxRequests, refillInterval, "login_limiter")
	}
}

func WithPasswordLimiter(maxRequests int64, refillInterval time.Duration) Config {
	return func(a *AuthRateLimiter) {
		a.passwordRateLimiter = NewLeakyBucketLimiter(maxRequests, refillInterval, "password_limiter")
	}
}

func WithIPLimiter(maxRequests int64, refillInterval time.Duration) Config {
	return func(a *AuthRateLimiter) {
		a.ipRateLimiter = NewLeakyBucketLimiter(maxRequests, refillInterval, "ip_limiter")
	}
}

func (a *AuthRateLimiter) AllowByLogin(ctx context.Context, login string) error {
	return a.loginRateLimiter.Allow(ctx, a.client, login)
}

func (a *AuthRateLimiter) AllowByPassword(ctx context.Context, password string) error {
	return a.passwordRateLimiter.Allow(ctx, a.client, password)
}

func (a *AuthRateLimiter) AllowByIP(ctx context.Context, ip string) error {
	return a.ipRateLimiter.Allow(ctx, a.client, ip)
}

func (a *AuthRateLimiter) Reset(ctx context.Context, ip, login string) error {
	// todo: implement this method
	return nil
}
