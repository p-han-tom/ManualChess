package userrepository

import (
	"fmt"
	"manual-chess/models/lobby"
)

type InMemPlayerRepository struct {
	players map[string]*lobby.User
}

func NewInMemPlayerRepository() *InMemPlayerRepository {
	return &InMemPlayerRepository{
		players: make(map[string]*lobby.User),
	}
}

func (r *InMemPlayerRepository) GetPlayerById(id string) (*lobby.User, error) {
	player, exists := r.players[id]
	if !exists {
		return nil, fmt.Errorf("player %s does not exist", id)
	}
	return player, nil
}

func (r *InMemPlayerRepository) SetPlayerById(id string, newPlayer *lobby.User) error {
	r.players[id] = newPlayer
	return nil
}
