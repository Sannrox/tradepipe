package tr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type PortfolioLoader struct {
	Client *APIClient
	Portfolio
}
type Portfolio struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	InstrumentID string  `json:"instrumentId"`
	NetSize      float64 `json:"netSize"`
	AverageBuyIn float64 `json:"averageBuyIn"`
}

func NewPortfolioLoader(client *APIClient) *PortfolioLoader {
	return &PortfolioLoader{
		Client: client,
		Portfolio: Portfolio{
			Positions: []Position{},
		},
	}
}

func (p *Portfolio) GetPositions() []Position {
	return p.Positions
}

func (p *Portfolio) GetPositionsAsBytes() ([]byte, error) {
	return json.Marshal(p.Positions)
}

func (p *PortfolioLoader) LoadPortfolio(ctx context.Context, data chan Message) error {
	_, err := p.Client.CompactPortfolio()
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-data:
			if msg.Subscription["type"] == "compactPortfolio" {
				var portfolio Portfolio
				logrus.Info(msg.Payload)
				b, err := json.Marshal(msg.Payload)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b, &portfolio)
				if err != nil {
					return fmt.Errorf("%w | %s", err, string(b))
				}
				p.Portfolio = portfolio
				return nil
			}
		}
	}
}
