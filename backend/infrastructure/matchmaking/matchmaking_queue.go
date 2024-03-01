package infrastructure

import (
	"github.com/redis/go-redis/v9"
)

type IMatchMakingQueue interface {
	GetPlayersByEloRange(eloMin int, eloMax int) ([]redis.Z, error)
	AddPlayer(id string, score int) error
	RemovePlayer(id string) error
}
