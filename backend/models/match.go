package models

import (
	"manual-chess/constants"
	"math/rand"

	"github.com/google/uuid"
)

type Match struct {
	ID      string               `json:"id"`
	State   constants.MatchState `json:"state"`
	Board   [][]Tile             `json:"board"`
	Action  string               `json:"action"` // string of the user's secret id
	Player1 Player               `json:"playerOne"`
	Player2 Player               `json:"playerTwo"`
	Roster  [][]string           `json:"roster"`
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

	// Generate roster
	return Match{
		ID:      uuid.New().String(),
		State:   constants.Pick,
		Board:   GenerateBoard(),
		Player1: Player{ID: id1, Colour: PlayerColour(Blue), Units: make(map[string]Unit), Gold: 6},
		Player2: Player{ID: id2, Colour: PlayerColour(Red), Units: make(map[string]Unit), Gold: 6},
		Action:  actionFirst,
		Roster:  GenerateRoster(),
	}
}

func GenerateBoard() [][]Tile {
	var board [][]Tile
	for i := 0; i < BoardHeight; i++ {
		board = append(board, []Tile{})
		for j := 0; j < BoardWidth; j++ {
			board[i] = append(board[i], Tile{Type: TileType(Grass), Status: TileStatus(Normal), Passable: true})
		}
	}
	return board
}
