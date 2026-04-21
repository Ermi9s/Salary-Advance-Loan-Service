package services

import (
	"testing"

	"salaryAdvance/internal/entity"
)

func TestRateCustomersReturnsBoundedScore(t *testing.T) {
	svc := &RatingService{AllowOverdraft: false}
	customers := []entity.BankCustomer{{
		CustomerName:    "A User",
		AccountNo:       "1050001035901",
		CustomerBalance: 10000,
	}}

	transactions := []entity.Transaction{
		{
			ID:              "tx-1",
			FromAccount:     "1050001035901",
			Amount:          "1500.00",
			TransactionDate: "1733741248663",
		},
		{
			ID:              "tx-2",
			FromAccount:     "1050001035901",
			Amount:          "500.00",
			TransactionDate: "1734741248663",
		},
	}

	ratings, err := svc.RateCustomers(customers, transactions)
	if err != nil {
		t.Fatalf("rate failed: %v", err)
	}
	if len(ratings) != 1 {
		t.Fatalf("expected 1 rating, got %d", len(ratings))
	}
	if ratings[0].Rating < 1 || ratings[0].Rating > 10 {
		t.Fatalf("rating out of bounds: %v", ratings[0].Rating)
	}
	if ratings[0].CalculationBreakdown.WeightedTotal <= 0 {
		t.Fatalf("expected positive weighted total")
	}
}
