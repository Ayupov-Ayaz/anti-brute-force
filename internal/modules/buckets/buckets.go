package buckets

import "context"

type Buckets struct {
}

func New() *Buckets {
	return &Buckets{}
}

func (b *Buckets) Check(ctx context.Context, ip, login, pass string) error {
	return nil
}

func (b *Buckets) Reset(ctx context.Context, ip, login string) error {
	return nil
}
