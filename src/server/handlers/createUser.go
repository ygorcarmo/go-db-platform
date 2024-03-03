package handlers

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/web"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Result struct {
	Message string
	Success bool
}

type filteredResults struct {
	Sucesses []string
	Errors   []string
}

var wg sync.WaitGroup

func LoadCreateUserForm(w http.ResponseWriter, r *http.Request) {
	dbs, err := db.GetDBsName()
	if err != nil {
		log.Fatal(err)
	}

	web.Templates["createUserForm"].ExecuteTemplate(w, "base-layout.tmpl", dbs)
}

func CreateUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	var results []Result

	c := make(chan Result)

	for _, database := range databases {
		wg.Add(1)

		dbDetail, err := db.GetDBByName(database)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("username: %s, wo: %s, database: %v\n", username, wo, dbDetail)
		go db.ConnectToDBAndCreateUser(dbDetail.Host, dbDetail.Port, dbDetail.DbType, dbDetail.SslMode, username, dbDetail.Name, c, &wg)
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
	// add this to template loader
	tmpl, err := template.ParseFiles("src/web/response.tmpl")

	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, fResponse)
}
