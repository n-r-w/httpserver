package httpserver

import (
	"net"
	"time"
)

type Option func(*Service)

func Address(host string, port string) Option {
	return func(s *Service) {
		s.server.Addr = net.JoinHostPort(host, port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Service) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Service) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Service) {
		s.shutdownTimeout = timeout
	}
}
