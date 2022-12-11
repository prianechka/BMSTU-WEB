package roomRepo

import "src/objects"

type RoomRepo interface {
	AddRoom(room objects.RoomDTO) error
	GetRooms(page, size int) ([]objects.Room, error)
	GetRoom(id int) (objects.Room, error)
	GetRoomThings(id int) ([]objects.Thing, error)
	DeleteRoom(id int) error
}
