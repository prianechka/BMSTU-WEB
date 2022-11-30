package roomController

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/roomRepo"
	"src/tests"
	"src/tests/mother"
	"testing"
)

const (
	InsertID          = 5
	RowsAffected      = 1
	DefaultRoomType   = "комната"
	DefaultRoomNumber = 1
)

func TestRoomController_AddRoom(t *testing.T) {
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

func TestRoomController_DeleteRoomPositive(t *testing.T) {
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
func TestRoomController_DeleteRoomNegative(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	execErr := controller.DeleteRoom(ID)

	tests.AssertErrors(t, execErr, RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestRoomController_GetRoomPositive(t *testing.T) {
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

func TestRoomController_GetRoomNegative(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	_, execErr := controller.GetRoom(ID)

	tests.AssertErrors(t, execErr, RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestRoomController_GetRooms(t *testing.T) {
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

func TestRoomController_GetRoomThingsPositive(t *testing.T) {
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

func TestRoomController_GetRoomThingsNegative(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := RoomController{Repo: &repo}

	// Act
	_, execErr := controller.GetRoomThings(ID)

	tests.AssertErrors(t, execErr, RoomNotFoundErr)
	tests.AssertMocks(t, mock)
}
