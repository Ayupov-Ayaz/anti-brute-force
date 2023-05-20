package limiter

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type LeakyBucketLimiter interface {
	Allow(ctx context.Context, client *redis.Client, key string) error
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
	// todo: implement this method
	return nil
}
