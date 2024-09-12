package storage

import (
	"database/sql"
	"errors"
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

	dbConnections := []models.ExternalDb{
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

func (db *MySQLStorage) UpdateApplicationUserPassword(i string, p string) error {
	_, err := db.connection.Exec(`
	 	UPDATE users SET password = ? WHERE id=UUID_TO_BIN(?);
	`, p, i)
	return err
}

func (db *MySQLStorage) UpdateApplicationUser(u models.AppUser) error {
	_, err := db.connection.Exec(`
	 	UPDATE users 
		SET username=?, supervisor=?, sector=?, isAdmin=?
		WHERE id=UUID_TO_BIN(?);
	`, u.Username, u.Supervisor, u.Sector, u.IsAdmin, u.Id)
	return err
}

func (db *MySQLStorage) UpdateApplicationUserCredentials(u string, p string, i string) error {
	_, err := db.connection.Exec(`
		UPDATE users
		SET username=?, password=?
		WHERE id=UUID_TO_BIN(?);
	`, u, p, i)
	return err
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

func (db *MySQLStorage) CreateApplicationUser(u models.AppUser) (string, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return "", err
	}

	// Insert user within the transaction
	_, err = tx.Exec(`
        INSERT INTO users 
        (username, password, supervisor, sector, isAdmin)
        VALUES (?, ?, ?, ?, ?);`,
		u.Username, u.Password, u.Supervisor, u.Sector, u.IsAdmin)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Retrieve the last inserted ID in the same transaction
	err = tx.QueryRow(`SELECT BIN_TO_UUID(id) FROM users WHERE id = LAST_INSERT_ID();`).Scan(&u.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Commit the transaction
	err = tx.Commit()

	return u.Id, err
}

func (db *MySQLStorage) GetDbById(i string) (*models.ExternalDb, error) {
	r := models.ExternalDb{Id: i}
	err := db.connection.QueryRow("SELECT name, host, port, type, sslMode, owner from external_databases WHERE id = UUID_TO_BIN(?);",
		i).Scan(&r.Name, &r.Host, &r.Port, &r.Type, &r.SslMode, &r.Owner)
	if err != nil {
		return nil, err
	}
	return &r, nil
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

func (db *MySQLStorage) GetAllDbs() ([]models.ExternalDb, error) {
	databasesConfig := []models.ExternalDb{}

	rows, err := db.connection.Query(`
		SELECT 
    		BIN_TO_UUID(ed.id) AS external_database_id,
    		ed.name AS external_database_name,
    		ed.host AS external_database_host,
    		ed.port AS external_database_port,
    		ed.type AS external_database_type,
    		ed.sslMode AS external_database_sslMode,
    		u.username AS user_username,
			ed.owner
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
		config := models.ExternalDb{}
		if err := rows.Scan(&config.Id, &config.Name, &config.Host, &config.Port, &config.Type, &config.SslMode, &config.CreatedBy, &config.Owner); err != nil {
			return nil, err
		}
		databasesConfig = append(databasesConfig, config)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return databasesConfig, nil
}

func (db *MySQLStorage) GetDbByName(name string) (*models.ExternalDb, error) {
	targetDb := models.ExternalDb{Name: name}
	err := db.connection.QueryRow(`
	SELECT 
	    ed.name AS external_database_name,
    	ed.host AS external_database_host,
    	ed.port AS external_database_port,
    	ed.type AS external_database_type,
    	ed.sslMode AS external_database_sslMode,
		BIN_TO_UUID(ed.id) as external_database_id,
		ed.username,
		ed.password,
		ed.owner
	FROM 
    	external_databases ed 
	WHERE 
		name=?;`, name).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode, &targetDb.Id, &targetDb.Username, &targetDb.Password, &targetDb.Owner)
	if err != nil {
		return nil, err
	}

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

	return &targetDb, nil
}

func (db *MySQLStorage) DeleteUserById(id string) error {
	_, err := db.connection.Exec("DELETE FROM users WHERE id=UUID_TO_BIN(?);", id)
	return err
}

func (db *MySQLStorage) UpdateExternalDb(e models.ExternalDb) error {
	_, err := db.connection.Exec(`
		UPDATE external_databases
		SET name=?, host=?, port=?, type=?, sslMode=?, owner=?
		WHERE id = UUID_TO_BIN(?);
	`, e.Name, e.Host, e.Port, e.Type, e.SslMode, e.Owner, e.Id)
	return err
}

func (db *MySQLStorage) UpdateExternalDbCredentials(i string, u string, p string) error {
	eService, err := utils.NewEncryptionService()
	if err != nil {
		return err
	}

	username, err := eService.Encrypt(u)
	if err != nil {
		return err
	}

	password, err := eService.Encrypt(p)
	if err != nil {
		return err
	}

	_, err = db.connection.Exec(`
		UPDATE external_databases
		SET username=?, password=?
		WHERE id = UUID_TO_BIN(?);
	`, username, password, i)
	return err
}

func (db *MySQLStorage) CreateExternalDb(edb models.ExternalDb) (string, error) {
	eService, err := utils.NewEncryptionService()
	if err != nil {
		return "", err
	}

	ePassword, err := eService.Encrypt(edb.Password)
	if err != nil {
		return "", err
	}
	edb.Password = ePassword

	eUsername, err := eService.Encrypt(edb.Username)
	if err != nil {
		return "", err
	}
	edb.Username = eUsername

	tx, err := db.connection.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`
	INSERT INTO external_databases
		(name, host, port, type, sslMode, username, password, createdBy, owner)
	VALUES 
		(?, ?, ?, ?, ?, ?, ?, UUID_TO_BIN(?), ?);`,
		edb.Name, edb.Host, edb.Port, edb.Type, edb.SslMode, edb.Username, edb.Password, edb.CreatedBy, edb.Owner)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.QueryRow(`SELECT BIN_TO_UUID(id) FROM users WHERE id = LAST_INSERT_ID();`).Scan(&edb.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	err = tx.Commit()

	return edb.Id, err
}

func (db *MySQLStorage) DeleteExternalDbById(id string) error {
	_, err := db.connection.Exec("DELETE FROM external_databases WHERE id=UUID_TO_BIN(?)", id)
	return err
}

func (db *MySQLStorage) CreateLog(log models.Log) error {
	_, err := db.connection.Exec("INSERT INTO logs (dbId, newUser, wo, createdBy, action, success) VALUES (UUID_TO_BIN(?),?,?,UUID_TO_BIN(?), ?, ?);", log.DbId, log.NewUser, log.WO, log.CreateBy, log.Action, log.Success)
	return err
}

func (db *MySQLStorage) GetAllLogs() ([]models.LogResponse, error) {
	logs := []models.LogResponse{}
	rows, err := db.connection.Query(`
	SELECT 
		l.newUser, 
		ed.name AS external_database_name, 
		l.wo, 
		u.username,
		l.createdAt,
		l.action,
		l.success
	FROM 
		logs l
	JOIN 
		users u ON l.createdBy = u.id 
	JOIN 
		external_databases ed ON l.dbId = ed.id 
	ORDER BY 
		l.createdAt DESC;`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		log := models.LogResponse{}
		err := rows.Scan(&log.NewUser, &log.Database, &log.WO, &log.CreatedBy, &log.CreatedAt, &log.Action, &log.Success)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (db *MySQLStorage) CreateAdminLog(l models.AdminLog) error {
	_, err := db.connection.Exec(`
	INSERT INTO admin_logs
		(action, resourceId, resourceType, userId)
	VALUES
		(?, UUID_TO_BIN(?), ?,UUID_TO_BIN(?));`, l.Action, l.ResourceId, l.ResourceType, l.UserId)
	return err
}

func (db *MySQLStorage) GetAllAdminLogs() ([]models.AdminLogResponse, error) {
	ls := []models.AdminLogResponse{}
	rows, err := db.connection.Query(`
		SELECT
			l.action,
			u.username,
			BIN_TO_UUID(l.resourceId),
			l.resourceType,
			l.createdAt
		FROM
			admin_logs l
		JOIN
			users u ON l.userId = u.id
		ORDER BY
			l.createdAt DESC;
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		l := models.AdminLogResponse{}
		err := rows.Scan(&l.Action, &l.Username, &l.ResourceId, &l.ResourceType, &l.CreatedAt)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ls, nil
}
