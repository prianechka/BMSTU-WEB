package roomManager

import (
	"src/logic/controllers/roomController"
	"src/objects"
)

type RoomManager struct {
	roomController roomController.RoomController
}

func CreateNewRoomManager(rc roomController.RoomController) *RoomManager {
	return &RoomManager{
		roomController: rc,
	}
}
func (rm *RoomManager) GetAllRooms(page, size int) ([]objects.Room, error) {
	return rm.roomController.GetRooms(page, size)
}

func (rm *RoomManager) GetRoom(id int) (objects.Room, error) {
	return rm.roomController.GetRoom(id)
}
