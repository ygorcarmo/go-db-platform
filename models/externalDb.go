package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	go_ora "github.com/sijms/go-ora/v2"
)

type dbType string

const (
	Postgres dbType = "postgres"
	MySQL    dbType = "mysql"
	Oracle   dbType = "oracle"
	OracleDG dbType = "oracle-dg"
)

type ExternalDb struct {
	Id               string
	Name             string
	Host             string
	Port             int
	Type             dbType
	SslMode          string
	Username         string
	Password         string
	Owner            string
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Protocol         string
	HostFallback     string
	PortFallback     int
	ProtocolFallback string
}

type NewDbUserProps struct {
	Username      string
	Password      string
	CurrentUserId string
	WO            int
}
type ExternalDbResponse struct {
	Message   string
	IsSuccess bool
	DbName    string
}

func ToDbType(t string) (dbType, error) {
	switch t {
	case string(MySQL):
		return MySQL, nil
	case string(Postgres):
		return Postgres, nil
	case string(Oracle):
		return Oracle, nil
	case string(OracleDG):
		return OracleDG, nil
	default:
		return "", fmt.Errorf("invalid dbType: %s", t)
	}
}

func (t *ExternalDb) ConnectAndCreateUser(user NewDbUserProps) ExternalDbResponse {
	switch t.Type {
	case Postgres:
		pg, err := t.connectToPostgresql()

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbName: t.Name}
		}
		defer pg.Close()

		_, err = pg.Exec("CREATE USER \"" + user.Username + "\" WITH PASSWORD '" + user.Password + "';")

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbName: t.Name}
		}

	case MySQL:
		mysql, err := t.connectToSQL()

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbName: t.Name}
		}
		defer mysql.Close()

		_, err = mysql.Exec("CREATE USER '" + user.Username + "'@'" + t.Host + "' IDENTIFIED BY '" + user.Password + "';")
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbName: t.Name}
		}

	case Oracle, OracleDG:
		var db *sql.DB
		var err error
		if t.Type == Oracle {
			db, err = t.connectToOracle()
		} else if t.Type == OracleDG {
			db, err = t.connectToOracleDG()
		}

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbName: t.Name}
		}
		defer db.Close()

		uppername := strings.ToUpper(user.Username)

		_, err = db.Exec("CREATE USER \"" + uppername + "\" IDENTIFIED BY \"" + user.Password + "\"")
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Create), IsSuccess: false, DbName: t.Name}
		}

	default:
		return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, fmt.Errorf("DB type %s not supported", t.Type), Create), IsSuccess: false, DbName: t.Name}
	}

	return ExternalDbResponse{Message: fmt.Sprintf("User %s has been created successfully at %s \n", user.Username, t.Name), IsSuccess: true, DbName: t.Name}
}

func (t *ExternalDb) ConnectAndDeleteUser(user NewDbUserProps) ExternalDbResponse {
	switch t.Type {
	case Postgres:
		pg, err := t.connectToPostgresql()

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}
		defer pg.Close()

		_, err = pg.Exec("DROP USER \"" + user.Username + "\";")

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}

	case MySQL:
		mysql, err := t.connectToSQL()

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}
		defer mysql.Close()

		_, err = mysql.Exec("DROP USER '" + user.Username + "'@'" + t.Host + "';")
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}

	case Oracle, OracleDG:
		var db *sql.DB
		var err error
		if t.Type == Oracle {
			db, err = t.connectToOracle()
		} else if t.Type == OracleDG {
			db, err = t.connectToOracleDG()
		}

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}
		defer db.Close()

		upperName := strings.ToUpper(user.Username)

		_, err = db.Exec("DROP USER \"" + upperName + "\" CASCADE ")
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}

	default:
		return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, fmt.Errorf("DB type %s not supported", t.Type), Delete), IsSuccess: false, DbName: t.Name}
	}

	return ExternalDbResponse{Message: fmt.Sprintf("User %s has been deleted successfully at %s \n", user.Username, t.Name), IsSuccess: true, DbName: t.Name}
}
func (t *ExternalDb) ConnectAndUpdateUserPassword(user NewDbUserProps) ExternalDbResponse {
	switch t.Type {
	case Postgres:
		pg, err := t.connectToPostgresql()

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}
		defer pg.Close()

		_, err = pg.Exec("ALTER ROLE \"" + user.Username + "\" WITH PASSWORD '" + user.Password + "';")

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}

	case MySQL:
		mysql, err := t.connectToSQL()

		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}
		defer mysql.Close()

		_, err = mysql.Exec("ALTER USER '" + user.Username + "'@'" + t.Host + "' IDENTIFIED BY '" + user.Password + "';")
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}

	case Oracle:
		db, err := t.connectToOracle()
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}
		defer db.Close()

		upperName := strings.ToUpper(user.Username)

		_, err = db.Exec("ALTER USER \"" + upperName + "\" IDENTIFIED BY \"" + user.Password + "\"")
		if err != nil {
			return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, err, Delete), IsSuccess: false, DbName: t.Name}
		}

	default:
		return ExternalDbResponse{Message: makeErrorMessage(user.Username, t.Name, fmt.Errorf("DB type %s not supported", t.Type), Delete), IsSuccess: false, DbName: t.Name}
	}

	return ExternalDbResponse{Message: fmt.Sprintf("User %s's password has been updated successfully at %s \n", user.Username, t.Name), IsSuccess: true, DbName: t.Name}
}

func (targetDb *ExternalDb) connectToPostgresql() (*sql.DB, error) {
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

func (targetDb *ExternalDb) connectToSQL() (*sql.DB, error) {
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

func (targetDb *ExternalDb) connectToOracle() (*sql.DB, error) {
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

func (targetDb *ExternalDb) connectToOracleDG() (*sql.DB, error) {

	connStr := fmt.Sprintf(`(DESCRIPTION=
    (ADDRESS_LIST=
    	(LOAD_BALANCE=OFF)
        (FAILOVER=ON)
    	(address=(PROTOCOL=%s)(host=%s)(PORT=%v))
    	(address=(protocol=%s)(host=%s)(port=%v))
    )
    (CONNECT_DATA=
    	(SERVICE_NAME=%s))
    )`, targetDb.Protocol, targetDb.Host, targetDb.Port, targetDb.ProtocolFallback, targetDb.HostFallback, targetDb.PortFallback, targetDb.Name)

	connectionStr := go_ora.BuildJDBC(targetDb.Username, targetDb.Password, connStr, nil)

	database, err := sql.Open(string(Oracle), connectionStr)
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
