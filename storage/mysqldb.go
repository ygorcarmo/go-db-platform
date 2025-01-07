package storage

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"db-platform/models"
	"db-platform/utils"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	connection *sql.DB
}

func NewMySQLStorage(user string, password string, address string, dbName string) *MySQLStorage {
	cfg := mysql.Config{
		User:      user,
		Passwd:    password,
		Net:       "tcp",
		Addr:      address,
		DBName:    dbName,
		ParseTime: true,
	}

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

func (db *MySQLStorage) Seed() error {

	dbConnections := []models.ExternalDb{
		{Name: "mysql", Host: "localhost", Port: 3001, Type: models.MySQL, SslMode: "disable", Username: "apt_db_platform", Password: "1qaz!EDC", CreatedBy: "admin"},
		{Name: "mysql-2", Host: "db-sql-02", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
		{Name: "mysql-3", Host: "db-sql-03", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
		{Name: "maria", Host: "db-maria", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
		{Name: "maria-2", Host: "db-maria-02", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
		{Name: "postgres", Host: "postgres", Port: 5432, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
		{Name: "postgres-2", Host: "postgres-02", Port: 5432, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
		{Name: "XEPDB1", Host: "172.21.192.1", Port: 1521, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: "admin"},
	}

	eService, err := utils.NewEncryptionService()
	if err != nil {
		log.Fatal("Error creating the encryption key: ", err)
	}

	// Prepare SQL query and arguments slice
	query := `
		INSERT INTO external_databases (name, host, port, type, sslMode, createdBy, username, password)
		VALUES
	`
	args := []interface{}{}

	for _, connection := range dbConnections {
		eu, err := eService.Encrypt(connection.Username)
		if err != nil {
			log.Fatal("Something went wrong when encrypting the username: ", err)
		}
		connection.Username = eu

		ep, err := eService.Encrypt(connection.Password)
		if err != nil {
			log.Fatal("Something went wrong when encrypting the password: ", err)
		}
		connection.Password = ep

		// Append each set of values to the query and args
		query += "(?, ?, ?, ?, ?, ?, ?, ?),"
		args = append(args, connection.Name, connection.Host, connection.Port, connection.Type, connection.SslMode, connection.CreatedBy, connection.Username, connection.Password)
	}

	// Remove the trailing comma
	query = query[:len(query)-1]

	// Execute the bulk insert
	_, err = db.connection.Exec(query, args...)
	if err != nil {
		fmt.Println("Failed to seed: ", err)
		return err
	}

	return nil
}
