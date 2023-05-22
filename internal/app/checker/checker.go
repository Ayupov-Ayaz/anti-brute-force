package checker

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"
)

//go:generate mockgen -source=checker.go -destination=mocks/mock.go
type IPList interface {
	Contains(ctx context.Context, ip string) (bool, error)
}

type Checker interface {
	AllowByLogin(ctx context.Context, login string) error
	AllowByPassword(ctx context.Context, login string) error
	AllowByIP(ctx context.Context, login string) error
}

type IPService interface {
	ParseIP(ip string) (net.IP, error)
	IPToUint32(ip net.IP) uint32
}

type App struct {
	whiteList IPList
	blackList IPList
	checker   Checker
	ip        IPService
}

func New(whiteList, blackList IPList, ip IPService, checker Checker) *App {
	return &App{
		whiteList: whiteList,
		blackList: blackList,
		checker:   checker,
		ip:        ip,
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

func (a *App) parseIPKey(ip string) (string, error) {
	parsedIP, err := a.ip.ParseIP(ip)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(uint64(a.ip.IPToUint32(parsedIP)), 10), nil
}

func (a *App) Check(ctx context.Context, dirtyIP, login, pass string) error {
	ipKey, err := a.parseIPKey(dirtyIP)
	if err != nil {
		return err
	}

	ok, err := a.whiteList.Contains(ctx, ipKey)
	if err != nil {
		return fmt.Errorf("check white list: %w", err)
	}

	if ok {
		return nil
	}

	ok, err = a.blackList.Contains(ctx, ipKey)
	if err != nil {
		return fmt.Errorf("check black list: %w", err)
	}

	if ok {
		return apperr.ErrUserIsBlocked
	}

	if err := a.authIsAllowed(ctx, ipKey, login, pass); err != nil {
		return fmt.Errorf("check auth: %w", err)
	}

	return nil
}
