package utils

import (
	"net/http"
	"src/objects"
	"strconv"
)

func GetIntParamByKey(r *http.Request, key string) (int, error) {
	paramByString := r.URL.Query().Get(key)

	return strconv.Atoi(paramByString)
}

func GetPageAndSizeFromQuery(r *http.Request) (page int, size int) {
	pageFromQuery, pageErr := GetIntParamByKey(r, "page")
	if pageErr != nil {
		page = objects.DefaultPage
	} else {
		page = pageFromQuery
	}

	sizeFromQuery, sizeErr := GetIntParamByKey(r, "size")
	if sizeErr != nil {
		size = objects.DefaultPageSize
	} else {
		size = sizeFromQuery
	}

	return page, size
}
