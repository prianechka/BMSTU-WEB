package studentController

import (
	"database/sql"
	"src/db/studentRepo"
	"src/objects"
)

type StudentController struct {
	repo studentRepo.StudentRepo
}

func (sc *StudentController) AddStudent(name, surname, group, studentNumber string, accID int) error {
	var err error
	if accID < 0 {
		err = BadAccIDErr
	} else if len(name) < 1 || len(surname) < 1 || len(group) < 1 || len(studentNumber) < 1 {
		err = BadParamsErr
	} else {
		allStudents, getStudentErr := sc.repo.GetAllStudents()
		if getStudentErr == nil {
			for _, tmpStudent := range allStudents {
				if tmpStudent.GetStudentNumber() == studentNumber {
					err = StudentAlreadyInBaseErr
					break
				}
			}
			if err == nil {
				studentDTO := objects.NewStudentDTO(name, surname, group, studentNumber)
				err = sc.repo.AddStudent(studentDTO, accID)
			}
		}
	}

	return err
}

func (sc *StudentController) GetAllStudents() ([]objects.Student, error) {
	return sc.repo.GetAllStudents()
}

func (sc *StudentController) GetStudentIDByNumber(studentNumber string) (int, error) {
	var result = objects.None
	allStudents, err := sc.repo.GetAllStudents()
	if err == nil {
		for _, tmpStudent := range allStudents {
			if tmpStudent.GetStudentNumber() == studentNumber {
				result = tmpStudent.GetID()
				break
			}
		}
	}
	if result == objects.None {
		err = StudentNotFoundErr
	}

	return result, err
}

func (sc *StudentController) GetStudent(id int) (objects.Student, error) {
	var err error
	var student objects.Student
	if id < 0 {
		err = BadParamsErr
	} else {
		student, err = sc.repo.GetStudent(id)
		if err != nil {
			err = StudentNotFoundErr
		}
	}
	return student, err
}

func (sc *StudentController) GetStudentRoom(id int) (int, error) {
	var result = objects.None
	allStudents, err := sc.repo.GetAllStudents()
	if err == nil {
		for _, tmpStudent := range allStudents {
			if tmpStudent.GetID() == id {
				result = tmpStudent.GetRoomID()
				break
			}
		}
	}
	if result == objects.None {
		err = StudentNotFoundErr
	}
	return result, err
}

func (sc *StudentController) SettleStudent(studentID, roomID int) error {
	student, err := sc.repo.GetStudent(studentID)
	if err == nil {
		if student.GetID() == objects.None {
			err = StudentNotFoundErr
		} else if student.GetRoomID() == objects.NotLiving {
			err = sc.repo.TransferStudent(studentID, roomID, objects.Get)
		} else {
			err = StudentAlreadyLiveErr
		}
	} else if err == sql.ErrNoRows {
		err = StudentNotFoundErr
	}
	return err
}

func (sc *StudentController) EvicStudent(studentID int) error {
	student, err := sc.repo.GetStudent(studentID)
	if err == nil {
		if student.GetRoomID() != objects.NotLiving {
			err = sc.repo.TransferStudent(studentID, student.GetRoomID(), objects.Ret)
		} else {
			err = StudentNotLivingErr
		}
	} else if err == sql.ErrNoRows {
		err = StudentNotFoundErr
	}
	return err
}

func (sc *StudentController) ChangeStudentGroup(studentID int, newGroup string) error {
	student, err := sc.repo.GetStudent(studentID)
	if err == nil {
		studentDTO := objects.NewStudentDTO(student.GetName(), student.GetSurname(),
			newGroup, student.GetStudentNumber())
		err = sc.repo.ChangeStudent(studentID, studentDTO)
	} else if err == sql.ErrNoRows {
		err = StudentNotFoundErr
	}
	return err
}

func (sc *StudentController) GetStudentThings(studentID int) ([]objects.Thing, error) {
	studentThings := make([]objects.Thing, 0)
	_, err := sc.repo.GetStudent(studentID)
	if err == nil {
		studentThings, err = sc.repo.GetStudentThings(studentID)
	} else {
		err = StudentNotFoundErr
	}
	return studentThings, err
}

func (sc *StudentController) TransferThing(studentID, thingID int) error {
	_, err := sc.repo.GetStudent(studentID)
	if err == nil {
		err = sc.repo.TransferThing(studentID, thingID, objects.Get)
	} else if err == sql.ErrNoRows {
		err = StudentNotFoundErr
	}
	return err
}

func (sc *StudentController) ReturnThing(studentID, thingID int) error {
	_, err := sc.repo.GetStudent(studentID)
	if err == nil {
		err = sc.repo.TransferThing(studentID, thingID, objects.Ret)
	} else if err == sql.ErrNoRows {
		err = StudentNotFoundErr
	}
	return err
}
