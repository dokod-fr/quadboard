package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dokod-fr/quadboard/internal/version"
)

func ServeVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"version": version.Version,
		"commit":  version.Commit,
		"date":    version.Date,
	})
}
