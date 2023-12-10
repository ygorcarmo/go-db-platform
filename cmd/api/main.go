package main

import (
	"custom-db-platform/internal/server"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
}

func main() {

	db, dberr := sql.Open("sqlite3", "./newDb.db")
	if dberr != nil {
		panic("Couldn't connect to db")
	}
	defer db.Close()
	createTable, tableerror := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTERGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if tableerror != nil {
		panic(tableerror)
	}
	createTable.Exec()

	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}
}
