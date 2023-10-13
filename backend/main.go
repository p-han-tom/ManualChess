package main

import (
	"fmt"
	"manual-chess/handlers"
	"manual-chess/services"
	"manual-chess/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.dev.yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading configuration:", err)
		panic(err)
	}

	utils.InitRedis()
	go services.InitMatchMaker()

	router := gin.Default()
	router.POST("/findMatch", handlers.FindMatch)
	router.POST("/login", handlers.Login)
	router.DELETE("/cancelMatch", handlers.CancelMatch)
	router.DELETE("/logout", handlers.Logout)
	router.Run("localhost:8080")
}
