package services

import (
	"fmt"
	"manual-chess/constants"
	dtos "manual-chess/dtos/socket"
	"manual-chess/models"
	repository "manual-chess/repository/match"
	"math/rand"

	"github.com/go-playground/validator/v10"
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
	randomNumber := rand.Float64()
	var actionFirst string
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

	validate := validator.New()

	fmt.Println("First pick: " + turn)
	for {
		var pick dtos.RosterPickDto
		var socket *websocket.Conn
		var player *models.Player
		roster := match.Roster

		conn1.WriteJSON(match)
		conn2.WriteJSON(match)

		if turn == id1 && canPlayerPick(roster, match.Player1.Gold) {
			socket = conn1
			player = &match.Player1
		} else if turn == id2 && canPlayerPick(roster, match.Player2.Gold) {
			socket = conn2
			player = &match.Player2
		} else {
			fmt.Println("Pick phase is over")
			g.matchRepo.SetMatch(matchId, match)
			go g.runGamePhase(matchId, id1, id2)
			break
		}

		fmt.Println("It's " + turn + "'s turn")
		fmt.Println("Roster: ", roster)
		fmt.Println("Player: ", *player)
		fmt.Println()

		// TODO: Add turn start verification for safer action processing

		// Phase 1 roster pick
		for {
			err := socket.ReadJSON(&pick)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			err = validate.Struct(pick)

			if err == nil && isValidRosterPick(roster, player.Gold, *pick.Row, *pick.Col) {
				break
			} else if err == nil {
				fmt.Println("Invalid roster pick, try again")
			} else {
				fmt.Println(err.Error())
			}
		}

		row, col := *pick.Row, *pick.Col
		unitType := roster[row][col]
		roster[row] = append(roster[row][:col], roster[row][col+1:]...)

		player.Gold -= (row + 1)
		unitId := uuid.New().String()
		player.Units[unitId] = models.UnitFactory(unitType)

		// end of turn
		if turn == id1 {
			turn = id2
		} else {
			turn = id1
		}
	}
}

func (g *GameService) runGamePhase(matchId string, id1 string, id2 string) {

	conn1 := g.socketService.GetConnection(id1)
	conn2 := g.socketService.GetConnection(id2)
	match, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Shutting down and closing sockets")
		conn1.Close()
		conn2.Close()
		return
	}

	// Determine first action
	randomNumber := rand.Float64()
	var actionFirst string
	if randomNumber < 0.5 {
		actionFirst = id1
	} else {
		actionFirst = id2
	}

	match.Action = actionFirst

	for {

	}

}

func canPlayerPick(roster [][]string, gold int) bool {
	for i := 0; i < gold; i++ {
		if len(roster[i]) > 0 {
			return true
		}
	}

	return false
}

func isValidRosterPick(roster [][]string, gold int, row int, col int) bool {
	return gold > row && row < len(roster) && row >= 0 && len(roster[row]) > 0 && col < len(roster[row]) && col >= 0
}
