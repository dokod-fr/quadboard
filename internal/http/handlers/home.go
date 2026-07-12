package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/http/views/pages"
)

type HomeHandler struct {
	discovery *app.Discovery
}

func NewHomeHandler(discovery *app.Discovery) *HomeHandler {
	return &HomeHandler{
		discovery: discovery,
	}
}

func (h *HomeHandler) Serve(w http.ResponseWriter, r *http.Request) {
	resources, err := h.discovery.Resources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	session, _ := auth.SessionFromContext(ctx)

	slog.Debug("Home handler", slog.Any("session", session))

	pages.Home(resources, session).Render(ctx, w)
}
