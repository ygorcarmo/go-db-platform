package handlers

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadCreateAppUser(w http.ResponseWriter, r *http.Request) {
	views.Templates["addAppUser"].Execute(w, nil)
}

func AddAppUserFormHanlder(w http.ResponseWriter, r *http.Request) {
	// TODO make sure passwords match and hash the password before creating user
	username := r.FormValue("username")
	password := r.FormValue("password")
	// rePassword := r.FormValue("re-password")
	supervisor := r.FormValue("supervisor")
	sector := r.FormValue("sector")
	admin := r.FormValue("admin")

	var isAdmin bool

	if admin == "" {
		isAdmin = false
	} else {
		isAdmin = true
	}

	encodedHash, err := hashParams.HashPasword(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encodedHash)

	user := models.AppUser{
		Username:   username,
		Password:   encodedHash,
		Supervisor: supervisor,
		Sector:     sector,
		IsAdmin:    isAdmin,
	}
	userErr := user.CreateUser()
	if userErr != nil {
		fmt.Println(userErr)
	}

	w.Write([]byte("user created"))

}

func LoadEditAppUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")

	var user models.AppUser
	user.GetUserById(userId)
	views.Templates["editAppUser"].Execute(w, user)
}

func DeleteAppUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")

	_, err := db.Database.Exec("DELETE FROM users WHERE id = UUID_TO_BIN(?);", userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("somthing went wrong"))
	}
	w.WriteHeader(http.StatusOK)
}
