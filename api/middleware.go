package api

import (
	"context"
	"log"
	"net/http"

	"github.com/ygorcarmo/db-platform/utils"
)

func authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := r.Cookie("token")

		if err != nil || jwtToken.Value == "" {
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		userId, err := utils.DecodeToken(jwtToken.Value)
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) adminsOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)

		// var user models.AppUser
		// user.GetUserById(userId)
		user, err := s.store.GetUserById(userId)
		if err != nil {
			log.Fatal("How Did this user get here without auth?")
			return
		}
		if !user.IsAdmin {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}
