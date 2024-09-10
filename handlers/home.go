package handlers

import (
	"log"
	"net/http"

	"db-platform/views/home"
)

func GetHomePage(w http.ResponseWriter, r *http.Request) {
	err := home.Index().Render(r.Context(), w)

	if err != nil {
		log.Fatal("Error when rendering the home page")
	}
}
