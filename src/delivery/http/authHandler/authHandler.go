package authHandler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"src/delivery/http/models"
	"src/logic/managers/appManager"
	"src/logic/managers/authManager"
	"src/objects"
	"src/utils"
	appErrors "src/utils/error"
)

type AuthHandler struct {
	logger      *logrus.Entry
	AuthManager authManager.AuthManager
	AppManager  appManager.AppManager
}

func CreateNewAuthHandler(logger *logrus.Entry, am authManager.AuthManager, ap appManager.AppManager) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		AuthManager: am,
		AppManager:  ap,
	}
}

func (h *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	var authParams models.AuthRequestMessage

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	unmarshalError := json.Unmarshal(body, &authParams)
	if unmarshalError != nil {
		http.Error(w, "invalid body params", http.StatusBadRequest)
		return
	}

	newRole, err := h.AuthManager.TryToAuth(authParams.Login, authParams.Password)
	switch err {
	case nil:
		if newRole == objects.NonAuth {
			h.AppManager.FoldState()
		} else {
			h.AppManager.SetNewState(authParams.Login, newRole)
		}
	case appErrors.UserNotFoundErr:
		statusCode = http.StatusForbidden
		handleMessage = objects.LoginErrorString
	case appErrors.PasswordNotEqualErr:
		statusCode = http.StatusForbidden
		handleMessage = objects.PasswordErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	if newRole == objects.NonAuth {
		h.AppManager.FoldState()
	} else {
		h.AppManager.SetNewState(authParams.Login, newRole)
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
	h.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s",
		r.Method, r.URL.Path, statusCode, handleMessage)
}
