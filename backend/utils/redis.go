package utils

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func GetAndUnmarshal[T any](redisClient *redis.Client, key string) (*T, error) {
	var item T
	raw, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(raw), &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func MarshalAndSet[T any](redisClient *redis.Client, key string, item *T) error {
	marshaledItem, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = redisClient.Set(context.Background(), key, marshaledItem, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
