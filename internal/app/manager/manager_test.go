package manager

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/ayupov-ayaz/anti-brute-force/internal/app/manager/mocks"
	"github.com/golang/mock/gomock"
)

const (
	ipToUint32 uint32 = 3232235520
	login             = "login"
)

var (
	_, ipNet, _ = net.ParseCIDR("192.1.1.0/26")
	ip          = net.ParseIP("192.1.1.65")
)

type mocks struct {
	ctx       context.Context
	blackList *mock.MockIPList
	whiteList *mock.MockIPList
	ip        *mock.MockIPService
	resetter  *mock.MockResetter
}

func new(ctx context.Context, ctrl *gomock.Controller) *mocks {
	return &mocks{
		ctx:       ctx,
		blackList: mock.NewMockIPList(ctrl),
		whiteList: mock.NewMockIPList(ctrl),
		ip:        mock.NewMockIPService(ctrl),
		resetter:  mock.NewMockResetter(ctrl),
	}
}
func (m *mocks) ipToUint32() *gomock.Call {
	return m.ip.EXPECT().IPToUint32(ip).Return(ipToUint32)
}

func (m *mocks) parseCIDR() *gomock.Call {
	return m.ip.EXPECT().ParseCIDR(ipNet.String())
}

func (m *mocks) parseIP() *gomock.Call {
	return m.ip.EXPECT().ParseIP(ip.String())
}

func (m *mocks) reset(err error) *gomock.Call {
	return m.resetter.EXPECT().Reset(m.ctx, login, fmt.Sprintf("%d", ipToUint32)).Return(err).
		After(m.ipToUint32())
}

func TestApp_Reset(t *testing.T) {
	ctx := context.Background()

	errParseIP := errors.New("parse ip error")
	errReset := errors.New("reset error")

	tests := []struct {
		name   string
		err    error
		before func(m *mocks)
	}{
		{
			name: "parse ip error",
			err:  errParseIP,
			before: func(m *mocks) {
				m.parseIP().Return(nil, errParseIP)
			},
		},
		{
			name: "reset error",
			err:  errReset,
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseIP().Return(ip, nil),
					m.reset(errReset),
				)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseIP().Return(ip, nil),
					m.reset(nil),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := new(ctx, ctrl)
			tt.before(m)

			app := New(m.whiteList, m.blackList, m.ip, m.resetter)
			err := app.Reset(ctx, login, ip.String())
			require.ErrorIs(t, err, tt.err)
		})
	}
}
