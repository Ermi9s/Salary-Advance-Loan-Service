package services

import (
	"math"
	"hash/fnv"
	"salaryAdvance/internal/entity"
)

func calculateBreakdown(customer entity.BankCustomer, tx []entity.EnrichedTransaction) entity.RatingBreakdown {
	countScore := clamp(float64(len(tx))/25*10, 0, 10)
	volume := totalVolume(tx)
	volumeScore := clamp(volume/max(customer.CustomerBalance, 1)*5, 0, 10)

	durationScore := 0.0
	if len(tx) > 1 {
		days := tx[len(tx)-1].TransactionDate.Sub(tx[0].TransactionDate).Hours() / 24
		durationScore = clamp(days/365*10, 0, 10)
	}

	stabilityScore := stabilityFromTransactions(customer.CustomerBalance, tx)
	weighted := 0.25*countScore + 0.30*volumeScore + 0.20*durationScore + 0.25*stabilityScore

	return entity.RatingBreakdown{
		CountScore:     countScore,
		VolumeScore:    volumeScore,
		DurationScore:  durationScore,
		StabilityScore: stabilityScore,
		WeightedTotal:  weighted,
	}
}

func stabilityFromTransactions(startBalance float64, tx []entity.EnrichedTransaction) float64 {
	if len(tx) == 0 {
		return 0
	}
	balances := make([]float64, 0, len(tx)+1)
	current := startBalance
	balances = append(balances, current)
	for _, item := range tx {
		if item.Direction == "debit" {
			current -= item.Amount
		} else {
			current += item.Amount
		}
		if current < 0 {
			current = 0
		}
		balances = append(balances, current)
	}

	mean := 0.0
	for _, b := range balances {
		mean += b
	}
	mean /= float64(len(balances))
	if mean == 0 {
		return 0
	}

	variance := 0.0
	for _, b := range balances {
		delta := b - mean
		variance += delta * delta
	}
	variance /= float64(len(balances))
	std := math.Sqrt(variance)
	cv := std / mean
	return clamp(10-(cv*10), 0, 10)
}

func totalVolume(tx []entity.EnrichedTransaction) float64 {
	total := 0.0
	for _, item := range tx {
		total += item.Amount
	}
	return total
}

func clamp(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func seededFloat(account string) float64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(account))
	return float64(h.Sum32()%1000) / 1000
}
