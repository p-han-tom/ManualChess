package repository

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisMatchmakingRepository struct {
	redisClient *redis.Client
}

func NewRedisMatchmakingRepository(r *redis.Client) *RedisMatchmakingRepository {
	return &RedisMatchmakingRepository{
		redisClient: r,
	}
}

func (r *RedisMatchmakingRepository) GetPlayersByEloRange(eloMin int, eloMax int) ([]redis.Z, error) {
	return r.redisClient.ZRangeByScoreWithScores(
		context.Background(),
		"mmList",
		&redis.ZRangeBy{Min: strconv.Itoa(eloMin), Max: strconv.Itoa(eloMax)}).Result()

}

func (r *RedisMatchmakingRepository) AddPlayer(id string, score int) error {
	r.redisClient.ZAdd(context.Background(), "mmList", redis.Z{Score: float64(score), Member: id})
	return nil
}

func (r *RedisMatchmakingRepository) RemovePlayer(id string) error {
	r.redisClient.ZRem(context.Background(), id)
	return nil
}
