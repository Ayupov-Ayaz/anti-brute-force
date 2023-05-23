package storage

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type Redis interface {
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
}

type Storage struct {
	redis Redis
}

func New(redis Redis) *Storage {
	return &Storage{
		redis: redis,
	}
}

func (s *Storage) Save(ctx context.Context, key string, val interface{}) error {
	return s.redis.SAdd(ctx, key, val).Err()
}

func (s *Storage) Load(ctx context.Context, key string) ([]string, error) {
	return s.redis.SMembers(ctx, key).Result()
}

func (s *Storage) Remove(ctx context.Context, key, member string) error {
	return s.redis.SRem(ctx, key, member).Err()
}
