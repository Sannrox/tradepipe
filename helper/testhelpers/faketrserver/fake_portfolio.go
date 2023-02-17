package testing

import (
	"encoding/json"
	"math/rand"

	"github.com/Sannrox/tradepipe/helper/testhelpers/random"
)

type FakePortfolio struct {
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

func NewFakePortfolio() *FakePortfolio {
	return &FakePortfolio{
		NetValue:                     0,
		Positions:                    []Position{},
		ReferenceChangeProfit:        0,
		ReferenceChangeProfitPercent: 0,
		UnrealisedCost:               0,
		UnrealisedProfit:             0,
		UnrealisedProfitPercent:      0,
	}
}

func (p *FakePortfolio) GetPositions() []Position {
	return p.Positions
}

func (p *FakePortfolio) GetNetValue() float64 {
	return p.NetValue
}

func (p *FakePortfolio) GetPositionsAsBytes() ([]byte, error) {
	return json.Marshal(p.Positions)
}

func (p *FakePortfolio) SetNetValue(value float64) {
	p.NetValue = value
}

func (p *FakePortfolio) SetPositions(positions []Position) {
	p.Positions = positions
}

func (p *FakePortfolio) SetReferenceChangeProfit(value int) {
	p.ReferenceChangeProfit = value
}

func (p *FakePortfolio) SetReferenceChangeProfitPercent(value int) {
	p.ReferenceChangeProfitPercent = value
}

func (p *FakePortfolio) SetUnrealisedCost(value float64) {
	p.UnrealisedCost = value
}

func (p *FakePortfolio) SetUnrealisedProfit(value float64) {
	p.UnrealisedProfit = value
}

func (p *FakePortfolio) SetUnrealisedProfitPercent(value float64) {
	p.UnrealisedProfitPercent = value
}

func (p *FakePortfolio) GetReferenceChangeProfit() int {
	return p.ReferenceChangeProfit
}

func (p *FakePortfolio) GetReferenceChangeProfitPercent() int {
	return p.ReferenceChangeProfitPercent
}

func (p *FakePortfolio) GetUnrealisedCost() float64 {
	return p.UnrealisedCost
}

func (p *FakePortfolio) GetUnrealisedProfit() float64 {
	return p.UnrealisedProfit
}

func (p *FakePortfolio) GetUnrealisedProfitPercent() float64 {
	return p.UnrealisedProfitPercent
}

func (p *FakePortfolio) GenerateFakePortfolio() {
	p.SetNetValue(100000)
	p.SetReferenceChangeProfit(0)
	p.SetReferenceChangeProfitPercent(0)
	p.SetUnrealisedCost(0)
	p.SetUnrealisedProfit(0)
	p.SetUnrealisedProfitPercent(0)
	p.SetPositions(p.GenerateFakePostions(10))
}

func (p *FakePortfolio) GenerateFakePostions(sets int) []Position {
	var positions []Position
	for i := 0; i < sets; i++ {
		positions = append(positions, Position{
			InstrumentID:          random.GenerateRandomeISIN(),
			NetSize:               rand.Float64(),
			NetValue:              rand.Float64(),
			RealisedProfit:        rand.Int(),
			UnrealisedAverageCost: rand.Float64(),
		})
	}

	return positions
}

func (p *FakePortfolio) GetPortfolio() FakePortfolio {
	return (*p)
}
