package main

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/server"
	"custom-db-platform/src/web"
	"fmt"
)

func init() {
	db.ConnectDB()
	web.LoadTemplates()
}

func main() {
	webServer := server.NewServer()
	err := webServer.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}

}
