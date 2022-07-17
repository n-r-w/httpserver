// Package httpserver ...
package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/n-r-w/lg"
	"github.com/n-r-w/nerr"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":8080"
	defaultShutdownTimeout = 3 * time.Second
)

type Service struct {
	server          *http.Server
	logger          lg.Logger
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, logger lg.Logger, opts ...Option) *Service {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Service{
		server:          httpServer,
		logger:          logger,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	go s.start()

	return s
}

func (s *Service) start() {
	l, err := net.Listen("tcp", s.server.Addr)
	if err == nil {
		s.logger.Info("http server started on %s", s.server.Addr)
		err = s.server.Serve(l)
	}
	if err != nil {
		s.notify <- nerr.New("net.Listen error", err)
	}
	close(s.notify)
}

func (s *Service) Notify() <-chan error {
	return s.notify
}

func (s *Service) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
