package models

import (
	"custom-db-platform/src/db"
	"time"
)

type AppUser struct {
	Id         string
	Username   string
	Password   string
	Supervisor interface{}
	Sector     interface{}
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

func (*AppUser) GetAllUsers() ([]AppUser, error) {
	var users []AppUser
	rows, err := db.Database.Query("SELECT BIN_TO_UUID(id), username, password, supervisor, sector, isAdmin, createdAt, updatedAt FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user AppUser
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Supervisor, &user.Sector, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
