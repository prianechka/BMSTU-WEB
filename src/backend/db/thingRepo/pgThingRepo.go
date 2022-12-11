package thingRepo

import (
	"database/sql"
	"log"
	"src/db/sql"
	"src/objects"
	"strconv"
)

type PgThingRepo struct {
	Conn *sql.DB
}

func (pg *PgThingRepo) AddThing(thing objects.ThingDTO) error {
	sqlString := pgsql.PostgreSQLAddThing{}.GetString()
	_, err := pg.Conn.Exec(sqlString, thing.GetMarkNumber(), thing.GetThingType())
	return err
}

func (pg *PgThingRepo) GetThings(page, size int) ([]objects.Thing, error) {
	var (
		resultThings                         = make([]objects.Thing, objects.Empty)
		thingID, markNumber, ownerID, roomID int
		thingType                            string
		err                                  error
		sizeParam                            string
	)
	sqlString := pgsql.PostgreSQLGetThings{}.GetString()
	if size != objects.Null {
		sizeParam = strconv.Itoa(size)
	} else {
		sizeParam = "ALL"
	}

	rows, execError := pg.Conn.Query(sqlString, sizeParam, page*size)
	if execError == nil {
		for rows.Next() {
			readRowErr := rows.Scan(&thingID, &markNumber, &thingType, &ownerID, &roomID)
			if err == nil {
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

func (pg *PgThingRepo) DeleteThing(id int) error {
	sqlString := pgsql.PostgreSQLDeleteThing{}.GetString()
	_, err := pg.Conn.Exec(sqlString, id)
	return err
}

func (pg *PgThingRepo) GetThing(id int) (objects.Thing, error) {
	var (
		thing                                = objects.NewEmptyThing()
		thingID, markNumber, ownerID, roomID int
		thingType                            string
		err                                  error
	)
	sqlString := pgsql.PostgreSQLGetThing{}.GetString()
	rows, execError := pg.Conn.Query(sqlString, id)
	if execError == nil {
		for rows.Next() {
			readRowErr := rows.Scan(&thingID, &markNumber, &thingType, &ownerID, &roomID)
			if err == nil {
				thing = objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
			} else {
				err = readRowErr
				break
			}
		}
	} else {
		err = execError
	}
	return thing, err
}

func (pg *PgThingRepo) TransferThingRoom(id, srcRoomID, dstRoomID int) error {
	sqlString := pgsql.PostgreSQLTransferThingRoom{}.GetString()
	log.Println(id, srcRoomID, dstRoomID)
	_, err := pg.Conn.Exec(sqlString, srcRoomID, dstRoomID, id)
	return err
}

func (pg *PgThingRepo) GetThingIDByMarkNumber(markNumber int) (int, error) {
	var result int
	sqlString := pgsql.PostgreSQLGetThingID{}.GetString()
	row := pg.Conn.QueryRow(sqlString, markNumber)
	err := row.Scan(&result)
	return result, err
}
