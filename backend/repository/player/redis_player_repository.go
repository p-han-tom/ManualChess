package repository

import (
	"manual-chess/models"
	"manual-chess/utils"

	"github.com/redis/go-redis/v9"
)

type RedisPlayerRepository struct {
	redisClient *redis.Client
}

func NewRedisPlayerRepository() *RedisPlayerRepository {
	return &RedisPlayerRepository{}
}

func (r *RedisPlayerRepository) GetPlayerById(id string) (*models.User, error) {
	return utils.GetAndUnmarshal[models.User](r.redisClient, id)
}

func (r *RedisPlayerRepository) SetPlayerById(id string, newPlayer *models.User) {
	utils.MarshalAndSet[models.User](r.redisClient, id, newPlayer)
}
