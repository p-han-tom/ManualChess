package utils

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.ADDR"),     // Redis server address
		Password: viper.GetString("redis.PASSWORD"), // Password (if set)
		DB:       0,                                 // Default DB
	})

	// Test the Redis connection here, and handle any connection errors if needed.
	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		// Handle the error, e.g., log it or exit the application.
		panic(err)
	}
}

func GetAndUnmarshal[T any](key string) (*T, error) {
	var item T
	raw, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(raw), &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func MarshalAndSet[T any](key string, item *T) error {
	marshaledItem, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = RedisClient.Set(context.Background(), key, marshaledItem, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
