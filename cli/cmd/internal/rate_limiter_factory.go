package internal

import (
	limitercfg "github.com/ayupov-ayaz/anti-brute-force/config/limiter"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/limiter"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

func NewAuthRateLimiter(config limitercfg.Limiter, redis *redis.Client, logger zerolog.Logger) *limiter.AuthRateLimiter {
	ip := config.IP
	login := config.Login
	password := config.Password

	return limiter.NewAuthRateLimiter(
		limiter.WithIPLimiter(ip.Count, ip.Interval, logger),
		limiter.WithLoginLimiter(login.Count, login.Interval, logger),
		limiter.WithRedisClient(redis),
		limiter.WithPasswordLimiter(password.Count, password.Interval, logger),
	)
}
