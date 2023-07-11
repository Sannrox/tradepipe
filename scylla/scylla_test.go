package scylla

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Sannrox/tradepipe/helper/testhelpers/utils"
	"github.com/Sannrox/tradepipe/scylla/testing/container"
	"github.com/gocql/gocql"
)

const ContactPoint = "127.0.0.1"

// StartPort and EndPort are the ports used to find a free port for the scylla container
// Different Ports are needed for different tests
const startPort = 9030
const endPort = 9040

func TestNewScyllaKeySpaceConnection(t *testing.T) {

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

	var keySpace = "test"

	if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	s, err := NewScyllaKeySpaceConnection("127.0.0.1", port, keySpace)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if s == nil {
		t.Errorf("Error: %s", err)
	}

	s.Close()

	// Check if the space exits

	cluster := gocql.NewCluster(ContactPoint)
	cluster.Keyspace = keySpace
	cluster.Timeout = 10 * time.Second
	cluster.Port = port

	session, err := cluster.CreateSession()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	session.Close()

}

func TestTryToConnectWithRetry(t *testing.T) {

	t.Run("Test TryToConnectWithRetry Successful", func(t *testing.T) {

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

		if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
			t.Errorf("Error: %s", err)
		}

	})

	t.Run("Test TryToConnectWithRetry Unsuccessful", func(t *testing.T) {

		freePort, err := utils.FindFreePort(9024, 9040)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if err := TryToConnectWithRetry(ContactPoint, freePort, 1, 10*time.Second); err == nil {
			t.Errorf("Error: %s", err)
		}
	})
}

func TestNewScyllaConnection(t *testing.T) {

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

	if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := CreateKeyspace(ContactPoint, port, "test"); err != nil {
		t.Errorf("Error: %s", err)
	}

	s, err := NewScyllaDbWithPool(ContactPoint, port, "test", 10)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if s == nil {
		t.Errorf("Error: %s", err)
	}

	s.Close()

}

func TestCreateKeySpace(t *testing.T) {

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

	if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := CreateKeyspace(ContactPoint, port, "test"); err != nil {
		t.Errorf("Error: %s", err)
	}

}

func TestCreateTable(t *testing.T) {

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
	if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := CreateKeyspace(ContactPoint, port, "test"); err != nil {
		t.Errorf("Error: %s", err)
	}

	s, err := NewScyllaDbWithPool(ContactPoint, port, "test", 10)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	testSchema := `
	id text,
	name text,
	PRIMARY KEY (id)
	`

	if err := s.CreateTable("test_table", testSchema); err != nil {
		t.Errorf("Error: %s", err)
	}

}

func TestClose(t *testing.T) {

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
	if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := CreateKeyspace(ContactPoint, port, "test"); err != nil {
		t.Errorf("Error: %s", err)
	}

	s, err := NewScyllaDbWithPool(ContactPoint, port, "test", 10)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	s.Close()

	if s.Session.Closed() == false {
		t.Errorf("Error: %s", err)
	}

}

func TestCheckIfTableExits(t *testing.T) {
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
	if err := TryToConnectWithRetry(ContactPoint, port, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	if err := CreateKeyspace(ContactPoint, port, "test"); err != nil {
		t.Errorf("Error: %s", err)
	}
	s, err := NewScyllaDbWithPool(ContactPoint, port, "test", 10)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	t.Run("Table does not exist", func(t *testing.T) {
		if s.CheckIfTableExits("test_table") {
			t.Errorf("Error: table should not exist")
		}
	})

	t.Run("Table does exist", func(t *testing.T) {
		if err := s.CreateTable("test_table", "id text, name text, PRIMARY KEY (id)"); err != nil {
			t.Errorf("Error: %s", err)
		}

		if !s.CheckIfTableExits("test_table") {
			t.Errorf("Error: table should exist")
		}

	})
}
