package roomManager

import (
	"src/logic/controllers/roomController"
	"src/objects"
)

type RoomManager struct {
	roomController roomController.RoomController
}

func (rm *RoomManager) GetAllRooms() ([]objects.Room, error) {
	return rm.roomController.GetRooms()
}
