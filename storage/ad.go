package storage

import "db-platform/models"

func (db *MySQLStorage) GetADConfig() (*models.LDAPConfig, error) {
	// TODO: Username and password should be encrypted?
	config := models.LDAPConfig{}

	err := db.connection.QueryRow(`
	SELECT
		connectionStr,
		username,
		passwd,
		topLevelDomain,
		secondLevelDomain
	FROM
		config;`).Scan(&config.ConnectionStr, &config.Username, &config.Password, &config.TopLevelDomain, &config.SecondLevelDomain)

	if err != nil {
		return nil, err
	}
	return &config, nil
}
