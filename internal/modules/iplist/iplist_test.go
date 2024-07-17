package iplist

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/ayupov-ayaz/anti-brute-force/internal/modules/iplist/mock"
	"github.com/golang/mock/gomock"
)

const (
	subnet  = "192.1.1.0/25"
	subnet2 = "192.1.1.0/26"
	subnet3 = "192.1.1.64/27"
	subnet4 = "192.1.1.128/28"
	ip      = "192.1.1.100"
	addr    = "blacklist"
)

var (
	_, subNet, _  = net.ParseCIDR(subnet)
	_, subNet2, _ = net.ParseCIDR(subnet2)
	_, subNet3, _ = net.ParseCIDR(subnet3)
	_, subNet4, _ = net.ParseCIDR(subnet4)
)

type mock struct {
	ctx       context.Context
	storage   *mocks.MockStorage
	ipService *mocks.MockIPService
}

func newMock(ctx context.Context, ctrl *gomock.Controller) *mock {
	return &mock{
		ctx:       ctx,
		storage:   mocks.NewMockStorage(ctrl),
		ipService: mocks.NewMockIPService(ctrl),
	}
}

func (m *mock) save(val string) *gomock.Call {
	return m.storage.EXPECT().Save(m.ctx, addr, val).Times(1)
}

func (m *mock) load() *gomock.Call {
	return m.storage.EXPECT().Load(m.ctx, addr).Times(1)
}

func (m *mock) remove() *gomock.Call {
	return m.storage.EXPECT().Remove(m.ctx, addr, subnet).Times(1)
}

func (m *mock) parseCIDR(net string) *gomock.Call {
	return m.ipService.EXPECT().ParseCIDR(net).Times(1)
}

func TestIPList_Add(t *testing.T) {
	ctx := context.Background()

	errSave := errors.New("save failed")
	table := []struct {
		name   string
		ip     string
		err    error
		before func(m *mock)
	}{
		{
			name: "invalid cidr",
			ip:   ip,
			err:  ErrInvalidCIDR,
			before: func(m *mock) {
				m.parseCIDR(ip).Return(nil, ErrInvalidCIDR)
			},
		},
		{
			name: "fail",
			ip:   subnet,
			err:  errSave,
			before: func(m *mock) {
				gomock.InOrder(
					m.parseCIDR(subnet).Return(subNet, nil),
					m.save(subnet).Return(errSave),
				)
			},
		},
		{
			name: "success",
			ip:   subnet,
			before: func(m *mock) {
				gomock.InOrder(
					m.parseCIDR(subnet).Return(subNet, nil),
					m.save(subnet).Return(nil),
				)
			},
		},
	}

	list := New(addr, nil, nil)

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctx, ctrl)
			tt.before(m)
			list.storage = m.storage
			list.ip = m.ipService

			err := list.Add(ctx, tt.ip)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestIPList_Remove(t *testing.T) {
	ctx := context.Background()

	errRemove := errors.New("remove failed")
	table := []struct {
		name   string
		err    error
		ip     string
		before func(m *mock)
	}{
		{
			name: "invalid cidr",
			ip:   ip,
			err:  ErrInvalidCIDR,
			before: func(m *mock) {
				m.parseCIDR(ip).Return(nil, ErrInvalidCIDR)
			},
		},
		{
			name: "remove failed",
			err:  errRemove,
			ip:   subnet,
			before: func(m *mock) {
				gomock.InOrder(
					m.parseCIDR(subnet).Return(subNet, nil),
					m.remove().Return(errRemove),
				)
			},
		},
		{
			name: "success",
			ip:   subnet,
			before: func(m *mock) {
				gomock.InOrder(
					m.parseCIDR(subnet).Return(subNet, nil),
					m.remove().Return(nil),
				)
			},
		},
	}

	list := New(addr, nil, nil)
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctx, ctrl)
			tt.before(m)
			list.storage = m.storage
			list.ip = m.ipService

			err := list.Remove(ctx, tt.ip)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestIPList_Contains(t *testing.T) {
	ctx := context.Background()
	errLoad := errors.New("load failed")
	table := []struct {
		name   string
		err    error
		exp    bool
		before func(m *mock)
	}{
		{
			name: "fail",
			err:  errLoad,
			before: func(m *mock) {
				m.load().Return(nil, errLoad)
			},
		},
		{
			name: "not found",
			before: func(m *mock) {
				m.load().Return([]string{}, nil)
			},
		},
		{
			name: "parse cidr failed",
			err:  ErrInvalidCIDR,
			before: func(m *mock) {
				m.load().Return([]string{ip}, ErrInvalidCIDR)
			},
		},
		{
			name: "ip is contains to the subnet",
			before: func(m *mock) {
				gomock.InOrder(
					m.load().Return([]string{subnet}, nil),
					m.parseCIDR(subnet).Return(subNet, nil),
				)
			},
			exp: true,
		},
		{
			name: "ip is not contains to the subnet",
			before: func(m *mock) {
				gomock.InOrder(
					m.load().Return([]string{subnet2, subnet3, subnet4}, nil),
					m.parseCIDR(subnet2).Return(subNet2, nil),
					m.parseCIDR(subnet3).Return(subNet3, nil),
					m.parseCIDR(subnet4).Return(subNet4, nil),
				)
			},
			exp: false,
		},
	}

	list := New(addr, nil, nil)
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctx, ctrl)
			tt.before(m)
			list.storage = m.storage
			list.ip = m.ipService

			ok, err := list.Contains(ctx, ip)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.exp, ok)
		})
	}
}
