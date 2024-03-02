package repository

import (
	"manual-chess/models"
)

type IUserRepository interface {
	GetUserById(id string) (*models.User, error)
	SetUserById(id string, newUser *models.User) error
}
