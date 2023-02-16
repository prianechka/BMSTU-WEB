package middleware

import (
	"net/http"
	"src/objects"
	"src/utils/access"
	"src/utils/jwtUtils"
)

func CheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/test" {
			next.ServeHTTP(w, r)
		} else if r.RequestURI == "/status" {
			next.ServeHTTP(w, r)
		} else if r.RequestURI == "/next/index1.html/" || r.RequestURI == "/next/style1.css/" ||
			r.RequestURI == "/next/image1.jpg/" || r.RequestURI == "/index1.html/" {
			next.ServeHTTP(w, r)
		} else {
			accessToken := r.Header.Get("access-token")
			role := jwtUtils.GetRoleFromJWT(accessToken)
			if access.CheckRoleAccess(r.RequestURI, r.Method, role) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, objects.ForbiddenErrorString, http.StatusForbidden)
			}
		}
	})
}
