package objects

import "fmt"

type Connection interface {
	GetUser() string
	GetHost() string
	GetDBName() string
	GetPassword() string
	GetPort() int

	GetURL() string
}

type PgConnection struct {
	user     string
	host     string
	db       string
	password string
	port     int
}

func InitPgConnectionByParams(user, host, db, password string, port int) PgConnection {
	return PgConnection{
		user:     user,
		host:     host,
		db:       db,
		password: password,
		port:     port,
	}
}

func (pc *PgConnection) GetUser() string {
	return pc.user
}

func (pc *PgConnection) GetHost() string {
	return pc.host
}

func (pc *PgConnection) GetDBName() string {
	return pc.db
}

func (pc *PgConnection) GetPassword() string {
	return pc.password
}

func (pc *PgConnection) GetPort() int {
	return pc.port
}

func (pc *PgConnection) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", pc.user, pc.password, pc.host, pc.port, pc.db)
}
