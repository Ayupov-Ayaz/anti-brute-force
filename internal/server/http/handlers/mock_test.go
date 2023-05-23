package handlers

import (
	"context"

	mocks "github.com/ayupov-ayaz/anti-brute-force/internal/server/http/handlers/mocks"
	"github.com/golang/mock/gomock"
)

type mock struct {
	ctx       context.Context
	validator *mocks.MockValidator
	checker   *mocks.MockChecker
	decoder   *mocks.MockDecoder
	manager   *mocks.MockManager
}

func newMock(ctrl *gomock.Controller) *mock {
	return &mock{
		validator: mocks.NewMockValidator(ctrl),
		checker:   mocks.NewMockChecker(ctrl),
		decoder:   mocks.NewMockDecoder(ctrl),
		manager:   mocks.NewMockManager(ctrl),
	}
}

type unmarshalCallback func(data []byte, v interface{}) error

func (m *mock) unmarshal(data []byte, callback unmarshalCallback) *gomock.Call {
	return m.decoder.EXPECT().Unmarshal(data, gomock.Any()).Times(1).DoAndReturn(callback)
}

func (m *mock) marshal(i interface{}) *gomock.Call {
	return m.decoder.EXPECT().Marshal(i).Times(1)
}

func (m *mock) validate(i interface{}) *gomock.Call {
	return m.validator.EXPECT().Validate(i).Times(1)
}

func (m *mock) check(ip, login, pass string) *gomock.Call {
	return m.checker.EXPECT().Check(gomock.Any(), ip, login, pass).Times(1)
}

func (m *mock) addToBlackList(ip, mask string) *gomock.Call {
	return m.manager.EXPECT().AddToBlackList(gomock.Any(), ip, mask).Times(1)
}

func (m *mock) addToWhiteList(ip, mask string) *gomock.Call {
	return m.manager.EXPECT().AddToWhiteList(gomock.Any(), ip, mask).Times(1)
}

func (m *mock) removeFromBlackList(ip, mask string) *gomock.Call {
	return m.manager.EXPECT().RemoveFromBlackList(gomock.Any(), ip, mask).Times(1)
}

func (m *mock) removeFromWhiteList(ip, mask string) *gomock.Call {
	return m.manager.EXPECT().RemoveFromWhiteList(gomock.Any(), ip, mask).Times(1)
}

func (m *mock) reset(login, ip string) *gomock.Call {
	return m.manager.EXPECT().Reset(gomock.Any(), login, ip).Times(1)
}
