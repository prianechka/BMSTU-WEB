package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func WriteInfoInLog(log *logrus.Entry, r *http.Request, statusCode int, handleMessage string, err error) {
	log.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}
