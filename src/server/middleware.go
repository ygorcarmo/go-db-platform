package server

import (
	"custom-db-platform/src/utils"
	"net/http"
)

// TODO: Make a custom middleware
func customMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtToken, err := r.Cookie("token")

			if err != nil {
				http.Redirect(w, r, "/sign-in", http.StatusFound)
				// w.WriteHeader(http.StatusUnauthorized)
				// fmt.Fprint(w, "token")
				return
			}

			err = utils.VerifyToken(jwtToken.Value)
			if err != nil {
				http.Redirect(w, r, "/sign-in", http.StatusFound)
				// w.WriteHeader(http.StatusUnauthorized)
				// fmt.Fprint(w, "Invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
