package models

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	go_ora "github.com/sijms/go-ora/v2"
)

type dbType string

const (
	postgres dbType = "postgres"
	mySQL    dbType = "mysql"
	oracle   dbType = "oracle"
	// This passwd should come from the user
	passwd string = "THISISASECRET"
)

type TargetDb struct {
	Id        string
	Name      string
	Host      string
	Port      int
	Type      dbType
	SslMode   string
	Username  string
	Password  string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewDbUserProps struct {
	Username      string
	CurrentUserId string
	WO            int
}
type TargetDbsResponse struct {
	Message string
	Success bool
	DbId    string
}

func (t *TargetDb) ConnectAndCreateUser(user NewDbUserProps, c chan TargetDbsResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	switch t.Type {
	case postgres:
		pg, err := t.connectToPostgresql()

		if err != nil {
			c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err), Success: false, DbId: t.Id}
			return
		}
		defer pg.Close()

		_, err = pg.Exec("CREATE USER " + user.Username + " WITH PASSWORD '" + passwd + "';")

		if err != nil {
			c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err), Success: false, DbId: t.Id}
			return
		}

	case mySQL:
		mysql, err := t.connectToSQL()

		if err != nil {
			c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err), Success: false, DbId: t.Id}
			return
		}
		defer mysql.Close()

		_, err = mysql.Exec("CREATE USER '" + user.Username + "'@'localhost' IDENTIFIED BY '" + passwd + "';")
		if err != nil {
			c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err), Success: false, DbId: t.Id}
			return
		}

	case oracle:
		db, err := t.connectToOracle()
		if err != nil {
			c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err), Success: false, DbId: t.Id}
			return
		}
		defer db.Close()

		_, err = db.Exec("CREATE USER " + user.Username + " IDENTIFIED BY " + passwd)
		if err != nil {
			c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err), Success: false, DbId: t.Id}
			return
		}

	default:
		c <- TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, fmt.Errorf("DB type %s not supported", t.Type)), Success: false, DbId: t.Id}
		return
	}

	// log := Log{
	// 	DbId:     t.Id,
	// 	NewUser:  user.Username,
	// 	WO:       user.WO,
	// 	CreateBy: user.CurrentUserId,
	// 	Action:   create,
	// }

	// create a log when the thing gets returned
	// go s.CreateLog(log)

	c <- TargetDbsResponse{Message: fmt.Sprintf("User %s has been created successfully at %s \n", user.Username, t.Name), Success: true, DbId: t.Id}

}

func (targetDb *TargetDb) connectToPostgresql() (*sql.DB, error) {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=%s", targetDb.Username, targetDb.Password, targetDb.Host, targetDb.Port, targetDb.SslMode)

	database, err := sql.Open(string(targetDb.Type), connectionStr)
	if err != nil {
		return nil, err
	}
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connected to %s!\n", targetDb.Name)
	return database, nil
}

func (targetDb *TargetDb) connectToSQL() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 targetDb.Username,
		Passwd:               targetDb.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", targetDb.Host, targetDb.Port),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	database, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connected to %s!\n", targetDb.Name)

	return database, nil
}

func (targetDb *TargetDb) connectToOracle() (*sql.DB, error) {
	connectionStr := go_ora.BuildUrl(targetDb.Host, targetDb.Port, targetDb.Name, targetDb.Username, targetDb.Password, nil)
	database, err := sql.Open(string(targetDb.Type), connectionStr)
	if err != nil {
		return nil, err
	}
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connected to %s!\n", targetDb.Name)
	return database, nil
}

func makeErrorMessage(username string, dbName string, err error) string {
	return fmt.Sprintf("Error when adding %s at %s: %v", username, dbName, err)
}
