package system

import (
	"fmt"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/scylla/keyspaces/system/users"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sirupsen/logrus"
)

const usersTableName = "users"

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

func (s *System) CreateUserTable() error {
	logrus.Debug("creating user table")
	if err := s.CreateTable(usersTableName, "id uuid, number text, pin text, PRIMARY KEY (id, number)"); err != nil {
		return err
	}
	logrus.Debug("created user table")
	tableMeta := table.Metadata{
		Name: usersTableName, Columns: []string{"id", "number", "pin"},
		PartKey: []string{"id"},
		SortKey: []string{"number"},
	}
	s.Users.Table = s.Scylla.NewTable(tableMeta)
	logrus.Debug("created user table")
	return nil
}

func (s *System) CreateTables() error {
	if err := s.CreateUserTable(); err != nil {
		return err
	}
	return nil
}

func (s *System) CreateUser(number, pin string) error {
	if s.Users.Users == nil {
		s.Users.Users = users.NewUsers()
	}
	if !s.CheckIfUserExists(number) {
		user, err := s.Users.CreateNewUser(number, pin)
		if err != nil {
			return err
		}
		if err := s.Users.AddUser(user); err != nil {
			return err
		}

		return s.Insert(s.Users.Table, &user)
	} else {
		return fmt.Errorf("User with number %v already exists", number)
	}
}

func (s *System) GetUsers() error {
	var allUsers []users.User
	s.Users.Users = users.NewUsers()
	exists := s.Scylla.CheckIfTableExits(usersTableName)
	if exists {
		if len(*s.Users.Users) == 0 {
			q := s.Scylla.GetAll(s.Users.Table)
			if err := q.SelectRelease(&allUsers); err != nil {
				return err
			}
			for _, user := range allUsers {
				if !s.Users.CheckIfUserExists(user.Number) {
					if err := s.Users.AddUser(&user); err != nil {
						return err
					}
				}
			}
		}

	} else {
		return fmt.Errorf("user table does not exist")
	}

	return nil
}

func (s *System) GetUser(number string) (*users.User, error) {
	if s.Users.Users == nil {
		if err := s.GetUsers(); err != nil {
			return nil, err
		}
	}
	if s.CheckIfUserExists(number) {
		return s.Users.ReadUser(number), nil
	} else {
		return nil, fmt.Errorf("User with number %v does not exist", number)
	}
}

func (s *System) UpdateUser(number, pin string) error {
	if s.Users.Users == nil {
		if err := s.GetUsers(); err != nil {
			return err
		}
	}

	if err := s.Users.UpdateUser(number, pin); err != nil {
		return err
	}

	user := s.Users.ReadUser(number)
	user.Pin = pin

	return s.Update(s.Users.Table, user, "pin")
}

func (s *System) DeleteUser(number string) error {

	user := s.Users.ReadUser(number)

	if err := s.Users.DeleteUser(number); err != nil {
		return err
	}

	return s.Delete(s.Users.Table, user)
}
