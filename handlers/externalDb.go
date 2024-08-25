package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/storage"
	"github.com/ygorcarmo/db-platform/views/components"
	externaldb "github.com/ygorcarmo/db-platform/views/externalDb"
	"github.com/ygorcarmo/db-platform/views/setting"
)

var wg sync.WaitGroup

func GetCreateDbUserPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetDbsName()
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
	data, derr := db.GetDbsName()
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

func ExternalDBUserHandler(w http.ResponseWriter, r *http.Request, s storage.Storage, a models.ActionType) {

	// TODO: add server side validation
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	wo := r.FormValue("wo")
	dbNames := r.Form["databases"]

	if len(dbNames) < 1 {
		externaldb.Response([]string{}, []string{"Please select a database"}).Render(r.Context(), w)
		return
	}

	woInt, err := strconv.Atoi(wo)
	if err != nil {
		fmt.Println(err)
	}

	user := r.Context().Value(models.UserCtx).(*models.AppUser)

	successr := []string{}
	failr := []string{}

	for _, dbName := range dbNames {
		wg.Add(1)

		currentDb, err := s.GetDbByName(dbName)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer wg.Done()
			result := models.TargetDbsResponse{}
			switch a {
			case models.Create:
				result = currentDb.ConnectAndCreateUser(models.NewDbUserProps{Username: username, CurrentUserId: user.Id, WO: woInt, Password: password})
			case models.Delete:
				result = currentDb.ConnectAndDeleteUser(models.NewDbUserProps{Username: username, CurrentUserId: user.Id, WO: woInt})
			default:
				fmt.Println("Action Type not supported")
				result = models.TargetDbsResponse{Message: "Action type not supported", IsSuccess: false, DbId: "NOTVALID"}
			}

			go s.CreateLog(models.Log{DbId: result.DbId, NewUser: username, WO: woInt, CreateBy: user.Id, Action: a, Sucess: result.IsSuccess})

			if result.IsSuccess {
				successr = append(successr, result.Message)
			} else {
				failr = append(failr, result.Message)
			}
		}()
	}

	wg.Wait()

	externaldb.Response(successr, failr).Render(r.Context(), w)
}

func GetDatabasesConfigPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	dbs, err := s.GetAllDbs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Somthing went wrong"))
		return
	}

	renderErr := setting.DatabasesPage(dbs).Render(r.Context(), w)
	if renderErr != nil {
		log.Fatalln("Something went wrong went trying to render the dbs config page")
	}
}

func GetCreateExternalDbPage(w http.ResponseWriter, r *http.Request) {
	err := externaldb.ExternalDbPage().Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when trying to render create external db page")
	}
}

func CreateExternalDbHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	user := r.Context().Value(models.UserCtx).(*models.AppUser)

	// TODO: add server side validation
	u := r.FormValue("username")
	p := r.FormValue("password")
	d := r.FormValue("name")
	h := r.FormValue("host")
	dp := r.FormValue("port")
	t := r.FormValue("type")
	m := r.FormValue("sslMode")

	dType, err := models.ToDbType(t)
	if err != nil {
		components.Response(models.Response{Message: err.Error(), IsSuccess: false}).Render(r.Context(), w)
		return
	}

	dPort, err := strconv.Atoi(dp)
	if err != nil {
		components.Response(models.Response{Message: "Invalid Port. It should be a valid number"}).Render(r.Context(), w)
		return
	}

	config := models.TargetDb{Username: u, Password: p, Name: d, Host: h, Port: dPort, Type: dType, SslMode: m, CreatedBy: user.Id}

	err = s.CreateExternalDb(config)

	if err != nil {
		components.Response(models.Response{Message: err.Error(), IsSuccess: false}).Render(r.Context(), w)
		return
	}

	components.Response(models.Response{Message: "DB Connection config has been created", IsSuccess: true}).Render(r.Context(), w)
}
