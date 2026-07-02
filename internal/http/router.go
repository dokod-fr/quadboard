package http

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handlers.Health)
	r.Get("/version", handlers.Version)

	return r
}
