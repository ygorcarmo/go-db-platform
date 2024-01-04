package main

import (
	"custom-db-platform/internal/server"
	"fmt"
)

func init() {

	// connStr := "postgres://postgres:test@localhost:5432/postgres?sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// utils.ConnectToDBAndCreateUser("postgres", "test", "localhost:5432", "postgres", "disable", "test3")
	// utils.ConnectToDBAndCreateUser("root", "test", "localhost:3306", "mysql", "disable", "test3")
}

func main() {
	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}

}
