package services

import (
	"fmt"
	dtos "manual-chess/dtos/socket"
	"manual-chess/models"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (g *GameService) runPickPhase(matchId string) {
	match, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Match " + matchId + " not found")
		return
	}

	id1, id2 := match.Player1.ID, match.Player2.ID
	turn := match.Action
	var nextAction string
	if turn == id1 {
		nextAction = id2
	} else {
		nextAction = id1
	}

	conn1 := g.socketService.GetConnection(id1)
	conn2 := g.socketService.GetConnection(id2)

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
			match.Action = nextAction
			g.matchRepo.SetMatch(matchId, match)
			go g.runDeployPhase(matchId)
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
