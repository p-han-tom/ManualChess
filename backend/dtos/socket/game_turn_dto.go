package dtos

import "manual-chess/models/match"

// data schema?
// EndTurn bool
// MoveTo Position
// Action int
// Targets []Position

type GameTurnDto struct {
	EndTurn *bool            `json:"endTurn" validate:"required"`
	MoveTo  match.Position   `json:"moveTo" validate:"required"`
	Action  *int             `json:"action" validate:"required"`
	Targets []match.Position `json:"targets" validate:"required"`
}
