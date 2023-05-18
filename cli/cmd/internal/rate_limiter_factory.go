package internal

import (
	limitercfg "github.com/ayupov-ayaz/anti-brute-force/config/limiter"
	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/limiter"
	redis "github.com/go-redis/redis/v8"
)

func NewAuthRateLimiter(config limitercfg.Limiter, redis *redis.Client) *limiter.AuthRateLimiter {
	ip := config.IP
	login := config.Login
	password := config.Password

	return limiter.NewAuthRateLimiter(
		limiter.WithIPLimiter(ip.Count, ip.Interval),
		limiter.WithLoginLimiter(login.Count, login.Interval),
		limiter.WithRedisClient(redis),
		limiter.WithPasswordLimiter(password.Count, password.Interval),
	)
}
