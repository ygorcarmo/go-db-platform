package models

import "time"

type AppUser struct {
	Id            string
	Username      string
	Password      string
	Supervisor    string
	Sector        string
	IsAdmin       bool
	LoginAttempts int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type contextKey string

const UserCtx contextKey = "user"
