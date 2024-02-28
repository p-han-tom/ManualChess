package constants

type MatchState int

const (
	Select MatchState = iota
	Deploy
	Play
)

type UserState int

const (
	InLobby UserState = iota
	InQueue
	InGame
)
