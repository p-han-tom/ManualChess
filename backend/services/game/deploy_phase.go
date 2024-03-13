package gameservice

import (
	"fmt"
	dtos "manual-chess/dtos/socket"
	"manual-chess/models/match"
	"sync"
)

func (g *GameService) runDeployPhase(matchId string) {
	var wg sync.WaitGroup
	match, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Match " + matchId + " not found")
		return
	}

	id1, id2 := match.Player1.ID, match.Player2.ID

	wg.Add(2)
	go g.processDeployForId(&wg, match, id1)
	go g.processDeployForId(&wg, match, id2)
	wg.Wait()

	g.runGamePhase(matchId)
}

func (g *GameService) processDeployForId(wg *sync.WaitGroup, game *match.Match, id string) {
	socket := g.socketService.GetConnection(id)
	var player *match.Player
	if game.Player1.ID == id {
		player = &game.Player1
	} else {
		player = &game.Player2
	}
	side := player.Colour
	for {
		var data dtos.DeploymentDto
		err := socket.ReadJSON(&data)
		for err != nil {
			fmt.Println("Invalid input, try again")
			err = socket.ReadJSON(&data)
		}

		confirmPlacement := *data.ConfirmPlacement
		unitId := data.UnitID
		row := *data.Row
		col := *data.Col

		if confirmPlacement {
			invalidDeployment := false
			for _, unit := range player.Units {
				if !unit.IsDeployed {
					invalidDeployment = true
					break
				}
			}
			if invalidDeployment {
				fmt.Println("Not all units are deployed")
				continue
			}
			break
		}

		if isValidUnitDeployment(game.Board, side, row, col) {
			if entry, ok := player.Units[unitId]; ok {
				entry.Pos.Row = row
				entry.Pos.Col = col
				entry.IsDeployed = true
			}
		}

	}

	wg.Done()
}

func isValidUnitDeployment(board [][]match.Tile, side match.PlayerColour, row int, col int) bool {
	validDeploy := row >= 0 && col >= 0 && row < match.BoardHeight && col < match.BoardWidth &&
		board[row][col].OccupantId == "" && board[row][col].Passable
	if side == match.PlayerColour(match.Blue) {
		return validDeploy && row < 3
	} else {
		return validDeploy && row >= 5
	}
}
