package keyspaces

import (
	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
)

type Portfolios struct {
	scylla.Scylla
}

func NewPortfolioKeyspace(contactPoint string) (*Portfolios, error) {
	var keyspace = "portfolio"
	s, err := scylla.NewScyllaKeySpaceConnection(contactPoint, keyspace)
	if err != nil {
		return nil, err
	}

	return &Portfolios{
		Scylla: *s,
	}, nil
}

func (p *Portfolios) CreateNewPortfolioTable(tableName string) error {
	tableName = p.CreateTablePath(tableName, "user")
	schema := "instrumentId text, " +
		"netSize double, " +
		"netValue double, " +
		"realisedProfit int, " +
		"unrealisedAverageCost double, " +
		"PRIMARY KEY (instrumentId)"

	if !p.TableExists(p.Keyspace, tableName) {
		if err := p.CreateTable(tableName, schema); err != nil {
			return err
		}
	}

	return nil
}

func (p *Portfolios) GetAllPositions(tableName string) ([]*tr.Position, error) {
	tablePath := p.CreateTablePath(tableName, "user")
	positions := []*tr.Position{}

	iter := p.Session.Query("SELECT * FROM " + tablePath).Iter()
	m := make(map[string]interface{})
	for iter.MapScan(m) {
		logrus.Debug(m)
		positions = append(positions, &tr.Position{
			InstrumentID:          m["instrumentid"].(string),
			NetSize:               m["netsize"].(float64),
			NetValue:              m["netvalue"].(float64),
			RealisedProfit:        m["realisedprofit"].(int),
			UnrealisedAverageCost: m["unrealisedaveragecost"].(float64),
		})
		m = make(map[string]interface{})
	}

	return positions, nil
}

func (p *Portfolios) AddPositions(tableName string, positions *[]tr.Position) error {
	for _, position := range *positions {
		err := p.AddPosition(tableName, &position)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Portfolios) AddPosition(tableName string, position *tr.Position) error {
	tablePath := p.CreateTablePath(tableName, "user")
	return p.Insert(tablePath, position)
}

func (p *Portfolios) UpdatePosition(tableName string, position *tr.Position) error {
	tablePath := p.CreateTablePath(tableName, "user")
	return p.Session.Query("UPDATE "+tablePath+" SET net_size = ?, net_value = ?, realised_profit = ?, unrealised_average_cost = ? WHERE instrument_id = ?",
		position.NetSize,
		position.NetValue,
		position.RealisedProfit,
		position.UnrealisedAverageCost,
		position.InstrumentID,
	).Exec()
}

func (p *Portfolios) DeletePosition(tableName string, position *tr.Position) error {
	tablePath := p.CreateTablePath(tableName, "user")
	return p.Session.Query("DELETE FROM "+tablePath+" WHERE instrument_id = ?", position.InstrumentID).Exec()
}

func (p *Portfolios) GetPosition(tableName string, instrumentID string) (*tr.Position, error) {
	tablePath := p.CreateTablePath(tableName, "user")
	var position tr.Position

	err := p.Session.Query("SELECT * FROM "+tablePath+" WHERE instrument_id = ?", instrumentID).Scan(
		&position.InstrumentID,
		&position.NetSize,
		&position.NetValue,
		&position.RealisedProfit,
		&position.UnrealisedAverageCost,
	)

	return &position, err
}
