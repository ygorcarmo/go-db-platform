package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ygorcarmo/db-platform/storage"
	"github.com/ygorcarmo/db-platform/views/login"
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

	router.Get("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login.Index().Render(r.Context(), w)
	}))

	// router.Use(authentication)

	router.Get("/test", make(s.handleHome))
	slog.Info("Server is running on: ", "listenAddr", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}
