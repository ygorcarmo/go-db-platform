package models

import (
	"database/sql"
	"fmt"
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
	Message   string
	IsSuccess bool
	DbId      string
}

func (t *TargetDb) ConnectAndCreateUser(user *NewDbUserProps, r *[]TargetDbsResponse) {
	switch t.Type {
	case postgres:
		pg, err := t.connectToPostgresql()

		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbId: t.Id})
			return
		}
		defer pg.Close()

		_, err = pg.Exec("CREATE USER " + user.Username + " WITH PASSWORD '" + t.Password + "';")

		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbId: t.Id})
			return
		}

	case mySQL:
		mysql, err := t.connectToSQL()

		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbId: t.Id})
			return
		}
		defer mysql.Close()

		_, err = mysql.Exec("CREATE USER '" + user.Username + "'@'" + t.Host + "' IDENTIFIED BY '" + t.Password + "';")
		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbId: t.Id})
			return
		}

	case oracle:
		db, err := t.connectToOracle()
		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbId: t.Id})
			return
		}
		defer db.Close()

		_, err = db.Exec("CREATE USER " + user.Username + " IDENTIFIED BY " + t.Password)
		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbId: t.Id})
			return
		}

	default:
		*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, fmt.Errorf("DB type %s not supported", t.Type), Create), IsSuccess: false, DbId: t.Id})
		return
	}

	*r = append(*r, TargetDbsResponse{Message: fmt.Sprintf("User %s has been created successfully at %s \n", user.Username, t.Name), IsSuccess: true, DbId: t.Id})
}

func (t *TargetDb) ConnectAndDeleteUser(user *NewDbUserProps, r *[]TargetDbsResponse) {
	switch t.Type {
	case postgres:
		pg, err := t.connectToPostgresql()

		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbId: t.Id})
			return
		}
		defer pg.Close()

		_, err = pg.Exec("DROP USER IF EXISTS " + user.Username + ";")

		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbId: t.Id})
			return
		}

	case mySQL:
		mysql, err := t.connectToSQL()

		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbId: t.Id})
			return
		}
		defer mysql.Close()

		_, err = mysql.Exec("DROP USER IF EXISTS '" + user.Username + "'@'" + t.Host + "';")
		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbId: t.Id})
			return
		}

	case oracle:
		db, err := t.connectToOracle()
		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbId: t.Id})
			return
		}
		defer db.Close()

		_, err = db.Exec("DROP USER " + user.Username + " CASCADE ")
		if err != nil {
			*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbId: t.Id})
			return
		}

	default:
		*r = append(*r, TargetDbsResponse{Message: makeErrorMessage(user.Username, t.Name, fmt.Errorf("DB type %s not supported", t.Type), Delete), IsSuccess: false, DbId: t.Id})
		return
	}

	*r = append(*r, TargetDbsResponse{Message: fmt.Sprintf("User %s has been deleted successfully at %s \n", user.Username, t.Name), IsSuccess: true, DbId: t.Id})
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

func makeErrorMessage(username string, dbName string, err error, a ActionType) string {
	return fmt.Sprintf("Error with action %s. Username: %s at %s. Error: %v", a, username, dbName, err)
}
