package checker

import (
	"context"
	"fmt"

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
}

func New(whiteList, blackList IPList, checker Checker) *App {
	return &App{
		whiteList: whiteList,
		blackList: blackList,
		checker:   checker,
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
	ok, err := a.whiteList.Contains(ctx, ip)
	if err != nil {
		return fmt.Errorf("check white list: %w", err)
	}

	if ok {
		return nil
	}

	ok, err = a.blackList.Contains(ctx, ip)
	if err != nil {
		return fmt.Errorf("check black list: %w", err)
	}

	if ok {
		return apperr.ErrUserIsBlocked
	}

	if err := a.authIsAllowed(ctx, ip, login, pass); err != nil {
		return fmt.Errorf("check auth: %w", err)
	}

	return nil
}
