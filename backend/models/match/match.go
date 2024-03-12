package match

import (
	"manual-chess/constants"
	"math/rand"
	"sort"

	"github.com/google/uuid"
)

type Turn struct {
	PlayerID   string `json:"playerId"`
	UnitID     string `json:"unitId"`
	Initiative int    `json:"initiative"`
}

type Match struct {
	ID        string               `json:"id"`
	State     constants.MatchState `json:"state"`
	Board     [][]Tile             `json:"board"`
	Action    string               `json:"action"` // string of the user's secret id
	TurnOrder []Turn               `json:"turnOrder"`
	Round     int                  `json:"round"`
	Player1   Player               `json:"playerOne"`
	Player2   Player               `json:"playerTwo"`
	Roster    [][]string           `json:"roster"`
}

const (
	BoardHeight int = 8
	BoardWidth  int = 8
)

func NewMatch(id1 string, id2 string) Match {
	// Determine first action
	randomNumber := rand.Float64()
	var actionFirst string
	if randomNumber < 0.5 {
		actionFirst = id1
	} else {
		actionFirst = id2
	}

	game := Match{
		ID:      uuid.New().String(),
		State:   constants.Pick,
		Player1: Player{ID: id1, Colour: PlayerColour(Blue), Units: make(map[string]Unit), Gold: 6},
		Player2: Player{ID: id2, Colour: PlayerColour(Red), Units: make(map[string]Unit), Gold: 6},
		Action:  actionFirst,
	}

	// Generate board
	var board [][]Tile
	for i := 0; i < BoardHeight; i++ {
		board = append(board, []Tile{})
		for j := 0; j < BoardWidth; j++ {
			board[i] = append(board[i], Tile{Type: TileType(Grass), Status: TileStatus(Normal), Passable: true})
		}
	}
	game.Board = board

	// Generate roster
	var roster [][]string
	for i := 0; i < 3; i++ {
		roster = append(roster, []string{})
	}
	for i := 0; i < 5; i++ {
		roster[0] = append(roster[0], OneCostPool[rand.Intn(len(OneCostPool))])
		roster[1] = append(roster[1], TwoCostPool[rand.Intn(len(TwoCostPool))])
		roster[2] = append(roster[2], ThreeCostPool[rand.Intn(len(ThreeCostPool))])
	}
	game.Roster = roster

	return game
}

func (m *Match) RollInitiative() {
	var turnOrder []Turn
	for key, unit := range m.Player1.Units {
		turnOrder = append(turnOrder, Turn{PlayerID: m.Player1.ID, UnitID: key, Initiative: rand.Intn(7) + unit.Speed})
	}

	for key, unit := range m.Player2.Units {
		turnOrder = append(turnOrder, Turn{PlayerID: m.Player2.ID, UnitID: key, Initiative: rand.Intn(7) + unit.Speed})
	}

	sort.Slice(turnOrder, func(i, j int) bool {
		return turnOrder[i].Initiative > turnOrder[j].Initiative
	})

	m.TurnOrder = turnOrder
}
