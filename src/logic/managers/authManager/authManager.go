package authManager

import (
	"src/logic/controllers/userController"
	"src/objects"
)

type AuthManager struct {
	currentLogin   string
	userController userController.UserController
}

func (am *AuthManager) TryToAuth(login, password string) (result objects.Levels, err error) {
	isUserExist := am.userController.UserExist(login)
	if isUserExist {
		tmpID, getUserIDErr := am.userController.GetUserID(login)
		if getUserIDErr == nil {
			if tmpID != objects.None {
				tmpUser, getUserErr := am.userController.GetUser(tmpID)
				if getUserErr == nil {
					if tmpUser.GetPassword() == password {
						result = tmpUser.GetPrivelegeLevel()
					} else {
						err = PasswordNotEqualErr
					}
				} else {
					err = getUserErr
				}
			}
		} else {
			err = getUserIDErr
		}
	} else {
		err = userController.UserNotFoundErr
	}
	return result, err
}

func (am *AuthManager) GetUserID(login string) (int, error) {
	var result = objects.None
	var err error
	isUserExist := am.userController.UserExist(login)
	if isUserExist {
		result, err = am.userController.GetUserID(login)
	} else {
		err = userController.UserNotFoundErr
	}
	return result, err
}

func (am *AuthManager) GetLogin() string {
	return am.currentLogin
}
