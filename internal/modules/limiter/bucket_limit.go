package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	"github.com/go-redis/redis/v8"
)

type LeakyBucketLimiter struct {
	maxRequests    int64
	refillInterval time.Duration
	field          string
}

func NewLeakyBucketLimiter(maxRequests int64, refillInterval time.Duration, field string) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		maxRequests:    maxRequests,
		refillInterval: refillInterval,
		field:          field,
	}
}

func (l *LeakyBucketLimiter) Allow(ctx context.Context, client *redis.Client, key string) error {
	pipe := client.Pipeline()
	//получаем количество запросов в бакете

	pipe.ZCard(ctx, key)
	//defer pipe.Exec(ctx)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute Redis pipeline: %v", err)
	}

	// Get the result of ZCard command
	zCardCmd := cmds[0].(*redis.IntCmd)
	count := zCardCmd.Val()
	if err != nil {
		return fmt.Errorf("failed to get ZCard result: %v", err)
	}

	log.Debug().Int64("count", count).Int64("max", l.maxRequests).Msg("")
	if count >= l.maxRequests {
		return fmt.Errorf("%w blocked by field %s (%d/%d)",
			apperr.ErrUserIsBlocked, l.field, count, l.maxRequests)
	}

	// Add current request to the bucket
	score := float64(time.Now().Unix())
	// Add current request to the bucket
	pipe.ZAdd(ctx, key, &redis.Z{Score: score, Member: score})
	pipe.Expire(ctx, key, l.refillInterval)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute Redis pipeline: %v", err)
	}

	return nil
}
