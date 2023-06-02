package users

import (
	"fmt"

	"github.com/gocql/gocql"
)

type Users map[string]User

type User struct {
	Id     gocql.UUID
	Number string
	Pin    string
}

func NewUsers() *Users {
	return &Users{}
}

func (u *User) SetPin(pin string) {
	u.Pin = pin
}

func (u *User) SetNumber(number string) {
	u.Number = number
}

func (u *Users) GetUserCount() int {
	if u == nil {
		return 0
	}

	users := *u
	return len(users)
}

func (u *Users) CheckIfUserExists(number string) bool {
	if u == nil {
		return false
	}

	users := *u
	_, ok := users[number]

	return ok
}

func (u *Users) AddUser(user *User) error {
	if u.CheckIfUserExists(user.Number) {
		return fmt.Errorf("User with number %v already exists", user.Number)
	}

	(*u)[user.Number] = *user

	return nil
}

func (u *Users) CreateNewUser(number, pin string) (*User, error) {
	switch {
	case pin == "":
		return nil, fmt.Errorf("Pin cannot be empty")
	case len(pin) != 4:
		return nil, fmt.Errorf("Pin must be 4 digits")

	case number == "":
		return nil, fmt.Errorf("Number cannot be empty")
	}
	var user = &User{
		Id:     gocql.TimeUUID(),
		Number: number,
		Pin:    pin,
	}

	return user, nil
}

func (u *Users) GetAllUsers() Users {
	if u == nil {
		return nil
	}

	users := *u
	return users
}

func (u *Users) ReadUser(number string) *User {
	if u == nil {
		return nil
	}

	users := *u
	user, ok := users[number]
	if !ok {
		return nil
	}

	return &user

}

func (u *Users) UpdateUser(number, pin string) error {

	if !u.CheckIfUserExists(number) {
		return fmt.Errorf("User with number %v does not exist", number)
	}

	user := u.ReadUser(number)

	switch {
	case user == nil:
		return fmt.Errorf("User with number %v does not exist", number)
	case user.Pin == "":
		return fmt.Errorf("User with number %v does not have a pin", number)
	case number == "":
		return fmt.Errorf("Number cannot be empty")
	case pin == "":
		return fmt.Errorf("Pin cannot be empty")
	case len(pin) != 4:
		return fmt.Errorf("Pin must be 4 digits")
	}
	user.SetPin(pin)

	user.SetNumber(number)

	(*u)[number] = *user

	return nil
}

func (u *Users) DeleteUser(number string) error {
	if !u.CheckIfUserExists(number) {
		return fmt.Errorf("User with number %v does not exist", number)
	}

	delete(*u, number)
	return nil
}

func (u *Users) All() (*Users, error) {
	if u == nil {
		return nil, fmt.Errorf("Users is nil")
	}

	users := *u

	return &users, nil

}
