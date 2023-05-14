package iplist

import "context"

type IPList struct {
	addr string
}

func New(addr string) *IPList {
	return &IPList{
		addr: addr,
	}
}

func (b *IPList) Add(ctx context.Context, key string) error {
	return nil
}

func (b *IPList) Remove(ctx context.Context, key string) error {
	return nil
}

func (b *IPList) Contains(ctx context.Context, key string) (bool, error) {
	return false, nil
}
