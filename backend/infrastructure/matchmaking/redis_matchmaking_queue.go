package infrastructure

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisMatchmakingQueue struct {
	redisClient *redis.Client
}

func NewRedisMatchmakingQueue(r *redis.Client) *RedisMatchmakingQueue {
	return &RedisMatchmakingQueue{
		redisClient: r,
	}
}

func (r *RedisMatchmakingQueue) GetPlayersByEloRange(eloMin int, eloMax int) ([]redis.Z, error) {
	return r.redisClient.ZRangeByScoreWithScores(
		context.Background(),
		"mmList",
		&redis.ZRangeBy{Min: strconv.Itoa(eloMin), Max: strconv.Itoa(eloMax)}).Result()

}

func (r *RedisMatchmakingQueue) AddPlayer(id string, score int) error {
	r.redisClient.ZAdd(context.Background(), "mmList", redis.Z{Score: float64(score), Member: id})
	return nil
}

func (r *RedisMatchmakingQueue) RemovePlayer(id string) error {
	r.redisClient.ZRem(context.Background(), id)
	return nil
}
