package storage

import "context"

//go:generate mockgen -destination=./mock/storage.go -package=storage github.com/ayupov-ayaz/anti-brute-force/internal/storage Storage
type Storage interface {
	Save(ctx context.Context, list, ip string, value interface{}) error
	Load(ctx context.Context, list, ip string) (string, error)
	Remove(ctx context.Context, list, ip string) error
}
