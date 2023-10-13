package handlers

import (
	dtos "manual-chess/dtos/request"
	"manual-chess/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @POST Add user to redis match making list
func FindMatch(c *gin.Context) {
	// Add user to redis store of active users looking for match
	var request dtos.MatchMakingRequestDto

	// Bind the JSON data from the request body to the User struct
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services.AddToMatchMakingQueue(request.ID)

	c.IndentedJSON(http.StatusCreated, request)
}

// @DELETE Remove user from redis match making list
func CancelMatch(c *gin.Context) {
	// Use go channels to signal
	var request dtos.MatchMakingRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services.RemoveFromMatchMakingQueue(request.ID)

	c.IndentedJSON(http.StatusOK, request)
}
