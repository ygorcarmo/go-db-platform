package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ygorcarmo/db-platform/handlers"
	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/storage"
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

	router.Handle("/*", public())

	router.Get("/login", handlers.GetLoginPage)
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) { handlers.HandleLogin(w, r, s.store) })
	router.Get("/logout", handlers.Logout)

	router.Group(func(r chi.Router) {
		r.Use(s.authentication)
		r.Get("/", handlers.GetHomePage)
		// THIS is only for DEV
		r.Get("/seed", func(w http.ResponseWriter, r *http.Request) { handlers.SeedHandler(w, r, s.store) })

		r.Route("/db", func(dbroute chi.Router) {
			dbroute.Get("/create-user", func(w http.ResponseWriter, r *http.Request) { handlers.GetCreateDbUserPage(w, r, s.store) })
			dbroute.Post("/create-user", func(w http.ResponseWriter, r *http.Request) {
				handlers.ExternalDBUserHandler(w, r, s.store, models.Create)
			})

			dbroute.Get("/delete-user", func(w http.ResponseWriter, r *http.Request) { handlers.GetDeleteDbUserPage(w, r, s.store) })
			dbroute.Post("/delete-user", func(w http.ResponseWriter, r *http.Request) {
				handlers.ExternalDBUserHandler(w, r, s.store, models.Delete)
			})

		})

		r.Route("/user", func(userRoute chi.Router) {
			userRoute.Get("/reset-password", handlers.GetResetPasswordPage)
		})

		r.Group(func(adminR chi.Router) {
			adminR.Use(s.adminsOnly)

			adminR.Route("/settings", func(settingsR chi.Router) {
				settingsR.Get("/", handlers.GetSettingsPage)
			})
		})

	})

	// router.Use(authentication)

	// router.Get("/test", func(w http.ResponseWriter, r *http.Request) { handlers.HandleHome(w, r, s.store) })
	slog.Info("Server is running on: ", "listenAddr", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}