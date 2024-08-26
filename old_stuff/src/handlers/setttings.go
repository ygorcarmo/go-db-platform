package handlers

import (
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"net/http"
)

func LoadSettings(w http.ResponseWriter, r *http.Request) {
	// views.Templates["settings"].Execute(w, nil)
	w.Write([]byte("GO"))
}

func LoadManageDbs(w http.ResponseWriter, r *http.Request) {
	var database models.ExternalDb
	databases, err := database.GetAll()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(databases)

	views.Templates["manageDbs"].Execute(w, databases)
}

func LoadManageUsers(w http.ResponseWriter, r *http.Request) {
	var users models.AppUser
	res, err := users.GetAllUsers()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	views.Templates["manageUsers"].Execute(w, res)
}

func LoadManageLogs(w http.ResponseWriter, r *http.Request) {
	var log models.Log
	logs, err := log.GetAllLogsPretty()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(logs)
	views.Templates["manageLogs"].Execute(w, logs)
}
