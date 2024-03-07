package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var Database *sql.DB

func ConnectDB() {

	// Capture connection properties.
	// cfg := mysql.Config{
	// 	User:   os.Getenv("DB_USER"),
	// 	Passwd: os.Getenv("DB_PASSWORD"),
	// 	Net:    "tcp",
	// 	Addr:   os.Getenv("DB_ADDRESS"),
	// 	DBName: os.Getenv("DB_NAME"),
	// }
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "test",
		Net:       "tcp",
		Addr:      "127.0.0.1:3001",
		DBName:    "db_platform",
		ParseTime: true,
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
