//go:build dev
// +build dev

package api

import (
	"fmt"
	"net/http"
	"os"
)

func public() http.Handler {
	fmt.Println("building static files for development")
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("api/public")))
}
