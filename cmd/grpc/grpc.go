package grpc

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	pb "github.com/Sannrox/tradepipe/grpc"
	"github.com/Sannrox/tradepipe/grpc/login"
	"github.com/Sannrox/tradepipe/grpc/portfolio"
	"github.com/Sannrox/tradepipe/grpc/timeline"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	*pb.UnimplementedTradePipeServer
	client map[string]*tr.APIClient
}

var port = flag.Int("port", 50051, "The server port")

func (s *GRPCServer) Login(ctx context.Context, in *login.Credentials) (*login.ProcessId, error) {
	client := tr.NewAPIClient()
	client.SetCredentials(in.Number, in.Pin)

	err := client.Login()
	if err != nil {
		return nil, err
	}

	s.client[client.GetProcessID()] = client

	return &login.ProcessId{ProcessId: client.GetProcessID()}, nil
}

func (s *GRPCServer) Verify(ctx context.Context, in *login.TwoFAAsks) (*login.TwoFAReturn, error) {
	client := s.client[in.ProcessId]
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

func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	logrus.Infof("server listening at %v", lis.Addr())
	pb.RegisterTradePipeServer(server, s)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{client: make(map[string]*tr.APIClient)}
}
