package storage

import (
	"github.com/ygorcarmo/db-platform/models"
)

type Storage interface {
	Seed() error // this is only for testing purposes
	GetUserById(string) (*models.AppUser, error)
	GetAllUsers() ([]models.AppUser, error)
	GetUserByUsername(string) (*models.AppUser, error)
	CreateApplicationUser(models.AppUser) error
	UpdateApplicationUserPassword(id string, password string) error
	DeleteUserById(string) error
	GetDbById(string) (*models.ExternalDb, error)
	GetDbsName() ([]string, error)
	GetDbByName(string) (*models.ExternalDb, error)
	GetAllDbs() ([]models.ExternalDb, error)
	UpdateExternalDb(models.ExternalDb) error
	UpdateExternalDbCredentials(id string, username string, password string) error
	DeleteExternalDbById(string) error
	CreateExternalDb(models.ExternalDb) error
	CreateLog(models.Log) error
	GetAllLogs() ([]models.LogResponse, error)
}
