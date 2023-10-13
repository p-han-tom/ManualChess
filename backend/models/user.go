package models

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
