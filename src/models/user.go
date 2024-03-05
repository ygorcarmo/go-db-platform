package models

import "custom-db-platform/src/db"

type AppUser struct {
	Username string
	Password string
}

func (user *AppUser) GetUser(username string) error {
	err := db.Database.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&user.Username, &user.Password)
	return err
}
