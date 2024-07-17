package handlers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	fiber "github.com/gofiber/fiber/v2"
)

//go:generate mockgen -source=checker.go -destination=mocks/mock_checker.go
type Checker interface {
	Check(ctx context.Context, ip, login, pass string) error
}

type Validator interface {
	Validate(i interface{}) error
}

type Decoder interface {
	Marshal(i interface{}) ([]byte, error)
	Unmarshal(data []byte, i interface{}) error
}

type CheckerHTTP struct {
	app       Checker
	validator Validator
	decoder   Decoder
	logger    zerolog.Logger
}

func NewChecker(app Checker, validator Validator, decoder Decoder, logger zerolog.Logger) *CheckerHTTP {
	return &CheckerHTTP{
		app:       app,
		validator: validator,
		decoder:   decoder,
		logger:    logger,
	}
}

func (c *CheckerHTTP) Register(app *fiber.App) {
	group := app.Group("/checker")

	group.Post("/check", c.check)
}

func (c *CheckerHTTP) check(ctx *fiber.Ctx) error {
	var auth CheckAuthRequest
	body := ctx.Body()
	if err := c.decoder.Unmarshal(body, &auth); err != nil {
		c.logger.Error().Bytes("req", body).Err(err).
			Msg("failed to parse request body")
		return err
	}

	if err := c.validator.Validate(auth); err != nil {
		c.logger.Error().Err(err).Msg("failed to validate request body")
		return err
	}

	err := c.app.Check(ctx.Context(), auth.IP, auth.Login, auth.Pass)
	status := fiber.StatusOK
	allowed := true
	if err != nil {
		if errors.Is(err, apperr.ErrUserIsBlocked) {
			status = fiber.StatusForbidden
			allowed = false
		} else {
			c.logger.Error().Err(err).Msg("failed to check user")
			return err
		}
	}

	respBody, err := c.decoder.Marshal(Response{Ok: allowed})
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal response")
		return err
	}

	if err := ctx.Status(status).Send(respBody); err != nil {
		c.logger.Error().Err(err).Msg("failed to send response")
		return err
	}

	return nil
}
