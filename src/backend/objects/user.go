package objects

type Levels int

const (
	NonAuth = iota
	StudentRole
	SupplyRole
	ComendRole
)

type User struct {
	id             int
	login          string
	password       string
	privelegeLevel Levels
}

func NewUserWithParams(id int, login, password string, privelegeLevel Levels) User {
	return User{
		id:             id,
		login:          login,
		password:       password,
		privelegeLevel: privelegeLevel,
	}
}

func NewEmptyUser() User {
	return User{id: None}
}

func (u *User) GetID() int {
	return u.id
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetPrivelegeLevel() Levels {
	return u.privelegeLevel
}
