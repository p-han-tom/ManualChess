package models

type Match struct {
	Board   [15][15]string  `json:"board"`
	Action  string          `json:"action"` // string of the user's secret id
	Player1 string          `json:"playerOne"`
	Player2 string          `json:"playerTwo"`
	P1Units map[string]Unit `json:"p1Units"`
	P2Units map[string]Unit `json:"p2Units"`
}
