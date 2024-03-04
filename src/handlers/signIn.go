package handlers

import (
	"custom-db-platform/src/views"
	"net/http"
)

func LoadSignInPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["signIn"].Execute(w, nil)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
