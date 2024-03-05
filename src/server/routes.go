package server

import (
	"custom-db-platform/src/db"
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
		r.Use(customMiddleware())

		r.Get("/", handlers.LoadHomePage)

		r.Get("/create-user", handlers.LoadExternalCreateUserPage)
		r.Post("/create-user", handlers.CreateExternalUserFormHandler)

		r.Get("/delete-user", deleteUserPageHandler)
		r.Post("/delete-user", deleteUserFormHandler)

		r.Get("/configuration", configPageHandler)
		r.Post("/db", addDBHandler)

	})

	return router
}

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {

	var dbNames models.TargetDb
	dbs, err := dbNames.GetAllNames()
	if err != nil {
		log.Fatal(err)
	}

	views.Templates["deleteUser"].Execute(w, dbs)
}

func deleteUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)
}

func configPageHandler(w http.ResponseWriter, r *http.Request) {
	views.Templates["config"].Execute(w, nil)
}

func addDBHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	host := r.FormValue("host")
	port := r.FormValue("port")
	dbType := r.FormValue("type")
	sslMode := r.FormValue("sslMode")

	_, err := db.Database.Exec("INSERT INTO db_connection_info (name, host, port, type, sslMode) VALUES (?, ?, ?, ?, ?)", name, host, port, dbType, sslMode)

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(fmt.Sprintf("<div class=\"border border-red-500 bg-red-300 w-fit p-2 rounded\">%v.</div>", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("<div class=\"border border-green-500 bg-green-300 w-fit p-2 rounded\">%s has been created successfully.</div>", name)))
}
