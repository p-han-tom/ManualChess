package constants

type MatchState int

const (
	Pick MatchState = iota
	Deploy
	Play
)

type UserState int

const (
	InLobby UserState = iota
	InQueue
	InGame
)
