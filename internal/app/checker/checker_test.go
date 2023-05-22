package checker

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	"github.com/stretchr/testify/require"

	mocks "github.com/ayupov-ayaz/anti-brute-force/internal/app/checker/mocks"
	"github.com/golang/mock/gomock"
)

const (
	login           = "login"
	ip              = "127.0.0.1"
	pass            = "pass"
	ipUint32 uint32 = 2130706433
	ipStr           = "2130706433"
)

type mock struct {
	ctx       context.Context
	whiteList *mocks.MockIPList
	blackList *mocks.MockIPList
	checker   *mocks.MockChecker
	ipService *mocks.MockIPService
}

func newMock(ctx context.Context, ctrl *gomock.Controller) *mock {
	return &mock{
		ctx:       ctx,
		whiteList: mocks.NewMockIPList(ctrl),
		blackList: mocks.NewMockIPList(ctrl),
		checker:   mocks.NewMockChecker(ctrl),
		ipService: mocks.NewMockIPService(ctrl),
	}
}

func (m *mock) app() *App {
	return New(m.whiteList, m.blackList, m.ipService, m.checker)
}

func (m *mock) checkLogin(err error) *gomock.Call {
	return m.checker.EXPECT().AllowByLogin(m.ctx, login).Times(1).Return(err)
}

func (m *mock) checkPassword(err error) *gomock.Call {
	return m.checker.EXPECT().AllowByPassword(m.ctx, pass).Times(1).Return(err).
		After(m.checkLogin(nil))
}

func (m *mock) checkIP(err error) *gomock.Call {
	return m.checker.EXPECT().AllowByIP(m.ctx, ipStr).Times(1).Return(err).
		After(m.checkPassword(nil))
}

func (m *mock) parseIP(err error) *gomock.Call {
	resp := net.IP{127, 0, 0, 1}
	return m.ipService.EXPECT().ParseIP(ip).Times(1).Return(resp, err)
}

func (m *mock) ipToUint32() *gomock.Call {
	return m.ipService.EXPECT().IPToUint32(net.IP{127, 0, 0, 1}).Times(1).Return(ipUint32).
		After(m.parseIP(nil))
}

func (m *mock) containsWhiteList(ok bool, err error) *gomock.Call {
	return m.whiteList.EXPECT().Contains(m.ctx, ipStr).
		Times(1).Return(ok, err)
}

func (m *mock) containsBlackList(ok bool, err error) *gomock.Call {
	return m.blackList.EXPECT().Contains(m.ctx, ipStr).
		Times(1).Return(ok, err)
}

func TestApp_authIsAllowed(t *testing.T) {
	ctx := context.Background()
	errLogin := errors.New("login error")
	errPass := errors.New("password error")
	errIP := errors.New("ip error")

	table := []struct {
		name   string
		err    error
		before func(m *mock)
	}{
		{
			name: "login error",
			err:  errLogin,
			before: func(m *mock) {
				m.checkLogin(errLogin)
			},
		},
		{
			name: "password error",
			err:  errPass,
			before: func(m *mock) {
				m.checkPassword(errPass)
			},
		},
		{
			name: "ip error",
			err:  errIP,
			before: func(m *mock) {
				m.checkIP(errIP)
			},
		},
		{
			name: "success",
			before: func(m *mock) {
				m.checkIP(nil)
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctx, ctrl)
			tt.before(m)

			err := m.app().authIsAllowed(ctx, ipStr, login, pass)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestApp_parseIPKey(t *testing.T) {
	ctx := context.Background()
	errFailed := errors.New("parse failed")

	table := []struct {
		name   string
		err    error
		before func(m *mock)
	}{
		{
			name: "failed",
			err:  errFailed,
			before: func(m *mock) {
				m.parseIP(errFailed)
			},
		},
		{
			name: "success",
			before: func(m *mock) {
				m.ipToUint32()
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctx, ctrl)
			tt.before(m)

			val, err := m.app().parseIPKey(ip)
			require.ErrorIs(t, err, tt.err)
			if tt.err != nil {
				require.Empty(t, val)
			} else {
				require.Equal(t, ipStr, val)
			}
		})
	}
}

func TestApp_Check(t *testing.T) {
	errParseIP := errors.New("parse failed")
	errWhiteList := errors.New("white list failed")
	errBlackList := errors.New("black list failed")
	errCheck := errors.New("check failed")

	ctx := context.Background()
	table := []struct {
		name   string
		err    error
		before func(m *mock)
	}{
		{
			name: "parse ip failed",
			err:  errParseIP,
			before: func(m *mock) {
				m.parseIP(errParseIP)
			},
		},
		{
			name: "check white list failed",
			err:  errWhiteList,
			before: func(m *mock) {
				gomock.InOrder(
					m.ipToUint32(),
					m.containsWhiteList(false, errWhiteList),
				)
			},
		},
		{
			name: "user is in white list",
			before: func(m *mock) {
				gomock.InOrder(
					m.ipToUint32(),
					m.containsWhiteList(true, nil),
				)
			},
		},
		{
			name: "check black list failed",
			err:  errBlackList,
			before: func(m *mock) {
				gomock.InOrder(
					m.ipToUint32(),
					m.containsWhiteList(false, nil),
					m.containsBlackList(false, errBlackList),
				)
			},
		},
		{
			name: "user is in black list",
			err:  apperr.ErrUserIsBlocked,
			before: func(m *mock) {
				gomock.InOrder(
					m.ipToUint32(),
					m.containsWhiteList(false, nil),
					m.containsBlackList(true, nil),
				)
			},
		},
		{
			name: "check failed",
			err:  errCheck,
			before: func(m *mock) {
				gomock.InOrder(
					m.ipToUint32(),
					m.containsWhiteList(false, nil),
					m.containsBlackList(false, nil),
					m.checkIP(errCheck),
				)
			},
		},
		{
			name: "success",
			before: func(m *mock) {
				gomock.InOrder(
					m.ipToUint32(),
					m.containsWhiteList(false, nil),
					m.containsBlackList(false, nil),
					m.checkIP(nil),
				)
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctx, ctrl)
			tt.before(m)
			err := m.app().Check(ctx, ip, login, pass)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
