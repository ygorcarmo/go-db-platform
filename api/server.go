package api

import (
	_ "embed"
	"log/slog"
	"net/http"
	"time"

	"db-platform/handlers"
	"db-platform/models"
	"db-platform/storage"

	"github.com/go-chi/chi"
)

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	router := chi.NewMux()
	router.Use(s.addHttpHeaders)

	router.Handle("/*", public())

	router.Get("/login", handlers.GetLoginPage)
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) { handlers.HandleLogin(w, r, s.store) })
	router.Get("/logout", handlers.Logout)

	router.Group(func(r chi.Router) {
		r.Use(s.authentication)
		r.Get("/", handlers.GetHomePage)

		r.Route("/db", func(dbroute chi.Router) {
			dbroute.Get("/create-user", func(w http.ResponseWriter, r *http.Request) { handlers.GetCreateDbUserPage(w, r, s.store) })
			dbroute.Post("/create-user", func(w http.ResponseWriter, r *http.Request) {
				handlers.ExternalDBUserHandler(w, r, s.store, models.Create)
			})

			dbroute.Get("/delete-user", func(w http.ResponseWriter, r *http.Request) { handlers.GetDeleteDbUserPage(w, r, s.store) })
			dbroute.Post("/delete-user", func(w http.ResponseWriter, r *http.Request) {
				handlers.ExternalDBUserHandler(w, r, s.store, models.Delete)
			})

			dbroute.Get("/update-user", func(w http.ResponseWriter, r *http.Request) {
				handlers.GetUpdateDbUserPasswordPage(w, r, s.store)
			})
			dbroute.Post("/update-user", func(w http.ResponseWriter, r *http.Request) {
				handlers.ExternalDBUserHandler(w, r, s.store, models.UPDATEPWD)
			})

		})

		r.Route("/user", func(userRoute chi.Router) {
			userRoute.Get("/reset-password", handlers.GetResetPasswordPage)
			userRoute.Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {
				handlers.ResetApplicationUserPasswordHandler(w, r, s.store)
			})
		})

		r.Group(func(adminR chi.Router) {
			adminR.Use(s.adminsOnly)

			adminR.Route("/settings", func(settingsR chi.Router) {
				settingsR.Get("/", handlers.GetSettingsPage)

				settingsR.Route("/users", func(sur chi.Router) {
					sur.Get("/", func(w http.ResponseWriter, r *http.Request) { handlers.GetAllUserSettingsPage(w, r, s.store) })
					sur.Get("/create", handlers.GetCreateUserPage)
					sur.Post("/create", func(w http.ResponseWriter, r *http.Request) { handlers.CreateUserHandler(w, r, s.store) })
					sur.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.DeleteUserById(w, r, s.store) })
					sur.Get("/update/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.GetEditUserSettingsPage(w, r, s.store) })
					sur.Post("/update/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.UpdateUserSettingsHandler(w, r, s.store) })
					sur.Get("/{id}/credentials", func(w http.ResponseWriter, r *http.Request) { handlers.GetUpdateAppUserCredentials(w, r, s.store) })
					sur.Post("/{id}/credentials", func(w http.ResponseWriter, r *http.Request) { handlers.UpdateAppUserCredentialsHandler(w, r, s.store) })
					sur.Put("/{id}/unlock", func(w http.ResponseWriter, r *http.Request) { handlers.UnlockUser(w, r, s.store) })
				})

				settingsR.Route("/dbs", func(sdr chi.Router) {
					sdr.Get("/", func(w http.ResponseWriter, r *http.Request) { handlers.GetDatabasesConfigPage(w, r, s.store) })
					sdr.Get("/create", handlers.GetCreateExternalDbPage)
					sdr.Post("/create", func(w http.ResponseWriter, r *http.Request) { handlers.CreateExternalDbHandler(w, r, s.store) })
					sdr.Get("/edit/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.GetEditExternalDbConfigPage(w, r, s.store) })
					sdr.Put("/edit/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.UpdateExternalDbHandler(w, r, s.store) })
					sdr.Get("/{id}/credentials", func(w http.ResponseWriter, r *http.Request) { handlers.GetUpdateExternalDbCredPage(w, r, s.store) })
					sdr.Post("/{id}/credentials", func(w http.ResponseWriter, r *http.Request) { handlers.UpdateExternalDbCredHandler(w, r, s.store) })
					sdr.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.DeleteExternalDbByIdHandler(w, r, s.store) })
				})

				settingsR.Get("/logs", func(w http.ResponseWriter, r *http.Request) { handlers.GetLogsPage(w, r, s.store) })
				settingsR.Get("/admin-logs", func(w http.ResponseWriter, r *http.Request) { handlers.GetAdminLogsPage(w, r, s.store) })
			})
			// THIS is only for DEV
			adminR.Get("/seed", func(w http.ResponseWriter, r *http.Request) { handlers.SeedHandler(w, r, s.store) })

		})

	})

	slog.Info("Server is running on: ", "listenAddr", s.listenAddr)

	server := &http.Server{
		Addr:              s.listenAddr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router,
	}

	return server.ListenAndServe()
}
