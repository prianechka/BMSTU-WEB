package studentRepo

import (
	"database/sql"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"src/objects"
	"testing"
)

var (
	DefaultStudentName          = "Ivan"
	DefaultStudentSurname       = "Ivanov"
	DefaultGroup                = "IU7-65B"
	InsertID              int64 = 5
	RowsAffected          int64 = 1
)

type StudentRepoObjectMother struct{}

func (m StudentRepoObjectMother) CreateRepo() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func (m StudentRepoObjectMother) CreateDefaultStudents(amount int) []objects.Student {
	resultStudents := make([]objects.Student, objects.Empty)
	for i := 1; i <= amount; i++ {
		resultStudents = append(resultStudents, objects.NewStudentWithParams(i, i, DefaultStudentName,
			DefaultStudentSurname, DefaultGroup, DefaultGroup+fmt.Sprintf("%d", i), i))
	}
	return resultStudents
}

func (m StudentRepoObjectMother) CreateRows(students []objects.Student) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"studentid", "webaccid", "studentname", "studentsurname",
		"studentgroup", "studentnumber", "roomid"})
	for _, student := range students {
		rows.AddRow(student.GetID(), student.GetAccID(), student.GetName(), student.GetSurname(),
			student.GetStudentGroup(), student.GetStudentNumber(), student.GetRoomID())
	}
	return rows
}

func (m StudentRepoObjectMother) CreateStudentDTO() objects.StudentDTO {
	return objects.NewStudentDTO(DefaultStudentName, DefaultStudentSurname, DefaultGroup, DefaultGroup+"0")
}

func TestPgStudentRepo_AddStudent(t *testing.T) {
	// Arrange
	objectMother := StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	studentDTO := objectMother.CreateStudentDTO()
	mock.ExpectExec("INSERT INTO").WithArgs(studentDTO.GetName(), studentDTO.GetSurname(),
		studentDTO.GetStudentGroup(), studentDTO.GetStudentNumber(), InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.AddStudent(studentDTO, int(InsertID))

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

func TestPgStudentRepo_GetAllStudents(t *testing.T) {
	// Arrange
	objectMother := StudentRepoObjectMother{}
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
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(realStudents, resultStudents) {
		t.Errorf("results not match, want %v, have %v", realStudents, resultStudents)
		return
	}
}

func TestPgStudentRepo_ChangeStudent(t *testing.T) {
	// Arrange
	objectMother := StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	studentDTO := objectMother.CreateStudentDTO()
	mock.ExpectExec("UPDATE").WithArgs(studentDTO.GetName(), studentDTO.GetSurname(),
		studentDTO.GetStudentGroup(), studentDTO.GetStudentNumber(), InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.ChangeStudent(int(InsertID), studentDTO)

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

// TestPgStudentRepo_GetStudentPositive проверяет, что если студент есть, он успешно вернётся.
func TestPgStudentRepo_GetStudentPositive(t *testing.T) {
	// Arrange
	objectMother := StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realRooms := objectMother.CreateDefaultStudents(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgStudentRepo{Conn: db}

	// Act
	room, execErr := repo.GetStudent(id)

	// Assert
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(room, realRooms[0]) {
		t.Errorf("results not match, want %v, have %v", realRooms[0], room)
		return
	}
}

// TestPgStudentRepo_GetStudentNegative проверяет, что если студента нет, то вернётся ошибка.
func TestPgStudentRepo_GetStudentNegative(t *testing.T) {
	// Arrange
	objectMother := StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(StudentNotFoundErr)
	repo := PgStudentRepo{Conn: db}

	// Act
	_, execErr := repo.GetStudent(int(InsertID))

	// Assert
	if execErr != StudentNotFoundErr {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestPgStudentRepo_TransferStudent(t *testing.T) {
	// Arrange
	var (
		roomID                              = int(InsertID)
		studentID                           = int(InsertID)
		dir       objects.TransferDirection = objects.Get
	)
	objectMother := StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(roomID, studentID, dir).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.TransferStudent(studentID, roomID, dir)

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

func TestPgStudentRepo_TransferThing(t *testing.T) {
	// Arrange
	var (
		thingID                             = int(InsertID)
		studentID                           = int(InsertID)
		dir       objects.TransferDirection = objects.Get
	)
	objectMother := StudentRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(thingID, studentID, dir).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgStudentRepo{Conn: db}

	// Act
	execErr := repo.TransferThing(studentID, thingID, dir)

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
