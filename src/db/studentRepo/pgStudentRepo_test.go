package studentRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	"testing"
)

var (
	InsertID     int64 = 5
	RowsAffected int64 = 1
)

func TestPgStudentRepo_AddStudent(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	studentDTO := objectMother.CreateStudentDTO()
	mock.ExpectExec("INSERT INTO").WithArgs(studentDTO.GetName(), studentDTO.GetSurname(),
		studentDTO.GetStudentGroup(), studentDTO.GetStudentNumber(), InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.AddStudent(studentDTO, int(InsertID))

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgStudentRepo_GetAllStudents(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows).WillReturnError(nil)

	repo := PgStudentRepo{Conn: db}

	// Act
	resultStudents, execErr := repo.GetAllStudents()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, resultStudents, realStudents)
}

func TestPgStudentRepo_ChangeStudent(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	studentDTO := objectMother.CreateStudentDTO()
	mock.ExpectExec("UPDATE").WithArgs(studentDTO.GetName(), studentDTO.GetSurname(),
		studentDTO.GetStudentGroup(), studentDTO.GetStudentNumber(), InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.ChangeStudent(int(InsertID), studentDTO)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

// TestPgStudentRepo_GetStudentPositive проверяет, что если студент есть, он успешно вернётся.
func TestPgStudentRepo_GetStudentPositive(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realStudents := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realStudents)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgStudentRepo{Conn: db}

	// Act
	student, execErr := repo.GetStudent(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, student, realStudents[0])
}

// TestPgStudentRepo_GetStudentNegative проверяет, что если студента нет, то вернётся ошибка.
func TestPgStudentRepo_GetStudentNegative(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(sql.ErrNoRows)
	repo := PgStudentRepo{Conn: db}

	// Act
	_, execErr := repo.GetStudent(int(InsertID))

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}

func TestPgStudentRepo_TransferStudent(t *testing.T) {
	// Arrange
	var (
		roomID                              = int(InsertID)
		studentID                           = int(InsertID)
		dir       objects.TransferDirection = objects.Get
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(roomID, studentID, dir).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.TransferStudent(studentID, roomID, dir)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgStudentRepo_TransferThing(t *testing.T) {
	// Arrange
	var (
		thingID                             = int(InsertID)
		studentID                           = int(InsertID)
		dir       objects.TransferDirection = objects.Get
	)
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(thingID, studentID, dir).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.TransferThing(studentID, thingID, dir)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgStudentRepo_GetStudentIDPositive(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	ID := 1
	rows := objectMother.CreateRowForID(ID)
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultStudentNumber).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgStudentRepo{Conn: db}

	// Act
	studentID, execErr := repo.GetStudentID(mother.DefaultStudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, studentID, ID)
}

func TestPgStudentRepo_GetStudentIDNegative(t *testing.T) {
	// Arrange
	objectMother := mother.StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mother.DefaultStudentNumber).
		WillReturnError(sql.ErrNoRows)
	repo := PgStudentRepo{Conn: db}

	// Act
	_, execErr := repo.GetStudentID(mother.DefaultStudentNumber)

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}

func TestPgStudentRepo_GetStudentThings(t *testing.T) {
	// Arrange
	studentObjectMother := mother.StudentRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}
	db, mock := studentObjectMother.CreateRepo()
	N := 3
	id := 1
	realThings := thingObjectMother.CreateDefaultThings(N)
	rows := thingObjectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgStudentRepo{Conn: db}

	// Act
	things, execErr := repo.GetStudentThings(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, things, realThings)
}
