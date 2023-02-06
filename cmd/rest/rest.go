package rest

import (
	"flag"
	"strconv"
	"sync"

	"github.com/Sannrox/tradepipe/pkg/rest"
	"github.com/Sannrox/tradepipe/pkg/tr"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

var Httpport = flag.Int("http-port", 8080, "The server port")

type RestServer struct {
	client map[string]*tr.APIClient
	Lock   sync.Mutex
}

func NewRestServer() *RestServer {
	return &RestServer{
		client: make(map[string]*tr.APIClient),
	}
}

func (s *RestServer) Login(ctx echo.Context) error {
	client := tr.NewAPIClient()

	var login rest.Login
	if err := ctx.Bind(&login); err != nil {
		return err
	}
	client.SetCredentials(*login.Number, *login.Pin)

	err := client.Login()
	if err != nil {
		return err
	}

	process := client.ProcessID

	s.Lock.Lock()
	s.client[process] = client
	s.Lock.Unlock()

	err = ctx.JSON(200, rest.ProcID{Token: &process})
	if err != nil {
		return err
	}

	return nil
}

func (s *RestServer) Verify(ctx echo.Context, processId string) error {
	s.Lock.Lock()
	client := s.client[processId]
	s.Lock.Unlock()

	var newVerify rest.Verify
	if err := ctx.Bind(&newVerify); err != nil {
		return err
	}

	intVar, err := strconv.Atoi(*newVerify.Token)
	if err != nil {
		return err
	}
	err = client.VerifyLogin(intVar)
	if err != nil {
		return err
	}

	err = ctx.JSON(200, rest.Verified{})
	if err != nil {
		return err
	}

	return nil
}

func (s *RestServer) DownloadAll(ctx echo.Context, processId string) error {
	s.Lock.Lock()
	client := s.client[processId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		return err
	}

	dl := tr.NewDownloader(*client)
	dl.DownloadAll(ctx.Request().Context(), data)

	err = ctx.JSON(200, rest.Downloaded{})
	if err != nil {
		return err
	}

	return nil
}

func (s *RestServer) Run() error {
	swagger, err := rest.GetSwagger()
	if err != nil {
		return err
	}

	swagger.Servers = nil

	api := NewRestServer()

	e := echo.New()

	e.Use(echomiddleware.Logger())

	e.Use(middleware.OapiRequestValidator(swagger))

	rest.RegisterHandlers(e, api)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(*Httpport)))

	return nil
}
