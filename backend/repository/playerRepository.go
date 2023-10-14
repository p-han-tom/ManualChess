package repository

import "manual-chess/models"

type PlayerRepository interface {
	GetPlayerById(id string) (models.User, error)
	GetPlayersByEloRange(eloMin float64, eloMax float64) ([]models.User, error)
	SetPlayerById(id string, newPlayer models.User) error
}
