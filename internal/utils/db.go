package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// TODO implement this later
type DBConfig struct {
	Username string
	Password string
	Address  string
	Type     string
	SSLMode  string
}

func ConnectToDBAndCreateUser(username string, password string, address string, dbType string, sslMode string, newUser string) {

	// username:password@protocol(address)/dbname?param=value
	// dsn := "%s:%s@%s/?sslmode=%s"
	var connectionStr string
	if dbType == "postgres" {
		connectionStr = fmt.Sprintf("postgres://%s:%s@%s/?sslmode=%s", username, password, address, sslMode)
	}
	if dbType == "mysql" {
		connectionStr = fmt.Sprintf("%s:%s@tcp(%s)/", username, password, address)
	}
	fmt.Println(connectionStr)

	db, err := sql.Open(dbType, connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Connect to docker postgres
	// Capture connection properties.
	// cfg := mysql.Config{
	// 	User:   "newuser",
	// 	Passwd: "123",
	// 	Net:    "tcp",
	// 	Addr:   "localhost:3306",
	// 	DBName: "test",
	// }
	// // Get a database handle.
	// var err error
	// db, err = sql.Open("mysql", cfg.FormatDSN())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// connStr := "postgres://postgres:test@localhost:5432/postgres?sslmode=disable"
	// db, err := sql.Open("postgres", connStr)

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	if dbType == "postgres" {

		_, err := db.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '1234'", newUser))
		if err != nil {
			log.Fatal(err)
		}
	}

	if dbType == "mysql" {

		_, err := db.Exec(fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY 'password';", newUser))
		if err != nil {
			log.Fatal(err)
		}
	}

}
