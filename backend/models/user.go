package models

import "manual-chess/constants"

type User struct {
	ID    string              `json:"id"`
	MMR   int                 `json:"mmr"`
	State constants.UserState `json:"state"`
}
