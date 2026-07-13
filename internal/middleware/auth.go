package middleware

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/auth"
)

// AuthMiddleware check OIDC session. If inactive, redirect to /login.
func AuthMiddleware(o *auth.OIDC) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := o.GetSessionFromRequest(r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Inject groups in the context for the downstream handlers
			ctx := auth.WithSession(r.Context(), *session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
