package gameservice

import (
	"fmt"
	dtos "manual-chess/dtos/socket"
	"manual-chess/models/match"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (g *GameService) runPickPhase(matchId string) {
	game, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Match " + matchId + " not found")
		return
	}

	id1, id2 := game.Player1.ID, game.Player2.ID
	turn := game.Action
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
		var player *match.Player
		roster := game.Roster

		conn1.WriteJSON(game)
		conn2.WriteJSON(game)

		if turn == id1 && canPlayerPick(roster, game.Player1.Gold) {
			socket = conn1
			player = &game.Player1
		} else if turn == id2 && canPlayerPick(roster, game.Player2.Gold) {
			socket = conn2
			player = &game.Player2
		} else {
			fmt.Println("Pick phase is over")
			game.Action = nextAction
			g.matchRepo.SetMatch(matchId, game)
			go g.runDeployPhase(matchId)
			break
		}

		fmt.Println("It's " + turn + "'s turn")
		fmt.Println("Roster: ", roster)
		fmt.Println("Player: ", *player)
		fmt.Println()

		// TODO: Add turn start verification for safer action processing

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
		player.Units[unitId] = match.UnitFactory(unitType)

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
