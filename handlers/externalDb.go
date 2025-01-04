package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"db-platform/models"
	"db-platform/storage"
	"db-platform/views/components"
	"db-platform/views/externalDb"
	"db-platform/views/setting"

	"github.com/go-chi/chi/v5"
)

var wg sync.WaitGroup

func GetCreateDbUserPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetDbsName()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Went Wrong"))
		return
	}

	err := externalDb.CreateUserPage(data).Render(r.Context(), w)

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

	err := externalDb.DeleteUserPage(data).Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when rendering delete db user page")
	}
}
func GetUpdateDbUserPasswordPage(w http.ResponseWriter, r *http.Request, db storage.Storage) {
	data, derr := db.GetDbsName()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something Went Wrong"))
		return
	}

	err := externalDb.UpdateDbUserPasswordPage(data).Render(r.Context(), w)
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
		externalDb.Response([]string{}, []string{"Please select a database"}).Render(r.Context(), w)
		return
	}

	woInt, err := strconv.Atoi(wo)
	if err != nil {
		fmt.Printf("Error when converting WO: %s\n", err)
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
			result := models.ExternalDbResponse{}
			switch a {
			case models.Create:
				result = currentDb.ConnectAndCreateUser(models.NewDbUserProps{Username: username, CurrentUserId: user.Id, WO: woInt, Password: password})
			case models.Delete:
				result = currentDb.ConnectAndDeleteUser(models.NewDbUserProps{Username: username, CurrentUserId: user.Id, WO: woInt})
			case models.UPDATEPWD:
				result = currentDb.ConnectAndUpdateUserPassword(models.NewDbUserProps{Username: username, CurrentUserId: user.Id, WO: woInt, Password: password})
			default:
				fmt.Println("Action Type not supported")
				result = models.ExternalDbResponse{Message: "Action type not supported", IsSuccess: false, DbId: "NOTVALID"}
			}

			go s.CreateLog(models.Log{DbId: result.DbId, NewUser: username, WO: woInt, CreateBy: user.Id, Action: a, Success: result.IsSuccess})

			if result.IsSuccess {
				successr = append(successr, result.Message)
			} else {
				failr = append(failr, result.Message)
			}
		}()
	}

	wg.Wait()

	externalDb.Response(successr, failr).Render(r.Context(), w)
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
	err := externalDb.ExternalDbPage().Render(r.Context(), w)
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
	o := r.FormValue("owner")
	protocol := r.FormValue("protocol")
	hostFallback := r.FormValue("host-fallback")
	pFallback := r.FormValue("port-fallback")
	protocolFallback := r.FormValue("protocol-fallback")

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

	portFallback, err := strconv.Atoi(pFallback)
	if err != nil {
		components.Response(models.Response{Message: "Invalid Port. It should be a valid number"}).Render(r.Context(), w)
		return
	}

	config := models.ExternalDb{Username: u, Password: p, Name: d, Host: h, Port: dPort, Type: dType, SslMode: m, CreatedBy: user.Id, Owner: o, Protocol: protocol, HostFallback: hostFallback, PortFallback: portFallback, ProtocolFallback: protocolFallback}

	id, err := s.CreateExternalDb(config)

	if err != nil {
		components.Response(models.Response{Message: err.Error(), IsSuccess: false}).Render(r.Context(), w)
		return
	}
	go s.CreateAdminLog(models.AdminLog{UserId: user.Id, Action: models.CreateAdminAction, ResourceType: models.DbConnection, ResourceId: id, ResourceName: d})

	w.Header().Add("HX-Redirect", "/settings/dbs")
	w.Write([]byte("DB Connection config has been created"))
}

func GetEditExternalDbConfigPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	i := chi.URLParam(r, "id")
	d, err := s.GetDbById(i)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	err = externalDb.UpdateDbConfigPage(d).Render(r.Context(), w)
	if err != nil {
		fmt.Println("something went wrong when rendering edit db config page")
	}
}

func UpdateExternalDbHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	user := r.Context().Value(models.UserCtx).(*models.AppUser)

	i := chi.URLParam(r, "id")
	n := r.FormValue("name")
	h := r.FormValue("host")
	p := r.FormValue("port")
	t := r.FormValue("type")
	m := r.FormValue("sslMode")
	o := r.FormValue("owner")
	protocol := r.FormValue("protocol")
	hostFallback := r.FormValue("host-fallback")
	pFallback := r.FormValue("port-fallback")
	protocolFallback := r.FormValue("protocol-fallback")

	port, err := strconv.Atoi(p)
	if err != nil {
		components.Response(models.CreateResponse("Invalid port number", false)).Render(r.Context(), w)
		return
	}

	dbType, err := models.ToDbType(t)
	if err != nil {
		components.Response(models.CreateResponse("Invalid db type", false)).Render(r.Context(), w)
		return
	}

	portFallback, err := strconv.Atoi(pFallback)
	if err != nil {
		components.Response(models.CreateResponse("Invalid port number", false)).Render(r.Context(), w)
		return
	}

	err = s.UpdateExternalDb(models.ExternalDb{Id: i, Name: n, Host: h, Port: port, Type: dbType, SslMode: m, Owner: o, Protocol: protocol, HostFallback: hostFallback, PortFallback: portFallback, ProtocolFallback: protocolFallback})
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	go s.CreateAdminLog(models.AdminLog{UserId: user.Id, Action: models.UpdateSettingsAdminAction, ResourceType: models.DbConnection, ResourceId: i, ResourceName: n})

	w.Header().Add("HX-Redirect", "/settings/dbs")
	w.Write([]byte("db config updated"))
}

func GetUpdateExternalDbCredPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	i := chi.URLParam(r, "id")

	d, err := s.GetDbById(i)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Db config not found"))
		return
	}

	err = externalDb.UpdateCredentials(d.Name, d.Id).Render(r.Context(), w)
	if err != nil {
		fmt.Println("Something went wrong when rendering update external db credentials page")
	}
}

func UpdateExternalDbCredHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	user := r.Context().Value(models.UserCtx).(*models.AppUser)

	i := chi.URLParam(r, "id")
	u := r.FormValue("username")
	p := r.FormValue("password")
	p2 := r.FormValue("re-password")

	if p != p2 {
		components.Response(models.CreateResponse("passwords do not match", false)).Render(r.Context(), w)
		return
	}

	err := s.UpdateExternalDbCredentials(i, u, p)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}
	go func() {
		db, err := s.GetDbById(i)
		if err != nil {
			return
		}

		s.CreateAdminLog(models.AdminLog{UserId: user.Id, Action: models.UpdateCredentialsAdminAction, ResourceType: models.DbConnection, ResourceId: i, ResourceName: db.Name})
	}()

	w.Header().Add("HX-Redirect", "/settings/dbs")
	w.Write([]byte("db credential updated"))
}

func DeleteExternalDbByIdHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	user := r.Context().Value(models.UserCtx).(*models.AppUser)
	id := chi.URLParam(r, "id")

	go func() {
		db, err := s.GetDbById(id)
		if err != nil {
			return
		}
		s.CreateAdminLog(models.AdminLog{UserId: user.Id, Action: models.DeleteAdminAction, ResourceType: models.DbConnection, ResourceId: id, ResourceName: db.Name})
	}()

	err := s.DeleteUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
