package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var Database *sql.DB

type dbDetails struct {
	Name    string
	Host    string
	Port    int
	DbType  string
	SslMode string
}

func ConnectDB() {

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDRESS"),
		DBName: os.Getenv("DB_NAME"),
	}
	// Get a database handle.
	var err error
	Database, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Database.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}
