package thingRepo

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/tests"
	mot "src/tests/mother"
	"testing"
	"time"
)

var (
	InsertID     int64 = 5
	RowsAffected int64 = 1
)

type TestPgThingRepo struct{}

func Test_PgThingRepo(t *testing.T) {
	testgroup.RunSerially(t, &TestPgThingRepo{})
}

func (*TestPgThingRepo) TestPgThingRepo_AddThing(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgThingRepo) TestPgThingRepo_DeleteThing(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgThingRepo) TestPgThingRepo_GetThings(t *testgroup.T) {
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

func (*TestPgThingRepo) TestPgThingRepo_TransferThingRoom(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
	// Arrange

	var (
		id        = int(InsertID)
		srcIDRoom = int(InsertID)
		dstIDRoom = int(InsertID)
	)
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(srcIDRoom, dstIDRoom, id).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := PgThingRepo{Conn: db}

	// Act
	execErr := repo.TransferThingRoom(id, srcIDRoom, dstIDRoom)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestPgThingRepo) TestPgThingRepo_GetThingPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgThingRepo) TestPgThingRepo_GetThingNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgThingRepo) TestPgThingRepo_GetThingIDByMarkNumberPositive(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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

func (*TestPgThingRepo) TestPgThingRepo_GetThingIDByMarkNumberNegative(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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
