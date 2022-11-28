package userRepo

import "src/objects"

type UserRepo interface {
	GetUserID(login string) (int, error)
	GetUser(id int) (objects.User, error)
	AddUser(login, password string, privelegeLevel objects.Levels) error
}
