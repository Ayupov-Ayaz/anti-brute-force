package bucket

import (
	"context"
	"testing"
	"time"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/logger"

	"github.com/stretchr/testify/require"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

func TestLeakyBucketLimiter(t *testing.T) {
	ctx := context.Background()
	mini := miniredis.NewMiniRedis()
	require.NoError(t, mini.Start())
	client := redis.NewClient(&redis.Options{Addr: mini.Addr()})
	require.NoError(t, client.Ping(ctx).Err())

	const (
		key            = "test"
		maxRequests    = 10
		refillInterval = 10 * time.Second
	)

	zLogger, err := logger.New("INFO")
	require.NoError(t, err)
	limiter := New(maxRequests, refillInterval, zLogger)

	t.Run("test refill", func(t *testing.T) {
		call := func() error {
			return limiter.Allow(ctx, client, key)
		}
		for i := 0; i < maxRequests; i++ {
			require.NoError(t, call())
		}

		require.ErrorIs(t, call(), apperr.ErrUserIsBlocked)
	})

	t.Run("Reset", func(t *testing.T) {
		require.NoError(t, limiter.Reset(ctx, client, key))
		require.NoError(t, limiter.Allow(ctx, client, key))
	})
}
