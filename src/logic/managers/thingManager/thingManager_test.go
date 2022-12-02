package thingManager

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/roomRepo"
	"src/db/studentRepo"
	"src/db/thingRepo"
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	"testing"
	"time"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

type TestThingManager struct{}

func Test_ThingManager(t *testing.T) {
	testgroup.RunSerially(t, &TestThingManager{})
}

func (*TestThingManager) TestThingManager_AddNewThingPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = mother.DefaultMarkNumber + 4
		ThingType  = mother.DefaultThingType
	)

	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()
	N := 3
	allThings := thingObjectMother.CreateDefaultThings(N)

	firstAllThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllThingsRows).WillReturnError(nil)
	mock.ExpectExec("INSERT INTO").WithArgs(MarkNumber, ThingType).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	_ = manager.AddNewThing(MarkNumber, ThingType)

	// Assert

	//tests.AssertErrors(t, execErr, nil)
	//tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_AddNewThingNegativeThingExist(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = mother.DefaultMarkNumber + 1
		ThingType  = mother.DefaultThingType
	)

	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()
	N := 3
	allThings := thingObjectMother.CreateDefaultThings(N)

	firstAllThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllThingsRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.AddNewThing(MarkNumber, ThingType)

	// Assert
	tests.AssertErrors(t, execErr, thingController.ThingAlreadyExistErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_AddNewThingNegativeBadParams(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = mother.DefaultMarkNumber + 4
		ThingType  = objects.EmptyString
	)

	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.AddNewThing(MarkNumber, ThingType)

	// Assert
	tests.AssertErrors(t, execErr, thingController.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_GetFreeThings(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	thingObjectMother := mother.ThingRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	N := 1
	allThings := thingObjectMother.CreateDefaultThings(N)
	allThings[0].SetOwnerID(objects.None)
	allRooms := roomObjectMother.CreateDefaultRooms(N)

	firstAllThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllThingsRows).WillReturnError(nil)
	firstRoomRows := roomObjectMother.CreateRows(allRooms)
	mock.ExpectQuery("SELECT").WillReturnRows(firstRoomRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.GetFreeThings()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_GetFullThingInfo(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	thingObjectMother := mother.ThingRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	N := 1
	allThings := thingObjectMother.CreateDefaultThings(N)
	allRooms := roomObjectMother.CreateDefaultRooms(N)

	firstAllThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllThingsRows).WillReturnError(nil)
	firstRoomRows := roomObjectMother.CreateRows(allRooms)
	mock.ExpectQuery("SELECT").WillReturnRows(firstRoomRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.GetFullThingInfo()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_GetStudentThingsPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	studentNumber := mother.DefaultStudentNumber + "1"
	StudentID := 1

	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	roomsN := 1

	allThings := thingObjectMother.CreateDefaultThings(thingsN)
	allRooms := roomObjectMother.CreateDefaultRooms(roomsN)
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	allThings[0].SetOwnerID(StudentID)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	secondAllStudentRows := studentObjectMother.CreateRows(allStudents[:1])

	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)
	mock.ExpectQuery("SELECT").WillReturnRows(secondAllStudentRows).WillReturnError(nil)

	firstStudentThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).
		WillReturnRows(firstStudentThingsRows)

	firstRoomRows := roomObjectMother.CreateRows(allRooms)
	mock.ExpectQuery("SELECT").WillReturnRows(firstRoomRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.GetStudentThings(studentNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_GetStudentThingsNegativeStudentNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	studentNumber := mother.DefaultStudentNumber + "6"

	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()
	studentsN := 4

	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.GetStudentThings(studentNumber)

	// Assert
	tests.AssertErrors(t, execErr, studentController.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_GetStudentThingsNegativeBadStudNumber(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	studentNumber := objects.EmptyString

	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.GetStudentThings(studentNumber)

	// Assert
	tests.AssertErrors(t, execErr, studentController.BadParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_TransferThingPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = 1
		ThingID    = 1
		SrcRoomID  = 1
		DstRoomID  = 2
	)
	thingObjectMother := mother.ThingRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	N := 1
	allThings := thingObjectMother.CreateDefaultThings(N)
	allRooms := roomObjectMother.CreateDefaultRooms(N)

	firstThingRow := thingObjectMother.CreateRowForID(ThingID)
	mock.ExpectQuery("SELECT").WillReturnRows(firstThingRow).WillReturnError(nil)

	firstRoomRows := roomObjectMother.CreateRows(allRooms)
	mock.ExpectQuery("SELECT").WillReturnRows(firstRoomRows).WillReturnError(nil)

	secondThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WillReturnRows(secondThingsRows).WillReturnError(nil)

	lastThingsRows := thingObjectMother.CreateRows(allThings)
	mock.ExpectQuery("SELECT").WillReturnRows(lastThingsRows).WillReturnError(nil)

	mock.ExpectExec("INSERT").WithArgs(ThingID, SrcRoomID, DstRoomID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.TransferThing(MarkNumber, DstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_TransferThingNegativeThingNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = 1
		DstRoomID  = 2
	)
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.TransferThing(MarkNumber, DstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, thingController.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_TransferThingNegativeRoomNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = 1
		ThingID    = 1
		DstRoomID  = 2
	)
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()
	firstThingRow := thingObjectMother.CreateRowForID(ThingID)
	mock.ExpectQuery("SELECT").WillReturnRows(firstThingRow).WillReturnError(nil)
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.TransferThing(MarkNumber, DstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, roomController.RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_TransferThingBadMarkNumber(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = -1
		DstRoomID  = 2
	)
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.TransferThing(MarkNumber, DstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, thingController.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestThingManager) TestThingManager_TransferThingBadRoom(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber = 1
		DstRoomID  = -2
	)
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := thingObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := ThingManager{studentController: studentC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.TransferThing(MarkNumber, DstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, thingController.BadDstRoomErr)
	tests.AssertMocks(t, mock)
}
