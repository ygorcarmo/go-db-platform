package handlers

import (
	"custom-db-platform/src/views"
	"net/http"
)

func LoadHomePage(w http.ResponseWriter, r *http.Request) {
	views.Templates["home"].Execute(w, nil)
}
