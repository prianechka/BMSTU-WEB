package utils

import (
	"net/http"
	"strconv"
)

func GetIntParamByKey(r *http.Request, key string) (int, error) {
	paramByString := r.URL.Query().Get(key)

	return strconv.Atoi(paramByString)
}
