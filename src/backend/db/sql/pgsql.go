package pgsql

import (
	"fmt"
)

type PostgreSQLChangeStudent struct{}
type PostgreSQLGetStudent struct{}
type PostgreSQLGetStudentID struct{}
type PostgreSQLGetStudentsThings struct{}
type PostgreSQLGetAllStudents struct{}
type PostgreSQLAddStudent struct{}
type PostgreSQLTransferStudent struct{}
type PostgreSQLTransferThing struct{}
type PostgreSQLAddRoom struct{}
type PostgreSQLGetRooms struct{}
type PostgreSQLGetRoom struct{}
type PostgreSQLGetRoomThings struct{}
type PostgreSQLDeleteRoom struct{}
type PostgreSQLTransferThingRoom struct{}
type PostgreSQLAddThing struct{}
type PostgreSQLGetThings struct{}
type PostgreSQLGetThing struct{}
type PostgreSQLGetThingID struct{}
type PostgreSQLDeleteThing struct{}
type PostgreSQLGetUserId struct{}
type PostgreSQLGetUser struct{}
type PostgreSQLAddUser struct{}

func (pg PostgreSQLChangeStudent) GetString() string {
	return "UPDATE  Student SET StudentName = $1, StudentSurname = $2, " +
		"StudentGroup = $3, StudentNumber = $4 WHERE StudentID = $5;"
}

func (pg PostgreSQLGetStudent) GetString() string {
	return "SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, " +
		"S.studentgroup, S.studentnumber,  FindStudentRoom(S.studentid) " +
		"FROM  Student as S WHERE S.studentid = $1;"
}
func (pg PostgreSQLGetStudentID) GetString() string {
	return "SELECT S.studentid FROM  Student as S WHERE StudentNumber = $1;"
}

func (pg PostgreSQLGetStudentsThings) GetString() string {
	return "SELECT T.thingid, T.marknumber, T.thingtype,  FindStudent(T.thingId), " +
		" FindRoom(T.thingid) FROM  Thing as T WHERE  FindStudent(T.ThingID) = $1 ORDER BY T.thingid LIMIT $2 OFFSET $3;"
}
func (pg PostgreSQLGetAllStudents) GetWithParamsString() string {
	return "SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, " +
		"S.studentgroup, S.studentnumber, " +
		" FindStudentRoom(S.studentid) " +
		"FROM  Student as S ORDER BY S.studentid LIMIT $1 OFFSET $2;"
}

func (pg PostgreSQLGetAllStudents) GetEmptyString() string {
	return "SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, " +
		"S.studentgroup, S.studentnumber, " +
		" FindStudentRoom(S.studentid) " +
		"FROM  Student as S;"
}

func (pg PostgreSQLAddStudent) GetString() string {
	return "INSERT INTO  Student(studentname, studentsurname, studentgroup, " +
		"studentnumber, settledate, webaccid) VALUES ($1, $2, $3, $4, current_date, $5);"
}
func (pg PostgreSQLTransferStudent) GetString() string {
	return "INSERT INTO  StudentRoomHistory (studentid, roomid, direction, transferdate) " +
		"VALUES ($1, $2, $3, current_date);"
}
func (pg PostgreSQLTransferThing) GetString() string {
	return "INSERT INTO  StudentThingHistory (studentid, thingid, direction, transferdate)" +
		"VALUES ($1, $2, $3, current_date);"
}

func (pg PostgreSQLAddRoom) GetString() string {
	return "INSERT INTO  Rooms(roomtype, roomnumber) VALUES ('$1', $2)"
}

func (pg PostgreSQLGetRooms) GetString() string {
	return "SELECT roomid, roomtype, roomnumber FROM rooms LIMIT $1 OFFSET $2;"
}

func (pg PostgreSQLGetRoom) GetString() string {
	return "SELECT roomid, roomtype, roomnumber FROM  rooms WHERE RoomID = $1;"
}

func (pg PostgreSQLGetRoomThings) GetString() string {
	return "SELECT T.thingid, T.marknumber, T.thingtype,  FindStudent(T.thingId), " +
		" FindRoom(T.thingid) FROM  Thing as T WHERE  FindRoom(T.ThingID) = $1"
}
func (pg PostgreSQLDeleteRoom) GetString() string {
	return "DELETE FROM  rooms WHERE RoomID = $1;"
}

func (pg PostgreSQLTransferThingRoom) GetString() string {
	return "INSERT INTO  ThingRoomHistory (srcroomid, dstroomid, thingid, transferdate) VALUES " +
		"($1, $2, $3, current_date);"
}

func (pg PostgreSQLAddThing) GetString() string {
	return "INSERT INTO  Thing(marknumber, creationdate, thingtype) VALUES  " +
		"($1, current_date, $2);"
}

func (pg PostgreSQLGetThings) GetString() string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, " +
		" FindStudent(T.thingId),  FindRoom(T.thingid) FROM  Thing as T ORDER BY T.thingid LIMIT $1 OFFSET $2;")
}

func (pg PostgreSQLGetThings) GetEmptyString() string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, " +
		" FindStudent(T.thingId),  FindRoom(T.thingid) FROM  Thing as T ORDER BY T.thingid;")
}

func (pg PostgreSQLGetThing) GetString() string {
	return "SELECT T.thingid, T.marknumber, T.thingtype, " +
		" FindStudent(T.thingId),  FindRoom(T.thingid) FROM  Thing as T " +
		"WHERE T.thingid = $1"
}

func (pg PostgreSQLGetThingID) GetString() string {
	return "SELECT T.thingid " +
		"FROM  Thing as T WHERE T.marknumber = $1;"
}

func (pg PostgreSQLDeleteThing) GetString() string {
	return "DELETE FROM  Thing WHERE ThingID = $1;"
}

func (pg PostgreSQLGetUserId) GetString() string {
	return "SELECT id FROM  Users WHERE UserLogin = $1;"
}

func (pg PostgreSQLGetUser) GetString() string {
	return "SELECT id, userlogin, userpassword, userrole FROM  Users WHERE ID = $1;"
}

func (pg PostgreSQLAddUser) GetString() string {
	return "INSERT INTO  Users(userlogin, userpassword, userrole) VALUES " +
		"($1, $2, $3);"
}
