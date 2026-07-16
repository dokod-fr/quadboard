package http

import (
	"io/fs"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/http/handlers"
	"github.com/dokod-fr/quadboard/internal/http/middleware"
	"github.com/dokod-fr/quadboard/internal/http/view"
	"github.com/go-chi/chi/v5"
)

func NewRouter(catalog *app.Catalog, oidc *auth.OIDC) http.Handler {
	r := chi.NewRouter()

	// --- Public routes ---
	r.Get("/health", handlers.ServeHealth)
	r.Get("/version", handlers.ServeVersion)

	// --- Assets ---
	// Error should never happen since the assets directory is embedded
	assetsFS, _ := fs.Sub(view.FS(), "assets")

	r.Handle(
		"/assets/*",
		middleware.StaticAssetsMiddleware(http.StripPrefix(
			"/assets/",
			http.FileServer(http.FS(assetsFS)),
		)),
	)

	// --- API routes ---
	catalogHandler := handlers.NewCatalogHandler(catalog)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/catalog", catalogHandler.Serve)
		// TODO: API serve
		// r.Post("/token", api.CreateToken)
		// r.Get("/resources", api.ListResources)
		// r.Post("/resources/{id}/start", api.StartResource)
	})

	// --- Protected routes ---
	homeHandler := handlers.NewHomeHandler()

	if oidc != nil {
		// This route is only used for OIDC login and callback, so we don't need to protect it with the AuthMiddleware
		r.Get("/login", oidc.LoginHandler)
		r.Get("/auth/callback", oidc.CallbackHandler)

		// Protected routes.
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(oidc))

			// UI routes
			r.Get("/", homeHandler.Serve)
		})
	} else {
		// Fallback to no authentication, all routes are public
		r.Get("/", homeHandler.Serve)
	}

	return r
}
