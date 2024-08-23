package storage

import (
	"github.com/ygorcarmo/db-platform/models"
)

type Storage interface {
	GetUserById(string) (*models.AppUser, error)
	GetUserByUsername(string) (*models.AppUser, error)
	GetAvailableDbs() ([]string, error)
}
