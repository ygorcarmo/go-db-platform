package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToDB() {
	db, err := sql.Open("sqlite3", "./internal/data/newDb.db")
	if err != nil {
		panic("Couldn't connect to db")
	}
	defer db.Close()
	createTable, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTERGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	createTable.Exec()
}
