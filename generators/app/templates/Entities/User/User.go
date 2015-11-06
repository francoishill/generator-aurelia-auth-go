package User

import (
	"encoding/json"
	"fmt"
)

type User interface {
	UUID() interface{}
	Id() int64
	FullName() string
	Email() string

	SetFullName(fullName string)

	Clone() User
}

type user struct {
	id       int64
	fullName string
	email    string
}

func (u *user) UUID() interface{} { return fmt.Sprintf("%d", u.id) }
func (u *user) Id() int64         { return u.id }
func (u *user) FullName() string  { return u.fullName }
func (u *user) Email() string     { return u.email }

func (u *user) SetFullName(fullName string) { u.fullName = fullName }

func (u *user) Clone() User { return NewUser(u.id, u.fullName, u.email) }

func (u *user) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id       int64
		FullName string
		Email    string
	}{u.id, u.fullName, u.email})
}

func NewUser(id int64, fullname, email string) *user {
	return &user{
		id,
		fullname,
		email,
	}
}
