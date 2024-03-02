package server

import (
	"custom-db-platform/src/web"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type userDetails struct {
	FirstName string
	LastName  string
}

type Result struct {
	Message string
	Success bool
}

type filteredResults struct {
	Sucesses []string
	Errors   []string
}

var wg sync.WaitGroup

func (s *Server) RegisterRoutes() http.Handler {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		// r.Use(utils.CustomMiddleware())

		r.Get("/", homeHandler)

		r.Get("/create-user", createUserPageHandler)
		r.Post("/create-user", createUserFormHandler)

		r.Get("/delete-user", deleteUserPageHandler)
		r.Post("/delete-user", deleteUserFormHandler)

		r.Get("/configuration", configPageHandler)
		r.Post("/db", addDBHandler)
	})

	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Do I need the user details? If yes where am I going to get it from?
	// the request context?
	// ctx := r.Context()
	// sessClaims, _ := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)

	// user, err := clerkClient.Users().Read(sessClaims.Subject)
	// if err != nil {
	// 	panic(err)
	// }

	web.Templates["home"].ExecuteTemplate(w, "base-layout.tmpl", nil)
}
func createUserPageHandler(w http.ResponseWriter, r *http.Request) {

	dbs, err := getDBsName()
	if err != nil {
		log.Fatal(err)
	}

	web.Templates["createUser"].ExecuteTemplate(w, "base-layout.tmpl", dbs)
}

func createUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]
	fmt.Println(databases)

	var results []Result

	c := make(chan Result)

	for _, database := range databases {
		wg.Add(1)

		dbDetail, err := getDBByName(database)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("username: %s, wo: %s, database: %v\n", username, wo, dbDetail)
		go ConnectToDBAndCreateUser(dbDetail.Host, dbDetail.Port, dbDetail.DbType, dbDetail.SslMode, username, dbDetail.Name, c)
		msg := <-c
		results = append(results, msg)
	}
	wg.Wait()

	var fResponse filteredResults

	for _, result := range results {
		if result.Success {
			fResponse.Sucesses = append(fResponse.Sucesses, result.Message)
		} else {
			fResponse.Errors = append(fResponse.Errors, result.Message)
		}
	}

	tmpl, err := template.ParseFiles("src/web/response.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, fResponse)
}

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {

	dbs, err := getDBsName()

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

	_, err := db.Exec("INSERT INTO db_connection_info (name, host, port, type, sslMode) VALUES (?, ?, ?, ?, ?)", name, host, port, dbType, sslMode)

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(fmt.Sprintf("<div class=\"border border-red-500 bg-red-300 w-fit p-2 rounded\">%v.</div>", err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("<div class=\"border border-green-500 bg-green-300 w-fit p-2 rounded\">%s has been created successfully.</div>", name)))
}
