package wsconnector

import (
	"errors"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
)

// The basis is this connector "github.com/togettoyou/wsc" with various extensions
type WSConnector struct {
	Config *Config

	WebSocket *WebSocket

	onConnected func()

	onConnectError func(err error)

	onDisconnected func(err error)

	onClose func(code int, text string)

	onTextMessageSent func(message string)

	onBinaryMessageSent func(data []byte)

	onSentError func(err error)

	onPingReceived func(appData string)

	onPongReceived func(appData string)

	onTextMessageReceived func(message string)

	onBinaryMessageReceived func(data []byte)
}

type Config struct {
	WriteWait time.Duration

	MaxMessageSize int64

	MinRecTime time.Duration

	MaxRecTime time.Duration

	RecFactor float64

	MessageBufferSize int
}

type WebSocket struct {
	Url    string
	Conn   *websocket.Conn
	Dialer *websocket.Dialer

	RequestHeader http.Header
	HttpResponse  *http.Response

	isConnected bool

	connMutex *sync.RWMutex

	sendMutex *sync.Mutex

	sendChan chan *webSocketMessage
}

type webSocketMessage struct {
	t   int
	msg []byte
}

func NewWSConnector(url string) *WSConnector {
	return &WSConnector{
		Config: &Config{
			WriteWait:         10 * time.Second,
			MaxMessageSize:    512,
			MinRecTime:        2 * time.Second,
			MaxRecTime:        60 * time.Second,
			RecFactor:         1.5,
			MessageBufferSize: 256,
		},
		WebSocket: &WebSocket{
			Url:           url,
			Dialer:        websocket.DefaultDialer,
			RequestHeader: http.Header{},
			isConnected:   false,
			connMutex:     &sync.RWMutex{},
			sendMutex:     &sync.Mutex{},
		},
	}
}

func (ws *WSConnector) SetConfig(config *Config) {
	ws.Config = config
}

func (ws *WSConnector) SetDialer(dialer *websocket.Dialer) {
	ws.WebSocket.Dialer = dialer
}

func (ws *WSConnector) SetHeader(header http.Header) {
	ws.WebSocket.RequestHeader = header
}

func (ws *WSConnector) OnConnected(f func()) {
	ws.onConnected = f
}

func (ws *WSConnector) OnConnectError(f func(err error)) {
	ws.onConnectError = f
}

func (ws *WSConnector) OnDisconnected(f func(err error)) {
	ws.onDisconnected = f
}

func (ws *WSConnector) OnClose(f func(code int, text string)) {
	ws.onClose = f
}

func (ws *WSConnector) OnTextMessageSent(f func(message string)) {
	ws.onTextMessageSent = f
}

func (ws *WSConnector) OnBinaryMessageSent(f func(data []byte)) {
	ws.onBinaryMessageSent = f
}

func (ws *WSConnector) OnSentError(f func(err error)) {
	ws.onSentError = f
}

func (ws *WSConnector) OnPingReceived(f func(appData string)) {
	ws.onPingReceived = f
}

func (ws *WSConnector) OnPongReceived(f func(appData string)) {
	ws.onPongReceived = f
}

func (ws *WSConnector) OnTextMessageReceived(f func(message string)) {
	ws.onTextMessageReceived = f
}

func (ws *WSConnector) OnBinaryMessageReceived(f func(data []byte)) {
	ws.onBinaryMessageReceived = f
}

func (ws *WSConnector) Closed() bool {
	ws.WebSocket.connMutex.RLock()
	defer ws.WebSocket.connMutex.RUnlock()
	return !ws.WebSocket.isConnected
}

func (ws *WSConnector) Connect() {
	ws.WebSocket.sendChan = make(chan *webSocketMessage, ws.Config.MessageBufferSize)
	ticker := &backoff.Backoff{
		Min:    ws.Config.MinRecTime,
		Max:    ws.Config.MaxRecTime,
		Factor: ws.Config.RecFactor,
		Jitter: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())
	go func() {
		for {
			var err error
			nextRec := ticker.Duration()
			ws.WebSocket.Conn, ws.WebSocket.HttpResponse, err = ws.WebSocket.Dialer.Dial(ws.WebSocket.Url, ws.WebSocket.RequestHeader)
			if err != nil {
				if ws.onConnectError != nil {
					ws.onConnectError(err)
				}
				time.Sleep(nextRec)
				continue
			}
			ws.WebSocket.connMutex.Lock()
			ws.WebSocket.isConnected = true
			ws.WebSocket.connMutex.Unlock()

			if ws.onConnected != nil {
				ws.onConnected()
			}

			// ws.WebSocket.Conn.SetReadLimit(ws.Config.MaxMessageSize)

			defaultCloseHandler := ws.WebSocket.Conn.CloseHandler()
			ws.WebSocket.Conn.SetCloseHandler(func(code int, text string) error {
				result := defaultCloseHandler(code, text)
				ws.clean()
				if ws.onClose != nil {
					ws.onClose(code, text)
				}
				return result
			})
			defaultPingHandler := ws.WebSocket.Conn.PingHandler()
			ws.WebSocket.Conn.SetPingHandler(func(appData string) error {
				if ws.onPingReceived != nil {
					ws.onPingReceived(appData)
				}
				return defaultPingHandler(appData)
			})

			defaultPongHandler := ws.WebSocket.Conn.PongHandler()
			ws.WebSocket.Conn.SetPongHandler(func(appData string) error {
				if ws.onPongReceived != nil {
					ws.onPongReceived(appData)
				}
				return defaultPongHandler(appData)
			})

			go ws.listen()
			go ws.write()
			return
		}
	}()
}

func (ws *WSConnector) listen() {
	logrus.Infof("Start listing for messages from %s", ws.WebSocket.Url)
	for {
		messageType, message, err := ws.WebSocket.Conn.ReadMessage()
		if err != nil {
			if ws.onDisconnected != nil {
				ws.onDisconnected(err)
			}
			return

		}
		switch messageType {
		case websocket.TextMessage:
			if ws.onTextMessageReceived != nil {
				ws.onTextMessageReceived(string(message))
			}
			break
		case websocket.BinaryMessage:
			if ws.onBinaryMessageReceived != nil {
				ws.onBinaryMessageReceived(message)
			}
			break
		}
	}
}

func (ws *WSConnector) write() {
	logrus.Infof("Start writing messages to %s", ws.WebSocket.Url)
	for {
		select {
		case wsMsg, ok := <-ws.WebSocket.sendChan:
			if !ok {
				return
			}
			err := ws.send(wsMsg.t, wsMsg.msg)
			if err != nil {
				if ws.onSentError != nil {
					ws.onSentError(err)
				}
				continue
			}
			switch wsMsg.t {
			case websocket.CloseMessage:
				return
			case websocket.TextMessage:
				if ws.onTextMessageSent != nil {
					ws.onTextMessageSent(string(wsMsg.msg))
				}
			case websocket.BinaryMessage:
				if ws.onBinaryMessageSent != nil {
					ws.onBinaryMessageSent(wsMsg.msg)
				}
				break
			}
		}
	}
}

var (
	CloseErr  = errors.New("Connection closed")
	BufferErr = errors.New("Message buffer is full")
)

func (ws *WSConnector) SendTextMessage(message string) error {
	if ws.Closed() {
		return CloseErr
	}

	select {
	case ws.WebSocket.sendChan <- &webSocketMessage{
		t:   websocket.TextMessage,
		msg: []byte(message),
	}:
	default:
		return BufferErr
	}
	return nil
}

func (ws *WSConnector) SendTextBinary(data []byte) error {
	if ws.Closed() {
		return CloseErr
	}

	select {
	case ws.WebSocket.sendChan <- &webSocketMessage{
		t:   websocket.TextMessage,
		msg: data,
	}:
	default:
		return BufferErr
	}
	return nil
}

func (ws *WSConnector) send(messageType int, data []byte) error {
	ws.WebSocket.sendMutex.Lock()
	defer ws.WebSocket.sendMutex.Unlock()
	if ws.Closed() {
		return CloseErr
	}
	_ = ws.WebSocket.Conn.SetWriteDeadline(time.Now().Add(ws.Config.WriteWait))
	err := ws.WebSocket.Conn.WriteMessage(messageType, data)
	return err
}

func (ws *WSConnector) closeAndReConnect() {
	if ws.Closed() {
		return
	}
	ws.clean()
	go ws.Connect()
}

func (ws *WSConnector) CloseWithMsg(message string) {
	if ws.Closed() {
		return
	}
	_ = ws.send(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, message))
	ws.clean()
	if ws.onClose != nil {
		ws.onClose(websocket.CloseNormalClosure, message)
	}
}

func (ws *WSConnector) Close() {
	ws.CloseWithMsg("")
}

func (ws *WSConnector) clean() {
	if ws.Closed() {
		return
	}
	ws.WebSocket.connMutex.Lock()
	ws.WebSocket.isConnected = false
	_ = ws.WebSocket.Conn.Close()
	close(ws.WebSocket.sendChan)
	ws.WebSocket.connMutex.Unlock()
}
