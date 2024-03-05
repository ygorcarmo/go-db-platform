package server

import (
	"custom-db-platform/src/handlers"
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/sign-in", handlers.LoadSignInPage)
	router.Post("/sign-in", handlers.HandleSignIn)

	router.Group(func(r chi.Router) {
		// TODO change this with jwt token decode and only make a few routes available(only admins are able to access)
		r.Use(customMiddleware())

		r.Get("/", handlers.LoadHomePage)

		r.Route("/db", func(dbRoute chi.Router) {

			dbRoute.Get("/", handlers.LoadAddDbPage)
			dbRoute.Post("/", handlers.AddDbFormHanlder)

			dbRoute.Get("/create-user", handlers.LoadCreateExternalUserPage)
			dbRoute.Post("/create-user", handlers.CreateExternalUserFormHandler)

			dbRoute.Get("/delete-user", deleteUserPageHandler)
			dbRoute.Post("/delete-user", deleteUserFormHandler)
		})
	})

	return router
}

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {

	var dbNames models.TargetDb
	dbs, err := dbNames.GetAllNames()
	if err != nil {
		log.Fatal(err)
	}

	views.Templates["deleteUserPage"].Execute(w, dbs)
}

func deleteUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)
}
