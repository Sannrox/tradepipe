package grpc

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	test "github.com/Sannrox/tradepipe/helper/testhelpers/fakegrpcclient"
	fake "github.com/Sannrox/tradepipe/helper/testhelpers/faketrserver"
)

var FakeTRServerPort string = "443"

func TestGrpcServer(t *testing.T) {
	done := make(chan struct{})
	s := NewGRPCServer()
	s.SetBaseURL("https://localhost:443")
	s.SetWsURL("wss://localhost:443")

	setClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	s.SetBaseHTTPClient(setClient)
	s.SetOverWriteTls(true)

	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	FakeServer.GenerateData()

	go FakeServer.Run(done, FakeTRServerPort, "../../test/ssl/cert.pem", "../../test/ssl/key.pem")
	go s.Run(done)

	time.Sleep(10 * time.Second)

	t.Run("Login test", Login)

	t.Run("Verify test", Verify)

	t.Run("Timeline test", Timeline)

	t.Run("TimelineDetail test", TimelineDetails)

	t.Run("Portfolio test", Portfolio)

	close(done)

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

	if resp.Error != "" {
		t.Fatal(resp.Error)
	}

	if resp.ProcessId == "" {
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

	if resp.Error != "" {
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

	if resp.Error != "" {
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

	if resp.Error != "" {
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

	if resp.Error != "" {
		t.Fatal(resp.Error)
	}

}
