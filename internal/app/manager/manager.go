package manager

import (
	"context"

	"go.uber.org/zap"
)

// todo: log

type IPList interface {
	Add(ctx context.Context, ip string) error
	Remove(ctx context.Context, ip string) error
}

type Resetter interface {
	Reset(ctx context.Context, login, ip string) error
}

type App struct {
	blackList IPList
	whiteList IPList
	resetter  Resetter
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

func WithBlackList(blackList IPList) Config {
	return func(app *App) {
		app.blackList = blackList
	}
}

func WithWhiteList(whiteList IPList) Config {
	return func(app *App) {
		app.whiteList = whiteList
	}
}

func WithResetter(resetter Resetter) Config {
	return func(app *App) {
		app.resetter = resetter
	}
}

func WithLogger(logger *zap.Logger) Config {
	return func(app *App) {
		app.logger = logger
	}
}

func joinIPMask(ip, mask string) string {
	return ip + "/" + mask
}

func (a *App) AddToBlackList(ctx context.Context, ip, mask string) error {
	key := joinIPMask(ip, mask)

	if err := a.blackList.Add(ctx, key); err != nil {
		a.logger.Error("add to black list failed",
			zap.String("ip/mask", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ip, mask string) error {
	key := joinIPMask(ip, mask)

	if err := a.whiteList.Add(ctx, key); err != nil {
		a.logger.Error("add to white list failed",
			zap.String("ip/mask", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (a *App) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	key := joinIPMask(ip, mask)

	if err := a.blackList.Remove(ctx, key); err != nil {
		a.logger.Error("remove from black list failed",
			zap.String("ip/mask", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (a *App) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	key := joinIPMask(ip, mask)

	if err := a.whiteList.Remove(ctx, key); err != nil {
		a.logger.Error("remove from white list failed",
			zap.String("ip/mask", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (a *App) Reset(ctx context.Context, login, ip string) error {
	if err := a.resetter.Reset(ctx, login, ip); err != nil {
		a.logger.Error("reset failed",
			zap.String("ip", ip),
			zap.String("login", login),
			zap.Error(err))
	}

	return nil
}
