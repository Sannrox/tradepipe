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
	NetValue  float64    `json:"netValue"`
}

func NewPortfolioLoader(client *APIClient) *PortfolioLoader {
	return &PortfolioLoader{
		Client: client,
		Portfolio: Portfolio{
			Positions: []Position{},
			NetValue:  0,
		},
	}
}

func (p *Portfolio) GetPositions() []Position {
	return p.Positions
}

func (p *Portfolio) GetNetValue() float64 {
	return p.NetValue
}

func (p *Portfolio) GetPositionsAsBytes() ([]byte, error) {
	return json.Marshal(p.Positions)
}

func (p *PortfolioLoader) LoadPortfolio(ctx context.Context, data chan Message) error {
	_, err := p.Client.Portfolio()
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-data:
			if msg.Subscription["type"] == "portfolio" {
				var portfolio RawPortfolio
				logrus.Info(msg.Payload)
				b, err := json.Marshal(msg.Payload)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b, &portfolio)
				if err != nil {
					return fmt.Errorf("%w | %s", err, string(b))
				}
				p.Positions = portfolio.Positions
				p.NetValue = portfolio.NetValue
				return nil
			}
		}
	}
}
