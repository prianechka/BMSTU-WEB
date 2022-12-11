package utils

import (
	"github.com/gorilla/mux"
	"net/http"
	"src/objects"
	appErrors "src/utils/error"
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

func GetMarkNumberFromPath(r *http.Request) (markNumber int, err error) {
	markNumberString, isFound := mux.Vars(r)["mark-number"]
	if isFound {
		markNumber, err = strconv.Atoi(markNumberString)
	} else {
		err = appErrors.ThingNotFoundErr
	}
	return markNumber, err
}
