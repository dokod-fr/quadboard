package middleware

import "net/http"

func StaticAssetsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cache d'un an (idéal si tes assets ont un hash ou changent peu)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		next.ServeHTTP(w, r)
	})
}
