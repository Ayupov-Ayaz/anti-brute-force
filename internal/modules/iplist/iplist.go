package iplist

import (
	"context"
	"errors"

	redis "github.com/go-redis/redis/v8"
)

const value = "+"

//go:generate mockgen -destination=./mock/storage.go -package=storage github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist Storage
type Storage interface {
	Save(ctx context.Context, key, field string, val interface{}) error
	Load(ctx context.Context, key, field string) (string, error)
	Remove(ctx context.Context, key, field string) error
}

type IPList struct {
	addr    string
	storage Storage
}

func New(addr string, storage Storage) *IPList {
	return &IPList{
		addr:    addr,
		storage: storage,
	}
}

func (b *IPList) Add(ctx context.Context, ip string) error {
	return b.storage.Save(ctx, b.addr, ip, value)
}

func (b *IPList) Remove(ctx context.Context, ip string) error {
	return b.storage.Remove(ctx, b.addr, ip)
}

func (b *IPList) Contains(ctx context.Context, ip string) (bool, error) {
	v, err := b.storage.Load(ctx, b.addr, ip)
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	return v == value, nil
}
