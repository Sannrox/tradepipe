package tr

import (
	"fmt"
	"os"
	"testing"

	fake "github.com/Sannrox/tradepipe/helper/testhelpers/faketrserver"
	"github.com/Sannrox/tradepipe/helper/testhelpers/utils"
	"github.com/sirupsen/logrus"
)

func TestClient(t *testing.T) {
	done := make(chan struct{})
	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	FakeServer.GenerateData()
	port, err := utils.FindFreePort(3443, 4443)
	if err != nil {
		t.Fatal(err)
	}

	go FakeServer.Run(done, port)

	url := fmt.Sprintf("https://localhost:%d", port)
	if err := utils.WaitForRestServerToBeUp(url, 10); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TR_SERVER_URL", url)

	t.Run("TestLogin", Login)
	t.Run("TestVerify", Verify)
	t.Run("TestTimeline", Timeline)

	close(done)

}
func Login(t *testing.T) {
	client := NewAPIClient()

	client.SetHTTPClient(fake.OverWriteClient())
	client.SetTLSConfig(fake.OverWriteTSLClientConfig())
	client.SetBaseURL(os.Getenv("TR_SERVER_URL"))
	client.SetWSBaseURL(os.Getenv("TR_SERVER_URL"))

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
	client.SetBaseURL(os.Getenv("TR_SERVER_URL"))
	client.SetWSBaseURL(os.Getenv("TR_SERVER_URL"))

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
	client.SetBaseURL(os.Getenv("TR_SERVER_URL"))
	client.SetWSBaseURL(os.Getenv("TR_SERVER_URL"))

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
