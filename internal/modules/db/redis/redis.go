package redisdb

import (
	"context"
	"time"

	"github.com/alicebob/miniredis"

	redis "github.com/go-redis/redis/v8"
)

func ping(cli *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return cli.Ping(ctx).Err()
}

func NewRedisClient(addr, user, pass string) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: pass,
	})

	if err := ping(cli); err != nil {
		return nil, err
	}

	return cli, nil
}

func NewMiniRedisClient() (*redis.Client, error) {
	mini, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	cli := redis.NewClient(&redis.Options{
		Addr: mini.Addr(),
	})

	if err := ping(cli); err != nil {
		return nil, err
	}

	return cli, nil
}
