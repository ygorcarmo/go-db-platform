package server

import (
	"net/http"
)

// TODO: Make a custom middleware
func CustomMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Replace clerk authentication logic
			// compare hashes, create JWT token, check if JWT token is valid
			// claims, ok := r.Context().Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
			// if !ok || claims == nil {
			// 	signInTemplate := filepath.Join("src", "web", "signIn.tmpl")
			// 	baseTemplate := filepath.Join("src", "web", "base.tmpl")

			// 	tmpl, err := template.ParseFiles(signInTemplate, baseTemplate)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	err = tmpl.Execute(w, nil)
			// 	return
			// }

			next.ServeHTTP(w, r)
			return
		})
	}
}
