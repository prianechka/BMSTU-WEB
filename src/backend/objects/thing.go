package objects

type Thing struct {
	id         int
	markNumber int
	thingType  string
	ownerID    int
	roomID     int
}

type ThingDTO struct {
	markNumber int
	thingType  string
}

type ThingResponseDTO struct {
	Marknumber int    `json:"mark-number"`
	ThingType  string `json:"thing-type"`
	RoomID     int    `json:"room-id"`
}

func NewThingWithParams(id, markNumber int, thingType string, ownerID, roomID int) Thing {
	return Thing{
		id:         id,
		markNumber: markNumber,
		thingType:  thingType,
		ownerID:    ownerID,
		roomID:     roomID,
	}
}

func NewEmptyThing() Thing {
	return Thing{id: None}
}

func NewThingDTO(markNumber int, thingType string) ThingDTO {
	return ThingDTO{
		markNumber: markNumber,
		thingType:  thingType,
	}
}

func (t *Thing) GetID() int {
	return t.id
}

func (t *Thing) GetMarkNumber() int {
	return t.markNumber
}

func (t *Thing) GetThingType() string {
	return t.thingType
}

func (t *Thing) GetOwnerID() int {
	return t.ownerID
}

func (t *Thing) GetRoomID() int {
	return t.roomID
}

func (t *Thing) SetRoomID(id int) {
	t.roomID = id
}

func (t *Thing) SetOwnerID(id int) {
	t.ownerID = id
}

func (t *ThingDTO) GetMarkNumber() int {
	return t.markNumber
}

func (t *ThingDTO) GetThingType() string {
	return t.thingType
}

func CreateThingResponse(things Thing) ThingResponseDTO {
	return ThingResponseDTO{
		Marknumber: things.GetMarkNumber(),
		ThingType:  things.GetThingType(),
		RoomID:     things.GetRoomID(),
	}
}
