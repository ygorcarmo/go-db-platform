package storage

import "db-platform/models"

func (db *MySQLStorage) GetADConfig() (*models.LDAP, error) {
	// TODO: Username and password should be encrypted?
	config := models.LDAP{}

	err := db.connection.QueryRow(`
	SELECT
		connectionStr,
		username,
		passwd,
		topLevelDomain,
		secondLevelDomain,
		baseGroup,
		adminGroup,
		isDefault,
		adminGroupOU,
		baseGroupOU
	FROM
		ldap_config;`).Scan(
		&config.ConnectionStr,
		&config.Username,
		&config.Password,
		&config.TopLevelDomain,
		&config.SecondLevelDomain,
		&config.BaseGroup,
		&config.AdminGroup,
		&config.IsDefault,
		&config.AdminGroupOU,
		&config.BaseGroupOU)

	if err != nil {
		return nil, err
	}
	return &config, nil
}
