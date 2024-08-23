package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/storage"
	externaldb "github.com/ygorcarmo/db-platform/views/externalDb"
)

type filteredResults struct {
	Sucesses []string
	Errors   []string
}

var wg sync.WaitGroup

func GetCreateDbUserPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetAvailableDbs()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Went Wrong"))
		return
	}

	err := externaldb.CreateUserPage(data).Render(r.Context(), w)

	if err != nil {
		log.Fatal("Error when rendering Create DB User Page")
	}
}

func GetDeleteDbUserPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetAvailableDbs()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Went Wrong"))
		return
	}

	err := externaldb.DeleteUserPage(data).Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when rendering delete db user page")
	}
}

func CreateDBUserHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	r.ParseForm()
	username := r.FormValue("username")
	wo := r.FormValue("wo")
	dbNames := r.Form["databases"]

	woInt, err := strconv.Atoi(wo)
	if err != nil {
		fmt.Println(err)
	}

	user := r.Context().Value(models.UserCtx).(*models.AppUser)

	var results []models.TargetDbsResponse

	c := make(chan models.TargetDbsResponse)

	for _, dbName := range dbNames {
		wg.Add(1)

		currentDb, err := s.GetDbByName(dbName)
		if err != nil {
			log.Fatal(err)
		}

		go currentDb.ConnectAndCreateUser(models.NewDbUserProps{Username: username, CurrentUserId: user.Id, WO: woInt}, c, &wg)
		msg := <-c
		results = append(results, msg)
	}

	wg.Wait()

	var fResponse filteredResults

	for _, result := range results {
		go s.CreateLog(models.Log{DbId: result.DbId, NewUser: username, WO: woInt, CreateBy: user.Id, Action: models.Create, Sucess: result.Success})
		if result.Success {
			fResponse.Sucesses = append(fResponse.Sucesses, result.Message)
		} else {
			fResponse.Errors = append(fResponse.Errors, result.Message)
		}
	}

	externaldb.Response(fResponse.Sucesses, fResponse.Errors).Render(r.Context(), w)
}
