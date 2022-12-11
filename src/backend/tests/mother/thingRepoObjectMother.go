package mother

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
)

type ThingRepoObjectMother struct{}

var (
	DefaultMarkNumber       = 123
	DefaultThingType        = "Табуретка"
	InsertID          int64 = 5
	RowsAffected      int64 = 1
)

func (m ThingRepoObjectMother) CreateRepo() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func (m ThingRepoObjectMother) CreateDefaultThings(amount int) []objects.Thing {
	resultStudents := make([]objects.Thing, objects.Empty)
	for i := 1; i <= amount; i++ {
		resultStudents = append(resultStudents, objects.NewThingWithParams(i, DefaultMarkNumber+i,
			DefaultThingType, i, i))
	}
	return resultStudents
}

func (m ThingRepoObjectMother) CreateRows(things []objects.Thing) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"thingid", "marknumber", "thingtype", "ownerID", "roomid"})
	for _, thing := range things {
		rows.AddRow(thing.GetID(), thing.GetMarkNumber(), thing.GetThingType(), thing.GetOwnerID(),
			thing.GetRoomID())
	}
	return rows
}

func (m ThingRepoObjectMother) CreateRowForID(id int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"thingid"})
	rows.AddRow(id)
	return rows
}

func (m ThingRepoObjectMother) CreateThingDTO() objects.ThingDTO {
	return objects.NewThingDTO(DefaultMarkNumber, DefaultThingType)
}
