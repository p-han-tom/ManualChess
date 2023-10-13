package handlers

import (
	"context"
	"encoding/json"
	"manual-chess/models"
	"manual-chess/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Later we need to actually implement auth logic, not sure how that will work
// For now login and logout just creates and destroys a user in Redis

// POST request to login
func Login(c *gin.Context) {
	var user models.User

	// Bind the JSON data from the request body to the User struct
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	val, err := json.Marshal(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.RedisClient.Set(context.Background(), user.ID, val, 0)
}

// DELETE request to logout
func Logout(c *gin.Context) {
	var user models.User

	// Bind the JSON data from the request body to the User struct
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.RedisClient.Del(context.Background(), user.ID)
}
