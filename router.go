package main

import (
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

func router() {
	clerkClient, clerkError = clerk.NewClient(os.Getenv("clerk"))
	authenticateSession := customRequireSessionV2(clerkClient)

	if clerkError != nil {
		fmt.Println(clerkError)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		r.Use(authenticateSession)

		r.Get("/", homeHandler)
		r.Get("/create-user", createUserHandler)
		r.Get("/edit-user", editUserHandler)
	})

	// TODO : Remove this route as the clerk authenticator will handle sign-ins
	router.Get("/sign-in", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/signIn.tmpl", "templates/base.tmpl")

		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
	})
	fmt.Println("Server running on http://localhost:3000")
	http.ListenAndServe("127.0.0.1:3000", router)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessClaims, _ := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)

	user, err := clerkClient.Users().Read(sessClaims.Subject)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome " + *user.FirstName)

	tmpl, err := template.ParseFiles("templates/index.tmpl", "templates/base.tmpl", "templates/userButtom.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	userDet := userDetails{*user.FirstName, *user.LastName}
	err = tmpl.Execute(w, userDet)
}
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/createUser.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/editUser.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
}

func customRequireSessionV2(client clerk.Client, verifyTokenOptions ...clerk.VerifyTokenOption) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
			if !ok || claims == nil {
				tmpl, err := template.ParseFiles("templates/signIn.tmpl", "templates/base.tmpl")

				if err != nil {
					log.Fatal(err)
				}
				err = tmpl.Execute(w, nil)
				return
			}

			next.ServeHTTP(w, r)
		})

		return clerk.WithSessionV2(client, verifyTokenOptions...)(f)
	}
}
