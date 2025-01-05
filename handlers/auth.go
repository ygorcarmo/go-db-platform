package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
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
	username := strings.TrimSpace(r.FormValue("username"))
	p := strings.TrimSpace(r.FormValue("password"))

	if username == "" || p == "" {
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

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

	tokenString, err := utils.CreateToken(user.Id, "", "")
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

func GetADLoginPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {

	config, err := s.GetADConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("AD is not configured. Please Contact your Administrator."))
		return
	}

	ad := models.LDAP{
		ConnectionStr:     config.ConnectionStr,
		Username:          config.Username,
		Password:          config.Password,
		TopLevelDomain:    config.TopLevelDomain,
		SecondLevelDomain: config.SecondLevelDomain,
	}

	conn, err := ad.Connect()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("AD is not configured. Please contact your Administrator."))
		return
	}

	defer conn.Close()

	err = login.AD().Render(r.Context(), w)
	if err != nil {
		fmt.Printf("Error when rendering the AD login page: %s\n", err)
	}
}

func HandleADLogin(w http.ResponseWriter, r *http.Request, s storage.Storage) {

	config, err := s.GetADConfig()
	if err != nil {
		components.Response(models.Response{Message: "Error Connecting to AD. Please contact your Administrator", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	ad := models.LDAP{
		ConnectionStr:     config.ConnectionStr,
		Username:          config.Username,
		Password:          config.Password,
		TopLevelDomain:    config.TopLevelDomain,
		SecondLevelDomain: config.SecondLevelDomain,
		BaseGroup:         config.BaseGroup,
		BaseGroupOU:       config.BaseGroupOU,
		AdminGroup:        config.AdminGroup,
		AdminGroupOU:      config.AdminGroupOU,
	}

	conn, err := ad.Connect()
	if err != nil {
		components.Response(models.Response{Message: "Error Connecting to AD. Please contact your Administrator", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	// check if this works throught the function without having to open another connection
	defer conn.Close()

	username := strings.TrimSpace(r.FormValue("username"))
	passwd := strings.TrimSpace(r.FormValue("password"))

	if username == "" || passwd == "" {
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	err = ad.Authenticate(conn, username, passwd)

	if err != nil {
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	admins, err := ad.GetGroupMembers(conn, ad.AdminGroup, ad.AdminGroupOU)
	if err != nil {
		components.Response(models.Response{Message: "Error Connecting to AD. Please contact your Administrator", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	// To avoid going through another loop just check if the user is part of the admin and redirect if true
	if slices.Contains(admins, username) {
		tokenStr, err := utils.CreateToken("", username, ad.AdminGroup)
		if err != nil {
			fmt.Printf("Unable to create token: %v", err)
			components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenStr,
			HttpOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Add("HX-Redirect", "/")
		w.Write([]byte("Logged In."))
	}

	baseUsers, err := ad.GetGroupMembers(conn, ad.BaseGroup, ad.BaseGroupOU)
	if err != nil {
		components.Response(models.Response{Message: "Error Connecting to AD. Please contact your Administrator", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	if !slices.Contains(baseUsers, username) {
		components.Response(models.Response{Message: "User does not have access to this feature. Please contact your Administrator", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	tokenStr, err := utils.CreateToken("", username, ad.BaseGroup)
	if err != nil {
		fmt.Printf("Unable to create token: %v", err)
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("Logged In."))

}
