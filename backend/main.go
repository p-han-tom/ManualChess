package main

import (
	"context"
	"fmt"
	"manual-chess/controllers"
	matchmaking "manual-chess/infrastructure/matchmaking"
	matchrepository "manual-chess/repository/match"
	userrepository "manual-chess/repository/user"
	authservice "manual-chess/services/auth"
	gameservice "manual-chess/services/game"
	matchmakingservice "manual-chess/services/matchmaking"
	socketservice "manual-chess/services/socket"

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

	// Set up infrastructure
	inMemMMQueue := matchmaking.NewInMemMatchmakingQueue()

	// Set up repositories
	// inMemPlayerRepo := playerRepository.NewInMemPlayerRepository()
	redisPlayerRepo := userrepository.NewRedisUserRepository(redisClient)
	// redisMatchRepo := matchRepository.NewRedisMatchRepository(redisClient)
	inMemMatchRepo := matchrepository.NewInMemMatchRepository()

	// Set up services
	socketService := socketservice.NewSocketService()
	gameService := gameservice.NewGameService(socketService, inMemMatchRepo)
	matchMakingService := matchmakingservice.NewMatchMakingService(gameService, redisPlayerRepo, inMemMatchRepo, inMemMMQueue)
	authService := authservice.NewAuthService(redisClient)

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
