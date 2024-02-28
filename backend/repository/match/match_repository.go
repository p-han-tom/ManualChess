package repository

import "manual-chess/models"

type MatchRepository interface {
	GetMatch(id string) (*models.Match, error)
	SetMatch(id string) error
}
