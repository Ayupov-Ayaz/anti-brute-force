package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Checker interface {
	Check(ctx context.Context, ip, login, pass string) error
}

type CheckerHTTP struct {
	app Checker
}

func NewChecker(app Checker) *CheckerHTTP {
	return &CheckerHTTP{
		app: app,
	}
}

func (c *CheckerHTTP) Register(app *fiber.App) {
	app.Post("/check", c.check)
}

func (c *CheckerHTTP) check(ctx *fiber.Ctx) error {
	return nil
}
