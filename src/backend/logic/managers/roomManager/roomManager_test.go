package roomManager

import (
	"github.com/bloomberg/go-testgroup"
	"src/db/roomRepo"
	"src/logic/controllers/roomController"
	"src/tests"
	"src/tests/mother"
	"testing"
	"time"
)

type TestRoomManager struct{}

func Test_RoomManager(t *testing.T) {
	testgroup.RunSerially(t, &TestRoomManager{})
}

func (*TestRoomManager) TestRoomManager_GetAllRooms(t *testgroup.T) {
	defer tests.TimeTrack(time.Now())
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
