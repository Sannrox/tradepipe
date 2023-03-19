package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Sannrox/tradepipe/grpc/pb"
	test "github.com/Sannrox/tradepipe/helper/testhelpers/fakegrpcclient"
	fake "github.com/Sannrox/tradepipe/helper/testhelpers/faketrserver"
	"github.com/Sannrox/tradepipe/helper/testhelpers/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var FakeTRServerPort string = "3443"

func TestGrpcServer(t *testing.T) {
	done := make(chan struct{})
	s := NewGRPCServer()
	s.SetBaseURL("https://localhost:" + FakeTRServerPort)
	s.SetWsURL("wss://localhost:" + FakeTRServerPort)

	setClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	s.SetBaseHTTPClient(setClient)
	s.SetOverWriteTls(true)

	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	FakeServer.GenerateData()

	go FakeServer.Run(done, FakeTRServerPort)
	go s.Run(done)

	if err := utils.WaitForRestServerToBeUp("https://localhost:"+FakeTRServerPort, 10); err != nil {
		t.Fatal(err)
	}

	if err := waitForGrpcServerToBeUp("localhost:50051", 10); err != nil {
		t.Fatal(err)
	}

	t.Run("Login test", Login)

	t.Run("Verify test", Verify)

	t.Run("Timeline test", Timeline)

	t.Run("TimelineDetail test", TimelineDetails)

	t.Run("Portfolio test", Portfolio)

	t.Run("Savingsplan test", SavingsPlans)

	close(done)
	if err := utils.WaitForPortToBeNotAttachedWithLimit(FakeTRServerPort, 10); err != nil {
		t.Fatal(err)
	}
	if err := utils.WaitForPortToBeNotAttachedWithLimit("50051", 10); err != nil {
		t.Fatal(err)
	}
}

func Login(t *testing.T) {
	c := test.NewFakeClient()
	c.SetCredentials("+49111111111", "1111")
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	resp, err := c.Login()
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Error) != 0 {
		t.Fatal(resp.Error)
	}

	if len(resp.ProcessId) == 0 {
		t.Fatal("ProcessID is empty")
	}
	if resp.ProcessId != "1234567890" {
		t.Fatal("ProcessID is not correct")
	}

}

func Verify(t *testing.T) {
	c := test.NewFakeClient()
	c.SetCredentials("+49111111111", "1111")
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	resp, err := c.Verify("1234567890", 1234)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Error) != 0 {
		t.Fatal(resp.Error)
	}

}

func Timeline(t *testing.T) {
	c := test.NewFakeClient()
	c.SetCredentials("+49111111111", "1111")
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	resp, err := c.Timeline("1234567890", 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Error) != 0 {
		t.Fatal(resp.Error)
	}

}

func TimelineDetails(t *testing.T) {
	c := test.NewFakeClient()
	c.SetCredentials("+49111111111", "1111")
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	resp, err := c.TimelineDetails("1234567890", 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Error) != 0 {
		t.Fatal(resp.Error)
	}

}

func Portfolio(t *testing.T) {
	c := test.NewFakeClient()
	c.SetCredentials("+49111111111", "1111")
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	resp, err := c.Positions("1234567890")
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Error) != 0 {
		t.Fatal(resp.Error)
	}

}

func SavingsPlans(t *testing.T) {
	c := test.NewFakeClient()
	c.SetCredentials("+49111111111", "1111")
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	resp, err := c.SavingsPlans("1234567890")
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Error) != 0 {
		t.Fatal(resp.Error)
	}

}

func waitForGrpcServerToBeUp(addr string, limit int) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTradePipeClient(conn)

	for i := 0; i < limit; i++ {
		_, err := client.Alive(context.Background(), &emptypb.Empty{})
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout waiting for gRPC server to be up")
}
