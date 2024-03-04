package handlers

import (
	"custom-db-platform/src/views"
	"fmt"
	"net/http"
)

func LoadSignInPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["signIn"].Execute(w, nil)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Printf("Username is %s, Password is %s\n", username, password)
	w.Write([]byte("Hello"))
}
