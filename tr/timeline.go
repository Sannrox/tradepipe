package tr

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type TimeLine struct {
	Client                    *APIClient
	ReceivedDetail            int
	RequestedDetail           int
	NumberofTimlineDetails    int
	NumberofTimline           int
	SinceTimestamp            int64
	TimeLineEvents            []TimeLineEvent
	TimeLineEventsWithDocs    []TimeLineEvent
	TimelineEventsWithoutDocs []TimeLineEvent
	TimelineDetails           []TimelineDetail
}

func NewTimeLine(client *APIClient) *TimeLine {
	return &TimeLine{
		Client:                    client,
		ReceivedDetail:            0,
		RequestedDetail:           0,
		NumberofTimlineDetails:    0,
		NumberofTimline:           0,
		SinceTimestamp:            0,
		TimeLineEvents:            []TimeLineEvent{},
		TimeLineEventsWithDocs:    []TimeLineEvent{},
		TimelineEventsWithoutDocs: []TimeLineEvent{},
		TimelineDetails:           []TimelineDetail{},
	}
}

func (t *TimeLine) SumTotalOfTimeLineEventsWDocsAndWithoutDocs() int {
	return len(t.TimeLineEventsWithDocs) + len(t.TimelineEventsWithoutDocs)
}

func (t *TimeLine) SetSinceTimestamp(sinceTimestamp int64) {
	t.SinceTimestamp = sinceTimestamp
}

func (t *TimeLine) GetTimeLineEventsWithDocs() []TimeLineEvent {
	return t.TimeLineEventsWithDocs
}

func (t *TimeLine) GetTimeLineEventsWithoutDocs() []TimeLineEvent {
	return t.TimelineEventsWithoutDocs
}

func (t *TimeLine) GetTimeLineEvents() []TimeLineEvent {
	return t.TimeLineEvents
}

func (t *TimeLine) GetTimeLineEventsAsBytes() ([]byte, error) {
	return json.Marshal(t.TimeLineEvents)
}

func (t *TimeLine) GetTimeLineDetails() []TimelineDetail {
	return t.TimelineDetails
}

func (t *TimeLine) GetTimeLineDetailsAsBytes() ([]byte, error) {
	return json.Marshal(t.TimelineDetails)
}

func (t *TimeLine) LoadTimeLine(ctx context.Context, data chan Message) error {
	_, err := t.LoadNextTimeline(nil, t.SinceTimestamp)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-data:
			if msg.Subscription["type"] == "timeline" {
				end, err := t.LoadNextTimeline(msg.Payload, t.SinceTimestamp)
				if err != nil {
					return err
				}
				if end == -1 {
					return nil
				}
			} else {
				continue
			}
		}
	}
}

func (t *TimeLine) LoadNextTimeline(response map[string]interface{}, maxAgeTimestamp int64) (int, error) {
	if response == nil {
		logrus.Info("Awaiting #1 timeline")

		id, err := t.Client.AllTimeline()
		if err != nil {
			return -1, err
		}
		return id, nil

	} else {
		var timeline RawTimeline
		b, err := json.Marshal(response)
		if err != nil {
			return -1, err
		}
		err = json.Unmarshal(b, &timeline)
		if err != nil {
			return -1, fmt.Errorf("%w | %s", err, string(b))
		}
		timelineData := timeline.Data
		timestamp := timelineData[len(timelineData)-1].Data.Timestamp
		t.NumberofTimline = t.NumberofTimline + 1
		t.NumberofTimlineDetails = t.NumberofTimlineDetails + len(response["data"].([]interface{}))
		logrus.Infof("Received timeline #%d with %d events", t.NumberofTimline, len(response["data"].([]interface{})))
		t.TimeLineEvents = append(t.TimeLineEvents, timelineData...)

		after := timeline.Cursors.After
		if after == "" {
			logrus.Info("No more timeline")
			return -1, nil
		} else if maxAgeTimestamp != 0 && timestamp < maxAgeTimestamp {
			logrus.Info("No more timeline")
			return -1, nil
		} else {
			logrus.Info("Requesting next timeline")
			id, err := t.Client.Timeline(after)
			if err != nil {
				return -1, err
			}
			return id, nil
		}
	}
}

func (t *TimeLine) LoadTimeLineDetails(ctx context.Context, data chan Message) error {

	re, err := t.RequestTimeLineDetails(5, t.SinceTimestamp)
	if err != nil {
		return err
	}
	if re == -1 {
		return nil
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-data:
			if msg.Subscription["type"] == "timelineDetail" {
				t.ReceivedDetail++
				b, err := json.Marshal(msg.Payload)
				if err != nil {
					return err
				}
				var timelineDetail TimelineDetail
				err = json.Unmarshal(b, &timelineDetail)
				if err != nil {
					return err
				}
				t.TimelineDetails = append(t.TimelineDetails, timelineDetail)
				end, err := t.loadMoreTimeLineDetails(t.SinceTimestamp)
				if err != nil {
					return err
				}
				if end == -1 {
					return nil
				}
			} else {
				continue
			}
		}
	}
}

func (t *TimeLine) RequestTimeLineDetails(numToRequest int, maxAgeTimestamp int64) (int, error) {
	for numToRequest > 0 {
		if len(t.TimeLineEvents) == t.SumTotalOfTimeLineEventsWDocsAndWithoutDocs() {
			logrus.Info("All timeline details requested")
			return -1, nil
		} else {
			event := t.TimeLineEvents[len(t.TimeLineEvents)-t.SumTotalOfTimeLineEventsWDocsAndWithoutDocs()-1]

			if WithDocs := t.FindTimeLineEventsWithDocs(event, maxAgeTimestamp); WithDocs {

				numToRequest--
				t.RequestedDetail++
				id := event.Data.ID
				return t.Client.TimelineDetail(id)
			} else {
				continue
			}

		}
	}

	return -1, nil
}

func (t *TimeLine) FindTimeLineEventsWithDocs(event TimeLineEvent, maxAgeTimestamp int64) bool {
	if t.TimeLineEvents == nil {
		return false
	}

	action := event.Data.Action

	var msg string
	if maxAgeTimestamp != 0 && event.Data.Timestamp > maxAgeTimestamp {
		msg = "Skipping timeline detail %s %s"
	} else if reflect.ValueOf(action).IsZero() {
		if event.Data.ActionLabel == "" {
			msg = "Skipping: no action"
		}
	} else if action.Type != "timelineDetail" {
		msg += fmt.Sprintf("Skipping timeline detail %s", action.Type)
	} else if str, ok := action.Payload.(string); !ok || str != event.Data.ID {
		msg += fmt.Sprintf("Skip: payload unmatched %s", action.Payload)
	} else if event.Data.Title == "Cash In" {
		// BUG:  Shows a timelineDetail but there is not detail to get - test from time to time
		msg += "Skip: 2023 Cash In, because of unknown error"
	}
	if msg == "" {
		t.TimeLineEventsWithDocs = append(t.TimeLineEventsWithDocs, event)
		return true
	} else {
		t.TimelineEventsWithoutDocs = append(t.TimelineEventsWithoutDocs, event)
		logrus.Debugf("%s %s %s: ", msg, event.Data.Title, event.Data.Body)
		t.NumberofTimlineDetails = t.NumberofTimlineDetails - 1
		return false

	}
}

func (t *TimeLine) loadMoreTimeLineDetails(maxAgeTimestamp int64) (int, error) {
	if t.ReceivedDetail == t.RequestedDetail {
		remaining := len(t.TimeLineEvents)
		var num int
		if remaining < 5 {
			num = remaining
		} else {
			num = 5
		}
		re, err := t.RequestTimeLineDetails(num, 0)
		if err != nil {
			return -1, err
		}
		return re, nil
	}
	return -1, nil
}
