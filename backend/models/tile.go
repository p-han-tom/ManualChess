package models

type TileStatus int

const (
	Normal TileStatus = iota
	Fire
	Acid
	Mud
)

type TileType int

const (
	Grass TileType = iota
	Forest
	Hill
)

type Tile struct {
	Type       TileType   `json:"TileType"`
	OccupantId string     `json:"occupantId"`
	Passable   bool       `json:"passable"`
	Status     TileStatus `json:"status"`
}
