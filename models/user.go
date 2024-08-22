package models

import "time"

type AppUser struct {
	Id         string
	Username   string
	Password   string
	Supervisor string
	Sector     string
	IsAdmin    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
