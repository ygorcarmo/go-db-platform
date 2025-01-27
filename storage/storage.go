package storage

import (
	"db-platform/models"
)

type Storage interface {
	Seed() error // this is only for testing purposes
	GetUserById(string) (*models.AppUser, error)
	GetAllUsers() ([]models.AppUser, error)
	GetUserByUsername(string) (*models.AppUser, error)
	CreateApplicationUser(models.AppUser) (string, error)
	UpdateApplicationUserPassword(id string, password string) error
	UpdateApplicationUser(models.AppUser) error
	UpdateApplicationUserCredentials(username string, password string, id string) error
	IncreaseUserLoginAttempts(id string, attempts int) error
	ResetUserLoginAttempts(id string) error
	DeleteUserById(string) error
	GetDbById(string) (*models.ExternalDb, error)
	GetDbsName() ([]string, error)
	GetDbByName(string) (*models.ExternalDb, error)
	GetAllDbs() ([]models.ExternalDb, error)
	UpdateExternalDb(models.ExternalDb) error
	UpdateExternalDbCredentials(id string, username string, password string) error
	DeleteExternalDbById(string) error
	CreateExternalDb(models.ExternalDb) (string, error)
	CreateLog(models.Log) error
	GetAllLogs() ([]models.Log, error)
	CreateAdminLog(models.AdminLog) error
	GetAllAdminLogs() ([]models.AdminLog, error)
	GetADConfigWithCredentials() (*models.LDAP, error)
	GetADConfig() (*models.LDAP, error)
	UpdateADConfig(models.LDAP) error
	UpdateADCredentials(username, password string) error
	GetIsADDefaultAndAdminGroup() (bool, string, error)
	UpdateADCACert(cert string) error
}
