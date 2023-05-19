package manager

import (
	"context"

	"github.com/rs/zerolog"
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
	logger    zerolog.Logger
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

func WithLogger(logger zerolog.Logger) Config {
	return func(app *App) {
		app.logger = logger
	}
}

func joinIPMask(ip, mask string) string {
	// todo: разобраться как работать с маской
	return ip
}

func (a *App) AddToBlackList(ctx context.Context, ip, mask string) error {
	if err := a.blackList.Add(ctx, joinIPMask(ip, mask)); err != nil {
		a.logger.Error().Err(err).Str("ip", ip).
			Str("mask", mask).Msg("add to black list failed")
		return err
	}

	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ip, mask string) error {
	if err := a.whiteList.Add(ctx, joinIPMask(ip, mask)); err != nil {
		a.logger.Error().Err(err).Str("ip", ip).
			Str("mask", mask).Msg("add to white list failed")
		return err
	}

	return nil
}

func (a *App) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	if err := a.blackList.Remove(ctx, joinIPMask(ip, mask)); err != nil {
		a.logger.Error().Err(err).Str("ip", ip).
			Str("mask", mask).Msg("remove from black list failed")
		return err
	}

	return nil
}

func (a *App) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	if err := a.whiteList.Remove(ctx, joinIPMask(ip, mask)); err != nil {
		a.logger.Error().Err(err).Str("ip", ip).
			Str("mask", mask).Msg("remove from white list failed")
		return err
	}

	return nil
}

func (a *App) Reset(ctx context.Context, login, ip string) error {
	if err := a.resetter.Reset(ctx, login, ip); err != nil {
		a.logger.Error().Err(err).Str("ip", ip).
			Msg("reset bucket by ip and login failed")
		return err
	}

	return nil
}
