package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/dokod-fr/quadboard/internal/config"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg config.Config, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:         cfg.Server.Address,
			Handler:      handler,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		},
	}
}

func (s *Server) Run() error {
	slog.Info("starting HTTP server", "addr", s.server.Addr)

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("stopping HTTP server")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
