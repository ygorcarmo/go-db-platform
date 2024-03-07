package models

import (
	"custom-db-platform/src/db"
	"time"
)

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

func (user *AppUser) GetUserByUsername(username string) error {
	err := db.Database.QueryRow("SELECT BIN_TO_UUID(id), username, password FROM users WHERE username=?;", username).Scan(&user.Id, &user.Username, &user.Password)
	return err
}
func (user *AppUser) GetUserById(id string) error {
	err := db.Database.QueryRow("SELECT username, password, isAdmin FROM users WHERE id=UUID_TO_BIN(?);", id).Scan(&user.Username, &user.Password, &user.IsAdmin)
	return err
}

func (user *AppUser) CreateUser() error {
	_, err := db.Database.Exec(`INSERT INTO users (username, password, isAdmin) VALUES ($1, $2, $3);`, user.Username, user.Password, user.IsAdmin)
	return err
}
