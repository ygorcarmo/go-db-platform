package storage

import (
	"db-platform/models"
	"db-platform/utils"
	"errors"
)

func (db *MySQLStorage) GetADConfigWithCredentials() (*models.LDAP, error) {
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
		baseGroupOU,
		timeOutInSecs,
		enableTLS,
		verifyCert,
		cert		
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
		&config.BaseGroupOU,
		&config.TimeOutInSecs,
		&config.EnableTLS,
		&config.VerifyCert,
		&config.Cert)

	if err != nil {
		return nil, err
	}

	eService, err := utils.NewEncryptionService()

	if err != nil {
		return nil, errors.New("something went wrong with the encryption service")
	}

	username, err := eService.Decrypt(config.Username)
	if err != nil {
		return nil, errors.New("unable to decrypt username")
	}

	config.Username = username

	passwd, err := eService.Decrypt(config.Password)
	if err != nil {
		return nil, errors.New("unable to decrypt password")
	}
	config.Password = passwd

	if config.VerifyCert {
		if len(config.Cert) == 0 {
			return nil, errors.New("please insert a certificate")
		}

		cert, err := eService.Decrypt(config.Cert)
		if err != nil {
			return nil, errors.New("unable to decrypt certificate")
		}
		config.Cert = cert
	}

	return &config, nil
}
func (db *MySQLStorage) GetADConfig() (*models.LDAP, error) {
	config := models.LDAP{}

	err := db.connection.QueryRow(`
	SELECT
		connectionStr,
		topLevelDomain,
		secondLevelDomain,
		baseGroup,
		adminGroup,
		isDefault,
		adminGroupOU,
		baseGroupOU,
		timeOutInSecs,
		enableTLS,
		verifyCert
	FROM
		ldap_config;`).Scan(
		&config.ConnectionStr,
		&config.TopLevelDomain,
		&config.SecondLevelDomain,
		&config.BaseGroup,
		&config.AdminGroup,
		&config.IsDefault,
		&config.AdminGroupOU,
		&config.BaseGroupOU,
		&config.TimeOutInSecs,
		&config.EnableTLS,
		&config.VerifyCert)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (db *MySQLStorage) UpdateADConfig(config models.LDAP) error {
	_, err := db.connection.Exec(`
		UPDATE 
			ldap_config
		SET
			connectionStr=?,
			topLevelDomain=?,
			secondLevelDomain=?,
			baseGroup=?,
			baseGroupOU=?,
			adminGroup=?,
			adminGroupOU=?,
			isDefault=?,
			timeOutInSecs=?,
			enableTLS=?,
			verifyCert=?;
	`,
		config.ConnectionStr,
		config.TopLevelDomain,
		config.SecondLevelDomain,
		config.BaseGroup,
		config.BaseGroupOU,
		config.AdminGroup,
		config.AdminGroupOU,
		config.IsDefault,
		config.TimeOutInSecs,
		config.EnableTLS,
		config.VerifyCert)
	if err != nil {
		return err
	}

	return nil
}

func (db *MySQLStorage) UpdateADCredentials(username, password string) error {
	eService, err := utils.NewEncryptionService()
	if err != nil {
		return err
	}

	eusername, err := eService.Encrypt(username)
	if err != nil {
		return err
	}

	epasswd, err := eService.Encrypt(password)
	if err != nil {
		return err
	}

	_, err = db.connection.Exec("UPDATE ldap_config	SET username=?, passwd=?;", eusername, epasswd)

	if err != nil {
		return err
	}

	return nil
}

func (db *MySQLStorage) GetIsADDefaultAndAdminGroup() (bool, string, error) {
	isDefault := false
	group := ""

	err := db.connection.QueryRow("SELECT isDefault, adminGroup FROM ldap_config;").Scan(&isDefault, &group)
	if err != nil {
		return isDefault, group, err
	}

	return isDefault, group, nil
}

func (db *MySQLStorage) UpdateADCACert(cert string) error {
	eService, err := utils.NewEncryptionService()
	if err != nil {
		return err
	}

	ecert, err := eService.Encrypt(cert)
	if err != nil {
		return err
	}

	_, err = db.connection.Exec("UPDATE ldap_config SET cert=?;", ecert)
	if err != nil {
		return err
	}

	return nil

}
