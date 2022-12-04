package roomHandler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"src/logic/managers/roomManager"
	"src/objects"
	"src/utils"
	"src/utils/logger"
)

type RoomHandler struct {
	logger  *logrus.Entry
	manager roomManager.RoomManager
}

func CreateNewRoomHandler(logger *logrus.Entry, manager roomManager.RoomManager) *RoomHandler {
	return &RoomHandler{
		logger:  logger,
		manager: manager,
	}
}

// GetAllRooms
// @Summary Get all rooms in dormitory
// @Description View full information about rooms in dormitory.
// @Produce json
// @Success 200 {object} objects.ThingResponseDTO
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/rooms/ [GET]
func (rh *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	allRooms, err := rh.manager.GetAllRooms()
	resultRooms := objects.CreateRoomResponseArr(allRooms)
	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
		bytes, _ := json.Marshal(&resultRooms)
		_, _ = w.Write(bytes)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	logger.WriteInfoInLog(rh.logger, r, statusCode, handleMessage, err)
}
