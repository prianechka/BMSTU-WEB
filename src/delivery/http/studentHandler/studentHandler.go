package studentHandler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"src/delivery/http/models"
	"src/logic/managers/studentManager"
	"src/objects"
	"src/utils"
	appErrors "src/utils/error"
)

type StudentHandler struct {
	logger  *logrus.Entry
	manager studentManager.StudentManager
}

func (sh *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	allStudents, err := sh.manager.ViewAllStudents()
	switch err {
	case nil:
		bytes, _ := json.Marshal(&allStudents)
		_, _ = w.Write(bytes)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
}

func (sh *StudentHandler) ChangeStudentGroup(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	var params models.ChangeStudentGroupRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if err != nil {
		err = readBodyErr
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	studentNumber := r.URL.Query().Get("studnumber")

	err = sh.manager.ChangeStudentGroup(studentNumber, params.NewGroup)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
}

func (sh *StudentHandler) SettleStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	var params models.SettleInRoomRequestMessage

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

	studentNumber := r.URL.Query().Get("studnumber")

	err = sh.manager.SettleStudent(studentNumber, params.RoomID)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.RoomNotFoundErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.RoomNotFoundErrorString
	case appErrors.StudentAlreadyLiveErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.SettleStudentErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
}

func (sh *StudentHandler) EvicStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	studentNumber := r.URL.Query().Get("studnumber")

	err = sh.manager.EvicStudent(studentNumber)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.RoomNotFoundErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.RoomNotFoundErrorString
	case appErrors.StudentNotLivingErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.EvicStudentErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

func (sh *StudentHandler) GiveStudentThing(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	markNumber, getMarkNumberErr := utils.GetIntParamByKey(r, "marknumber")

	if getMarkNumberErr != nil {
		err = getMarkNumberErr
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	studentNumber := r.URL.Query().Get("studnumber")

	err = sh.manager.GiveStudentThing(studentNumber, markNumber)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.ThingNotFound
	case appErrors.ThingHasOwnerErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.GiveThingErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

func (sh *StudentHandler) ReturnThingFromStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	markNumber, getMarkNumberErr := utils.GetIntParamByKey(r, "marknumber")

	if getMarkNumberErr != nil {
		err = getMarkNumberErr
		statusCode = http.StatusBadRequest
		handleMessage = objects.MustBeIntErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		return
	}

	studentNumber := r.URL.Query().Get("studnumber")

	err = sh.manager.ReturnStudentThing(studentNumber, markNumber)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.ThingNotFound
	case appErrors.StudentIsNotOwnerErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.ReturnThingErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

func (sh *StudentHandler) AddNewStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	var params models.AddNewStudentRequestMessage

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

	err = sh.manager.AddNewStudent(params.Name, params.Surname, params.Group, params.StudentNumber,
		params.Login, params.Password)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.BadUserParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.StudentAlreadyInBaseErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.StudentAlreadyExistErrorString
	case appErrors.LoginOccupedErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.UserAlreadyExistErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
}
