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
	socketService := services.NewSocketService()
	matchMakingService := services.NewMatchMakingService(socketService, redisClient)
	authService := services.NewAuthService(redisClient)

	// Set up controllers
	matchMakingController := controllers.NewMatchMakingController(matchMakingService, socketService)

	authController := controllers.AuthController{
		AuthService: authService,
	}

	go matchMakingService.RunMatchMaker()

	router := gin.Default()
	router.GET("/ws/findMatch/:id", matchMakingController.FindMatch)
	router.POST("/login", authController.Login)
	router.DELETE("/cancelMatch", matchMakingController.CancelMatch)
	router.DELETE("/logout", authController.Logout)
	router.Run("localhost:8080")
}
