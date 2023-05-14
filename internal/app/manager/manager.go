package manager

import "context"

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

func (a *App) AddToBlackList(ctx context.Context, ip, mask string) error {
	return a.blackList.Add(ctx, ip)
}

func joinIPMask(ip, mask string) string {
	return ip + "/" + mask
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
