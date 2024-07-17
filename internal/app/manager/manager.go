package manager

import (
	"context"
	"net"
	"strconv"
)

//go:generate mockgen -source=manager.go -destination=mocks/mock.go
type IPList interface {
	Add(ctx context.Context, ip string) error
	Remove(ctx context.Context, ip string) error
}

type Resetter interface {
	Reset(ctx context.Context, login string, ip string) error
}

type IPService interface {
	ParseCIDR(ip string) (*net.IPNet, error)
	ParseIP(ip string) (net.IP, error)
	IPToUint32(ip net.IP) uint32
}

type App struct {
	blackList IPList
	whiteList IPList
	ip        IPService
	resetter  Resetter
}

func New(whiteList, blackList IPList, ip IPService, resetter Resetter) *App {
	return &App{
		blackList: blackList,
		whiteList: whiteList,
		resetter:  resetter,
		ip:        ip,
	}
}

func (a *App) AddToBlackList(ctx context.Context, ipNet string) error {
	return a.blackList.Add(ctx, ipNet)
}

func (a *App) AddToWhiteList(ctx context.Context, ipNet string) error {
	return a.whiteList.Add(ctx, ipNet)
}

func (a *App) RemoveFromBlackList(ctx context.Context, ipNet string) error {
	return a.blackList.Remove(ctx, ipNet)
}

func (a *App) RemoveFromWhiteList(ctx context.Context, ipNet string) error {
	return a.whiteList.Remove(ctx, ipNet)
}

func (a *App) Reset(ctx context.Context, login, dirtyIP string) error {
	ip, err := a.ip.ParseIP(dirtyIP)
	if err != nil {
		return err
	}

	key := strconv.FormatUint(uint64(a.ip.IPToUint32(ip)), 10)

	return a.resetter.Reset(ctx, login, key)
}
