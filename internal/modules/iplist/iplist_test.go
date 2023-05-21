package iplist

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	storage "github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist/mock"
	"github.com/golang/mock/gomock"
)

const (
	ip   = "127.0.0.1"
	addr = "blacklist"
)

func TestIPList_Add(t *testing.T) {
	ctx := context.Background()

	save := func(m *storage.MockStorage, err error) *gomock.Call {
		return m.EXPECT().Save(ctx, addr, ip, value).Times(1).Return(err)
	}

	errSave := errors.New("save failed")
	table := []struct {
		name   string
		err    error
		before func(m *storage.MockStorage)
	}{
		{
			name: "fail",
			err:  errSave,
			before: func(m *storage.MockStorage) {
				save(m, errSave)
			},
		},
		{
			name: "success",
			before: func(m *storage.MockStorage) {
				save(m, nil)
			},
		},
	}

	list := New(addr, nil)

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := storage.NewMockStorage(ctrl)
			tt.before(m)
			list.storage = m

			err := list.Add(ctx, ip)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestIPList_Remove(t *testing.T) {
	ctx := context.Background()

	remove := func(m *storage.MockStorage, err error) *gomock.Call {
		return m.EXPECT().Remove(ctx, addr, ip).Times(1).Return(err)
	}

	errRemove := errors.New("remove failed")
	table := []struct {
		name   string
		err    error
		before func(m *storage.MockStorage)
	}{
		{
			name: "fail",
			err:  errRemove,
			before: func(m *storage.MockStorage) {
				remove(m, errRemove)
			},
		},
		{
			name: "success",
			before: func(m *storage.MockStorage) {
				remove(m, nil)
			},
		},
	}

	list := New(addr, nil)
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := storage.NewMockStorage(ctrl)
			tt.before(m)
			list.storage = m

			err := list.Remove(ctx, ip)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestIPList_Contains(t *testing.T) {
	ctx := context.Background()
	load := func(m *storage.MockStorage, val string, err error) *gomock.Call {
		return m.EXPECT().Load(ctx, addr, ip).Times(1).Return(val, err)
	}

	errLoad := errors.New("load failed")
	table := []struct {
		name   string
		err    error
		exp    bool
		before func(m *storage.MockStorage)
	}{
		{
			name: "fail",
			err:  errLoad,
			before: func(m *storage.MockStorage) {
				load(m, "", errLoad)
			},
		},
		{
			name: "another value",
			before: func(m *storage.MockStorage) {
				load(m, "another value", nil)
			},
		},
		{
			name: "success",
			before: func(m *storage.MockStorage) {
				load(m, value, nil)
			},
			exp: true,
		},
	}

	list := New(addr, nil)
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := storage.NewMockStorage(ctrl)
			tt.before(m)
			list.storage = m

			ok, err := list.Contains(ctx, ip)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.exp, ok)
		})
	}
}
