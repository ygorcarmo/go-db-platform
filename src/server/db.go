package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

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
		User:   "root",
		Passwd: "test",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		Addr:   "127.0.0.1:3001",
		DBName: "db_platform",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func getDBsName() ([]string, error) {
	var dbs []string

	rows, err := db.Query("SELECT name FROM db_connection_info")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}

		dbs = append(dbs, dbName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dbs, nil

}

func getDBByName(name string) (*dbDetails, error) {

	var newDB dbDetails

	row := db.QueryRow("SELECT * FROM db_connection_info WHERE name=?;", name)

	if err := row.Scan(&newDB.Name, &newDB.Host, &newDB.Port, &newDB.DbType, &newDB.SslMode); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &newDB, nil
}
