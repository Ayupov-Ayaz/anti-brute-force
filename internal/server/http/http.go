package httpserver

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/ayupov-ayaz/anti-brute-force/internal/server/http/handlers"

	fiber "github.com/gofiber/fiber/v2"
)

type shutdown func() error

type Server struct {
	manager  *handlers.ManagerHTTP
	checker  *handlers.CheckerHTTP
	shutdown shutdown
	logger   zerolog.Logger
}

func New(manager *handlers.ManagerHTTP, checker *handlers.CheckerHTTP, logger zerolog.Logger) *Server {
	return &Server{
		manager: manager,
		checker: checker,
		logger:  logger,
	}
}

func (s *Server) Register(app *fiber.App) {
	s.manager.Register(app)
	s.checker.Register(app)
}

func (s *Server) Start(app *fiber.App, port int) error {
	s.shutdown = app.Shutdown
	s.Register(app)

	if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
		return fmt.Errorf("listen port=%d failed: %w", port, err)
	}

	return nil
}

func (s *Server) Stop() error {
	if s.shutdown != nil {
		return s.shutdown()
	}

	return nil
}
