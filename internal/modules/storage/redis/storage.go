package redisstorage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	HSet(ctx context.Context, key string, value ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
}

type Storage struct {
	redis Redis
}

func New(redis Redis) *Storage {
	return &Storage{
		redis: redis,
	}
}

func (s *Storage) Save(ctx context.Context, key, field string, val interface{}) error {
	return s.redis.HSet(ctx, key, field, val).Err()
}

func (s *Storage) Load(ctx context.Context, key, field string) (string, error) {
	return s.redis.HGet(ctx, key, field).Result()
}

func (s *Storage) Remove(ctx context.Context, key, field string) error {
	return s.redis.HDel(ctx, key, field).Err()
}
