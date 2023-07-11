package timeline

import (
	"encoding/json"

	"github.com/Sannrox/tradepipe/scylla"
	"github.com/Sannrox/tradepipe/tr"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"github.com/sirupsen/logrus"
)

var (
	schemaTimeline string = `
    type text,
    id text,
    timestamp bigint,
    icon text,
    title text,
    body text,
    action_type text,
    action_payload blob,
    action_label text,
    attributes text,
    month text,
    cash_change_amount double,
    PRIMARY KEY (type, id)
	`
	timelineColumns = []string{"type", "id", "timestamp", "icon", "title", "body", "action_type", "action_payload", "action_label", "attributes", "month", "cash_change_amount"}
	timelinePartKey = []string{"type", "id"}

	schemaTimelineDetail  string = "id text," + "title_text text," + "subtitle_text text," + "sections list<frozen<tuple<text, text>>>, " + "PRIMARY KEY (id)"
	timelineDetailColumns        = []string{"id", "title_text", "subtitle_text", "sections"}
	timelineDetailPartKey        = []string{"id"}
)

type Timeline struct {
	scylla.Scylla
}

func NewTimelineKeyspace(contactPoint string, port int) (*Timeline, error) {
	var keyspace = "timeline"

	s, err := scylla.NewScyllaKeySpaceConnection(contactPoint, port, keyspace)
	if err != nil {
		return nil, err
	}

	return &Timeline{
		Scylla: *s,
	}, nil
}

func (t *Timeline) CreateTable(tableName string) (*table.Table, error) {
	tablePath := t.CreateTablePath(tableName, "user")
	if err := t.Scylla.CreateTable(tablePath, schemaTimeline); err != nil {
		return nil, err
	}
	tableMeta := table.Metadata{
		Name:    tablePath,
		Columns: timelineColumns,
		PartKey: timelinePartKey,
	}
	return t.Scylla.NewTable(tableMeta), nil
}

func (t *Timeline) CreateDetailTable(tableName string) (*table.Table, error) {
	tablePath := t.CreateTablePath(tableName, "details_user")
	if err := t.Scylla.CreateTable(tablePath, schemaTimelineDetail); err != nil {
		return nil, err
	}

	tableMeta := table.Metadata{
		Name:    tablePath,
		Columns: timelineDetailColumns,
		PartKey: timelineDetailPartKey,
	}
	return t.Scylla.NewTable(tableMeta), nil
}

func (t *Timeline) GetCompleteTimeline(currentTable *table.Table) (*[]tr.TimeLineEvent, error) {
	var flatTimelineEvents []FlatTimeLineEvent
	q := t.Scylla.GetAll(currentTable)
	if err := q.SelectRelease(&flatTimelineEvents); err != nil {
		return nil, err
	}

	var timelineEvents []tr.TimeLineEvent
	for _, flatTimelineEvent := range flatTimelineEvents {
		unflattenEvent, err := UnflattenTimelineEvent(&flatTimelineEvent)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		timelineEvents = append(timelineEvents, *unflattenEvent)
	}

	return &timelineEvents, nil
}

func (t *Timeline) GetAllTimelineDetails(currentTable *table.Table) (*[]tr.TimelineDetail, error) {
	var flatTimelineDetails []FlatTimelineDetail
	q := t.Scylla.GetAll(currentTable)
	if err := q.SelectRelease(&flatTimelineDetails); err != nil {
		return nil, err
	}

	var timelineDetails []tr.TimelineDetail
	for _, flatTimelineDetail := range flatTimelineDetails {
		unflattenDetail, err := UnflattenTimelineDetail(&flatTimelineDetail)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		timelineDetails = append(timelineDetails, *unflattenDetail)
	}
	return &timelineDetails, nil
}

func (t *Timeline) AddTimelineEvent(currentTable *table.Table, te tr.TimeLineEvent) error {
	logrus.Info("Adding timeline event")
	flattedTimeLineEvent, err := FlattenTimelineEvent(&te)
	if err != nil {
		return err
	}
	return t.Scylla.Insert(currentTable, flattedTimeLineEvent)
}

func (t *Timeline) AddTimelineDetails(currentTable *table.Table, td tr.TimelineDetail) error {
	logrus.Info("Adding timeline details")
	flattedTimelineDetail, err := FlattenTimelineDetail(&td)
	if err != nil {
		return err
	}
	return t.Scylla.Insert(currentTable, flattedTimelineDetail)
}

func (t *Timeline) UpdateTimelineEvent(currentTable *table.Table, te tr.TimeLineEvent) error {
	return t.Scylla.Update(currentTable, &te)
}

func (t *Timeline) UpdateTimelineDetails(currentTable *table.Table, td tr.TimelineDetail) error {
	return t.Scylla.Update(currentTable, &td)
}

func (t *Timeline) DeleteTimelineEvent(currentTable *table.Table, te tr.TimeLineEvent) error {
	return t.Scylla.Delete(currentTable, &te)
}

func (t *Timeline) DeleteTimelineDetails(currentTable *table.Table, td tr.TimelineDetail) error {
	return t.Scylla.Delete(currentTable, &td)
}

func (t *Timeline) GetTimelineEvent(currentTable *table.Table, te tr.TimeLineEvent) (*tr.TimeLineEvent, error) {
	flatTimelineEvent, err := FlattenTimelineEvent(&te)
	var returnFlattedEvent FlatTimeLineEvent
	if err != nil {
		return nil, err
	}
	q := t.Scylla.GetByKeys(currentTable, flatTimelineEvent)
	if err := q.GetRelease(&returnFlattedEvent); err != nil {
		return nil, err
	}

	returnEvent, err := UnflattenTimelineEvent(&returnFlattedEvent)
	if err != nil {
		return nil, err
	}

	return returnEvent, nil
}

func (t *Timeline) CheckIfTimelineEventExists(currentTable *table.Table, te tr.TimeLineEvent) bool {
	_, err := t.GetTimelineEvent(currentTable, te)
	return err != gocql.ErrNotFound
}

func (t *Timeline) GetTimelineDetail(currentTable *table.Table, td tr.TimelineDetail) (*tr.TimelineDetail, error) {
	flatTimelineDetail, err := FlattenTimelineDetail(&td)
	var returnFlattedDetail FlatTimelineDetail
	if err != nil {
		return nil, err
	}

	q := t.Scylla.GetByKeys(currentTable, flatTimelineDetail)
	if err := q.GetRelease(&returnFlattedDetail); err != nil {
		return nil, err
	}

	returnDetail, err := UnflattenTimelineDetail(&returnFlattedDetail)
	if err != nil {
		return nil, err
	}

	return returnDetail, nil
}

func (t *Timeline) CheckIfTimelineDetailExists(currentTable *table.Table, td tr.TimelineDetail) bool {
	_, err := t.GetTimelineDetail(currentTable, td)
	return err != gocql.ErrNotFound
}

type FlatTimeLineEvent struct {
	Type             string  `json:"type"`
	ID               string  `json:"id"`
	Timestamp        int64   `json:"timestamp"`
	Icon             string  `json:"icon"`
	Title            string  `json:"title"`
	Body             string  `json:"body"`
	ActionType       string  `json:"action_type,omitempty"`
	ActionPayload    string  `json:"action_payload,omitempty"`
	ActionLabel      string  `json:"action_label,omitempty"`
	Attributes       string  `json:"attributes"`
	Month            string  `json:"month"`
	CashChangeAmount float64 `json:"cashChangeAmount,omitempty"`
}

func FlattenTimelineEvent(tle *tr.TimeLineEvent) (*FlatTimeLineEvent, error) {
	actionPayloadBytes, err := json.Marshal(tle.Data.Action.Payload)
	if err != nil {
		return nil, err
	}

	attributesBytes, err := json.Marshal(tle.Data.Attributes)
	if err != nil {
		return nil, err
	}

	return &FlatTimeLineEvent{
		Type:             tle.Type,
		ID:               tle.Data.ID,
		Timestamp:        tle.Data.Timestamp,
		Icon:             tle.Data.Icon,
		Title:            tle.Data.Title,
		Body:             tle.Data.Body,
		ActionType:       tle.Data.Action.Type,
		ActionPayload:    string(actionPayloadBytes),
		ActionLabel:      tle.Data.ActionLabel,
		Attributes:       string(attributesBytes),
		Month:            tle.Data.Month,
		CashChangeAmount: tle.Data.CashChangeAmount,
	}, nil
}

func UnflattenTimelineEvent(ftle *FlatTimeLineEvent) (*tr.TimeLineEvent, error) {
	var actionPayload interface{}
	err := json.Unmarshal([]byte(ftle.ActionPayload), &actionPayload)
	if err != nil {
		return nil, err
	}

	var attributes []interface{}
	err = json.Unmarshal([]byte(ftle.Attributes), &attributes)
	if err != nil {
		return nil, err
	}

	var timelineEvent tr.TimeLineEvent
	timelineEvent.Type = ftle.Type
	timelineEvent.Data.ID = ftle.ID
	timelineEvent.Data.Timestamp = ftle.Timestamp
	timelineEvent.Data.Icon = ftle.Icon
	timelineEvent.Data.Title = ftle.Title
	timelineEvent.Data.Body = ftle.Body
	timelineEvent.Data.Action.Type = ftle.ActionType
	timelineEvent.Data.Action.Payload = actionPayload
	timelineEvent.Data.ActionLabel = ftle.ActionLabel
	timelineEvent.Data.Attributes = attributes
	timelineEvent.Data.Month = ftle.Month
	timelineEvent.Data.CashChangeAmount = ftle.CashChangeAmount
	return &timelineEvent, nil

}

func FlattenTimelineDetail(detail *tr.TimelineDetail) (*FlatTimelineDetail, error) {
	sectionsBytes, err := json.Marshal(detail.Sections)
	if err != nil {
		return nil, err
	}

	return &FlatTimelineDetail{
		ID:           detail.ID,
		TitleText:    detail.TitleText,
		SubtitleText: detail.SubtitleText,
		Sections:     string(sectionsBytes),
	}, nil
}

type FlatTimelineDetail struct {
	ID           string `json:"id"`
	TitleText    string `json:"title_text"`
	SubtitleText string `json:"subtitle_text"`
	Sections     string `json:"sections"`
}

func UnflattenTimelineDetail(flatDetail *FlatTimelineDetail) (*tr.TimelineDetail, error) {
	var sections []struct {
		Data      []interface{}
		Type      string   `json:"type"`
		Title     string   `json:"title"`
		Documents []tr.Doc `json:"documents"`
	}
	err := json.Unmarshal([]byte(flatDetail.Sections), &sections)
	if err != nil {
		return nil, err
	}

	return &tr.TimelineDetail{
		ID:           flatDetail.ID,
		TitleText:    flatDetail.TitleText,
		SubtitleText: flatDetail.SubtitleText,
		Sections:     sections,
	}, nil
}
