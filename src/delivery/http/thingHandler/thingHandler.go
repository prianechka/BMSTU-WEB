package thingHandler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"src/delivery/http/models"
	"src/logic/managers/thingManager"
	"src/objects"
	"src/utils"
	appErrors "src/utils/error"
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

// ViewStudentThings
// @Summary Get student things
// @Description View full information about all current things of student.
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Success 200 {object} models.ThingFullInfo
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустым."
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number}/things/ [GET]
func (th *ThingHandler) ViewStudentThings(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	studentNumber := r.URL.Query().Get("stud-number")

	allThings, err := th.manager.GetStudentThings(studentNumber)
	switch err {
	case nil:
		bytes, _ := json.Marshal(&allThings)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
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
	}

	if err != nil {
		utils.SendShortResponse(w, statusCode, handleMessage)
	}

	th.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}

// ViewFreeThings
// @Summary Get all free things in dormitory
// @Description View full information about free things (things without owner) .
// @Produce json
// @Success 200 {object} models.ThingFullInfo
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/things/free/ [GET]
func (th *ThingHandler) ViewFreeThings(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	allFreeThings, err := th.manager.GetFreeThings()
	switch err {
	case nil:
		bytes, _ := json.Marshal(&allFreeThings)
		_, _ = w.Write(bytes)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}

	th.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}

// ViewAllThings
// @Summary Get all things in dormitory
// @Description View full information about all things dormitory.
// @Produce json
// @Success 200 {object} models.ThingFullInfo
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/things/ [GET]
func (th *ThingHandler) ViewAllThings(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	allThings, err := th.manager.GetFullThingInfo()
	switch err {
	case nil:
		bytes, _ := json.Marshal(&allThings)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	th.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}

// TransferThingBetweenRooms
// @Summary Transfer thing to another room.
// @Description Transfer thing to another room.
// @Produce json
// @Param  mark-number path int true "Thing mark number"
// @Param  room-id body  models.TransferThingRequestMessage true "Dst room in which thing will be transferred."
// @Success 200 {object} models.ShortResponseMessage "Вещь успешно перемещена!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 404 {object} models.ShortResponseMessage "Вещь не найдена" | "Комната не найдена"
// @Failure 422 {object} models.ShortResponseMessage "Вещь уже находится в этой комнате!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/things/{mark-number}/ [PATCH]
func (th *ThingHandler) TransferThingBetweenRooms(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer th.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	var params models.TransferThingRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	markNumber, getMarkNumberErr := utils.GetIntParamByKey(r, "mark-number")

	if getMarkNumberErr != nil {
		err = getMarkNumberErr
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
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
}

// AddNewThing
// @Summary Add new thing
// @Description Add new thing with params in base.
// @Produce json
// @Param params body models.AddNewThingRequestMessage true "body for buy service"
// @Success 200 {object} models.ShortResponseMessage "Операция успешно проведена!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 422 {object} models.ShortResponseMessage "Вещь с таким же уникальным номером уже есть в базе"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/things/ [POST]
func (th *ThingHandler) AddNewThing(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer th.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	var params models.AddNewThingRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
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
}
