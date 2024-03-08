package handlers

import (
	"custom-db-platform/src/views"
	"net/http"
)

func LoadCreateAppUser(w http.ResponseWriter, r *http.Request) {
	views.Templates["addAppUser"].Execute(w, nil)
}
