package storage

import "db-platform/models"

func (db *MySQLStorage) CreateAdminLog(l models.AdminLog) error {
	_, err := db.connection.Exec(`
	INSERT INTO admin_logs
		(action, resourceType, username, resourceName)
	VALUES
		(?, ?, ?,?);`, l.Action, l.ResourceType, l.Username, l.ResourceName)
	return err
}

func (db *MySQLStorage) GetAllAdminLogs() ([]models.AdminLog, error) {
	ls := []models.AdminLog{}
	rows, err := db.connection.Query(`
		SELECT
			action,
			username,
			resourceType,
			createdAt,
			resourceName
		FROM
			admin_logs
		ORDER BY
			createdAt DESC;
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		l := models.AdminLog{}
		err := rows.Scan(&l.Action, &l.Username, &l.ResourceType, &l.CreatedAt, &l.ResourceName)
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
