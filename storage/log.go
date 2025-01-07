package storage

import "db-platform/models"

func (db *MySQLStorage) CreateLog(log models.Log) error {
	_, err := db.connection.Exec("INSERT INTO logs (dbName, newUser, wo, createdBy, action, success) VALUES (?,?,?,?, ?, ?);", log.DBName, log.NewUser, log.WO, log.CreatedBy, log.Action, log.Success)
	return err
}

func (db *MySQLStorage) GetAllLogs() ([]models.Log, error) {
	logs := []models.Log{}
	rows, err := db.connection.Query(`
	SELECT 
		newUser, 
		dbName, 
		wo, 
		createdBy,
		createdAt,
		action,
		success
	FROM 
		logs
	ORDER BY 
		createdAt DESC;`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		log := models.Log{}
		err := rows.Scan(&log.NewUser, &log.DBName, &log.WO, &log.CreatedBy, &log.CreatedAt, &log.Action, &log.Success)
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
