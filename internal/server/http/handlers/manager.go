package handlers

import (
	"context"

	"go.uber.org/zap"

	fiber "github.com/gofiber/fiber/v2"
)

type Manager interface {
	AddToBlackList(ctx context.Context, ip, mask string) error
	AddToWhiteList(ctx context.Context, ip, mask string) error
	RemoveFromBlackList(ctx context.Context, ip, mask string) error
	RemoveFromWhiteList(ctx context.Context, ip, mask string) error
	Reset(ctx context.Context, login, pass string) error
}

type ManagerHTTP struct {
	manager Manager
	logger  *zap.Logger
}

func NewManager(app Manager, logger *zap.Logger) *ManagerHTTP {
	return &ManagerHTTP{
		manager: app,
		logger:  logger,
	}
}

func (m *ManagerHTTP) Register(app *fiber.App) {
	bl := app.Group("/black-list")
	bl.Post("/add", m.addToBlackList)
	bl.Delete("/remove", m.removeFromBlackList)

	wl := app.Group("/white-list")
	wl.Post("/add", m.addToWhiteList)
	wl.Delete("/remove", m.removeFromWhiteList)

	app.Post("/buckets/reset", m.reset)
}

type AddToList func(ctx context.Context, ip, mask string) error

func (m *ManagerHTTP) addToList(ctx *fiber.Ctx, addToList AddToList) error {
	var model IP

	if err := ctx.BodyParser(&model); err != nil {
		m.logger.Error("parse body failed",
			zap.ByteString("body", ctx.Body()),
			zap.Error(err))
		return err
	}

	if err := addToList(ctx.Context(), model.IP, model.Mask); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (m *ManagerHTTP) addToBlackList(ctx *fiber.Ctx) error {
	return m.addToList(ctx, m.manager.AddToBlackList)
}

func (m *ManagerHTTP) addToWhiteList(ctx *fiber.Ctx) error {
	return m.addToList(ctx, m.manager.AddToWhiteList)
}

type RemoveFromList func(ctx context.Context, ip, mask string) error

func (m *ManagerHTTP) removeFromList(ctx *fiber.Ctx, remove RemoveFromList) error {
	var model IP

	if err := ctx.BodyParser(&model); err != nil {
		m.logger.Error("parse body failed",
			zap.ByteString("body", ctx.Body()),
			zap.Error(err))
		return err
	}

	if err := remove(ctx.Context(), model.IP, model.Mask); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (m *ManagerHTTP) removeFromBlackList(ctx *fiber.Ctx) error {
	return m.removeFromList(ctx, m.manager.RemoveFromBlackList)
}

func (m *ManagerHTTP) removeFromWhiteList(ctx *fiber.Ctx) error {
	return m.removeFromList(ctx, m.manager.RemoveFromWhiteList)
}

func (m *ManagerHTTP) reset(ctx *fiber.Ctx) error {
	var auth Auth
	if err := ctx.BodyParser(&auth); err != nil {
		m.logger.Error("parse body failed",
			zap.ByteString("body", ctx.Body()),
			zap.Error(err))
		return err
	}

	if err := m.manager.Reset(ctx.Context(), auth.Login, auth.Pass); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
