package tr

import (
	"context"
	"encoding/json"
	"fmt"
)

type SavingsPlaner struct {
	Client *APIClient
	Plans  []SavingsPlan
}

func NewSavingsPlan(client *APIClient) *SavingsPlaner {
	return &SavingsPlaner{Client: client, Plans: []SavingsPlan{}}
}

func (s *SavingsPlaner) LoadSavingsplans(ctx context.Context, data chan Message) error {
	_, err := s.Client.SavingsPlanOverview()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case msg := <-data:
			if msg.Subscription["type"] == "savingsPlans" {
				var savingsPlan RawSavingsPlans

				b, err := json.Marshal(msg.Payload)
				if err != nil {
					return fmt.Errorf("%w | %s", err, string(b))
				}
				err = json.Unmarshal(b, &savingsPlan)
				if err != nil {
					return fmt.Errorf("%w | %s", err, string(b))
				}
				for _, plan := range savingsPlan.SavingsPlans {
					s.Plans = append(s.Plans, plan)
				}

				return nil
			}
		}
	}
}

func (s *SavingsPlaner) GetSavingsPlans() []SavingsPlan {
	return s.Plans
}

func (s *SavingsPlaner) GetSavingsPlansAsBytes() ([]byte, error) {
	return json.Marshal(s.Plans)
}

func (s *SavingsPlaner) GetSavingsPlan(id string) (SavingsPlan, error) {
	for _, plan := range s.Plans {
		if plan.ID == id {
			return plan, nil
		}
	}
	return SavingsPlan{}, fmt.Errorf("no plan found with id %s", id)
}

func (s *SavingsPlaner) GetAllSavingPlansIDs() []string {
	var ids []string
	for _, plan := range s.Plans {
		ids = append(ids, plan.ID)
	}
	return ids
}

func (s *SavingsPlaner) GetSavingsPlanAsBytes(id string) ([]byte, error) {
	plan, err := s.GetSavingsPlan(id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(plan)
}

func (s *SavingsPlaner) GetSavingsPlanAmount(id string) (float64, error) {
	plan, err := s.GetSavingsPlan(id)
	if err != nil {
		return 0, err
	}
	return plan.Amount, nil
}

func (s *SavingsPlaner) GetAddedAmountForMonth() (float64, error) {
	var amount float64
	for _, plan := range s.Plans {
		amountForMonth, err := s.GetAmountForMonth(plan.ID)
		if err != nil {
			return 0, err
		}
		amount += amountForMonth
	}
	return amount, nil
}

func (s *SavingsPlaner) GetAmountForMonth(id string) (float64, error) {
	plan, err := s.GetSavingsPlan(id)
	var multiplikator float64
	if err != nil {
		return 0, err
	}
	switch {
	case plan.Interval == "weekly":
		multiplikator = 4.0
	case plan.Interval == "monthly":
		multiplikator = 1.0
	case plan.Interval == "quarterly":
		multiplikator = 1.0 / 3.0
	case plan.Interval == "halfYearly":
		multiplikator = 1.0 / 6.0
	case plan.Interval == "yearly":
		multiplikator = 1.0 / 12.0

	}
	return float64(multiplikator) * plan.Amount, nil
}
