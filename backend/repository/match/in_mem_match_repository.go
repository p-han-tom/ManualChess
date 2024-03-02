package repository

import (
	"fmt"
	"manual-chess/models"
)

type InMemMatchRepository struct {
	matches map[string]*models.Match
}

func NewInMemMatchRepository() *InMemMatchRepository {
	return &InMemMatchRepository{matches: make(map[string]*models.Match)}
}

func (r *InMemMatchRepository) GetMatch(id string) (*models.Match, error) {
	res, exists := r.matches[id]
	if !exists {
		return nil, fmt.Errorf("match id: %s does not exist", id)
	}
	return res, nil
}

func (r *InMemMatchRepository) SetMatch(id string, match *models.Match) error {
	r.matches[id] = match
	return nil
}
