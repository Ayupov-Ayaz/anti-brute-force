package iplist

import (
	"context"
	"testing"

	redis "github.com/go-redis/redis/v8"

	redisstore "github.com/ayupov-ayaz/anti-brute-force/internal/modules/db/redis"
	redisstorage "github.com/ayupov-ayaz/anti-brute-force/internal/modules/storage"
	"github.com/stretchr/testify/require"
)

func newList(t *testing.T) (*IPList, *redis.Client) {
	db, err := redisstore.NewMiniRedisClient()
	require.NoError(t, err)

	storage := redisstorage.New(db)
	list := New(addr, storage)

	return list, db
}

func TestIPList_WithRedis_Add(t *testing.T) {
	const (
		addr     = "blacklist"
		bannedIP = "254.13.2.11"
	)

	list, db := newList(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	ctx := context.Background()
	err := list.Add(ctx, bannedIP)
	require.NoError(t, err)

	v, err := db.HGet(ctx, addr, bannedIP).Result()
	require.NoError(t, err)
	require.Equal(t, value, v)
}

func TestIPList_WithRedis_Remove(t *testing.T) {
	const (
		bannedIP  = "123.113.23.1"
		anotherIP = "213.41.13.5"
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
			ip:     anotherIP,
			after:  func(t *testing.T) {},
			before: func(t *testing.T) {},
		},
		{
			name: "remove banned ip",
			ip:   bannedIP,
			before: func(t *testing.T) {
				require.NoError(t, db.HSet(ctx, addr, bannedIP, value).Err())
			},
			after: func(t *testing.T) {
				_, err := db.HGet(ctx, addr, bannedIP).Result()
				require.ErrorIs(t, err, redis.Nil)
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
		addr           = "blacklist"
		bannedIP       = "254.13.2.10"
		anotherValueIP = "235.0.40.2"
		notFoundIP     = "225.10.243.3"
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
			ip:     notFoundIP,
			before: func(t *testing.T) {},
		},
		{
			name: "another value",
			ip:   anotherValueIP,
			before: func(t *testing.T) {
				require.NoError(t, db.HSet(ctx, addr, anotherValueIP, "another value").Err())
			},
		},
		{
			name: "same value",
			ip:   bannedIP,
			ok:   true,
			before: func(t *testing.T) {
				require.NoError(t, db.HSet(ctx, addr, bannedIP, value).Err())
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
