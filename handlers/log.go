package handlers

import (
	"fmt"
	"net/http"

	"db-platform/storage"
	"db-platform/views/setting"
)

func GetLogsPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	logs, err := s.GetAllLogs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	err = setting.LogsPage(logs).Render(r.Context(), w)
	if err != nil {
		fmt.Println("Something went wrong when trying to render logs config page")
	}
}
