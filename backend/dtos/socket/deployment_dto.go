package dtos

type DeploymentDto struct {
	ConfirmPlacement *bool  `json:"confirmPlacement" validate:"required"`
	UnitID           string `json:"unitId" validate:"required"`
	Row              *int   `json:"row" validate:"required"`
	Col              *int   `json:"col" validate:"required"`
}
