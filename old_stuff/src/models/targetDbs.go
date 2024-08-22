package models

import (
	"custom-db-platform/src/db"
	"database/sql"
	"fmt"
	"os"
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

type NewDbUserProps struct {
	Username      string
	CurrentUserId string
	WO            int
}

type externalDBUser struct {
	Username   string
	Password   string
	NewUserPWD string
}

var eDbUser externalDBUser

func init() {
	eDbUser = externalDBUser{
		Username:   os.Getenv("GLOBAL_DB_USER"),
		Password:   os.Getenv("GLOBAL_USER_PWD"),
		NewUserPWD: os.Getenv("NEW_USER_PWD"),
	}
}

func (targetDb *TargetDb) GetByName(name string) (*TargetDb, error) {
	err := db.Database.QueryRow(`
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

	rows, err := db.Database.Query(`
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

func (targetDb *TargetDb) ConnectToDBAndCreateUser(newUserProps NewDbUserProps, c chan TargetDbsRepose, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(eDbUser)
	if targetDb.Type == "postgres" {
		pg, err := targetDb.connectToPostgre()
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
		defer pg.Close()

		_, err = pg.Exec("CREATE USER " + newUserProps.Username + " WITH PASSWORD '" + eDbUser.NewUserPWD + "';")
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
	} else if targetDb.Type == "mysql" {
		database, err := targetDb.connectToSQL()
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
		defer database.Close()

		_, err = database.Exec("CREATE USER '" + newUserProps.Username + "'@'localhost' IDENTIFIED BY '" + eDbUser.NewUserPWD + "';")
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
	} else if targetDb.Type == "oracle" {
		oracledb, err := targetDb.connectToOracle()
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
		defer oracledb.Close()

		_, err = oracledb.Exec("CREATE USER " + newUserProps.Username + " IDENTIFIED BY " + eDbUser.NewUserPWD)
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
	} else {
		c <- TargetDbsRepose{Message: fmt.Sprintf("Error when adding %s at %s: DB Type not Supported", newUserProps.Username, targetDb.Name), Success: false}
		return
	}

	log := Log{
		DbId:    targetDb.Id,
		NewUser: newUserProps.Username,
		WO:      newUserProps.WO,
		UserId:  newUserProps.CurrentUserId,
		Action:  "CREATE",
	}
	go log.CreateLog()

	c <- TargetDbsRepose{Message: fmt.Sprintf("User %s has been created successfully at %s \n", newUserProps.Username, targetDb.Name), Success: true}
}

func (targetDb *TargetDb) ConnectToDBAndDeleteUser(newUserProps NewDbUserProps, c chan TargetDbsRepose, wg *sync.WaitGroup) {
	defer wg.Done()
	if targetDb.Type == "postgres" {
		pg, err := targetDb.connectToPostgre()
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
		defer pg.Close()

		_, err = pg.Exec("DROP USER IF EXISTS " + newUserProps.Username + ";")
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
	} else if targetDb.Type == "mysql" {
		database, err := targetDb.connectToSQL()
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
		defer database.Close()

		_, err = database.Exec("DROP USER IF EXISTS '" + newUserProps.Username + "'@'" + targetDb.Host + "';")
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
	} else if targetDb.Type == "oracle" {
		oracledb, err := targetDb.connectToOracle()
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
		defer oracledb.Close()

		_, err = oracledb.Exec("DROP USER " + newUserProps.Username + " CASCADE ")
		if err != nil {
			c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: %v", newUserProps.Username, targetDb.Name, err), Success: false}
			return
		}
	} else {
		c <- TargetDbsRepose{Message: fmt.Sprintf("Error when deleting %s at %s: DB Type not Supported", newUserProps.Username, targetDb.Name), Success: false}
		return
	}

	log := Log{
		DbId:    targetDb.Id,
		NewUser: newUserProps.Username,
		WO:      newUserProps.WO,
		UserId:  newUserProps.CurrentUserId,
		Action:  "DROP",
	}
	go log.CreateLog()

	c <- TargetDbsRepose{Message: fmt.Sprintf("User %s has been deleted successfully at %s \n", newUserProps.Username, targetDb.Name), Success: true}
}

func (targetDb *TargetDb) connectToSQL() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 eDbUser.Username,
		Passwd:               eDbUser.Password,
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

func (targetDb *TargetDb) connectToPostgre() (*sql.DB, error) {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=%s", eDbUser.Username, eDbUser.Password, targetDb.Host, targetDb.Port, targetDb.SslMode)

	database, err := sql.Open(targetDb.Type, connectionStr)
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
	connectionStr := go_ora.BuildUrl(targetDb.Host, targetDb.Port, targetDb.Name, eDbUser.Username, eDbUser.Password, nil)
	database, err := sql.Open(targetDb.Type, connectionStr)
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
