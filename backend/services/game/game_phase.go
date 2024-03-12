package gameservice

import (
	"fmt"
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

	for {
		for _, turn := range game.TurnOrder {
			var playerId string = turn.PlayerID
			var unitId string = turn.UnitID
		}
	}

}
