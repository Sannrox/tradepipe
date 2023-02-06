package tr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Sannrox/tradepipe/pkg/tls"
	wsconnector "github.com/Sannrox/tradepipe/pkg/websocket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Look into
// https://gitlab.com/marzzzello/pytr/-/blob/master/pytr/account.py
// https://github.com/Zarathustra2/TradeRepublicApi
const (
	TRADEREPUBLICURL = "https://api.traderepublic.com"
	TRADEREPUBLIWSS  = "wss://api.traderepublic.com"
)

type APIClient struct {
	Creds                 *LoginCreds
	ProcessID             string
	Session               []*http.Cookie
	Client                *http.Client
	WebSocket             *wsconnector.WSConnector
	local                 string `default:"de"`
	subscriptionIdCounter int
	subscribstions        map[int]interface{}
	previousRespones      map[int]interface{}
	URL                   string
	WSURL                 string
}

type LoginCreds struct {
	Number string `json:"phoneNumber"`
	Pin    string `json:"pin"`
}

// Example: "0 C {payload}"
type Message struct {
	SubscriptionId int
	Subscription   map[string]interface{}
	Payload        map[string]interface{}
}

type Response struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

func NewAPIClient() *APIClient {
	logrus.Debugf("Initialize client")

	return &APIClient{
		Creds:                 nil,
		ProcessID:             "",
		Session:               nil,
		Client:                &http.Client{},
		WebSocket:             nil,
		local:                 "",
		subscriptionIdCounter: 0,
		subscribstions:        make(map[int]interface{}),
		previousRespones:      make(map[int]interface{}),
		URL:                   TRADEREPUBLICURL,
		WSURL:                 TRADEREPUBLIWSS,
	}
}

func (api *APIClient) SetBaseURL(url string) {
	api.URL = url
}

func (api *APIClient) SetWSBaseURL(url string) {
	api.WSURL = url
}

func (api *APIClient) SetLocal(local string) {
	api.local = local
}

func (api *APIClient) SetCredentials(number, pin string) {
	creds := &LoginCreds{
		Number: number,
		Pin:    pin,
	}
	api.Creds = creds
}

func (api *APIClient) GetProcessID() string {
	return api.ProcessID
}

func (api *APIClient) Login() error {
	if api.Creds == nil {
		return fmt.Errorf("no credentials set")
	}
	url := fmt.Sprintf("%s/api/v1/auth/web/login", api.URL)

	payload, err := json.Marshal(api.Creds)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	res, err := api.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("code: %d : %+v", res.StatusCode, string(body))
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp interface{}

	json.Unmarshal(body, &resp)

	respBody := resp.(map[string]interface{})

	api.ProcessID = respBody["processId"].(string)

	return nil
}

func (api *APIClient) VerifyLogin(verifyCode int) error {
	url := fmt.Sprintf("%s/api/v1/auth/web/login/%s/%d", api.URL, api.ProcessID, verifyCode)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	res, err := api.Client.Do(req)
	if err != nil {
		return err
	}
	api.Session = res.Cookies()
	defer res.Body.Close()

	return nil
}

func (api *APIClient) NewWebSocketConnection(dataChan chan Message) error {
	ready := make(chan bool)

	tls, err := tls.CreateTSLClientConfig()
	if err != nil {
		return err
	}

	ws := wsconnector.NewWSConnector(api.WSURL)

	dialer := websocket.Dialer{
		TLSClientConfig: tls,
	}

	ws.SetDialer(&dialer)

	var cookieStr []string

	for _, cookie := range api.Session {
		cookieStr = append(cookieStr, fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value))
	}

	extraHeaders := &http.Header{
		"Cookie": cookieStr,
	}

	ws.SetHeader(*extraHeaders)

	ws.OnConnected(func() {
		logrus.Infof("OnConnected %s", ws.WebSocket.Url)
		connectionId := 26
		connectionMessage := map[string]string{
			"locale":          api.local,
			"platformId":      "webtrading",
			"platformVersion": "chrome - 109.0.0",
			"clientId":        "app.traderepublic.com",
			"clientVersion":   "1.20.2",
		}
		connected := make(chan bool)
		go func() {
			for {
				_, message, err := ws.WebSocket.Conn.ReadMessage()
				if err != nil {
					logrus.Fatal(err)
					return
				}
				if strings.Contains(string(message), "connected") {
					logrus.Debug("Connected:", string(message))
					connected <- true
					return
				} else {
					logrus.Fatal(fmt.Errorf("something went wrong %s", string(message)))
					return
				}
			}
		}()
		msg, _ := json.Marshal(connectionMessage)
		ws.WebSocket.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("connect %d %s", connectionId, string(msg))))

		connect := <-connected
		if connect {
			logrus.Info("Connection successful established!")
			ready <- true
			return
		}
	})
	ws.OnConnectError(func(err error) {
		logrus.Info("OnConnectError: ", err.Error())
	})
	ws.OnDisconnected(func(err error) {
		logrus.Info("OnDisconnected: ", err.Error())
	})
	ws.OnClose(func(code int, text string) {
		logrus.Info("OnClose: ", code, text)
	})
	ws.OnTextMessageSent(func(message string) {
		logrus.Info("OnTextMessageSent: ", message)
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		logrus.Info("OnBinaryMessageSent: ", string(data))
	})
	ws.OnSentError(func(err error) {
		logrus.Info("OnSentError: ", err.Error())
	})
	ws.OnPingReceived(func(appData string) {
		logrus.Info("OnPingReceived: ", appData)
	})
	ws.OnPongReceived(func(appData string) {
		logrus.Info("OnPongReceived: ", appData)
	})
	ws.OnTextMessageReceived(func(message string) {
		logrus.Debug("OnTextMessageReceived: ", message)

		parseMessage := strings.SplitAfterN(message, " ", 3)
		subscriptionId, err := strconv.Atoi(strings.TrimSpace(parseMessage[0]))
		if err != nil {
			logrus.Error(err)
		}

		code := strings.TrimSpace(parseMessage[1])
		subscription, ok := api.subscribstions[subscriptionId].(map[string]interface{})
		if !ok {
			if code != "C" {
				logrus.Debugf("No active subscription for id %d, dropping message", subscriptionId)
				return
			}
		}

		if code == "A" {
			api.previousRespones[subscriptionId] = parseMessage[2]
			payload := map[string]interface{}{}
			if err := json.Unmarshal([]byte(parseMessage[2]), &payload); err != nil {
				logrus.Warnf("Received Message not unmarshaled: %+v", message)
			}
			logrus.Debugf("Send following data to Channel: %d, %s, %+v", subscriptionId, code, payload)
			dataChan <- Message{
				subscriptionId,
				subscription,
				payload,
			}
			return
		} else if code == "D" {
			response := api.calculateDelta(subscriptionId, parseMessage[2])
			logrus.Debugf("Payload is %+v", response)

			api.previousRespones[subscriptionId] = response
			payload := map[string]interface{}{}
			if err := json.Unmarshal([]byte(response), &payload); err != nil {
				logrus.Warnf("Delta Message not unmarshaled: %+v", message)
			}

			logrus.Debugf("Send following data to Channel: %d, %s, %+v", subscriptionId, code, payload)
			dataChan <- Message{
				subscriptionId,
				subscription,
				payload,
			}
			return
		}

		if code == "C" {
			delete(api.subscribstions, subscriptionId)
			delete(api.previousRespones, subscriptionId)
			return
		} else if code == "E" {
			logrus.Errorf("Received error message: %+v", message)
			api.Unsubscribe(subscriptionId)
			logrus.Fatalf("Error for %d, %+v , %+v", subscriptionId, subscription, message)
		}
	})
	ws.OnBinaryMessageReceived(func(data []byte) {
		logrus.Info("OnBinaryMessageReceived: ", string(data))
	})
	go ws.Connect()

	r := <-ready
	if r {
		api.WebSocket = ws
		return nil
	}

	return nil
}

func (api *APIClient) calculateDelta(subscriptionId int, deltaPayload string) string {
	previousResponse := api.previousRespones[subscriptionId].(string)
	i := 0
	result := []string{}
	for _, val := range strings.Split(deltaPayload, "\t") {
		sign := val[0]
		if sign == '+' {
			unescapedVal, err := url.QueryUnescape(val)
			if err != nil {
				logrus.Error(err)
			}
			result = append(result, strings.TrimSpace(unescapedVal))
		} else if sign == '-' || sign == '=' {
			if sign == '=' {
				runes := utf8.RuneCountInString(val[1:])
				result = append(result, previousResponse[i:i+runes])

			}
			i = +utf8.RuneCountInString(val[1:])

		}

	}
	return strings.Join(result, "")
}

func (api *APIClient) NextSubscriptionId() int {
	subscriptionId := api.subscriptionIdCounter
	api.subscriptionIdCounter += 1
	return subscriptionId
}

func (api *APIClient) Subscribe(payload map[string]interface{}) (int, error) {
	subscriptionId := api.NextSubscriptionId()
	api.subscribstions[subscriptionId] = payload
	msg, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	api.WebSocket.SendTextMessage(fmt.Sprintf("sub %d %s", subscriptionId, string(msg)))

	return subscriptionId, nil
}

func (api *APIClient) Unsubscribe(subscriptionId int) {
	api.WebSocket.SendTextMessage(fmt.Sprintf("unsub %d", subscriptionId))
	delete(api.subscribstions, subscriptionId)
}

func (api *APIClient) CashAvailableforOrder() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "availableCash"})
}

func (api *APIClient) CashAvailableforPayout() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "availableCashForPayout"})
}

func (api *APIClient) PortfolioStatus() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "portfolioStatus"})
}
func (api *APIClient) PortfolioHistory(timeframe string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "portfolioAggregateHistory", "range": timeframe})
}

func (api *APIClient) InstrumentDetails(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "instrument", "id": isin})
}

func (api *APIClient) InstrumentSuitability(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "instrumentSuitability", "instrumentId": isin})
}

func (api *APIClient) Portfolio() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "portfolio"})
}

func (api *APIClient) Timeline(after string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "timeline", "after": after})
}

func (api *APIClient) Watchlist() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "watchlist"})
}

func (api *APIClient) AllTimeline() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "timeline"})
}

func (api *APIClient) TimelineDetail(timelineId string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "timelineDetail", "id": timelineId})
}

func (api *APIClient) TimelineDetailOrder(order_id int) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "timelineDetail", "orderId": order_id})
}

func (api *APIClient) TimelineDetailSavingsPlan(savings_plan_id int) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "timelineDetail", "savingsPlanId": savings_plan_id})
}

func (api *APIClient) StockDetails(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "stockDetails", "id": isin})
}

func (api *APIClient) AddWatchList(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "addToWatchlist", "instrumentId": isin})
}

func (api *APIClient) RemoveWatchList(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "removeFromWatchlist", "instrumentId": isin})
}

func (api *APIClient) Ticker(isin string, exchange string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "ticker", "id": fmt.Sprintf("%s.%s", isin, exchange)})
}

func (api *APIClient) Performance(isin string, exchange string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "performance", "id": fmt.Sprintf("%s.%s", isin, exchange)})
}

func (api *APIClient) PerformanceHistory(isin string, timeframe string, exchange string, resolution string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "aggregateHistory", "id": fmt.Sprintf("%s.%s", isin, exchange), "range": timeframe, "resolution": resolution})
}

func (api *APIClient) Experience() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "experience"})
}

func (api *APIClient) MessageOfTheDay() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "messageOfTheDay"})
}

func (api *APIClient) NeonCards() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "neonCards"})
}

func (api *APIClient) NeonSearchTags() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "neonSearchTags"})
}

func (api *APIClient) NeonSearchSuggestedTags(query string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "neonSearchSuggestedTags", "data": map[string]interface{}{"q": query}})
}

func (api *APIClient) NeonSearch(query string, tags []string, page int, pageSize int) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "neonSearch", "data": map[string]interface{}{"q": query, "tags": tags, "page": page, "pageSize": pageSize}})
}

func (api *APIClient) SearchDerivative(underlying_isin string, product_type string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "derivatives", "underlying": underlying_isin, "productCategory": product_type})
}

func (api *APIClient) OrderOverview() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "orders"})
}

func (api *APIClient) OrderPrice(isin string, exchange string, order_type string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "priceForOrder", "parameters": map[string]interface{}{"exchangeId": exchange, "instrumentId": isin, "type": order_type}})
}

func (api *APIClient) OrderSize(isin string, exchange string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "availableSize", "parameters": map[string]interface{}{"exchangeId": exchange, "instrumentId": isin}})
}

func (api *APIClient) LimitedOrder(isin string, exchange string, order_type string, size int, limit float64, expiry string, expiry_date string, warnings_shown []string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "simpleCreateOrder", "clientProcessId": uuid.New().String(), "warningsShown": warnings_shown, "parameters": map[string]interface{}{"instrumentId": isin, "exchangeId": exchange, "expiry": map[string]interface{}{"type": expiry}, "limit": limit, "mode": "limit", "size": size, "type": order_type}})
}

func (api *APIClient) MarketOrder(isin string, exchange string, order_type string, size int, expiry string, expiry_date string, warnings_shown []string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "simpleCreateOrder", "clientProcessId": uuid.New().String(), "warningsShown": warnings_shown, "parameters": map[string]interface{}{"instrumentId": isin, "exchangeId": exchange, "expiry": map[string]interface{}{"type": expiry}, "mode": "market", "size": size, "type": order_type}})
}

func (api *APIClient) StopMarketOrder(isin string, exchange string, order_type string, size int, limit float64, stop float64, expiry string, expiry_date string, warnings_shown []string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "simpleCreateOrder", "clientProcessId": uuid.New().String(), "warningsShown": warnings_shown, "parameters": map[string]interface{}{"instrumentId": isin, "exchangeId": exchange, "expiry": map[string]interface{}{"type": expiry}, "limit": limit, "mode": "stopLimit", "size": size, "stop": stop, "type": order_type}})
}

func (api *APIClient) CancelOrder(order_id string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "cancelOrder", "orderId": order_id})
}

func (api *APIClient) SavingsPlanOverview() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "savingsPlans"})
}

func (api *APIClient) SavingsPlanParameters(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "savingsPlanParameters", "instrumentId": isin})
}

func (api *APIClient) CreateSavingsPlan(isin string, amount float64, interval string, start_date string, start_date_type string, start_date_value string, warnings_shown []string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "createSavingsPlan", "warningsShown": warnings_shown, "parameters": map[string]interface{}{"amount": amount, "instrumentId": isin, "interval": interval, "startDate": map[string]interface{}{"nextExecutionDate": start_date, "type": start_date_type, "value": start_date_value}}})
}

func (api *APIClient) ChangeSavingsPlan(savings_plan_id string, isin string, amount float64, interval string, start_date string, start_date_type string, start_date_value string, warnings_shown []string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "changeSavingsPlan", "id": savings_plan_id, "warningsShown": warnings_shown, "parameters": map[string]interface{}{"amount": amount, "instrumentId": isin, "interval": interval, "startDate": map[string]interface{}{"nextExecutionDate": start_date, "type": start_date_type, "value": start_date_value}}})
}

func (api *APIClient) CancelSavingsPlan(savings_plan_id string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "cancelSavingsPlan", "savingsPlanId": savings_plan_id})
}

func (api *APIClient) PriceAlarmOverview() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "priceAlarms"})
}

func (api *APIClient) CreatePriceAlarm(isin string, price float64) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "createPriceAlarm", "instrumentId": isin, "targetPrice": price})
}

func (api *APIClient) CancelPriceAlarm(price_alarm_id string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "cancelPriceAlarm", "id": price_alarm_id})
}

func (api *APIClient) News(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "neonNews", "isin": isin})
}

func (api *APIClient) NewsSubscriptions() (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "newsSubscriptions"})
}

func (api *APIClient) SubscribeNews(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "subscribeNews", "instrumentId": isin})
}

func (api *APIClient) UnsubscribeNews(isin string) (int, error) {
	return api.Subscribe(map[string]interface{}{"type": "unsubscribeNews", "instrumentId": isin})
}
