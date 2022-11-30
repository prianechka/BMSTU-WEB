package roomRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/tests"
	"src/tests/mother"
	"testing"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

func TestPgRoomRepo_GetRooms(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realRooms := objectMother.CreateDefaultRooms(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows).WillReturnError(nil)
	repo := PgRoomRepo{Conn: db}

	// Act
	resultRooms, execErr := repo.GetRooms()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, resultRooms, realRooms)
}

func TestPgRoomRepo_AddRoom(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	roomDTO := objectMother.CreateDTORoom()
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(roomDTO.GetRoomType(), roomDTO.GetRoomNumber()).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))
	repo := PgRoomRepo{Conn: db}

	// Act
	execErr := repo.AddRoom(roomDTO)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

// TestPgRoomRepo_GetRoomPositive проверяет, что если комната есть, она успешно вернётся.
func TestPgRoomRepo_GetRoomPositive(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 1
	id := 1
	realRooms := objectMother.CreateDefaultRooms(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgRoomRepo{Conn: db}

	// Act
	room, execErr := repo.GetRoom(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, room, realRooms[0])
}

// TestPgRoomRepo_GetRoomNegative проверяет, что если комнаты нет, то вернётся ошибка.
func TestPgRoomRepo_GetRoomNegative(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(sql.ErrNoRows)
	repo := PgRoomRepo{Conn: db}

	// Act
	_, execErr := repo.GetRoom(InsertID)

	// Assert
	tests.AssertErrors(t, execErr, sql.ErrNoRows)
	tests.AssertMocks(t, mock)
}

func TestPgRoomRepo_DeleteRoom(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("DELETE").WithArgs(ID).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(int64(ID), RowsAffected))
	repo := PgRoomRepo{Conn: db}

	// Act
	execErr := repo.DeleteRoom(ID)

	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestPgRoomRepo_GetRoomThings(t *testing.T) {
	// Arrange
	roomObjectMother := mother.RoomRepoObjectMother{}
	thingObjectMother := mother.ThingRepoObjectMother{}
	db, mock := roomObjectMother.CreateRepo()
	N := 1
	id := 1
	realThings := thingObjectMother.CreateDefaultThings(N)
	rows := thingObjectMother.CreateRows(realThings)
	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(nil).WillReturnRows(rows)
	repo := PgRoomRepo{Conn: db}

	// Act
	things, execErr := repo.GetRoomThings(id)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, things, realThings)
}
