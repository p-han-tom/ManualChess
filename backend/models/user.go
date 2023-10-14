package models

// Down the line need to add units to the User struct
// For now we hard code the units that a User spawns in the matchMakingService
type UserState int

const (
	InLobby UserState = iota
	InQueue
	InGame
)

type User struct {
	ID    string    `json:"id"`
	MMR   int       `json:"mmr"`
	State UserState `json:"state"`
}
