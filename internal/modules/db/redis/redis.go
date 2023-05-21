package redisdb

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func NewRedisClient(addr, user, pass string) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: pass,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return cli, nil
}
