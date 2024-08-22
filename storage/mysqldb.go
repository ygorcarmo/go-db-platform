package storage

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/ygorcarmo/db-platform/models"
)

type MySQLStorage struct {
	connection *sql.DB
}

func NewMySQLStorage(user string, password string, address string, dbName string) *MySQLStorage {
	log.Println("Creating db")
	cfg := mysql.Config{
		User:      user,
		Passwd:    password,
		Net:       "tcp",
		Addr:      address,
		DBName:    dbName,
		ParseTime: true,
	}
	fmt.Println(cfg)
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	slog.Info("Connected To DB")
	return &MySQLStorage{connection: db}
}

func (db *MySQLStorage) GetUserById(id string) (*models.AppUser, error) {
	fmt.Println("Getting User")
	user := models.AppUser{Id: id}
	err := db.connection.QueryRow("SELECT username, password, isAdmin, supervisor, sector FROM users WHERE id=UUID_TO_BIN(?);", id).Scan(&user.Username, &user.Password, &user.IsAdmin, &user.Supervisor, &user.Sector)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
