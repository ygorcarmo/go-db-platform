package server

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/handlers"
	"custom-db-platform/src/models"
	"custom-db-platform/src/web"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", handlers.LoadHomePage)

	router.Group(func(r chi.Router) {
		// r.Use(utils.CustomMiddleware())

		r.Get("/", handlers.LoadHomePage)

		r.Get("/create-user", handlers.LoadCreateUserForm)
		r.Post("/create-user", handlers.CreateUserFormHandler)

		r.Get("/delete-user", deleteUserPageHandler)
		r.Post("/delete-user", deleteUserFormHandler)

		r.Get("/configuration", configPageHandler)
		r.Post("/db", addDBHandler)

		r.Get("/sign-in", handleSignIn)
	})

	return router
}

// func createUserFormHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.FormValue("username")
// 	wo := r.FormValue("wo")
// 	databases := r.Form["databases"]
// 	fmt.Println(databases)

// 	var results []Result

// 	c := make(chan Result)

// 	for _, database := range databases {
// 		wg.Add(1)

// 		dbDetail, err := getDBByName(database)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Printf("username: %s, wo: %s, database: %v\n", username, wo, dbDetail)
// 		go ConnectToDBAndCreateUser(dbDetail.Host, dbDetail.Port, dbDetail.DbType, dbDetail.SslMode, username, dbDetail.Name, c)
// 		msg := <-c
// 		results = append(results, msg)
// 	}
// 	wg.Wait()

// 	var fResponse filteredResults

// 	for _, result := range results {
// 		if result.Success {
// 			fResponse.Sucesses = append(fResponse.Sucesses, result.Message)
// 		} else {
// 			fResponse.Errors = append(fResponse.Errors, result.Message)
// 		}
// 	}

// 	tmpl, err := template.ParseFiles("src/web/response.tmpl")

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = tmpl.Execute(w, fResponse)
// }

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {

	var dbNames models.TargetDb
	dbs, err := dbNames.GetAllNames()
	if err != nil {
		log.Fatal(err)
	}

	web.Templates["deleteUser"].ExecuteTemplate(w, "base-layout.tmpl", dbs)
}

func deleteUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)
}

func configPageHandler(w http.ResponseWriter, r *http.Request) {
	web.Templates["config"].ExecuteTemplate(w, "base-layout.tmpl", nil)
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

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	web.Templates["signIn"].ExecuteTemplate(w, "base-layout.tmpl", nil)
}
