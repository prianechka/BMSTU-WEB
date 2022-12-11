package thingManager

import (
	"src/logic/controllers/roomController"
	"src/logic/controllers/studentController"
	"src/logic/controllers/thingController"
	"src/logic/managers/models"
	"src/objects"
	appErrors "src/utils/error"
)

type ThingManager struct {
	thingController   thingController.ThingController
	roomController    roomController.RoomController
	studentController studentController.StudentController
}

func CreateNewThingManager(rc roomController.RoomController, sc studentController.StudentController,
	tc thingController.ThingController) *ThingManager {
	return &ThingManager{
		roomController:    rc,
		studentController: sc,
		thingController:   tc,
	}
}

func (tm *ThingManager) GetFullThingInfo(page, size int) ([]models.ThingFullInfo, error) {
	thingsFullInfo := make([]models.ThingFullInfo, 0)
	allThings, err := tm.thingController.GetThings(page, size)
	if err == nil {
		for _, tmpThing := range allThings {
			roomID := tmpThing.GetRoomID()
			tmpRoom, getRoomErr := tm.roomController.GetRoom(roomID)
			if getRoomErr == nil {
				newRoomFullInfo := models.ThingFullInfo{Thing: tmpThing, Room: tmpRoom}
				thingsFullInfo = append(thingsFullInfo, newRoomFullInfo)
			} else {
				err = getRoomErr
				break
			}
		}
	}
	return thingsFullInfo, err
}

func (tm *ThingManager) GetFreeThings(page, size int) ([]models.ThingFullInfo, error) {
	thingsFullInfo := make([]models.ThingFullInfo, 0)
	allThings, err := tm.thingController.GetFreeThings(page, size)
	if err == nil {
		for _, tmpThing := range allThings {
			roomID := tmpThing.GetRoomID()
			tmpRoom, getRoomErr := tm.roomController.GetRoom(roomID)
			if getRoomErr == nil {
				newRoomFullInfo := models.ThingFullInfo{Thing: tmpThing, Room: tmpRoom}
				thingsFullInfo = append(thingsFullInfo, newRoomFullInfo)
			} else {
				err = getRoomErr
				break
			}
		}
	}
	return thingsFullInfo, err
}

func (tm *ThingManager) GetStudentThings(studentNumber string, page, size int) ([]models.ThingFullInfo, error) {
	if studentNumber == objects.EmptyString {
		return nil, appErrors.BadStudentParamsErr
	}

	thingsFullInfo := make([]models.ThingFullInfo, 0)
	studentID, err := tm.studentController.GetStudentIDByNumber(studentNumber)
	if err == nil {
		allThings, getStudentThingsErr := tm.studentController.GetStudentThings(studentID, page, size)
		if getStudentThingsErr == nil {
			for _, tmpThing := range allThings {
				roomID := tmpThing.GetRoomID()
				tmpRoom, getRoomErr := tm.roomController.GetRoom(roomID)
				if getRoomErr == nil {
					newRoomFullInfo := models.ThingFullInfo{Thing: tmpThing, Room: tmpRoom}
					thingsFullInfo = append(thingsFullInfo, newRoomFullInfo)
				} else {
					err = getRoomErr
					break
				}
			}
		} else {
			err = getStudentThingsErr
		}
	}
	return thingsFullInfo, err
}

func (tm *ThingManager) AddNewThing(markNumber int, thingType string) error {
	if thingType == objects.EmptyString || markNumber <= objects.None {
		return appErrors.BadThingParamsErr
	} else {
		return tm.thingController.AddThing(markNumber, thingType)
	}
}

func (tm *ThingManager) TransferThing(markNumber int, roomID int) error {
	if markNumber <= objects.None {
		return appErrors.BadThingParamsErr
	}

	if roomID <= objects.None {
		return appErrors.BadRoomParamsErr
	}

	thingID, err := tm.thingController.GetThingIDByMarkNumber(markNumber)
	if err == nil {
		_, getRoomErr := tm.roomController.GetRoom(roomID)
		if getRoomErr == nil {
			tmpThing, getThingErr := tm.thingController.GetThing(thingID)
			if getThingErr == nil {
				err = tm.thingController.TransferThing(thingID, tmpThing.GetRoomID(), roomID)
			} else {
				err = getThingErr
			}
		} else {
			err = getRoomErr
		}
	}
	return err
}
