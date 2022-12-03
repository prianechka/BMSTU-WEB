package studentController

import (
	"database/sql"
	"fmt"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/studentRepo"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	appErrors "src/utils/error"
	"testing"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

type TestStudentController struct{}

func Test_StudentController(t *testing.T) {
	testgroup.RunSerially(t, &TestStudentController{})
}

func (*TestStudentController) TestStudentController_AddStudentPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, InsertID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_AddStudentNegativeAlreadyLive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, InsertID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentAlreadyInBaseErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_AddStudentNegativeBadID(t *testgroup.T) {
	// Arrange
	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + fmt.Sprintf("%d", 1)
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, -InsertID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadAccIDErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_AddStudentNegativeBadStudentGroup(t *testgroup.T) {
	// Arrange
	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = objects.EmptyString
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.AddStudent(Name, Surname, StudentGroup, StudentNumber, InsertID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_GetAllStudents(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows).WillReturnError(nil)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	resultStudents, execErr := controller.GetAllStudents()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, resultStudents, realStudents)
}

func (*TestStudentController) TestStudentController_GetStudentPositive(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	student, execErr := controller.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, student, realStudents[0])
}

func (*TestStudentController) TestStudentController_GetStudentNegativeNotFound(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	id := 1
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	_, execErr := controller.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_GetStudentNegativeBadID(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	id := -1

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	_, execErr := controller.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_GetStudentIDByNumberPositive(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	N := 1
	StudNumber := mother.DefaultStudentNumber + fmt.Sprintf("%d", 1)
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	studentID, execErr := controller.GetStudentIDByNumber(StudNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, studentID, ID)
}

func (*TestStudentController) TestStudentController_GetStudentIDByNumberNegative(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	StudNumber := mother.DefaultStudentNumber + fmt.Sprintf("%d", 6)
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	_, execErr := controller.GetStudentIDByNumber(StudNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_GetStudentRoomPositive(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	N := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	roomID, execErr := controller.GetStudentRoom(ID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, roomID, ID)
}

func (*TestStudentController) TestStudentController_GetStudentRoomNegative(t *testgroup.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 5
	N := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs().WillReturnError(nil).WillReturnRows(rows)
	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	_, execErr := controller.GetStudentRoom(ID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_GetStudentThingsPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	things, execErr := controller.GetStudentThings(ID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, things, realThings)
}

func (*TestStudentController) TestStudentController_GetStudentThingsNegative(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	ID := 1
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	_, execErr := controller.GetStudentThings(ID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_SettleStudentPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.SettleStudent(studentID, roomID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_SettleStudentNegativeStudentNotFound(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	roomID := 2

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.SettleStudent(studentID, roomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_SettleStudentNegativeStudentLiveNow(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1
	roomID := 2

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.SettleStudent(studentID, roomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentAlreadyLiveErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_EvicStudentPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.EvicStudent(studentID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_EvicStudentNegativeStudentNotFound(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.EvicStudent(studentID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_EvicStudentStudentDoesNotLive(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 1
	studentID := 1

	realStudents := studentObjectMother.CreateDefaultStudents(N)
	realStudents[0].SetRoomID(objects.NotLiving)
	studentRows := studentObjectMother.CreateRows(realStudents)

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(nil).WillReturnRows(studentRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.EvicStudent(studentID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotLivingErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_ChangeStudentGroupPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.ChangeStudentGroup(studentID, newGroup)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_ChangeStudentGroupNegative(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	newGroup := "iu7-86"
	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.ChangeStudentGroup(studentID, newGroup)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_TransferThingPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.TransferThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_TransferThingNegative(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	thingID := 2

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.TransferThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_ReturnThingPositive(t *testgroup.T) {
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

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.ReturnThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentController) TestStudentController_ReturnThingNegative(t *testgroup.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	studentID := 1
	thingID := 2

	mock.ExpectQuery("SELECT").WithArgs(studentID).WillReturnError(sql.ErrNoRows)

	Repo := studentRepo.PgStudentRepo{Conn: db}
	controller := StudentController{Repo: &Repo}

	// Act
	execErr := controller.ReturnThing(studentID, thingID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}
