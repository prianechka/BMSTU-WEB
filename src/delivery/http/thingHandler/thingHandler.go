package thingHandler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"src/delivery/http/models"
	roomErr "src/logic/controllers/roomController"
	studErr "src/logic/controllers/studentController"
	thingErr "src/logic/controllers/thingController"
	"src/logic/managers/thingManager"
	"src/objects"
	"src/utils"
	"strconv"
)

type ThingHandler struct {
	logger  *logrus.Entry
	manager thingManager.ThingManager
}

func (th *ThingHandler) ViewStudentThings(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string

	studentNumber := r.URL.Query().Get("studnumber")

	allThings, err := th.manager.GetStudentThings(studentNumber)
	switch err {
	case nil:
		bytes, _ := json.Marshal(&allThings)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
	case studErr.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case studErr.BadParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case roomErr.RoomNotFoundErr:
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

func (th *ThingHandler) TransferThingBetweenRooms(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer th.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	markNumber, getMarkNumberErr := utils.GetIntParamByKey(r, "marknumber")

	if getMarkNumberErr != nil {
		err = getMarkNumberErr
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	roomID, getRoomIDErr := utils.GetIntParamByKey(r, "roomID")

	if getRoomIDErr != nil {
		err = getRoomIDErr
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	err = th.manager.TransferThing(markNumber, roomID)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.TransferThingOK
	case thingErr.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.ThingNotFound
	case thingErr.BadDstRoomErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
	case thingErr.BadSrcRoomErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.ThingAlreadyInRoomErrorString
	case roomErr.RoomNotFoundErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.RoomNotFoundErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

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

	markNumber, castErr := strconv.Atoi(params.MarkNumber)
	if castErr != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	err = th.manager.AddNewThing(markNumber, params.ThingType)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case thingErr.ThingAlreadyExistErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.UniqueMarkNumberErrorString
	case studErr.BadParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
}
