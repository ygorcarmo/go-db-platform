package models

import "time"

type TargetDb struct {
	Id        string
	Name      string
	Host      string
	Port      int
	Type      string
	SslMode   string
	UserId    string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}
