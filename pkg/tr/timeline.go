package tr

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type TimeLine struct {
	Client                    *APIClient
	ReceivedDetail            int
	RequestedDetail           int
	NumberofTimlineDetails    int
	NumberofTimline           int
	SinceTimestamp            float64
	TimeLineEvents            []interface{}
	TimeLineEventsWithDocs    []interface{}
	TimelineEventsWithoutDocs []interface{}
	TimelineDetails           []map[string]interface{}
}

func NewTimeLine(client *APIClient) *TimeLine {
	return &TimeLine{
		Client:                    client,
		ReceivedDetail:            0,
		RequestedDetail:           0,
		NumberofTimlineDetails:    0,
		NumberofTimline:           0,
		SinceTimestamp:            0,
		TimeLineEvents:            []interface{}{},
		TimeLineEventsWithDocs:    []interface{}{},
		TimelineEventsWithoutDocs: []interface{}{},
		TimelineDetails:           []map[string]interface{}{},
	}
}

func (t *TimeLine) SumTotalOfTimeLineEventsWDocsAndWithoutDocs() int {
	return len(t.TimeLineEventsWithDocs) + len(t.TimelineEventsWithoutDocs)
}

func (t *TimeLine) SetSinceTimestamp(sinceTimestamp float64) {
	t.SinceTimestamp = sinceTimestamp
}

func (t *TimeLine) GetTimeLineEventsWithDocs() []interface{} {
	return t.TimeLineEventsWithDocs
}

func (t *TimeLine) GetTimeLineEventsWithoutDocs() []interface{} {
	return t.TimelineEventsWithoutDocs
}

func (t *TimeLine) GetTimeLineEvents() []interface{} {
	return t.TimeLineEvents
}

func (t *TimeLine) GetTimeLineDetails() []map[string]interface{} {
	return t.TimelineDetails
}

func (t *TimeLine) LoadTimeLine(ctx context.Context, data chan Message) error {
	_, err := t.LoadNextTimeline(nil, t.SinceTimestamp)
	if err != nil {
		return err
	}
	for {
		select {
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
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (t *TimeLine) LoadNextTimeline(response map[string]interface{}, maxAgeTimestamp float64) (int, error) {
	if response == nil {
		logrus.Info("Awaiting #1 timeline")

		id, err := t.Client.AllTimeline()
		if err != nil {
			return -1, err
		}
		return id, nil

	} else {
		timestamp := response["data"].([]interface{})[len(response["data"].([]interface{}))-1].(map[string]interface{})["data"].(map[string]interface{})["timestamp"].(float64)
		t.NumberofTimline = t.NumberofTimline + 1
		t.NumberofTimlineDetails = t.NumberofTimlineDetails + len(response["data"].([]interface{}))
		logrus.Infof("Received timeline #%d with %d events", t.NumberofTimline, len(response["data"].([]interface{})))
		t.TimeLineEvents = append(t.TimeLineEvents, response["data"].([]interface{})...)

		after, ok := response["cursors"].(map[string]interface{})["after"].(string)
		if !ok {
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

	_, err := t.RequestTimeLineDetails(5, t.SinceTimestamp)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-data:
			if msg.Subscription["type"] == "timelineDetail" {
				t.ReceivedDetail++
				t.TimelineDetails = append(t.TimelineDetails, msg.Payload)
				if end, err := t.loadMoreTimeLineDetails(t.SinceTimestamp); err != nil {
					return err
				} else if end == -1 {
					return nil
				}
			} else {
				continue
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (t *TimeLine) RequestTimeLineDetails(numToRequest int, maxAgeTimestamp float64) (int, error) {
	for numToRequest > 0 {
		if len(t.TimeLineEvents) == t.SumTotalOfTimeLineEventsWDocsAndWithoutDocs() {
			logrus.Info("All timeline details requested")
			return -1, nil
		} else {
			event := t.TimeLineEvents[len(t.TimeLineEvents)-t.SumTotalOfTimeLineEventsWDocsAndWithoutDocs()-1].(map[string]interface{})

			if WithDocs := t.FindTimeLineEventsWithDocs(event, maxAgeTimestamp); WithDocs {

				numToRequest--
				t.RequestedDetail++
				id := event["data"].(map[string]interface{})["id"].(string)
				return t.Client.TimelineDetail(id)
			} else {
				continue
			}

		}
	}

	return -1, nil
}

func (t *TimeLine) FindTimeLineEventsWithDocs(event map[string]interface{}, maxAgeTimestamp float64) bool {
	if t.TimeLineEvents == nil {
		return false
	}

	action, ok := event["data"].(map[string]interface{})["action"].(map[string]interface{})

	var msg string
	if maxAgeTimestamp != 0 && event["data"].(map[string]interface{})["timestamp"].(float64) > maxAgeTimestamp {
		msg = "Skipping timeline detail %s %s"
	} else if !ok {
		if _, ok := event["data"].(map[string]interface{})["actionLabel"].(string); !ok {
			msg = "Skipping: no action"
		}
	} else if action["type"].(string) != "timelineDetail" {
		msg += fmt.Sprintf("Skipping timeline detail %s", action["type"].(string))
	}
	if msg == "" {
		t.TimeLineEventsWithDocs = append(t.TimeLineEventsWithDocs, event)
		return true
	} else {
		t.TimelineEventsWithoutDocs = append(t.TimelineEventsWithoutDocs, event)
		logrus.Debugf("%s %s %s", msg, event["data"].(map[string]interface{})["title"].(string), event["data"].(map[string]interface{})["body"].(string))
		t.NumberofTimlineDetails = t.NumberofTimlineDetails - 1
		return false

	}
}

func (t *TimeLine) loadMoreTimeLineDetails(maxAgeTimestamp float64) (int, error) {
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

func (t *TimeLine) isSavingsPlan(response map[string]interface{}) bool {
	if response["subtitleText"] == "Sparplan" {
		return true
	} else {
		for _, section := range response["sections"].([]interface{}) {
			if section.(map[string]interface{})["type"] == "actionButtons" {
				for _, button := range section.(map[string]interface{})["data"].([]interface{}) {
					if button.(map[string]interface{})["action"].(map[string]interface{})["type"] == "editSavingPlan" || button.(map[string]interface{})["action"].(map[string]interface{})["type"] == "deleteSavingPlan" {
						return true
					}
				}
			}
		}
	}
	return false
}

func (t *TimeLine) getSavingsPlanFMT(response map[string]interface{}, ifSavingPlan bool) string {
	if response["subtitleText"] != "Sparplan" && ifSavingPlan {
		return " -- SPARPLAN"
	} else {
		return ""
	}
}
