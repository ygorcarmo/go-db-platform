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
			if err != nil {
				http.Redirect(w, r, "/sign-in", http.StatusFound)
				// w.WriteHeader(http.StatusUnauthorized)
				// fmt.Fprint(w, "token")
				return
			}

			userId, err := utils.DecodeToken(jwtToken.Value)
			if err != nil {
				http.Redirect(w, r, "/sign-in", http.StatusFound)
				// w.WriteHeader(http.StatusUnauthorized)
				// fmt.Fprint(w, "token")
				return
			}
			fmt.Println(userId)

			ctx := context.WithValue(r.Context(), "userId", userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func adminsOnlyMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if user is admin
			userId := r.Context().Value("userId").(string)

			fmt.Println("admin")

			var user models.AppUser
			user.GetUserById(userId)
			fmt.Println(user)
			if !user.IsAdmin {
				http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
