package handlers

import (
	"db-platform/models"
	"db-platform/storage"
	"db-platform/views/components"
	"db-platform/views/settings"
	"fmt"
	"net/http"
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

func TestConnectionHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	config, err := s.GetADConfig()
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
