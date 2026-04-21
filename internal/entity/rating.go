package entity

type RatingBreakdown struct {
	CountScore     float64 `json:"count_score"`
	VolumeScore    float64 `json:"volume_score"`
	DurationScore  float64 `json:"duration_score"`
	StabilityScore float64 `json:"stability_score"`
	WeightedTotal  float64 `json:"weighted_total"`
}

type CustomerRating struct {
	AccountNo            string          `json:"account_no"`
	CustomerName         string          `json:"customer_name"`
	Rating               float64         `json:"rating"`
	GeneratedTxCount     int             `json:"generated_transaction_count"`
	TransactionCount     int             `json:"transaction_count"`
	TotalVolume          float64         `json:"total_volume"`
	CalculationBreakdown RatingBreakdown `json:"calculation_breakdown"`
}
