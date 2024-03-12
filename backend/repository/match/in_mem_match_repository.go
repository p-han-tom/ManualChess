package matchrepository

import (
	"fmt"
	"manual-chess/models/match"
)

type InMemMatchRepository struct {
	matches map[string]*match.Match
}

func NewInMemMatchRepository() *InMemMatchRepository {
	return &InMemMatchRepository{matches: make(map[string]*match.Match)}
}

func (r *InMemMatchRepository) GetMatch(id string) (*match.Match, error) {
	res, exists := r.matches[id]
	if !exists {
		return nil, fmt.Errorf("match id: %s does not exist", id)
	}
	return res, nil
}

func (r *InMemMatchRepository) SetMatch(id string, match *match.Match) error {
	r.matches[id] = match
	return nil
}
