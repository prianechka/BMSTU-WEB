package roomRepo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"src/db/sql"
	"src/objects"
)

type PgRoomRepo struct {
	ConnectParams objects.PgConnection
}

func (pg *PgRoomRepo) AddRoom(room objects.RoomDTO) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLAddRoom{}.GetString(room)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgRoomRepo) GetRooms() ([]objects.Room, error) {
	var resultRooms = make([]objects.Room, objects.Empty)
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetRooms{}.GetString()
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					id := int(values[0].(int32))
					roomType := values[1].(string)
					roomNumber := int(values[2].(int32))
					tmpRoom := objects.NewRoomWithParams(id, roomType, roomNumber)
					resultRooms = append(resultRooms, tmpRoom)
				}
			}
		}
	}
	return resultRooms, err
}

func (pg *PgRoomRepo) GetRoom(id int) (objects.Room, error) {
	var resultRoom = objects.NewEmptyRoom()
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetRoom{}.GetString(id)
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					roomID := int(values[0].(int32))
					roomType := values[1].(string)
					roomNumber := int(values[2].(int32))
					resultRoom = objects.NewRoomWithParams(roomID, roomType, roomNumber)
				}
			}
		}
	}
	return resultRoom, err
}

func (pg *PgRoomRepo) GetRoomThings(id int) ([]objects.Thing, error) {
	var resultRoomThings = make([]objects.Thing, objects.Empty)
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetRoomThings{}.GetString(id)
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					thingID := int(values[0].(int32))
					markNumber := int(values[1].(int32))
					thingType := values[2].(string)
					ownerID := int(values[3].(int32))
					roomID := int(values[4].(int32))

					tmpRoom := objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
					resultRoomThings = append(resultRoomThings, tmpRoom)
				}
			}
		}
	}
	return resultRoomThings, err
}

func (pg *PgRoomRepo) DeleteRoom(id int) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLDeleteRoom{}.GetString(id)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}
