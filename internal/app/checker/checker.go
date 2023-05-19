package checker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"
)

type IPList interface {
	Contains(ctx context.Context, ip string) (bool, error)
}

type Checker interface {
	AllowByLogin(ctx context.Context, login string) error
	AllowByPassword(ctx context.Context, login string) error
	AllowByIP(ctx context.Context, login string) error
}

type App struct {
	whiteList IPList
	blackList IPList
	checker   Checker
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

func WithCheckers(checker Checker) Config {
	return func(app *App) {
		app.checker = checker
	}
}

func WithLogger(logger zerolog.Logger) Config {
	return func(app *App) {
		app.logger = logger
	}
}

func (a *App) authIsAllowed(ctx context.Context, ip, login, pass string) error {
	if err := a.checker.AllowByLogin(ctx, login); err != nil {
		return fmt.Errorf("check login: %w", err)
	}

	if err := a.checker.AllowByPassword(ctx, pass); err != nil {
		return fmt.Errorf("check password: %w", err)
	}

	if err := a.checker.AllowByIP(ctx, ip); err != nil {
		return fmt.Errorf("check ip: %w", err)
	}

	return nil
}

func (a *App) Check(ctx context.Context, ip, login, pass string) error {
	logger := a.logger.With().Str("ip", ip).Str("login", login).Logger()

	ok, err := a.whiteList.Contains(ctx, ip)
	if err != nil {
		logger.Error().Err(err).Msg("error while checking ip in white list")
		return err
	}

	if ok {
		logger.Info().Msg("ip is in white list")
		return nil
	}

	ok, err = a.blackList.Contains(ctx, ip)
	if err != nil {
		logger.Error().Err(err).Msg("error while checking ip in black list")
		return err
	}

	if ok {
		logger.Info().Msg("ip is in black list")
		return apperr.ErrUserIsBlocked
	}

	if err := a.authIsAllowed(ctx, ip, login, pass); err != nil {
		logger.Error().Err(err).Msg("error while checking auth is allowed")
		return err
	}

	return nil
}
