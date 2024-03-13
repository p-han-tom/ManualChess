package gameservice

import (
	"fmt"
	"manual-chess/models/match"

	"github.com/gorilla/websocket"
)

func (g *GameService) runGamePhase(matchId string) {

	game, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Match " + matchId + " not found")
		return
	}
	id1, id2 := game.Player1.ID, game.Player2.ID

	conn1 := g.socketService.GetConnection(id1)
	conn2 := g.socketService.GetConnection(id2)

	turn := game.Action
	fmt.Println("It's " + turn + "'s turn to move")

	game.RollInitiative()

	var round int = 1

	for {
		for _, turn := range game.TurnOrder {
			var socket *websocket.Conn
			var player *match.Player

			moved, acted := false, false

			conn1.WriteJSON(game)
			conn2.WriteJSON(game)

			if turn.PlayerID == id1 {
				socket = conn1
				player = &game.Player1
			} else if turn.PlayerID == id2 {
				socket = conn2
				player = &game.Player2
			}

			fmt.Println("It's " + turn.PlayerID + "'s turn")
			fmt.Println("Unit's turn:")
			fmt.Println(player.Units[turn.UnitID])

			var data map[string]interface{}
			err := socket.ReadJSON(&data)
			for err != nil {
				fmt.Println("Invalid input, try again")
				err = socket.ReadJSON(&data)
			}

			// data schema?
			// EndTurn bool
			// MoveTo Position
			// Action ActionChoice
			// Targets []Position

			for !data["EndTurn"].(bool) {
				if !moved {

					moved = true
				}

				if !acted {

					acted = true
				}

				socket.ReadJSON(&data)
				for err != nil {
					fmt.Println("Invalid input, try again")
					err = socket.ReadJSON(&data)
				}
			}

			// Check end of action status of units, tiles they occupy, etc

		}
		fmt.Printf("Round %d is over\n", round)
	}

}

func pathExists(board [][]match.Tile, start match.Position, end match.Position, moveRange int) bool {
	// BFS to find path from start to end on board
	var directions = [4]match.Position{{Row: 1, Col: 0}, {Row: -1, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: -1}}
	rows, cols := len(board), len(board[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	queue := []match.Position{start}
	visited[start.Row][start.Col] = true

	for len(queue) > 0 {
		size := len(queue)
		if moveRange == 0 {
			return false
		}

		for i := 0; i < size; i++ {
			current := queue[0]
			queue = queue[1:]
			if start == end {
				return true
			}

			for _, dir := range directions {
				newRow, newCol := current.Row+dir.Row, current.Col+dir.Col
				if isValidTile(board, visited, newRow, newCol) {
					visited[newRow][newCol] = true
					queue = append(queue, match.Position{Row: newRow, Col: newCol})
				}
			}
		}
		moveRange--
	}

	return false
}

func isValidTile(board [][]match.Tile, visited [][]bool, row int, col int) bool {
	rows, cols := len(board), len(board[0])
	return row >= 0 && col >= 0 && row < rows && col < cols && !visited[row][col] && board[row][col].OccupantId == "" && board[row][col].Passable
}
