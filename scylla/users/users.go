package users

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

const (
	tableName = "users"
)

type User struct {
	scylla.Scylla
	Users    *EntryMap
	keyspace string
}

type EntryMap struct {
	entries map[string]*entry
}

type entry struct {
	ID     gocql.UUID `json:"id"`
	Number string     `json:"number"`
	Pin    string     `json:"pin,omitempty"`
}

func NewUserKeyspace(contactPoint string, keyspace string) *User {
	if err := scylla.CreateKeyspace(contactPoint, keyspace); err != nil {
		panic(err)
	}

	s, err := scylla.NewScyllaDbWithPool(contactPoint, keyspace, 10)
	if err != nil {
		panic(err)
	}

	return &User{
		Scylla:   *s,
		Users:    nil,
		keyspace: keyspace,
	}
}

func (u *User) CreateNewUserTable() error {
	if u.TableExists(u.keyspace, "users") {
		u.GetAllUsers()
		return nil
	}
	return u.CreateTable(tableName, "id uuid, number text, pin text, PRIMARY KEY (id)")
}

func (u *EntryMap) ID(number string) *string {
	entry, ok := u.entries[number]
	if !ok {
		return nil
	}
	id := entry.ID.String()
	return &id
}

func (u *EntryMap) Pin(number string) *string {
	entry, ok := u.entries[number]
	if !ok {
		return nil
	}
	pin := entry.Pin
	return &pin
}

func (u *EntryMap) String() string {
	lines := make([]string, 0, len(u.entries))
	for k, e := range u.entries {
		lines = append(lines, fmt.Sprintf("%s{%s,%s};", k, e.ID, e.Number))
	}
	sort.Strings(lines)
	return strings.Join(lines, "")
}

func (u *User) CheckIfUserExists(number string) bool {
	logrus.Debug("Checking if user exists: ", number)
	if u.Users == nil {
		return false
	}
	_, ok := u.Users.entries[number]
	return ok
}

func (u *User) AddUser(number, pin string) error {
	logrus.Debug("Adding user: ", number, pin)
	return u.Scylla.Insert(u.keyspace+"."+tableName, &entry{
		ID:     gocql.TimeUUID(),
		Number: number,
		Pin:    pin,
	})
}

func (u *User) ReadUser(number string) (*entry, error) {
	logrus.Debug("Reading user: ", number)
	err := u.GetAllUsers()
	if err != nil {
		return nil, err
	}

	e := u.Users.entries[number]
	return e, nil
}

func (u *User) UpdateUser(number, pin string) error {
	logrus.Debug("Updating user: ", number, pin)
	entry, err := u.ReadUser(number)
	if err != nil {
		return err
	}

	if pin != "" {
		entry.Pin = pin
	}

	if number != "" {
		entry.Number = number
	}
	return u.Scylla.Update(u.keyspace+"."+tableName, entry.ID, entry)
}

func (u *User) DeleteUser(number string) error {
	logrus.Debug("Deleting user: ", number)
	entry, err := u.ReadUser(number)
	if err != nil {
		return err
	}
	return u.Scylla.Delete(tableName, entry.ID, entry)
}

func (u *User) GetAllUsers() error {
	logrus.Debug("Getting all users")
	users := &EntryMap{entries: make(map[string]*entry)}
	iter := u.Scylla.GetAllValues(u.keyspace + "." + tableName)

	defer iter.Close()

	for {
		entry := &entry{}
		var idStr string
		if !iter.Scan(&idStr, &entry.Number, &entry.Pin) {
			break
		}
		entry.ID, _ = gocql.ParseUUID(idStr)
		users.entries[entry.Number] = entry
	}

	if err := iter.Close(); err != nil {
		return err
	}

	u.Users = users

	return nil
}
