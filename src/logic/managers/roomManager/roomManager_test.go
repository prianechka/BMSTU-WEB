package roomManager

import (
	"src/db/roomRepo"
	"src/logic/controllers/roomController"
	"src/tests"
	"src/tests/mother"
	"testing"
)

func TestRoomManager_GetAllRooms(t *testing.T) {
	// Arrange
	objectMother := mother.RoomRepoObjectMother{}
	N := 3
	db, mock := objectMother.CreateRepo()
	realRooms := objectMother.CreateDefaultRooms(N)
	rows := objectMother.CreateRows(realRooms)
	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)
	repo := roomRepo.PgRoomRepo{Conn: db}
	controller := roomController.RoomController{Repo: &repo}
	manager := RoomManager{roomController: controller}

	// Act
	resultRooms, execErr := manager.GetAllRooms()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, resultRooms, realRooms)
}
