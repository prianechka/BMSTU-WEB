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
	return "UPDATE PPO.Student SET StudentName = $1, StudentSurname = $2, " +
		"StudentGroup = $3, StudentNumber = $4 WHERE StudentID = $5;"
}

func (pg PostgreSQLGetStudent) GetString() string {
	return "SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, " +
		"S.studentgroup, S.studentnumber, PPO.FindStudentRoom(S.studentid) " +
		"FROM PPO.Student as S WHERE S.studentid = $1;"
}
func (pg PostgreSQLGetStudentID) GetString() string {
	return "SELECT S.studentid FROM PPO.Student as S WHERE StudentNumber = $1;"
}

func (pg PostgreSQLGetStudentsThings) GetString() string {
	return "SELECT T.thingid, T.marknumber, T.thingtype, PPO.FindStudent(T.thingId), " +
		"PPO.FindRoom(T.thingid) FROM PPO.Thing as T WHERE PPO.FindStudent(T.ThingID) = $1;"
}
func (pg PostgreSQLGetAllStudents) GetString() string {
	return "SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, " +
		"S.studentgroup, S.studentnumber, " +
		"PPO.FindStudentRoom(S.studentid) " +
		"FROM PPO.Student as S;"
}

func (pg PostgreSQLAddStudent) GetString() string {
	return "INSERT INTO PPO.Student(studentname, studentsurname, studentgroup, " +
		"studentnumber, settledate, webaccid) VALUES ($1, $2, $3, $4, current_date, $5);"
}
func (pg PostgreSQLTransferStudent) GetString() string {
	return "INSERT INTO PPO.StudentRoomHistory (studentid, roomid, direction, transferdate) " +
		"VALUES ($1, $2, $3, current_date);"
}
func (pg PostgreSQLTransferThing) GetString() string {
	return "INSERT INTO PPO.StudentThingHistory (studentid, thingid, direction, transferdate)" +
		"VALUES ($1, $2, $3, current_date);"
}

func (pg PostgreSQLAddRoom) GetString() string {
	return "INSERT INTO PPO.Rooms(roomtype, roomnumber) VALUES ('$1', $2)"
}

func (pg PostgreSQLGetRooms) GetString() string {
	return "SELECT roomid, roomtype, roomnumber FROM PPO.rooms;"
}

func (pg PostgreSQLGetRoom) GetString() string {
	return "SELECT roomid, roomtype, roomnumber FROM PPO.rooms WHERE RoomID = $1;"
}

func (pg PostgreSQLGetRoomThings) GetString() string {
	return "SELECT T.thingid, T.marknumber, T.thingtype, PPO.FindStudent(T.thingId), " +
		"PPO.FindRoom(T.thingid) FROM PPO.Thing as T WHERE PPO.FindRoom(T.ThingID) = $1"
}
func (pg PostgreSQLDeleteRoom) GetString() string {
	return "DELETE FROM PPO.rooms WHERE RoomID = $1;"
}

func (pg PostgreSQLTransferThingRoom) GetString() string {
	return "INSERT INTO PPO.ThingRoomHistory (srcroomid, dstroomid, thingid, transferdate) VALUES " +
		"($1, $2, $3, current_date);"
}

func (pg PostgreSQLAddThing) GetString() string {
	return "INSERT INTO PPO.Thing(marknumber, creationdate, thingtype) VALUES  " +
		"($1, current_date, $2);"
}

func (pg PostgreSQLGetThings) GetString() string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, " +
		"PPO.FindStudent(T.thingId), PPO.FindRoom(T.thingid) FROM PPO.Thing as T;")
}

func (pg PostgreSQLGetThing) GetString() string {
	return "SELECT T.thingid, T.marknumber, T.thingtype, " +
		"PPO.FindStudent(T.thingId), PPO.FindRoom(T.thingid) FROM PPO.Thing as T " +
		"WHERE T.thingid = $1"
}

func (pg PostgreSQLGetThingID) GetString() string {
	return "SELECT T.thingid " +
		"FROM PPO.Thing as T WHERE T.marknumber = $1;"
}

func (pg PostgreSQLDeleteThing) GetString() string {
	return "DELETE FROM PPO.Thing WHERE ThingID = $1;"
}

func (pg PostgreSQLGetUserId) GetString() string {
	return "SELECT id FROM PPO.Users WHERE UserLogin = $1;"
}

func (pg PostgreSQLGetUser) GetString() string {
	return "SELECT id, userlogin, userpassword, userrole FROM PPO.Users WHERE ID = $1;"
}

func (pg PostgreSQLAddUser) GetString() string {
	return "INSERT INTO PPO.Users(userlogin, userpassword, userrole) VALUES " +
		"($1, $2, $3);"
}
