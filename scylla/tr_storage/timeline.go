package tr_storage

import (
	"reflect"
	"strings"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
)

type Timeline struct {
	scylla.Scylla
}

func NewTimelineKeyspace(contactPoint string) *Timeline {
	var keyspace = "timeline"
	s, err := scylla.NewScyllaKeySpaceConnection(contactPoint, keyspace)
	if err != nil {
		panic(err)
	}

	return &Timeline{
		Scylla: *s,
	}
}

func (t *Timeline) CreateNewTable(tableName string) error {
	tablePath := t.CreateTablePath(tableName, "user")

	// logrus.Info(t.Session.Query("CREATE TYPE " + t.Keyspace + ".action (type text, payload text)").String())
	// if err := t.Session.Query("CREATE TYPE IF NOT EXISTS " + t.Keyspace + ".action (type text, payload text)").Exec(); err != nil {
	// 	return err
	// }
	// Type string `json:"type"`
	// Data struct {
	// 	ID        string `json:"id"`
	// 	Timestamp int64  `json:"timestamp"`
	// 	Icon      string `json:"icon"`
	// 	Title     string `json:"title"`
	// 	Body      string `json:"body"`
	// 	Action    struct {
	// 		Type    string      `json:"type,omitempty"`
	// 		Payload interface{} `json:"payload,omitempty"`
	// 	} `json:"action,omitempty"`
	// 	ActionLabel      string        `json:"actionLabel,omitempty"`
	// 	Attributes       []interface{} `json:"attributes"`
	// 	Month            string        `json:"month"`
	// 	CashChangeAmount float64       `json:"cashChangeAmount,omitempty"`
	// } `json:"data"`
	schema := "type text," + "id text," + "timestamp bigint," +
		"icon text," + "title text," + "body text," +
		"action frozen<tuple<text, text>>," + "actionlabel text," +
		"attributes list<text>," + "month text," +
		"cashchangeamount double," + "PRIMARY KEY (type, id)"
	if !t.TableExists(t.Keyspace, tablePath) {
		if err := t.CreateTable(tablePath, schema); err != nil {
			return err
		}
	}

	return nil
}

func (t *Timeline) All(tableName string) ([]*tr.TimeLineEvent, error) {
	tablePath := t.CreateTablePath(tableName, "user")

	iter := t.Session.Query("SELECT * FROM " + tablePath).Iter()

	timelineEvents := make([]*tr.TimeLineEvent, 0, iter.NumRows())
	var te *tr.TimeLineEvent
	// Table timeline
	// alphabetically sorted (except for id)
	// sort is specified by the database
	for iter.Scan(&te.Type, &te.Data.ID, &te.Data.Timestamp, &te.Data.Icon, &te.Data.Title, &te.Data.Body, &te.Data.Action.Type, &te.Data.Action.Payload, &te.Data.ActionLabel, &te.Data.Attributes, &te.Data.Month, &te.Data.CashChangeAmount) {
		timelineEvents = append(timelineEvents, te)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return timelineEvents, nil

}

func (t *Timeline) GetByID(tableName string, id string) (*tr.TimeLineEvent, error) {
	tablePath := t.CreateTablePath(tableName, "user")

	iter := t.Session.Query("SELECT * FROM "+tablePath+" WHERE id = ?", id).Iter()

	var te tr.TimeLineEvent
	// Table timeline
	// alphabetically sorted (except for id)
	// sort is specified by the database
	for iter.Scan(&te.Type, &te.Data.ID, &te.Data.Timestamp, &te.Data.Icon, &te.Data.Title, &te.Data.Body, &te.Data.Action.Type, &te.Data.Action.Payload, &te.Data.ActionLabel, &te.Data.Attributes, &te.Data.Month, &te.Data.CashChangeAmount) {
		return &te, nil
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *Timeline) getFields(model interface{}) []string {
	tt := reflect.TypeOf(model)

	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}

	var fields []string
	for i := 0; i < tt.NumField(); i++ {
		field := tt.Field(i)
		if field.Name == "ID" {
			fields = append(fields, "id")
		} else {
			fields = append(fields, strings.ToLower(field.Name))
		}
	}

	return fields
}
func (t *Timeline) InsertOne(tableName string, te *tr.TimeLineEvent) error {
	tablePath := t.CreateTablePath(tableName, "user")
	fields := t.getFields(te.Data)

	//type, id, timestamp, icon, title, body, action, actionlabel, attributes, month, cashchangeamount
	query := "INSERT INTO " + tablePath + "(" + "type, " + strings.Join(fields, ", ") + ") " + "VALUES (?, ?, ?, ?, ?, ?, (?, ?), ?, ?, ?, ?)"
	return t.Query(query, te.Type, te.Data.ID, te.Data.Timestamp, te.Data.Icon, te.Data.Title, te.Data.Body, te.Data.Action.Type, te.Data.Action.Payload, te.Data.ActionLabel, te.Data.Attributes, te.Data.Month, te.Data.CashChangeAmount)
}

func (t *Timeline) InsertMany(tableName string, tes *[]tr.TimeLineEvent) error {
	for _, te := range *tes {
		if err := t.InsertOne(tableName, &te); err != nil {
			return err
		}
	}
	return nil
}

func (t *Timeline) DeleteOne(tableName string, id string) error {
	tablePath := t.CreateTablePath(tableName, "user")
	return t.Delete(tablePath, id)
}
