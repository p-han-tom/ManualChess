package models

import "manual-chess/constants"

type Match struct {
	ID          string               `json:"id"`
	State       constants.MatchState `json:"state"`
	Board       [15][15]string       `json:"board"`
	Action      string               `json:"action"` // string of the user's secret id
	Player1     string               `json:"playerOne"`
	Player2     string               `json:"playerTwo"`
	Player1Pool string               `json:"playerOnePool"`
	Player2Pool string               `json:"playerTwoPool"`

	P1Units map[string]Unit `json:"p1Units"`
	P2Units map[string]Unit `json:"p2Units"`
}
