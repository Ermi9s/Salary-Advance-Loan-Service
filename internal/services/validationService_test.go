package services

import (
	"fmt"
	"testing"

	"salaryAdvance/internal/entity"
)

func TestProcessDataDetectsFaultyRecords(t *testing.T) {
	customers := make([]entity.BankCustomer, 0, 50)
	samples := make([]entity.InputCustomer, 0, 50)

	for i := 0; i < 50; i++ {
		account := fmt.Sprintf("100000000%03d", i)
		name := fmt.Sprintf("Customer %d", i)
		customers = append(customers, entity.BankCustomer{CustomerName: name, AccountNo: account})
		samples = append(samples, entity.InputCustomer{CustomerName: name, AccountNo: account})
	}

	// Intentional name mismatch
	samples[12].CustomerName = "Wrong Name"
	// Intentional account number format/type mismatch
	samples[37].AccountNo = "12AB"

	batches, verified, err := ProcessData(customers, samples)
	if err != nil {
		t.Fatalf("process data failed: %v", err)
	}

	if len(batches) != 5 {
		t.Fatalf("expected 5 batches, got %d", len(batches))
	}

	if len(verified) != 48 {
		t.Fatalf("expected 48 verified records, got %d", len(verified))
	}

	errorsFound := 0
	for _, batch := range batches {
		for _, record := range batch.Records {
			if !record.Verified {
				errorsFound++
			}
		}
	}
	if errorsFound != 2 {
		t.Fatalf("expected exactly 2 invalid records, got %d", errorsFound)
	}
}
