package roomController

import (
	"src/db/roomRepo"
	"src/objects"
	appErrors "src/utils/error"
)

type RoomController struct {
	Repo roomRepo.RoomRepo
}

func (rc *RoomController) AddRoom(roomType string, number int) error {
	roomDTO := objects.NewRoomDTO(roomType, number)
	err := rc.Repo.AddRoom(roomDTO)
	return err
}

func (rc *RoomController) GetRooms(page, size int) ([]objects.Room, error) {
	return rc.Repo.GetRooms(page, size)
}

func (rc *RoomController) GetRoom(id int) (objects.Room, error) {
	tmpRoom, err := rc.Repo.GetRoom(id)
	if err != nil || tmpRoom.GetID() == objects.None {
		err = appErrors.RoomNotFoundErr
	}
	return tmpRoom, err
}

func (rc *RoomController) DeleteRoom(id int) error {
	tmpRoom, err := rc.Repo.GetRoom(id)
	if err == nil && tmpRoom.GetID() > objects.None {
		err = rc.Repo.DeleteRoom(id)
	} else {
		err = appErrors.RoomNotFoundErr
	}
	return err
}

func (rc *RoomController) GetRoomThings(id int) ([]objects.Thing, error) {
	tmpRoom, err := rc.Repo.GetRoom(id)
	if err == nil {
		if tmpRoom.GetID() != objects.None {
			return rc.Repo.GetRoomThings(id)
		}
	} else {
		err = appErrors.RoomNotFoundErr
	}
	return nil, err
}
