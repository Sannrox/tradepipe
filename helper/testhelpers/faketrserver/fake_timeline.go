package testing

import (
	"fmt"
	"math/rand"

	"github.com/Sannrox/tradepipe/helper/testhelpers/random"
)

type FakeRawTimeline struct {
	Data    []TimeLineEvent `json:"data"`
	Cursors struct {
		Before string `json:"before,omitempty"`
		After  string `json:"after,omitempty"`
	} `json:"cursors"`
}

type TimeLineEvent struct {
	Type string `json:"type"`
	Data struct {
		ID          string        `json:"id"`
		Timestamp   int64         `json:"timestamp"`
		Icon        string        `json:"icon"`
		Title       string        `json:"title"`
		Body        string        `json:"body"`
		Action      Action        `json:"action"`
		ActionLabel string        `json:"actionLabel,omitempty"`
		Attributes  []interface{} `json:"attributes"`
		Month       string        `json:"month"`
	} `json:"data"`
}

type Action struct {
	Type    string `json:"type,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type TimelineDetail struct {
	ID           string `json:"id"`
	TitleText    string `json:"titleText"`
	SubtitleText string `json:"subtitleText"`
	Sections     []struct {
		Data      []interface{}
		Type      string `json:"type"`
		Title     string `json:"title"`
		Documents []Doc  `json:"documents"`
	} `json:"sections"`
}

type Doc struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Action struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	} `json:"action"`
	ID          string `json:"id"`
	PostboxType string `json:"postboxType"`
}

type FakeRawTimelines []FakeRawTimeline
type FakeTimelineDetails []TimelineDetail

func NewFakeTimelines() *FakeRawTimelines {
	return &FakeRawTimelines{}
}

func (t *FakeRawTimelines) GenerateRawTimelines(sets int) {
	for i := 0; i <= sets; i++ {
		timeline := FakeRawTimeline{}
		timeline.GenerateTimeline(10)
		t.Add(timeline)
	}
}

func (t *FakeRawTimelines) Add(timeline FakeRawTimeline) {
	*t = append(*t, timeline)
}

func (t *FakeRawTimelines) Len() int {
	return len(*t)
}

func (t *FakeRawTimelines) Last() FakeRawTimeline {
	return (*t)[t.Len()-1]
}

func (t *FakeRawTimelines) First() FakeRawTimeline {
	return (*t)[0]
}

func (t *FakeRawTimelines) Next(after string) FakeRawTimeline {
	for i, time := range *t {
		if time.Cursors.After == after {
			return (*t)[i+1]
		}
	}
	return FakeRawTimeline{}
}

func (t *FakeRawTimeline) GenerateTimeline(sets int) {
	timeline := FakeRawTimeline{}
	var typ string
	switch rand.Intn(2) {
	case 0:
		typ = "timeline"
	case 1:
		typ = "timelineDetail"
	}

	for j := 0; j < sets; j++ {
		timeline.Data = append(timeline.Data, TimeLineEvent{
			Type: typ,
			Data: struct {
				ID          string        `json:"id"`
				Timestamp   int64         `json:"timestamp"`
				Icon        string        `json:"icon"`
				Title       string        `json:"title"`
				Body        string        `json:"body"`
				Action      Action        `json:"action"`
				ActionLabel string        `json:"actionLabel,omitempty"`
				Attributes  []interface{} `json:"attributes"`
				Month       string        `json:"month"`
			}{
				ID:        random.GenerateRandomString(10),
				Timestamp: int64(j),
				Icon:      random.GenerateRandomString(10),
				Title:     random.GenerateRandomString(10),
				Body:      random.GenerateRandomString(10),
				Action: Action{
					Type:    "type",
					Payload: "payload",
				},
				ActionLabel: "actionLabel",
				Attributes:  []interface{}{},
				Month:       fmt.Sprintf("%d-%d", random.GenerateRandomYear(), random.GenerateRandomMonth()),
			},
		})

	}
	*t = timeline
}

func (t *FakeRawTimeline) Add(detail TimeLineEvent) {
	t.Data = append(t.Data, detail)
}

func (t *FakeRawTimeline) Get() FakeRawTimeline {
	return *t
}

func (t *FakeRawTimeline) Clear() {
	t.Data = []TimeLineEvent{}
}

func NewFakeTimelineDetails() *FakeTimelineDetails {
	return &FakeTimelineDetails{}
}
func (t *FakeTimelineDetails) GenerateTimelineDetailById(id string) TimelineDetail {
	for _, detail := range *t {
		if detail.ID == id {
			return detail
		}
	}
	return TimelineDetail{}
}

func (t *FakeTimelineDetails) GenerateTimelineDetail(event *TimeLineEvent) {

	var detail TimelineDetail
	if event.Type == "timelineDetail" {
		detail = t.GenerateTimelineDetailWithDoc(event.Data.ID, rand.Intn(5))
		t.Add(detail)
	} else {
		detail = t.GenerateTimelineDetailWithoutDoc(event.Data.ID)
		t.Add(detail)
	}

}

func (t *FakeTimelineDetails) GenerateTimelineDetailWithDoc(id string, sections int) TimelineDetail {
	detail := TimelineDetail{}
	detail.ID = id
	detail.TitleText = random.GenerateRandomString(10)
	detail.SubtitleText = random.GenerateRandomString(10)
	for i := 0; i < sections; i++ {
		detail.Sections = append(detail.Sections, struct {
			Data      []interface{}
			Type      string `json:"type"`
			Title     string `json:"title"`
			Documents []Doc  `json:"documents"`
		}{
			Data:      []interface{}{},
			Type:      "section",
			Title:     random.GenerateRandomString(10),
			Documents: []Doc{},
		})
		for j := 0; j < rand.Intn(5); j++ {
			detail.Sections[i].Documents = append(detail.Sections[i].Documents, Doc{
				Title:  random.GenerateRandomString(10),
				Detail: random.GenerateRandomString(10),
				Action: struct {
					Type    string `json:"type"`
					Payload string `json:"payload"`
				}{
					Type:    "action",
					Payload: random.GenerateRandomString(10),
				},
				ID:          random.GenerateRandomString(10),
				PostboxType: "postboxType",
			})
		}
	}
	return detail

}

func (t *FakeTimelineDetails) GenerateTimelineDetailWithoutDoc(id string) TimelineDetail {
	detail := TimelineDetail{}
	detail.ID = id
	detail.TitleText = random.GenerateRandomString(10)
	detail.SubtitleText = random.GenerateRandomString(10)
	for i := 0; i < rand.Intn(5); i++ {
		detail.Sections = append(detail.Sections, struct {
			Data      []interface{}
			Type      string `json:"type"`
			Title     string `json:"title"`
			Documents []Doc  `json:"documents"`
		}{
			Data:      []interface{}{},
			Type:      "section",
			Title:     random.GenerateRandomString(10),
			Documents: []Doc{},
		})
	}
	return detail
}

func (t *FakeTimelineDetails) Add(detail TimelineDetail) {
	*t = append(*t, detail)
}

func (t *FakeTimelineDetails) Len() int {
	return len(*t)
}
