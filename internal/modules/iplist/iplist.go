package iplist

import (
	"context"
	"errors"
	"net"

	redis "github.com/go-redis/redis/v8"
)

var (
	ErrIPInvalid   = errors.New("invalid ip")
	ErrInvalidCIDR = errors.New("invalid cidr")
)

//go:generate mockgen -destination=./mock/storage.go -package=mocks github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist Storage
type Storage interface {
	Save(ctx context.Context, key string, val interface{}) error
	Load(ctx context.Context, key string) ([]string, error)
	Remove(ctx context.Context, key, field string) error
}

//go:generate mockgen -destination=./mock/ip_service.go -package=mocks github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist IPService
type IPService interface {
	ParseCIDR(ip string) (*net.IPNet, error)
}

type IPList struct {
	addr    string
	storage Storage
	ip      IPService
}

func New(addr string, storage Storage, ip IPService) *IPList {
	return &IPList{
		addr:    addr,
		storage: storage,
		ip:      ip,
	}
}

func (b *IPList) Add(ctx context.Context, ipNet string) error {
	if _, err := b.ip.ParseCIDR(ipNet); err != nil {
		return err
	}

	return b.storage.Save(ctx, b.addr, ipNet)
}

func (b *IPList) Remove(ctx context.Context, ipNet string) error {
	if _, err := b.ip.ParseCIDR(ipNet); err != nil {
		return err
	}

	return b.storage.Remove(ctx, b.addr, ipNet)
}

func (b *IPList) Contains(ctx context.Context, ip string) (bool, error) {
	ipV4 := net.ParseIP(ip)
	if ipV4.String() == "" {
		return false, ErrIPInvalid
	}

	subnets, err := b.storage.Load(ctx, b.addr)
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	var ipNet *net.IPNet
	for _, subnet := range subnets {
		ipNet, err = b.ip.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}

		if ipNet.Contains(ipV4) {
			return true, nil
		}
	}

	return false, nil
}
