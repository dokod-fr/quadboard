package middleware

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/auth"
)

// AuthMiddleware vérifie la session OIDC. Si inactive, redirige vers /login ou renvoie un 401.
func AuthMiddleware(o *auth.OIDC) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Anti-loops detection
			if r.URL.Path == "/login" || r.URL.Path == "/auth/callback" {
				next.ServeHTTP(w, r)
				return
			}

			session, err := o.GetSessionFromRequest(r)
			if err != nil {
				// Detect which kind of redirection is
				// HTMX detection
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Classical API JSON requests redirection
				if r.Header.Get("Accept") == "application/json" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Standard redirection
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Inject session into the context
			ctx := auth.WithSession(r.Context(), *session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
