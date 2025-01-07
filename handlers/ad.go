package handlers

import (
	"db-platform/models"
	"db-platform/storage"
	"db-platform/views/components"
	"db-platform/views/settings"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetADConfigPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	config, err := s.GetADConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong: %v", err)))
		return
	}

	err = settings.LDAPSettingsPage(*config).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong: %v", err)))
		return
	}
}

func UpdateADConfigHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	// TODO: Add validation
	connectionStr := r.FormValue("connectionStr")
	domain := strings.Split(r.FormValue("domain"), ".")
	baseGroup := r.FormValue("baseGroup")
	baseGroupOU := r.FormValue("baseGroupOU")
	adminGroup := r.FormValue("adminGroup")
	adminGroupOU := r.FormValue("adminGroupOU")
	t := r.FormValue("timeOutInSecs")
	timeOutInSecs, err := strconv.Atoi(t)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}
	isDefault := r.FormValue("isDefault")

	config := models.LDAP{
		ConnectionStr:     connectionStr,
		TopLevelDomain:    domain[0],
		SecondLevelDomain: domain[1],
		BaseGroup:         baseGroup,
		BaseGroupOU:       baseGroupOU,
		AdminGroup:        adminGroup,
		AdminGroupOU:      adminGroupOU,
		TimeOutInSecs:     timeOutInSecs,
		IsDefault:         isDefault != "",
	}

	err = s.UpdateADConfig(config)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	err = components.Response(models.CreateResponse("AD configuration has been updated", true)).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong: %v", err)))
		return
	}

}

func TestConnectionHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	config, err := s.GetADConfigWithCredentials()
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	conn, err := config.Connect()
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}
	defer conn.Close()

	err = components.Response(models.CreateResponse("The connection was successful", true)).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong: %v", err)))
		return
	}
}

func GetADCredentialsPage(w http.ResponseWriter, r *http.Request) {
	err := settings.UpdateLDAPCredentialsPage().Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong: %v", err)))
		return
	}
}

func UpdateADCredentialsHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	username := strings.TrimSpace(r.FormValue("username"))
	passwd := strings.TrimSpace(r.FormValue("password"))
	passwd2 := strings.TrimSpace(r.FormValue("re-password"))

	if username == "" || passwd == "" || passwd2 == "" {
		components.Response(models.Response{Message: "Invalid Username or Password.", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	if passwd != passwd2 {
		components.Response(models.Response{Message: "Passwords do not match", IsSuccess: false}).Render(r.Context(), w)
		return
	}

	err := s.UpdateADCredentials(username, passwd)
	if err != nil {
		components.Response(models.Response{Message: err.Error(), IsSuccess: false}).Render(r.Context(), w)
		return
	}

	w.Header().Add("HX-Redirect", "/settings/ldap")
	w.Write([]byte("success"))
}
