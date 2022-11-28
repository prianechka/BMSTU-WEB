package studentRepo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"src/db/sql"
	"src/objects"
)

type PgStudentRepo struct {
	ConnectParams objects.PgConnection
}

func (pg *PgStudentRepo) AddStudent(newStudent objects.StudentDTO, accID int) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLAddStudent{}.GetString(newStudent, accID)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgStudentRepo) GetAllStudents() ([]objects.Student, error) {
	var resultStudents = make([]objects.Student, objects.Empty)
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer func() {
		conn.Close(context.Background())
	}()
	if err == nil {
		sqlString := pgsql.PostgreSQLGetAllStudents{}.GetString()
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					studentID := int(values[0].(int32))
					accID := int(values[1].(int32))
					studentName := values[2].(string)
					studentSurname := values[3].(string)
					studentGroup := values[4].(string)
					studentNumber := values[5].(string)
					roomID := int(values[6].(int32))

					tmpThings := objects.NewStudentWithParams(studentID, accID, studentName, studentSurname,
						studentGroup, studentNumber, roomID)
					resultStudents = append(resultStudents, tmpThings)
				} else {
					err = readRowError
				}
			}
		} else {
			err = execError
		}
	}
	return resultStudents, err
}

func (pg *PgStudentRepo) GetStudentID(studentNumber string) (int, error) {
	var result int
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetStudentID{}.GetString(studentNumber)
		row := conn.QueryRow(context.Background(), sqlString)
		err = row.Scan(&result)
	}
	return result, err
}

func (pg *PgStudentRepo) GetStudent(id int) (objects.Student, error) {
	var student = objects.NewEmptyStudent()
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetStudent{}.GetString(id)
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					studentID := int(values[0].(int32))
					accID := int(values[1].(int32))
					studentName := values[2].(string)
					studentSurname := values[3].(string)
					studentGroup := values[4].(string)
					studentNumber := values[5].(string)
					roomID := int(values[6].(int32))

					student = objects.NewStudentWithParams(studentID, accID, studentName, studentSurname,
						studentGroup, studentNumber, roomID)
				} else {
					err = readRowError
				}
			}
		} else {
			err = execError
		}
	}
	return student, err
}

func (pg *PgStudentRepo) TransferStudent(studentID, roomID int, direct objects.TransferDirection) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLTransferStudent{}.GetString(studentID, roomID, direct)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgStudentRepo) ChangeStudent(studentID int, studentInfo objects.StudentDTO) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLChangeStudent{}.GetString(studentID, studentInfo)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgStudentRepo) TransferThing(studentID, thingID int, direct objects.TransferDirection) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLTransferThing{}.GetString(studentID, thingID, direct)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}

func (pg *PgStudentRepo) GetStudentThings(id int) ([]objects.Thing, error) {
	var resultThings = make([]objects.Thing, objects.Empty)
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetStudentsThings{}.GetString(id)
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					thingID := int(values[0].(int32))
					markNumber := int(values[1].(int32))
					thingType := values[2].(string)
					ownerID := int(values[3].(int32))
					roomID := int(values[4].(int32))
					tmpThings := objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
					resultThings = append(resultThings, tmpThings)
				} else {
					err = readRowError
				}
			}
		} else {
			err = execError
		}
	}
	return resultThings, err
}
