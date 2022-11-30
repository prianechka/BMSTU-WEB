package thingRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/tests"
	mot "src/tests/mother"
	"testing"
)

var (
	InsertID     int64 = 5
	RowsAffected int64 = 1
)

func TestPgThingRepo_AddThing(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	thingDTO := objectMother.CreateThingDTO()
	mock.ExpectExec("INSERT INTO").WithArgs(thingDTO.GetMarkNumber(), thingDTO.GetThingType()).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgThingRepo{Conn: db}

	// Act
	execErr := repo.AddThing(thingDTO)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgThingRepo_DeleteThing(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("DELETE").WithArgs(InsertID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgThingRepo{Conn: db}

	// Act
	execErr := repo.DeleteThing(int(InsertID))

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgThingRepo_GetThings(t *testing.T) {
	objectMother := mot.ThingRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows).WillReturnError(nil)

	repo := PgThingRepo{Conn: db}

	// Act
	resultThings, execErr := repo.GetThings()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, realThings, resultThings)
}

func TestPgThingRepo_TransferThingRoom(t *testing.T) {
	// Arrange
	var (
		id        = int(InsertID)
		srcIDRoom = int(InsertID)
		dstIDRoom = int(InsertID)
	)
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(id, srcIDRoom, dstIDRoom).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgThingRepo{Conn: db}

	// Act
	execErr := repo.TransferThingRoom(id, srcIDRoom, dstIDRoom)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgThingRepo_GetThingPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgThingRepo{Conn: db}

	// Act
	thing, execErr := repo.GetThing(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, thing, realThings[0])
}

func TestPgThingRepo_GetThingNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(sql.ErrNoRows)
	repo := PgThingRepo{Conn: db}

	// Act
	_, execErr := repo.GetThing(int(InsertID))

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}

func TestPgThingRepo_GetThingIDByMarkNumberPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	id := 1
	rows := objectMother.CreateRowForID(id)
	mock.ExpectQuery("SELECT").WithArgs(mot.DefaultMarkNumber).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgThingRepo{Conn: db}

	// Act
	thingID, execErr := repo.GetThingIDByMarkNumber(mot.DefaultMarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, thingID, id)
}

func TestPgThingRepo_GetThingIDByMarkNumberNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(mot.DefaultMarkNumber).
		WillReturnError(sql.ErrNoRows)
	repo := PgThingRepo{Conn: db}

	// Act
	_, execErr := repo.GetThingIDByMarkNumber(mot.DefaultMarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}
