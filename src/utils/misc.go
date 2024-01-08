package utils

import (
	"custom-db-platform/src/datatypes"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func ConnectToDBAndCreateUser(host string, port int, dbType string, sslMode string, newUser string, c chan datatypes.Result, wg *sync.WaitGroup) {
	var connectionStr string
	if dbType == "postgres" {
		connectionStr = fmt.Sprintf("postgres://postgres:test@%s:%d/?sslmode=%s", host, port, sslMode)
	}
	if dbType == "mysql" {
		connectionStr = fmt.Sprintf("root:test@tcp(%s:%d)/", host, port)
	}
	fmt.Println(connectionStr)

	db, err := sql.Open(dbType, connectionStr)
	if err != nil {
		// log.Fatal(err)
		c <- datatypes.Result{Message: fmt.Sprintf("Error: %v", err), Success: false}
	}

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		// log.Fatal(pingErr)
		c <- datatypes.Result{Message: fmt.Sprintf("Error: %v", err), Success: false}
	}
	fmt.Println("Connected!")

	if dbType == "postgres" {

		_, err := db.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '1234';", newUser))
		if err != nil {
			// log.Fatal(err)
			c <- datatypes.Result{Message: fmt.Sprintf("Error: %v", err), Success: false}
		}
	}

	if dbType == "mysql" {

		_, err := db.Exec(fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY 'password';", newUser))
		if err != nil {
			// log.Fatal(err)
			c <- datatypes.Result{Message: fmt.Sprintf("Error: %v", err), Success: false}
		}
	}
	c <- datatypes.Result{Message: fmt.Sprintf("User %s has been created successfully\n", newUser), Success: true}
	wg.Done()
}
