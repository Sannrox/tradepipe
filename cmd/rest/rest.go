package rest

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Sannrox/tradepipe/rest"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type RestServer struct {
	client         map[string]*tr.APIClient
	Lock           sync.Mutex
	baseURL        string
	wsURL          string
	overWriteTls   bool
	baseHTTPClient *http.Client
}

func NewRestServer() *RestServer {
	return &RestServer{
		client:       make(map[string]*tr.APIClient),
		baseURL:      "",
		wsURL:        "",
		overWriteTls: false,
	}
}

func (s *RestServer) SetBaseURL(url string) {
	s.baseURL = url
}

func (s *RestServer) SetWsURL(url string) {
	s.wsURL = url
}
func (s *RestServer) SetBaseHTTPClient(client *http.Client) {
	s.baseHTTPClient = client
}

func (s *RestServer) SetOverWriteTls(overWriteTls bool) {
	s.overWriteTls = overWriteTls
}

func (s *RestServer) Alive(ctx echo.Context) error {
	res := rest.AliveResponse{}
	time := time.Now().Unix()
	status := "OK"
	alive := rest.Alive{
		ServerTime: &time,
		Status:     &status,
	}

	res.JSON200 = &alive
	return ctx.JSON(200, res)
}

func (s *RestServer) Login(ctx echo.Context) error {
	client := tr.NewAPIClient()

	if s.baseHTTPClient != nil {
		client.SetHTTPClient(s.baseHTTPClient)
	}
	if s.baseURL != "" {
		client.SetBaseURL(s.baseURL)
	}

	if s.wsURL != "" {
		client.SetWSBaseURL(s.wsURL)
	}
	if s.overWriteTls {
		logrus.Info("Overwriting TLS")
		client.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	var login rest.Login
	if err := ctx.Bind(&login); err != nil {
		return err
	}
	client.SetCredentials(login.Number, login.Pin)

	err := client.Login()
	if err != nil {
		return err
	}

	s.Lock.Lock()
	s.client[client.ProcessID] = client
	s.Lock.Unlock()

	return ctx.JSON(200, client.ProcessID)

}

func (s *RestServer) Verify(ctx echo.Context, processId string) error {
	s.Lock.Lock()
	client := s.client[processId]
	s.Lock.Unlock()

	var newVerify rest.Verify
	if err := ctx.Bind(&newVerify); err != nil {
		return err
	}

	intVar, err := strconv.Atoi(newVerify.Token)
	if err != nil {
		return err
	}
	err = client.VerifyLogin(intVar)
	if err != nil {
		return err
	}

	return ctx.JSON(200, rest.Verified{})

}

func (s *RestServer) Timeline(ctx echo.Context, processId string, params rest.TimelineParams) error {
	s.Lock.Lock()
	client := s.client[processId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		return ctx.JSON(500, err)
	}

	tl := tr.NewTimeLine(client)
	if params.Since != nil {
		tl.SetSinceTimestamp(int64(*params.Since))
	}

	err = tl.LoadTimeLine(context.Background(), data)
	if err != nil {
		return ctx.JSON(500, err)
	}

	timeline := tl.GetTimeLineEvents()
	resptimeline := rest.Timeline{}

	b, err := json.Marshal(timeline)
	if err != nil {
		return ctx.JSON(500, err)
	}
	err = json.Unmarshal(b, &resptimeline.Timeline)
	if err != nil {
		return ctx.JSON(500, err)
	}

	return ctx.JSON(200, resptimeline)
}

func (s *RestServer) TimelineDetails(ctx echo.Context, processId string, params rest.TimelineDetailsParams) error {
	s.Lock.Lock()
	client := s.client[processId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		logrus.Debug(err)
		return ctx.JSON(500, err)
	}

	tl := tr.NewTimeLine(client)
	if params.Since != nil {
		tl.SetSinceTimestamp(int64(*params.Since))
	}

	err = tl.LoadTimeLine(context.Background(), data)
	if err != nil {
		logrus.Debug(err)
		return ctx.JSON(500, err)
	}

	err = tl.LoadTimeLineDetails(context.Background(), data)
	if err != nil {
		logrus.Debug(err)
		return ctx.JSON(500, err)
	}

	timelineDetails := tl.GetTimeLineDetails()

	response := rest.TimelineDetails{}

	b, err := json.Marshal(timelineDetails)
	if err != nil {
		logrus.Debug(err)
		return ctx.JSON(500, err)
	}
	err = json.Unmarshal(b, &response.TimelineDetails)
	if err != nil {
		logrus.Debug(err)
		return ctx.JSON(500, err)
	}

	return ctx.JSON(200, response)
}

func (s *RestServer) Positions(ctx echo.Context, processId string) error {
	s.Lock.Lock()
	client := s.client[processId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		logrus.Debug(err)
		return ctx.JSON(500, err)
	}
	postions := tr.NewPortfolio(client)

	err = postions.LoadPortfolio(context.Background(), data)
	if err != nil {
		return ctx.JSON(500, err)
	}

	positions, err := postions.GetPositionsAsBytes()
	if err != nil {
		return ctx.JSON(500, err)
	}

	response := rest.Positions{}
	err = json.Unmarshal(positions, &response.Positions)
	if err != nil {
		return ctx.JSON(500, err)
	}

	return ctx.JSON(200, response)
}

func (s *RestServer) Run(done chan struct{}, port string) error {
	swagger, err := rest.GetSwagger()
	if err != nil {
		return err
	}

	swagger.Servers = nil

	e := echo.New()

	e.Use(echoMiddleware.Logger())

	e.Use(middleware.OapiRequestValidator(swagger))

	rest.RegisterHandlers(e, s)

	var errChan chan error
	go func(err chan error) {
		err <- e.Start(":" + port)
	}(errChan)

	if err := <-errChan; err != nil {
		return err
	}
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logrus.Error(err)
	}
	return nil
}
