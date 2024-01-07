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

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	if dbType == "postgres" {

		_, err := db.Exec("CREATE USER ? PASSWORD '1234'", newUser)
		if err != nil {
			log.Fatal(err)
		}
	}

	if dbType == "mysql" {

		_, err := db.Exec("CREATE USER '?'@'localhost' IDENTIFIED BY 'password';", newUser)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("User %s has been created successfully", newUser)

}
