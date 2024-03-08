package handlers

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LoadAddDb(w http.ResponseWriter, r *http.Request) {
	views.Templates["addDb"].Execute(w, nil)
}

func AddDbFormHanlder(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	host := r.FormValue("host")
	port := r.FormValue("port")
	dbType := r.FormValue("type")
	sslMode := r.FormValue("sslMode")

	_, err := db.Database.Exec("INSERT INTO external_databases (name, host, port, type, sslMode) VALUES (?, ?, ?, ?, ?)", name, host, port, dbType, sslMode)

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(fmt.Sprintf("<div class=\"border border-red-500 bg-red-300 w-fit p-2 rounded\">%v.</div>", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("<div class=\"border border-green-500 bg-green-300 w-fit p-2 rounded\">%s has been created successfully.</div>", name)))
}

func LoadEditDb(w http.ResponseWriter, r *http.Request) {
	dbId := chi.URLParam(r, "id")
	var db models.TargetDb
	db.GetByid(dbId)
	views.Templates["editDb"].Execute(w, db)
}

func DeleteDb(w http.ResponseWriter, r *http.Request) {
	dbId := chi.URLParam(r, "id")

	_, err := db.Database.Exec("DELETE FROM external_databases WHERE id=UUID_TO_BIN(?)", dbId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("somthing went wrong"))
	}

	w.WriteHeader(http.StatusOK)

}
