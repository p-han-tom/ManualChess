package repository

import "manual-chess/models"

type IMatchRepository interface {
	GetMatch(id string) (*models.Match, error)
	SetMatch(id string, match *models.Match) error
}
