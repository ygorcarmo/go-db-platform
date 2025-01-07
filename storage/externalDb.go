package storage

import (
	"db-platform/models"
	"db-platform/utils"
	"errors"
)

func (db *MySQLStorage) GetDbById(i string) (*models.ExternalDb, error) {
	r := models.ExternalDb{Id: i}
	err := db.connection.QueryRow(`
		SELECT
		 	name, host, port, type, sslMode, owner, protocol, host_fallback, port_fallback, protocol_fallback 
		FROM
			external_databases WHERE id = UUID_TO_BIN(?);`,
		i).Scan(&r.Name, &r.Host, &r.Port, &r.Type, &r.SslMode, &r.Owner, &r.Protocol, &r.HostFallback, &r.PortFallback, &r.ProtocolFallback)
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
    		BIN_TO_UUID(id),
    		name,
    		host,
    		port,
    		type,
    		sslMode,
    		createdBy,
			owner
		FROM 
    		external_databases
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
		ed.owner,
		ed.protocol,
		ed.host_fallback,
		ed.port_fallback,
		ed.protocol_fallback
	FROM 
    	external_databases ed 
	WHERE 
		name=?;`, name).Scan(&targetDb.Name, &targetDb.Host, &targetDb.Port, &targetDb.Type, &targetDb.SslMode, &targetDb.Id, &targetDb.Username, &targetDb.Password, &targetDb.Owner, &targetDb.Protocol, &targetDb.HostFallback, &targetDb.PortFallback, &targetDb.ProtocolFallback)
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

func (db *MySQLStorage) UpdateExternalDb(e models.ExternalDb) error {
	if e.Type == models.OracleDG {
		_, err := db.connection.Exec(`
		UPDATE external_databases
		SET name=?, host=?, port=?, type=?, sslMode=?, owner=?, protocol=?, host_fallback=?, port_fallback=?, protocol_fallback=?
		WHERE id = UUID_TO_BIN(?);
	`, e.Name, e.Host, e.Port, e.Type, e.SslMode, e.Owner, e.Protocol, e.HostFallback, e.PortFallback, e.ProtocolFallback, e.Id)
		return err
	}
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
		(name, host, port, type, sslMode, username, password, createdBy, owner, protocol, host_fallback, port_fallback, protocol_fallback)
	VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		edb.Name, edb.Host, edb.Port, edb.Type, edb.SslMode, edb.Username, edb.Password, edb.CreatedBy, edb.Owner, edb.Protocol, edb.HostFallback, edb.PortFallback, edb.ProtocolFallback)
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
