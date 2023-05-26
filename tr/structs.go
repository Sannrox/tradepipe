package tr

type RawTimeline struct {
	Data    []TimeLineEvent `json:"data"`
	Cursors struct {
		Before string `json:"before,omitempty"`
		After  string `json:"after,omitempty"`
	} `json:"cursors"`
}

type TimeLineEvent struct {
	Type string `json:"type"`
	Data struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
		Icon      string `json:"icon"`
		Title     string `json:"title"`
		Body      string `json:"body"`
		Action    struct {
			Type    string      `json:"type,omitempty"`
			Payload interface{} `json:"payload,omitempty"`
		} `json:"action,omitempty"`
		ActionLabel      string        `json:"actionLabel,omitempty"`
		Attributes       []interface{} `json:"attributes"`
		Month            string        `json:"month"`
		CashChangeAmount float64       `json:"cashChangeAmount,omitempty"`
	} `json:"data"`
}
type Action struct {
	Type    string      `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type TimelineDetail struct {
	ID           string `json:"id"`
	TitleText    string `json:"titleText"`
	SubtitleText string `json:"subtitleText"`
	Sections     []struct {
		Data      []interface{}
		Type      string `json:"type"`
		Title     string `json:"title"`
		Documents []Doc  `json:"documents"`
	} `json:"sections"`
}

type Doc struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Action struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	} `json:"action"`
	ID          string `json:"id"`
	PostboxType string `json:"postboxType"`
}

type RawPortfolio struct {
	NetValue                     float64    `json:"netValue"`
	Positions                    []Position `json:"positions"`
	ReferenceChangeProfit        int        `json:"referenceChangeProfit"`
	ReferenceChangeProfitPercent int        `json:"referenceChangeProfitPercent"`
	UnrealisedCost               float64    `json:"unrealisedCost"`
	UnrealisedProfit             float64    `json:"unrealisedProfit"`
	UnrealisedProfitPercent      float64    `json:"unrealisedProfitPercent"`
}

type Position struct {
	InstrumentID          string  `json:"instrumentId"`
	NetSize               float64 `json:"netSize"`
	NetValue              float64 `json:"netValue"`
	RealisedProfit        int     `json:"realisedProfit"`
	UnrealisedAverageCost float64 `json:"unrealisedAverageCost"`
}

type SavingsPlan struct {
	ID           string  `json:"id"`
	CreatedAt    int64   `json:"createdAt"`
	InstrumentID string  `json:"instrumentId"`
	Amount       float64 `json:"amount"`
	Interval     string  `json:"interval"`
	StartDate    struct {
		Type              string `json:"type"`
		Value             int    `json:"value"`
		NextExecutionDate string `json:"nextExecutionDate"`
	} `json:"startDate"`
	FirstExecutionDate           interface{} `json:"firstExecutionDate"`
	NextExecutionDate            string      `json:"nextExecutionDate"`
	PreviousExecutionDate        string      `json:"previousExecutionDate"`
	VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
	FinalExecutionDate           interface{} `json:"finalExecutionDate"`
	PaymentMethodID              interface{} `json:"paymentMethodId"`
	PaymentMethodCode            interface{} `json:"paymentMethodCode"`
	LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
	Paused                       bool        `json:"paused"`
}

type RawSavingsPlans struct {
	SavingsPlans []struct {
		ID           string  `json:"id"`
		CreatedAt    int64   `json:"createdAt"`
		InstrumentID string  `json:"instrumentId"`
		Amount       float64 `json:"amount"`
		Interval     string  `json:"interval"`
		StartDate    struct {
			Type              string `json:"type"`
			Value             int    `json:"value"`
			NextExecutionDate string `json:"nextExecutionDate"`
		} `json:"startDate"`
		FirstExecutionDate           interface{} `json:"firstExecutionDate"`
		NextExecutionDate            string      `json:"nextExecutionDate"`
		PreviousExecutionDate        string      `json:"previousExecutionDate"`
		VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
		FinalExecutionDate           interface{} `json:"finalExecutionDate"`
		PaymentMethodID              interface{} `json:"paymentMethodId"`
		PaymentMethodCode            interface{} `json:"paymentMethodCode"`
		LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
		Paused                       bool        `json:"paused"`
	} `json:"savingsPlans"`
}
