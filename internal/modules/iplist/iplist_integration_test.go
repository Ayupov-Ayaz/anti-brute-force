package iplist

import (
	"context"
	"fmt"
	"testing"
	"time"

	ipService "github.com/ayupov-ayaz/anti-brute-force/internal/modules/ip"

	"github.com/alicebob/miniredis"

	redis "github.com/go-redis/redis/v8"

	redisstorage "github.com/ayupov-ayaz/anti-brute-force/internal/modules/storage"
	"github.com/stretchr/testify/require"
)

func newMiniRedisClient() (*redis.Client, error) {
	mini, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	cli := redis.NewClient(&redis.Options{
		Addr: mini.Addr(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return cli, nil
}

func newList(t *testing.T) (*IPList, *redis.Client) {
	db, err := newMiniRedisClient()
	require.NoError(t, err)

	storage := redisstorage.New(db)
	ip := ipService.New()
	list := New(addr, storage, ip)

	return list, db
}

func TestIPList_WithRedis_Add(t *testing.T) {
	const (
		addr        = "blacklist"
		bannedIPNet = "192.1.1.128/28"
	)

	list, db := newList(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	ctx := context.Background()
	err := list.Add(ctx, bannedIPNet)
	require.NoError(t, err)

	v, err := db.SMembers(ctx, addr).Result()
	require.NoError(t, err)
	require.Equal(t, []string{bannedIPNet}, v)
}

func TestIPList_WithRedis_Remove(t *testing.T) {
	const (
		bannedIPNet  = "192.1.1.128/28"
		anotherIPNet = "192.1.1.64/27"
	)

	list, db := newList(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	ctx := context.Background()

	table := []struct {
		name   string
		ip     string
		after  func(t *testing.T)
		before func(t *testing.T)
	}{
		{
			name:   "remove not banned ip",
			ip:     anotherIPNet,
			after:  func(t *testing.T) {},
			before: func(t *testing.T) {},
		},
		{
			name: "remove banned ip",
			ip:   bannedIPNet,
			before: func(t *testing.T) {
				require.NoError(t, db.SAdd(ctx, addr, bannedIPNet).Err())
			},
			after: func(t *testing.T) {
				res, err := db.SMembers(ctx, addr).Result()
				require.NoError(t, err)
				require.Empty(t, res)
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			err := list.Remove(ctx, tt.ip)
			require.NoError(t, err)
			tt.after(t)
		})
	}
}

func TestIPList_WithRedis_Contains(t *testing.T) {
	const (
		addr        = "blacklist"
		bannedIPNet = "192.1.1.0/26"
		bannedIP    = "192.1.1.63"
		notBannedIP = "192.1.1.143"
	)

	list, db := newList(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	ctx := context.Background()

	tables := []struct {
		name   string
		ip     string
		ok     bool
		before func(t *testing.T)
	}{
		{
			name:   "not found",
			ip:     notBannedIP,
			before: func(t *testing.T) {},
		},
		{
			name: "found ip",
			ip:   bannedIP,
			ok:   true,
			before: func(t *testing.T) {
				require.NoError(t, db.SAdd(ctx, addr, bannedIPNet).Err())
			},
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			ok, err := list.Contains(ctx, tt.ip)
			require.NoError(t, err)
			require.Equal(t, tt.ok, ok)
		})
	}
}
