package testing

import "github.com/Sannrox/tradepipe/helper/testhelpers/random"

type FakeSavingsplans struct {
	SavingsPlans []struct {
		ID           string  `json:"id"`
		CreatedAt    int64   `json:"createdAt"`
		InstrumentID string  `json:"instrumentId"`
		Amount       float64 `json:"amount"`
		Interval     string  `json:"interval"`
		StartDate    struct {
			Type              string `json:"type"`
			Value             int    `json:"value"`
			NextExecutionDate string `json:"nextExecutionDate"`
		} `json:"startDate"`
		FirstExecutionDate           interface{} `json:"firstExecutionDate"`
		NextExecutionDate            string      `json:"nextExecutionDate"`
		PreviousExecutionDate        string      `json:"previousExecutionDate"`
		VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
		FinalExecutionDate           interface{} `json:"finalExecutionDate"`
		PaymentMethodID              interface{} `json:"paymentMethodId"`
		PaymentMethodCode            interface{} `json:"paymentMethodCode"`
		LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
		Paused                       bool        `json:"paused"`
	} `json:"savingsPlans"`
}

func NewFakeSavingsplans() *FakeSavingsplans {
	return &FakeSavingsplans{}
}

func (s *FakeSavingsplans) GenerateFakeSavingsPlans(sets int) {
	var savingsPlan struct {
		ID           string  `json:"id"`
		CreatedAt    int64   `json:"createdAt"`
		InstrumentID string  `json:"instrumentId"`
		Amount       float64 `json:"amount"`
		Interval     string  `json:"interval"`
		StartDate    struct {
			Type              string `json:"type"`
			Value             int    `json:"value"`
			NextExecutionDate string `json:"nextExecutionDate"`
		} `json:"startDate"`
		FirstExecutionDate           interface{} `json:"firstExecutionDate"`
		NextExecutionDate            string      `json:"nextExecutionDate"`
		PreviousExecutionDate        string      `json:"previousExecutionDate"`
		VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
		FinalExecutionDate           interface{} `json:"finalExecutionDate"`
		PaymentMethodID              interface{} `json:"paymentMethodId"`
		PaymentMethodCode            interface{} `json:"paymentMethodCode"`
		LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
		Paused                       bool        `json:"paused"`
	}

	for i := 0; i < sets; i++ {
		savingsPlan = s.GenerateFakeSavingsPlan()
	}
	if len(s.SavingsPlans) > 0 {
		s.SavingsPlans = append(s.SavingsPlans, savingsPlan)
	} else {
		s.SavingsPlans = []struct {
			ID           string  `json:"id"`
			CreatedAt    int64   `json:"createdAt"`
			InstrumentID string  `json:"instrumentId"`
			Amount       float64 `json:"amount"`
			Interval     string  `json:"interval"`
			StartDate    struct {
				Type              string `json:"type"`
				Value             int    `json:"value"`
				NextExecutionDate string `json:"nextExecutionDate"`
			} `json:"startDate"`
			FirstExecutionDate           interface{} `json:"firstExecutionDate"`
			NextExecutionDate            string      `json:"nextExecutionDate"`
			PreviousExecutionDate        string      `json:"previousExecutionDate"`
			VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
			FinalExecutionDate           interface{} `json:"finalExecutionDate"`
			PaymentMethodID              interface{} `json:"paymentMethodId"`
			PaymentMethodCode            interface{} `json:"paymentMethodCode"`
			LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
			Paused                       bool        `json:"paused"`
		}{savingsPlan}
	}

}

func (s *FakeSavingsplans) GenerateFakeSavingsPlan() struct {
	ID           string  `json:"id"`
	CreatedAt    int64   `json:"createdAt"`
	InstrumentID string  `json:"instrumentId"`
	Amount       float64 `json:"amount"`
	Interval     string  `json:"interval"`
	StartDate    struct {
		Type              string `json:"type"`
		Value             int    `json:"value"`
		NextExecutionDate string `json:"nextExecutionDate"`
	} `json:"startDate"`
	FirstExecutionDate           interface{} `json:"firstExecutionDate"`
	NextExecutionDate            string      `json:"nextExecutionDate"`
	PreviousExecutionDate        string      `json:"previousExecutionDate"`
	VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
	FinalExecutionDate           interface{} `json:"finalExecutionDate"`
	PaymentMethodID              interface{} `json:"paymentMethodId"`
	PaymentMethodCode            interface{} `json:"paymentMethodCode"`
	LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
	Paused                       bool        `json:"paused"`
} {
	return struct {
		ID           string  `json:"id"`
		CreatedAt    int64   `json:"createdAt"`
		InstrumentID string  `json:"instrumentId"`
		Amount       float64 `json:"amount"`
		Interval     string  `json:"interval"`
		StartDate    struct {
			Type              string `json:"type"`
			Value             int    `json:"value"`
			NextExecutionDate string `json:"nextExecutionDate"`
		} `json:"startDate"`
		FirstExecutionDate           interface{} `json:"firstExecutionDate"`
		NextExecutionDate            string      `json:"nextExecutionDate"`
		PreviousExecutionDate        string      `json:"previousExecutionDate"`
		VirtualPreviousExecutionDate string      `json:"virtualPreviousExecutionDate"`
		FinalExecutionDate           interface{} `json:"finalExecutionDate"`
		PaymentMethodID              interface{} `json:"paymentMethodId"`
		PaymentMethodCode            interface{} `json:"paymentMethodCode"`
		LastPaymentExecutionDate     interface{} `json:"lastPaymentExecutionDate"`
		Paused                       bool        `json:"paused"`
	}{
		ID:           "fake-savingsplan-id",
		CreatedAt:    1610000000,
		InstrumentID: random.GenerateRandomeISIN(),
		Amount:       random.GenerateRandomAmount(),
		Interval:     random.GenerateRandomInterval(),
		StartDate: struct {
			Type              string `json:"type"`
			Value             int    `json:"value"`
			NextExecutionDate string `json:"nextExecutionDate"`
		}{
			Type:              "dayOfMonth",
			Value:             1,
			NextExecutionDate: "2021-01-01",
		},
		FirstExecutionDate:           nil,
		NextExecutionDate:            "2021-01-01",
		PreviousExecutionDate:        "2021-01-01",
		VirtualPreviousExecutionDate: "2021-01-01",
		FinalExecutionDate:           nil,
		PaymentMethodID:              nil,
		PaymentMethodCode:            nil,
		LastPaymentExecutionDate:     nil,
		Paused:                       false,
	}
}

func (s *FakeSavingsplans) GetSavingsPlans() *FakeSavingsplans {
	return s
}
