package handlers

import (
	"fmt"
	"net/http"

	"db-platform/storage"
	"db-platform/views/setting"
)

func GetAdminLogsPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	ls, err := s.GetAllAdminLogs()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	err = setting.AdminLogsPage(ls).Render(r.Context(), w)
	if err != nil {
		fmt.Println("Something went wrong went loading the Admin Logs Page")
	}
}
