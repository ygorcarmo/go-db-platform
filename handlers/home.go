package handlers

import (
	"fmt"
	"net/http"

	"db-platform/views/home"
)

func GetHomePage(w http.ResponseWriter, r *http.Request) {
	err := home.Index().Render(r.Context(), w)

	if err != nil {
		fmt.Printf("Error when rendering the home page: %v\n", err)
	}
}
