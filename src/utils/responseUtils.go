package utils

import (
	"encoding/json"
	"net/http"
	"src/delivery/http/models"
)

func SendShortResponse(w http.ResponseWriter, code int, comment string) {
	var msg = models.ShortResponseMessage{Comment: comment}
	result, err := json.Marshal(msg)
	if err == nil {
		w.WriteHeader(code)
		_, err = w.Write(result)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
