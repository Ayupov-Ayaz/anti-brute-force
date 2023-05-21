package manager

import (
	"context"
	"net"
	"strconv"
)

type IPList interface {
	Add(ctx context.Context, ip string) error
	Remove(ctx context.Context, ip string) error
}

type Resetter interface {
	Reset(ctx context.Context, login string, ip string) error
}

type IPService interface {
	ParseMaskedIP(ip, mask string) (net.IP, error)
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

func (a *App) makeMaskedKey(ip, mask string) (string, error) {
	maskedIP, err := a.ip.ParseMaskedIP(ip, mask)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(uint64(a.ip.IPToUint32(maskedIP)), 10), nil
}

func (a *App) AddToBlackList(ctx context.Context, ip, mask string) error {
	key, err := a.makeMaskedKey(ip, mask)
	if err != nil {
		return err
	}

	if err := a.blackList.Add(ctx, key); err != nil {
		return err
	}

	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ip, mask string) error {
	key, err := a.makeMaskedKey(ip, mask)
	if err != nil {
		return err
	}

	return a.whiteList.Add(ctx, key)
}

func (a *App) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	key, err := a.makeMaskedKey(ip, mask)
	if err != nil {
		return err
	}

	return a.blackList.Remove(ctx, key)
}

func (a *App) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	key, err := a.makeMaskedKey(ip, mask)
	if err != nil {
		return err
	}

	return a.whiteList.Remove(ctx, key)
}

func (a *App) Reset(ctx context.Context, login, dirtyIP string) error {
	ip, err := a.ip.ParseIP(dirtyIP)
	if err != nil {
		return err
	}

	key := strconv.FormatUint(uint64(a.ip.IPToUint32(ip)), 10)

	return a.resetter.Reset(ctx, login, key)
}
