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
	GetDbById(string) (*models.ExternalDb, error)
	GetDbsName() ([]string, error)
	GetDbByName(string) (*models.ExternalDb, error)
	GetAllDbs() ([]models.ExternalDb, error)
	UpdateExternalDb(models.ExternalDb) error
	CreateExternalDb(models.ExternalDb) error
	CreateLog(models.Log) error
}
