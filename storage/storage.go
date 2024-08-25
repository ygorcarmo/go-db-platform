package storage

import (
	"github.com/ygorcarmo/db-platform/models"
)

type Storage interface {
	Seed() error // this is only for testing purposes
	GetUserById(string) (*models.AppUser, error)
	GetAllUsers() ([]models.AppUser, error)
	GetUserByUsername(string) (*models.AppUser, error)
	GetDbsName() ([]string, error)
	GetDbByName(string) (*models.TargetDb, error)
	GetAllDbs() ([]models.TargetDb, error)
	CreateExternalDb(models.TargetDb) error
	CreateLog(models.Log) error
}
