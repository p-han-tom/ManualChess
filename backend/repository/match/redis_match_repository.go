package repository

import (
	"manual-chess/models"
	"manual-chess/utils"

	"github.com/redis/go-redis/v9"
)

type RedisMatchRepository struct {
	redisClient *redis.Client
}

func NewRedisMatchRepository(r *redis.Client) *RedisMatchRepository {
	return &RedisMatchRepository{redisClient: r}
}

func (r *RedisMatchRepository) GetMatch(id string) (*models.Match, error) {
	return utils.GetAndUnmarshal[models.Match](r.redisClient, id)
}

func (r *RedisMatchRepository) SetMatch(id string, match *models.Match) error {
	return utils.MarshalAndSet[models.Match](r.redisClient, id, match)
}
