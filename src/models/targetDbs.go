package models

import (
	"custom-db-platform/src/db"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	go_ora "github.com/sijms/go-ora/v2"
)

type TargetDb struct {
	Id        string
	Name      string
	Host      string
	Port      int
	Type      string
	SslMode   string
	UserId    string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TargetDbsRepose struct {
	Message string
	Success bool
}

func (targetDb *TargetDb) GetByName(name string) (*TargetDb, error) {
	err := db.Database.QueryRow(`SELECT 
	    ed.name AS external_database_name,
    ed.host AS external_database_host,
    ed.port AS external_database_port,
    ed.type AS external_database_type,
    ed.sslMode AS external_database_sslMode,
	BIN_TO_UUID(ed.id) as external_database_id
FROM 
    external_databases ed WHERE name=?;`, name).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode, &targetDb.Id)
	return targetDb, err
}

func (targetDb *TargetDb) GetByid(id string) (*TargetDb, error) {
	targetDb.Id = id
	err := db.Database.QueryRow(`SELECT name, host,	port, type,	sslMode FROM external_databases WHERE id=UUID_TO_BIN(?);`,
		id).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode)
	return targetDb, err
}

func (*TargetDb) GetAllNames() ([]string, error) {
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

func (targetDb *TargetDb) Update() error {
	_, err := db.Database.Exec(`
		UPDATE external_databases
		SET name = ?, host = ?, port = ?, type = ?, sslMode = ?
		WHERE id = UUID_TO_BIN(?);
	`, targetDb.Name, targetDb.Host, targetDb.Port, targetDb.Type, targetDb.SslMode, targetDb.Id)
	return err
}

func (*TargetDb) GetAll() ([]TargetDb, error) {
	var databases []TargetDb

	rows, err := db.Database.Query(`SELECT 
    BIN_TO_UUID(ed.id) AS external_database_id,
    ed.name AS external_database_name,
    ed.host AS external_database_host,
    ed.port AS external_database_port,
    ed.type AS external_database_type,
    ed.sslMode AS external_database_sslMode,
    u.username AS user_username
FROM 
    external_databases ed
JOIN 
    users u ON ed.userId = u.id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var database TargetDb
		if err := rows.Scan(&database.Id, &database.Name, &database.Host, &database.Port, &database.Type, &database.SslMode, &database.CreatedBy); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return databases, nil
}

func (*TargetDb) DeleteDbById(dbId string) error {
	_, err := db.Database.Exec("DELETE FROM external_databases WHERE id=UUID_TO_BIN(?)", dbId)
	return err
}

func (targetDb *TargetDb) ConnectToDBAndCreateUser(newUser, currentUserId string, wo int, c chan TargetDbsRepose, wg *sync.WaitGroup) {
	defer wg.Done()
	if targetDb.Type == "postgres" {
		c <- targetDb.connectToPostgreAndCreateUser(newUser, currentUserId, wo)
		return
	} else if targetDb.Type == "mysql" {
		c <- targetDb.connectToSQLAndCreateUser(newUser, currentUserId, wo)
		return
	} else if targetDb.Type == "oracle" {
		c <- targetDb.connectToOracleAndCreateUser(newUser, currentUserId, wo)
	} else {
		c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: DB Type not Supported", newUser, targetDb.Name), Success: false}
		return
	}
}

func (targetdb *TargetDb) connectToSQLAndCreateUser(newUser, currentUserId string, wo int) TargetDbsRepose {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "test",
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", targetdb.Host, targetdb.Port),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	database, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetdb.Name, err), Success: false}
	}
	err = database.Ping()
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetdb.Name, err), Success: false}
	}
	fmt.Println("Connected!")
	defer database.Close()

	_, err = database.Exec("CREATE USER '" + newUser + "'@'localhost' IDENTIFIED BY 'password';")
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetdb.Name, err), Success: false}
	}

	log := Log{
		DbId:    targetdb.Id,
		NewUser: newUser,
		WO:      wo,
		UserId:  currentUserId,
	}

	go log.CreateLog()

	return TargetDbsRepose{Message: fmt.Sprintf("User %s has been created successfully at %s \n", newUser, targetdb.Name), Success: true}
}

func (targetDb *TargetDb) connectToPostgreAndCreateUser(newUser, currentUserId string, wo int) TargetDbsRepose {
	connectionStr := fmt.Sprintf("postgres://postgres:test@%s:%d/?sslmode=%s", targetDb.Host, targetDb.Port, targetDb.SslMode)

	database, err := sql.Open(targetDb.Type, connectionStr)
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
	}
	err = database.Ping()
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
	}
	fmt.Println("Connected!")
	defer database.Close()
	_, err = database.Exec("CREATE USER " + newUser + " WITH PASSWORD '1234';")
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
	}

	log := Log{
		DbId:    targetDb.Id,
		NewUser: newUser,
		WO:      wo,
		UserId:  currentUserId,
	}

	go log.CreateLog()

	return TargetDbsRepose{Message: fmt.Sprintf("User %s has been created successfully at %s \n", newUser, targetDb.Name), Success: true}
}

func (targetDb *TargetDb) connectToOracleAndCreateUser(newUser, currentUserId string, wo int) TargetDbsRepose {
	connectionStr := go_ora.BuildUrl(targetDb.Host, targetDb.Port, targetDb.Name, "teste", "teste", nil)
	database, err := sql.Open(targetDb.Type, connectionStr)
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
	}
	err = database.Ping()
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
	}
	fmt.Println("Connected!")
	defer database.Close()

	_, err = database.Exec("CREATE USER " + newUser + " IDENTIFIED BY new_password")
	if err != nil {
		return TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUser, targetDb.Name, err), Success: false}
	}

	log := Log{
		DbId:    targetDb.Id,
		NewUser: newUser,
		WO:      wo,
		UserId:  currentUserId,
	}

	go log.CreateLog()

	return TargetDbsRepose{Message: fmt.Sprintf("User %s has been created successfully at %s \n", newUser, targetDb.Name), Success: true}
}
