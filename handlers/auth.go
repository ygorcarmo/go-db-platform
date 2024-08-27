package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/storage"
	"github.com/ygorcarmo/db-platform/utils"
	"github.com/ygorcarmo/db-platform/views/components"
	"github.com/ygorcarmo/db-platform/views/login"
)

// should we get this from a enviroment variable?
var hashParams = utils.HashParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func GetLoginPage(w http.ResponseWriter, r *http.Request) {
	err := login.Index().Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when rendering the login page")
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	// TODO: add validation
	username := r.FormValue("username")
	p := r.FormValue("password")

	user, err := s.GetUserByUsername(username)
	if err != nil {
		fmt.Println(err)
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}
	match, err := hashParams.ComparePasswordAndHash(p, user.Password)
	if err != nil || !match {
		fmt.Printf("Unable to compare hash: %v", err)
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	tokenString, err := utils.CreateToken(user.Id)
	if err != nil {
		fmt.Printf("Unable to create token: %v", err)
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(6 * time.Hour)}

	http.SetCookie(w, &cookie)

	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Logged In."))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
