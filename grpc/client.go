package grpc

import (
	"context"

	pb "github.com/Sannrox/tradepipe/grpc/pb"
	"github.com/Sannrox/tradepipe/grpc/pb/login"
	portfolio "github.com/Sannrox/tradepipe/grpc/pb/portfolio"
	"github.com/Sannrox/tradepipe/grpc/pb/savingsplan"
	"github.com/Sannrox/tradepipe/grpc/pb/timeline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	opts []grpc.DialOption
	conn *grpc.ClientConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(serverAddr string) error {
	c.opts = append(c.opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(serverAddr, c.opts...)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Login(number, pin string) (*login.ProcessId, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Login(context.Background(), &login.Credentials{Number: number, Pin: pin})
}

func (c *Client) Verify(processId string, verifyCode int32) (*login.TwoFAReturn, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Verify(context.Background(), &login.TwoFAAsks{ProcessId: processId, VerifyCode: verifyCode})
}

func (c *Client) Positions(processId string) (*portfolio.ResponsePositions, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Positions(context.Background(), &portfolio.RequestPositions{ProcessId: processId})
}

func (c *Client) SavingsPlans(processId string) (*savingsplan.ResponseSavingsplan, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.SavingsPlans(context.Background(), &savingsplan.RequestSavingsplan{ProcessId: processId})
}

func (c *Client) Timeline(processId string) (*timeline.ResponseTimeline, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Timeline(context.Background(), &timeline.RequestTimeline{ProcessId: processId})
}
