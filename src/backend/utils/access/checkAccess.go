package access

import "src/objects"

func CheckRoleAccess(requestURI string, method string, role objects.Levels) bool {
	return true
}
