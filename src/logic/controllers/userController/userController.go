package userController

import (
	"database/sql"
	"src/db/userRepo"
	"src/objects"
)

type UserController struct {
	Repo userRepo.UserRepo
}

func (uc *UserController) GetUserID(login string) (int, error) {
	id, err := uc.Repo.GetUserID(login)
	if err == sql.ErrNoRows {
		err = UserNotFoundErr
	}
	return id, err
}

func (uc *UserController) GetUser(id int) (objects.User, error) {
	user, err := uc.Repo.GetUser(id)
	if err == sql.ErrNoRows {
		err = UserNotFoundErr
	}
	return user, err
}
func (uc *UserController) AddUser(login, password string, privelegeLevel objects.Levels) error {
	if login == objects.EmptyString || password == objects.EmptyString {
		return BadParamsErr
	}
	_, err := uc.Repo.GetUserID(login)
	if err == sql.ErrNoRows {
		return uc.Repo.AddUser(login, password, privelegeLevel)
	} else {
		return LoginOccupedErr
	}
}

func (uc *UserController) UserExist(login string) bool {
	var result bool
	id, err := uc.Repo.GetUserID(login)
	if err == nil && id > objects.None {
		result = true
	}
	return result
}
