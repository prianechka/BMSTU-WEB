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
	"src/utils/logger"
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
// @Tags students
// @Security JWT-Token
// @Produce json
// @Param page query int false "Page param for pagination"
// @Param size query int false "Size param for pagination"
// @Success 200 {object} objects.StudentResponseDTO
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students [GET]
func (sh *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	page, size := utils.GetPageAndSizeFromQuery(r)

	allStudents, err := sh.manager.ViewAllStudents(page, size)
	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
		resultStudents := objects.CreateStudentResponse(allStudents)
		bytes, _ := json.Marshal(&resultStudents)
		_, _ = w.Write(bytes)
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		utils.SendResponseWithInternalErr(w)
	}
	logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
}

// ChangeStudentGroup
// @Summary Change student group for student
// @Description Change in database group information about student.
// @Produce json
// @Tags students
// @Security JWT-Token
// @Param  stud-number path string true "Student Number"
// @Param  new-group  body models.ChangeStudentGroupRequestMessage true "New student group"
// @Success 200 {object} models.ShortResponseMessage "Данные о студенте успешно обновлены!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number} [PUT]
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
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
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
	logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
}

// TransferStudent
// @Summary Settle/evic student in dormitory
// @Description Settle/evic student in certain room.
// @Produce json
// @Tags students
// @Security JWT-Token
// @Param  requestParams body models.StudentLiveActsRequestMessage true "Параметры запроса. Если roomID == 0, то студент выселяется."
// @Success 200 {object} models.ShortResponseMessage "Данные о студенте успешно обновлены!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден" | "Комната не найдена"
// @Failure 422 {object} models.ShortResponseMessage "Студент уже живёт в другой комнате!" | "Студент уже нигде не живёт!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/student-live-acts [POST]
func (sh *StudentHandler) TransferStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	var params models.StudentLiveActsRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
		return
	}

	switch params.RoomID {
	case objects.Null:
		err = sh.manager.EvicStudent(params.StudentNumber)
	default:
		err = sh.manager.SettleStudent(params.StudentNumber, params.RoomID)
	}

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.StudentChangeOKString
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
		statusCode = http.StatusBadRequest
		handleMessage = objects.RoomNotFoundErrorString
	case appErrors.StudentAlreadyLiveErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.SettleStudentErrorString
	case appErrors.StudentNotLivingErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.EvicStudentErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
}

// TransferThingFromToStudents
// @Summary Action with things by students
// @Description Give thing to student without changing its location.
// @Produce json
// @Tags students
// @Security JWT-Token
// @Param TransferParams body models.StudentThingsActsRequestMessage true "Параметры запроса. У поля status 2 значения: give(выдать) и return (забрать)"
// @Success 200 {object} models.ShortResponseMessage "Вещь передана!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден" | "Вещь не найдена!"
// @Failure 422 {object} models.ShortResponseMessage "Вещь уже у другого студента!" | "Вещь и так не у студента."
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/student-things-acts [POST]
func (sh *StudentHandler) TransferThingFromToStudents(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	var params models.StudentThingsActsRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
		return
	}
	switch params.Status {
	case Give:
		err = sh.manager.GiveStudentThing(params.StudentNumber, params.MarkNumber)
	case Return:
		err = sh.manager.ReturnStudentThing(params.StudentNumber, params.MarkNumber)
	default:
		err = appErrors.WrongRequestParamsErr
	}

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.GiveThingOK
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.WrongRequestParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	case appErrors.ThingNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.ThingNotFound
	case appErrors.ThingHasOwnerErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.GiveThingErrorString
	case appErrors.StudentIsNotOwnerErr:
		statusCode = http.StatusUnprocessableEntity
		handleMessage = objects.ReturnThingErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}
	utils.SendShortResponse(w, statusCode, handleMessage)
	sh.logger.Infof("Request: method - %s,  url - %s, Result: status_code = %d, text = %s, err = %v",
		r.Method, r.URL.Path, statusCode, handleMessage, err)
}

// AddNewStudent
// @Summary Add new student to base
// @Description Add new student in user and student base.
// @Tags students
// @Param  user-params body models.AddNewStudentRequestMessage true "student base information"
// @Success 200 {object} models.ShortResponseMessage "Операция успешно проведена!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 422 {object} models.ShortResponseMessage "Студент с таким же студенческим билетом уже существует!"
// @Failure 422 {object} models.ShortResponseMessage "Пользователь с таким логином уже существует!"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Security JWT-Token
// @Produce json
// @Router /api/v1/students [POST]
func (sh *StudentHandler) AddNewStudent(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	var params models.AddNewStudentRequestMessage

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
		err = readBodyErr
		utils.SendResponseWithInternalErr(w)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		statusCode = http.StatusBadRequest
		handleMessage = objects.WrongParamsErrorString
		utils.SendShortResponse(w, statusCode, handleMessage)
		logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
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
	logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
}

// ViewStudentInfo
// @Summary View information about student
// @Description View full information about student.
// @Tags students
// @Security JWT-Token
// @Produce json
// @Param  stud-number path string true "Student Number"
// @Success 200 {object} models.ShortResponseMessage "Данные о студенте успешно обновлены!"
// @Failure 400 {object} models.ShortResponseMessage "Параметр не должен быть пустой" | "Параметр должен быть числом!"
// @Failure 403 {object} models.ShortResponseMessage "У вас нет достаточно прав!"
// @Failure 404 {object} models.ShortResponseMessage "Студент не найден"
// @Failure 500 {object} models.ShortResponseMessage "Проблемы на стороне сервера."
// @Router /api/v1/students/{stud-number} [GET]
func (sh *StudentHandler) ViewStudentInfo(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var handleMessage string
	var err error

	studentNumber, _ := mux.Vars(r)["stud-number"]

	studentInfo, err := sh.manager.ViewStudent(studentNumber)

	switch err {
	case nil:
		statusCode = http.StatusOK
		handleMessage = objects.AddOK
		resultStudents := models.CreateStudentFullInfoResponse(studentInfo)
		bytes, _ := json.Marshal(&resultStudents)
		_, _ = w.Write(bytes)
		return
	case appErrors.BadStudentParamsErr:
		statusCode = http.StatusBadRequest
		handleMessage = objects.EmptyParamsErrorString
	case appErrors.StudentNotFoundErr:
		statusCode = http.StatusNotFound
		handleMessage = objects.StudentNotFoundErrorString
	default:
		statusCode = http.StatusInternalServerError
		handleMessage = objects.InternalServerErrorString
	}

	utils.SendShortResponse(w, statusCode, handleMessage)
	logger.WriteInfoInLog(sh.logger, r, statusCode, handleMessage, err)
}
