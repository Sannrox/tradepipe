package savingsplans

import (
	"fmt"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	schema string = "id text," +
		"created_at bigint," +
		"instrument_id text," +
		"amount double," +
		"interval text," +
		"start_date_type text," +
		"start_date_value int," +
		"start_date_next_execution_date text," +
		"first_execution_date blob," +
		"next_execution_date text," +
		"previous_execution_date text," +
		"virtual_previous_execution_date text," +
		"final_execution_date blob," +
		"payment_method_id blob," +
		"payment_method_code blob," +
		"last_payment_execution_date blob," +
		"paused boolean," +
		"primary key ((id), created_at)"
	savingsPlanColumns = []string{"id", "created_at", "instrument_id", "amount", "interval", "start_date_type", "start_date_value", "start_date_next_execution_date", "first_execution_date", "next_execution_date", "previous_execution_date", "virtual_previous_execution_date", "final_execution_date", "payment_method_id", "payment_method_code", "last_payment_execution_date", "paused"}
	savingsPlanPartKey = []string{"id", "created_at"}
)

type SavingsPlans struct {
	scylla.Scylla
}

func NewSavingsPlanKeyspace(contactPoint string, port int) (*SavingsPlans, error) {
	var keyspace = "savingsplans"
	s, err := scylla.NewScyllaKeySpaceConnection(contactPoint, port, keyspace)
	if err != nil {
		return nil, err
	}

	return &SavingsPlans{
		Scylla: *s,
	}, nil
}
func (s *SavingsPlans) CreateTable(tableName string) (*table.Table, error) {
	tablePath := s.CreateTablePath(tableName, "user")
	if err := s.Scylla.CreateTable(tablePath, schema); err != nil {
		return nil, err
	}
	tableMeta := table.Metadata{
		Name:    tablePath,
		Columns: savingsPlanColumns,
		PartKey: savingsPlanPartKey,
	}
	return s.Scylla.NewTable(tableMeta), nil
}

func (s *SavingsPlans) GetAllPlans(currentTable *table.Table) (*[]tr.SavingsPlan, error) {
	var flatSavingsplans []SavingsPlanFlat

	q := s.Scylla.GetAll(currentTable)
	if err := q.SelectRelease(&flatSavingsplans); err != nil {
		return nil, err
	}
	var plans []tr.SavingsPlan

	for _, flatSavingsplan := range flatSavingsplans {
		unflattedSavingsplan := UnflattenSavingsPlan(&flatSavingsplan)
		plans = append(plans, *unflattedSavingsplan)
	}

	return &plans, nil
}

func (s *SavingsPlans) AddPlan(currentTable *table.Table, plan tr.SavingsPlan) error {
	flattedSavingsplan := FlattenSavingsPlan(&plan)
	return s.Scylla.Insert(currentTable, &flattedSavingsplan)
}

func (s *SavingsPlans) UpdatePlan(currentTable *table.Table, plan tr.SavingsPlan) error {
	return s.Scylla.Update(currentTable, &plan)
}

func (s *SavingsPlans) DeletePlan(currentTable *table.Table, plan tr.SavingsPlan) error {
	return s.Scylla.Delete(currentTable, &plan)
}

func (s *SavingsPlans) GetPlan(currentTable *table.Table, plan tr.SavingsPlan) (*tr.SavingsPlan, error) {
	flatSavingsplan := FlattenSavingsPlan(&plan)
	var returnFlattedPlan SavingsPlanFlat

	q := s.Scylla.GetByKeys(currentTable, flatSavingsplan)
	if err := q.GetRelease(&returnFlattedPlan); err != nil {
		return nil, err
	}

	returnPlan := UnflattenSavingsPlan(&returnFlattedPlan)

	return returnPlan, nil
}

func (s *SavingsPlans) CheckIfPlanExists(currentTable *table.Table, plan tr.SavingsPlan) bool {
	_, err := s.GetPlan(currentTable, plan)
	return err != gocql.ErrNotFound
}

type SavingsPlanFlat struct {
	ID                           string  `json:"id"`
	CreatedAt                    int64   `json:"created_at"`
	InstrumentID                 string  `json:"instrument_id"`
	Amount                       float64 `json:"amount"`
	Interval                     string  `json:"interval"`
	StartDateType                string  `json:"start_date_type"`
	StartDateValue               int     `json:"start_date_value"`
	StartDateNextExecutionDate   string  `json:"start_date_next_execution_date"`
	FirstExecutionDate           string  `json:"first_execution_date"`
	NextExecutionDate            string  `json:"next_execution_date"`
	PreviousExecutionDate        string  `json:"previous_execution_date"`
	VirtualPreviousExecutionDate string  `json:"virtual_previous_execution_date"`
	FinalExecutionDate           string  `json:"final_execution_date"`
	PaymentMethodID              string  `json:"payment_method_id"`
	PaymentMethodCode            string  `json:"payment_method_code"`
	LastPaymentExecutionDate     string  `json:"last_payment_execution_date"`
	Paused                       bool    `json:"paused"`
}

func FlattenSavingsPlan(plan *tr.SavingsPlan) *SavingsPlanFlat {
	// Ensure nullable fields have a default value
	firstExecutionDate := ""
	if plan.FirstExecutionDate != nil {
		firstExecutionDate = fmt.Sprint(plan.FirstExecutionDate)
	}

	finalExecutionDate := ""
	if plan.FinalExecutionDate != nil {
		finalExecutionDate = fmt.Sprint(plan.FinalExecutionDate)
	}

	paymentMethodID := ""
	if plan.PaymentMethodID != nil {
		paymentMethodID = fmt.Sprint(plan.PaymentMethodID)
	}

	paymentMethodCode := ""
	if plan.PaymentMethodCode != nil {
		paymentMethodCode = fmt.Sprint(plan.PaymentMethodCode)
	}

	lastPaymentExecutionDate := ""
	if plan.LastPaymentExecutionDate != nil {
		lastPaymentExecutionDate = fmt.Sprint(plan.LastPaymentExecutionDate)
	}

	return &SavingsPlanFlat{
		ID:                           plan.ID,
		CreatedAt:                    plan.CreatedAt,
		InstrumentID:                 plan.InstrumentID,
		Amount:                       plan.Amount,
		Interval:                     plan.Interval,
		StartDateType:                plan.StartDate.Type,
		StartDateValue:               plan.StartDate.Value,
		StartDateNextExecutionDate:   plan.StartDate.NextExecutionDate,
		FirstExecutionDate:           firstExecutionDate,
		NextExecutionDate:            plan.NextExecutionDate,
		PreviousExecutionDate:        plan.PreviousExecutionDate,
		VirtualPreviousExecutionDate: plan.VirtualPreviousExecutionDate,
		FinalExecutionDate:           finalExecutionDate,
		PaymentMethodID:              paymentMethodID,
		PaymentMethodCode:            paymentMethodCode,
		LastPaymentExecutionDate:     lastPaymentExecutionDate,
		Paused:                       plan.Paused,
	}
}

func UnflattenSavingsPlan(flat *SavingsPlanFlat) *tr.SavingsPlan {
	return &tr.SavingsPlan{
		ID:           flat.ID,
		CreatedAt:    flat.CreatedAt,
		InstrumentID: flat.InstrumentID,
		Amount:       flat.Amount,
		Interval:     flat.Interval,
		StartDate: struct {
			Type              string `json:"type"`
			Value             int    `json:"value"`
			NextExecutionDate string `json:"nextExecutionDate"`
		}{
			Type:              flat.StartDateType,
			Value:             flat.StartDateValue,
			NextExecutionDate: flat.StartDateNextExecutionDate,
		},
		FirstExecutionDate:           flat.FirstExecutionDate,
		NextExecutionDate:            flat.NextExecutionDate,
		PreviousExecutionDate:        flat.PreviousExecutionDate,
		VirtualPreviousExecutionDate: flat.VirtualPreviousExecutionDate,
		FinalExecutionDate:           flat.FinalExecutionDate,
		PaymentMethodID:              flat.PaymentMethodID,
		PaymentMethodCode:            flat.PaymentMethodCode,
		LastPaymentExecutionDate:     flat.LastPaymentExecutionDate,
		Paused:                       flat.Paused,
	}
}
