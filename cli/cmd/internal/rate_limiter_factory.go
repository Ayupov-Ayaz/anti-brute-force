package internal

import (
	limitercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/limiter"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/limiter"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/limiter/bucket"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

func NewAuthRateLimiter(config limitercfg.Limiter, redis *redis.Client, logger zerolog.Logger) *limiter.AuthRateLimiter {
	loginBucket := bucket.New(config.Login.Count, config.Login.Interval, logger)
	passwordBucket := bucket.New(config.Password.Count, config.Password.Interval, logger)
	ipBucket := bucket.New(config.IP.Count, config.IP.Interval, logger)

	return limiter.New(redis, loginBucket, passwordBucket, ipBucket)
}
