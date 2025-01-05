package handlers

import (
	"db-platform/storage"
	"db-platform/views/setting"
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

	err = setting.LDAPSettingsPage(*config).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong: %v", err)))
		return
	}
}
