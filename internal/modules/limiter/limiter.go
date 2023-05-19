package limiter

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	redis "github.com/go-redis/redis/v8"
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

func subLoggerByLimiter(logger zerolog.Logger, limiter string) zerolog.Logger {
	return logger.With().Str("limiter", limiter).Logger()
}

func WithLoginLimiter(maxRequests int64, refillInterval time.Duration, logger zerolog.Logger) Config {
	return func(a *AuthRateLimiter) {
		a.loginRateLimiter = NewLeakyBucketLimiter(maxRequests, refillInterval,
			subLoggerByLimiter(logger, "login"))
	}
}

func WithPasswordLimiter(maxRequests int64, refillInterval time.Duration, logger zerolog.Logger) Config {
	return func(a *AuthRateLimiter) {
		a.passwordRateLimiter = NewLeakyBucketLimiter(maxRequests, refillInterval,
			subLoggerByLimiter(logger, "password"))
	}
}

func WithIPLimiter(maxRequests int64, refillInterval time.Duration, logger zerolog.Logger) Config {
	return func(a *AuthRateLimiter) {
		a.ipRateLimiter = NewLeakyBucketLimiter(maxRequests, refillInterval,
			subLoggerByLimiter(logger, "ip"))
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
