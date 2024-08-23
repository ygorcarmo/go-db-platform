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

func (db *MySQLStorage) Seed() {
	var userId []byte

	// Get the user ID from the database
	err := db.connection.QueryRow(`
        SELECT id
        FROM users
        WHERE username = "admin"
    `).Scan(&userId)
	if err != nil {
		fmt.Println("Failed to get userId: ", err)
		return
	}

	// Insert into external_databases table
	_, err = db.connection.Exec(`
        INSERT INTO
            external_databases (name, host, port, type, sslMode, userId)
        VALUES
            ("mysql","dbsql",3306,"mysql","disable",?),
            ("mysql-2","db-sql-02",3306,"mysql","disable",?),
            ("mysql-3","db-sql-03",3306,"mysql","disable",?),
            ("maria","db-maria",3306,"mysql","disable",?),
            ("maria-2","db-maria-02",3306,"mysql","disable",?),
            ("postgres","postgres",5432,"postgres","disable",?),
            ("XEPDB1","172.21.192.1",1521,"oracle","disable",?),
            ("postgres-2","pg",5432,"postgres","disable",?)
    `, userId, userId, userId, userId, userId, userId, userId, userId)
	if err != nil {
		fmt.Println("Failed to seed: ", err)
	}
}

func (db *MySQLStorage) GetUserById(id string) (*models.AppUser, error) {
	user := models.AppUser{Id: id}
	err := db.connection.QueryRow("SELECT username, password, isAdmin, supervisor, sector FROM users WHERE id=UUID_TO_BIN(?);", id).Scan(&user.Username, &user.Password, &user.IsAdmin, &user.Supervisor, &user.Sector)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *MySQLStorage) GetUserByUsername(username string) (*models.AppUser, error) {
	user := models.AppUser{Username: username}
	err := db.connection.QueryRow("SELECT BIN_TO_UUID(id), username, password FROM users WHERE username=?;", username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *MySQLStorage) GetAvailableDbs() ([]string, error) {
	var names []string

	rows, err := db.connection.Query("SELECT name FROM external_databases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		names = append(names, dbName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}
