package handlers

import (
	"net/http"

	"db-platform/storage"
)

func SeedHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	err := s.Seed()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DB Seeded"))

}
