package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

const (
	ip    = "192.168.23.3"
	mask  = "255.255.255.0.0"
	login = "login"
)

var successUnmarshalIPCallBack = func(data []byte, v interface{}) error {
	model, ok := v.(*IP)
	if !ok {
		return fmt.Errorf("invalid type %T", v)
	}

	model.IP = ip
	model.Mask = mask

	return nil
}

func TestManagerHTTP_addToBlackList(t *testing.T) {
	const route = "/black-list/add"
	reqBody := []byte(`{json}`)

	model := IP{
		IP:   ip,
		Mask: mask,
	}

	errUnmarshal := errors.New("unmarshal failed")
	errValidate := errors.New("validation failed")
	errAdd := errors.New("add failed")

	table := []struct {
		name    string
		expBody string
		expCode int
		before  func(m *mock)
	}{
		{
			name:    "unmarshal failed",
			expCode: http.StatusInternalServerError,
			expBody: errUnmarshal.Error(),
			before: func(m *mock) {
				m.unmarshal(reqBody, func(data []byte, v interface{}) error {
					return errUnmarshal
				})
			},
		},
		{
			name:    "validation failed",
			expCode: http.StatusInternalServerError,
			expBody: errValidate.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(errValidate),
				)
			},
		},
		{
			name:    "add failed",
			expCode: http.StatusInternalServerError,
			expBody: errAdd.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.addToBlackList(ip, mask).Return(errAdd),
				)
			},
		},
		{
			name:    "success",
			expCode: http.StatusOK,
			expBody: "OK",
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.addToBlackList(ip, mask).Return(nil),
				)
			},
		},
	}

	app := fiber.New()
	zLogger, err := logger.New("DISABLED")
	require.NoError(t, err)
	handler := NewManager(nil, nil, nil, zLogger)
	handler.Register(app)

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)
			handler.validator = m.validator
			handler.decoder = m.decoder
			handler.app = m.manager

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(reqBody))
			resp, err := app.Test(req, 1)
			require.NoError(t, err)
			require.Equal(t, tt.expCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expBody, string(respBody))
		})
	}
}

func TestManagerHTTP_addToWhiteList(t *testing.T) {
	const route = "/white-list/add"

	reqBody := []byte(`{json}`)
	model := IP{
		IP:   ip,
		Mask: mask,
	}

	errUnmarshal := errors.New("unmarshal failed")
	errValidate := errors.New("validation failed")
	errAdd := errors.New("add failed")

	table := []struct {
		name    string
		expBody string
		expCode int
		before  func(m *mock)
	}{
		{
			name:    "unmarshal failed",
			expCode: http.StatusInternalServerError,
			expBody: errUnmarshal.Error(),
			before: func(m *mock) {
				m.unmarshal(reqBody, func(data []byte, v interface{}) error {
					return errUnmarshal
				})
			},
		},
		{
			name:    "validation failed",
			expCode: http.StatusInternalServerError,
			expBody: errValidate.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(errValidate),
				)
			},
		},
		{
			name:    "add failed",
			expCode: http.StatusInternalServerError,
			expBody: errAdd.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.addToWhiteList(ip, mask).Return(errAdd),
				)
			},
		},
		{
			name:    "success",
			expCode: http.StatusOK,
			expBody: "OK",
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.addToWhiteList(ip, mask).Return(nil),
				)
			},
		},
	}

	app := fiber.New()
	zLogger, err := logger.New("DISABLED")
	require.NoError(t, err)
	handler := NewManager(nil, nil, nil, zLogger)
	handler.Register(app)

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)
			handler.validator = m.validator
			handler.decoder = m.decoder
			handler.app = m.manager

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(reqBody))
			resp, err := app.Test(req, 1)
			require.NoError(t, err)
			require.Equal(t, tt.expCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expBody, string(respBody))
		})
	}
}

func TestManagerHTTP_removeFromBlackList(t *testing.T) {
	reqBody := []byte(`{json}`)

	errUnmarshal := errors.New("unmarshal failed")
	errValidate := errors.New("validation failed")
	errRemove := errors.New("remove failed")

	model := IP{
		IP:   ip,
		Mask: mask,
	}

	tests := []struct {
		name    string
		expCode int
		expBody string
		before  func(m *mock)
	}{
		{
			name:    "unmarshal failed",
			expBody: errUnmarshal.Error(),
			expCode: http.StatusInternalServerError,
			before: func(m *mock) {
				m.unmarshal(reqBody, func(data []byte, v interface{}) error {
					return errUnmarshal
				})
			},
		},
		{
			name:    "validation failed",
			expBody: errValidate.Error(),
			expCode: http.StatusInternalServerError,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(errValidate),
				)
			},
		},
		{
			name:    "remove failed",
			expCode: http.StatusInternalServerError,
			expBody: errRemove.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.removeFromBlackList(ip, mask).Return(errRemove),
				)
			},
		},
		{
			name:    "success",
			expBody: "OK",
			expCode: http.StatusOK,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.removeFromBlackList(ip, mask).Return(nil),
				)
			},
		},
	}

	const route = "/black-list/remove"
	app := fiber.New()
	zLogger, err := logger.New("DISABLED")
	require.NoError(t, err)
	handler := NewManager(nil, nil, nil, zLogger)
	handler.Register(app)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := newMock(ctrl)
			tt.before(m)

			handler.validator = m.validator
			handler.decoder = m.decoder
			handler.app = m.manager

			req := httptest.NewRequest(http.MethodDelete, route, bytes.NewBuffer(reqBody))
			resp, err := app.Test(req, 1)
			require.NoError(t, err)
			require.Equal(t, tt.expCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expBody, string(respBody))
		})
	}
}

func TestManagerHTTP_removeFromWhiteList(t *testing.T) {
	reqBody := []byte(`{json}`)

	errUnmarshal := errors.New("unmarshal failed")
	errValidate := errors.New("validation failed")
	errRemove := errors.New("remove failed")

	model := IP{
		IP:   ip,
		Mask: mask,
	}

	tests := []struct {
		name    string
		expCode int
		expBody string
		before  func(m *mock)
	}{
		{
			name:    "unmarshal failed",
			expBody: errUnmarshal.Error(),
			expCode: http.StatusInternalServerError,
			before: func(m *mock) {
				m.unmarshal(reqBody, func(data []byte, v interface{}) error {
					return errUnmarshal
				})
			},
		},
		{
			name:    "validation failed",
			expBody: errValidate.Error(),
			expCode: http.StatusInternalServerError,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(errValidate),
				)
			},
		},
		{
			name:    "remove failed",
			expCode: http.StatusInternalServerError,
			expBody: errRemove.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.removeFromWhiteList(ip, mask).Return(errRemove),
				)
			},
		},
		{
			name:    "success",
			expBody: "OK",
			expCode: http.StatusOK,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalIPCallBack),
					m.validate(model).Return(nil),
					m.removeFromWhiteList(ip, mask).Return(nil),
				)
			},
		},
	}

	const route = "/white-list/remove"
	app := fiber.New()
	zLogger, err := logger.New("DISABLED")
	require.NoError(t, err)
	handler := NewManager(nil, nil, nil, zLogger)
	handler.Register(app)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := newMock(ctrl)
			tt.before(m)

			handler.validator = m.validator
			handler.decoder = m.decoder
			handler.app = m.manager

			req := httptest.NewRequest(http.MethodDelete, route, bytes.NewBuffer(reqBody))
			resp, err := app.Test(req, 1)
			require.NoError(t, err)
			require.Equal(t, tt.expCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expBody, string(respBody))
		})
	}
}

func TestManagerHTTP_reset(t *testing.T) {
	const route = "/buckets"

	reqBody := []byte(`{json}`)

	errUnmarshal := errors.New("unmarshal failed")
	errValidate := errors.New("validation failed")
	errReset := errors.New("reset failed")

	var successUnmarshalCallback = func(data []byte, v interface{}) error {
		req, ok := v.(*BaseRequest)
		if !ok {
			return fmt.Errorf("invalid type %v", v)
		}

		req.IP = ip
		req.Login = login

		return nil
	}

	model := BaseRequest{
		IP:    ip,
		Login: login,
	}

	tests := []struct {
		name    string
		expCode int
		expBody string
		before  func(m *mock)
	}{
		{
			name:    "unmarshal failed",
			expBody: errUnmarshal.Error(),
			expCode: http.StatusInternalServerError,
			before: func(m *mock) {
				m.unmarshal(reqBody, func(data []byte, v interface{}) error {
					return errUnmarshal
				})
			},
		},
		{
			name:    "validation failed",
			expBody: errValidate.Error(),
			expCode: http.StatusInternalServerError,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallback),
					m.validate(model).Return(errValidate),
				)
			},
		},
		{
			name:    "reset failed",
			expCode: http.StatusInternalServerError,
			expBody: errReset.Error(),
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallback),
					m.validate(model).Return(nil),
					m.reset(login, ip).Return(errReset),
				)
			},
		},
		{
			name:    "success",
			expBody: "OK",
			expCode: http.StatusOK,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallback),
					m.validate(model).Return(nil),
					m.reset(login, ip).Return(nil),
				)
			},
		},
	}

	app := fiber.New()
	zLogger, err := logger.New("DISABLED")
	require.NoError(t, err)
	handler := NewManager(nil, nil, nil, zLogger)
	handler.Register(app)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := newMock(ctrl)
			tt.before(m)

			handler.validator = m.validator
			handler.decoder = m.decoder
			handler.app = m.manager

			req := httptest.NewRequest(http.MethodDelete, route, bytes.NewBuffer(reqBody))
			resp, err := app.Test(req, 1)
			require.NoError(t, err)
			require.Equal(t, tt.expCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expBody, string(respBody))
		})
	}
}
