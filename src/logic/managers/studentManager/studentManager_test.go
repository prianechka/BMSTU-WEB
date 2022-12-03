package studentManager

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/roomRepo"
	"src/db/studentRepo"
	"src/db/thingRepo"
	"src/db/userRepo"
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/logic/controllers/userController"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	appErrors "src/utils/error"
	"testing"
	"time"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

type TestStudentManager struct{}

func Test_StudentManager(t *testing.T) {
	testgroup.RunSerially(t, &TestStudentManager{})
}

func (*TestStudentManager) TestStudentManager_AddNewStudentPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		Login         = mother.DefaultLogin + "123"
		Password      = mother.DefaultPassword
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	userObjectMother := mother.UserRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	N := 4
	allStudents := studentObjectMother.CreateDefaultStudents(N)

	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WithArgs(Login, Password, objects.StudentRole).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))
	thirdUsersRows := userObjectMother.CreateRowForID(InsertID)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(thirdUsersRows)

	allStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(allStudentRows).WillReturnError(nil)
	mock.ExpectExec("INSERT INTO").WithArgs(Name, Surname, StudentGroup, StudentNumber, InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.AddNewStudent(Name, Surname, StudentGroup, StudentNumber, Login, Password)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_AddNewStudentNegativeLoginOccuped(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		Login         = mother.DefaultLogin
		Password      = mother.DefaultPassword
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	userObjectMother := mother.UserRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(userObjectMother.CreateRowForID(InsertID))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.AddNewStudent(Name, Surname, StudentGroup, StudentNumber, Login, Password)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.LoginOccupedErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_AddNewStudentNegativeBadLoginParam(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		Login         = objects.EmptyString
		Password      = mother.DefaultPassword
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.AddNewStudent(Name, Surname, StudentGroup, StudentNumber, Login, Password)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadUserParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_AddNewStudentNegativeBadName(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		Login         = mother.DefaultLogin
		Password      = mother.DefaultPassword
		Name          = objects.EmptyString
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.AddNewStudent(Name, Surname, StudentGroup, StudentNumber, Login, Password)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ChangeStudentGroupPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		Name          = mother.DefaultStudentName
		Surname       = mother.DefaultStudentSurname
		StudentGroup  = mother.DefaultGroup + "12"
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	ID := 1
	N := 4
	allStudents := studentObjectMother.CreateDefaultStudents(N)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	secondStudentRow := studentObjectMother.CreateRows(allStudents[:1])
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(secondStudentRow)
	mock.ExpectExec("UPDATE").WithArgs(Name, Surname, StudentGroup, StudentNumber, ID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(int64(ID), RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ChangeStudentGroup(StudentNumber, StudentGroup)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ChangeStudentGroupNegativeBadParams(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentGroup  = objects.EmptyString
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ChangeStudentGroup(StudentNumber, StudentGroup)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ChangeStudentGroupNegativeStudentNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentGroup  = mother.DefaultGroup + "12"
		StudentNumber = mother.DefaultStudentNumber + "12"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	N := 4
	allStudents := studentObjectMother.CreateDefaultStudents(N)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ChangeStudentGroup(StudentNumber, StudentGroup)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_SettleStudentPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentID     = 1
		RoomID        = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	roomsN := 1
	realRooms := roomObjectMother.CreateDefaultRooms(roomsN)
	rows := roomObjectMother.CreateRows(realRooms)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)
	mock.ExpectQuery("SELECT").WithArgs(RoomID).WillReturnError(nil).WillReturnRows(rows)

	allStudents[0].SetRoomID(objects.NotLiving)
	secondStudentRows := studentObjectMother.CreateRows(allStudents[:1])
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondStudentRows)
	mock.ExpectExec("INSERT").WithArgs(StudentID, RoomID, objects.Get).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.SettleStudent(StudentNumber, RoomID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_SettleStudentNegativeStudentNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		RoomID        = 1
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.SettleStudent(StudentNumber, RoomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_SettleStudentNegativeStudentIsLiving(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentID     = 1
		RoomID        = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	roomsN := 1
	realRooms := roomObjectMother.CreateDefaultRooms(roomsN)
	rows := roomObjectMother.CreateRows(realRooms)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)
	mock.ExpectQuery("SELECT").WithArgs(RoomID).WillReturnError(nil).WillReturnRows(rows)

	secondStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.SettleStudent(StudentNumber, RoomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentAlreadyLiveErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_SettleStudentNegativeRoomNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		RoomID        = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)
	mock.ExpectQuery("SELECT").WithArgs(RoomID).WillReturnError(sql.ErrNoRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.SettleStudent(StudentNumber, RoomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_SettleStudentBadRoomIDParam(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		RoomID        = -1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.SettleStudent(StudentNumber, RoomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_SettleStudentBadStudNumberErr(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		RoomID        = 1
		StudentNumber = objects.EmptyString
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.SettleStudent(StudentNumber, RoomID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_EvicStudentPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentID     = 1
		RoomID        = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	secondStudentRows := studentObjectMother.CreateRows(allStudents[:1])
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondStudentRows)
	mock.ExpectExec("INSERT").WithArgs(StudentID, RoomID, objects.Ret).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.EvicStudent(StudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_EvicStudentNegativeStudentNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentNumber = mother.DefaultStudentNumber + "9"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.EvicStudent(StudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_EvicStudentNegativeStudentNotLiving(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentID     = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	allStudents[0].SetRoomID(objects.NotLiving)
	secondStudentRows := studentObjectMother.CreateRows(allStudents[:1])
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.EvicStudent(StudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotLivingErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_EvicStudentBadParam(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentNumber = objects.EmptyString
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.EvicStudent(StudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GetStudentByAccIDPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		AccID         = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	realStudentNumber, execErr := manager.GetStudentByAccID(AccID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, realStudentNumber, StudentNumber)
}

func (*TestStudentManager) TestStudentManager_GetStudentByAccIDNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		AccID = 7
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.GetStudentByAccID(AccID)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ViewAllStudents(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	realStudents, execErr := manager.ViewAllStudents()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, realStudents, allStudents)
}

func (*TestStudentManager) TestStudentManager_ViewStudentPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentID     = 1
		RoomID        = 1
		StudentNumber = mother.DefaultStudentNumber + "1"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	roomObjectMother := mother.RoomRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	roomsN := 1
	realRooms := roomObjectMother.CreateDefaultRooms(roomsN)
	rows := roomObjectMother.CreateRows(realRooms)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	secondStudentRows := studentObjectMother.CreateRows(allStudents[:1])
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondStudentRows)
	mock.ExpectQuery("SELECT").WithArgs(RoomID).WillReturnError(nil).WillReturnRows(rows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.ViewStudent(StudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ViewStudentNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)

	firstStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnError(nil).WillReturnRows(firstStudentRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	_, execErr := manager.ViewStudent(StudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GiveStudentThingPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		ThingID       = 1
		StudentID     = 3
		StudentNumber = mother.DefaultStudentNumber + "3"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(objects.None)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(nil).
		WillReturnRows(thingObjectMother.CreateRowForID(ThingID))

	realThingsRow := thingObjectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(ThingID).WillReturnError(nil).WillReturnRows(realThingsRow)

	secondAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondAllStudentRows)

	mock.ExpectExec("INSERT").WithArgs(StudentID, ThingID, objects.Get).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.GiveStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GiveStudentThingNegativeStudentNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(objects.None)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.GiveStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GiveStudentThingNegativeThingNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		StudentNumber = mother.DefaultStudentNumber + "3"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(objects.None)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(sql.ErrNoRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.GiveStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GiveStudentThingNegativeThingHasOwner(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		ThingID       = 1
		StudentNumber = mother.DefaultStudentNumber + "3"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(nil).
		WillReturnRows(thingObjectMother.CreateRowForID(ThingID))

	realThingsRow := thingObjectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(ThingID).WillReturnError(nil).WillReturnRows(realThingsRow)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.GiveStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.ThingHasOwnerErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GiveStudentThingNegativeBadStudentNumber(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		StudentNumber = objects.EmptyString
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.GiveStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_GiveStudentThingNegativeBadMarkNumber(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = -2
		StudentNumber = mother.DefaultStudentNumber
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.GiveStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ReturnStudentThingPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		ThingID       = 1
		StudentID     = 3
		StudentNumber = mother.DefaultStudentNumber + "3"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(StudentID)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(nil).
		WillReturnRows(thingObjectMother.CreateRowForID(ThingID))

	realThingsRow := thingObjectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(ThingID).WillReturnError(nil).WillReturnRows(realThingsRow)

	secondAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WithArgs(StudentID).WillReturnError(nil).WillReturnRows(secondAllStudentRows)

	mock.ExpectExec("INSERT").WithArgs(StudentID, ThingID, objects.Ret).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ReturnStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ReturnStudentThingNegativeStudentNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		StudentNumber = mother.DefaultStudentNumber + "7"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(objects.None)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ReturnStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ReturnStudentThingNegativeThingNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		StudentNumber = mother.DefaultStudentNumber + "3"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(objects.None)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(sql.ErrNoRows)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ReturnStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ReturnStudentThingNegativeThingHasNotOwner(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		ThingID       = 1
		StudentNumber = mother.DefaultStudentNumber + "3"
	)
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()
	studentsN := 4
	thingsN := 1
	allStudents := studentObjectMother.CreateDefaultStudents(studentsN)
	realThings := thingObjectMother.CreateDefaultThings(thingsN)
	realThings[0].SetOwnerID(objects.None)

	firstAllStudentRows := studentObjectMother.CreateRows(allStudents)
	mock.ExpectQuery("SELECT").WillReturnRows(firstAllStudentRows).WillReturnError(nil)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(nil).
		WillReturnRows(thingObjectMother.CreateRowForID(ThingID))

	realThingsRow := thingObjectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(ThingID).WillReturnError(nil).WillReturnRows(realThingsRow)

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ReturnStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.StudentIsNotOwnerErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ReturnStudentThingNegativeBadMarkNumber(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = -2
		StudentNumber = mother.DefaultStudentNumber
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ReturnStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestStudentManager) TestStudentManager_ReturnStudentThingNegativeBadStudentNumber(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		MarkNumber    = 123
		StudentNumber = objects.EmptyString
	)
	studentObjectMother := mother.StudentRepoObjectMother{}

	db, mock := studentObjectMother.CreateRepo()

	studentRepository := studentRepo.PgStudentRepo{Conn: db}
	userRepository := userRepo.PgUserRepo{Conn: db}
	thingRepository := thingRepo.PgThingRepo{Conn: db}
	roomRepository := roomRepo.PgRoomRepo{Conn: db}

	studentC := studentController.StudentController{Repo: &studentRepository}
	userC := userController.UserController{Repo: &userRepository}
	thingC := thingController.ThingController{Repo: &thingRepository}
	roomC := roomController.RoomController{Repo: &roomRepository}

	manager := StudentManager{studentController: studentC, userController: userC, thingController: thingC,
		roomController: roomC}

	// Act
	execErr := manager.ReturnStudentThing(StudentNumber, MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadStudentParamsErr)
	tests.AssertMocks(t, mock)
}
