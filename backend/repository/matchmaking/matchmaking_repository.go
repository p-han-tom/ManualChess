package repository

import (
	"github.com/redis/go-redis/v9"
)

type MatchMakingRepository interface {
	GetPlayersByEloRange(eloMin int, eloMax int) ([]redis.Z, error)
	AddPlayer(id string, score int) error
	RemovePlayer(id string) error
}
