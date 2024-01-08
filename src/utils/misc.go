package utils

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func ConnectToDBAndCreateUser(host string, port int, dbType string, sslMode string, newUser string, c chan string, wg *sync.WaitGroup) {
	var res string
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
		res = fmt.Sprintf("Error: %v", err)
		return
	}

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		// log.Fatal(pingErr)
		res = fmt.Sprintf("Error: %v", err)
		return
	}
	fmt.Println("Connected!")

	if dbType == "postgres" {

		_, err := db.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '1234';", newUser))
		if err != nil {
			// log.Fatal(err)
			res = fmt.Sprintf("Error: %v", err)
			return
		}
	}

	if dbType == "mysql" {

		_, err := db.Exec(fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY 'password';", newUser))
		if err != nil {
			// log.Fatal(err)
			res = fmt.Sprintf("Error: %v", err)
			return
		}
	}

	res = fmt.Sprintf("User %s has been created successfully\n", newUser)

	c <- res
	wg.Done()
}
