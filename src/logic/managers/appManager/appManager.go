package appManager

import "src/objects"

type AppManager struct {
	currentLogin string
	currentRole  objects.Levels
}

func (am *AppManager) SetNewState(login string, role objects.Levels) {
	am.currentLogin = login
	am.currentRole = role
}

func (am *AppManager) SetNewLogin(login string) {
	am.currentLogin = login
}

func (am *AppManager) SetCurrentRole(role objects.Levels) {
	am.currentRole = role
}

func (am *AppManager) GetCurrentLogin() string {
	return am.currentLogin
}

func (am *AppManager) GetCurrentRole() objects.Levels {
	return am.currentRole
}

func (am *AppManager) FoldState() {
	am.currentRole = objects.NonAuth
	am.currentLogin = objects.EmptyString
}
