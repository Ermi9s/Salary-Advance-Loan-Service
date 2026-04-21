package entity

type ValidationLogRecord struct {
	RecordIndex      int          `json:"record_index"`
	Verified         bool         `json:"verified"`
	Errors           []string     `json:"errors,omitempty"`
	NormalizedRecord BankCustomer `json:"normalized_record,omitempty"`
}

type BatchValidationLog struct {
	BatchID        int                   `json:"batch_id"`
	Records        []ValidationLogRecord `json:"records"`
	ContainsErrors bool                  `json:"contains_errors"`
	FailureReason  string                `json:"failure_reason,omitempty"`
}

type ValidationReport struct {
	ProcessedAt       string               `json:"processed_at"`
	TotalRecords      int                  `json:"total_records"`
	VerifiedRecords   int                  `json:"verified_records"`
	UnverifiedRecords int                  `json:"unverified_records"`
	Batches           []BatchValidationLog `json:"batches"`
}
