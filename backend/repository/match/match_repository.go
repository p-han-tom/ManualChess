package matchrepository

import "manual-chess/models/match"

type IMatchRepository interface {
	GetMatch(id string) (*match.Match, error)
	SetMatch(id string, match *match.Match) error
}
