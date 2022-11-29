package roomRepo

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"src/objects"
	"testing"
)

const (
	Type         = "Комната"
	InsertID     = 5
	RowsAffected = 1
)

type RoomRepoObjectMother struct{}

func (m RoomRepoObjectMother) CreateRepo() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func (m RoomRepoObjectMother) CreateDefaultRooms(amount int) []objects.Room {
	resultRooms := make([]objects.Room, objects.Empty)
	roomType := Type
	for i := 1; i <= amount; i++ {
		resultRooms = append(resultRooms, objects.NewRoomWithParams(i, roomType, i))
	}
	return resultRooms
}

func (m RoomRepoObjectMother) CreateRows(rooms []objects.Room) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"roomid", "roomtype", "roomnumber"})
	for _, room := range rooms {
		rows.AddRow(room.GetID(), room.GetRoomType(), room.GetRoomNumber())
	}
	return rows
}

func (m RoomRepoObjectMother) CreateDTORoom() objects.RoomDTO {
	return objects.NewRoomDTO(Type, InsertID)
}

func TestPgRoomRepo_GetRooms(t *testing.T) {
	// Arrange
	objectMother := RoomRepoObjectMother{}
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
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(realRooms, resultRooms) {
		t.Errorf("results not match, want %v, have %v", realRooms, resultRooms)
		return
	}
}

func TestPgRoomRepo_AddRoom(t *testing.T) {
	// Arrange
	objectMother := RoomRepoObjectMother{}
	roomDTO := objectMother.CreateDTORoom()
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("INSERT INTO").WithArgs(roomDTO.GetRoomType(), roomDTO.GetRoomNumber()).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))
	repo := PgRoomRepo{Conn: db}

	// Act
	execErr := repo.AddRoom(roomDTO)

	// Assert
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestPgRoomRepo_GetRoomPositive проверяет, что если комната есть, она успешно вернётся.
func TestPgRoomRepo_GetRoomPositive(t *testing.T) {
	// Arrange
	objectMother := RoomRepoObjectMother{}
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
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(room, realRooms[0]) {
		t.Errorf("results not match, want %v, have %v", realRooms[0], room)
		return
	}
}

// TestPgRoomRepo_GetRoomNegative проверяет, что если комната нет, то вернётся ошибка.
func TestPgRoomRepo_GetRoomNegative(t *testing.T) {
	// Arrange
	objectMother := RoomRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	mock.ExpectQuery("SELECT").WithArgs(InsertID).WillReturnError(RoomNotFoundErr)
	repo := PgRoomRepo{Conn: db}

	// Act
	_, execErr := repo.GetRoom(InsertID)

	// Assert
	if execErr != RoomNotFoundErr {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

func TestPgRoomRepo_DeleteRoom(t *testing.T) {
	// Arrange
	objectMother := RoomRepoObjectMother{}
	ID := 1
	db, mock := objectMother.CreateRepo()
	mock.ExpectExec("DELETE").WithArgs(ID).WillReturnError(nil).
		WillReturnResult(sqlmock.NewResult(int64(ID), RowsAffected))
	repo := PgRoomRepo{Conn: db}

	// Act
	execErr := repo.DeleteRoom(ID)

	// Assert
	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}
