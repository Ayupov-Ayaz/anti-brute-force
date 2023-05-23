package handlers

import (
	"context"

	"github.com/rs/zerolog"

	fiber "github.com/gofiber/fiber/v2"
)

//go:generate mockgen -source=manager.go -destination=mocks/mock_manager.go
type Manager interface {
	AddToBlackList(ctx context.Context, ip, mask string) error
	AddToWhiteList(ctx context.Context, ip, mask string) error
	RemoveFromBlackList(ctx context.Context, ip, mask string) error
	RemoveFromWhiteList(ctx context.Context, ip, mask string) error
	Reset(ctx context.Context, login, ip string) error
}

type ManagerHTTP struct {
	app       Manager
	validator Validator
	decoder   Decoder
	logger    zerolog.Logger
}

func NewManager(app Manager, validator Validator, decoder Decoder, logger zerolog.Logger) *ManagerHTTP {
	return &ManagerHTTP{
		app:       app,
		validator: validator,
		decoder:   decoder,
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

	app.Delete("/buckets", m.reset)
}

type AddToList func(ctx context.Context, ip, mask string) error

func (m *ManagerHTTP) parseIP(body []byte) (ip string, mask string, err error) {
	var model IP

	if err := m.decoder.Unmarshal(body, &model); err != nil {
		m.logger.Error().Err(err).Bytes("body", body).Msg("parse body failed")
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
	ip, mask, err := m.parseIP(ctx.Body())
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
	return m.addToList(ctx, m.app.AddToBlackList)
}

func (m *ManagerHTTP) addToWhiteList(ctx *fiber.Ctx) error {
	return m.addToList(ctx, m.app.AddToWhiteList)
}

type RemoveFromList func(ctx context.Context, ip, mask string) error

func (m *ManagerHTTP) removeFromList(ctx *fiber.Ctx, remove RemoveFromList) error {
	ip, mask, err := m.parseIP(ctx.Body())
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
	return m.removeFromList(ctx, m.app.RemoveFromBlackList)
}

func (m *ManagerHTTP) removeFromWhiteList(ctx *fiber.Ctx) error {
	return m.removeFromList(ctx, m.app.RemoveFromWhiteList)
}

func (m *ManagerHTTP) reset(ctx *fiber.Ctx) error {
	var model BaseRequest

	if err := m.decoder.Unmarshal(ctx.Body(), &model); err != nil {
		m.logger.Error().Err(err).Bytes("body", ctx.Body()).Msg("parse body failed")
		return err
	}

	if err := m.validator.Validate(model); err != nil {
		m.logger.Error().Err(err).Str("login", model.Login).Str("ip", model.IP).Msg("validate auth failed")
		return err
	}

	if err := m.app.Reset(ctx.Context(), model.Login, model.IP); err != nil {
		m.logger.Error().Err(err).Str("login", model.Login).Msg("reset failed")
		return err
	}

	if err := ctx.SendStatus(fiber.StatusOK); err != nil {
		m.logger.Error().Err(err).Msg("send status failed")
		return err
	}

	return nil
}
