package dtos

type RosterPickDto struct {
	Row *int `json:"row" validate:"required"`
	Col *int `json:"col" validate:"required"`
}
