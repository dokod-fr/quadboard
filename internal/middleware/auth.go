package middleware

import (
	"context"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/auth"
)

type contextKey string

const GroupsKey contextKey = "groups"

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
			ctx := context.WithValue(r.Context(), GroupsKey, session.Groups)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
