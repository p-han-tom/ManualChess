package userrepository

import (
	"manual-chess/models/lobby"
	"manual-chess/utils"

	"github.com/redis/go-redis/v9"
)

type RedisUserRepository struct {
	redisClient *redis.Client
}

func NewRedisUserRepository(redisClient *redis.Client) *RedisUserRepository {
	return &RedisUserRepository{redisClient: redisClient}
}

func (r *RedisUserRepository) GetUserById(id string) (*lobby.User, error) {
	return utils.GetAndUnmarshal[lobby.User](r.redisClient, id)
}

func (r *RedisUserRepository) SetUserById(id string, newUser *lobby.User) error {
	return utils.MarshalAndSet[lobby.User](r.redisClient, id, newUser)
}
