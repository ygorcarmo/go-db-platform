package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/ygorcarmo/db-platform/models"
	"github.com/ygorcarmo/db-platform/utils"
)

type MySQLStorage struct {
	connection *sql.DB
}

// schCreateExternalDb implements Storage.
func (db *MySQLStorage) schCreateExternalDb(models.TargetDb) error {
	panic("unimplemented")
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
	var userId string

	// Get the user ID from the database
	err := db.connection.QueryRow(`
        SELECT BIN_TO_UUID(id)
        FROM users
        WHERE username = "admin"
    `).Scan(&userId)
	if err != nil {
		fmt.Println("Failed to get userId: ", err)
		return err
	}

	dbConnections := []models.TargetDb{
		{Name: "mysql", Host: "localhost", Port: 3001, Type: models.MySQL, SslMode: "disable", Username: "apt_db_platform", Password: "1qaz!EDC", CreatedBy: userId},
		{Name: "mysql-2", Host: "db-sql-02", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
		{Name: "mysql-3", Host: "db-sql-03", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
		{Name: "maria", Host: "db-maria", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
		{Name: "maria-2", Host: "db-maria-02", Port: 3306, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
		{Name: "postgres", Host: "postgres", Port: 5432, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
		{Name: "postgres-2", Host: "postgres-02", Port: 5432, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
		{Name: "XEPDB1", Host: "172.21.192.1", Port: 1521, Type: models.MySQL, SslMode: "disable", Username: "user", Password: "CHANGEME", CreatedBy: userId},
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
		query += "(?, ?, ?, ?, ?, UUID_TO_BIN(?), ?, ?),"
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

func (db *MySQLStorage) GetAllUsers() ([]models.AppUser, error) {
	users := []models.AppUser{}
	rows, err := db.connection.Query("SELECT BIN_TO_UUID(id), username, password, supervisor, sector, isAdmin, createdAt, updatedAt FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.AppUser{}
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Supervisor, &user.Sector, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *MySQLStorage) CreateApplicationUser(u models.AppUser) error {
	_, err := db.connection.Exec("INSERT INTO users (username, password, supervisor, sector, isAdmin) VALUES (?, ?, ?, ?, ?);",
		u.Username, u.Password, u.Supervisor, u.Sector, u.IsAdmin)

	return err
}

func (db *MySQLStorage) GetDbsName() ([]string, error) {
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

func (db *MySQLStorage) GetAllDbs() ([]models.TargetDb, error) {
	databasesConfig := []models.TargetDb{}

	rows, err := db.connection.Query(`
		SELECT 
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
    		users u ON ed.createdBy = u.id;
	`)

	if err != nil {
		return databasesConfig, err
	}
	defer rows.Close()

	for rows.Next() {
		config := models.TargetDb{}
		if err := rows.Scan(&config.Id, &config.Name, &config.Host, &config.Port, &config.Type, &config.SslMode, &config.CreatedBy); err != nil {
			return nil, err
		}
		databasesConfig = append(databasesConfig, config)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return databasesConfig, nil
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
		BIN_TO_UUID(ed.id) as external_database_id,
		ed.username,
		ed.password
	FROM 
    	external_databases ed 
	WHERE 
		name=?;`, name).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode, &targetDb.Id, &targetDb.Username, &targetDb.Password)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Username: %s, Password: %s\n", targetDb.Username, targetDb.Password)

	eService, err := utils.NewEncryptionService()
	if err != nil {
		return nil, errors.New("unable to create encryption key")
	}

	password, err := eService.Decrypt(targetDb.Password)
	if err != nil {
		return nil, errors.New("unable to decrypt password")
	}
	targetDb.Password = password

	username, err := eService.Decrypt(targetDb.Username)
	if err != nil {
		return nil, errors.New("unable to decrypt username")
	}
	targetDb.Username = username
	fmt.Printf("Username: %s, Password: %s\n", targetDb.Username, targetDb.Password)

	return &targetDb, nil
}
func (db *MySQLStorage) CreateExternalDb(edb models.TargetDb) error {
	eService, err := utils.NewEncryptionService()
	if err != nil {
		return err
	}

	ePassword, err := eService.Encrypt(edb.Password)
	if err != nil {
		return err
	}
	edb.Password = ePassword

	eUsername, err := eService.Encrypt(edb.Username)
	if err != nil {
		return err
	}
	edb.Username = eUsername

	_, err = db.connection.Exec("INSERT INTO external_databases (name, host, port, type, sslMode, username, password, createdBy) VALUES (?, ?, ?, ?, ?, ?, ?, UUID_TO_BIN(?));",
		edb.Name, edb.Host, edb.Port, edb.Type, edb.SslMode, edb.Username, edb.Password, edb.CreatedBy)
	return err
}

func (db *MySQLStorage) CreateLog(log models.Log) error {
	_, err := db.connection.Exec("INSERT INTO logs (dbId, newUser, wo, createdBy, action, success) VALUES (UUID_TO_BIN(?),?,?,UUID_TO_BIN(?), ?, ?);", log.DbId, log.NewUser, log.WO, log.CreateBy, log.Action, log.Sucess)
	return err
}
