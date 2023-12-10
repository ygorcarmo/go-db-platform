package main

import (
	"custom-db-platform/internal/server"
	"custom-db-platform/internal/utils"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	utils.ConnectToDB()
}

func main() {
	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}
}
