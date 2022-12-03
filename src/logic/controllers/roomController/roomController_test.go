package roomController

import (
	"database/sql"
	"github.com/bloomberg/go-testgroup"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/roomRepo"
	"src/tests"
	"src/tests/mother"
	appErrors "src/utils/error"
	"testing"
)

const (
	InsertID          = 5
	RowsAffected      = 1
	DefaultRoomType   = "комната"
	DefaultRoomNumber = 1
)

type TestRoomController struct{}

func Test_RoomController(t *testing.T) {
	testgroup.RunSerially(t, &TestRoomController{})
}

func (*TestRoomController) TestRoomController_AddRoom(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(DefaultRoomType, DefaultRoomNumber).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	execErr := controller.AddRoom(DefaultRoomType, DefaultRoomNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func (*TestRoomController) TestRoomController_DeleteRoomPositive(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	N := 1
	realRooms := objectMother.CreateDefaultRooms(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(rows)
	mock.ExpectExec("DELETE").WithArgs(ID).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(int64(ID), RowsAffected))
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	execErr := controller.DeleteRoom(ID)

	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

// TestRoomController_DeleteRoomNegative проверяет, что если комнаты не существует, то удаления не произойдет.
func (*TestRoomController) TestRoomController_DeleteRoomNegative(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	execErr := controller.DeleteRoom(ID)

	tests.AssertErrors(t, execErr, appErrors.RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestRoomController) TestRoomController_GetRoomPositive(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	N := 1
	realRooms := objectMother.CreateDefaultRooms(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(rows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	room, execErr := controller.GetRoom(ID)

	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, room, realRooms[0])
}

func (*TestRoomController) TestRoomController_GetRoomNegative(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	_, execErr := controller.GetRoom(ID)

	tests.AssertErrors(t, execErr, appErrors.RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func (*TestRoomController) TestRoomController_GetRooms(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realRooms := objectMother.CreateDefaultRooms(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	resultRooms, execErr := controller.GetRooms()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, resultRooms, realRooms)
}

func (*TestRoomController) TestRoomController_GetRoomThingsPositive(t *testgroup.T) {
	// Arrange
	roomObjectMother := mother.RoomRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}
	db, mock := roomObjectMother.CreateRepo()
	N := 1
	id := 1
	realRooms := roomObjectMother.CreateDefaultRooms(N)
	roomRows := roomObjectMother.CreateRows(realRooms)

	realThings := thingObjectMother.CreateDefaultThings(N + 1)
	thingRows := thingObjectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(nil).WillReturnRows(roomRows)
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(nil).WillReturnRows(thingRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	things, execErr := controller.GetRoomThings(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, things, realThings)
}

func (*TestRoomController) TestRoomController_GetRoomThingsNegative(t *testgroup.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	_, execErr := controller.GetRoomThings(ID)

	tests.AssertErrors(t, execErr, appErrors.RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}
