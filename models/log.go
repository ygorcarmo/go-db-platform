package models

import "time"

type ActionType string

const (
	Create    ActionType = "CREATE"
	Delete    ActionType = "DELETE"
	UPDATEPWD ActionType = "UPDATE_PASSWORD"
)

type Log struct {
	Id        string
	NewUser   string
	WO        int
	CreatedBy string
	CreatedAt time.Time
	DBName    string
	Action    ActionType
	Success   bool
}

type LogResponse struct {
	NewUser   string
	WO        int
	CreatedBy string
	Database  string
	CreatedAt time.Time
	Action    ActionType
	Success   bool
}
