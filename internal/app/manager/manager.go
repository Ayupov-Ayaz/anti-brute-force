package manager

import (
	"context"
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
}

func New(whiteList, blackList IPList, resetter Resetter) *App {
	return &App{
		blackList: blackList,
		whiteList: whiteList,
		resetter:  resetter,
	}
}

func joinIPMask(ip, mask string) string {
	// todo: разобраться как работать с маской
	return ip
}

func (a *App) AddToBlackList(ctx context.Context, ip, mask string) error {
	if err := a.blackList.Add(ctx, joinIPMask(ip, mask)); err != nil {
		return err
	}

	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ip, mask string) error {
	return a.whiteList.Add(ctx, joinIPMask(ip, mask))
}

func (a *App) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	return a.blackList.Remove(ctx, joinIPMask(ip, mask))
}

func (a *App) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	return a.whiteList.Remove(ctx, joinIPMask(ip, mask))
}

func (a *App) Reset(ctx context.Context, login, ip string) error {
	return a.resetter.Reset(ctx, login, ip)
}
