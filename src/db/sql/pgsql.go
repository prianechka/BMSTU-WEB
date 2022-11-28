package pgsql

import (
	"fmt"
	"src/objects"
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

func (pg PostgreSQLChangeStudent) GetString(studentID int, studentInfo objects.StudentDTO) string {
	return fmt.Sprintf("UPDATE PPO.Student SET StudentName = '%s', StudentSurname = '%s', "+
		"', StudentGroup = '%s', StudentNumber = '%s' WHERE StudentID = %d;", studentInfo.GetName(),
		studentInfo.GetSurname(), studentInfo.GetStudentGroup(), studentInfo.GetStudentNumber(),
		studentID)
}

func (pg PostgreSQLGetStudent) GetString(id int) string {
	return fmt.Sprintf("SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, "+
		"S.studentgroup, S.studentnumber, PPO.FindStudentRoom(S.studentid) "+
		"FROM PPO.Student as S WHERE S.studentid = %d;", id)
}
func (pg PostgreSQLGetStudentID) GetString(studentNumber string) string {
	return fmt.Sprintf("SELECT S.studentid FROM PPO.Student as S WHERE StudentNumber = '%s';",
		studentNumber)
}

func (pg PostgreSQLGetStudentsThings) GetString(id int) string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, PPO.FindStudent(T.thingId), "+
		"PPO.FindRoom(T.thingid) FROM PPO.Thing as T WHERE PPO.FindStudent(T.ThingID) = %d;", id)
}
func (pg PostgreSQLGetAllStudents) GetString() string {
	return fmt.Sprintf("SELECT S.studentid, S.webaccid, S.studentname, S.studentsurname, " +
		"S.studentgroup, S.studentnumber, " +
		"PPO.FindStudentRoom(S.studentid) " +
		"FROM PPO.Student as S;")
}

func (pg PostgreSQLAddStudent) GetString(newStudent objects.StudentDTO, accID int) string {
	return fmt.Sprintf("INSERT INTO PPO.Student(studentname, studentsurname, studentgroup, "+
		"studentnumber, settledate, webaccid) VALUES ('%s', '%s', '%s', '%s', current_date, %d);",
		newStudent.GetName(), newStudent.GetSurname(), newStudent.GetStudentGroup(),
		newStudent.GetStudentNumber(), accID)
}
func (pg PostgreSQLTransferStudent) GetString(studentID, roomID int, direct objects.TransferDirection) string {
	return fmt.Sprintf("INSERT INTO PPO.StudentRoomHistory (studentid, roomid, direction, transferdate) "+
		"VALUES (%d, %d, %d, current_date);", studentID, roomID, direct)
}
func (pg PostgreSQLTransferThing) GetString(studentID, thingID int, direct objects.TransferDirection) string {
	return fmt.Sprintf("INSERT INTO PPO.StudentThingHistory (studentid, thingid, direction, transferdate)"+
		"VALUES (%d, %d, %d, current_date);", studentID, thingID, direct)
}

func (pg PostgreSQLAddRoom) GetString(room objects.RoomDTO) string {
	return fmt.Sprintf("INSERT INTO PPO.Rooms(roomtype, roomnumber) VALUES ('%s', %d)",
		room.GetRoomType(), room.GetRoomNumber())
}

func (pg PostgreSQLGetRooms) GetString() string {
	return fmt.Sprintf("SELECT * FROM PPO.rooms;")
}

func (pg PostgreSQLGetRoom) GetString(id int) string {
	return fmt.Sprintf("SELECT * FROM PPO.rooms WHERE RoomID = %d;", id)
}

func (pg PostgreSQLGetRoomThings) GetString(id int) string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, PPO.FindStudent(T.thingId), "+
		"PPO.FindRoom(T.thingid) FROM PPO.Thing as T WHERE PPO.FindRoom(T.ThingID) = %d", id)
}
func (pg PostgreSQLDeleteRoom) GetString(id int) string {
	return fmt.Sprintf("DELETE FROM PPO.rooms WHERE RoomID = %d;", id)
}

func (pg PostgreSQLTransferThingRoom) GetString(id, srcRoomID, dstRoomID int) string {
	return fmt.Sprintf("INSERT INTO PPO.ThingRoomHistory (srcroomid, dstroomid, thingid, transferdate) VALUES "+
		"(%d, %d, %d, current_date);", srcRoomID, dstRoomID, id)
}

func (pg PostgreSQLAddThing) GetString(thing objects.ThingDTO) string {
	return fmt.Sprintf("INSERT INTO PPO.Thing(marknumber, creationdate, thingtype) VALUES  "+
		"(%d, current_date, %s", thing.GetMarkNumber(), thing.GetThingType())
}

func (pg PostgreSQLGetThings) GetString() string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, " +
		"PPO.FindStudent(T.thingId), PPO.FindRoom(T.thingid) FROM PPO.Thing as T;")
}

func (pg PostgreSQLGetThing) GetString(id int) string {
	return fmt.Sprintf("SELECT T.thingid, T.marknumber, T.thingtype, "+
		"PPO.FindStudent(T.thingId), PPO.FindRoom(T.thingid) FROM PPO.Thing as T "+
		"WHERE T.thingid = %d", id)
}

func (pg PostgreSQLGetThingID) GetString(marknumber int) string {
	return fmt.Sprintf("SELECT T.thingid "+
		"FROM PPO.Thing as T WHERE T.marknumber = %d;", marknumber)
}

func (pg PostgreSQLDeleteThing) GetString(id int) string {
	return fmt.Sprintf("DELETE FROM PPO.Thing WHERE ThingID = %d;", id)
}

func (pg PostgreSQLGetUserId) GetString(login string) string {
	return fmt.Sprintf("SELECT * FROM PPO.Users WHERE UserLogin = '%s';", login)
}

func (pg PostgreSQLGetUser) GetString(id int) string {
	return fmt.Sprintf("SELECT * FROM PPO.Users WHERE ID = %d;", id)
}

func (pg PostgreSQLAddUser) GetString(login, password string, privelegeLevel objects.Levels) string {
	return fmt.Sprintf("INSERT INTO PPO.Users(userlogin, userpassword, userrole) VALUES"+
		"('%s', '%s', %d);", login, password, privelegeLevel)
}
