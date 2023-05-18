package handlers

import (
	"context"
	"errors"

	"github.com/ayupov-ayaz/anti-brute-force/internal/apperr"

	fiber "github.com/gofiber/fiber/v2"
)

type Checker interface {
	Check(ctx context.Context, ip, login, pass string) error
}

type Validator interface {
	Validate(i interface{}) error
}

type CheckerHTTP struct {
	app       Checker
	validator Validator
}

func NewChecker(app Checker, validator Validator) *CheckerHTTP {
	return &CheckerHTTP{
		app:       app,
		validator: validator,
	}
}

func (c *CheckerHTTP) Register(app *fiber.App) {
	group := app.Group("/checker")

	group.Post("/check", c.check)
	group.Post("/reset", c.reset)
}

func (c *CheckerHTTP) check(ctx *fiber.Ctx) error {
	var auth Auth
	if err := ctx.BodyParser(&auth); err != nil {
		return err
	}

	if err := c.validator.Validate(auth); err != nil {
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
			return err
		}
	}

	return ctx.Status(status).JSON(Response{Ok: allowed})
}

func (c *CheckerHTTP) reset(ctx *fiber.Ctx) error {
	return errors.New("reset not implemented")
}
