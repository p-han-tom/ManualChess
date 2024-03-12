package matchrepository

import (
	"manual-chess/models/match"
	"manual-chess/utils"

	"github.com/redis/go-redis/v9"
)

type RedisMatchRepository struct {
	redisClient *redis.Client
}

func NewRedisMatchRepository(r *redis.Client) *RedisMatchRepository {
	return &RedisMatchRepository{redisClient: r}
}

func (r *RedisMatchRepository) GetMatch(id string) (*match.Match, error) {
	return utils.GetAndUnmarshal[match.Match](r.redisClient, id)
}

func (r *RedisMatchRepository) SetMatch(id string, matchInstance *match.Match) error {
	return utils.MarshalAndSet[match.Match](r.redisClient, id, matchInstance)
}
