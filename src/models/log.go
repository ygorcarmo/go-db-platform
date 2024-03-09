package models

import (
	"custom-db-platform/src/db"
	"time"
)

type Log struct {
	Id        string
	DbId      string
	NewUser   string
	WO        int
	UserId    string
	CreatedAt time.Time
}

type LogResponse struct {
	NewUser   string
	WO        int
	CreatedBy string
	Database  string
	CreatedAt time.Time
}

func (log *Log) CreateLog() (*Log, error) {
	return nil, nil
}

func (*Log) GetAllLogsPretty() ([]LogResponse, error) {
	var logs []LogResponse

	rows, err := db.Database.Query("SELECT l.newUser, ed.name AS external_database_name, l.wo, u.username, l.createdAt FROM logs l JOIN users u ON l.userId = u.id JOIN external_databases ed ON l.dbId = ed.id")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var log LogResponse
		err := rows.Scan(&log.NewUser, &log.Database, &log.WO, &log.CreatedBy, &log.CreatedAt)
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
