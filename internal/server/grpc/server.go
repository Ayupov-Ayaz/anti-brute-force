package grpcserver

import "context"

type Manager interface {
	AddToBlackList(ctx context.Context, ip, mask string) error
	AddToWhiteList(ctx context.Context, ip, mask string) error
	RemoveFromBlackList(ctx context.Context, ip, mask string) error
	RemoveFromWhiteList(ctx context.Context, ip, mask string) error
	Reset(ctx context.Context, login, p string) error
}

type Checker interface {
	Check(ctx context.Context, ip, login, pass string) error
}

type Server struct {
	manager Manager
	checker Checker
}

type Config func(s *Server)

func New(configs ...Config) *Server {
	s := &Server{}

	for _, config := range configs {
		config(s)
	}

	return s
}

func WithManager(manager Manager) Config {
	return func(s *Server) {
		s.manager = manager
	}
}

func WithChecker(checker Checker) Config {
	return func(s *Server) {
		s.checker = checker
	}
}

func (s *Server) Start(port int) error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}
