package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/config"
)

type MeHandler struct {
	cfg *config.Config
}

func NewMeHandler(cfg *config.Config) *MeHandler {
	return &MeHandler{cfg: cfg}
}

func (h *MeHandler) Serve(w http.ResponseWriter, r *http.Request) {
	username := ""
	groups := make([]string, 0)

	if h.cfg.Auth.OIDC != nil && h.cfg.Auth.OIDC.Issuer != "" {
		if session, ok := auth.SessionFromContext(r.Context()); ok {
			username = session.Username
			groups = session.Groups
		}
	}

	response := struct {
		Username       string   `json:"username"`
		Groups         []string `json:"groups"`
		GroupByDefault bool     `json:"groupByDefault"`
	}{
		Username:       username,
		Groups:         groups,
		GroupByDefault: h.cfg.UI.GroupByDefault,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
