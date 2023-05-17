package scylla

import (
	"context"
	"testing"
	"time"

	st "github.com/Sannrox/tradepipe/scylla/testing"
	"github.com/gocql/gocql"
)

const ContactPoint = "127.0.0.1"

func TestNewScyllaKeySpaceConnection(t *testing.T) {

	context := context.Background()

	if err := st.SetUpScylla(context); err != nil {
		t.Errorf("Error: %s", err)
	}

	var keySpace = "test"

	if err := TryToConnectWithRetry(ContactPoint, 10, 10*time.Second); err != nil {
		t.Errorf("Error: %s", err)
	}

	s, err := NewScyllaKeySpaceConnection("127.0.0.1", keySpace)
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

	session, err := cluster.CreateSession()
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	session.Close()
	defer st.TearDownScylla(context)

}

func TestTryToConnectWithRetry(t *testing.T) {

	t.Run("Test TryToConnectWithRetry Successful", func(t *testing.T) {

		context := context.Background()

		if err := st.SetUpScylla(context); err != nil {
			t.Errorf("Error: %s", err)
		}

		if err := TryToConnectWithRetry(ContactPoint, 10, 10*time.Second); err != nil {
			t.Errorf("Error: %s", err)
		}

		defer st.TearDownScylla(context)
	})

	t.Run("Test TryToConnectWithRetry Unsuccessful", func(t *testing.T) {

		if err := TryToConnectWithRetry(ContactPoint, 1, 10*time.Second); err == nil {
			t.Errorf("Error: %s", err)
		}
	})
}
