package testing

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sannrox/tradepipe/logger"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type FakeServer struct {
	Timeline        *FakeRawTimelines
	TimelineDetails *FakeTimelineDetails
	Subscriptios    *FakeSubscriptionStore
	Portfolio       *FakePortfolio
	SavingsPlans    *FakeSavingsplans
	verifyCode      string
	number          string
	pin             string
	ProcessId       string
	CertFile        string
	KeyFile         string
}

func NewFakeServer(number, pin, processId, verifyCode string) *FakeServer {
	return &FakeServer{
		Timeline:        NewFakeTimelines(),
		TimelineDetails: NewFakeTimelineDetails(),
		Subscriptios:    NewFakeSubscriptionStore(),
		Portfolio:       NewFakePortfolio(),
		SavingsPlans:    NewFakeSavingsplans(),
		number:          number,
		pin:             pin,
		ProcessId:       processId,
		verifyCode:      verifyCode,
		KeyFile:         "key.pem",
		CertFile:        "cert.pem",
	}
}
func OverWriteClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func OverWriteTSLClientConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func (s *FakeServer) GenerateData() {
	s.Timeline.GenerateRawTimelines(1)
	for _, timeline := range *s.Timeline {
		for _, event := range timeline.Data {
			s.TimelineDetails.GenerateTimelineDetail(&event)
		}
	}
	s.Portfolio.GenerateFakePortfolio()
	s.SavingsPlans.GenerateFakeSavingsPlans(1)
	logrus.Info("Fake Data generated")
}

func (s *FakeServer) Run(done chan struct{}, port int) {
	logger.Enable()
	logger.SetLogFile("fakeserver.log")
	http.HandleFunc("/", s.WebSocket)
	http.HandleFunc("/api/v1/auth/web/login", s.Login)
	http.HandleFunc("/api/v1/auth/web/login/", s.Verify)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	if err := s.CreateCertAndKeyForFakeServer(); err != nil {
		logrus.Fatal(err)
	}

	go func() {
		logrus.Infof("Fake Server starting on port %d", port)
		err := server.ListenAndServeTLS(s.CertFile, s.KeyFile)
		if err != nil && err != http.ErrServerClosed {
			logrus.Error(err)
		}
		if err != nil {
			logrus.Error(err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Error(err)
	}
	s.RemoveCertAndKeyForFakeServer()
	if err := os.Remove("fakeserver.log"); err != nil {
		logrus.Error(err)
	}
	logrus.Info("Fake Server stopped")
}

func (s *FakeServer) CreateCertAndKeyForFakeServer() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a self-signed certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return err
	}

	// Write the private key and certificate to files
	keyOut := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certOut := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// Save the files to disk
	if err := os.WriteFile(s.CertFile, certOut, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(s.KeyFile, keyOut, 0600); err != nil {
		return err
	}
	return nil
}

func (s *FakeServer) RemoveCertAndKeyForFakeServer() error {
	if err := os.Remove(s.CertFile); err != nil {
		return err
	}
	if err := os.Remove(s.KeyFile); err != nil {
		return err
	}
	return nil
}

//nolint:all
func (s *FakeServer) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	logrus.Info("Login with number and pin")
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if data["phoneNumber"] != s.number || data["pin"] != s.pin {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Got number: %s, pin: %s", s.number, s.pin)))
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    http.StatusOK,
		"processId": s.ProcessId})

}

func (s *FakeServer) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path

	pathElements := strings.Split(path, "/")

	if len(pathElements) != 8 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	processId := pathElements[6]
	verifyCode := pathElements[7]

	if processId != s.ProcessId {
		logrus.Info("Wrong processId: " + processId + " expected: " + s.ProcessId)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("wrong processId!"))
		return
	}

	if verifyCode != s.verifyCode {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}

type code int

const (
	A code = iota
	D
	C
	E
)

func (c code) String() string {
	switch c {
	case A:
		return "A"
	case D:
		return "D"
	case C:
		return "C"
	case E:
		return "E"
	default:
		return "Unknown"
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *FakeServer) Status(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is running!"))
}

func (s *FakeServer) WebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Warn(err)
		return
	}

	defer conn.Close()

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Error(err)
			return
		}

		parsedMessage := strings.SplitAfterN(string(msg), " ", 3)

		var returner string

		logrus.Info("Message received: " + string(msg))

		switch strings.TrimSpace(parsedMessage[0]) {
		case "connect":
			logrus.Info("Connection Incoming")
			err = conn.WriteMessage(websocket.TextMessage, []byte("connected"))
			if err != nil {
				logrus.Error(err)
				return
			}
			continue
		case "sub":
			logrus.Info("Subscription Incoming")
			subscriptionId, err := strconv.Atoi(strings.TrimSpace(parsedMessage[1]))
			if err != nil {
				logrus.Error(err)
				return
			}
			subscriptionType := strings.TrimSpace(parsedMessage[2])
			if !s.Subscriptios.Contains(subscriptionId) {
				s.Subscriptios.AddSubscription(subscriptionId)
				var message map[string]interface{}
				err := json.Unmarshal([]byte(subscriptionType), &message)
				if err != nil {
					err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
					if err != nil {
						logrus.Error(err)
						return
					}
					continue

				}
				switch message["type"] {
				case "timeline":
					logrus.Info("Timeline Subscription")
					var data FakeRawTimeline
					if message["after"] == nil {
						data = s.Timeline.First()
					} else {
						data = s.Timeline.Next(message["after"].(string))
					}
					timelineJSON, err := json.Marshal(data)

					if err != nil {
						err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
						if err != nil {
							logrus.Error(err)
							return
						}
						continue
					}
					returner = fmt.Sprintf("%d %s %s", subscriptionId, A.String(), string(timelineJSON))
				case "timelineDetail":
					logrus.Info("Timeline Detail Subscription")
					id := message["id"]
					if id == nil {
						err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), "No id provided")))
						if err != nil {
							logrus.Error(err)
							return
						}
						continue
					} else {
						detail := s.TimelineDetails.GenerateTimelineDetailById(id.(string))
						detailJSON, err := json.Marshal(detail)
						if err != nil {
							err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
							if err != nil {
								logrus.Error(err)
								return
							}
							continue
						}
						returner = fmt.Sprintf("%d %s %s", subscriptionId, A.String(), string(detailJSON))
					}

				case "compactPortfolio":
					logrus.Info("Portfolio Subscription")
					portfolio := s.Portfolio.GetPortfolio()
					logrus.Info(portfolio)
					portfolioJSON, err := json.Marshal(portfolio)
					if err != nil {
						err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
						if err != nil {
							logrus.Error(err)
							return
						}
						continue
					}
					returner = fmt.Sprintf("%d %s %s", subscriptionId, A.String(), string(portfolioJSON))

				case "savingsPlans":
					logrus.Info("Savings Plans Subscription")
					savingsPlans := s.SavingsPlans.GetSavingsPlans()
					savingsPlansJSON, err := json.Marshal(savingsPlans)
					if err != nil {
						err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
						if err != nil {
							logrus.Error(err)
							return
						}
						continue
					}
					returner = fmt.Sprintf("%d %s %s", subscriptionId, A.String(), string(savingsPlansJSON))

				default:
					returner = fmt.Sprintf("%d %s %s", subscriptionId, E.String(), "Unknown subscription type")
				}
			}
		case "unsub":
			logrus.Info("Unsubscription Incoming")
			subscriptionId, err := strconv.Atoi(strings.TrimSpace(parsedMessage[1]))
			if err != nil {
				logrus.Error(err)
				return
			}
			if s.Subscriptios.Contains(subscriptionId) {
				s.Subscriptios.Remove(subscriptionId)
			}

		}

		// Write message back to browser
		logrus.Info(returner)
		err = conn.WriteMessage(1, []byte(returner))
		if err != nil {
			logrus.Error(err)
			return
		}
	}

}
