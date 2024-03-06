package handlers

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/views"
	"fmt"
	"net/http"
)

func LoadAddDbPage(w http.ResponseWriter, r *http.Request) {
	views.Templates["addDb"].Execute(w, nil)
}

func AddDbFormHanlder(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	host := r.FormValue("host")
	port := r.FormValue("port")
	dbType := r.FormValue("type")
	sslMode := r.FormValue("sslMode")

	_, err := db.Database.Exec("INSERT INTO db_connection_info (name, host, port, type, sslMode) VALUES (?, ?, ?, ?, ?)", name, host, port, dbType, sslMode)

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(fmt.Sprintf("<div class=\"border border-red-500 bg-red-300 w-fit p-2 rounded\">%v.</div>", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("<div class=\"border border-green-500 bg-green-300 w-fit p-2 rounded\">%s has been created successfully.</div>", name)))
}
