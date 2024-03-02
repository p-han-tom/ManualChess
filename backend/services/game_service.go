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

	// Determine first action
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomNumber := r.Float64()
	actionFirst := ""
	if randomNumber < 0.5 {
		actionFirst = id1
	} else {
		actionFirst = id2
	}

	// Generate roster
	match := models.Match{
		ID:      uuid.New().String(),
		State:   constants.Select,
		Player1: models.Player{ID: id1, Units: make(map[string]models.Unit), Gold: 6},
		Player2: models.Player{ID: id2, Units: make(map[string]models.Unit), Gold: 6},
		Action:  actionFirst,
		Roster:  models.GenerateRoster(),
	}

	g.matchRepo.SetMatch(match.ID, &match)
	g.socketService.GetConnection(id1).WriteJSON(map[string]interface{}{"matchId": match.ID})
	g.socketService.GetConnection(id2).WriteJSON(map[string]interface{}{"matchId": match.ID})

	go g.runPickPhase(match.ID, id1, id2)
}

func (g *GameService) runPickPhase(matchId string, id1 string, id2 string) {
	conn1 := g.socketService.GetConnection(id1)
	conn2 := g.socketService.GetConnection(id2)
	match, _ := g.matchRepo.GetMatch(matchId)
	turn := match.Action

	match, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Shutting down and closing sockets")
		conn1.Close()
		conn2.Close()
		return
	}

	fmt.Println("First pick: " + turn)
	for {
		conn1.WriteJSON(match)
		conn2.WriteJSON(match)

		var data map[string]interface{}
		var socket *websocket.Conn
		var player *models.Player
		if turn == id1 {
			socket = conn1
			player = &match.Player1
		} else {
			socket = conn2
			player = &match.Player2
		}

		// TODO: Add turn start verification for safer action processing

		err := socket.ReadJSON(&data)
		for err != nil {
			fmt.Println(err.Error())
			err = socket.ReadJSON(&data)
		}

		fmt.Println(turn + ": " + data["text"].(string))

		// Phase 1: Roster pick
		roster := match.Roster
		for !isValidRosterPick(roster, data["rosterRow"].(int), data["rosterCol"].(int)) {
			socket.ReadJSON(&data)
		}
		rosterRow, rosterCol := data["rosterRow"].(int), data["rosterCol"].(int)
		unitType := roster[rosterRow][rosterCol]
		roster[rosterRow] = append(roster[rosterRow][:rosterCol], roster[rosterRow][rosterCol+1:]...)

		player.Gold -= (rosterRow + 1)
		unitId := uuid.New().String()
		fmt.Println(unitId, unitType)

		// end of turn
		if turn == id1 {
			turn = id2
		} else {
			turn = id1
		}

	}
}

func isValidRosterPick(roster [][]string, row int, col int) bool {
	return row < len(roster) && row >= 0 && len(roster[row]) > 0 && col < len(roster[row]) && col >= 0
}
