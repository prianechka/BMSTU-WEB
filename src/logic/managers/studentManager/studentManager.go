package studentManager

import (
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/logic/controllers/userController"
	"src/logic/managers/models"
	"src/objects"
)

type StudentManager struct {
	roomController    roomController.RoomController
	studentController studentController.StudentController
	userController    userController.UserController
	thingController   thingController.ThingController
}

func (sm *StudentManager) AddNewStudent(name, surname, studentGroup, studentNumber,
	login, password string) error {
	err := sm.userController.AddUser(login, password, objects.StudentRole)
	if err == nil {
		accID, getUserErr := sm.userController.GetUserID(login)
		if getUserErr == nil {
			err = sm.studentController.AddStudent(name, surname, studentGroup, studentNumber, accID)
		} else {
			err = getUserErr
		}
	}
	return err
}

func (sm *StudentManager) ViewStudent(studentNumber string) (models.StudentFullInfo, error) {
	var studentInfo models.StudentFullInfo
	studID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		student, getStudentErr := sm.studentController.GetStudent(studID)
		if getStudentErr == nil {
			room, getRoomErr := sm.roomController.GetRoom(student.GetRoomID())
			if getRoomErr == nil {
				studentInfo.Student = student
				studentInfo.Room = room
			} else {
				err = getRoomErr
			}
		} else {
			err = getStudentErr
		}
	}
	return studentInfo, err
}

func (sm *StudentManager) ViewAllStudents() ([]objects.Student, error) {
	return sm.studentController.GetAllStudents()
}

func (sm *StudentManager) ChangeStudentGroup(studentNumber, newGroup string) error {
	studID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		err = sm.studentController.ChangeStudentGroup(studID, newGroup)
	}
	return err
}

func (sm *StudentManager) SettleStudent(studentNumber string, roomID int) error {
	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		err = sm.studentController.SettleStudent(studentID, roomID)
	}
	return err
}

func (sm *StudentManager) EvicStudent(studentNumber string) error {
	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		err = sm.studentController.EvicStudent(studentID)
	}
	return err
}

func (sm *StudentManager) GiveStudentThing(studentNumber string, markNumber int) error {
	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		thingID, getThingErr := sm.thingController.GetThingIDByMarkNumber(markNumber)
		if getThingErr == nil {
			ownerID, getOwnerErr := sm.thingController.GetCurrentOwner(thingID)
			if getOwnerErr == nil {
				if ownerID == objects.None {
					err = sm.studentController.TransferThing(studentID, thingID)
				}
			} else {
				err = getOwnerErr
			}
		} else {
			err = getThingErr
		}
	}
	return err
}

func (sm *StudentManager) ReturnStudentThing(studentNumber string, markNumber int) error {
	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		thingID, getThingErr := sm.thingController.GetThingIDByMarkNumber(markNumber)
		if getThingErr == nil {
			ownerID, getOwnerErr := sm.thingController.GetCurrentOwner(thingID)
			if getOwnerErr == nil {
				if ownerID == studentID {
					err = sm.studentController.ReturnThing(studentID, thingID)
				}
			} else {
				err = getOwnerErr
			}
		} else {
			err = getThingErr
		}
	}
	return err
}

func (sm *StudentManager) GetStudentByAccID(accID int) (string, error) {
	var resultStudentNumber string
	allStudents, err := sm.studentController.GetAllStudents()
	if err == nil {
		for _, tmpStudent := range allStudents {
			if tmpStudent.GetAccID() == accID {
				resultStudentNumber = tmpStudent.GetStudentNumber()
				break
			}
		}
	}
	if resultStudentNumber == objects.EmptyString {
		err = studentController.StudentNotFoundErr
	}
	return resultStudentNumber, err
}
