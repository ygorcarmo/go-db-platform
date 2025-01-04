package storage

import "db-platform/models"

func (db *MySQLStorage) CreateAdminLog(l models.AdminLog) error {
	_, err := db.connection.Exec(`
	INSERT INTO admin_logs
		(action, resourceId, resourceType, userId, resourceName)
	VALUES
		(?, UUID_TO_BIN(?), ?,UUID_TO_BIN(?),?);`, l.Action, l.ResourceId, l.ResourceType, l.UserId, l.ResourceName)
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
			l.createdAt,
			l.resourceName
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
		err := rows.Scan(&l.Action, &l.Username, &l.ResourceId, &l.ResourceType, &l.CreatedAt, &l.ResourceName)
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
