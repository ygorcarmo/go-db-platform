package handlers

import (
	"custom-db-platform/src/web"
	"net/http"
)

func LoadHomePage(w http.ResponseWriter, r *http.Request) {
	web.Templates["Home"].ExecuteTemplate(w, "base-layout.tmpl", nil)
}
