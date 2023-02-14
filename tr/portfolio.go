package tr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Portfolio struct {
	Client   *APIClient
	Postions []Position
	NetValue float64
}

func NewPortfolio(client *APIClient) *Portfolio {
	return &Portfolio{
		Client:   client,
		Postions: []Position{},
		NetValue: 0,
	}
}

func (p *Portfolio) GetPositions() []Position {
	return p.Postions
}

func (p *Portfolio) GetNetValue() float64 {
	return p.NetValue
}

func (p *Portfolio) GetPositionsAsBytes() ([]byte, error) {
	return json.Marshal(p.Postions)
}

func (p *Portfolio) LoadPortfolio(ctx context.Context, data chan Message) error {
	_, err := p.Client.Portfolio()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-data:
			logrus.Info(msg)
			if msg.Subscription["type"] == "portfolio" {
				var portfolio RawPortfolio
				b, err := json.Marshal(msg.Payload)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b, &portfolio)
				if err != nil {
					return fmt.Errorf("%w | %s", err, string(b))
				}
				p.NetValue = portfolio.NetValue
				p.Postions = portfolio.Positions
				return nil
			}
		}
	}
}
