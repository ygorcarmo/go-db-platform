package server

import (
	"custom-db-platform/src/handlers"
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed assets/**
var assets embed.FS

func (s *Server) RegisterRoutes() http.Handler {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// assetsFolder, err := fs.Sub(assets, "assets")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	fs := http.FileServer(http.FS(assets))

	router.Handle("/assets/*", fs)

	router.Get("/sign-in", handlers.LoadSignInPage)
	router.Post("/sign-in", handlers.HandleSignIn)

	router.Group(func(r chi.Router) {
		// TODO change this with jwt token decode and only make a few routes available(only admins are able to access)
		r.Use(verifyUserMiddleware())

		r.Get("/", handlers.LoadHomePage)

		r.Get("/reset-password", handlers.LoadResetPasswordPage)

		r.Route("/db", func(r chi.Router) {

			r.Get("/", handlers.LoadAddDbPage)
			r.Post("/", handlers.AddDbFormHanlder)

			r.Get("/create-user", handlers.LoadCreateExternalUser)
			r.Post("/create-user", handlers.CreateExternalUserFormHandler)

			r.Get("/delete-user", deleteUserPageHandler)
			r.Post("/delete-user", deleteUserFormHandler)
		})

		r.Route("/settings", func(adminsOnlyRoute chi.Router) {
			adminsOnlyRoute.Use(adminsOnlyMiddleware())
			adminsOnlyRoute.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Welcome")) })
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

	views.Templates["deleteUser "].Execute(w, dbs)
}

func deleteUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)
}
