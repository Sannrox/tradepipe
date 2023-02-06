package testing

type FakeSubscriptionStore []int

func NewFakeSubscriptionStore() *FakeSubscriptionStore {
	return &FakeSubscriptionStore{}

}

func (s *FakeSubscriptionStore) AddSubscription(subscriptionId int) {
	*s = append(*s, subscriptionId)
}

func (s *FakeSubscriptionStore) Remove(subscriptionId int) {
	for i, v := range *s {
		if v == subscriptionId {
			*s = append((*s)[:i], (*s)[i+1:]...)
			break
		}
	}
}

func (s *FakeSubscriptionStore) GetSubscriptions() []int {
	return *s
}

func (s *FakeSubscriptionStore) GetSubscriptionCount() int {
	return len(*s)
}

func (s *FakeSubscriptionStore) Clear() {
	*s = []int{}
}

func (s *FakeSubscriptionStore) Contains(subscriptionId int) bool {
	for _, v := range *s {
		if v == subscriptionId {
			return true
		}
	}
	return false
}
