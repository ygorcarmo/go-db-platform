package storage

import "db-platform/models"

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
