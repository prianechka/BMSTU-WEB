package mother

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"src/objects"
)

type UserRepoObjectMother struct{}

var (
	DefaultLogin                   = "Ivan"
	DefaultPassword                = "12345678"
	DefaultRole     objects.Levels = objects.StudentRole
)

func (m UserRepoObjectMother) CreateRepo() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func (m UserRepoObjectMother) CreateUser() objects.User {
	return objects.NewUserWithParams(int(InsertID), DefaultLogin, DefaultPassword, DefaultRole)
}

func (m UserRepoObjectMother) CreateDefaultUsers(amount int) []objects.User {
	resultUsers := make([]objects.User, objects.Empty)
	for i := 1; i <= amount; i++ {
		resultUsers = append(resultUsers, objects.NewUserWithParams(i, DefaultLogin,
			DefaultPassword, DefaultRole))
	}
	return resultUsers
}

func (m UserRepoObjectMother) CreateRows(users []objects.User) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "userlogin", "userpassword", "userrole"})
	for _, user := range users {
		rows.AddRow(user.GetID(), user.GetLogin(), user.GetPassword(), user.GetPrivelegeLevel())
	}
	return rows
}

func (m UserRepoObjectMother) CreateRowForID(id int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(id)
	return rows
}
