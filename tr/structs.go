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
		ActionLabel string        `json:"actionLabel,omitempty"`
		Attributes  []interface{} `json:"attributes"`
		Month       string        `json:"month"`
	} `json:"data"`
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
