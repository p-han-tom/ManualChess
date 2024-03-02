package models

type Player struct {
	ID    string          `json:"id"`
	Units map[string]Unit `json:"units"`
	Gold  int             `json:"gold"`
}
