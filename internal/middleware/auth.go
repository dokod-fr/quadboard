package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dokod-fr/quadboard/internal/auth"
)

// AuthMiddleware gère à la fois l'authentification par en-têtes (Reverse Proxy / Authelia)
// et l'authentification OIDC applicative en fallback.
func AuthMiddleware(o *auth.OIDC) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" || r.URL.Path == "/auth/callback" {
				next.ServeHTTP(w, r)
				return
			}

			// --- Use Authelia headers is already connected ---
			username := r.Header.Get("Remote-User")
			if username != "" {
				groupsHeader := r.Header.Get("Remote-Groups")
				var groups []string
				if groupsHeader != "" {
					// Authelia send groups with comma separator
					for _, g := range strings.Split(groupsHeader, ",") {
						groups = append(groups, strings.TrimSpace(g))
					}
				}

				slog.Info("Authentification successed via reverse-proxy (headers)",
					"user", username,
					"groups", groups,
				)

				session := auth.Session{
					Username: username,
					Groups:   groups,
					Source:   "headers",
				}

				ctx := auth.WithSession(r.Context(), session)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// --- OIDC configuration ---
			if o != nil {
				session, err := o.GetSessionFromRequest(r)
				if err == nil {
					slog.Info("Authentification successed via local session OIDC",
						"user", session.Username,
					)
					session.Source = "oidc"
					ctx := auth.WithSession(r.Context(), *session)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				slog.Info("OIDC session expired or missing, redirection to /login", "error", err)

				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if r.Header.Get("Accept") == "application/json" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// --- No headers and no OIDC conf ---
			slog.Warn("Access denied",
				"path", r.URL.Path,
				"ip", r.RemoteAddr,
			)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: No authentication method available"))
		})
	}
}
