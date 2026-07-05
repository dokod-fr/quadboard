package http

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/http/handlers"
	"github.com/dokod-fr/quadboard/web"
	"github.com/go-chi/chi/v5"
)

func NewRouter(discovery *app.Discovery) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handlers.ServeHealth)
	r.Get("/version", handlers.ServeVersion)

	// UI
	r.Get("/", handlers.NewHomeHandler(discovery).Serve)

	// Manage assets
	r.Handle(
		"/assets/*",
		http.StripPrefix(
			"/assets/",
			http.FileServer(http.FS(web.FS())),
		),
	)
	return r
}
