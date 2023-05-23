// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// AddToBlackList mocks base method.
func (m *MockManager) AddToBlackList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBlackList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBlackList indicates an expected call of AddToBlackList.
func (mr *MockManagerMockRecorder) AddToBlackList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBlackList", reflect.TypeOf((*MockManager)(nil).AddToBlackList), ctx, ip, mask)
}

// AddToWhiteList mocks base method.
func (m *MockManager) AddToWhiteList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToWhiteList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToWhiteList indicates an expected call of AddToWhiteList.
func (mr *MockManagerMockRecorder) AddToWhiteList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToWhiteList", reflect.TypeOf((*MockManager)(nil).AddToWhiteList), ctx, ip, mask)
}

// RemoveFromBlackList mocks base method.
func (m *MockManager) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromBlackList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromBlackList indicates an expected call of RemoveFromBlackList.
func (mr *MockManagerMockRecorder) RemoveFromBlackList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromBlackList", reflect.TypeOf((*MockManager)(nil).RemoveFromBlackList), ctx, ip, mask)
}

// RemoveFromWhiteList mocks base method.
func (m *MockManager) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromWhiteList", ctx, ip, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromWhiteList indicates an expected call of RemoveFromWhiteList.
func (mr *MockManagerMockRecorder) RemoveFromWhiteList(ctx, ip, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromWhiteList", reflect.TypeOf((*MockManager)(nil).RemoveFromWhiteList), ctx, ip, mask)
}

// Reset mocks base method.
func (m *MockManager) Reset(ctx context.Context, login, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", ctx, login, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockManagerMockRecorder) Reset(ctx, login, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockManager)(nil).Reset), ctx, login, ip)
}
