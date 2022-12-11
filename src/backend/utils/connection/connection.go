package connection

import (
	"database/sql"
	"fmt"
	"src/configs"
)

func NewPgSQLConnection(conn configs.PgSQLConnectionParams) *sql.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conn.User, conn.Password,
		conn.Host, conn.Port, conn.Database)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	return db
}
