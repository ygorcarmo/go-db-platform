package server

import (
	"custom-db-platform/src/datatypes"
	"custom-db-platform/src/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

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

		r.Get("/configuration", configPageHandler)
		r.Post("/db", addDBHandler)
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

	fmt.Printf("Welcome %s\n", *user.FirstName)

	tmpl, err := template.ParseFiles("src/web/index.tmpl", "src/web/base.tmpl", "src/web/userButtom.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	userDet := userDetails{*user.FirstName, *user.LastName}
	err = tmpl.Execute(w, userDet)
}
func createUserPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/web/createUser.tmpl")

	if err != nil {
		log.Fatal(err)
	}

	dbs, err := getDBsName()

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, dbs)
}

func createUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	var results []datatypes.Result

	c := make(chan datatypes.Result)
	var wg sync.WaitGroup

	for _, database := range databases {
		wg.Add(1)

		dbDetail, err := getDBByName(database)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dbDetail)
		fmt.Printf("username: %s, wo: %s, database: %v\n", username, wo, dbDetail)
		go utils.ConnectToDBAndCreateUser(dbDetail.Host, dbDetail.Port, dbDetail.DbType, dbDetail.SslMode, username, c, &wg)
		fmt.Println(<-c)
		results = append(results, <-c)
	}
	wg.Wait()
	fmt.Println(results)
}

func deleteUserPageHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("src/web/deleteUser.tmpl")

	if err != nil {
		log.Fatal(err)
	}

	dbs, err := getDBsName()

	if err != nil {
		log.Fatal(err)
	}

	// TODO maybe have a different layout for partial and full reloads.
	hxHeader := r.Header.Get("HX-Request")
	if hxHeader != "" {
		fmt.Println("Partial Reload")
		// return
	}

	err = tmpl.Execute(w, dbs)
	fmt.Println("Full Reload")
}

func deleteUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)
}

func configPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/web/config.tmpl")

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
}

func addDBHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	host := r.FormValue("host")
	port := r.FormValue("port")
	dbType := r.FormValue("type")
	sslMode := r.FormValue("sslMode")

	_, err := db.Exec("INSERT INTO db_connection_info (name, host, port, type, sslMode) VALUES (?, ?, ?, ?, ?)", name, host, port, dbType, sslMode)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success")
}
