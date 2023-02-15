package rest

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	fakeClient "github.com/Sannrox/tradepipe/helper/testhelpers/fakerestclient"
	fake "github.com/Sannrox/tradepipe/helper/testhelpers/faketrserver"
	"github.com/Sannrox/tradepipe/rest"
)

var FakeTRServerPort string = "443"
var FakeHTTPServer string = "8080"

func TestRestServer(t *testing.T) {
	done := make(chan struct{})
	s := NewRestServer()
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
	go s.Run(done, FakeHTTPServer)

	time.Sleep(10 * time.Second)
	t.Run("Login test", Login)
	t.Run("Verify test", Verify)

	t.Run("Timeline test", Timeline)
	close(done)
}

func Login(t *testing.T) {
	fakeC := fakeClient.NewFakeClient()
	fakeC.SetCredentials("+49111111111", "1111")

	setClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	fakeC.SetBaseURL("http://localhost:" + FakeHTTPServer)
	client, err := fakeC.SetupClient(setClient)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fakeC.Username)
	t.Log(fakeC.Password)
	resp, err := client.LoginWithResponse(context.Background(), rest.Login{
		Number: fakeC.Username,
		Pin:    fakeC.Password,
	})
	if err != nil {
		t.Fatal(err)
	}

	if *resp.JSON200 != "1234567890" {
		t.Fatal("ProcessId not set")
	}
}

func Verify(t *testing.T) {
	fakeC := fakeClient.NewFakeClient()
	fakeC.SetCredentials("+49111111111", "1111")

	setClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	fakeC.SetBaseURL("http://localhost:" + FakeHTTPServer)
	client, err := fakeC.SetupClient(setClient)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.VerifyWithResponse(context.Background(), "1234567890", rest.Verify{
		Token: "1234",
	})
	if err != nil {
		t.Fatal(err)
	}

}

func Timeline(t *testing.T) {
	fakeC := fakeClient.NewFakeClient()
	fakeC.SetCredentials("+49111111111", "1111")

	setClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	fakeC.SetBaseURL("http://localhost:" + FakeHTTPServer)
	client, err := fakeC.SetupClient(setClient)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.TimelineWithResponse(context.Background(), "1234567890", &rest.TimelineParams{})
	if err != nil {
		t.Fatal(err)
	}

	if *resp.JSON200.Message != "" {
		t.Fatal(*resp.JSON200.Message)
	}
	if *resp.JSON200.Timeline != nil {
		t.Fatal(*resp.JSON200.Timeline)
	}

}
