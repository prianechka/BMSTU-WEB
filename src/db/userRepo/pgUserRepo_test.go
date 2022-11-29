package userRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"src/objects"
	"testing"
)

type UserRepoObjectMother struct{}

var (
	DefaultLogin                   = "Ivan"
	DefaultPassword                = "12345678"
	DefaultRole     objects.Levels = objects.StudentRole
	InsertID        int64          = 5
	RowsAffected    int64          = 1
)

func (m UserRepoObjectMother) CreateRepo() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func (m UserRepoObjectMother) CreateUser() objects.User {
	return objects.NewUserWithParams(int(InsertID), DefaultLogin, DefaultPassword, DefaultRole)
}

func (m UserRepoObjectMother) CreateDefaultUsers(amount int) []objects.User {
	resultUsers := make([]objects.User, objects.Empty)
	for i := 1; i <= amount; i++ {
		resultUsers = append(resultUsers, objects.NewUserWithParams(i, DefaultLogin,
			DefaultPassword, DefaultRole))
	}
	return resultUsers
}

func (m UserRepoObjectMother) CreateRows(users []objects.User) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "userlogin", "userpassword", "userrole"})
	for _, user := range users {
		rows.AddRow(user.GetID(), user.GetLogin(), user.GetPassword(), user.GetPrivelegeLevel())
	}
	return rows
}

func TestPgUserRepo_AddUser(t *testing.T) {
	// Arrange
	objectMother := UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	userDTO := objectMother.CreateUser()
	mock.ExpectExec("INSERT INTO").WithArgs(userDTO.GetLogin(), userDTO.GetPassword(),
		userDTO.GetPrivelegeLevel()).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgUserRepo{Conn: db}

	// Act
	execErr := repo.AddUser(userDTO.GetLogin(), userDTO.GetPassword(), userDTO.GetPrivelegeLevel())

	// Assert
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestPgUserRepo_GetUserPositive(t *testing.T) {
	// Arrange
	objectMother := UserRepoObjectMother{}
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
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(user, realUsers[0]) {
		t.Errorf("results not match, want %v, have %v", realUsers[0], user)
		return
	}
}

func TestPgUserRepo_GetUserNegative(t *testing.T) {
	// Arrange
	objectMother := UserRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(UserNotFoundErr)
	repo := PgUserRepo{Conn: db}

	// Act
	_, execErr := repo.GetUser(int(InsertID))

	// Assert
	if execErr != UserNotFoundErr {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}
