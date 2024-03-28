package gameservice

import (
	"fmt"
	dtos "manual-chess/dtos/socket"
	"manual-chess/models/match"

	"github.com/go-playground/validator/v10"
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
	validate := validator.New()

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

			var unit match.Unit = player.Units[turn.UnitID]

			fmt.Println("It's " + turn.PlayerID + "'s turn")
			fmt.Println("Unit's turn:")
			fmt.Println(player.Units[turn.UnitID])

			var data dtos.GameTurnDto

			for {

				socket.WriteJSON(player)

				// input validation
				err := socket.ReadJSON(&data)
				for {
					if err != nil {
						continue
					}
					err = validate.Struct(data)
					if err != nil {
						continue
					}
					break
				}

				var endTurn bool = *data.EndTurn
				var moveTo match.Position = data.MoveTo
				var actionChoice int = *data.Action
				var targets []match.Position = data.Targets

				if endTurn {
					break
				}

				if !moved {
					if pathExists(game.Board, unit.Pos, moveTo, unit.MoveRange) {
						unit.Pos = moveTo
						moved = true
					} else {
						fmt.Println("No path exists between (", unit.Pos.Row, unit.Pos.Col, ") and (", moveTo.Row, moveTo.Col, ")")
					}
				}

				if !acted {
					switch actionChoice {
					case 0: // no ability input
						break
					case 1: // primary ability
						match.PrimaryAbilityLookupTable[unit.Type](player.ID, turn.UnitID, game, targets)
					case 2: // secondary ability
						match.PrimaryAbilityLookupTable[unit.Type](player.ID, turn.UnitID, game, targets)
					default:
						fmt.Println("Invalid action choice")
						// error out
					}
					acted = true
				}

				player.Units[turn.UnitID] = unit
			}

			// Check end of action status of units, tiles they occupy, etc
			processEndOfTurn(game)
		}
		fmt.Printf("Round %d is over\n", round)
	}
}

func processEndOfTurn(game *match.Match) {
	for id, unit := range game.Player1.Units {
		if unit.IsAlive && unit.HP <= 0 {
			game.Board[unit.Pos.Row][unit.Pos.Col].OccupantId = ""
			unit.IsAlive = false
			game.Player1.Units[id] = unit
		}
	}

	for id, unit := range game.Player2.Units {
		if unit.IsAlive && unit.HP <= 0 {
			game.Board[unit.Pos.Row][unit.Pos.Col].OccupantId = ""
			unit.IsAlive = false
			game.Player2.Units[id] = unit
		}
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
		if moveRange < 0 {
			return false
		}

		for i := 0; i < size; i++ {
			current := queue[i]
			if current == end {
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
		queue = queue[size:]
		moveRange--
	}

	return false
}

func isValidTile(board [][]match.Tile, visited [][]bool, row int, col int) bool {
	rows, cols := len(board), len(board[0])
	return row >= 0 && col >= 0 && row < rows && col < cols && !visited[row][col] && board[row][col].OccupantId == "" && board[row][col].Passable
}
