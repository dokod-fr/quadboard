package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/domain"
)

type CatalogHandler struct {
	catalog *app.Catalog
}

func NewCatalogHandler(c *app.Catalog) *CatalogHandler {
	return &CatalogHandler{
		catalog: c,
	}
}

func (h *CatalogHandler) Serve(w http.ResponseWriter, r *http.Request) {
	resources := h.catalog.Resources()

	if resources == nil {
		resources = make([]domain.Resource, 0)
	}

	w.Header().Set("Content-Type", "application/json")

	// Encoding directly into ResponseWriter improve performance
	if err := json.NewEncoder(w).Encode(resources); err != nil {
		http.Error(w, "Failed to encode catalog", http.StatusInternalServerError)
		return
	}
}
