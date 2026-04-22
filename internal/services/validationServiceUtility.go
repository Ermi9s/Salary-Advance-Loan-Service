package services

import (
	"fmt"
	"regexp"
	"salaryAdvance/internal/entity"
	"sort"
	"strings"
	"sync"
)

var accountNumberFormat = regexp.MustCompile(`^\d{10,13}$`)

type ValidationUtility struct {
	CustomerMap map[string]entity.BankCustomer
	SampleData  []entity.InputCustomer
}

func (v *ValidationUtility) convertCustomerListToMap(customerList []entity.BankCustomer) {
	customerMap := make(map[string]entity.BankCustomer)
	for _, customer := range customerList {
		accountNo := strings.TrimSpace(customer.AccountNo)
		customer.CustomerName = strings.TrimSpace(customer.CustomerName)
		customerMap[accountNo] = customer
	}
	v.CustomerMap = customerMap
}

func normalizeName(name string) string {
	trimmed := strings.TrimSpace(strings.ToLower(name))
	return strings.Join(strings.Fields(trimmed), " ")
}

func isValidAccountFormat(accountNo string) bool {
	return accountNumberFormat.MatchString(strings.TrimSpace(accountNo))
}

func (v *ValidationUtility) compareData(record entity.InputCustomer, recordIndex int) entity.ValidationLogRecord {
	var validationLog entity.ValidationLogRecord
	validationLog.RecordIndex = recordIndex

	accountNumber := strings.TrimSpace(record.AccountNo)
	if !isValidAccountFormat(accountNumber) {
		validationLog.Verified = false
		validationLog.Errors = append(validationLog.Errors, "account number invalid format/type")
		return validationLog
	}

	customer, ok := v.CustomerMap[accountNumber]
	if !ok {
		validationLog.Verified = false
		validationLog.Errors = append(validationLog.Errors, "account number does not match")
		return validationLog
	}

	if normalizeName(record.CustomerName) != normalizeName(customer.CustomerName) {
		validationLog.Verified = false
		validationLog.Errors = append(validationLog.Errors, "name does not match")
		return validationLog
	}

	validationLog.Verified = true
	validationLog.NormalizedRecord = customer

	return validationLog
}

func (v *ValidationUtility) processBatch(batchID int, start int, end int) entity.BatchValidationLog {
	var batch entity.BatchValidationLog
	batch.BatchID = batchID

	for index := start; index < end; index++ {
		record := v.compareData(v.SampleData[index], index+1)
		batch.Records = append(batch.Records, record)
		if !record.Verified {
			batch.ContainsErrors = true
		}
	}

	if batch.ContainsErrors {
		batch.FailureReason = "contains one or more records with invalid name or account number"
	}

	return batch
}

func BuildValidationReport(batches []entity.BatchValidationLog) entity.ValidationReport {
	report := entity.ValidationReport{Batches: batches}

	for _, batch := range batches {
		for _, record := range batch.Records {
			report.TotalRecords++
			if record.Verified {
				report.VerifiedRecords++
			} else {
				report.UnverifiedRecords++
			}
		}
	}

	return report
}

func ProcessData(customerList []entity.BankCustomer, sampleData []entity.InputCustomer) ([]entity.BatchValidationLog, []entity.BankCustomer, error) {
	var vu ValidationUtility
	vu.convertCustomerListToMap(customerList)
	vu.SampleData = sampleData
	if len(vu.SampleData) == 0 {
		return nil, nil, fmt.Errorf("sample data is empty")
	}

	totalRecords := len(vu.SampleData)
	if totalRecords > 50 {
		totalRecords = 50
	}
	batchCount := (totalRecords + 9) / 10

	var wg sync.WaitGroup
	resultsChan := make(chan entity.BatchValidationLog, batchCount)

	// launch goroutines
	for batchID := 0; batchID < batchCount; batchID++ {
		start := 10 * batchID
		end := start + 10
		if end > totalRecords {
			end = totalRecords
		}

		wg.Add(1)

		go func(id int, start int, end int) {
			defer wg.Done()
			result := vu.processBatch(id, start, end)
			resultsChan <- result
		}(batchID, start, end)
	}

	// close channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []entity.BatchValidationLog
	for res := range resultsChan {
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].BatchID < results[j].BatchID
	})

	verifiedByAccount := make(map[string]entity.BankCustomer)
	for _, batch := range results {
		for _, record := range batch.Records {
			if record.Verified {
				verifiedByAccount[record.NormalizedRecord.AccountNo] = record.NormalizedRecord
			}
		}
	}

	verifiedCustomers := make([]entity.BankCustomer, 0, len(verifiedByAccount))
	for _, customer := range verifiedByAccount {
		verifiedCustomers = append(verifiedCustomers, customer)
	}

	return results, verifiedCustomers, nil
}
