package studentManager

import (
	"database/sql"
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
	"testing"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

func TestStudentManager_AddNewStudentPositive(t *testing.T) {
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

func TestStudentManager_AddNewStudentNegativeLoginOccuped(t *testing.T) {
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
	tests.AssertErrors(t, execErr, userController.LoginOccupedErr)
	tests.AssertMocks(t, mock)
}

func TestStudentManager_AddNewStudentNegativeBadLoginParam(t *testing.T) {
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
	tests.AssertErrors(t, execErr, studentController.BadParamsErr)
	tests.AssertMocks(t, mock)
}

func TestStudentManager_AddNewStudentNegativeBadLogin(t *testing.T) {
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
	tests.AssertErrors(t, execErr, studentController.BadParamsErr)
	tests.AssertMocks(t, mock)
}

func TestStuden
