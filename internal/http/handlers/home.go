package handlers

import (
	"net/http"

	"github.com/dokod-fr/quadboard/internal/domain"
	"github.com/dokod-fr/quadboard/internal/http/views/pages"
	"github.com/dokod-fr/quadboard/internal/mock"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resources := []domain.Resource{
			mock.Grafana(),
		}

		if err := pages.Home(resources).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
