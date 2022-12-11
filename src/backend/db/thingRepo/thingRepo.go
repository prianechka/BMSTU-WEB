package thingRepo

import "src/objects"

type ThingRepo interface {
	AddThing(thing objects.ThingDTO) error
	GetThings(page, size int) ([]objects.Thing, error)
	DeleteThing(id int) error
	GetThing(id int) (objects.Thing, error)
	TransferThingRoom(id, srcRoomID int, dstRoomID int) error
	GetThingIDByMarkNumber(markNumber int) (int, error)
}
