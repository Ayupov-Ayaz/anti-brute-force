package httpserver

import (
	"fmt"
	"strconv"

	"github.com/ayupov-ayaz/anti-brute-force/internal/server/http/handlers"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	manager  *handlers.ManagerHTTP
	checker  *handlers.CheckerHTTP
	shutdown func() error
}

type Config func(s *Server)

func New(configs ...Config) *Server {
	s := &Server{}
	for _, config := range configs {
		config(s)
	}

	return s
}

func WithManager(manager *handlers.ManagerHTTP) Config {
	return func(s *Server) {
		s.manager = manager
	}
}

func WithChecker(checker *handlers.CheckerHTTP) Config {
	return func(s *Server) {
		s.checker = checker
	}
}

func (s *Server) Start(app *fiber.App, port int) error {
	s.shutdown = app.Shutdown

	s.manager.Register(app)
	s.checker.Register(app)

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
