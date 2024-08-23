package storage

import (
	"github.com/ygorcarmo/db-platform/models"
)

type Storage interface {
	Seed() error // this is only for testing purposes
	GetUserById(string) (*models.AppUser, error)
	GetUserByUsername(string) (*models.AppUser, error)
	GetAvailableDbs() ([]string, error)
	GetDbByName(string) (*models.TargetDb, error)
	CreateLog(models.Log) error
}
