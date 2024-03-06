package models

import (
	"custom-db-platform/src/db"
	"database/sql"
	"fmt"
	"sync"
)

type TargetDb struct {
	Host    string
	Port    int
	Type    string
	SslMode string
	Name    string
}

func (targetDb *TargetDb) GetByName(name string) (*TargetDb, error) {
	err := db.Database.QueryRow("SELECT * FROM external_databases WHERE name=?;", name).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode)
	return targetDb, err
}

func (targetDb *TargetDb) GetAllNames() ([]string, error) {
	var dbs []string

	rows, err := db.Database.Query("SELECT name FROM external_databases")

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

func (targetDb TargetDb) ConnectToDBAndCreateUser(newUser string, c chan TargetDbsRepose, wg *sync.WaitGroup) {
	defer wg.Done()
	var connectionStr string
	if targetDb.Type == "postgres" {
		connectionStr = fmt.Sprintf("postgres://postgres:test@%s:%d/?sslmode=%s", targetDb.Host, targetDb.Port, targetDb.SslMode)
	} else if targetDb.Type == "mysql" {
		connectionStr = fmt.Sprintf("root:test@tcp(%s:%d)/", targetDb.Host, targetDb.Port)
	} else {
		c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: DB Type not Supported", newUser, targetDb.Name), Success: false}
		return
	}

	localdb, err := sql.Open(targetDb.Type, connectionStr)
	if err != nil {
		// log.Fatal(err)
		c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
		return
	}

	defer localdb.Close()

	pingErr := localdb.Ping()
	if pingErr != nil {
		// log.Fatal(pingErr)
		c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}

		return
	}
	fmt.Println("Connected!")

	if targetDb.Type == "postgres" {

		_, err := localdb.Exec(fmt.Sprintf("CREATE USER %s PASSWORD '1234';", newUser))
		if err != nil {
			// log.Fatal(err)
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}

			return
		}
	}

	if targetDb.Type == "mysql" {

		_, err := localdb.Exec(fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY 'password';", newUser))
		if err != nil {
			// log.Fatal(err)
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}

			return
		}
	}

	c <- TargetDbsRepose{Message: fmt.Sprintf("User %s has been created successfully at %s \n", newUser, targetDb.Name), Success: true}
}
