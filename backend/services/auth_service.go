package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	redisClient *redis.Client
}

func NewAuthService(r *redis.Client) *AuthService {
	return &AuthService{
		redisClient: r,
	}
}

// TODO: Should retrieve user from persistent database
func (a *AuthService) Login(id string, val []byte) {
	a.redisClient.Set(context.Background(), id, val, 0)
}

func (a *AuthService) Logout(id string) {
	a.redisClient.Del(context.Background(), id)
}
