package handlers

import (
	"html/template"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/domain"
	"github.com/dokod-fr/quadboard/internal/http/view"
)

type HomeHandler struct {
	discovery *app.Discovery
	tmpl      *template.Template
}

func NewHomeHandler(discovery *app.Discovery) *HomeHandler {
	tmpl := template.Must(
		template.ParseFS(view.FS(),
			"templates/layout.html",
			"templates/home.html",
		),
	)

	return &HomeHandler{
		discovery: discovery,
		tmpl:      tmpl,
	}
}

func (h *HomeHandler) Serve(w http.ResponseWriter, r *http.Request) {
	resources, err := h.discovery.Resources()
	if err != nil {
		http.Error(w, "Failed to load resources", http.StatusInternalServerError)
		return
	}

	var username string
	if session, ok := auth.SessionFromContext(r.Context()); ok {
		username = session.Username
	}

	data := struct {
		Resources []domain.Resource
		Username  string
	}{
		Resources: resources,
		Username:  username,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
