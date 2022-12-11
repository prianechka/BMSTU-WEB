package studentRepo

import "src/objects"

type StudentRepo interface {
	AddStudent(newStudent objects.StudentDTO, accID int) error
	GetAllStudents(page, size int) ([]objects.Student, error)
	GetStudentID(studentNumber string) (int, error)
	GetStudent(id int) (objects.Student, error)
	TransferStudent(studentID, roomID int, direct objects.TransferDirection) error
	ChangeStudent(studentID int, studentInfo objects.StudentDTO) error
	TransferThing(studentID, thingID int, direct objects.TransferDirection) error
	GetStudentThings(id int, page, size int) ([]objects.Thing, error)
}
