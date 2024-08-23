package models

import "time"

type actionType string

const (
	Create actionType = "CREATE"
	Delete actionType = "DELETE"
)

type Log struct {
	Id        string
	DbId      string
	NewUser   string
	WO        int
	CreateBy  string
	CreatedAt time.Time
	Action    actionType
	Sucess    bool
}
