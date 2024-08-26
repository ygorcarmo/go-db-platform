package handlers

import (
	"log"
	"net/http"

	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/storage"
	"github.com/ygorcarmo/db-platform/views/components"
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

func GetCreateUserPage(w http.ResponseWriter, r *http.Request) {
	err := appUser.CreateUserPage().Render(r.Context(), w)
	if err != nil {
		log.Fatal("error when trying to render create user page")
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	u := r.FormValue("username")
	p := r.FormValue("password")
	p2 := r.FormValue("re-password")
	sup := r.FormValue("supervisor")
	sec := r.FormValue("sector")
	isA := r.FormValue("admin")

	if p != p2 {
		components.Response(models.CreateResponse("Passwords do not match", false)).Render(r.Context(), w)
		return
	}

	isAdmin := false
	if isA != "" {
		isAdmin = true
	}

	encodedHash, err := hashParams.HashPasword(p)
	if err != nil {
		components.Response(models.CreateResponse("Something went wrong when hashing the password.", false)).Render(r.Context(), w)
		return
	}

	user := models.AppUser{Username: u, Password: encodedHash, Supervisor: sup, Sector: sec, IsAdmin: isAdmin}

	err = s.CreateApplicationUser(user)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	components.Response(models.CreateResponse("User created successfully.", true)).Render(r.Context(), w)

}
