package services

import "fmt"

func (g *GameService) runGamePhase(matchId string) {

	match, err := g.matchRepo.GetMatch(matchId)
	if err != nil {
		fmt.Println("Match " + matchId + " not found")
		return
	}
	id1, id2 := match.Player1.ID, match.Player2.ID

	conn1 := g.socketService.GetConnection(id1)
	conn2 := g.socketService.GetConnection(id2)

	turn := match.Action
	fmt.Println("It's " + turn + "'s turn to move")

	for {

	}

}
