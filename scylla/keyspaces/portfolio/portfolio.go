package portfolio

import (
	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/scylladb/gocqlx/v2/table"
)

type Portfolios struct {
	scylla.Scylla
}

var (
	schema string = "instrument_id text, " +
		"net_size double, " +
		"average_buy_in double, " +
		"PRIMARY KEY (instrument_id)"
	portfolioColumns = []string{"instrument_id", "net_size", "average_buy_in"}
	portfolioPartKey = []string{"instrument_id"}
)

func NewPortfolioKeyspace(contactPoint string, port int) (*Portfolios, error) {
	var keyspace = "portfolio"
	s, err := scylla.NewScyllaKeySpaceConnection(contactPoint, port, keyspace)
	if err != nil {
		return nil, err
	}

	return &Portfolios{
		Scylla: *s,
	}, nil
}

func (p *Portfolios) CreateTable(tableName string) (*table.Table, error) {
	tablePath := p.CreateTablePath(tableName, "user")
	if err := p.Scylla.CreateTable(tablePath, schema); err != nil {
		return nil, err
	}
	tableMeta := table.Metadata{
		Name:    tablePath,
		Columns: portfolioColumns,
		PartKey: portfolioPartKey,
	}
	return p.Scylla.NewTable(tableMeta), nil
}
func (p *Portfolios) GetAllPositions(currentTable *table.Table) (*[]tr.Position, error) {
	var positions []tr.Position
	q := p.Scylla.GetAll(currentTable)
	if err := q.SelectRelease(&positions); err != nil {
		return nil, err
	}
	return &positions, nil
}

func (p *Portfolios) AddPosition(currentTable *table.Table, position tr.Position) error {
	return p.Scylla.Insert(currentTable, &position)
}

func (p *Portfolios) UpdatePosition(currentTable *table.Table, position tr.Position) error {
	return p.Scylla.Update(currentTable, &position)
}

func (p *Portfolios) DeletePosition(currentTable *table.Table, position tr.Position) error {
	return p.Scylla.Delete(currentTable, &position)
}

func (p *Portfolios) GetPosition(currentTable *table.Table, position tr.Position) (*tr.Position, error) {
	var returnPosition *tr.Position
	q := p.Scylla.GetByKeys(currentTable, position)
	if err := q.GetRelease(&position); err != nil {
		return nil, err
	}
	return returnPosition, nil
}

func (p *Portfolios) CheckIfPositionExists(currentTable *table.Table, position tr.Position) bool {
	_, err := p.GetPosition(currentTable, position)
	return err == nil
}
