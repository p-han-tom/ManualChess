package gameservice

import (
	"manual-chess/models/match"
	repository "manual-chess/repository/match"
	socketServices "manual-chess/services/socket"
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
	match := match.NewMatch(id1, id2)

	g.matchRepo.SetMatch(match.ID, &match)
	g.socketService.GetConnection(id1).WriteJSON(map[string]interface{}{"matchId": match.ID})
	g.socketService.GetConnection(id2).WriteJSON(map[string]interface{}{"matchId": match.ID})

	g.runPickPhase(match.ID)
}
