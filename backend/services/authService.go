package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	RedisClient *redis.Client
}

// TODO: Should retrieve user from persistent database
func (a *AuthService) Login(id string, val []byte) {
	a.RedisClient.Set(context.Background(), id, val, 0)
}

func (a *AuthService) Logout(id string) {
	a.RedisClient.Del(context.Background(), id)
}
