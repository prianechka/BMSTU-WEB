package userController

import "errors"

var (
	UserNotFoundErr = errors.New("user not found")
	LoginOccupedErr = errors.New("login is already exist")
	BadParamsErr    = errors.New("bad params")
)
