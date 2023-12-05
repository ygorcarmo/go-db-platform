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
	injectActiveSession := clerk.WithSessionV2(client)

	if err != nil {
		// handle error
		fmt.Println(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/index.html", "templates/base.html")

		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
		// w.Write([]byte("Welcome"))
	})
	router.Get("/sign-in",	func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/signIn.html", "templates/base.html")

		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
		// w.Write([]byte("Welcome"))
	})


	router.Handle("/admin", injectActiveSession(helloUserHandler(client)))
	http.ListenAndServe("127.0.0.1:3000", router)
}
func helloUserHandler(client clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
		if !ok {
			// w.WriteHeader(http.StatusUnauthorized)
			// w.Write([]byte("Unauthorized"))
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		user, err := client.Users().Read(sessClaims.Subject)
		if err != nil {
			panic(err)
		}

		w.Write([]byte("Welcome " + *user.FirstName))
	}
}
