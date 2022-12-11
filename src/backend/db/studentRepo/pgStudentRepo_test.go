package studentRepo

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
	"src/tests"
	"src/tests/mother"
	"testing"
	"time"
)

var (
	InsertID     int64 = 5
	RowsAffected int64 = 1
)

type TestPgStudentRepo struct{}

func Test_PgStudentRepo(t *testing.T) {
	testgroup.RunSerially(t, &TestPgStudentRepo{})
}

func (*TestPgStudentRepo) TestPgStudentRepo_AddStudent(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_GetAllStudents(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_ChangeStudent(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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
func (*TestPgStudentRepo) TestPgStudentRepo_GetStudentPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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
func (*TestPgStudentRepo) TestPgStudentRepo_GetStudentNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_TransferStudent(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_TransferThing(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_GetStudentIDPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_GetStudentIDNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgStudentRepo) TestPgStudentRepo_GetStudentThings(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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
