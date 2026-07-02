package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/dokod-fr/quadboard/internal/config"
)

type Server struct {
	cfg    config.Config
	server *http.Server
}

func New(cfg config.Config) *Server {
	r := NewRouter(cfg)

	s := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		cfg:    cfg,
		server: s,
	}
}

func (s *Server) Start(ctx context.Context) error {
	slog.Info("http server starting",
		"addr", s.cfg.Server.Address,
		"read_timeout", s.cfg.Server.ReadTimeout,
		"write_timeout", s.cfg.Server.WriteTimeout,
	)
	go func() {
		<-ctx.Done()
		_ = s.server.Shutdown(context.Background())
	}()

	return s.server.ListenAndServe()
}
