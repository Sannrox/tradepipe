package tr

import (
	"testing"

	fake "github.com/Sannrox/tradepipe/helper/testhelpers/faketrserver"
	"github.com/sirupsen/logrus"
)

const FakeServerPort string = "443"

func TestClient(t *testing.T) {
	done := make(chan struct{})
	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	FakeServer.GenerateData()
	go FakeServer.Run(done, FakeServerPort, "../test/ssl/cert.pem", "../test/ssl/key.pem")

	t.Run("TestLogin", Login)
	t.Run("TestVerify", Verify)
	t.Run("TestTimeline", Timeline)

	close(done)

}
func Login(t *testing.T) {
	client := NewAPIClient()

	client.SetHTTPClient(fake.OverWriteClient())
	client.SetTLSConfig(fake.OverWriteTSLClientConfig())
	client.SetBaseURL("https://localhost:443")
	client.SetWSBaseURL("wss://localhost:443")

	client.SetCredentials("+49111111111", "1111")

	err := client.Login()
	if err != nil {
		t.Fatal(err)
	}

	if client.ProcessID != "1234567890" {
		t.Fatal("ProcessId not set")
	}
	logrus.Info("Successful")

}

func Verify(t *testing.T) {
	client := NewAPIClient()

	client.SetHTTPClient(fake.OverWriteClient())
	client.SetTLSConfig(fake.OverWriteTSLClientConfig())
	client.SetBaseURL("https://localhost:443")
	client.SetWSBaseURL("wss://localhost:443")

	client.SetCredentials("+49111111111", "1111")
	err := client.Login()
	if err != nil {
		t.Fatal(err)
	}

	err = client.VerifyLogin(1234)
	if err != nil {
		t.Fatal(err)
	}

}

func Timeline(t *testing.T) {
	client := NewAPIClient()

	client.SetHTTPClient(fake.OverWriteClient())
	client.SetTLSConfig(fake.OverWriteTSLClientConfig())
	client.SetBaseURL("https://localhost:443")
	client.SetWSBaseURL("wss://localhost:443")

	client.SetCredentials("+49111111111", "1111")

	err := client.Login()
	if err != nil {
		t.Fatal(err)
	}

	err = client.VerifyLogin(1234)
	if err != nil {
		t.Fatal(err)
	}

	data := make(chan Message)

	err = client.NewWebSocketConnection(data)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Timeline("")
	if err != nil {
		t.Fatal(err)
	}

	<-data

}
