package userRepo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"src/db/sql"

	"src/objects"
)

type PgUserRepo struct {
	ConnectParams objects.PgConnection
}

func (pg *PgUserRepo) GetUserID(login string) (int, error) {
	var result int
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetUserId{}.GetString(login)
		row := conn.QueryRow(context.Background(), sqlString)
		err = row.Scan(&result)
	}
	return result, err
}

func (pg *PgUserRepo) GetUser(id int) (objects.User, error) {
	var result = objects.NewEmptyUser()
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLGetUser{}.GetString(id)
		rows, execError := conn.Query(context.Background(), sqlString)
		if execError == nil {
			for rows.Next() {
				values, readRowError := rows.Values()
				if readRowError == nil {
					userID := int(values[0].(int32))
					login := values[1].(string)
					password := values[2].(string)
					level := objects.Levels(values[3].(int32))
					result = objects.NewUserWithParams(userID, login, password, level)
				}
			}
		}
	}
	return result, err
}

func (pg *PgUserRepo) AddUser(login, password string, privelegeLevel objects.Levels) error {
	conn, err := pgx.Connect(context.Background(), pg.ConnectParams.GetURL())
	defer conn.Close(context.Background())
	if err == nil {
		sqlString := pgsql.PostgreSQLAddUser{}.GetString(login, password, privelegeLevel)
		_, err = conn.Query(context.Background(), sqlString)
	}
	return err
}
