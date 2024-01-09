package server

import (
	"database/sql"
	"fmt"
)

func ConnectToDBAndCreateUser(host string, port int, dbType string, sslMode string, newUser string, dbName string, c chan Result) {
	defer wg.Done()
	var connectionStr string
	if dbType == "postgres" {
		connectionStr = fmt.Sprintf("postgres://postgres:test@%s:%d/?sslmode=%s", host, port, sslMode)
	} else if dbType == "mysql" {
		connectionStr = fmt.Sprintf("root:test@tcp(%s:%d)/", host, port)
	} else {
		c <- Result{Message: fmt.Sprintf("Error when adding %s at %s: DB Type not Supported", newUser, dbName), Success: false}
		return
	}

	fmt.Println(connectionStr)

	localdb, err := sql.Open(dbType, connectionStr)
	if err != nil {
		// log.Fatal(err)
		c <- Result{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, dbName, err), Success: false}
		return
	}

	defer localdb.Close()

	pingErr := localdb.Ping()
	if pingErr != nil {
		// log.Fatal(pingErr)
		c <- Result{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, dbName, err), Success: false}

		return
	}
	fmt.Println("Connected!")

	if dbType == "postgres" {

		_, err := localdb.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '1234';", newUser))
		if err != nil {
			// log.Fatal(err)
			c <- Result{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, dbName, err), Success: false}

			return
		}
	}

	if dbType == "mysql" {

		_, err := localdb.Exec(fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY 'password';", newUser))
		if err != nil {
			// log.Fatal(err)
			c <- Result{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, dbName, err), Success: false}

			return
		}
	}

	c <- Result{Message: fmt.Sprintf("User %s has been created successfully at %s \n", newUser, dbName), Success: true}
}
