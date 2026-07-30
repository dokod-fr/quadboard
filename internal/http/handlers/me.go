package handlers

import (
	"encoding/json"
	"log/slog"
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

	if h.cfg.Auth.OIDC != nil && h.cfg.Auth.OIDC.Issuer != "" {
		if session, ok := auth.SessionFromContext(r.Context()); ok {
			username = session.Username
		}
	}

	// On essaie juste de récupérer la session depuis le contexte
	if session, ok := auth.SessionFromContext(r.Context()); ok {
		slog.Warn("username", slog.String("username", session.Username))
	} else {
		// Si on voit ce log dans la console, c'est que le middleware n'a pas tourné
		slog.Warn("Session not found in context for /api/v1/me")
	}

	response := struct {
		Username       string `json:"username"`
		GroupByDefault bool   `json:"groupByDefault"`
	}{
		Username:       username,
		GroupByDefault: h.cfg.UI.GroupByDefault,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
