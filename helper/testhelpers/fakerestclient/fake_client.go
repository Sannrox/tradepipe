package fakerestclient

import (
	"context"
	"net/http"

	"github.com/Sannrox/tradepipe/rest/api"
)

type FakeClient struct {
	Client   http.Client
	BaseURL  string
	Username string
	Password string
}

func NewFakeClient() *FakeClient {
	return &FakeClient{
		Client:   http.Client{},
		BaseURL:  "",
		Username: "",
		Password: "",
	}

}

func (c *FakeClient) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *FakeClient) SetCredentials(username, password string) {
	c.Username = username
	c.Password = password
}

func (c *FakeClient) SetupClient(restClient *http.Client) (*api.ClientWithResponses, error) {
	return api.NewClientWithResponses(c.BaseURL, api.WithHTTPClient(restClient))
}

func (c *FakeClient) Login(ctx context.Context) (string, error) {
	restClient, err := c.SetupClient(&c.Client)
	if err != nil {
		return "", err
	}

	resp, err := restClient.LoginWithResponse(ctx, api.Login{
		Number: c.Username,
		Pin:    c.Password,
	})
	if err != nil {
		return "", err
	}

	return *resp.JSON200, nil
}

func (c *FakeClient) Verify(ctx context.Context, processId string, token string) error {
	restClient, err := c.SetupClient(&c.Client)
	if err != nil {
		return err
	}

	_, err = restClient.VerifyWithResponse(ctx, processId, api.Verify{
		Token: token,
	},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *FakeClient) Timeline(ctx context.Context, processId string, sinceTimestamp float64) (*api.Timeline, error) {
	restClient, err := c.SetupClient(&c.Client)
	if err != nil {
		return nil, err
	}

	var timestamp *api.TimelineParams
	if sinceTimestamp != 0 {
		timestamp.Since = &sinceTimestamp
	}
	resp, err := restClient.TimelineWithResponse(ctx, processId, timestamp)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *FakeClient) TimelineDetails(ctx context.Context, processId string, sinceTimestamp float64) (*api.TimelineDetails, error) {
	restClient, err := c.SetupClient(&c.Client)
	if err != nil {
		return nil, err
	}

	var timestamp *api.TimelineDetailsParams
	if sinceTimestamp != 0 {
		timestamp.Since = &sinceTimestamp
	}

	resp, err := restClient.TimelineDetailsWithResponse(ctx, processId, timestamp)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil

}
