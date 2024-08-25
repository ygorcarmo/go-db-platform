package handlers

import (
	"net/http"

	"github.com/ygorcarmo/db-platform/views/setting"
)

func GetSettingsPage(w http.ResponseWriter, r *http.Request) {
	setting.Index().Render(r.Context(), w)
}
