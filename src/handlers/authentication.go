package handlers

import (
	"custom-db-platform/src/models"
	"custom-db-platform/src/utils"
	"custom-db-platform/src/views"
	"fmt"
	"net/http"
)

var hashParams = models.HashParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func LoadSignInPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["signIn"].Execute(w, nil)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var currentUser models.AppUser
	err := currentUser.GetUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid Username or Password."))
		return
	}

	match, err := hashParams.ComparePasswordAndHash(password, currentUser.Password)
	if err != nil || !match {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid Username or Password."))
		return
	}

	tokenString, err := utils.CreateToken(currentUser.Id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid Username or Password."))
		return
	}

	cookie := http.Cookie{Name: "token", Value: tokenString}
	http.SetCookie(w, &cookie)

	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Loged In."))
}
