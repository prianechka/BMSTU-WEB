package mother

import (
	"database/sql"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
)

var (
	DefaultStudentName    = "Ivan"
	DefaultStudentSurname = "Ivanov"
	DefaultGroup          = "IU7-65B"
	DefaultStudentNumber  = "19u609"
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
			DefaultStudentSurname, DefaultGroup, DefaultStudentNumber+fmt.Sprintf("%d", i), i))
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

func (m StudentRepoObjectMother) CreateRowForID(id int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"studentid"})
	rows.AddRow(id)
	return rows
}

func (m StudentRepoObjectMother) CreateStudentDTO() objects.StudentDTO {
	return objects.NewStudentDTO(DefaultStudentName, DefaultStudentSurname, DefaultGroup, DefaultGroup+"0")
}
