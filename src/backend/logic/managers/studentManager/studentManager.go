package studentManager

import (
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/logic/controllers/userController"
	"src/logic/managers/models"
	"src/objects"
	appErrors "src/utils/error"
)

type StudentManager struct {
	roomController    roomController.RoomController
	studentController studentController.StudentController
	userController    userController.UserController
	thingController   thingController.ThingController
}

func CreateNewStudentManager(rc roomController.RoomController, sc studentController.StudentController,
	uc userController.UserController, tc thingController.ThingController) *StudentManager {
	return &StudentManager{
		roomController:    rc,
		studentController: sc,
		userController:    uc,
		thingController:   tc,
	}
}

func (sm *StudentManager) AddNewStudent(name, surname, studentGroup, studentNumber, login, password string) (err error) {
	if name == objects.EmptyString || surname == objects.EmptyString || studentGroup == objects.EmptyString ||
		studentNumber == objects.EmptyString {
		err = appErrors.BadStudentParamsErr
	} else if login == objects.EmptyString || password == objects.EmptyString {
		err = appErrors.BadUserParamsErr
	} else {
		isUserExist := sm.userController.UserExist(login)
		if !isUserExist {
			_, getStudentErr := sm.studentController.GetStudentIDByNumber(studentNumber)
			if getStudentErr == appErrors.StudentNotFoundErr {
				addUserErr := sm.userController.AddUser(login, password, objects.StudentRole)
				if addUserErr == nil {
					accID, getUserErr := sm.userController.GetUserID(login)
					if getUserErr == nil {
						err = sm.studentController.AddStudent(name, surname, studentGroup, studentNumber, accID)
					} else {
						err = getUserErr
					}
				} else {
					err = addUserErr
				}
			} else {
				err = appErrors.StudentAlreadyInBaseErr
			}
		} else {
			err = appErrors.LoginOccupedErr
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
			studentInfo.Student = student
		} else {
			err = getStudentErr
		}
	}
	return studentInfo, err
}

func (sm *StudentManager) ViewAllStudents(page, size int) ([]objects.Student, error) {
	return sm.studentController.GetAllStudents(page, size)
}

func (sm *StudentManager) ChangeStudentGroup(studentNumber, newGroup string) (err error) {
	if newGroup == objects.EmptyString {
		err = appErrors.BadStudentParamsErr
	} else {
		studID, getStudentErr := sm.studentController.GetStudentIDByNumber(studentNumber)
		if getStudentErr == nil {
			err = sm.studentController.ChangeStudentGroup(studID, newGroup)
		} else {
			err = getStudentErr
		}
	}
	return err
}

func (sm *StudentManager) SettleStudent(studentNumber string, roomID int) error {
	if studentNumber == objects.EmptyString {
		return appErrors.BadStudentParamsErr
	}

	if roomID <= objects.NotLiving {
		return appErrors.RoomNotFoundErr
	}

	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		_, err = sm.roomController.GetRoom(roomID)
		if err == nil {
			err = sm.studentController.SettleStudent(studentID, roomID)
		}
	}
	return err
}

func (sm *StudentManager) EvicStudent(studentNumber string) error {
	if studentNumber == objects.EmptyString {
		return appErrors.BadStudentParamsErr
	}

	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		err = sm.studentController.EvicStudent(studentID)
	}
	return err
}

func (sm *StudentManager) GiveStudentThing(studentNumber string, markNumber int) error {
	if studentNumber == objects.EmptyString {
		return appErrors.BadStudentParamsErr
	}

	if markNumber <= objects.None {
		return appErrors.BadThingParamsErr
	}

	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		thingID, getThingErr := sm.thingController.GetThingIDByMarkNumber(markNumber)
		if getThingErr == nil {
			ownerID, getOwnerErr := sm.thingController.GetCurrentOwner(thingID)
			if getOwnerErr == nil {
				if ownerID == objects.None {
					err = sm.studentController.TransferThing(studentID, thingID)
				} else {
					err = appErrors.ThingHasOwnerErr
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
	if studentNumber == objects.EmptyString {
		return appErrors.BadStudentParamsErr
	}

	if markNumber <= objects.None {
		return appErrors.BadThingParamsErr
	}

	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		thingID, getThingErr := sm.thingController.GetThingIDByMarkNumber(markNumber)
		if getThingErr == nil {
			ownerID, getOwnerErr := sm.thingController.GetCurrentOwner(thingID)
			if getOwnerErr == nil {
				if ownerID == studentID {
					err = sm.studentController.ReturnThing(studentID, thingID)
				} else {
					err = appErrors.StudentIsNotOwnerErr
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
	allStudents, err := sm.studentController.GetAllStudents(objects.Null, objects.Null)
	if err == nil {
		for _, tmpStudent := range allStudents {
			if tmpStudent.GetAccID() == accID {
				resultStudentNumber = tmpStudent.GetStudentNumber()
				break
			}
		}
	}
	if resultStudentNumber == objects.EmptyString {
		err = appErrors.StudentNotFoundErr
	}
	return resultStudentNumber, err
}

func (sm *StudentManager) GetCurrentRoom(studentNumber string) (int, error) {
	var result = objects.None
	if studentNumber == objects.EmptyString {
		return result, appErrors.BadStudentParamsErr
	}

	studentID, err := sm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		student, getStudentErr := sm.studentController.GetStudent(studentID)
		if getStudentErr == nil {
			result = student.GetRoomID()
		} else {
			err = getStudentErr
		}
	}
	return result, err
}
