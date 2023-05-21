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
	ip       = net.IP([]byte{192, 168, 0, 1})
	mask     = net.IPMask([]byte{255, 255, 255, 0})
	maskedIP = net.IP([]byte{192, 168, 0, 0})
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

func (m *mocks) parseMaskedIP(err error) *gomock.Call {
	return m.ip.EXPECT().ParseMaskedIP(ip.String(), mask.String()).Return(maskedIP, err)
}

func (m *mocks) maskedIPToUint32() *gomock.Call {
	return m.ip.EXPECT().IPToUint32(maskedIP).Return(ipToUint32)
}

func (m *mocks) ipToUint32() *gomock.Call {
	return m.ip.EXPECT().IPToUint32(ip).Return(ipToUint32)
}

func (m *mocks) addToList(list *mock.MockIPList, err error) *gomock.Call {
	return list.EXPECT().Add(m.ctx, fmt.Sprintf("%d", ipToUint32)).Return(err).
		After(m.maskedIPToUint32())
}

func (m *mocks) removeFromList(list *mock.MockIPList, err error) *gomock.Call {
	return list.EXPECT().Remove(m.ctx, fmt.Sprintf("%d", ipToUint32)).Return(err).
		After(m.maskedIPToUint32())
}

func (m *mocks) parseIP(err error) *gomock.Call {
	return m.ip.EXPECT().ParseIP(ip.String()).Return(ip, err)
}

func (m *mocks) reset(err error) *gomock.Call {
	return m.resetter.EXPECT().Reset(m.ctx, login, fmt.Sprintf("%d", ipToUint32)).Return(err).
		After(m.ipToUint32())
}

func TestApp_makeMaskedKey(t *testing.T) {
	ctx := context.Background()
	errParseMaskedIP := errors.New("parse masked ip error")

	tests := []struct {
		name   string
		before func(m *mocks)
		err    error
	}{
		{
			name: "parse masked ip error",
			err:  errParseMaskedIP,
			before: func(m *mocks) {
				m.parseMaskedIP(errParseMaskedIP)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.maskedIPToUint32(),
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
			key, err := app.makeMaskedKey(ip.String(), mask.String())
			require.Equal(t, tt.err, err)
			if tt.err == nil {
				require.Equal(t, fmt.Sprintf("%d", ipToUint32), key)
			}
		})
	}
}

func TestApp_AddToBlackList(t *testing.T) {
	ctx := context.Background()
	errParseMaskedIP := errors.New("parse masked ip error")
	errAddToBlackList := errors.New("add to black list error")

	tests := []struct {
		name   string
		before func(m *mocks)
		err    error
	}{
		{
			name: "parse masked ip error",
			err:  errParseMaskedIP,
			before: func(m *mocks) {
				m.parseMaskedIP(errParseMaskedIP)
			},
		},
		{
			name: "add to black list error",
			err:  errAddToBlackList,
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.addToList(m.blackList, errAddToBlackList),
				)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.addToList(m.blackList, nil),
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

			a := New(m.whiteList, m.blackList, m.ip, m.resetter)

			err := a.AddToBlackList(ctx, ip.String(), mask.String())
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestApp_AddToWhiteList(t *testing.T) {
	ctx := context.Background()
	errParseMaskedIP := errors.New("parse masked ip error")
	addToWhiteListFailed := errors.New("add to white list error")

	tests := []struct {
		name   string
		before func(m *mocks)
		err    error
	}{
		{
			name: "parse masked ip error",
			err:  errParseMaskedIP,
			before: func(m *mocks) {
				m.parseMaskedIP(errParseMaskedIP)
			},
		},
		{
			name: "add to white list error",
			err:  addToWhiteListFailed,
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.addToList(m.whiteList, addToWhiteListFailed),
				)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.addToList(m.whiteList, nil),
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

			a := New(m.whiteList, m.blackList, m.ip, m.resetter)

			err := a.AddToWhiteList(ctx, ip.String(), mask.String())
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestApp_RemoveFromBlackList(t *testing.T) {
	ctx := context.Background()

	errParseIP := errors.New("parse ip error")
	errRemoveFromList := errors.New("remove from list error")

	tests := []struct {
		name   string
		before func(m *mocks)
		err    error
	}{
		{
			name: "make masked key error",
			err:  errParseIP,
			before: func(m *mocks) {
				m.parseMaskedIP(errParseIP)
			},
		},
		{
			name: "remove from black list failed",
			err:  errRemoveFromList,
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.removeFromList(m.blackList, errRemoveFromList),
				)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.removeFromList(m.blackList, nil),
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
			err := app.RemoveFromBlackList(ctx, ip.String(), mask.String())
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestApp_RemoveFromWhiteList(t *testing.T) {
	ctx := context.Background()

	errParseIP := errors.New("parse ip error")
	errRemoveFromList := errors.New("remove from list error")

	tests := []struct {
		name   string
		before func(m *mocks)
		err    error
	}{
		{
			name: "make masked key error",
			err:  errParseIP,
			before: func(m *mocks) {
				m.parseMaskedIP(errParseIP)
			},
		},
		{
			name: "remove from black list failed",
			err:  errRemoveFromList,
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.removeFromList(m.whiteList, errRemoveFromList),
				)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseMaskedIP(nil),
					m.removeFromList(m.whiteList, nil),
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
			err := app.RemoveFromWhiteList(ctx, ip.String(), mask.String())
			require.ErrorIs(t, err, tt.err)
		})
	}
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
				m.parseIP(errParseIP)
			},
		},
		{
			name: "reset error",
			err:  errReset,
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseIP(nil),
					m.reset(errReset),
				)
			},
		},
		{
			name: "success",
			before: func(m *mocks) {
				gomock.InOrder(
					m.parseIP(nil),
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
