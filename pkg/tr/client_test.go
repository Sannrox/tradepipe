package tr

import (
	"testing"

	fake "github.com/Sannrox/tradepipe/pkg/tr/testing"
)

func TestLogin(t *testing.T) {
	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	go FakeServer.Run()
	client := NewAPIClient()

	client.SetBaseURL("http://localhost:8080")
	client.SetWSBaseURL("wss://localhost:8080")

	client.SetCredentials("+49111111111", "1111")

	err := client.Login()
	if err != nil {
		t.Fatal(err)
	}

	if client.ProcessID != "1234567890" {
		t.Fatal("ProcessId not set")
	}

}

func TestVerify(t *testing.T) {
	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
	go FakeServer.Run()
	client := NewAPIClient()

	client.SetBaseURL("http://localhost:8080")
	client.SetWSBaseURL("wss://localhost:8080")

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

// func TestTimeline(t *testing.T) {
// 	FakeServer := fake.NewFakeServer("+49111111111", "1111", "1234567890", "1234")
// 	go FakeServer.Run()
// 	client := NewAPIClient()

// 	client.SetBaseURL("http://localhost:8080")
// 	client.SetWSBaseURL("wss://localhost:8080")

// 	client.SetCredentials("+49111111111", "1111")

// 	err := client.Login()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err = client.VerifyLogin(1234)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	timeline, err := client.AllTimeline()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// }
