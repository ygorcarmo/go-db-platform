package handlers

import (
	"custom-db-platform/src/models"
	"custom-db-platform/src/utils"
	"custom-db-platform/src/views"
	"fmt"
	"log"
	"net/http"
)

func LoadHomePage(w http.ResponseWriter, r *http.Request) {

	jwtToken, err := r.Cookie("token")

	if err != nil {
		log.Fatal(err)
	}

	userId, err := utils.DecodeToken(jwtToken.Value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userId)
	var user models.AppUser
	user.GetUserById(userId)
	fmt.Println(user)

	views.Templates["home"].Execute(w, user)
}
