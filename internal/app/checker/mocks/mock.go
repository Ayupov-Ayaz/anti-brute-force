// Code generated by MockGen. DO NOT EDIT.
// Source: checker.go

// Package mock_checker is a generated GoMock package.
package mock_checker

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPList is a mock of IPList interface.
type MockIPList struct {
	ctrl     *gomock.Controller
	recorder *MockIPListMockRecorder
}

// MockIPListMockRecorder is the mock recorder for MockIPList.
type MockIPListMockRecorder struct {
	mock *MockIPList
}

// NewMockIPList creates a new mock instance.
func NewMockIPList(ctrl *gomock.Controller) *MockIPList {
	mock := &MockIPList{ctrl: ctrl}
	mock.recorder = &MockIPListMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPList) EXPECT() *MockIPListMockRecorder {
	return m.recorder
}

// Contains mocks base method.
func (m *MockIPList) Contains(ctx context.Context, ip string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Contains", ctx, ip)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Contains indicates an expected call of Contains.
func (mr *MockIPListMockRecorder) Contains(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Contains", reflect.TypeOf((*MockIPList)(nil).Contains), ctx, ip)
}

// MockChecker is a mock of Checker interface.
type MockChecker struct {
	ctrl     *gomock.Controller
	recorder *MockCheckerMockRecorder
}

// MockCheckerMockRecorder is the mock recorder for MockChecker.
type MockCheckerMockRecorder struct {
	mock *MockChecker
}

// NewMockChecker creates a new mock instance.
func NewMockChecker(ctrl *gomock.Controller) *MockChecker {
	mock := &MockChecker{ctrl: ctrl}
	mock.recorder = &MockCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChecker) EXPECT() *MockCheckerMockRecorder {
	return m.recorder
}

// AllowByIP mocks base method.
func (m *MockChecker) AllowByIP(ctx context.Context, login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowByIP", ctx, login)
	ret0, _ := ret[0].(error)
	return ret0
}

// AllowByIP indicates an expected call of AllowByIP.
func (mr *MockCheckerMockRecorder) AllowByIP(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowByIP", reflect.TypeOf((*MockChecker)(nil).AllowByIP), ctx, login)
}

// AllowByLogin mocks base method.
func (m *MockChecker) AllowByLogin(ctx context.Context, login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowByLogin", ctx, login)
	ret0, _ := ret[0].(error)
	return ret0
}

// AllowByLogin indicates an expected call of AllowByLogin.
func (mr *MockCheckerMockRecorder) AllowByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowByLogin", reflect.TypeOf((*MockChecker)(nil).AllowByLogin), ctx, login)
}

// AllowByPassword mocks base method.
func (m *MockChecker) AllowByPassword(ctx context.Context, login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowByPassword", ctx, login)
	ret0, _ := ret[0].(error)
	return ret0
}

// AllowByPassword indicates an expected call of AllowByPassword.
func (mr *MockCheckerMockRecorder) AllowByPassword(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowByPassword", reflect.TypeOf((*MockChecker)(nil).AllowByPassword), ctx, login)
}

// MockIPService is a mock of IPService interface.
type MockIPService struct {
	ctrl     *gomock.Controller
	recorder *MockIPServiceMockRecorder
}

// MockIPServiceMockRecorder is the mock recorder for MockIPService.
type MockIPServiceMockRecorder struct {
	mock *MockIPService
}

// NewMockIPService creates a new mock instance.
func NewMockIPService(ctrl *gomock.Controller) *MockIPService {
	mock := &MockIPService{ctrl: ctrl}
	mock.recorder = &MockIPServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPService) EXPECT() *MockIPServiceMockRecorder {
	return m.recorder
}

// IPToUint32 mocks base method.
func (m *MockIPService) IPToUint32(ip net.IP) uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IPToUint32", ip)
	ret0, _ := ret[0].(uint32)
	return ret0
}

// IPToUint32 indicates an expected call of IPToUint32.
func (mr *MockIPServiceMockRecorder) IPToUint32(ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IPToUint32", reflect.TypeOf((*MockIPService)(nil).IPToUint32), ip)
}

// ParseIP mocks base method.
func (m *MockIPService) ParseIP(ip string) (net.IP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseIP", ip)
	ret0, _ := ret[0].(net.IP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseIP indicates an expected call of ParseIP.
func (mr *MockIPServiceMockRecorder) ParseIP(ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseIP", reflect.TypeOf((*MockIPService)(nil).ParseIP), ip)
}
