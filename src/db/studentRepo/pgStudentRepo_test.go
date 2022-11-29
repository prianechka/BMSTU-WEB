package studentRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
	"testing"
)

var (
	DefaultStudentName    = "Ivan"
	DefaultStudentSurname = "Ivanov"
	DefaultGroup          = "IU7-65B"
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
			DefaultStudentSurname, DefaultGroup, DefaultGroup+string(i), i))
	}
	return resultStudents
}

func (m StudentRepoObjectMother) CreateRows([]objects.Student) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"roomid", "roomtype", "roomnumber"})
	for _, room := range rooms {
		rows.AddRow(room.GetID(), room.GetRoomType(), room.GetRoomNumber())
	}
	return rows
}

func TestPgStudentRepo_AddStudent(t *testing.T) {

}
