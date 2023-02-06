package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Sannrox/tradepipe/pkg/rest"
)

type FakeClient struct {
	Client http.Client
}

func NewFakeClient() *FakeClient {
	return &FakeClient{
		Client: http.Client{},
	}

}

func (c *FakeClient) Login(number, pin string) (*rest.ProcID, error) {

	req, err := http.NewRequest("POST", "http://localhost:8080/login", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(number, pin)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	var pId rest.ProcID
	json.NewDecoder(resp.Body).Decode(&pId)
	return &pId, nil

}

func (c *FakeClient) Verify(processId, token string) (*rest.ProcID, error) {

	req, err := http.NewRequest("GET", "http://localhost:8080/verify/"+token, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	var procID rest.ProcID
	json.NewDecoder(resp.Body).Decode(&procID)
	return &procID, nil

}

func (c *FakeClient) DownloadAll(processId string) (*rest.Downloaded, error) {

	req, err := http.NewRequest("GET", "http://localhost:8080/download-all/"+processId, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	var downloaded rest.Downloaded
	json.NewDecoder(resp.Body).Decode(&downloaded)
	return &downloaded, nil

}
