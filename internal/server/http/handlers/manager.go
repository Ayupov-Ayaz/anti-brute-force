package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Manager interface {
	AddToBlackList(ctx context.Context, ip, mask string) error
	AddToWhiteList(ctx context.Context, ip, mask string) error
	RemoveFromBlackList(ctx context.Context, ip, mask string) error
	RemoveFromWhiteList(ctx context.Context, ip, mask string) error
	Reset(ctx context.Context, login, p string) error
}

type ManagerHTTP struct {
	manager Manager
}

func NewManager(app Manager) *ManagerHTTP {
	return &ManagerHTTP{
		manager: app,
	}
}

func (m *ManagerHTTP) Register(app *fiber.App) {
	app.Post("/add-to-black-list", m.addToBlackList)
	app.Post("/add-to-white-list", m.addToWhiteList)
	app.Post("/remove-from-black-list", m.removeFromBlackList)
	app.Post("/remove-from-white-list", m.removeFromWhiteList)
	app.Post("/reset", m.reset)
}

func (m *ManagerHTTP) addToBlackList(ctx *fiber.Ctx) error {
	return nil
}

func (m *ManagerHTTP) addToWhiteList(ctx *fiber.Ctx) error {
	return nil
}

func (m *ManagerHTTP) removeFromBlackList(ctx *fiber.Ctx) error {
	return nil
}

func (m *ManagerHTTP) removeFromWhiteList(ctx *fiber.Ctx) error {
	return nil
}

func (m *ManagerHTTP) reset(ctx *fiber.Ctx) error {
	return nil
}
