package handlers

import (
	"context"

	"github.com/rs/zerolog"

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
	manager   Manager
	validator Validator
	logger    zerolog.Logger
}

func NewManager(app Manager, validator Validator, logger zerolog.Logger) *ManagerHTTP {
	return &ManagerHTTP{
		manager:   app,
		validator: validator,
		logger:    logger,
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

func (m *ManagerHTTP) parseIP(ctx *fiber.Ctx) (ip string, mask string, err error) {
	var model IP

	if err := ctx.BodyParser(&model); err != nil {
		m.logger.Error().Err(err).Bytes("body", ctx.Body()).Msg("parse body failed")
		return "", "", err
	}

	if err := m.validator.Validate(model); err != nil {
		m.logger.Error().Err(err).Str("mask", model.Mask).
			Str("ip", model.IP).Msg("validate ip failed")
		return "", "", err
	}

	return model.IP, model.Mask, nil
}

func (m *ManagerHTTP) addToList(ctx *fiber.Ctx, addToList AddToList) error {
	ip, mask, err := m.parseIP(ctx)
	if err != nil {
		return err
	}

	if err := addToList(ctx.Context(), ip, mask); err != nil {
		m.logger.Error().Err(err).Str("mask", mask).
			Str("ip", ip).Msg("add to list failed")
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
	ip, mask, err := m.parseIP(ctx)
	if err != nil {
		return err
	}

	if err := remove(ctx.Context(), ip, mask); err != nil {
		m.logger.Error().Err(err).Str("mask", mask).
			Str("ip", ip).Msg("remove from list failed")
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
	var model Auth

	if err := ctx.BodyParser(&model); err != nil {
		m.logger.Error().Err(err).Bytes("body", ctx.Body()).Msg("parse body failed")
		return err
	}

	if err := m.validator.Validate(model); err != nil {
		m.logger.Error().Err(err).Str("login", model.Login).
			Str("pass", model.Pass).Msg("validate auth failed")
		return err
	}

	if err := m.manager.Reset(ctx.Context(), model.Login, model.Pass); err != nil {
		m.logger.Error().Err(err).Str("login", model.Login).Msg("reset failed")
		return err
	}

	if err := ctx.SendStatus(fiber.StatusOK); err != nil {
		m.logger.Error().Err(err).Msg("send status failed")
		return err
	}

	return nil
}
