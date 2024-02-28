package repository

import (
	"manual-chess/models"
)

type PlayerRepository interface {
	GetPlayerById(id string) (*models.User, error)
	SetPlayerById(id string, newPlayer *models.User)
}
