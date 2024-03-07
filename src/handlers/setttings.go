package handlers

import (
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"net/http"
)

func LoadSettings(w http.ResponseWriter, r *http.Request) {
	views.Templates["settings"].Execute(w, nil)
}

func LoadManageDbs(w http.ResponseWriter, r *http.Request) {
	views.Templates["manageDbs"].Execute(w, nil)
}

func LoadManageUsers(w http.ResponseWriter, r *http.Request) {
	var users models.AppUser
	res, err := users.GetAllUsers()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	views.Templates["manageUsers"].Execute(w, nil)
}
