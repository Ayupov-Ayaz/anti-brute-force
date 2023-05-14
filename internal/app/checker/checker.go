package checker

import (
	"context"
	"errors"
	"fmt"
)

//todo: log

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

func (a *App) Check(ctx context.Context, ip, login, pass string) error {
	ok, err := a.whiteList.Contains(ctx, ip)
	if err != nil {
		return fmt.Errorf("white list: %w", err)
	}

	ok, err = a.blackList.Contains(ctx, ip)
	if err != nil {
		return err
	}

	if ok {
		return ErrUserIsBlocked
	}

	if err := a.buckets.Check(ctx, ip, login, pass); err != nil {
		return err
	}

	return nil
}
