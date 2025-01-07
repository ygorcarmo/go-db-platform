package handlers

import (
	"net/http"

	"db-platform/views/settings"
)

func GetSettingsPage(w http.ResponseWriter, r *http.Request) {
	settings.Index().Render(r.Context(), w)
}
