package userRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/tests"
	"src/tests/mother"
	"testing"
)

var (
	InsertID     int64 = 5
	RowsAffected int64 = 1
)

func TestPgUserRepo_AddUser(t *testing.T) {
	// Arrange
	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	userDTO := objectMother.CreateUser()
	mock.ExpectExec("INSERT INTO").WithArgs(userDTO.GetLogin(), userDTO.GetPassword(),
		userDTO.GetPrivelegeLevel()).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgUserRepo{Conn: db}

	// Act
	execErr := repo.AddUser(userDTO.GetLogin(), userDTO.GetPassword(), userDTO.GetPrivelegeLevel())

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgUserRepo_GetUserPositive(t *testing.T) {
	// Arrange
	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realUsers := objectMother.CreateDefaultUsers(N)
	rows := objectMother.CreateRows(realUsers)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgUserRepo{Conn: db}

	// Act
	user, execErr := repo.GetUser(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, user, realUsers[0])
}

func TestPgUserRepo_GetUserNegative(t *testing.T) {
	// Arrange
	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(sql.ErrNoRows)
	repo := PgUserRepo{Conn: db}

	// Act
	_, execErr := repo.GetUser(int(InsertID))

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}

func TestPgUserRepo_GetUserIDPositive(t *testing.T) {
	// Arrange
	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	rows := objectMother.CreateRowForID(ID)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgUserRepo{Conn: db}

	// Act
	userID, execErr := repo.GetUserID(mother.DefaultLogin)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, userID, ID)
}

func TestPgUserRepo_GetUserIDNegative(t *testing.T) {
	// Arrange
	objectMother := mother.UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultLogin).
		WillReturnError(sql.ErrNoRows)
	repo := PgUserRepo{Conn: db}

	// Act
	_, execErr := repo.GetUserID(mother.DefaultLogin)

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}
