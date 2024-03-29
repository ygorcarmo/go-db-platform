package server

import (
	"custom-db-platform/src/handlers"
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed assets/**
var assets embed.FS

func (s *Server) RegisterRoutes() http.Handler {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	fs := http.FileServer(http.FS(assets))
	router.Handle("/assets/*", fs)

	router.Get("/sign-in", handlers.LoadSignInPage)
	router.Post("/sign-in", handlers.HandleSignIn)

	router.Group(func(r chi.Router) {
		r.Use(verifyUserMiddleware())

		r.Get("/", handlers.LoadHomePage)

		r.Route("/user", func(userRoute chi.Router) {
			userRoute.Get("/reset-password", handlers.LoadResetPasswordPage)
			userRoute.Patch("/reset-password", handlers.ResetPasswordFormHandler)
			userRoute.Get("/logout", handlers.Logout)
		})

		r.Route("/db", func(r chi.Router) {
			r.Get("/create-user", handlers.LoadCreateExternalUser)
			r.Post("/create-user", handlers.CreateExternalUserFormHandler)

			r.Get("/delete-user", handlers.LoadDeleteExternalUser)
			r.Post("/delete-user", handlers.DeleteExternalUserFormHandler)
		})

		r.Group(func(adminsOnlyRoute chi.Router) {
			adminsOnlyRoute.Use(adminsOnlyMiddleware())

			adminsOnlyRoute.Route("/settings", func(settingsRoute chi.Router) {
				settingsRoute.Get("/", handlers.LoadSettings)
				settingsRoute.Get("/users", handlers.LoadManageUsers)
				settingsRoute.Delete("/user/{id}", handlers.DeleteAppUser)
				settingsRoute.Get("/create-user", handlers.LoadCreateAppUser)
				settingsRoute.Post("/create-user", handlers.AddAppUserFormHanlder)
				settingsRoute.Get("/update-user/{id}", handlers.LoadEditAppUser)
				settingsRoute.Put("/update-user/{id}", handlers.UpdateAppUser)

				settingsRoute.Get("/dbs", handlers.LoadManageDbs)
				settingsRoute.Delete("/db/{id}", handlers.DeleteDb)
				settingsRoute.Get("/create-db", handlers.LoadAddDb)
				settingsRoute.Post("/create-db", handlers.AddDbFormHanlder)
				settingsRoute.Get("/update-db/{id}", handlers.LoadEditDb)
				settingsRoute.Put("/update-db/{id}", handlers.UpdateDb)

				settingsRoute.Get("/logs", handlers.LoadManageLogs)
			})
		})
	})

	return router
}
