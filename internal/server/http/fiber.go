package httpserver

import (
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		return ctx.Status(http.StatusInternalServerError).JSON(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
	}
}

func useCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
	})
}

func useRecover() fiber.Handler {
	return recover.New(recover.Config{EnableStackTrace: true})
}

func NewFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:   ErrorHandler(),
		AppName:        "anti-brute-force",
		RequestMethods: []string{http.MethodPost},
	})

	app.Use(useCors(), useRecover())

	return app
}
