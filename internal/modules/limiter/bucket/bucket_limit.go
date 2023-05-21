package bucket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	redis "github.com/go-redis/redis/v8"
)

type LeakyBucketLimiter struct {
	maxRequests    int64
	refillInterval time.Duration
	logger         zerolog.Logger
}

func New(maxRequests int64, refillInterval time.Duration, logger zerolog.Logger) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		maxRequests:    maxRequests,
		refillInterval: refillInterval,
		logger:         logger,
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
	zCardCmd, ok := cmds[0].(*redis.IntCmd)
	if !ok {
		return fmt.Errorf("failed to cast ZCard command result to IntCmd")
	}

	count := zCardCmd.Val()
	if err != nil {
		return fmt.Errorf("failed to get ZCard result: %v", err)
	}

	l.logger.Debug().Int64("count", count).Int64("max", l.maxRequests).Msg("")

	if count >= l.maxRequests {
		return fmt.Errorf("%w (bucket: %d/%d)",
			apperr.ErrUserIsBlocked, count, l.maxRequests)
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

func (l *LeakyBucketLimiter) Reset(ctx context.Context, client *redis.Client, key string) error {
	max := strconv.FormatInt(time.Now().Unix(), 10)
	count, err := client.ZRemRangeByScore(ctx, key, "-inf", max).Result()
	l.logger.Debug().Int64("count", count).Str("key", key).Msg("reset")
	return err
}
