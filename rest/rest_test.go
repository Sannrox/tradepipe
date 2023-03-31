package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	fakeClient "github.com/Sannrox/tradepipe/helper/testhelpers/fakerestclient"
	fake "github.com/Sannrox/tradepipe/helper/testhelpers/faketrserver"
	"github.com/Sannrox/tradepipe/helper/testhelpers/utils"
	"github.com/Sannrox/tradepipe/rest/api"
)

var FakeHTTPServer string = "8088"

func TestRestServer(t *testing.T) {
	done := make(chan struct{})
	s := NewRestServer()

	fakeTrServerPort, err := utils.FindFreePort(3443, 4443)
	if err != nil {
		t.Fatal(err)
	}
	// Set the server to use the fake server
	s.SetBaseURL(fmt.Sprintf("https://localhost:%d", fakeTrServerPort))
	s.SetWsURL(fmt.Sprintf("wss://localhost:%d", fakeTrServerPort))

	setClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	s.SetBaseHTTPClient(setClient)
	s.SetOverWriteTls(true)

	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	FakeServer.GenerateData()

	go FakeServer.Run(done, fakeTrServerPort)
	go s.Run(done, FakeHTTPServer)

	fakeTrServerUrl := fmt.Sprintf("https://localhost:%d", fakeTrServerPort)

	if err := utils.WaitForRestServerToBeUp(fakeTrServerUrl, 10); err != nil {
		t.Fatal(err)
	}

	if err := utils.WaitForRestServerToBeUp("http://localhost:"+FakeHTTPServer, 10); err != nil {
		t.Fatal(err)
	}

	t.Run("Login test", Login)
	t.Run("Verify test", Verify)

	t.Run("Timeline test", Timeline)

	t.Run("TimelineDetail test", TimelineDetails)

	t.Run("Portfolio test", Portfolio)
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
	resp, err := client.LoginWithResponse(context.Background(), api.Login{
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

	_, err = client.VerifyWithResponse(context.Background(), "1234567890", api.Verify{
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

	resp, err := client.TimelineWithResponse(context.Background(), "1234567890", &api.TimelineParams{})
	if err != nil {
		t.Fatal(err)
	}

	if resp.JSON500 != nil {
		t.Fatal("Error in response")
	}

	if resp.JSON401 != nil {
		t.Fatal("Error in response")
	}

	timeline := resp.JSON200.Timeline

	if len(timeline) == 0 {
		t.Fatal("No timeline data")
	}
	t.Log(timeline...)

}

func TimelineDetails(t *testing.T) {
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

	resp, err := client.TimelineDetailsWithResponse(context.Background(), "1234567890", &api.TimelineDetailsParams{})
	if err != nil {
		t.Fatal(err)
	}

	if resp.JSON500 != nil {
		t.Fatalf("Error in response %+v | %+v", resp.JSON500.Message, string(resp.Body))
	}

	if resp.JSON401 != nil {
		t.Fatal("Error in response")
	}

	timelineDetail := resp.JSON200.TimelineDetails

	if len(timelineDetail) == 0 {
		t.Fatal("No timeline data")
	}
	t.Log(timelineDetail...)
}

func Portfolio(t *testing.T) {
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

	resp, err := client.PositionsWithResponse(context.Background(), "1234567890")
	if err != nil {
		t.Fatal(err)
	}

	if resp.JSON500 != nil {
		t.Fatal("Error in response:", resp.JSON500)
	}

	if resp.JSON401 != nil {
		t.Fatal("Error in response")
	}

	positions := resp.JSON200.Positions

	if len(positions) == 0 {
		t.Fatal("No portfolio data")
	}
	t.Log(positions...)
}
