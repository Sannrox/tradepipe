package testing

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type FakeServer struct {
	Timeline     *FakeTimeline
	Subscriptios *FakeSubscriptionStore
	verifyCode   string
	number       string
	pin          string
	ProcessId    string
}

func NewFakeServer(number, pin, processId, verifyCode string) *FakeServer {
	return &FakeServer{
		Timeline:     NewFakeTimeline(),
		Subscriptios: NewFakeSubscriptionStore(),
		number:       number,
		pin:          pin,
		ProcessId:    processId,
		verifyCode:   verifyCode,
	}
}

func (s *FakeServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/web/login", s.Login)
	mux.HandleFunc("/", s.WebSocket)
	logrus.Info("Fake Server started")
	http.ListenAndServe(":8080", mux)
}

func (s *FakeServer) Login(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Fake Login")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	if r.Body != nil && path != "/api/v1/auth/web/login" {

		pathElements := strings.Split(path, "/")

		if len(pathElements) != 5 {
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}

		processId := pathElements[3]
		verifyCode := pathElements[4]

		if processId != s.ProcessId {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if verifyCode != s.verifyCode {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {

		if r.Body != nil {
			var data map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if data["phoneNumber"] != s.number || data["pin"] != s.pin {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf("Got number: %s, pin: %s", s.number , s.pin)))
				return
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    http.StatusOK,
			"processId": s.ProcessId})
	}

}

func (s *FakeServer) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	path := r.URL.Path

	pathElements := strings.Split(path, "/")

	if len(pathElements) != 5 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	processId := pathElements[3]
	verifyCode := pathElements[4]

	if processId != s.ProcessId {
		w.WriteHeader(http.StatusUnauthorized)
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
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *FakeServer) WebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		parsedMessage := strings.SplitAfterN(string(msg), " ", 3)
		subscriptionId, err := strconv.Atoi(strings.TrimSpace(parsedMessage[1]))
		if err != nil {
			log.Println(err)
			return
		}

		subscriptionType := strings.TrimSpace(parsedMessage[2])

		var returner string
		if parsedMessage[0] == "sub" {
			if !s.Subscriptios.Contains(subscriptionId) {
				s.Subscriptios.AddSubscription(subscriptionId)
				var message map[string]interface{}
				err := json.Unmarshal([]byte(subscriptionType), &message)
				if err != nil {
					err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
					if err != nil {
						log.Println(err)
						return
					}
					continue

				}
				switch message["type"] {
				case "timeline":
					data := s.Timeline.First()
					timelineJSON, err := json.Marshal(data)
					if err != nil {
						err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d %s %s", subscriptionId, E.String(), err)))
						if err != nil {
							log.Println(err)
							return
						}
						continue
					}
					returner = fmt.Sprintf("%d %s %s", subscriptionId, A.String(), string(timelineJSON))
				}
			}
		}

		if parsedMessage[0] == "unsub" {
			if s.Subscriptios.Contains(subscriptionId) {
				s.Subscriptios.Remove(subscriptionId)
			}
		}

		// Write message back to browser
		err = conn.WriteMessage(1, []byte(returner))
		if err != nil {
			log.Println(err)
			return
		}
	}

}
