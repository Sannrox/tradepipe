package tr_storage

import (
	"strings"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/sirupsen/logrus"
)

type SavingsPlans struct {
	scylla.Scylla
	keyspace string
}

func NewSavingsPlanKeyspace(contactPoint string, keyspace string) *SavingsPlans {
	err := scylla.CreateKeyspace(contactPoint, keyspace)
	if err != nil {
		panic(err)
	}
	s, err := scylla.NewScyllaDbWithPool(contactPoint, keyspace, 10)
	if err != nil {
		panic(err)
	}

	return &SavingsPlans{
		Scylla:   *s,
		keyspace: keyspace,
	}
}

func (s *SavingsPlans) CreateTableName(tableName string) string {
	return s.keyspace + "." + "user" + strings.ReplaceAll(tableName, "-", "_")
}

func (s *SavingsPlans) CreateNewSavingsPlanTable(tableName string) error {
	tableName = s.CreateTableName(tableName)
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
	if !s.TableExists(s.keyspace, tableName) {
		if err := s.CreateTable(tableName, schema); err != nil {
			return err
		}
	}

	return nil
}
func (s *SavingsPlans) GetAllSavingsPlans(tableName string) ([]*tr.SavingsPlan, error) {
	tableName = s.CreateTableName(tableName)

	iter := s.Session.Query("SELECT * FROM " + tableName).Iter()

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

// func (s *SavingsPlans) GetAllSavingsPlans(tableName string) ([]*tr.SavingsPlan, error) {
// 	tableName = s.CreateTableName(tableName)

// 	savingsPlans := []*tr.SavingsPlan{}

// 	iter := s.Session.Query("SELECT * FROM " + tableName).Iter()
// 	m := make(map[string]interface{})

// 	for iter.MapScan(m) {
// 		logrus.Info(m["startdate"])
// 		gocql.Unmarshal()
// 		savingsPlans = append(savingsPlans, &tr.SavingsPlan{
// 			ID:           m["id"].(string),
// 			CreatedAt:    int64(m["createdat"].(int)),
// 			InstrumentID: m["instrumentid"].(string),
// 			Amount:       m["amount"].(float64),
// 			Interval:     m["interval"].(string),
// 			StartDate: struct {
// 				Type              string "json:\"type\""
// 				Value             int    "json:\"value\""
// 				NextExecutionDate string "json:\"nextExecutionDate\""
// 			}{
// 				Type:              m["startdate"].([]interface{})[0].(string),
// 				Value:             m["startdate"].([]interface{})[1].(int),
// 				NextExecutionDate: m["startdate"].([]interface{})[2].(string),
// 			},
// 			FirstExecutionDate:           m["firstexecutiondate"].(string),
// 			NextExecutionDate:            m["nextexecutiondate"].(string),
// 			PreviousExecutionDate:        m["previousexecutiondate"].(string),
// 			VirtualPreviousExecutionDate: m["virtualpreviousexecutiondate"].(string),
// 			FinalExecutionDate:           m["finalexecutiondate"].(string),
// 			PaymentMethodID:              m["paymentmethodid"].(string),
// 			PaymentMethodCode:            m["paymentmethodcode"].(string),
// 			LastPaymentExecutionDate:     m["lastpaymentexecutiondate"].(string),
// 			Paused:                       m["paused"].(bool),
// 		})
// 	}
// 	return savingsPlans, nil
// }

func (s *SavingsPlans) GetSavingsPlanByID(tableName string, id string) (*tr.SavingsPlan, error) {
	tableName = s.CreateTableName(tableName)
	var savingsPlan tr.SavingsPlan

	err := s.Session.Query("SELECT * FROM "+tableName+" WHERE id = ?", id).Scan(
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

func (s *SavingsPlans) AddSavingsPlans(tableName string, savingsplans *[]tr.SavingsPlan) error {
	for _, savingsplan := range *savingsplans {
		if err := s.AddSavingsPlan(tableName, &savingsplan); err != nil {
			return err
		}
	}
	return nil
}

func (s *SavingsPlans) AddSavingsPlan(tableName string, savingsplan *tr.SavingsPlan) error {
	tableName = s.CreateTableName(tableName)
	logrus.Info("Saving savingsplan: ", savingsplan)
	return s.Insert(tableName, savingsplan)
}

func (s *SavingsPlans) UpdateSavingsPlan(tableName string, savingsplan *tr.SavingsPlan) error {
	tableName = s.CreateTableName(tableName)
	return s.Session.Query("UPDATE "+tableName+" SET "+"id = ?, created_at = ?, instrument_id = ?, amount = ?, interval = ?, start_date= ?, first_execution_date = ?, next_execution_date = ?, previous_execution_date = ?, virtual_previous_execution_date = ?, final_execution_date = ?, payment_method_id = ?, payment_method_code = ?, last_payment_execution_date = ?, paused = ? WHERE id = ?",
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
