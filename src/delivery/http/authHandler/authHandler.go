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
	"src/utils/jwtUtils"
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

// Authorize
// @Summary Try to authorize in system
// @Description Try to authorize in system. JWT-Token send with success
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Success 200 {object} models.ShortResponseMessage "Операция прошла успешно!"
// @Failure 403 {object} models.ShortResponseMessage "Пароль введен неверно!"
// @Failure 404 {object} models.ShortResponseMessage "Пользователь не не найден"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/auth/ [GET]
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
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
		jwtToken := jwtUtils.CreateJWTToken(authParams.Login, newRole)
		w.Header().Set("access-token", jwtToken)
	case appErrors.UserNotFoundErr:
		statusCode = http.StatusNotFound
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
