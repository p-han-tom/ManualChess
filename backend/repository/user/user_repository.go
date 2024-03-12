package userrepository

import "manual-chess/models/lobby"

type IUserRepository interface {
	GetUserById(id string) (*lobby.User, error)
	SetUserById(id string, newUser *lobby.User) error
}
