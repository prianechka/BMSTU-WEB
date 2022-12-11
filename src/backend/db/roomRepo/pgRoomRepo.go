package roomRepo

import (
	"database/sql"
	"src/db/sql"
	"src/objects"
	"strconv"
)

type PgRoomRepo struct {
	Conn *sql.DB
}

func (pg *PgRoomRepo) AddRoom(room objects.RoomDTO) error {
	sqlString := pgsql.PostgreSQLAddRoom{}.GetString()
	_, err := pg.Conn.Exec(sqlString, room.GetRoomType(), room.GetRoomNumber())
	return err
}

func (pg *PgRoomRepo) GetRooms(page, size int) ([]objects.Room, error) {
	var (
		resultRooms    = make([]objects.Room, objects.Empty)
		id, roomNumber int
		roomType       string
		err            error
		sizeParam      string
	)
	sqlString := pgsql.PostgreSQLGetRooms{}.GetString()
	if size != objects.Null {
		sizeParam = strconv.Itoa(size)
	} else {
		sizeParam = "ALL"
	}
	rows, execError := pg.Conn.Query(sqlString, sizeParam, page*size)
	if execError == nil {
		for rows.Next() {
			scanErr := rows.Scan(&id, &roomType, &roomNumber)
			if scanErr == nil {
				tmpRoom := objects.NewRoomWithParams(id, roomType, roomNumber)
				resultRooms = append(resultRooms, tmpRoom)
			}
		}
	} else {
		err = execError
	}
	return resultRooms, err
}

func (pg *PgRoomRepo) GetRoom(id int) (objects.Room, error) {
	var (
		resultRoom         = objects.NewEmptyRoom()
		roomID, roomNumber int
		roomType           string
		err                error
	)
	sqlString := pgsql.PostgreSQLGetRoom{}.GetString()
	rows, execError := pg.Conn.Query(sqlString, id)
	if execError == nil {
		for rows.Next() {
			scanErr := rows.Scan(&roomID, &roomType, &roomNumber)
			if scanErr == nil {
				resultRoom = objects.NewRoomWithParams(id, roomType, roomNumber)
			}
		}
	} else {
		err = execError
	}
	return resultRoom, err
}

func (pg *PgRoomRepo) GetRoomThings(id int) ([]objects.Thing, error) {
	var (
		resultThings                         = make([]objects.Thing, objects.Empty)
		thingID, markNumber, ownerID, roomID int
		thingType                            string
		err                                  error
	)
	sqlString := pgsql.PostgreSQLGetRoomThings{}.GetString()
	rows, execError := pg.Conn.Query(sqlString, id)
	if execError == nil {
		for rows.Next() {
			readRowErr := rows.Scan(&thingID, &markNumber, &thingType, &ownerID, &roomID)
			if readRowErr == nil {
				tmpThings := objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
				resultThings = append(resultThings, tmpThings)
			} else {
				err = readRowErr
				break
			}
		}
	} else {
		err = execError
	}
	return resultThings, err
}

func (pg *PgRoomRepo) DeleteRoom(id int) error {
	sqlString := pgsql.PostgreSQLDeleteRoom{}.GetString()
	_, err := pg.Conn.Exec(sqlString, id)
	return err
}
