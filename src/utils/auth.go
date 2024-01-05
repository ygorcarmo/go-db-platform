package utils

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// TODO: put clerk functionality here
func CustomRequireSessionV2(client clerk.Client, verifyTokenOptions ...clerk.VerifyTokenOption) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
			if !ok || claims == nil {
				signInTemplate := filepath.Join("src", "web", "signIn.tmpl")
				baseTemplate := filepath.Join("src", "web", "base.tmpl")

				tmpl, err := template.ParseFiles(signInTemplate, baseTemplate)

				if err != nil {
					log.Fatal(err)
				}
				err = tmpl.Execute(w, nil)
				return
			}

			next.ServeHTTP(w, r)
		})

		return clerk.WithSessionV2(client, verifyTokenOptions...)(f)
	}
}
