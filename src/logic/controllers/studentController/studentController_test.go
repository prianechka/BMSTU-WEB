package studentController

import (
	"database/sql"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/studentRepo"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	"testing"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

func TestStudentController_AddStudentPositive(t *testing.T) {
	// Arrange
	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + "6"
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 5
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)
	mock.ExpectExec("INSERT INTO").WithArgs(Name, Surname, StudentGroup, StudentNumber, InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, InsertID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestStudentController_AddStudentNegativeAlreadyLive(t *testing.T) {
	// Arrange
	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + fmt.Sprintf("%d", 1)
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 5
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, InsertID)

	// Assert
	tests.AssertErrors(t, execErr, StudentAlreadyInBaseErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_AddStudentNegativeBadID(t *testing.T) {
	// Arrange
	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + fmt.Sprintf("%d", 1)
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, -InsertID)

	// Assert
	tests.AssertErrors(t, execErr, BadAccIDErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_AddStudentNegativeBadStudentGroup(t *testing.T) {
	// Arrange
	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = objects.EmptyString
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, InsertID)

	// Assert
	tests.AssertErrors(t, execErr, BadParamsErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_GetAllStudents(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows).WillReturnError(nil)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	resultStudents, execErr := controller.GetAllStudents()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, resultStudents, realStudents)
}

func TestStudentController_GetStudentPositive(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	student, execErr := controller.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, student, realStudents[0])
}

func TestStudentController_GetStudentNegativeNotFound(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	id := 1
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	_, execErr := controller.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_GetStudentNegativeBadID(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	id := -1

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	_, execErr := controller.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, BadParamsErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_GetStudentIDByNumberPositive(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	N := 1
	StudNumber := mother.DefaultStudentNumber + fmt.Sprintf("%d", 1)
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	studentID, execErr := controller.GetStudentIDByNumber(StudNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, studentID, ID)
}

func TestStudentController_GetStudentIDByNumberNegative(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	StudNumber := mother.DefaultStudentNumber + fmt.Sprintf("%d", 6)
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	_, execErr := controller.GetStudentIDByNumber(StudNumber)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_GetStudentRoomPositive(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	N := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	roomID, execErr := controller.GetStudentRoom(ID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, roomID, ID)
}

func TestStudentController_GetStudentRoomNegative(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 5
	N := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)
	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	_, execErr := controller.GetStudentRoom(ID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_GetStudentThingsPositive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 3
	ID := 1
	realThings := thingObjectMother.CreateDefaultThings(N)
	thingRows := thingObjectMother.CreateRows(realThings)

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(studentRows)
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(thingRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	things, execErr := controller.GetStudentThings(ID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, things, realThings)
}

func TestStudentController_GetStudentThingsNegative(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	ID := 1
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	_, execErr := controller.GetStudentThings(ID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_SettleStudentPositive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	roomID := 2

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	realStudents[0].SetRoomID(objects.NotLiving)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)
	mock.ExpectExec("INSERT").WithArgs(studentID, roomID, objects.Get).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.SettleStudent(studentID, roomID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestStudentController_SettleStudentNegativeStudentNotFound(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	roomID := 2

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.SettleStudent(studentID, roomID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_SettleStudentNegativeStudentLiveNow(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	roomID := 2

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.SettleStudent(studentID, roomID)

	// Assert
	tests.AssertErrors(t, execErr, StudentAlreadyLiveErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_EvicStudentPositive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	roomID := 1

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)
	mock.ExpectExec("INSERT").WithArgs(studentID, roomID, objects.Ret).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.EvicStudent(studentID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestStudentController_EvicStudentNegativeStudentNotFound(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.EvicStudent(studentID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_EvicStudentStudentDoesNotLive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	realStudents[0].SetRoomID(objects.NotLiving)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.EvicStudent(studentID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotLivingErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_ChangeStudentGroupPositive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	newGroup := "iu7-86"

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	realStudents[0].SetRoomID(objects.NotLiving)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)
	mock.ExpectExec("UPDATE").WithArgs(realStudents[0].GetName(), realStudents[0].GetSurname(),
		newGroup, realStudents[0].GetStudentNumber(), realStudents[0].GetID()).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.ChangeStudentGroup(studentID, newGroup)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestStudentController_ChangeStudentGroupNegative(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	newGroup := "iu7-86"
	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.ChangeStudentGroup(studentID, newGroup)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_TransferThingPositive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	thingID := 2

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)
	mock.ExpectExec("INSERT").WithArgs(studentID, thingID, objects.Get).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.TransferThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestStudentController_TransferThingNegative(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	thingID := 2

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.TransferThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestStudentController_ReturnThingPositive(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	thingID := 2

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)
	mock.ExpectExec("INSERT").WithArgs(studentID, thingID, objects.Ret).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.ReturnThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestStudentController_ReturnThingNegative(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	thingID := 2

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{repo: &repo}

	// Act
	execErr := controller.ReturnThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}
