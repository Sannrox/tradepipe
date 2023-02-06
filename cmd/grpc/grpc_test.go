package grpc

import (
	"testing"
	"time"

	test "github.com/Sannrox/tradepipe/pkg/grpc/testing"
)

func TestGrpcServer(t *testing.T) {
	s := NewGRPCServer()

	go s.Run()
	time.Sleep(1 * time.Second)

	c := test.NewFakeClient()
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

}
