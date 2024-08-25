package handlers

import (
	"log"
	"net/http"

	"github.com/ygorcarmo/db-platform/storage"
	"github.com/ygorcarmo/db-platform/views/setting"
	appUser "github.com/ygorcarmo/db-platform/views/user"
)

func GetResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	err := appUser.ResetPassword().Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when trying to render reset password page")
	}
}

func GetAllUserSettingsPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	users, err := s.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	setting.GetUsersPage(users).Render(r.Context(), w)
}
