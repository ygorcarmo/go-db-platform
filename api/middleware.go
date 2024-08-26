package api

import (
	"context"
	"net/http"

	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/utils"
)

func (s *Server) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := r.Cookie("token")

		if err != nil || jwtToken.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		userId, err := utils.DecodeToken(jwtToken.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/"})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := s.store.GetUserById(userId)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), models.UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) adminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.UserCtx).(*models.AppUser)

		if !user.IsAdmin {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}
