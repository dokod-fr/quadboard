package handlers

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/http/views"
)

func Home(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		views.Home(cfg).Render(r.Context(), w)
	}
}
