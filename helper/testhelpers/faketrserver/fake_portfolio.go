package testing

import (
	"encoding/json"
	"math/rand"

	"github.com/Sannrox/tradepipe/helper/testhelpers/random"
)

type FakePortfolio struct {
	Positions []Position `json:"positions"`
}

type Position struct {
	InstrumentID string  `json:"instrumentid"`
	NetSize      float64 `json:"netSize"`
	AverageBuyIn float64 `json:"averageBuyIn"`
}

func NewFakePortfolio() *FakePortfolio {
	return &FakePortfolio{}
}

func (p *FakePortfolio) GetPositions() []Position {
	return p.Positions
}

func (p *FakePortfolio) GetPositionsAsBytes() ([]byte, error) {
	return json.Marshal(p.Positions)
}

func (p *FakePortfolio) SetPositions(positions []Position) {
	p.Positions = positions
}

func (p *FakePortfolio) GenerateFakePortfolio() {
	p.SetPositions(p.GenerateFakePostions(10))
}

func (p *FakePortfolio) GenerateFakePostions(sets int) []Position {
	var positions []Position
	for i := 0; i < sets; i++ {
		positions = append(positions, Position{
			InstrumentID: random.GenerateRandomeISIN(),
			NetSize:      rand.Float64(),
			AverageBuyIn: rand.Float64(),
		})
	}

	return positions
}

func (p *FakePortfolio) GetPortfolio() FakePortfolio {
	return (*p)
}
