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

func main() {

	client, err := clerk.NewClient(os.Getenv("clerk"))
	authenticateSession := customRequireSessionV2(client)

	if err != nil {
		fmt.Println(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		r.Use(authenticateSession)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			sessClaims, _ := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)

			user, err := client.Users().Read(sessClaims.Subject)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Welcome " + *user.FirstName)

			tmpl, err := template.ParseFiles("templates/index.html", "templates/base.html", "templates/userButtom.html")

			if err != nil {
				log.Fatal(err)
			}
			err = tmpl.Execute(w, nil)

		})
	})

	router.Get("/sign-in", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/signIn.html", "templates/base.html")

		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
	})
	http.ListenAndServe("127.0.0.1:3000", router)
}

func customRequireSessionV2(client clerk.Client, verifyTokenOptions ...clerk.VerifyTokenOption) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
			if !ok || claims == nil {
				tmpl, err := template.ParseFiles("templates/signIn.html", "templates/base.html")

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
