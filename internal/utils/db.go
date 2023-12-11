package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectToDB() {
	// TODO: Connect to docker postgres
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "newuser",
		Passwd: "123",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "test",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	res, err := db.Exec("CREATE TABLE IF NOT EXISTS test (name varchar(25))")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
