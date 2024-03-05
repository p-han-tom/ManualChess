package services

import (
	"fmt"
	"manual-chess/models"
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

func (g *GameService) processDeployForId(wg *sync.WaitGroup, match *models.Match, id string) {
	socket := g.socketService.GetConnection(id)
	var player *models.Player
	if match.Player1.ID == id {
		player = &match.Player1
	} else {
		player = &match.Player2
	}
	side := player.Colour
	for {
		var data map[string]interface{}
		err := socket.ReadJSON(&data)
		if err != nil {
			fmt.Println("Invalid input, try again")
			err = socket.ReadJSON(&data)
		}

		unitId := data["unitId"].(string)
		row := int(data["row"].(float64))
		col := int(data["col"].(float64))

		if isValidUnitDeployment(match.Board, side, row, col) {
			if entry, ok := player.Units[unitId]; ok {
				entry.Row = row
				entry.Col = col
				entry.IsDeployed = true
			}
		}

		if data["confirmPlacement"].(bool) {
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
	}

	wg.Done()
}

func isValidUnitDeployment(board [][]models.Tile, side models.PlayerColour, row int, col int) bool {
	validDeploy := row >= 0 && col >= 0 && row < models.BoardHeight && col < models.BoardWidth &&
		board[row][col].OccupantId == "" && board[row][col].Passable == true
	if side == models.PlayerColour(models.Blue) {
		return validDeploy && row < 3
	} else {
		return validDeploy && row >= 5
	}
}
