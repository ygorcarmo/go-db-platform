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

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoadResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["resetPassword"].Execute(w, nil)
}

func ResetPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(string)

	var user models.AppUser
	user.GetUserById(userId)

	currentPassword := r.FormValue("password")
	newPassword := r.FormValue("new-password")
	reEnteredPassword := r.FormValue("re-password")

	match, err := hashParams.ComparePasswordAndHash(currentPassword, user.Password)
	if err != nil || !match {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid Password."))
		return
	}

	if newPassword != reEnteredPassword {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Passwords do not match."))
		return
	}

	hashedNewPassword, err := hashParams.HashPasword(newPassword)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hashedNewPassword)

	passwdErr := user.UpdatePassword(hashedNewPassword)
	if passwdErr != nil {
		fmt.Println(passwdErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	w.Write([]byte("password updated"))
}
