package handlers

import (
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"net/http"
)

func LoadResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["resetPassword"].Execute(w, nil)
}

func ResetPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)

	var user models.AppUser
	user.GetUserById(userId)

	currentPassword := r.FormValue("password")
	fmt.Println(currentPassword)
	newPassword := r.FormValue("new-password")
	fmt.Println(newPassword)
	w.Write([]byte("gogo"))
}
