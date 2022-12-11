package thingHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"src/delivery/http/models"
	models2 "src/logic/managers/models"
	"src/logic/managers/thingManager"
	"src/objects"
	"src/utils"
	appErrors "src/utils/error"
	"src/utils/logger"
	"strconv"
)

type ThingHandler struct {
	logger  *logrus.Entry
	manager thingManager.ThingManager
}

func CreateNewThingHandler(logger *logrus.Entry, man thingManager.ThingManager) *ThingHandler {
	return &ThingHandler{
		logger:  logger,
		manager: man,
	}
}

// GetThings
// @Summary Get things in dormitory with params
// @Description Get full information about different things.
// @Param page query int false "Page param for pagination"
// @Param size query int false "Size param for pagination"
// @Param status query string false "Status defines mode: all things, free things or student things. Possible values: all, free, student"
// @Param stud-number query string false "Student number for searching in Student mode"
// @Tags things
// @Security JWT-Token
// @Produce json
// @Success 200 {object} models.ThingFullInfo
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметры указаны неверно"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера"
// @Router /api/v1/things [GET]
func (th *ThingHandler) GetThings(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var allThings []models2.ThingFullInfo
	var err error

	page, size := utils.GetPageAndSizeFromQuery(r)

	status := r.URL.Query().Get("status")
	switch status {
	case All:
		allThings, err = th.manager.GetFullThingInfo(page, size)
	case Free:
		allThings, err = th.manager.GetFreeThings(page, size)
	case OnStudent:
		studentNumber := r.URL.Query().Get("stud-number")
		allThings, err = th.manager.GetStudentThings(studentNumber, page, size)
	default:
		err = appErrors.WrongRequestParamsErr
	}

	switch err {
	case nil:
		resultThings := models.CreateThingFullInfoResponse(allThings)
		bytes, _ := json.Marshal(&resultThings)
		_, _ = w.Write(bytes)
		return
	case appErrors.BadThingParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.WrongRequestParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.RoomNotFoundErr:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
}

// GetThing
// @Summary Get thing info
// @Description Get full information about thing by mark-number.
// @Tags things
// @Security JWT-Token
// @Produce json
// @Param  mark-number path string true "Mark number for thing"
// @Success 200 {object} models.ThingFullInfo
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметры указаны неверно"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Вещь не найдена"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера"
// @Router /api/v1/things/{mark-number} [GET]
func (th *ThingHandler) GetThing(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	markNumberString, _ := mux.Vars(r)["mark-number"]

	markNumber, atoiErr := strconv.Atoi(markNumberString)
	if atoiErr != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, atoiErr)
		return
	}

	thingInfo, err := th.manager.GetThingInfo(markNumber)

	switch err {
	case nil:
		resultThings := models.CreateThingInfoResponse(thingInfo)
		bytes, _ := json.Marshal(&resultThings)
		_, _ = w.Write(bytes)
		return
	case appErrors.BadThingParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.ThingNotFound
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
}

// TransferThingBetweenRooms
// @Summary Transfer thing to another room.
// @Description Transfer thing to another room.
// @Produce json
// @Tags things
// @Security JWT-Token
// @Param  mark-number path int true "Thing mark number"
// @Param  room-id body  models.TransferThingRequestMessage true "Dst room in which thing will be transferred."
// @Success 200 {object} models.ShortResponseMessage "Вещь успешно перемещена!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Вещь не найдена" | "Комната не найдена"
// @Failure 422 {object} models.ShortResponseMessage "Вещь уже находится в этой комнате!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/things/{mark-number} [PATCH]
func (th *ThingHandler) TransferThingBetweenRooms(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	var params models.TransferThingRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
		return
	}

	markNumberString, _ := mux.Vars(r)["mark-number"]

	markNumber, atoiErr := strconv.Atoi(markNumberString)
	if atoiErr != nil {
		err = atoiErr
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
		return
	}

	err = th.manager.TransferThing(markNumber, params.NewRoomID)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.TransferThingOK
	case appErrors.BadThingParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.ThingNotFound
	case appErrors.BadDstRoomErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.ThingAlreadyInRoomErrorString
	case appErrors.RoomNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.RoomNotFoundErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
}

// AddNewThing
// @Summary Add new thing
// @Description Add new thing with params in base.
// @Produce json
// @Tags things
// @Security JWT-Token
// @Param params body models.AddNewThingRequestMessage true "body for buy service"
// @Success 200 {object} models.ShortResponseMessage "Операция успешно проведена!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 422 {object} models.ShortResponseMessage "Вещь с таким же уникальным номером уже есть в базе"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/things [POST]
func (th *ThingHandler) AddNewThing(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	var params models.AddNewThingRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
		return
	}

	markNumber := params.MarkNumber

	err = th.manager.AddNewThing(markNumber, params.ThingType)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
	case appErrors.ThingAlreadyExistErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.UniqueMarkNumberErrorString
	case appErrors.BadThingParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
}

// ViewThingHistory
// @Summary View history of thing owners
// @Description View history of thing owners
// @Tags students
// @Security JWT-Token
// @Produce json
// @Param  mark-number path int true "Маркировочный номер"
// @Param  status query string true "Параметр того, как выводить историю: текущего владельца(current) или общую историю (all)"
// @Success 200 {object} models.ThingOwnerHistoryResponseMessage
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Вещь не найдена"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Failure 501 {object} models.ShortResponseMessage "Пока функционал на стадии реализации"
// @Router /api/v1/student-things-acts/{mark-number} [GET]
func (th *ThingHandler) ViewThingHistory(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error
	var studentNumber string

	status := r.URL.Query().Get("status")
	switch status {
	case All:
		err = appErrors.NotImplementedErr
	case Current:
		markNumber, getMarkNumberErr := utils.GetMarkNumberFromPath(r)
		if getMarkNumberErr == nil {
			studentNumber, err = th.manager.GetOwner(markNumber)
		} else {
			err = appErrors.WrongRequestParamsErr
		}
	default:
		err = appErrors.WrongRequestParamsErr
	}

	switch err {
	case nil:
		result := models.ThingOwnerHistoryResponseMessage{OwnerStudentNumber: studentNumber}
		bytes, _ := json.Marshal(&result)
		_, _ = w.Write(bytes)
		return
	case appErrors.NotImplementedErr:
		statusCode = http.StatusNotImplemented
		handleMessage = objects.NotImplementedErrorString
	case appErrors.WrongRequestParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
	case appErrors.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.RoomNotFoundErr:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(th.logger, r, statusCode, handleMessage, err)
}
