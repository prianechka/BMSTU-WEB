package roomHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"src/logic/managers/roomManager"
	"src/objects"
	"src/utils"
	appErrors "src/utils/error"
	"src/utils/logger"
	"strconv"
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
// @Tags rooms
// @Security JWT-Token
// @param access-token header string true "JWT Token"
// @Produce json
// @Param page query int false "Page param for pagination"
// @Param size query int false "Size param for pagination"
// @Success 200 {object} objects.RoomResponseDTO
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/rooms [GET]
func (rh *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	page, size := utils.GetPageAndSizeFromQuery(r)

	allRooms, err := rh.manager.GetAllRooms(page, size)
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

// GetRoom
// @Summary Get room information
// @Description View information about room in dormitory.
// @Tags rooms
// @Security JWT-Token
// @param access-token header string true "JWT Token"
// @Produce json
// @Param  room-id path int true "Room id"
// @Success 200 {object} objects.RoomResponseDTO
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Комната не найдена!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/rooms/{room-id} [GET]
func (rh *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	roomIDString, _ := mux.Vars(r)["room-id"]

	roomID, atoiErr := strconv.Atoi(roomIDString)
	if atoiErr != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(rh.logger, r, statusCode, handleMessage, atoiErr)
		return
	}

	room, err := rh.manager.GetRoom(roomID)

	resultRoomInfo := objects.CreateRoomResponse(room)
	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
		bytes, _ := json.Marshal(&resultRoomInfo)
		_, _ = w.Write(bytes)
		return
	case appErrors.RoomNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.RoomNotFoundErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(rh.logger, r, statusCode, handleMessage, err)
}
