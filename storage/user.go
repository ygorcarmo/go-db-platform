package storage

import (
	"db-platform/models"
)

func (db *MySQLStorage) GetUserById(id string) (*models.AppUser, error) {
	user := models.AppUser{Id: id}
	err := db.connection.QueryRow("SELECT username, password, isAdmin, supervisor, sector FROM users WHERE id=UUID_TO_BIN(?);", id).Scan(&user.Username, &user.Password, &user.IsAdmin, &user.Supervisor, &user.Sector)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *MySQLStorage) UpdateApplicationUserPassword(i string, p string) error {
	_, err := db.connection.Exec(`
	 	UPDATE users SET password = ? WHERE id=UUID_TO_BIN(?);
	`, p, i)
	return err
}

func (db *MySQLStorage) UpdateApplicationUser(u models.AppUser) error {
	_, err := db.connection.Exec(`
	 	UPDATE users 
		SET username=?, supervisor=?, sector=?, isAdmin=?
		WHERE id=UUID_TO_BIN(?);
	`, u.Username, u.Supervisor, u.Sector, u.IsAdmin, u.Id)
	return err
}

func (db *MySQLStorage) UpdateApplicationUserCredentials(u string, p string, i string) error {
	_, err := db.connection.Exec(`
		UPDATE users
		SET username=?, password=?
		WHERE id=UUID_TO_BIN(?);
	`, u, p, i)
	return err
}

func (db *MySQLStorage) GetUserByUsername(username string) (*models.AppUser, error) {
	user := models.AppUser{Username: username}
	err := db.connection.QueryRow("SELECT BIN_TO_UUID(id), username, password, loginAttempts FROM users WHERE username=?;", username).Scan(&user.Id, &user.Username, &user.Password, &user.LoginAttempts)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *MySQLStorage) IncreaseUserLoginAttempts(id string, attempts int) error {
	_, err := db.connection.Exec(`
		UPDATE users
		SET loginAttempts=? 
		WHERE id=UUID_TO_BIN(?);
	`, attempts, id)
	return err
}

func (db *MySQLStorage) ResetUserLoginAttempts(id string) error {
	_, err := db.connection.Exec(`
	 	UPDATE users
		SET loginAttempts=0
		WHERE id=UUID_TO_BIN(?);
	`, id)
	return err
}

func (db *MySQLStorage) GetAllUsers() ([]models.AppUser, error) {
	users := []models.AppUser{}
	rows, err := db.connection.Query("SELECT BIN_TO_UUID(id), username, password, supervisor, sector, isAdmin, createdAt, updatedAt, loginAttempts FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.AppUser{}
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Supervisor, &user.Sector, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt, &user.LoginAttempts); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *MySQLStorage) CreateApplicationUser(u models.AppUser) (string, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return "", err
	}

	// Insert user within the transaction
	_, err = tx.Exec(`
        INSERT INTO users 
        (username, password, supervisor, sector, isAdmin)
        VALUES (?, ?, ?, ?, ?);`,
		u.Username, u.Password, u.Supervisor, u.Sector, u.IsAdmin)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Retrieve the last inserted ID in the same transaction
	err = tx.QueryRow(`SELECT BIN_TO_UUID(id) FROM users WHERE id = LAST_INSERT_ID();`).Scan(&u.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Commit the transaction
	err = tx.Commit()

	return u.Id, err
}

func (db *MySQLStorage) DeleteUserById(id string) error {
	_, err := db.connection.Exec("DELETE FROM users WHERE id=UUID_TO_BIN(?);", id)
	return err
}
