package access

import (
	"src/objects"
	"strings"
)

func CheckRoleAccess(requestURI string, method string, role objects.Levels) (result bool) {
	if role > objects.NonAuth || requestURI == objects.AuthURI || strings.Contains(requestURI, "swagger") {
		result = true
	}
	return result
}
