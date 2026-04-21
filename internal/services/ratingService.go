package services

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
	"salaryAdvance/internal/entity"
)

type RatingService struct {
	AllowOverdraft bool
}

func (s *RatingService) RateCustomers(customers []entity.BankCustomer, rawTransactions []entity.Transaction) ([]entity.CustomerRating, error) {
	if len(customers) == 0 {
		return nil, fmt.Errorf("no customers provided")
	}

	mapped, err := s.mapTransactions(customers, rawTransactions)
	if err != nil {
		return nil, err
	}

	ratings := make([]entity.CustomerRating, 0, len(customers))
	for _, customer := range customers {
		tx := mapped[customer.AccountNo]
		generated := 0
		if len(tx) == 0 {
			tx = s.generateSyntheticTransactions(customer)
			generated = len(tx)
		}

		sort.Slice(tx, func(i, j int) bool {
			return tx[i].TransactionDate.Before(tx[j].TransactionDate)
		})

		breakdown := calculateBreakdown(customer, tx)
		ratings = append(ratings, entity.CustomerRating{
			AccountNo:            customer.AccountNo,
			CustomerName:         customer.CustomerName,
			Rating:               clamp(breakdown.WeightedTotal, 1, 10),
			GeneratedTxCount:     generated,
			TransactionCount:     len(tx),
			TotalVolume:          totalVolume(tx),
			CalculationBreakdown: breakdown,
		})
	}

	return ratings, nil
}

func (s *RatingService) mapTransactions(customers []entity.BankCustomer, raw []entity.Transaction) (map[string][]entity.EnrichedTransaction, error) {
	accountSet := make(map[string]struct{}, len(customers))
	for _, customer := range customers {
		accountSet[customer.AccountNo] = struct{}{}
	}

	result := make(map[string][]entity.EnrichedTransaction)
	for _, tx := range raw {
		amount, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			continue
		}

		timeMs, err := strconv.ParseInt(tx.TransactionDate, 10, 64)
		if err != nil {
			continue
		}
		timeValue := time.UnixMilli(timeMs).UTC()

		if _, ok := accountSet[tx.FromAccount]; ok {
			result[tx.FromAccount] = append(result[tx.FromAccount], entity.EnrichedTransaction{
				ID:              tx.ID,
				AccountNo:       tx.FromAccount,
				Amount:          amount,
				Direction:       "debit",
				TransactionDate: timeValue,
				Source:          "source_file",
			})
		}

		if tx.ToAccount != nil {
			if _, ok := accountSet[*tx.ToAccount]; ok {
				result[*tx.ToAccount] = append(result[*tx.ToAccount], entity.EnrichedTransaction{
					ID:              tx.ID,
					AccountNo:       *tx.ToAccount,
					Amount:          amount,
					Direction:       "credit",
					TransactionDate: timeValue,
					Source:          "source_file",
				})
			}
		}
	}

	return result, nil
}

func (s *RatingService) generateSyntheticTransactions(customer entity.BankCustomer) []entity.EnrichedTransaction {
	seed := seededFloat(customer.AccountNo)
	base := customer.CustomerBalance
	if base <= 0 {
		base = 1000
	}

	items := make([]entity.EnrichedTransaction, 0, 4)
	current := base
	start := time.Now().UTC().AddDate(0, -4, 0)
	for i := 0; i < 4; i++ {
		debit := (0.08 + seed*0.06) * base
		credit := (0.05 + seed*0.04) * base
		if i%2 == 0 {
			amount := math.Min(debit, max(current, 50))
			current -= amount
			if current < 0 && !s.AllowOverdraft {
				amount += current
				current = 0
			}
			items = append(items, entity.EnrichedTransaction{
				ID:              fmt.Sprintf("syn-%s-%d", customer.AccountNo, i),
				AccountNo:       customer.AccountNo,
				Amount:          amount,
				Direction:       "debit",
				TransactionDate: start.AddDate(0, i, 0),
				Source:          "synthetic",
			})
		} else {
			current += credit
			items = append(items, entity.EnrichedTransaction{
				ID:              fmt.Sprintf("syn-%s-%d", customer.AccountNo, i),
				AccountNo:       customer.AccountNo,
				Amount:          credit,
				Direction:       "credit",
				TransactionDate: start.AddDate(0, i, 0),
				Source:          "synthetic",
			})
		}
	}

	return items
}

