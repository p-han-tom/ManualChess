package repository

import (
	"fmt"
	"manual-chess/models"
)

type InMemPlayerRepository struct {
	players map[string]*models.User
}

func NewInMemPlayerRepository() *InMemPlayerRepository {
	return &InMemPlayerRepository{
		players: make(map[string]*models.User),
	}
}

func (r *InMemPlayerRepository) GetPlayerById(id string) (*models.User, error) {
	player, exists := r.players[id]
	if !exists {
		return nil, fmt.Errorf("player %s does not exist", id)
	}
	return player, nil
}

func (r *InMemPlayerRepository) SetPlayerById(id string, newPlayer *models.User) {
	r.players[id] = newPlayer
}
