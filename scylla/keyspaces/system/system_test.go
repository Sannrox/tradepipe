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

	if err := system.CreateNewUser("test", "test"); err != nil {
		t.Errorf("Error: %s", err)
	}
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
		if err := system.CreateNewUser(v.username, v.password); err != nil {
			t.Errorf("Error: %s", err)
		}
	}

	if err := system.GetUsers(); err != nil {
		t.Errorf("Error: %s", err)
	}

	users := *system.Users.Users

	if len(users) != len(test) {
		t.Errorf("Error: Expected %d, got %d", len(test), len(users))
	}

	for i, v := range users {
		if v.Number != test[i].username {
			t.Errorf("Error: Expected %s, got %s", test[i].username, v.Number)
		}
		if v.Pin != test[i].password {
			t.Errorf("Error: Expected %s, got %s", test[i].password, v.Pin)
		}
	}
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
		if err := system.CreateNewUser(v.username, v.password); err != nil {
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
		if err := system.CreateNewUser(v.username, v.password); err != nil {
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
		if err := system.CreateNewUser(v.username, v.password); err != nil {
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
