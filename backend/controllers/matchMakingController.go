package controllers

import (
	. "manual-chess/dtos/request"
	"manual-chess/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatchMakingController struct {
	MatchMakingService *services.MatchMakingService
}

// @POST Add user to redis match making list
func (m *MatchMakingController) FindMatch(c *gin.Context) {
	// Add user to redis store of active users looking for match
	var request MatchMakingRequestDto

	// Bind the JSON data from the request body to the User struct
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.MatchMakingService.AddToMatchMakingQueue(request.ID)

	c.IndentedJSON(http.StatusCreated, request)
}

// @DELETE Remove user from redis match making list
func (m *MatchMakingController) CancelMatch(c *gin.Context) {
	// Use go channels to signal
	var request MatchMakingRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.MatchMakingService.RemoveFromMatchMakingQueue(request.ID)

	c.IndentedJSON(http.StatusOK, request)
}
