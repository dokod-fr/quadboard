package http

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/http/handlers"
	"github.com/dokod-fr/quadboard/internal/middleware"
	"github.com/dokod-fr/quadboard/web"
	"github.com/go-chi/chi/v5"
)

func NewRouter(discovery *app.Discovery, oidc *auth.OIDC) http.Handler {
	r := chi.NewRouter()

	// --- Public routes ---
	r.Get("/health", handlers.ServeHealth)
	r.Get("/version", handlers.ServeVersion)

	// Manage assets
	r.Handle(
		"/assets/*",
		http.StripPrefix(
			"/assets/",
			http.FileServer(http.FS(web.FS())),
		),
	)

	// --- Authentification configuration ---
	if oidc != nil {
		// This route is only used for OIDC login and callback, so we don't need to protect it with the AuthMiddleware
		r.Get("/login", oidc.LoginHandler)
		r.Get("/auth/callback", oidc.CallbackHandler)

		// Protected routes.
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(oidc))

			// UI routes
			r.Get("/", handlers.NewHomeHandler(discovery).Serve)
		})
	} else {
		// Fallback to no authentication, all routes are public
		r.Get("/", handlers.NewHomeHandler(discovery).Serve)
	}

	return r
}
