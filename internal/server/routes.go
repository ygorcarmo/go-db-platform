package server

import (
	"custom-db-platform/internal/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type userDetails struct {
	FirstName string
	LastName  string
}

var clerkClient clerk.Client
var clerkError error

func (s *Server) RegisterRoutes() http.Handler {
	clerkClient, clerkError = clerk.NewClient(os.Getenv("clerk"))
	authenticateSession := utils.CustomRequireSessionV2(clerkClient)

	if clerkError != nil {
		fmt.Println(clerkError)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		r.Use(authenticateSession)

		r.Get("/", homeHandler)

		r.Get("/create-user", createUserPageHandler)
		r.Post("/create-user", createUserFormHandler)

		r.Get("/delete-user", deleteUserPageHandler)
		r.Post("/delete-user", deleteUserFormHandler)
	})

	// TODO : Remove this route as the clerk authenticator will handle sign-ins
	router.Get("/sign-in", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("../web/signIn.tmpl", "../web/base.tmpl")

		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
	})

	return router

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessClaims, _ := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)

	user, err := clerkClient.Users().Read(sessClaims.Subject)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome " + *user.FirstName)

	tmpl, err := template.ParseFiles("internal/web/index.tmpl", "internal/web/base.tmpl", "internal/web/userButtom.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	userDet := userDetails{*user.FirstName, *user.LastName}
	err = tmpl.Execute(w, userDet)
}
func createUserPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/web/createUser.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
}

func createUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, password: %s, wo: %s, databases: %v\n", username, password, wo, databases)
}

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("internal/web/deleteUser.tmpl")

	if err != nil {
		log.Fatal(err)
	}

	// TODO maybe have a different layout for partial and full reloads.
	hxHeader := r.Header.Get("HX-Request")
	if hxHeader != "" {
		fmt.Println("Partial Reload")
		// return
	}

	err = tmpl.Execute(w, nil)
	fmt.Println("Full Reload")
}

func deleteUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)

}