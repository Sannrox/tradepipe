package system

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/scylla/testing/container"
)

const (
	startPort    = 9050
	endPort      = 9060
	ContactPoint = "127.0.0.1"
)

func TestNewSystemKeySpace(t *testing.T) {

	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	_, err = NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

}

func TestCreateUserTable(t *testing.T) {
	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	system, err := NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	err = system.CreateUserTable()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	system.Close()
}

func TestNewUser(t *testing.T) {
	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	system, err := NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := system.CreateUserTable(); err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := system.CreateUser("test", "0000"); err != nil {
		t.Errorf("Error: %s", err)
	}

	system.Close()
}

func TestGetUsers(t *testing.T) {
	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	system, err := NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := system.CreateUserTable(); err != nil {
		t.Errorf("Error: %s", err)
	}

	test := map[string]struct{ username, password string }{
		"test1": {"test1", "0000"},
		"test2": {"test2", "0000"},
		"test3": {"test3", "0000"},
	}

	for _, v := range test {
		if err := system.CreateUser(v.username, v.password); err != nil {
			t.Errorf("Error: %s", err)
		}
	}

	*system.Users.Users = nil

	if err := system.GetUsers(); err != nil {
		t.Errorf("Error: %s", err)
	}

	users := system.GetAllUsers()

	if len(users) != len(test) {
		t.Errorf("Error: Expected %d, got %d", len(test), len(*system.Users.Users))
	}

	equalUsers := 0
	for i, v := range users {
		t.Log(i, v)
		if system.ReadUser(i).Number == test[i].username && system.ReadUser(i).Pin == test[i].password {
			equalUsers++
		}
	}

	if equalUsers != len(test) {
		t.Errorf("Error: Expected %d users to be equal, got %d", len(test), equalUsers)
	}

	system.Close()
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	system, err := NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := system.CreateUserTable(); err != nil {
		t.Errorf("Error: %s", err)
	}

	test := map[string]struct{ username, password string }{
		"test1": {"test1", "0000"},
		"test2": {"test2", "0000"},
		"test3": {"test3", "0000"},
	}

	for _, v := range test {
		if err := system.CreateUser(v.username, v.password); err != nil {
			t.Errorf("Error: %s", err)
		}
	}

	user, err := system.GetUser("test1")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if user.Number != test["test1"].username {
		t.Errorf("Error: Expected %s, got %s", test["test1"].username, user.Number)
	}
	if user.Pin != test["test1"].password {
		t.Errorf("Error: Expected %s, got %s", test["test1"].password, user.Pin)
	}

	system.Close()
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	system, err := NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := system.CreateUserTable(); err != nil {
		t.Errorf("Error: %s", err)
	}

	test := map[string]struct{ username, password string }{
		"test1": {"test1", "0000"},
		"test2": {"test2", "0000"},
		"test3": {"test3", "0000"},
	}

	for _, v := range test {
		if err := system.CreateUser(v.username, v.password); err != nil {
			t.Errorf("Error: %s", err)
		}
	}

	if err := system.UpdateUser("test1", "0001"); err != nil {
		t.Errorf("Error: %s", err)
	}

	user, err := system.GetUser("test1")
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if user.Number != test["test1"].username {
		t.Errorf("Error: Expected %s, got %s", test["test1"].username, user.Number)
	}
	if user.Pin != "0001" {
		t.Errorf("Error: Expected %s, got %s", "0001", user.Pin)
	}

	system.Close()

}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	containerName, port, err := container.SetUpScylla(ctx, startPort, endPort)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Cleanup(func() {
		if err := container.TearDownScylla(containerName, ctx); err != nil {
			t.Fatal(fmt.Errorf("failed to tear down scylla container: %w", err))
		}
	})
	if err := scylla.TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	system, err := NewSystemKeyspace(ContactPoint, port)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := system.CreateUserTable(); err != nil {
		t.Errorf("Error: %s", err)
	}

	test := map[string]struct{ username, password string }{
		"test1": {"test1", "0000"},
		"test2": {"test2", "0000"},
		"test3": {"test3", "0000"},
	}

	for _, v := range test {
		if err := system.CreateUser(v.username, v.password); err != nil {
			t.Errorf("Error: %s", err)
		}
	}

	if err := system.DeleteUser("test1"); err != nil {
		t.Errorf("Error: %s", err)
	}

	_, err = system.GetUser("test1")
	if err == nil {
		t.Errorf("Error found User")

	}

}
