package userController

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/userRepo"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	appErrors "src/utils/error"
	"testing"
	"time"
)

var (
	InsertID     int64 = 5
	RowsAffected int64 = 1
)

type TestUserController struct{}

func Test_UserController(t *testing.T) {
	testgroup.RunSerially(t, &TestUserController{})
}

func (*TestUserController) TestUserController_AddUserPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WithArgs(mother.DefaultLogin, mother.DefaultPassword, objects.StudentRole).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	execErr := controller.AddUser(mother.DefaultLogin, mother.DefaultPassword, objects.StudentRole)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestUserController) TestUserController_AddUserNegativeAlreadyExist(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	rows := objectMother.CreateRowForID(ID)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(rows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	execErr := controller.AddUser(mother.DefaultLogin, mother.DefaultPassword, objects.StudentRole)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.LoginOccupedErr)
	tests.AssertMocks(t, mock)
}

func (*TestUserController) TestUserController_AddUserNegativeBadParams(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	execErr := controller.AddUser(objects.EmptyString, mother.DefaultPassword, objects.StudentRole)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.BadUserParamsErr)
	tests.AssertMocks(t, mock)
}

func (*TestUserController) TestUserController_UserExistTrue(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	rows := objectMother.CreateRowForID(ID)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).
		WillReturnError(nil).WillReturnRows(rows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	result := controller.UserExist(mother.DefaultLogin)

	// Assert
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, result, true)
}

func (*TestUserController) TestUserController_UserExistFalse(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(sql.ErrNoRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	result := controller.UserExist(mother.DefaultLogin)

	// Assert
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, result, false)
}

func (*TestUserController) TestUserController_GetUserPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realUsers := objectMother.CreateDefaultUsers(N)
	rows := objectMother.CreateRows(realUsers)
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(nil).WillReturnRows(rows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	user, execErr := controller.GetUser(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, user, realUsers[0])
}

func (*TestUserController) TestUserController_GetUserNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	id := 1
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(sql.ErrNoRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	_, execErr := controller.GetUser(id)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.UserNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestUserController) TestUserController_GetUserIDPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	rows := objectMother.CreateRowForID(ID)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).
		WillReturnError(nil).WillReturnRows(rows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	userID, execErr := controller.GetUserID(mother.DefaultLogin)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, userID, ID)
}

func (*TestUserController) TestUserController_GetUserIDNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).
		WillReturnError(sql.ErrNoRows)
	repo := userRepo.PgUserRepo{Conn: db}
	controller := UserController{Repo: &repo}

	// Act
	_, execErr := controller.GetUserID(mother.DefaultLogin)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.UserNotFoundErr)
	tests.AssertMocks(t, mock)
}
