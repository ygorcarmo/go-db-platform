package handlers

import (
	"log"
	"net/http"

	"github.com/ygorcarmo/db-platform/storage"
	externaldb "github.com/ygorcarmo/db-platform/views/externalDb"
)

func GetCreateDbUserPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetAvailableDbs()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Went Wrong"))
		return
	}

	err := externaldb.CreateUserPage(data).Render(r.Context(), w)

	if err != nil {
		log.Fatal("Error when rendering Create DB User Page")
	}
}

func GetDeleteDbUserPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetAvailableDbs()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Went Wrong"))
		return
	}

	err := externaldb.DeleteUserPage(data).Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when rendering delete db user page")
	}
}
