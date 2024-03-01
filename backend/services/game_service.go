package services

import (
	"fmt"
	"manual-chess/constants"
	"manual-chess/models"
	repository "manual-chess/repository/match"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameService struct {
	socketService *SocketService
	matchRepo     repository.IMatchRepository
}

func NewGameService(socketService *SocketService, matchRepo repository.IMatchRepository) *GameService {
	return &GameService{
		socketService: socketService,
		matchRepo:     matchRepo,
	}
}

func (g *GameService) SetupMatch(id1 string, id2 string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomNumber := r.Float64()
	actionFirst := ""
	if randomNumber < 0.5 {
		actionFirst = id1
	} else {
		actionFirst = id2
	}

	match := models.Match{
		ID:      uuid.New().String(),
		State:   constants.Select,
		Player1: id1,
		Player2: id2,
		Action:  actionFirst,
	}

	g.matchRepo.SetMatch(match.ID, &match)
	g.socketService.GetConnection(id1).WriteJSON(map[string]interface{}{"matchId": match.ID})
	g.socketService.GetConnection(id2).WriteJSON(map[string]interface{}{"matchId": match.ID})

	go g.runMatch(match.ID, id1, id2)
}

func (g *GameService) runMatch(matchId string, id1 string, id2 string) {
	conn1 := g.socketService.GetConnection(id1)
	conn2 := g.socketService.GetConnection(id2)
	match, _ := g.matchRepo.GetMatch(matchId)
	turn := match.Action

	fmt.Println("First action: " + turn)
	for {
		var data map[string]interface{}
		var player *websocket.Conn
		if turn == id1 {
			player = conn1
		} else {
			player = conn2
		}
		// TODO: Add turn start verification for safer action processing

		err := player.ReadJSON(&data)
		for err != nil {
			fmt.Println(err.Error())
			err = player.ReadJSON(&data)
		}

		fmt.Println(turn + ": " + data["text"].(string))

		// Phase 1: Roster pick

		// end of turn
		if turn == id1 {
			turn = id2
		} else {
			turn = id1
		}

	}
}
