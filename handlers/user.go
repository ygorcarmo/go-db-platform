package handlers

import (
	"log"
	"net/http"

	"github.com/ygorcarmo/db-platform/views/user"
)

func GetResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	err := user.ResetPassword().Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when trying to render reset password page")
	}
}
