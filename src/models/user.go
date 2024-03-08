package models

import (
	"custom-db-platform/src/db"
	"fmt"
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
	err := db.Database.QueryRow("SELECT username, password, isAdmin, supervisor, sector FROM users WHERE id=UUID_TO_BIN(?);", id).Scan(&user.Username, &user.Password, &user.IsAdmin, &user.Supervisor, &user.Sector)
	user.Id = id
	return err
}

func (user *AppUser) CreateUser() error {
	_, err := db.Database.Exec(`INSERT INTO users (username, password, supervisor, sector, isAdmin) VALUES (?, ?, ?, ?, ?);`, user.Username, user.Password, user.Supervisor, user.Sector, user.IsAdmin)
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

func (user *AppUser) UpdatePassword(hashedNewPassword string) error {
	fmt.Println("userid")
	fmt.Println(user.Id)
	_, err := db.Database.Exec("UPDATE users SET password = ? WHERE id = UUID_TO_BIN(?);", hashedNewPassword, user.Id)
	return err
}
