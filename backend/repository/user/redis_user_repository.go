package repository

import (
	"manual-chess/models"
	"manual-chess/utils"

	"github.com/redis/go-redis/v9"
)

type RedisUserRepository struct {
	redisClient *redis.Client
}

func NewRedisUserRepository(redisClient *redis.Client) *RedisUserRepository {
	return &RedisUserRepository{redisClient: redisClient}
}

func (r *RedisUserRepository) GetUserById(id string) (*models.User, error) {
	return utils.GetAndUnmarshal[models.User](r.redisClient, id)
}

func (r *RedisUserRepository) SetUserById(id string, newUser *models.User) error {
	return utils.MarshalAndSet[models.User](r.redisClient, id, newUser)
}
