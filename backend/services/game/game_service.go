package gameservice

import (
	"manual-chess/models/match"
	repository "manual-chess/repository/match"
	socketServices "manual-chess/services/socket"

	"github.com/google/uuid"
)

type GameService struct {
	socketService *socketServices.SocketService
	matchRepo     repository.IMatchRepository
}

func NewGameService(socketService *socketServices.SocketService, matchRepo repository.IMatchRepository) *GameService {
	return &GameService{
		socketService: socketService,
		matchRepo:     matchRepo,
	}
}

func (g *GameService) SetupMatch(id1 string, id2 string) {
	game := match.NewMatch(id1, id2)

	game.Player1.Colour = match.Blue
	game.Player2.Colour = match.Red
	// Hard coding game for testing

	p1Unit1, p1Unit2, p2Unit1, p2Unit2 := match.UnitFactory("necromancer"), match.UnitFactory("hedge_knight"), match.UnitFactory("necromancer"), match.UnitFactory("hedge_knight")
	p1Unit1.IsDeployed = true
	p1Unit2.IsDeployed = true
	p2Unit1.IsDeployed = true
	p2Unit2.IsDeployed = true
	p1Unit1.Pos = match.Position{Row: 7, Col: 2}
	p1Unit2.Pos = match.Position{Row: 7, Col: 6}
	p2Unit1.Pos = match.Position{Row: 2, Col: 2}
	p2Unit2.Pos = match.Position{Row: 2, Col: 6}

	game.Player1.Units[uuid.New().String()] = p1Unit1
	game.Player1.Units[uuid.New().String()] = p1Unit2
	game.Player2.Units[uuid.New().String()] = p2Unit1
	game.Player2.Units[uuid.New().String()] = p2Unit2

	g.matchRepo.SetMatch(game.ID, &game)
	g.socketService.GetConnection(id1).WriteJSON(map[string]interface{}{"matchId": game.ID})
	g.socketService.GetConnection(id2).WriteJSON(map[string]interface{}{"matchId": game.ID})

	// g.runPickPhase(game.ID)
	g.runGamePhase(game.ID)
}
