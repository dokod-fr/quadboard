package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/http/view"
)

type HomeHandler struct {
	tmpl *template.Template
}

func NewHomeHandler() *HomeHandler {
	tmpl := template.Must(
		template.ParseFS(view.FS(),
			"templates/layout.html",
			"templates/home.html",
		),
	)

	return &HomeHandler{
		tmpl: tmpl,
	}
}

func (h *HomeHandler) Serve(w http.ResponseWriter, r *http.Request) {
	var username string
	if session, ok := auth.SessionFromContext(r.Context()); ok {
		username = session.Username
	}

	data := struct {
		Username string
	}{
		Username: username,
	}

	// 1. Rendu dans le buffer
	buf := new(bytes.Buffer)
	if err := h.tmpl.ExecuteTemplate(buf, "layout.html", data); err != nil {
		slog.Error("Failed to render template layout.html", "error", err)

		// 2. Au lieu de http.Error, on appelle notre helper HTML d'erreur
		serveErrorPage(w, http.StatusInternalServerError, "Impossible de charger la page d'accueil.")
		return
	}

	// Tout est OK, on envoie le résultat
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = buf.WriteTo(w)
}

// serveErrorPage renvoie une page HTML d'erreur minimaliste mais propre,
// stylisée directement, sans dépendre du système de templates externe.
func serveErrorPage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	// Un petit HTML "inline" propre avec un style CSS basique et moderne.
	// Pas besoin de fichier externe, comme ça si l'embed FS ou le parseur a un problème,
	// cette page fonctionnera TOUJOURS.
	html := `
	<!DOCTYPE html>
	<html lang="fr">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Erreur - QuadBoard</title>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
				background-color: #0f172a;
				color: #f8fafc;
				display: flex;
				align-items: center;
				justify-content: center;
				height: 100vh;
				margin: 0;
			}
			.card {
				background-color: #1e293b;
				padding: 2rem;
				border-radius: 12px;
				box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
				text-align: center;
				max-width: 400px;
				border: 1px solid #334155;
			}
			h1 { color: #f43f5e; margin-top: 0; }
			p { color: #94a3b8; line-height: 1.5; }
			.btn {
				display: inline-block;
				margin-top: 1.5rem;
				padding: 0.5rem 1rem;
				background-color: #3b82f6;
				color: white;
				text-decoration: none;
				border-radius: 6px;
				font-weight: 500;
				transition: background-color 0.2s;
			}
			.btn:hover { background-color: #2563eb; }
		</style>
	</head>
	<body>
		<div class="card">
			<h1>Oups !</h1>
			<p>` + message + `</p>
			<a href="/" class="btn">Retour à l'accueil</a>
		</div>
	</body>
	</html>`

	_, _ = w.Write([]byte(html))
}
