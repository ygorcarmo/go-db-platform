package handlers

import (
	"custom-db-platform/src/views"
	"net/http"
)

func LoadSettings(w http.ResponseWriter, r *http.Request) {
	views.Templates["settings"].Execute(w, nil)
}

func LoadManageDbs(w http.ResponseWriter, r *http.Request) {
	views.Templates["manageDbs"].Execute(w, nil)
}

func LoadManageUsers(w http.ResponseWriter, r *http.Request) {
	views.Templates["manageUsers"].Execute(w, nil)
}
