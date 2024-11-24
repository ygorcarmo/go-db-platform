package handlers

import (
	"fmt"
	"net/http"
	"time"

	"db-platform/models"
	"db-platform/storage"
	"db-platform/utils"
	"db-platform/views/components"
	"db-platform/views/login"
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
		fmt.Printf("Error when rendering the login page: %s\n", err)
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	// TODO: add validation
	username := r.FormValue("username")
	p := r.FormValue("password")

	user, err := s.GetUserByUsername(username)
	if err != nil {
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	if user.LoginAttempts > 5 {
		components.Response(models.Response{Message: "Account is locked.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	match, err := hashParams.ComparePasswordAndHash(p, user.Password)
	if err != nil || !match {
		fmt.Printf("Unable to compare hash: %v", err)
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		go s.IncreaseUserLoginAttempts(user.Id, user.LoginAttempts+1)
		return
	}

	tokenString, err := utils.CreateToken(user.Id)
	if err != nil {
		fmt.Printf("Unable to create token: %v", err)
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	// reset login attempts on successfull login
	go func() {
		if user.LoginAttempts > 0 {
			s.ResetUserLoginAttempts(user.Id)
		}
	}()

	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Logged In."))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
