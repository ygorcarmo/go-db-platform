package handlers

import (
	"fmt"
	"net/http"

	"db-platform/storage"
	"db-platform/views/settings"
)

func GetAdminLogsPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	ls, err := s.GetAllAdminLogs()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = settings.AdminLogsPage(ls).Render(r.Context(), w)
	if err != nil {
		fmt.Println("Something went wrong went loading the Admin Logs Page")
	}
}
