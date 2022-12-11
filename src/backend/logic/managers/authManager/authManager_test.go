package authManager

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"src/db/userRepo"
	"src/logic/controllers/userController"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	appErrors "src/utils/error"
	"testing"
	"time"
)

type TestAuthManager struct{}

func Test_AuthManager(t *testing.T) {
	testgroup.RunSerially(t, &TestAuthManager{})
}

func (*TestAuthManager) TestAuthManager_GetUserIDPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	firstRows := objectMother.CreateRowForID(ID)
	secondRows := objectMother.CreateRowForID(ID)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(firstRows)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(secondRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := userController.UserController{Repo: &repo}
	manager := AuthManager{
		currentLogin:   mother.DefaultLogin,
		userController: controller,
	}

	// Act
	userID, execErr := manager.GetUserID(mother.DefaultLogin)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, userID, ID)
}

func (*TestAuthManager) TestAuthManager_GetUserIDNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(sql.ErrNoRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := userController.UserController{Repo: &repo}
	manager := AuthManager{
		currentLogin:   mother.DefaultLogin,
		userController: controller,
	}

	// Act
	_, execErr := manager.GetUserID(mother.DefaultLogin)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.UserNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestAuthManager) TestAuthManager_TryToAuthPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	N := 1
	firstRows := objectMother.CreateRowForID(ID)
	secondRows := objectMother.CreateRowForID(ID)
	users := objectMother.CreateDefaultUsers(N)
	thirdRows := objectMother.CreateRows(users)

	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(firstRows)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(secondRows)
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(thirdRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := userController.UserController{Repo: &repo}
	manager := AuthManager{
		currentLogin:   mother.DefaultLogin,
		userController: controller,
	}

	// Act
	role, execErr := manager.TryToAuth(mother.DefaultLogin, mother.DefaultPassword)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, role, objects.Levels(objects.StudentRole))
}

func (*TestAuthManager) TestAuthManager_TryToAuthNegativeUserNotFound(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()

	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(sql.ErrNoRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := userController.UserController{Repo: &repo}
	manager := AuthManager{
		currentLogin:   mother.DefaultLogin,
		userController: controller,
	}

	// Act
	_, execErr := manager.TryToAuth(mother.DefaultLogin, mother.DefaultPassword)

	// Assert
	tests.AssertErrors(t, execErr, appErrors.UserNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestAuthManager) TestAuthManager_TryToAuthNegativeBadPassword(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	N := 1
	firstRows := objectMother.CreateRowForID(ID)
	secondRows := objectMother.CreateRowForID(ID)
	users := objectMother.CreateDefaultUsers(N)
	thirdRows := objectMother.CreateRows(users)

	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(firstRows)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).WillReturnError(nil).WillReturnRows(secondRows)
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(thirdRows)

	repo := userRepo.PgUserRepo{Conn: db}
	controller := userController.UserController{Repo: &repo}
	manager := AuthManager{
		currentLogin:   mother.DefaultLogin,
		userController: controller,
	}

	// Act
	_, execErr := manager.TryToAuth(mother.DefaultLogin, mother.DefaultPassword+"1")

	// Assert
	tests.AssertErrors(t, execErr, appErrors.PasswordNotEqualErr)
	tests.AssertMocks(t, mock)
}
