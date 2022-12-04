package studentHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

func CreateNewStudentHandler(logger *logrus.Entry, manager studentManager.StudentManager) *StudentHandler {
	return &StudentHandler{
		logger:  logger,
		manager: manager,
	}
}

// GetAllStudents
// @Summary Get all students in dormitory
// @Description View full information about students have lived in dormitory.
// @Produce json
// @Success 200 {object} objects.StudentResponseDTO
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/ [GET]
func (sh *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	allStudents, err := sh.manager.ViewAllStudents()
	resultStudents := objects.CreateStudentResponse(allStudents)
	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
		bytes, _ := json.Marshal(&resultStudents)
		_, _ = w.Write(bytes)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}

// ChangeStudentGroup
// @Summary Change student group for student
// @Description Change in database group information about student.
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Param  new-group  body models.ChangeStudentGroupRequestMessage true "New student group"
// @Success 200 {object} models.ShortResponseMessage "Данные о студенте успешно обновлены!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number}/ [PUT]
func (sh *StudentHandler) ChangeStudentGroup(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	var params models.ChangeStudentGroupRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if err != nil {
		err = readBodyErr
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
			r.Method, r.URL.Path, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
			r.Method, r.URL.Path, statusCode, handleMessage, err)
		return
	}

	studentNumber, _ := mux.Vars(r)["stud-number"]

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
	sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}

// SettleStudent
// @Summary Settle student in room
// @Description Settle student in certain room.
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Param  room-id     body int true "New student room ID"
// @Success 200 {object} models.ShortResponseMessage "Данные о студенте успешно обновлены!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден" | "Комната не найдена"
// @Failure 422 {object} models.ShortResponseMessage "Студент уже живёт в другой комнате!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number}/rooms/ [POST]
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

	studentNumber := r.URL.Query().Get("stud-number")

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

// EvicStudent
// @Summary Evic student from current room
// @Description Settle student in certain room.
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Success 200 {object} models.ShortResponseMessage "Данные о студенте успешно обновлены!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден"
// @Failure 422 {object} models.ShortResponseMessage "Студент уже нигде не живёт!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number}/rooms/ [DELETE]
func (sh *StudentHandler) EvicStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	defer sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)

	studentNumber := r.URL.Query().Get("stud-number")

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
	case appErrors.StudentNotLivingErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.EvicStudentErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

// GiveStudentThing
// @Summary Give thing to student
// @Description Give thing to student without changing its location.
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Param  mark-number path int true "Mark number of thing"
// @Success 200 {object} models.ShortResponseMessage "Вещь передана!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден" | "Вещь не найдена!"
// @Failure 422 {object} models.ShortResponseMessage "Вещь уже у другого студента!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number}/things/{thing-id}/ [POST]
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
		handleMessage = objects.GiveThingOK
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
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.GiveThingErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

// ReturnThingFromStudent
// @Summary Return thing from student.
// @Description Return thing from student without changing thing location.
// @Param  stud-number path string true "Student Number"
// @Param  mark-number path int true "Mark number of thing"
// @Success 200 {object} models.ShortResponseMessage "Вещь передана!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден" | "Вещь не найдена!"
// @Failure 422 {object} models.ShortResponseMessage "Вещь и так была не у студента!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number}/things/{thing-id}/ [DELETE]
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

	studentNumber := r.URL.Query().Get("stud-number")

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
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.ReturnThingErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
}

// AddNewStudent
// @Summary Add new student to base
// @Description Add new student in user and student base.
// @Param  user-params body models.AddNewStudentRequestMessage true "student base information"
// @Success 200 {object} models.ShortResponseMessage "Операция успешно проведена!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 422 {object} models.ShortResponseMessage "Студент с таким же студенческим билетом уже существует!"
// @Failure 422 {object} models.ShortResponseMessage "Пользователь с таким логином уже существует!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Produce json
// @Success 200 {object} objects.Student
// @Failure 500 {object} models.ShortResponseMessage "internal server error"
// @Router /api/v1/students/ [POST]
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
		handleMessage = objects.AddOK
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
