package objects

type Room struct {
	id         int
	roomType   string
	roomNumber int
}

type RoomResponseDTO struct {
	RoomID     int    `json:"room-id"`
	RoomType   string `json:"room-type"`
	RoomNumber int    `json:"room-number"`
}

func NewRoomWithParams(id int, roomType string, roomNumber int) Room {
	return Room{id, roomType, roomNumber}
}

func NewEmptyRoom() Room {
	return Room{id: None}
}

func (r *Room) GetID() int {
	return r.id
}

func (r *Room) GetRoomType() string {
	return r.roomType
}

func (r *Room) GetRoomNumber() int {
	return r.roomNumber
}

type RoomDTO struct {
	roomType   string
	roomNumber int
}

func NewRoomDTO(roomType string, roomNumber int) RoomDTO {
	return RoomDTO{
		roomType:   roomType,
		roomNumber: roomNumber,
	}
}

func (rd *RoomDTO) GetRoomType() string {
	return rd.roomType
}

func (rd *RoomDTO) GetRoomNumber() int {
	return rd.roomNumber
}

func CreateRoomResponseDTO(arr Room) RoomResponseDTO {
	return RoomResponseDTO{
		RoomID:     arr.GetID(),
		RoomType:   arr.GetRoomType(),
		RoomNumber: arr.GetRoomNumber(),
	}
}

func CreateRoomResponse(room Room) RoomResponseDTO {
	return RoomResponseDTO{
		RoomID:     room.GetID(),
		RoomType:   room.GetRoomType(),
		RoomNumber: room.GetRoomNumber(),
	}
}

func CreateRoomResponseArr(arr []Room) []RoomResponseDTO {
	result := make([]RoomResponseDTO, Empty)
	for _, room := range arr {
		result = append(result, RoomResponseDTO{
			RoomID:     room.GetID(),
			RoomType:   room.GetRoomType(),
			RoomNumber: room.GetRoomNumber(),
		})
	}
	return result
}
