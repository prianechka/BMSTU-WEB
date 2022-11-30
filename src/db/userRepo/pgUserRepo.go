package userRepo

import (
	"database/sql"
	"src/db/sql"

	"src/objects"
)

type PgUserRepo struct {
	Conn *sql.DB
}

func (pg *PgUserRepo) GetUserID(login string) (int, error) {
	var result int
	sqlString := pgsql.PostgreSQLGetUserId{}.GetString()
	row := pg.Conn.QueryRow(sqlString, login)
	err := row.Scan(&result)
	return result, err
}

func (pg *PgUserRepo) GetUser(id int) (objects.User, error) {
	var (
		result   = objects.NewEmptyUser()
		userID   int32
		login    string
		password string
		level    objects.Levels
		err      error
	)
	sqlString := pgsql.PostgreSQLGetUser{}.GetString()
	rows, execError := pg.Conn.Query(sqlString, id)
	if execError == nil {
		for rows.Next() {
			readRowError := rows.Scan(&userID, &login, &password, &level)
			if readRowError != nil {
				err = readRowError
				break
			} else {
				result = objects.NewUserWithParams(int(userID), login, password, level)
			}
		}
	} else {
		err = execError
	}
	return result, err
}

func (pg *PgUserRepo) AddUser(login, password string, privelegeLevel objects.Levels) error {
	sqlString := pgsql.PostgreSQLAddUser{}.GetString()
	_, err := pg.Conn.Exec(sqlString, login, password, int(privelegeLevel))
	return err
}
