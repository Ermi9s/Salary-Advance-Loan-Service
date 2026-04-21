package dto

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"salaryAdvance/internal/entity"
	"strconv"
	"strings"
)

type DTOMethodes struct {
	CustomersFilePath      string
	TransactionFilePath    string
	SampleCustomerFilePath string
}

func readJSONFile(path string, target any) error {
	openedFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = openedFile.Close()
	}()

	decoder := json.NewDecoder(openedFile)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}

func (d *DTOMethodes) ReadCustomerData() ([]entity.BankCustomer, error) {
	var records []entity.BankCustomer
	if err := readJSONFile(d.CustomersFilePath, &records); err != nil {
		log.Printf("Error Reading data from Json %v", err)
		return nil, err
	}
	return records, nil
}

func (d *DTOMethodes) ReadTransactionData() ([]entity.Transaction, error) {
	var records []entity.Transaction
	if err := readJSONFile(d.TransactionFilePath, &records); err != nil {
		log.Printf("Error Reading data from Json %v", err)
		return nil, err
	}
	return records, nil
}

func (d *DTOMethodes) ReadSampleData() ([]entity.InputCustomer, error) {
	ext := strings.ToLower(filepath.Ext(d.SampleCustomerFilePath))
	if ext == ".json" {
		return d.readSampleJSON()
	}
	return d.ReadCSV()
}

func (d *DTOMethodes) readSampleJSON() ([]entity.InputCustomer, error) {
	var rawRecords []map[string]any
	if err := readJSONFile(d.SampleCustomerFilePath, &rawRecords); err != nil {
		return nil, err
	}

	customers := make([]entity.InputCustomer, 0, len(rawRecords))
	for _, raw := range rawRecords {
		name := strings.TrimSpace(fmt.Sprintf("%v", raw["name"]))

		var account string
		switch v := raw["account_number"].(type) {
		case string:
			account = strings.TrimSpace(v)
		case float64:
			account = strconv.FormatInt(int64(v), 10)
		default:
			account = strings.TrimSpace(fmt.Sprintf("%v", v))
		}

		customers = append(customers, entity.InputCustomer{
			CustomerName: name,
			AccountNo:    account,
		})
	}

	return customers, nil
}

func (d *DTOMethodes) ReadCSV() ([]entity.InputCustomer, error) {
	openedFile, err := os.Open(d.SampleCustomerFilePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := openedFile.Close(); closeErr != nil {
			log.Printf("Error closing file: %v", closeErr)
		}
	}()

	reader := csv.NewReader(openedFile)
	var customers []entity.InputCustomer

	// skip header
	if _, err := reader.Read(); err != nil {
		log.Printf("Error reading header: %v", err)
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			return nil, err
		}

		// safety check
		if len(record) < 2 {
			log.Printf("Invalid record: %v", record)
			continue
		}

		customer := entity.InputCustomer{
			CustomerName: record[0],
			AccountNo:    record[1],
		}
		customers = append(customers, customer)
	}

	return customers, nil
}
