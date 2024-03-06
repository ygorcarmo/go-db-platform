package handlers

import (
	"custom-db-platform/src/views"
	"net/http"
)

func LoadResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["resetPassword"].Execute(w, nil)
}
