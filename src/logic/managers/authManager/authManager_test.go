package authManager

import (
	"database/sql"
	"src/db/userRepo"
	"src/logic/controllers/userController"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	"testing"
)

func TestAuthManager_GetUserIDPositive(t *testing.T) {
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

func TestAuthManager_GetUserIDNegative(t *testing.T) {
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
	tests.AssertErrors(t, execErr, userController.UserNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestAuthManager_TryToAuthPositive(t *testing.T) {
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

func TestAuthManager_TryToAuthNegativeUserNotFound(t *testing.T) {
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
	tests.AssertErrors(t, execErr, userController.UserNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestAuthManager_TryToAuthNegativeBadPassword(t *testing.T) {
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
	tests.AssertErrors(t, execErr, PasswordNotEqualErr)
	tests.AssertMocks(t, mock)
}
