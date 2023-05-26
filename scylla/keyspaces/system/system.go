package system

import (
	"fmt"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/scylla/keyspaces/system/users"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sirupsen/logrus"
)

type System struct {
	scylla.Scylla
	Users
}

type Users struct {
	*table.Table
	*users.Users
}

func NewSystemKeyspace(contactPoint string, port int) (*System, error) {
	if err := scylla.CreateKeyspace(contactPoint, port, "tradepipe_system"); err != nil {
		return nil, err
	}

	s, err := scylla.NewScyllaDbWithPool(contactPoint, port, "tradepipe_system", 10)
	if err != nil {
		return nil, err
	}

	return &System{
		Scylla: *s,
		Users:  Users{},
	}, nil
}

func (s *System) CreateUserTable(scy *scylla.Scylla) error {
	logrus.Debug("creating user table")
	if err := s.CreateTable("users", "id uuid, number text, pin text, PRIMARY KEY (id)"); err != nil {
		return err
	}
	logrus.Debug("created user table")
	tableMeta := table.Metadata{
		Name: "users", Columns: []string{"id", "number", "pin"},
		PartKey: []string{"id"},
		SortKey: []string{"number"},
	}
	s.Users.Table = scy.NewTable(tableMeta)
	logrus.Debug("created user table")
	return nil
}

func (s *System) CreateTables() error {
	if err := s.CreateUserTable(&s.Scylla); err != nil {
		return err
	}
	return nil
}

func (s *System) CreateNewUser(number, pin string) error {
	if s.Users.Users == nil {
		s.Users.Users = users.NewUsers()
	}
	user := s.Users.CreateNewUser(number, pin)
	if err := s.Users.AddUser(&user); err != nil {
		return err
	}

	return s.Insert(s.Users.Table, &user)
}

func (s *System) getUsers() error {
	var allUsers []users.User
	if s.Users.Users == nil {
		s.Users.Users = users.NewUsers()
		exists, err := s.Scylla.CheckIfTableExits("user")
		if err != nil {
			return err
		}
		if exists {
			q := s.Scylla.GetAll(s.Users.Table)
			if err := q.SelectRelease(&allUsers); err != nil {
				return err
			}

		}
		for _, user := range allUsers {
			s.Users.AddUser(&user)
		}
	}
	q := s.Scylla.GetAll(s.Users.Table)
	if err := q.SelectRelease(&allUsers); err != nil {
		return err
	}

	for _, user := range allUsers {
		s.Users.AddUser(&user)
	}

	return nil
}

func (s *System) GetUser(number string) (*users.User, error) {
	if err := s.getUsers(); err != nil {
		return nil, err
	}
	if s.CheckIfUserExists(number) {
		return s.Users.ReadUser(number), nil
	} else {
		return nil, fmt.Errorf("User with number %v does not exist", number)
	}
}

func (s *System) UpdateUser(number, pin string) error {
	if err := s.getUsers(); err != nil {
		return err
	}

	if err := s.Users.UpdateUser(number, pin); err != nil {
		return err
	}

	user := s.Users.ReadUser(number)

	return s.Update(s.Users.Table, &user)
}

func (s *System) DeleteUser(number string) error {
	if err := s.getUsers(); err != nil {
		return err
	}

	if err := s.Users.DeleteUser(number); err != nil {
		return err
	}
	user := s.Users.ReadUser(number)
	return s.Delete(s.Users.Table, &user)
}
