package utils

import (
	"context"

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
