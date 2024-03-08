package main

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/server"
	"custom-db-platform/src/views"
	"log"
)

func init() {
	db.ConnectDB()
	views.LoadTemplates()
}

func main() {
	webServer := server.NewServer()
	err := webServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
