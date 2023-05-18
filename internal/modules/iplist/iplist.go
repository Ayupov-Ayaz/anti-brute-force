package iplist

import (
	"context"
	"errors"

	redis "github.com/go-redis/redis/v8"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/storage"
)

const value = "+"

type IPList struct {
	addr    string
	storage storage.Storage
}

func New(addr string, storage storage.Storage) *IPList {
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
