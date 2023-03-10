package grpc

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Sannrox/tradepipe/grpc/pb"

	"github.com/Sannrox/tradepipe/grpc/pb/login"
	"github.com/Sannrox/tradepipe/grpc/pb/portfolio"
	"github.com/Sannrox/tradepipe/grpc/pb/timeline"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	*pb.UnimplementedTradePipeServer
	client         map[string]*tr.APIClient
	Lock           sync.Mutex
	baseURL        string
	wsURL          string
	overWriteTls   bool
	baseHTTPClient *http.Client
}

var port = flag.Int("port", 50051, "The server port")

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{client: make(map[string]*tr.APIClient),
		baseURL:      "",
		wsURL:        "",
		overWriteTls: false,
	}
}

func (s *GRPCServer) SetBaseURL(url string) {
	s.baseURL = url
}

func (s *GRPCServer) SetWsURL(url string) {
	s.wsURL = url
}

func (s *GRPCServer) SetBaseHTTPClient(client *http.Client) {
	s.baseHTTPClient = client
}

func (s *GRPCServer) SetOverWriteTls(overWriteTls bool) {
	s.overWriteTls = overWriteTls
}

func (s *GRPCServer) Alive(ctx context.Context, in *emptypb.Empty) (*pb.Alive, error) {
	res := pb.Alive{}
	time := time.Now().Unix()
	status := "OK"
	res.ServerTime = time
	res.Status = status
	return &res, nil
}

func (s *GRPCServer) Login(ctx context.Context, in *login.Credentials) (*login.ProcessId, error) {
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
	client.SetCredentials(in.GetNumber(), in.GetPin())

	err := client.Login()
	if err != nil {
		return nil, err
	}

	s.client[client.GetProcessID()] = client

	return &login.ProcessId{ProcessId: client.GetProcessID()}, nil
}

func (s *GRPCServer) Verify(ctx context.Context, in *login.TwoFAAsks) (*login.TwoFAReturn, error) {
	client := s.client[in.ProcessId]
	logrus.Debugf("%+v", client)
	err := client.VerifyLogin(int(in.VerifyCode))
	if err != nil {
		return nil, err
	}

	return &login.TwoFAReturn{}, nil
}

func (s *GRPCServer) Timeline(ctx context.Context, in *timeline.RequestTimeline) (*timeline.ResponseTimeline, error) {
	client := s.client[in.ProcessId]
	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)
	tl := tr.NewTimeLine(client)

	tl.SetSinceTimestamp(int64(in.GetSinceTimestamp()))
	err = tl.LoadTimeLine(ctx, data)
	if err != nil {
		return nil, err
	}
	bytes, err := tl.GetTimeLineEventsAsBytes()
	if err != nil {
		return nil, err
	}
	logrus.Debug(fmt.Sprintf("Timeline: %s", bytes))
	return &timeline.ResponseTimeline{
		ProcessId: in.ProcessId,
		Items:     bytes,
	}, nil
}

func (s *GRPCServer) TimelineDetails(ctx context.Context, in *timeline.RequestTimeline) (*timeline.ResponseTimeline, error) {
	client := s.client[in.ProcessId]
	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)
	tl := tr.NewTimeLine(client)

	tl.SetSinceTimestamp(int64(in.GetSinceTimestamp()))
	err = tl.LoadTimeLine(ctx, data)
	if err != nil {
		return nil, err
	}
	err = tl.LoadTimeLineDetails(ctx, data)
	if err != nil {
		return nil, err
	}
	bytes, err := tl.GetTimeLineDetailsAsBytes()
	if err != nil {
		return nil, err
	}
	return &timeline.ResponseTimeline{
		ProcessId: in.ProcessId,
		Items:     bytes,
	}, nil
}

func (s *GRPCServer) Positions(ctx context.Context, in *portfolio.RequestPositions) (*portfolio.ResponsePositions, error) {
	client := s.client[in.ProcessId]
	data := make(chan tr.Message)

	err := client.NewWebSocketConnection(data)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)
	p := tr.NewPortfolio(client)
	err = p.LoadPortfolio(ctx, data)
	if err != nil {
		return nil, err
	}
	bytes, err := p.GetPositionsAsBytes()
	if err != nil {
		return nil, err
	}
	return &portfolio.ResponsePositions{
		ProcessId: in.ProcessId,
		Postions:  bytes,
	}, nil
}

func (s *GRPCServer) Run(done chan struct{}) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	logrus.Infof("server listening at %v", lis.Addr())
	pb.RegisterTradePipeServer(server, s)
	var errChan chan error
	go func(err chan error) {
		err <- server.Serve(lis)

	}(errChan)
	if err := <-errChan; err != nil {
		return err
	}

	<-done
	server.GracefulStop()

	return nil
}
