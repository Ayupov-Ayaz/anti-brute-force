package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/alicebob/miniredis"

	servercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/server"

	limitercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/limiter"
	listcfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/list"
	loggercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/logger"
	rediscfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/redis"

	"github.com/ayupov-ayaz/anti-brute-force/cli/config"

	"github.com/stretchr/testify/require"

	run "github.com/ayupov-ayaz/anti-brute-force/cli/cmd"
	"github.com/gofiber/fiber/v2"
)

const (
	ipNet    = "192.1.1.64/27"
	maskedIP = "192.1.1.65"
)

func testConfig(t *testing.T) config.Config {
	t.Helper()
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())

	return config.Config{
		Server: servercfg.Server{Port: 8080},
		Redis: rediscfg.Redis{
			Addr: mr.Addr(),
		},
		Logger: loggercfg.Logger{
			Level: "INFO",
		},
		IPList: listcfg.IPList{
			BlackListAddr: "black-list",
			WhiteListAddr: "white-list",
		},
		Limiter: limitercfg.Limiter{
			Login: limitercfg.ValueLimiter{
				Count:    2,
				Interval: 1 * time.Minute,
			},
			Password: limitercfg.ValueLimiter{
				Count:    4,
				Interval: 1 * time.Minute,
			},
			IP: limitercfg.ValueLimiter{
				Count:    6,
				Interval: 1 * time.Minute,
			},
		},
	}
}

func TestApp(t *testing.T) {
	const (
		// routes
		addToBlackList      = "/black-list/add"
		addToWhiteList      = "/white-list/add"
		removeFromBlackList = "/black-list/remove"
		removeFromWhiteList = "/white-list/remove"
		check               = "/checker/check"
		reset               = "/buckets"
		// responses
		forbiddenResp = `{"ok":false}`
		okResp        = `{"ok":true}`
	)
	cfg := testConfig(t)
	server, err := run.MakeServer(cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})
	require.NoError(t, client.Ping(ctx).Err())

	app := fiber.New()
	server.Register(app)

	listBody := fmt.Sprintf(`{"ip_net":"%s"}`, ipNet)
	checkBody := func(login, pass, ip string) string {
		return fmt.Sprintf(`{"login":"%s","password":"%s","ip":"%s"}`, login, pass, ip)
	}

	resetBody := func(login, ip string) string {
		return fmt.Sprintf(`{"login":"%s","ip":"%s"}`, login, ip)
	}

	table := []struct {
		name    string
		method  string
		url     string
		body    string
		expCode int
		expBody string
	}{
		{
			name:    "add to black list",
			method:  http.MethodPost,
			url:     addToBlackList,
			body:    listBody,
			expCode: http.StatusOK,
			expBody: "OK",
		},
		{
			name:    "ip is in black list",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login", "pass", maskedIP),
			expCode: http.StatusForbidden,
			expBody: forbiddenResp,
		},
		{
			name:    "remove from black list",
			method:  http.MethodDelete,
			url:     removeFromBlackList,
			body:    listBody,
			expCode: http.StatusOK,
			expBody: "OK",
		},
		{
			name:    "check #1 (by login)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login", "pass", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "check #2 (by login)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login", "pass", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "check #3 (by login)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login", "pass", maskedIP),
			expCode: http.StatusForbidden,
			expBody: forbiddenResp,
		},
		{
			name:    "check #4 (by password)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login2", "pass", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "check #5 (by password)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login2", "pass", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "check #6 (by password)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login3", "pass", maskedIP),
			expCode: http.StatusForbidden,
			expBody: forbiddenResp,
		},
		{
			name:    "check #5 (by ip)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login3", "pass2", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "check #5 (by ip)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login4", "pass2", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "check #6 (by ip)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login4", "pass3", maskedIP),
			expCode: http.StatusForbidden,
			expBody: forbiddenResp,
		},
		{
			name:    "add ip to white list",
			method:  http.MethodPost,
			url:     addToWhiteList,
			body:    listBody,
			expCode: http.StatusOK,
			expBody: "OK",
		},
		{
			name:    "check #7 (ip is in white list)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login4", "pass3", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
		{
			name:    "remove from white list",
			method:  http.MethodDelete,
			url:     removeFromWhiteList,
			body:    listBody,
			expCode: http.StatusOK,
			expBody: "OK",
		},
		{
			name:    "check #7 (ip is in white list)",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login4", "pass3", maskedIP),
			expCode: http.StatusForbidden,
			expBody: forbiddenResp,
		},
		{
			name:    "reset buckets",
			method:  http.MethodDelete,
			url:     reset,
			body:    resetBody("login4", maskedIP),
			expCode: http.StatusOK,
			expBody: "OK",
		},
		{
			name:    "check #8",
			method:  http.MethodPost,
			url:     check,
			body:    checkBody("login4", "pass3", maskedIP),
			expCode: http.StatusOK,
			expBody: okResp,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(tt.body)))
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expBody, string(respBody))
		})
	}
}
