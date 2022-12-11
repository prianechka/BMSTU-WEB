package middleware

import (
	"net/http"
	"src/objects"
	"src/utils/access"
	"src/utils/jwtUtils"
)

func CheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("jwt-token")
		role := jwtUtils.GetRoleFromJWT(accessToken)
		if access.CheckRoleAccess(r.RequestURI, r.Method, role) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, objects.InternalServerErrorString, http.StatusInternalServerError)
		}
	})
}
