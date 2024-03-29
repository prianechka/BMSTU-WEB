package authManager

import (
	"src/logic/controllers/userController"
	"src/objects"
	appErrors "src/utils/error"
)

type AuthManager struct {
	currentLogin   string
	userController userController.UserController
}

func CreateNewAuthManager(uc userController.UserController) *AuthManager {
	return &AuthManager{
		userController: uc,
	}
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
						err = appErrors.PasswordNotEqualErr
					}
				} else {
					err = getUserErr
				}
			}
		} else {
			err = getUserIDErr
		}
	} else {
		err = appErrors.UserNotFoundErr
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
		err = appErrors.UserNotFoundErr
	}
	return result, err
}

func (am *AuthManager) GetLogin() string {
	return am.currentLogin
}
