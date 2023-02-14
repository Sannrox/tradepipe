package testing

import (
	"fmt"
	"strconv"
	"time"

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
		timeline.GenerateTimeline(5)
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

func NewFakeTimelineDetail() *FakeTimelineDetails {
	return &FakeTimelineDetails{}
}

func (t *FakeRawTimeline) GenerateTimeline(sets int) {
	timeline := FakeRawTimeline{}
	for j := 0; j < sets; j++ {
		timeline.Data = append(timeline.Data, TimeLineEvent{
			Type: "timeline_event",
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
				ID:          random.GenerateRandomString(10),
				Timestamp:   int64(j),
				Icon:        random.GenerateRandomString(10),
				Title:       random.GenerateRandomString(10),
				Body:        random.GenerateRandomString(10),
				Action:      Action{},
				ActionLabel: "actionLabel",
				Attributes:  []interface{}{},
				Month:       fmt.Sprintf("%d-%d", random.GenerateRandomYear(), random.GenerateRandomMonth()),
			},
		})

	}
	timeline.Cursors.After = strconv.FormatInt(time.Now().Unix(), 36)

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

func (t *FakeTimelineDetails) GenerateTimelineDetail() {
}

func (t *FakeTimelineDetails) GenerateTimelineDetailWithDoc() {

}

func (t *FakeTimelineDetails) GenerateTimelineDetailWithDocAndAction() {

}
