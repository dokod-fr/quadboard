package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/domain"
)

type CatalogHandler struct {
	catalog *app.Catalog
	config  *config.Config
}

// Front-end resources
type ResourceView struct {
	domain.Resource
	Authorized  bool   `json:"Authorized"`
	DisplayName string `json:"DisplayName"`
}

func NewCatalogHandler(c *app.Catalog, cfg *config.Config) *CatalogHandler {
	return &CatalogHandler{
		catalog: c,
		config:  cfg,
	}
}

func (h *CatalogHandler) Serve(w http.ResponseWriter, r *http.Request) {
	rawResources := h.catalog.Resources()
	userGroups := []string{}

	// Get user groups only if OIDC is activated
	if h.config.Auth.OIDC != nil {
		if session, ok := auth.SessionFromContext(r.Context()); ok {
			userGroups = session.Groups
		}
	}

	viewResources := make([]ResourceView, 0, len(rawResources))

	for _, res := range rawResources {
		// 1. Mapping du libellé
		displayName := res.Group
		if displayName == "" {
			displayName = "Default"
		}
		if label, ok := h.config.UI.GroupLabels[res.Group]; ok {
			displayName = label
		}

		// Check right: No group or no authen means public
		isAuthorized := true
		if h.config.Auth.OIDC != nil && res.Group != "" {
			isAuthorized = contains(userGroups, res.Group)
		}

		if h.config.UI.HideUnauthorizedGroups && !isAuthorized {
			continue
		}

		// Show it self or not
		isQuadboard := strings.EqualFold(res.Name, "quadboard") || strings.Contains(strings.ToLower(res.Name), "quadboard")
		if isQuadboard && !h.config.UI.ShowItSelf {
			continue
		}

		viewResources = append(viewResources, ResourceView{
			Resource:    res,
			Authorized:  isAuthorized,
			DisplayName: displayName,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	// Encoding directly into ResponseWriter improve performance
	if err := json.NewEncoder(w).Encode(viewResources); err != nil {
		http.Error(w, "Failed to encode catalog", http.StatusInternalServerError)
		return
	}
}

/* ---- Helper ---- */
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
