package checker

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrUserIsBlocked = errors.New("user is blocked")
)

type IPList interface {
	Contains(ctx context.Context, ip string) (bool, error)
}

type Checker interface {
	Check(ctx context.Context, ip, login, pass string) error
}

type App struct {
	whiteList IPList
	blackList IPList
	buckets   Checker
	logger    *zap.Logger
}

type Config func(app *App)

func New(configs ...Config) *App {
	app := &App{}
	for _, config := range configs {
		config(app)
	}
	return app
}

func WithWhiteList(whiteList IPList) Config {
	return func(app *App) {
		app.whiteList = whiteList
	}
}

func WithBlackList(blackList IPList) Config {
	return func(app *App) {
		app.blackList = blackList
	}
}

func WithBuckets(buckets Checker) Config {
	return func(app *App) {
		app.buckets = buckets
	}
}

func WithLogger(logger *zap.Logger) Config {
	return func(app *App) {
		app.logger = logger
	}
}

func (a *App) Check(ctx context.Context, ip, login, pass string) error {
	logger := a.logger.Named("check").
		With(zap.String("ip", ip),
			zap.String("login", login))

	ok, err := a.whiteList.Contains(ctx, ip)
	if err != nil {
		logger.Error("error while checking ip in white list", zap.Error(err))
		return err
	}

	if ok {
		a.logger.Info("ip is in white list")
		return nil
	}

	ok, err = a.blackList.Contains(ctx, ip)
	if err != nil {
		logger.Error("error while checking ip in black list", zap.Error(err))
		return err
	}

	if ok {
		logger.Info("ip is in black list")
		return ErrUserIsBlocked
	}

	if err := a.buckets.Check(ctx, ip, login, pass); err != nil {
		logger.Error("error while checking buckets", zap.Error(err))
		return err
	}

	return nil
}
