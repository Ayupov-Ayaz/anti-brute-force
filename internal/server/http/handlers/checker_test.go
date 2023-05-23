package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	"github.com/stretchr/testify/require"

	"github.com/ayupov-ayaz/anti-brute-force/internal/modules/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
)

const (
	successResp   = `{"ok":true}`
	forbiddenResp = `{"ok":false}`
)

func TestResponse(t *testing.T) {
	successRespBody, err := jsoniter.Marshal(Response{Ok: true})
	require.NoError(t, err)
	require.Equal(t, []byte(successResp), successRespBody)

	forbiddenRespBody, err := jsoniter.Marshal(Response{Ok: false})
	require.NoError(t, err)
	require.Equal(t, []byte(forbiddenResp), forbiddenRespBody)
}

func TestCheckerHTTP_Check(t *testing.T) {
	const (
		ip    = "127.0.0.1"
		login = "login"
		pass  = "pass"
	)
	const route = "/checker/check"

	errUnmarshal := errors.New("unmarshal failed")
	errValidate := errors.New("validation failed")
	errCheck := errors.New("check failed")
	errMarshal := errors.New("marshal failed")

	reqBody := []byte(`{json}`)
	successUnmarshalCallBack := func(data []byte, v interface{}) error {
		auth, ok := v.(*CheckAuthRequest)
		if !ok {
			return fmt.Errorf("invalid type %T", v)
		}

		auth.IP = ip
		auth.Login = login
		auth.Pass = pass
		return nil
	}

	model := CheckAuthRequest{
		BaseRequest: BaseRequest{
			IP:    ip,
			Login: login,
		},
		Pass: pass,
	}

	okResp := Response{Ok: true}
	notOkResp := Response{Ok: false}

	table := []struct {
		name    string
		expCode int
		expBody string
		before  func(m *mock)
	}{
		{
			name: "unmarshal failed",
			before: func(m *mock) {
				m.unmarshal(reqBody, func(data []byte, v interface{}) error {
					return errUnmarshal
				})
			},
			expCode: http.StatusInternalServerError,
			expBody: errUnmarshal.Error(),
		},
		{
			name: "validation failed",
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallBack),
					m.validate(model).Return(errValidate),
				)
			},
			expCode: http.StatusInternalServerError,
			expBody: errValidate.Error(),
		},
		{
			name: "check failed",
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallBack),
					m.validate(model).Return(nil),
					m.check(ip, login, pass).Return(errCheck),
				)
			},
			expCode: http.StatusInternalServerError,
			expBody: errCheck.Error(),
		},
		{
			name: "marshal failed",
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallBack),
					m.validate(model).Return(nil),
					m.check(ip, login, pass).Return(nil),
					m.marshal(okResp).Return(nil, errMarshal),
				)
			},
			expCode: http.StatusInternalServerError,
			expBody: errMarshal.Error(),
		},
		{
			name:    "success",
			expCode: http.StatusOK,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallBack),
					m.validate(model).Return(nil),
					m.check(ip, login, pass).Return(nil),
					m.marshal(okResp).Return([]byte(successResp), nil),
				)
			},
			expBody: successResp,
		},
		{
			name:    "user is blocked",
			expCode: http.StatusForbidden,
			before: func(m *mock) {
				gomock.InOrder(
					m.unmarshal(reqBody, successUnmarshalCallBack),
					m.validate(model).Return(nil),
					m.check(ip, login, pass).Return(apperr.ErrUserIsBlocked),
					m.marshal(notOkResp).Return([]byte(forbiddenResp), nil),
				)
			},
			expBody: forbiddenResp,
		},
	}

	app := fiber.New()
	zLogger, err := logger.New("DISABLED")
	require.NoError(t, err)
	if err != nil {

	}
	handler := NewChecker(nil, nil, nil, zLogger)
	handler.Register(app)

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			handler.validator = m.validator
			handler.app = m.checker
			handler.decoder = m.decoder

			tt.before(m)

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
