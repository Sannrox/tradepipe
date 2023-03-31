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
	"github.com/Sannrox/tradepipe/logger"
	"github.com/Sannrox/tradepipe/scylla/tr_storage"
	"github.com/Sannrox/tradepipe/scylla/users"

	"github.com/Sannrox/tradepipe/grpc/pb/login"
	"github.com/Sannrox/tradepipe/grpc/pb/portfolio"
	"github.com/Sannrox/tradepipe/grpc/pb/savingsplan"
	"github.com/Sannrox/tradepipe/grpc/pb/timeline"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	user_keyspace        = "user"
	portfolio_keyspace   = "portfolio"
	savingsPlan_keyspace = "savingsplan"
)

type GRPCServer struct {
	*pb.UnimplementedTradePipeServer
	client         map[string]*tr.APIClient
	Lock           sync.Mutex
	baseURL        string
	wsURL          string
	overWriteTls   bool
	baseHTTPClient *http.Client
	Keyspaces
}

type Keyspaces struct {
	User         *users.User
	Portfolio    *tr_storage.Portfolios
	Savingsplans *tr_storage.SavingsPlans
}

var port = flag.Int("port", 50051, "The server port")

func NewGRPCServer(dbhost string) *GRPCServer {
	return &GRPCServer{client: make(map[string]*tr.APIClient),
		baseURL:      "",
		wsURL:        "",
		overWriteTls: false,
		Keyspaces: Keyspaces{
			User:         users.NewUserKeyspace(dbhost, user_keyspace),
			Portfolio:    tr_storage.NewPortfolioKeyspace(dbhost, portfolio_keyspace),
			Savingsplans: tr_storage.NewSavingsPlanKeyspace(dbhost, savingsPlan_keyspace),
		},
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
	logrus.Debug("Run Login for ", in.Number, "with pin", in.Pin)
	client := tr.NewAPIClient()

	if s.baseHTTPClient != nil {
		client.SetHTTPClient(s.baseHTTPClient)
	}
	if len(s.baseURL) != 0 {
		client.SetBaseURL(s.baseURL)
	}

	if len(s.wsURL) != 0 {
		client.SetWSBaseURL(s.wsURL)
	}
	if s.overWriteTls {
		logrus.Info("Overwriting TLS")
		client.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}
	client.SetCredentials(in.Number, in.Pin)

	err := client.Login()
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Login failed")
	}

	s.Lock.Lock()
	s.client[client.ProcessID] = client
	s.Lock.Unlock()

	return &login.ProcessId{ProcessId: client.GetProcessID()}, nil
}

func (s *GRPCServer) Verify(ctx context.Context, in *login.TwoFAAsks) (*login.TwoFAReturn, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	err := client.VerifyLogin(int(in.VerifyCode))
	if err != nil {
		return nil, err
	}

	if s.User.CheckIfUserExists(client.Creds.Number) {
		currentPin := s.User.Users.Pin(client.Creds.Number)
		if *currentPin != client.Creds.Pin {
			logrus.Debug("Updating pin")
			s.User.UpdateUser(client.Creds.Number, client.Creds.Pin)
		}
	} else {
		s.User.AddUser(client.Creds.Number, client.Creds.Pin)
	}

	return &login.TwoFAReturn{}, nil
}

func (s *GRPCServer) Timeline(ctx context.Context, in *timeline.RequestTimeline) (*timeline.ResponseTimeline, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()
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
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

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
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()
	data := make(chan tr.Message)

	user, err := s.User.ReadUser(client.Creds.Number)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error reading user from database")
	}

	err = s.Portfolio.CreateNewPortfolioTable(user.ID.String())
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error creating new portfolio table")
	}

	positions, err := s.Portfolio.GetAllPositions(user.ID.String())
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error getting all positions from database")
	}

	err = client.NewWebSocketConnection(data)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error creating new websocket connection")
	}

	time.Sleep(10 * time.Second)
	p := tr.NewPortfolioLoader(client)
	err = p.LoadPortfolio(ctx, data)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	bytes, err := p.GetPositionsAsBytes()
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error getting positions as bytes")
	}

	newPositions := p.GetPositions()
	err = s.Portfolio.AddPositions(user.ID.String(), &newPositions)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error adding positions to database")
	}

	for _, pos := range newPositions {
		positions = append(positions, &pos)
	}

	return &portfolio.ResponsePositions{
		ProcessId: in.ProcessId,
		Postions:  bytes,
	}, nil
}

func (s *GRPCServer) SavingsPlans(ctx context.Context, in *savingsplan.RequestSavingsplan) (*savingsplan.ResponseSavingsplan, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	user, err := s.User.ReadUser(client.Creds.Number)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error reading user from database")
	}

	err = s.Savingsplans.CreateNewSavingsPlanTable(user.ID.String())
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error creating new savingsplan table")
	}

	savingsplans, err := s.Savingsplans.GetAllSavingsPlans(user.ID.String())
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error getting all savingsplans from database")
	}

	err = client.NewWebSocketConnection(data)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error creating new websocket connection")
	}
	time.Sleep(10 * time.Second)
	p := tr.NewSavingsPlan(client)
	err = p.LoadSavingsplans(ctx, data)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error loading savingsplans")
	}
	bytes, err := p.GetSavingsPlansAsBytes()
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error getting savingsplans as bytes")
	}

	newSavingsplans := p.GetSavingsPlans()
	err = s.Savingsplans.AddSavingsPlans(user.ID.String(), &newSavingsplans)
	if err != nil {
		return nil, logger.ErrorWrapper(err, "Error adding savingsplans to database")
	}

	for _, pos := range newSavingsplans {
		savingsplans = append(savingsplans, &pos)
	}

	return &savingsplan.ResponseSavingsplan{
		ProcessId:    in.ProcessId,
		Savingsplans: bytes,
	}, nil

}

func (s *GRPCServer) Run(done chan struct{}) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	s.User.CreateNewUserTable()
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
