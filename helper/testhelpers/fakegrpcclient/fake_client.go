package fakegrpcclient

import (
	"context"

	"github.com/Sannrox/tradepipe/grpc/pb"
	"github.com/Sannrox/tradepipe/grpc/pb/login"
	"github.com/Sannrox/tradepipe/grpc/pb/portfolio"
	"github.com/Sannrox/tradepipe/grpc/pb/timeline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr = "localhost:50051"

type FakeClient struct {
	opts   []grpc.DialOption
	conn   *grpc.ClientConn
	Number string
	PIN    string
}

func NewFakeClient() *FakeClient {
	return &FakeClient{}
}

func (c *FakeClient) SetCredentials(number, pin string) {
	c.Number = number
	c.PIN = pin
}

func (c *FakeClient) Connect() error {
	c.opts = append(c.opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(serverAddr, c.opts...)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *FakeClient) Close() {
	c.conn.Close()
}

func (c *FakeClient) Login() (*login.ProcessId, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Login(context.Background(), &login.Credentials{Number: c.Number, Pin: c.PIN})
}

func (c *FakeClient) Verify(processId string, verifyCode int32) (*login.TwoFAReturn, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Verify(context.Background(), &login.TwoFAAsks{ProcessId: processId, VerifyCode: verifyCode})
}

func (c *FakeClient) Timeline(processId string, sinceTimestamp float64) (*timeline.ResponseTimeline, error) {
	client := pb.NewTradePipeClient(c.conn)
	if sinceTimestamp == 0 {
		return client.Timeline(context.Background(), &timeline.RequestTimeline{ProcessId: processId})
	} else {
		return client.Timeline(context.Background(), &timeline.RequestTimeline{ProcessId: processId,
			Timeline: &timeline.RequestTimeline_SinceTimestamp{SinceTimestamp: sinceTimestamp}})
	}
}

func (c *FakeClient) TimelineDetails(processId string, sinceTimestamp float64) (*timeline.ResponseTimeline, error) {
	client := pb.NewTradePipeClient(c.conn)
	if sinceTimestamp == 0 {
		return client.TimelineDetails(context.Background(), &timeline.RequestTimeline{ProcessId: processId})
	} else {
		return client.TimelineDetails(context.Background(), &timeline.RequestTimeline{ProcessId: processId,
			Timeline: &timeline.RequestTimeline_SinceTimestamp{SinceTimestamp: sinceTimestamp}})
	}
}

func (c *FakeClient) Positions(processId string) (*portfolio.ResponsePositions, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Positions(context.Background(), &portfolio.RequestPositions{ProcessId: processId})
}
