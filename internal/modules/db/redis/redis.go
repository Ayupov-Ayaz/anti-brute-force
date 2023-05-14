package redisstore

import (
	"context"
	"time"

	"github.com/alicebob/miniredis"

	"github.com/redis/go-redis/v9"
)

type Config interface {
	Address() string
	Username() string
	Password() string
}

func ping(cli *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return cli.Ping(ctx).Err()
}

func NewRedisClient(config Config) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     config.Address(),
		Username: config.Username(),
		Password: config.Password(),
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
