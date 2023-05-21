package limiter

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

type LeakyBucketLimiter interface {
	Allow(ctx context.Context, client *redis.Client, key string) error
	Reset(ctx context.Context, client *redis.Client, key string) error
}

type AuthRateLimiter struct {
	client              *redis.Client
	loginRateLimiter    LeakyBucketLimiter
	passwordRateLimiter LeakyBucketLimiter
	ipRateLimiter       LeakyBucketLimiter
}

func New(redis *redis.Client, login, pass, ip LeakyBucketLimiter) *AuthRateLimiter {
	return &AuthRateLimiter{
		client:              redis,
		loginRateLimiter:    login,
		passwordRateLimiter: pass,
		ipRateLimiter:       ip,
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
	if err := a.loginRateLimiter.Reset(ctx, a.client, login); err != nil {
		return fmt.Errorf("reset login rate limiter: %w", err)
	}

	if err := a.ipRateLimiter.Reset(ctx, a.client, ip); err != nil {
		return fmt.Errorf("reset ip rate limiter: %w", err)
	}

	return nil
}
