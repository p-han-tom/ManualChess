package controllers

import (
	"encoding/json"
	. "manual-chess/models"
	"manual-chess/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

// Later we need to actually implement auth logic, not sure how that will work
// For now login and logout just creates and destroys a user in Redis

// POST request to login
func (a *AuthController) Login(c *gin.Context) {

	// TODO: Technically should just be receiving an AuthRequestDto
	// Since we don't have a persistent db set up, we just take a whole user object
	var user User

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

	a.AuthService.Login(user.ID, val)
}

// DELETE request to logout
func (a *AuthController) Logout(c *gin.Context) {
	var user User

	// Bind the JSON data from the request body to the User struct
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.AuthService.Logout(user.ID)
}
