package testing

type FakeTimeline []FakeTimelineDetail

type FakeTimelineDetail map[string]interface{}

func NewFakeTimeline() *FakeTimeline {
	return &FakeTimeline{}
}

func (t *FakeTimeline) GenerateTimeline() {
	// TODO: Generate timeline
}

func (t *FakeTimeline) GenerateTimelineDetail() {
}

func (t *FakeTimeline) Add(detail FakeTimelineDetail) {
	*t = append(*t, detail)
}

func (t *FakeTimeline) Get() []FakeTimelineDetail {
	return *t
}

func (t *FakeTimeline) Clear() {
	*t = []FakeTimelineDetail{}
}

func (t *FakeTimeline) Len() int {
	return len(*t)
}

func (t *FakeTimeline) Last() FakeTimelineDetail {
	return (*t)[t.Len()-1]
}

func (t *FakeTimeline) First() FakeTimelineDetail {
	return (*t)[0]
}

func (t *FakeTimeline) GetBy(key string, value interface{}) []FakeTimelineDetail {
	var result []FakeTimelineDetail
	for _, detail := range *t {
		if detail.Get(key) == value {
			result = append(result, detail)
		}
	}
	return result
}

func (t *FakeTimeline) GetByIndex(key string, value interface{}) int {
	for i, detail := range *t {
		if detail.Get(key) == value {
			return i
		}
	}
	return -1
}

func (td *FakeTimelineDetail) Get(key string) interface{} {
	return (*td)[key]
}

func (td *FakeTimelineDetail) Set(key string, value interface{}) {
	(*td)[key] = value
}

func NewFakeTimelineDetail() *FakeTimelineDetail {
	return &FakeTimelineDetail{}
}
