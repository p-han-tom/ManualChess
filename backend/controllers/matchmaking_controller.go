package controllers

import (
	"fmt"
	. "manual-chess/dtos/request"
	"manual-chess/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MatchMakingController struct {
	matchMakingService *services.MatchMakingService
	socketService      *services.SocketService
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewMatchMakingController(m *services.MatchMakingService, s *services.SocketService) *MatchMakingController {
	return &MatchMakingController{
		matchMakingService: m,
		socketService:      s,
	}
}

// @POST Add user to redis match making list and establish socket connection
func (m *MatchMakingController) FindMatch(c *gin.Context) {
	// Bind the JSON data from the request body to the User struct
	id := c.Param("id")

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	m.socketService.SetConnection(id, conn)

	m.matchMakingService.AddToMatchMakingQueue(id)
}

// @DELETE Remove user from redis match making list
func (m *MatchMakingController) CancelMatch(c *gin.Context) {
	// Use go channels to signal
	var request MatchMakingRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.socketService.RemoveConnection(request.ID)

	m.matchMakingService.RemoveFromMatchMakingQueue(request.ID)

	c.IndentedJSON(http.StatusOK, request)
}
