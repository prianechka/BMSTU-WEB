package studentRepo

import (
	"database/sql"
	"src/db/sql"
	"src/objects"
)

type PgStudentRepo struct {
	Conn *sql.DB
}

func (pg *PgStudentRepo) AddStudent(newStudent objects.StudentDTO, accID int) error {
	sqlString := pgsql.PostgreSQLAddStudent{}.GetString()
	_, err := pg.Conn.Exec(sqlString, newStudent.GetName(), newStudent.GetSurname(), newStudent.GetStudentGroup(),
		newStudent.GetStudentNumber(), accID)
	return err
}

func (pg *PgStudentRepo) GetAllStudents() ([]objects.Student, error) {
	var (
		resultStudents                                           = make([]objects.Student, objects.Empty)
		studentID, accID, roomID                                 int
		studentName, studentSurname, studentGroup, studentNumber string
		err                                                      error
	)
	sqlString := pgsql.PostgreSQLGetAllStudents{}.GetString()
	rows, execError := pg.Conn.Query(sqlString)
	if execError == nil {
		for rows.Next() {
			scanErr := rows.Scan(&studentID, &accID, &studentName, &studentSurname, &studentGroup, &studentNumber, &roomID)
			if scanErr == nil {
				tmpStudent := objects.NewStudentWithParams(studentID, accID, studentName, studentSurname,
					studentGroup, studentNumber, roomID)
				resultStudents = append(resultStudents, tmpStudent)
			} else {
				err = scanErr
			}
		}
	} else {
		err = execError
	}
	return resultStudents, err
}

func (pg *PgStudentRepo) GetStudentID(studentNumber string) (int, error) {
	var result = objects.None
	sqlString := pgsql.PostgreSQLGetStudentID{}.GetString()
	row := pg.Conn.QueryRow(sqlString, studentNumber)
	err := row.Scan(&result)
	return result, err
}

func (pg *PgStudentRepo) GetStudent(id int) (objects.Student, error) {
	var (
		student                                                  = objects.NewEmptyStudent()
		studentID, accID, roomID                                 int
		studentName, studentSurname, studentGroup, studentNumber string
		err                                                      error
	)
	sqlString := pgsql.PostgreSQLGetStudent{}.GetString()
	rows, execError := pg.Conn.Query(sqlString, id)
	if execError == nil {
		for rows.Next() {
			scanErr := rows.Scan(&studentID, &accID, &studentName, &studentSurname, &studentGroup, &studentNumber, &roomID)
			if scanErr == nil {
				student = objects.NewStudentWithParams(studentID, accID, studentName, studentSurname,
					studentGroup, studentNumber, roomID)
			} else {
				err = scanErr
			}
		}
	} else {
		err = execError
	}
	return student, err
}

func (pg *PgStudentRepo) TransferStudent(studentID, roomID int, direct objects.TransferDirection) error {
	sqlString := pgsql.PostgreSQLTransferStudent{}.GetString()
	_, err := pg.Conn.Exec(sqlString, studentID, roomID, int(direct))
	return err
}

func (pg *PgStudentRepo) ChangeStudent(studentID int, studentInfo objects.StudentDTO) error {
	sqlString := pgsql.PostgreSQLChangeStudent{}.GetString()
	_, err := pg.Conn.Exec(sqlString, studentInfo.GetName(), studentInfo.GetSurname(),
		studentInfo.GetStudentGroup(), studentInfo.GetStudentNumber(), studentID)
	return err
}

func (pg *PgStudentRepo) TransferThing(studentID, thingID int, direct objects.TransferDirection) error {
	sqlString := pgsql.PostgreSQLTransferThing{}.GetString()
	_, err := pg.Conn.Exec(sqlString, studentID, thingID, int(direct))
	return err
}

func (pg *PgStudentRepo) GetStudentThings(id int) ([]objects.Thing, error) {
	var (
		resultThings                         = make([]objects.Thing, objects.Empty)
		thingID, markNumber, ownerID, roomID int
		thingType                            string
		err                                  error
	)
	sqlString := pgsql.PostgreSQLGetStudentsThings{}.GetString()
	rows, execError := pg.Conn.Query(sqlString, id)
	if execError == nil {
		for rows.Next() {
			readRowErr := rows.Scan(&thingID, &markNumber, &thingType, &ownerID, &roomID)
			if err == nil {
				tmpThings := objects.NewThingWithParams(thingID, markNumber, thingType, ownerID, roomID)
				resultThings = append(resultThings, tmpThings)
			} else {
				err = readRowErr
				break
			}
		}
	} else {
		err = execError
	}
	return resultThings, err
}
