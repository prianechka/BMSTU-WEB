package thingRepo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"src/db/sql"
	"src/objects"
)

type PgThingRepo struct {
	ConnectParams objects.PgConnection
}

func (pg *PgThingRepo) AddThing(thing objects.ThingDTO) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLAddThing{}.GetString(thing)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgThingRepo) GetThings() ([]objects.Thing, error) {
	var resultThings = make([]objects.Thing, objects.Empty)
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetThings{}.GetString()
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

					tmpThings := objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
					resultThings = append(resultThings, tmpThings)
				} else {
					err = readRowError
				}
			}
		} else {
			err = execError
		}
	}
	return resultThings, err
}

func (pg *PgThingRepo) DeleteThing(id int) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLDeleteThing{}.GetString(id)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgThingRepo) GetThing(id int) (objects.Thing, error) {
	var thing = objects.NewEmptyThing()
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetThing{}.GetString(id)
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
					thing = objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
				} else {
					err = readRowError
				}
			}
		} else {
			err = execError
		}
	}
	return thing, err
}

func (pg *PgThingRepo) TransferThingRoom(id, srcRoomID, dstRoomID int) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLTransferThingRoom{}.GetString(id, srcRoomID, dstRoomID)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgThingRepo) GetThingIDByMarkNumber(markNumber int) (int, error) {
	var result int
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetThingID{}.GetString(markNumber)
		row := conn.QueryRow(context.Background(), sqlString)
		err = row.Scan(&result)
	}
	return result, err
}
