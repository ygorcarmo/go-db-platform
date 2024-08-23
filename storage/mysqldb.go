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

func (db *MySQLStorage) Seed() error {
	var userId []byte

	// Get the user ID from the database
	err := db.connection.QueryRow(`
        SELECT id
        FROM users
        WHERE username = "admin"
    `).Scan(&userId)
	if err != nil {
		fmt.Println("Failed to get userId: ", err)
		return err
	}

	// Insert into external_databases table
	_, err = db.connection.Exec(`
        INSERT INTO
            external_databases (name, host, port, type, sslMode, createdBy, username, password)
        VALUES
            ("mysql","dbsql",3306,"mysql","disable",?,"user", "CHANGEME"),
            ("mysql-2","db-sql-02",3306,"mysql","disable",?,"user", "CHANGEME"),
            ("mysql-3","db-sql-03",3306,"mysql","disable",?,"user", "CHANGEME"),
            ("maria","db-maria",3306,"mysql","disable",?,"user", "CHANGEME"),
            ("maria-2","db-maria-02",3306,"mysql","disable",?,"user", "CHANGEME"),
            ("postgres","postgres",5432,"postgres","disable",?,"user", "CHANGEME"),
            ("XEPDB1","172.21.192.1",1521,"oracle","disable",?,"user", "CHANGEME"),
            ("postgres-2","pg",5432,"postgres","disable",?,"user", "CHANGEME")
    `, userId, userId, userId, userId, userId, userId, userId, userId)

	if err != nil {
		fmt.Println("Failed to seed: ", err)
		return err
	}

	return nil
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

func (db *MySQLStorage) GetDbByName(name string) (*models.TargetDb, error) {
	targetDb := models.TargetDb{Name: name}
	err := db.connection.QueryRow(`
	SELECT 
	    ed.name AS external_database_name,
    	ed.host AS external_database_host,
    	ed.port AS external_database_port,
    	ed.type AS external_database_type,
    	ed.sslMode AS external_database_sslMode,
		BIN_TO_UUID(ed.id) as external_database_id
	FROM 
    	external_databases ed 
	WHERE 
		name=?;`, name).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode, &targetDb.Id)
	if err != nil {
		return nil, err
	}
	return &targetDb, nil
}

func (db *MySQLStorage) CreateLog(log models.Log) error {
	_, err := db.connection.Exec("INSERT INTO logs (dbId, newUser, wo, createdBy, action, success) VALUES (UUID_TO_BIN(?),?,?,UUID_TO_BIN(?), ?, ?);", log.DbId, log.NewUser, log.WO, log.CreateBy, log.Action, log.Sucess)
	return err
}
