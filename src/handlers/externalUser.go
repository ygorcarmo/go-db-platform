package handlers

import (
	"custom-db-platform/src/models"
	"custom-db-platform/src/views"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type filteredResults struct {
	Sucesses []string
	Errors   []string
}

var wg sync.WaitGroup

var targetDbs models.TargetDb

func LoadCreateExternalUser(w http.ResponseWriter, r *http.Request) {
	dbs, err := targetDbs.GetAllNames()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dbs)

	views.Templates["createDbUser"].Execute(w, dbs)
}

func CreateExternalUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	woInt, err := strconv.Atoi(wo)
	if err != nil {
		fmt.Println(err)
	}

	userId := r.Context().Value("userId").(string)

	var results []models.TargetDbsRepose

	c := make(chan models.TargetDbsRepose)

	for _, database := range databases {
		wg.Add(1)
		var currentDb models.TargetDb
		_, err := currentDb.GetByName(database)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("username: %s, wo: %s, database: %v\n", username, wo, currentDb)
		go currentDb.ConnectToDBAndCreateUser(username, userId, woInt, c, &wg)
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
	views.Templates["dbUserFormResponse"].Execute(w, fResponse)
}

func LoadDeleteExternalUser(w http.ResponseWriter, r *http.Request) {
	var dbNames models.TargetDb
	dbs, err := dbNames.GetAllNames()
	if err != nil {
		log.Fatal(err)
	}

	views.Templates["deleteDbUser"].Execute(w, dbs)
}

func DeleteExternalUserFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	databases := r.Form["databases"]

	fmt.Printf("username: %s, wo: %s, databases: %v\n", username, wo, databases)
}
