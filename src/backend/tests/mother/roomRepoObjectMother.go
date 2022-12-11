package mother

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
)

const (
	Type = "Комната"
)

type RoomRepoObjectMother struct{}

func (m RoomRepoObjectMother) CreateRepo() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func (m RoomRepoObjectMother) CreateDefaultRooms(amount int) []objects.Room {
	resultRooms := make([]objects.Room, objects.Empty)
	roomType := Type
	for i := 1; i <= amount; i++ {
		resultRooms = append(resultRooms, objects.NewRoomWithParams(i, roomType, i))
	}
	return resultRooms
}

func (m RoomRepoObjectMother) CreateRows(rooms []objects.Room) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"roomid", "roomtype", "roomnumber"})
	for _, room := range rooms {
		rows.AddRow(room.GetID(), room.GetRoomType(), room.GetRoomNumber())
	}
	return rows
}

func (m RoomRepoObjectMother) CreateDTORoom() objects.RoomDTO {
	return objects.NewRoomDTO(Type, int(InsertID))
}
