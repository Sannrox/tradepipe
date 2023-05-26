package client

import (
	"context"

	pb "github.com/Sannrox/tradepipe/gear/protobuf"
	"github.com/Sannrox/tradepipe/gear/protobuf/login"
	portfolio "github.com/Sannrox/tradepipe/gear/protobuf/portfolio"
	"github.com/Sannrox/tradepipe/gear/protobuf/savingsplan"
	"github.com/Sannrox/tradepipe/gear/protobuf/timeline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	Opts []grpc.DialOption
	Conn *grpc.ClientConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(serverAddr string) error {
	c.Opts = append(c.Opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	Conn, err := grpc.Dial(serverAddr, c.Opts...)
	if err != nil {
		return err
	}
	c.Conn = Conn
	return nil
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Login(number, pin string) (*login.ProcessId, error) {
	client := pb.NewTradePipeClient(c.Conn)
	return client.Login(context.Background(), &login.Credentials{Number: number, Pin: pin})
}

func (c *Client) Verify(processId string, verifyCode int32) (*login.TwoFAReturn, error) {
	client := pb.NewTradePipeClient(c.Conn)
	return client.Verify(context.Background(), &login.TwoFAAsks{ProcessId: processId, VerifyCode: verifyCode})
}

func (c *Client) Positions(processId string) (*portfolio.ResponsePortfolio, error) {
	client := pb.NewTradePipeClient(c.Conn)
	return client.ReadPortfolio(context.Background(), &portfolio.RequestPortfolio{ProcessId: processId})
}

func (c *Client) UpdatePositions(processId string) error {
	client := pb.NewTradePipeClient(c.Conn)
	_, err := client.UpdatePortfolio(context.Background(), &portfolio.RequestPortfolioUpdate{ProcessId: processId})
	return err
}

func (c *Client) SavingsPlans(processId string) (*savingsplan.ResponseSavingsplan, error) {
	client := pb.NewTradePipeClient(c.Conn)
	return client.ReadSavingsPlans(context.Background(), &savingsplan.RequestSavingsplan{ProcessId: processId})
}

func (c *Client) UpdateSavingsPlans(processId string) error {
	client := pb.NewTradePipeClient(c.Conn)
	_, err := client.UpdateSavingsPlans(context.Background(), &savingsplan.RequestSavingsplanUpdate{ProcessId: processId})
	return err
}

func (c *Client) Timeline(processId string) (*timeline.ResponseTimeline, error) {
	client := pb.NewTradePipeClient(c.Conn)
	return client.ReadTimeline(context.Background(), &timeline.RequestTimeline{ProcessId: processId})
}

func (c *Client) UpdateTimeline(processId string) error {
	client := pb.NewTradePipeClient(c.Conn)
	_, err := client.UpdateTimeline(context.Background(), &timeline.RequestTimelineUpdate{ProcessId: processId})
	return err
}

func (c *Client) TimelineDetails(processId string) (*timeline.ResponseTimelineDetails, error) {
	client := pb.NewTradePipeClient(c.Conn)
	return client.ReadTimelineDetails(context.Background(), &timeline.RequestTimelineDetails{ProcessId: processId})
}

func (c *Client) Status() error {
	client := pb.NewTradePipeClient(c.Conn)
	_, err := client.Status(context.Background(), &emptypb.Empty{})
	return err
}
