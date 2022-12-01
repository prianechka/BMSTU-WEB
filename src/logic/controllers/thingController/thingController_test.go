package thingController

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/db/thingRepo"
	"src/objects"
	"src/tests"
	mot "src/tests/mother"
	"testing"
)

const (
	InsertID     = 5
	RowsAffected = 1
)

func TestThingController_AddThingPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N          = 3
		MarkNumber = mot.DefaultMarkNumber + 12
		ThingType  = mot.DefaultThingType
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)
	mock.ExpectExec("INSERT INTO").WithArgs(MarkNumber, ThingType).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.AddThing(MarkNumber, ThingType)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestThingController_AddThingNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N          = 3
		MarkNumber = mot.DefaultMarkNumber + N - 1
		ThingType  = mot.DefaultThingType
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.AddThing(MarkNumber, ThingType)

	// Assert
	tests.AssertErrors(t, execErr, ThingAlreadyExistErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_DeleteThingPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N  = 3
		ID = 1
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(rows)
	mock.ExpectExec("DELETE").WithArgs(ID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.DeleteThing(ID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestThingController_DeleteThingNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		ID = 5
	)

	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.DeleteThing(ID)

	// Assert
	tests.AssertErrors(t, execErr, ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_GetThings(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 3
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	things, execErr := controller.GetThings()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, realThings, things)
}

func TestThingController_GetThingPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N  = 1
		ID = 1
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	thing, execErr := controller.GetThing(ID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, thing, realThings[0])
}

func TestThingController_GetThingNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		ID = 1
	)

	mock.ExpectQuery("SELECT").WithArgs(ID).WillReturnError(sql.ErrNoRows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	_, execErr := controller.GetThing(ID)

	// Assert
	tests.AssertErrors(t, execErr, ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_GetThingIDByMarkNumberPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		ID         = 1
		MarkNumber = mot.DefaultMarkNumber
	)
	rows := objectMother.CreateRowForID(ID)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	thingID, execErr := controller.GetThingIDByMarkNumber(MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, thingID, ID)
}

func TestThingController_GetThingIDByMarkNumberNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		MarkNumber = mot.DefaultMarkNumber
	)

	mock.ExpectQuery("SELECT").WithArgs(MarkNumber).WillReturnError(sql.ErrNoRows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	_, execErr := controller.GetThingIDByMarkNumber(MarkNumber)

	// Assert
	tests.AssertErrors(t, execErr, ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_GetThingRoomPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		roomID  = 1
		N       = 1
		thingID = 1
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	realRoomID, execErr := controller.GetThingRoom(thingID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, roomID, realRoomID)
}

func TestThingController_GetThingRoomNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		thingID = 1
	)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(sql.ErrNoRows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	_, execErr := controller.GetThingRoom(thingID)

	// Assert
	tests.AssertErrors(t, execErr, ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_GetCurrentOwnerPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		ownerID = 1
		N       = 1
		thingID = 1
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	realOwnerID, execErr := controller.GetCurrentOwner(thingID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, ownerID, realOwnerID)
}

func TestThingController_GetCurrentOwnerNegative(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		thingID = 1
	)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(sql.ErrNoRows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	_, execErr := controller.GetCurrentOwner(thingID)

	// Assert
	tests.AssertErrors(t, execErr, ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_GetFreeThings(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	N := 3
	realThings := objectMother.CreateDefaultThings(N)
	realThings[0].SetOwnerID(objects.None)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WillReturnRows(rows).WillReturnError(nil)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	things, execErr := controller.GetFreeThings()

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
	tests.AssertResult(t, things, realThings[0:1])
}

func TestThingController_TransferThingPositive(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N         = 1
		thingID   = 1
		srcRoomID = 1
		dstRoomID = 2
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(nil).WillReturnRows(rows)
	mock.ExpectExec("INSERT").WithArgs(thingID, srcRoomID, dstRoomID).
		WillReturnError(nil).WillReturnResult(sqlmock.NewResult(InsertID, RowsAffected))

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.TransferThing(thingID, srcRoomID, dstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, nil)
	tests.AssertMocks(t, mock)
}

func TestThingController_TransferThingNegativeThingNotFound(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		thingID   = 1
		srcRoomID = 1
		dstRoomID = 2
	)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(sql.ErrNoRows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.TransferThing(thingID, srcRoomID, dstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, ThingNotFoundErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_TransferThingNegativeBadSrcID(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N         = 1
		thingID   = 1
		srcRoomID = 3
		dstRoomID = 2
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.TransferThing(thingID, srcRoomID, dstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, BadSrcRoomErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_TransferThingNegativeBadDstID(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N         = 1
		thingID   = 1
		srcRoomID = 1
		dstRoomID = objects.None
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.TransferThing(thingID, srcRoomID, dstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, BadDstRoomErr)
	tests.AssertMocks(t, mock)
}

func TestThingController_TransferThingNegativeEqualSrcDstID(t *testing.T) {
	// Arrange
	objectMother := mot.ThingRepoObjectMother{}
	db, mock := objectMother.CreateRepo()
	var (
		N         = 1
		thingID   = 1
		srcRoomID = 1
		dstRoomID = 1
	)
	realThings := objectMother.CreateDefaultThings(N)
	rows := objectMother.CreateRows(realThings)

	mock.ExpectQuery("SELECT").WithArgs(thingID).WillReturnError(nil).WillReturnRows(rows)

	repo := thingRepo.PgThingRepo{Conn: db}
	controller := ThingController{Repo: &repo}

	// Act
	execErr := controller.TransferThing(thingID, srcRoomID, dstRoomID)

	// Assert
	tests.AssertErrors(t, execErr, BadDstRoomErr)
	tests.AssertMocks(t, mock)
}
