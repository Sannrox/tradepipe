package fakegrpcclient

import (
	"context"

	pb "github.com/Sannrox/tradepipe/grpc"
	"github.com/Sannrox/tradepipe/grpc/login"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr = "localhost:50051"

type FakeClient struct {
	opts []grpc.DialOption
	conn *grpc.ClientConn
}

func NewFakeClient() *FakeClient {
	return &FakeClient{}
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

func (c *FakeClient) Login(number, pin string) (*login.ProcessId, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Login(context.Background(), &login.Credentials{Number: number, Pin: pin})
}

func (c *FakeClient) Verify(processId string, verifyCode int32) (*login.TwoFAReturn, error) {
	client := pb.NewTradePipeClient(c.conn)
	return client.Verify(context.Background(), &login.TwoFAAsks{ProcessId: processId, VerifyCode: verifyCode})
}
