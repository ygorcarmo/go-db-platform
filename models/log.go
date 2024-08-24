package models

import "time"

type ActionType string

const (
	Create ActionType = "CREATE"
	Delete ActionType = "DELETE"
)

type Log struct {
	Id        string
	DbId      string
	NewUser   string
	WO        int
	CreateBy  string
	CreatedAt time.Time
	Action    ActionType
	Sucess    bool
}
