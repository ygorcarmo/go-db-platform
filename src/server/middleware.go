package server

import (
	"context"
	"custom-db-platform/src/models"
	"custom-db-platform/src/utils"
	"fmt"
	"net/http"
)

// TODO: Make a custom middleware
func verifyUserMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
}

func adminsOnlyMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value("userId").(string)

			var user models.AppUser
			user.GetUserById(userId)
			fmt.Println(user.IsAdmin)
			if !user.IsAdmin {
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
