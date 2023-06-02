package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/Sannrox/tradepipe/gear/protobuf"
	"github.com/Sannrox/tradepipe/gear/protobuf/login"
	portfolio "github.com/Sannrox/tradepipe/gear/protobuf/portfolio"
	"github.com/Sannrox/tradepipe/gear/protobuf/savingsplan"
	"github.com/Sannrox/tradepipe/gear/protobuf/timeline"
	"github.com/Sannrox/tradepipe/scylla"
	kp "github.com/Sannrox/tradepipe/scylla/keyspaces/portfolio"
	ks "github.com/Sannrox/tradepipe/scylla/keyspaces/savingsplans"
	"github.com/Sannrox/tradepipe/scylla/keyspaces/system"
	kt "github.com/Sannrox/tradepipe/scylla/keyspaces/timeline"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var port = flag.Int("port", 50051, "The server port")

type Server struct {
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
	System       *system.System
	Portfolio    *kp.Portfolios
	Savingsplans *ks.SavingsPlans
	Timelines    *kt.Timeline
}

func NewServer() *Server {

	return &Server{client: make(map[string]*tr.APIClient),
		baseURL:      "",
		wsURL:        "",
		overWriteTls: false,
		Keyspaces:    Keyspaces{},
	}
}

func (s *Server) SetBaseURL(url string) {
	s.baseURL = url
}

func (s *Server) SetWsURL(url string) {
	s.wsURL = url
}

func (s *Server) SetBaseHTTPClient(client *http.Client) {
	s.baseHTTPClient = client
}

func (s *Server) SetOverWriteTls(overWriteTls bool) {
	s.overWriteTls = overWriteTls
}

func (s *Server) Login(ctx context.Context, in *login.Credentials) (*login.ProcessId, error) {
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
		client.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}
	client.SetCredentials(in.Number, in.Pin)

	err := client.Login()
	if err != nil {
		logrus.Error("Login failed:", err)
		return nil, err
	}

	s.Lock.Lock()
	s.client[client.ProcessID] = client
	s.Lock.Unlock()

	return &login.ProcessId{ProcessId: client.GetProcessID()}, nil
}

func (s *Server) Verify(ctx context.Context, in *login.TwoFAAsks) (*login.TwoFAReturn, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	err := client.VerifyLogin(int(in.VerifyCode))
	if err != nil {
		return nil, err
	}

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Info("User not found, creating new user")
		if err := s.System.CreateUser(client.Creds.Number, client.Creds.Pin); err != nil {
			logrus.Error("Failed to create new user: ", err)
			return nil, err
		}
	} else if user.Pin != client.Creds.Pin {
		if err := s.System.UpdateUser(client.Creds.Number, client.Creds.Pin); err != nil {
			logrus.Error("Failed to update user: ", err)
			return nil, err
		}
	}
	return &login.TwoFAReturn{}, nil
}

func (s *Server) ReadTimeline(ctx context.Context, in *timeline.RequestTimeline) (*timeline.ResponseTimeline, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Failed to read user ", client.Creds.Number)
		return nil, err
	}

	timelineTable, err := s.Timelines.CreateTable(user.Id.String())
	if err != nil {
		logrus.Error("Failed to create timeline table: ", err)
		return nil, err
	}

	userTimeline, err := s.Timelines.GetCompleteTimeline(timelineTable)
	if err != nil {
		logrus.Error("Failed to read timeline: ", err)
		return nil, err
	}

	bytes, err := json.Marshal(userTimeline)
	if err != nil {
		logrus.Error("Failed to marshal timeline: ", err)
		return nil, err
	}

	return &timeline.ResponseTimeline{
		ProcessId: in.ProcessId,
		Items:     bytes,
	}, nil
}

func (s *Server) UpdateTimeline(ctx context.Context, in *timeline.RequestTimelineUpdate) (*emptypb.Empty, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Failed to read user ", client.Creds.Number)
		return nil, err
	}

	timelineTable, err := s.Timelines.CreateTable(user.Id.String())
	if err != nil {
		logrus.Error("Failed to create timeline table: ", err)
		return nil, err
	}

	timlineDetailTable, err := s.Timelines.CreateDetailTable(user.Id.String())
	if err != nil {
		logrus.Error("Failed to create timeline detail table: ", err)
		return nil, err
	}

	if err := client.NewWebSocketConnection(data); err != nil {
		logrus.Error("Failed to connect to websocket: ", err)
		return nil, err
	}

	time.Sleep(10 * time.Second)
	tl := tr.NewTimeLine(client)

	tl.SetSinceTimestamp(int64(in.GetSinceTimestamp()))
	if err = tl.LoadTimeLine(ctx, data); err != nil {
		logrus.Error("Failed to load timeline: ", err)
		return nil, err
	}

	newTimelines := tl.GetTimeLineEvents()
	if len(newTimelines) == 0 {
		logrus.Info("No new timeline events")
		return &emptypb.Empty{}, nil
	}
	for _, v := range newTimelines {
		exits := s.Timelines.CheckIfTimelineEventExists(timelineTable, v)
		if exits {
			if err := s.Timelines.UpdateTimelineEvent(timelineTable, v); err != nil {
				logrus.Error("Failed to update timeline: ", err)
				return nil, err
			}
		} else {
			if err := s.Timelines.AddTimelineEvent(timelineTable, v); err != nil {
				logrus.Error("Failed to insert new timeline: ", err)
				return nil, err
			}
		}
	}

	err = tl.LoadTimeLineDetails(ctx, data)
	if err != nil {
		logrus.Errorf("Failed to load timeline details: %v", err)
		return nil, err
	}

	newDetails := tl.GetTimeLineDetails()

	if len(newDetails) == 0 {
		logrus.Info("No new timeline details")
		return &emptypb.Empty{}, nil
	}
	for _, v := range newDetails {
		exists := s.Timelines.CheckIfTimelineDetailExists(timlineDetailTable, v)
		if exists {
			if err := s.Timelines.UpdateTimelineDetails(timlineDetailTable, v); err != nil {
				logrus.Error("Failed to update timeline details: ", err)
				return nil, err
			}
		} else {
			if err := s.Timelines.AddTimelineDetails(timlineDetailTable, v); err != nil {
				logrus.Error("Failed to insert new timeline details: ", err)
				return nil, err
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ReadTimelineDetails(ctx context.Context, in *timeline.RequestTimelineDetails) (*timeline.ResponseTimelineDetails, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Reading user from database ", err)
		return nil, err
	}

	timelineDetailTable, err := s.Timelines.CreateDetailTable(user.Id.String())
	if err != nil {
		logrus.Error("Creating timeline detail table ", err)
		return nil, err
	}

	userTimelineDetails, err := s.Timelines.GetAllTimelineDetails(timelineDetailTable)

	if err != nil {
		logrus.Error("Getting all timeline details from database ", err)
		return nil, err
	}

	bdetails, err := json.Marshal(userTimelineDetails)
	if err != nil {
		logrus.Error("Marshalling timeline details ", err)
		return nil, err
	}

	return &timeline.ResponseTimelineDetails{
		ProcessId: in.ProcessId,
		Items:     bdetails,
	}, nil
}

func (s *Server) ReadPortfolio(ctx context.Context, in *portfolio.RequestPortfolio) (*portfolio.ResponsePortfolio, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Reading user from database ", err)
		return nil, err
	}

	userPortfolio, err := s.Portfolio.CreateTable(user.Id.String())
	if err != nil {
		logrus.Error("Creating portfolio table ", err)
		return nil, err
	}

	positions, err := s.Portfolio.GetAllPositions(userPortfolio)
	if err != nil {
		logrus.Error("Getting all positions from database ", err)
		return nil, err
	}

	bpositions, err := json.Marshal(positions)
	if err != nil {
		logrus.Error("Marshalling positions ", err)
		return nil, err
	}

	return &portfolio.ResponsePortfolio{
		ProcessId: in.ProcessId,
		Positions: bpositions,
	}, nil
}

func (s *Server) UpdatePortfolio(ctx context.Context, in *portfolio.RequestPortfolioUpdate) (*emptypb.Empty, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	data := make(chan tr.Message)
	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Reading user from database ", err)
		return nil, err
	}

	userPortfolio, err := s.Portfolio.CreateTable(user.Id.String())
	if err != nil {
		logrus.Error("Creating portfolio table ", err)
		return nil, err
	}

	err = client.NewWebSocketConnection(data)
	if err != nil {
		logrus.Error("Creating new websocket connection ", err)
		return nil, err
	}

	time.Sleep(10 * time.Second)
	p := tr.NewPortfolioLoader(client)
	err = p.LoadPortfolio(ctx, data)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	newPositions := p.GetPositions()
	if len(newPositions) == 0 {
		logrus.Info("No new positions")
		return &emptypb.Empty{}, nil
	}

	for _, v := range newPositions {
		exists := s.Portfolio.CheckIfPositionExists(userPortfolio, v)
		if exists {
			if err := s.Portfolio.UpdatePosition(userPortfolio, v); err != nil {
				logrus.Error("Failed to update position: ", err)
				return nil, err
			} else {
				if err := s.Portfolio.AddPosition(userPortfolio, v); err != nil {
					logrus.Error("Failed to insert new position: ", err)
					return nil, err
				}
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ReadSavingsPlans(ctx context.Context, in *savingsplan.RequestSavingsplan) (*savingsplan.ResponseSavingsplan, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Reading user from database ", err)
		return nil, err
	}
	userSavingsplans, err := s.Savingsplans.CreateTable(user.Id.String())
	if err != nil {
		logrus.Error("Creating savingsplans table ", err)
		return nil, err
	}

	savingsplans, err := s.Savingsplans.GetAllPlans(userSavingsplans)
	if err != nil {
		logrus.Error("Getting all savingsplans from database ", err)
		return nil, err
	}

	bsavingsplans, err := json.Marshal(savingsplans)
	if err != nil {
		logrus.Error("Marshalling savingsplans ", err)
		return nil, err
	}

	return &savingsplan.ResponseSavingsplan{
		ProcessId:    in.ProcessId,
		Savingsplans: bsavingsplans,
	}, nil

}

func (s *Server) UpdateSavingsPlans(ctx context.Context, in *savingsplan.RequestSavingsplanUpdate) (*emptypb.Empty, error) {
	s.Lock.Lock()
	client := s.client[in.ProcessId]
	s.Lock.Unlock()

	data := make(chan tr.Message)

	user, err := s.System.GetUser(client.Creds.Number)
	if err != nil {
		logrus.Error("Reading user from database ", err)
		return nil, err
	}

	userSavingsplans, err := s.Savingsplans.CreateTable(user.Id.String())
	if err != nil {
		logrus.Error("Creating savingsplans table ", err)
		return nil, err
	}

	err = client.NewWebSocketConnection(data)
	if err != nil {
		logrus.Error("Creating new websocket connection ", err)
		return nil, err
	}
	time.Sleep(10 * time.Second)
	p := tr.NewSavingsPlan(client)
	err = p.LoadSavingsplans(ctx, data)
	if err != nil {
		logrus.Error("Failed to load savingsplans: ", err)
		return nil, err
	}

	newSavingsplans := p.GetSavingsPlans()
	if len(newSavingsplans) == 0 {
		logrus.Info("No new savingsplans")
		return &emptypb.Empty{}, nil
	}
	for _, v := range newSavingsplans {
		exists := s.Savingsplans.CheckIfPlanExists(userSavingsplans, v)
		if err != nil {
			logrus.Error("Failed to check if savingsplan exists: ", err)
			return nil, err
		}
		if exists {
			if err := s.Savingsplans.UpdatePlan(userSavingsplans, v); err != nil {
				logrus.Error("Failed to update savingsplan: ", err)
				return nil, err
			}
		} else {
			if err := s.Savingsplans.AddPlan(userSavingsplans, v); err != nil {
				logrus.Error("Failed to insert new savingsplan: ", err)
				return nil, err
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Status(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) CreateKeySpaceConnection(contactPoint string, port int, dbAttempts int, delay time.Duration) error {
	if err := scylla.TryToConnectWithRetry(contactPoint, port, dbAttempts, delay); err != nil {
		return err
	}

	var err error
	s.System, err = system.NewSystemKeyspace(contactPoint, port)
	if err != nil {

		return err
	}

	s.Portfolio, err = kp.NewPortfolioKeyspace(contactPoint, port)
	if err != nil {
		return err
	}

	s.Savingsplans, err = ks.NewSavingsPlanKeyspace(contactPoint, port)
	if err != nil {
		return err
	}

	s.Timelines, err = kt.NewTimelineKeyspace(contactPoint, port)
	if err != nil {
		return err
	}

	return nil

}

func (s *Server) CloseKeySpaceConnection() {
	s.System.Close()
	s.Portfolio.Close()
	s.Savingsplans.Close()
	s.Timelines.Close()
}

func (s *Server) Run(done chan struct{}, dbAddress string, dbPort, dbAttempts int, dbTimeOut time.Duration) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	if err := s.CreateKeySpaceConnection(dbAddress, dbPort, dbAttempts, dbTimeOut); err != nil {
		return err
	}

	if err := s.System.CreateTables(); err != nil {
		return err
	}

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
	s.CloseKeySpaceConnection()
	server.GracefulStop()

	return nil
}
