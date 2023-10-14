package main

import (
	"context"
	"fmt"
	"manual-chess/controllers"
	"manual-chess/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func main() {

	// Read config file
	viper.SetConfigFile("config.dev.yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading configuration:", err)
		panic(err)
	}

	// Set up redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.ADDR"),     // Redis server address
		Password: viper.GetString("redis.PASSWORD"), // Password (if set)
		DB:       0,                                 // Default DB
	})

	redisClient.FlushAll(context.Background())

	// Set up websockets
	// Idea: make map of socket connections (gorilla websockets)
	// Inject map of socket connections as needed per service or controller

	// Set up services
	matchMakingService := services.MatchMakingService{
		RedisClient: redisClient,
	}

	authService := services.AuthService{
		RedisClient: redisClient,
	}

	// Set up controllers
	matchMakingController := controllers.MatchMakingController{
		MatchMakingService: &matchMakingService,
	}

	authController := controllers.AuthController{
		AuthService: &authService,
	}

	go matchMakingService.RunMatchMaker()

	router := gin.Default()
	router.POST("/findMatch", matchMakingController.FindMatch)
	router.POST("/login", authController.Login)
	router.DELETE("/cancelMatch", matchMakingController.CancelMatch)
	router.DELETE("/logout", authController.Logout)
	router.Run("localhost:8080")
}
