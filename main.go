package main

import (
	"custom-db-platform/src/server"
	"fmt"
)

func init() {
	server.ConnectDB()
}

func main() {
	webServer := server.NewServer()
	err := webServer.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}

}
