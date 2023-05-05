package tr_storage

import (
	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
)

type SavingsPlans struct {
	scylla.Scylla
}

func NewSavingsPlanKeyspace(contactPoint string) *SavingsPlans {
	var keyspace = "savingsplans"
	s, err := scylla.NewScyllaKeySpaceConnection(contactPoint, keyspace)
	if err != nil {
		panic(err)
	}

	return &SavingsPlans{
		Scylla: *s,
	}
}

func (s *SavingsPlans) CreateNewTable(tableName string) error {
	tablePath := s.CreateTablePath(tableName, "user")
	schema :=
		"id text," +
			"createdAt bigint," +
			"instrumentId text," +
			"amount double," +
			"interval text," +
			"startDate frozen<tuple<text, int, text>>," +
			"firstExecutionDate blob," +
			"nextExecutionDate text," +
			"previousExecutionDate text," +
			"virtualPreviousExecutionDate text," +
			"finalExecutionDate blob," +
			"paymentMethodId blob," +
			"paymentMethodCode blob," +
			"lastPaymentExecutionDate blob," +
			"paused boolean," +
			"PRIMARY KEY ((id), createdAt)"
	if !s.TableExists(s.Keyspace, tablePath) {
		if err := s.CreateTable(tablePath, schema); err != nil {
			return err
		}
	}

	return nil
}
func (s *SavingsPlans) All(tableName string) ([]*tr.SavingsPlan, error) {
	tablePath := s.CreateTablePath(tableName, "user")

	iter := s.Session.Query("SELECT * FROM " + tablePath).Iter()

	savingsPlans := make([]*tr.SavingsPlan, 0, iter.NumRows())
	var sp tr.SavingsPlan
	// Table savingsplans
	// alphabetically sorted (except for id)
	// sort is specified by the database
	//  id  | createdat  | amount | finalexecutiondate | firstexecutiondate | instrumentid | interval | lastpaymentexecutiondate | nextexecutiondate | paused | paymentmethodcode | paymentmethodid | previousexecutiondate | startdate  | virtualpreviousexecutiondate
	for iter.Scan(&sp.ID, &sp.CreatedAt, &sp.Amount, &sp.FinalExecutionDate, &sp.FirstExecutionDate, &sp.InstrumentID, &sp.Interval, &sp.LastPaymentExecutionDate, &sp.NextExecutionDate, &sp.Paused, &sp.PaymentMethodCode, &sp.PaymentMethodID, &sp.PreviousExecutionDate, &sp.StartDate.Type, &sp.StartDate.Value, &sp.StartDate.NextExecutionDate, &sp.VirtualPreviousExecutionDate) {
		savingsPlans = append(savingsPlans, &sp)
		sp = tr.SavingsPlan{}
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return savingsPlans, nil
}

func (s *SavingsPlans) GetByID(tableName string, id string) (*tr.SavingsPlan, error) {
	tablePath := s.CreateTablePath(tableName, "user")
	var savingsPlan tr.SavingsPlan

	err := s.Session.Query("SELECT * FROM "+tablePath+" WHERE id = ?", id).Scan(
		&savingsPlan.ID,
		&savingsPlan.CreatedAt,
		&savingsPlan.InstrumentID,
		&savingsPlan.Amount,
		&savingsPlan.Interval,
		&savingsPlan.StartDate.Type,
		&savingsPlan.StartDate.Value,
		&savingsPlan.StartDate.NextExecutionDate,
		&savingsPlan.FirstExecutionDate,
		&savingsPlan.NextExecutionDate,
		&savingsPlan.PreviousExecutionDate,
		&savingsPlan.VirtualPreviousExecutionDate,
		&savingsPlan.FinalExecutionDate,
		&savingsPlan.PaymentMethodID,
		&savingsPlan.PaymentMethodCode,
		&savingsPlan.LastPaymentExecutionDate,
		&savingsPlan.Paused,
	)

	return &savingsPlan, err
}

func (s *SavingsPlans) InsertMany(tableName string, savingsplans *[]tr.SavingsPlan) error {
	for _, savingsplan := range *savingsplans {
		if err := s.InsertOne(tableName, &savingsplan); err != nil {
			return err
		}
	}
	return nil
}

func (s *SavingsPlans) InsertOne(tableName string, savingsplan *tr.SavingsPlan) error {
	tablePath := s.CreateTablePath(tableName, "user")
	return s.Insert(tablePath, savingsplan)
}

func (s *SavingsPlans) Update(tableName string, savingsplan *tr.SavingsPlan) error {
	tablePath := s.CreateTablePath(tableName, "user")
	return s.Session.Query("UPDATE "+tablePath+" SET "+"id = ?, created_at = ?, instrument_id = ?, amount = ?, interval = ?, start_date= ?, first_execution_date = ?, next_execution_date = ?, previous_execution_date = ?, virtual_previous_execution_date = ?, final_execution_date = ?, payment_method_id = ?, payment_method_code = ?, last_payment_execution_date = ?, paused = ? WHERE id = ?",
		savingsplan.ID,
		savingsplan.CreatedAt,
		savingsplan.InstrumentID,
		savingsplan.Amount,
		savingsplan.Interval,
		savingsplan.StartDate,
		savingsplan.FirstExecutionDate,
		savingsplan.NextExecutionDate,
		savingsplan.PreviousExecutionDate,
		savingsplan.VirtualPreviousExecutionDate,
		savingsplan.FinalExecutionDate,
		savingsplan.PaymentMethodID,
		savingsplan.PaymentMethodCode,
		savingsplan.LastPaymentExecutionDate,
		savingsplan.Paused,
	).Exec()

}
